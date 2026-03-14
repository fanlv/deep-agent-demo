import { MessageRoleEnum, MessageStatusEnum, ToolCallStatusEnum } from './protocol';

export interface BaseMessage {
  id: string;
  role: MessageRoleEnum;
  createdAt: number;
  status: MessageStatusEnum;
  content: string;
}

export interface UserMessage extends BaseMessage {
  role: MessageRoleEnum.USER;
}

export interface AssistantMessage extends BaseMessage {
  role: MessageRoleEnum.ASSISTANT;
  name?: string;
  thinkingContent?: string;
  isThinking?: boolean;
}

export interface ToolMessage extends BaseMessage {
  role: MessageRoleEnum.TOOL;
  toolCallId: string;
  toolCallName: string;
  toolCallArgs: string;
  toolCallStatus: ToolCallStatusEnum;
  parentMessageId?: string;
  finishedAt?: number;
}

export type Message = UserMessage | AssistantMessage | ToolMessage;

export interface RunAgentInput {
  sessionId: string;
  messages: Array<{
    id: string;
    type: string;
    content: string;
    timestamp: number;
    role: string;
  }>;
  context?: {
    timestamp: number;
    timezone: string;
    userId?: string;
  };
  state?: Record<string, unknown>;
}
