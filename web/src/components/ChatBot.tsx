import { useCallback, useEffect, useRef, useState } from 'react';
import { useAgentChat } from '../hooks/useAgentChat';
import { MessageList } from './MessageList';
import { ChatInput } from './ChatInput';
import { fetchModelOptions, ModelOption } from '../utils/models';
import './ChatBot.css';

interface ChatBotProps {
  baseUrl?: string;
  initialMessage?: string | null;
  existingSessionId: string;
  onBack?: () => void;
  onSessionCreated?: (sessionId: string) => void;
  onToggleSidebar?: () => void;
}

export function ChatBot(props: ChatBotProps) {
  const { baseUrl, initialMessage, existingSessionId, onBack, onSessionCreated } = props;
  const { messages, isLoading, isLoadingHistory, error, sessionId, sessionModelId, totalTokens, sendMessage, stopGeneration, clearMessages } = useAgentChat({
    baseUrl,
    existingSessionId,
  });
  const [models, setModels] = useState<ModelOption[]>([]);
  const [selectedModelId, setSelectedModelId] = useState<number | null>(() => {
    const stored = localStorage.getItem(`session_model_${existingSessionId}`);
    const n = stored ? Number(stored) : NaN;
    return Number.isFinite(n) ? n : null;
  });
  const [hasUserSelectedModel, setHasUserSelectedModel] = useState(false);
  const initialMessageSent = useRef(false);

  useEffect(() => {
    if (sessionId && onSessionCreated) {
      onSessionCreated(sessionId);
    }
  }, [sessionId, onSessionCreated]);

  useEffect(() => {
    fetchModelOptions().then((allModels) => {
      setModels(allModels);
      setSelectedModelId((prev) => {
        if (prev || allModels.length === 0) {
          return prev;
        }
        const next = allModels[0].model.id;
        localStorage.setItem(`session_model_${existingSessionId}`, String(next));
        return next;
      });
    }).catch((err) => {
      console.error('Failed to fetch model options:', err);
    });
  }, [existingSessionId]);

  useEffect(() => {
    if (!hasUserSelectedModel && sessionModelId) {
      localStorage.setItem(`session_model_${existingSessionId}`, String(sessionModelId));
    }
  }, [existingSessionId, hasUserSelectedModel, sessionModelId]);

  const effectiveModelId = hasUserSelectedModel ? selectedModelId : sessionModelId ?? selectedModelId;

  const handleSendMessage = useCallback(
    (content: string) => {
      sendMessage(content, effectiveModelId);
    },
    [sendMessage, effectiveModelId]
  );

  useEffect(() => {
    console.log('ChatBot useEffect:', { initialMessage, initialMessageSent: initialMessageSent.current, isLoadingHistory, messagesLength: messages.length, error });
    if (initialMessage && !initialMessageSent.current && !isLoadingHistory && !error && messages.length === 0) {
      console.log('Sending initial message:', initialMessage);
      initialMessageSent.current = true;
      sendMessage(initialMessage, effectiveModelId).catch((err) => {
        console.error('Failed to send initial message:', err);
      });
    }
  }, [effectiveModelId, error, initialMessage, isLoadingHistory, messages.length, sendMessage]);

  const handleNewChat = () => {
    clearMessages();
    if (onBack) {
      onBack();
    }
  };

  return (
    <div className="chatbot-container">
      <header className="chatbot-header">
        <div className="header-left">
          <button className="back-button" onClick={handleNewChat} title="Back to Home">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M19 12H5M12 19l-7-7 7-7" />
            </svg>
          </button>
          <span className="header-logo" onClick={handleNewChat} style={{ cursor: 'pointer' }}>🤖 Deep Agent</span>
        </div>
        <nav className="header-nav">
          <div className="nav-search">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="11" cy="11" r="8" />
              <path d="M21 21l-4.35-4.35" />
            </svg>
            <span>Search Docs</span>
            <kbd>⌘K</kbd>
          </div>
          <div className="nav-links">
            <span className="nav-link" onClick={handleNewChat}>New Chat</span>
            <span className="nav-link" onClick={clearMessages}>Clear</span>
          </div>
        </nav>
      </header>

      {error && (
        <div className="error-banner">
          <span>{error}</span>
          <button onClick={() => {
            window.location.href = '/';
          }}>返回首页</button>
        </div>
      )}

      {isLoadingHistory && (
        <div className="chatbot-main">
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%', color: '#888' }}>
            Loading...
          </div>
        </div>
      )}

      {!isLoadingHistory && (
      <div className="chatbot-main">
        <MessageList
          messages={messages}
          isLoading={isLoading}
          onSendMessage={handleSendMessage}
        />
        <ChatInput
          onSend={handleSendMessage}
          onStop={stopGeneration}
          isLoading={isLoading}
          totalTokens={totalTokens}
          models={models}
          selectedModelId={effectiveModelId}
          onSelectModel={(modelId) => {
            setHasUserSelectedModel(true);
            setSelectedModelId(modelId);
            localStorage.setItem(`session_model_${existingSessionId}`, String(modelId));
          }}
        />
      </div>
      )}
    </div>
  );
}
