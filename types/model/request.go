package model

type InitRequest struct {
	ModelID int64 `json:"modelId"`
}

type RunAgentRequest struct {
	SessionID string           `json:"sessionId"`
	ModelID  int64            `json:"modelId,omitempty"`
	Messages []RequestMessage `json:"messages"`
	Context  *RequestContext  `json:"context,omitempty"`
	State    map[string]any   `json:"state,omitempty"`
}

type RequestMessage struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Role      string `json:"role"`
}

type RequestContext struct {
	Timestamp int64  `json:"timestamp"`
	Timezone  string `json:"timezone"`
	UserID    string `json:"userId,omitempty"`
}

type GetPromptRequest struct {
	Key string `json:"key"`
}

type SavePromptRequest struct {
	Key    string `json:"key"`
	Prompt string `json:"prompt"`
}
