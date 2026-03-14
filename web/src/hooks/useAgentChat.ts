import { useState, useCallback, useRef, useEffect } from 'react';
import {
  Message,
  AssistantMessage,
  ToolMessage,
  AgentEvent,
  EventTypeEnum,
  MessageRoleEnum,
  MessageStatusEnum,
  ToolCallStatusEnum,
} from '../types';
import { SSEClient } from '../utils/sse-client';

interface HistoryToolCall {
  id: string;
  name: string;
  arguments: string;
}

interface HistoryMessage {
  id: string;
  role: 'user' | 'assistant' | 'tool';
  content: string;
  reasoningContent?: string;
  toolCallId?: string;
  toolCalls?: HistoryToolCall[];
}

interface UseAgentChatOptions {
  baseUrl?: string;
  existingSessionId: string;
}

export function useAgentChat(options: UseAgentChatOptions) {
  const { baseUrl = '/api/v1/agent', existingSessionId } = options;
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingHistory, setIsLoadingHistory] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [sessionModelId, setSessionModelId] = useState<number | null>(null);
  const [totalTokens, setTotalTokens] = useState<number>(0);
  const sessionId = existingSessionId;
  const sseClientRef = useRef<SSEClient | null>(null);
  const historyLoadedRef = useRef(false);

  const handleEvent = useCallback((event: AgentEvent) => {
    switch (event.type) {
      case EventTypeEnum.RUN_STARTED:
        setIsLoading(true);
        setError(null);
        break;

      case EventTypeEnum.RUN_FINISHED:
        setIsLoading(false);
        setMessages((prev) =>
          prev.map((msg) =>
            msg.status !== MessageStatusEnum.Finished
              ? { ...msg, status: MessageStatusEnum.Finished }
              : msg
          )
        );
        break;

      case EventTypeEnum.RUN_ERROR:
        setIsLoading(false);
        setError(event.message);
        break;

      case EventTypeEnum.TEXT_MESSAGE_START: {
        const isThinking = event.external?.isThinking ?? false;
        const newMessage: AssistantMessage = {
          id: event.messageId,
          role: MessageRoleEnum.ASSISTANT,
          content: '',
          createdAt: event.timestamp,
          status: MessageStatusEnum.Started,
          name: event.name,
          thinkingContent: '',
          isThinking,
        };
        setMessages((prev) => {
          const existingIndex = prev.findIndex((m) => m.id === event.messageId);
          if (existingIndex >= 0) {
            const updated = [...prev];
            updated[existingIndex] = {
              ...updated[existingIndex],
              ...newMessage,
            };
            return updated;
          }
          return [...prev, newMessage];
        });
        break;
      }

      case EventTypeEnum.TEXT_MESSAGE_CONTENT: {
        const isThinking = event.external?.isThinking ?? false;
        setMessages((prev) =>
          prev.map((msg) => {
            if (msg.id === event.messageId && msg.role === MessageRoleEnum.ASSISTANT) {
              const assistantMsg = msg as AssistantMessage;
              if (isThinking) {
                return {
                  ...assistantMsg,
                  thinkingContent: (assistantMsg.thinkingContent || '') + event.delta,
                  isThinking: true,
                };
              }
              return {
                ...assistantMsg,
                content: assistantMsg.content + event.delta,
                isThinking: false,
              };
            }
            return msg;
          })
        );
        break;
      }

      case EventTypeEnum.TEXT_MESSAGE_END: {
        setMessages((prev) =>
          prev.map((msg) =>
            msg.id === event.messageId
              ? { ...msg, status: MessageStatusEnum.Finished, isThinking: false }
              : msg
          )
        );
        break;
      }

      case EventTypeEnum.TOOL_CALL_START: {
        const toolMessage: ToolMessage = {
          id: event.toolCallId,
          role: MessageRoleEnum.TOOL,
          content: '',
          createdAt: event.timestamp,
          status: MessageStatusEnum.Started,
          toolCallId: event.toolCallId,
          toolCallName: event.toolCallName,
          toolCallArgs: '',
          toolCallStatus: event.toolCallStatus || ToolCallStatusEnum.Processing,
          parentMessageId: event.parentMessageId,
        };
        setMessages((prev) => [...prev, toolMessage]);
        break;
      }

      case EventTypeEnum.TOOL_CALL_ARGS: {
        setMessages((prev) =>
          prev.map((msg) => {
            if (msg.id === event.toolCallId && msg.role === MessageRoleEnum.TOOL) {
              const toolMsg = msg as ToolMessage;
              return {
                ...toolMsg,
                toolCallArgs: toolMsg.toolCallArgs + event.delta,
                toolCallStatus: event.toolCallStatus || toolMsg.toolCallStatus,
              };
            }
            return msg;
          })
        );
        break;
      }

      case EventTypeEnum.TOOL_CALL_RESULT: {
        setMessages((prev) =>
          prev.map((msg) => {
            if (msg.id === event.toolCallId && msg.role === MessageRoleEnum.TOOL) {
              const toolMsg = msg as ToolMessage;
              return {
                ...toolMsg,
                content: toolMsg.content + event.delta,
                toolCallStatus: event.toolCallStatus,
              };
            }
            return msg;
          })
        );
        break;
      }

      case EventTypeEnum.TOOL_CALL_END: {
        setMessages((prev) =>
          prev.map((msg) => {
            if (msg.id === event.toolCallId && msg.role === MessageRoleEnum.TOOL) {
              return {
                ...msg,
                status: MessageStatusEnum.Finished,
                toolCallStatus: event.toolCallStatus || ToolCallStatusEnum.Success,
                finishedAt: event.timestamp,
              };
            }
            return msg;
          })
        );
        break;
      }

      case EventTypeEnum.CUSTOM: {
        if (event.name === 'token_usage') {
          const usage = event.value as { totalTokens?: number };
          if (usage?.totalTokens) {
            setTotalTokens(usage.totalTokens);
          }
        }
        break;
      }

      default:
        break;
    }
  }, []);

  const sendMessage = useCallback(
    async (content: string, modelId?: number | null) => {
      console.log('[sendMessage] called, content:', content, 'modelId:', modelId, 'sessionId:', sessionId);
      historyLoadedRef.current = true;

      const userMessageId = crypto.randomUUID?.() ?? `${Date.now()}-${Math.random().toString(36).slice(2, 11)}`;
      const userMessage: Message = {
        id: userMessageId,
        role: MessageRoleEnum.USER,
        content,
        createdAt: Date.now(),
        status: MessageStatusEnum.Finished,
      };

      setMessages((prev) => [...prev, userMessage]);
      setIsLoading(true);
      setError(null);

      try {
        sseClientRef.current = new SSEClient();
        const payload: Record<string, unknown> = {
          sessionId,
          ...(modelId ? { modelId } : {}),
          messages: [
            {
              id: userMessageId,
              type: 'text',
              content,
              timestamp: Date.now(),
              role: 'user',
            },
          ],
          context: {
            timestamp: Date.now(),
            timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
          },
        };
        const body = JSON.stringify(payload);
        console.log('[sendMessage] connecting SSE, url:', `${baseUrl}/run`);

        await sseClientRef.current.connect({
          url: `${baseUrl}/run`,
          body,
          onEvent: handleEvent,
          onError: (err) => {
            console.error('[sendMessage] SSE onError:', err);
            setError(err.message);
            setIsLoading(false);
          },
          onComplete: () => {
            console.log('[sendMessage] SSE onComplete');
            setIsLoading(false);
          },
        });
        console.log('[sendMessage] connect returned');
      } catch (err) {
        console.error('[sendMessage] unexpected error:', err);
        setError(err instanceof Error ? err.message : 'Failed to send message');
        setIsLoading(false);
      }
    },
    [baseUrl, handleEvent, sessionId]
  );

  const stopGeneration = useCallback(() => {
    sseClientRef.current?.disconnect();
    setIsLoading(false);
  }, []);

  const clearMessages = useCallback(() => {
    setMessages([]);
    setError(null);
  }, []);

  const loadHistory = useCallback(async () => {
    console.log('[loadHistory] called, sessionId:', sessionId, 'historyLoaded:', historyLoadedRef.current);
    if (!sessionId || historyLoadedRef.current) {
      return;
    }

    setIsLoadingHistory(true);
    try {
      console.log('[loadHistory] fetching messages...');
      const response = await fetch(`/api/v1/sessions/${sessionId}/messages`);
      console.log('[loadHistory] response status:', response.status);
      if (!response.ok) {
        if (response.status === 404) {
          historyLoadedRef.current = true;
          return;
        }
        const errData = await response.json().catch(() => null);
        const errMsg = errData?.error || `Failed to load history (HTTP ${response.status})`;
        throw new Error(errMsg);
      }

      const data = await response.json();
      const modelId = Number(data.modelId);
      setSessionModelId(Number.isFinite(modelId) && modelId > 0 ? modelId : null);
      if (data.tokenUsage?.totalTokens) {
        setTotalTokens(data.tokenUsage.totalTokens);
      }
      const historyMessages: HistoryMessage[] = data.messages || [];

      const convertedMessages: Message[] = [];
      const now = Date.now();

      for (const msg of historyMessages) {
        if (msg.role === 'user') {
          convertedMessages.push({
            id: msg.id,
            role: MessageRoleEnum.USER,
            content: msg.content,
            createdAt: now,
            status: MessageStatusEnum.Finished,
          });
        } else if (msg.role === 'assistant') {
          convertedMessages.push({
            id: msg.id,
            role: MessageRoleEnum.ASSISTANT,
            content: msg.content,
            createdAt: now,
            status: MessageStatusEnum.Finished,
            thinkingContent: msg.reasoningContent || '',
            isThinking: false,
          });

          if (msg.toolCalls && msg.toolCalls.length > 0) {
            for (const tc of msg.toolCalls) {
              convertedMessages.push({
                id: tc.id,
                role: MessageRoleEnum.TOOL,
                content: '',
                createdAt: now,
                status: MessageStatusEnum.Started,
                toolCallId: tc.id,
                toolCallName: tc.name,
                toolCallArgs: tc.arguments,
                toolCallStatus: ToolCallStatusEnum.Processing,
                parentMessageId: msg.id,
              });
            }
          }
        } else if (msg.role === 'tool') {
          const existingToolIndex = convertedMessages.findIndex(
            (m) => m.role === MessageRoleEnum.TOOL && (m as ToolMessage).toolCallId === msg.toolCallId
          );
          if (existingToolIndex >= 0) {
            const existingTool = convertedMessages[existingToolIndex] as ToolMessage;
            convertedMessages[existingToolIndex] = {
              ...existingTool,
              content: msg.content,
              status: MessageStatusEnum.Finished,
              toolCallStatus: ToolCallStatusEnum.Success,
            };
          }
        }
      }

      if (convertedMessages.length > 0) {
        setMessages(convertedMessages);
      }
      historyLoadedRef.current = true;
    } catch (err) {
      console.error('Failed to load history:', err);
      setError(err instanceof Error ? err.message : 'Failed to load history');
    } finally {
      historyLoadedRef.current = true;
      setIsLoadingHistory(false);
    }
  }, [sessionId]);

  useEffect(() => {
    if (sessionId) {
      loadHistory();
    }
  }, [sessionId, loadHistory]);

  return {
    messages,
    isLoading,
    isLoadingHistory,
    error,
    sessionId,
    sessionModelId,
    totalTokens,
    sendMessage,
    stopGeneration,
    clearMessages,
    loadHistory,
  };
}
