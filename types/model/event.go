package model

type EventType string

const (
	EventTypeTextMessageStart   EventType = "TEXT_MESSAGE_START"
	EventTypeTextMessageContent EventType = "TEXT_MESSAGE_CONTENT"
	EventTypeTextMessageEnd     EventType = "TEXT_MESSAGE_END"
	EventTypeToolCallStart      EventType = "TOOL_CALL_START"
	EventTypeToolCallArgs       EventType = "TOOL_CALL_ARGS"
	EventTypeToolCallResult     EventType = "TOOL_CALL_RESULT"
	EventTypeToolCallEnd        EventType = "TOOL_CALL_END"
	EventTypeCustom             EventType = "CUSTOM"
	EventTypeRunStarted         EventType = "RUN_STARTED"
	EventTypeRunFinished        EventType = "RUN_FINISHED"
	EventTypeRunError           EventType = "RUN_ERROR"
	EventTypeArtifactStart      EventType = "ARTIFACT_START"
	EventTypeArtifactContent    EventType = "ARTIFACT_CONTENT"
	EventTypeArtifactEnd        EventType = "ARTIFACT_END"
	EventTypeStateSnapshot      EventType = "STATE_SNAPSHOT"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
	MessageRoleTool      MessageRole = "tool"
	MessageRoleCustom    MessageRole = "custom"
)

type ToolCallStatus string

const (
	ToolCallStatusProcessing ToolCallStatus = "Processing"
	ToolCallStatusSuccess    ToolCallStatus = "Success"
	ToolCallStatusError      ToolCallStatus = "Error"
)

type BaseEvent struct {
	Type      EventType      `json:"type"`
	SessionID string         `json:"sessionId"`
	RunID     string         `json:"runId"`
	StepID    string         `json:"stepId,omitempty"`
	Timestamp int64          `json:"timestamp"`
	External  map[string]any `json:"external,omitempty"`
}

type RunStartedEvent struct {
	BaseEvent
}

type TokenUsage struct {
	TotalTokens int `json:"totalTokens"`
}

type RunFinishedEvent struct {
	BaseEvent
}

type RunErrorEvent struct {
	BaseEvent
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type TextMessageStartEvent struct {
	BaseEvent
	MessageID   string      `json:"messageId"`
	Role        MessageRole `json:"role"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
}

type TextMessageContentEvent struct {
	BaseEvent
	MessageID   string      `json:"messageId"`
	Role        MessageRole `json:"role"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Delta       string      `json:"delta"`
}

type TextMessageEndEvent struct {
	BaseEvent
	MessageID   string      `json:"messageId"`
	Role        MessageRole `json:"role"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
}

type ToolCallStartEvent struct {
	BaseEvent
	ToolCallID      string         `json:"toolCallId"`
	ToolCallName    string         `json:"toolCallName"`
	ParentMessageID string         `json:"parentMessageId,omitempty"`
	ToolCallStatus  ToolCallStatus `json:"toolCallStatus,omitempty"`
}

type ToolCallArgsEvent struct {
	BaseEvent
	ToolCallID      string         `json:"toolCallId"`
	ToolCallName    string         `json:"toolCallName,omitempty"`
	ParentMessageID string         `json:"parentMessageId,omitempty"`
	Delta           string         `json:"delta"`
	ToolCallStatus  ToolCallStatus `json:"toolCallStatus,omitempty"`
}

type ToolCallResultEvent struct {
	BaseEvent
	ToolCallID      string         `json:"toolCallId"`
	ToolCallName    string         `json:"toolCallName,omitempty"`
	ParentMessageID string         `json:"parentMessageId,omitempty"`
	Delta           string         `json:"delta"`
	ToolCallStatus  ToolCallStatus `json:"toolCallStatus"`
}

type ToolCallEndEvent struct {
	BaseEvent
	ToolCallID      string         `json:"toolCallId"`
	ToolCallName    string         `json:"toolCallName,omitempty"`
	ParentMessageID string         `json:"parentMessageId,omitempty"`
	ToolCallStatus  ToolCallStatus `json:"toolCallStatus,omitempty"`
}

type ArtifactStartEvent struct {
	BaseEvent
	ArtifactID   string `json:"artifactId"`
	ArtifactType string `json:"artifactType"`
	Title        string `json:"title,omitempty"`
}

type ArtifactContentEvent struct {
	BaseEvent
	ArtifactID string `json:"artifactId"`
	Delta      string `json:"delta"`
}

type ArtifactEndEvent struct {
	BaseEvent
	ArtifactID string `json:"artifactId"`
}

type StateSnapshotEvent struct {
	BaseEvent
	Messages []any `json:"messages,omitempty"`
}

type CustomEvent struct {
	BaseEvent
	Name  string `json:"name"`
	Value any    `json:"value"`
}
