package model

type InitResponse struct {
	SessionID string `json:"sessionId"`
	CreatedAt int64  `json:"createdAt"`
}

type SessionInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type ListSessionsResponse struct {
	Sessions []SessionInfo `json:"sessions"`
}

type HistoryMessage struct {
	ID               string         `json:"id"`
	Role             MessageRole    `json:"role"`
	Content          string         `json:"content"`
	ReasoningContent string         `json:"reasoningContent,omitempty"`
	ToolCallID       string         `json:"toolCallId,omitempty"`
	ToolCalls        []ToolCallInfo `json:"toolCalls,omitempty"`
}

type ToolCallInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type GetMessagesResponse struct {
	ModelID    int64            `json:"modelId"`
	Messages   []HistoryMessage `json:"messages"`
	TokenUsage *TokenUsage      `json:"tokenUsage,omitempty"`
}

type GetPromptResponse struct {
	Code   int    `json:"code"`
	Prompt string `json:"prompt"`
}

type SavePromptResponse struct {
	Code int `json:"code"`
}
