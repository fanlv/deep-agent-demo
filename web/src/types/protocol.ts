export enum EventTypeEnum {
  TEXT_MESSAGE_START = 'TEXT_MESSAGE_START',
  TEXT_MESSAGE_CONTENT = 'TEXT_MESSAGE_CONTENT',
  TEXT_MESSAGE_END = 'TEXT_MESSAGE_END',
  TOOL_CALL_START = 'TOOL_CALL_START',
  TOOL_CALL_ARGS = 'TOOL_CALL_ARGS',
  TOOL_CALL_RESULT = 'TOOL_CALL_RESULT',
  TOOL_CALL_END = 'TOOL_CALL_END',
  CUSTOM = 'CUSTOM',
  RUN_STARTED = 'RUN_STARTED',
  RUN_FINISHED = 'RUN_FINISHED',
  RUN_ERROR = 'RUN_ERROR',
  ARTIFACT_START = 'ARTIFACT_START',
  ARTIFACT_CONTENT = 'ARTIFACT_CONTENT',
  ARTIFACT_END = 'ARTIFACT_END',
  STATE_SNAPSHOT = 'STATE_SNAPSHOT',
}

export enum MessageRoleEnum {
  USER = 'user',
  ASSISTANT = 'assistant',
  SYSTEM = 'system',
  TOOL = 'tool',
  CUSTOM = 'custom',
}

export enum MessageStatusEnum {
  Loading = 'Loading',
  Started = 'Started',
  Finished = 'Finished',
  Error = 'Error',
}

export enum ToolCallStatusEnum {
  Processing = 'Processing',
  Success = 'Success',
  Error = 'Error',
}

export interface BaseEvent {
  type: EventTypeEnum;
  sessionId: string;
  runId: string;
  stepId?: string;
  timestamp: number;
  external?: Record<string, unknown>;
}

export interface RunStartedEvent extends BaseEvent {
  type: EventTypeEnum.RUN_STARTED;
}

export interface TokenUsage {
  totalTokens: number;
}

export interface RunFinishedEvent extends BaseEvent {
  type: EventTypeEnum.RUN_FINISHED;
}

export interface RunErrorEvent extends BaseEvent {
  type: EventTypeEnum.RUN_ERROR;
  message: string;
  code?: string;
}

export interface TextMessageStartEvent extends BaseEvent {
  type: EventTypeEnum.TEXT_MESSAGE_START;
  messageId: string;
  role: MessageRoleEnum;
  name?: string;
  description?: string;
  external?: {
    isThinking?: boolean;
    [key: string]: unknown;
  };
}

export interface TextMessageContentEvent extends BaseEvent {
  type: EventTypeEnum.TEXT_MESSAGE_CONTENT;
  messageId: string;
  role: MessageRoleEnum;
  name?: string;
  description?: string;
  delta: string;
  external?: {
    isThinking?: boolean;
    [key: string]: unknown;
  };
}

export interface TextMessageEndEvent extends BaseEvent {
  type: EventTypeEnum.TEXT_MESSAGE_END;
  messageId: string;
  role: MessageRoleEnum;
  name?: string;
  description?: string;
  external?: {
    isThinking?: boolean;
    [key: string]: unknown;
  };
}

export interface ToolCallStartEvent extends BaseEvent {
  type: EventTypeEnum.TOOL_CALL_START;
  toolCallId: string;
  toolCallName: string;
  parentMessageId?: string;
  toolCallStatus?: ToolCallStatusEnum;
}

export interface ToolCallArgsEvent extends BaseEvent {
  type: EventTypeEnum.TOOL_CALL_ARGS;
  toolCallId: string;
  toolCallName?: string;
  parentMessageId?: string;
  delta: string;
  toolCallStatus?: ToolCallStatusEnum;
}

export interface ToolCallResultEvent extends BaseEvent {
  type: EventTypeEnum.TOOL_CALL_RESULT;
  toolCallId: string;
  toolCallName?: string;
  parentMessageId?: string;
  delta: string;
  toolCallStatus: ToolCallStatusEnum;
}

export interface ToolCallEndEvent extends BaseEvent {
  type: EventTypeEnum.TOOL_CALL_END;
  toolCallId: string;
  toolCallName?: string;
  parentMessageId?: string;
  toolCallStatus?: ToolCallStatusEnum;
}

export interface ArtifactStartEvent extends BaseEvent {
  type: EventTypeEnum.ARTIFACT_START;
  artifactId: string;
  artifactType: string;
  title?: string;
}

export interface ArtifactContentEvent extends BaseEvent {
  type: EventTypeEnum.ARTIFACT_CONTENT;
  artifactId: string;
  delta: string;
}

export interface ArtifactEndEvent extends BaseEvent {
  type: EventTypeEnum.ARTIFACT_END;
  artifactId: string;
}

export interface StateSnapshotEvent extends BaseEvent {
  type: EventTypeEnum.STATE_SNAPSHOT;
  messages?: unknown[];
}

export interface CustomEvent extends BaseEvent {
  type: EventTypeEnum.CUSTOM;
  name: string;
  value: unknown;
}

export type AgentEvent =
  | RunStartedEvent
  | RunFinishedEvent
  | RunErrorEvent
  | TextMessageStartEvent
  | TextMessageContentEvent
  | TextMessageEndEvent
  | ToolCallStartEvent
  | ToolCallArgsEvent
  | ToolCallResultEvent
  | ToolCallEndEvent
  | ArtifactStartEvent
  | ArtifactContentEvent
  | ArtifactEndEvent
  | StateSnapshotEvent
  | CustomEvent;
