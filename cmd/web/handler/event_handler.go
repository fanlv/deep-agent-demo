package handler

import (
	"context"

	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/services/agent"
	"github.com/fanlv/deep-agent-demo/types/model"
	"github.com/google/uuid"
)

type sseEventHandler struct {
	stream    *SSEStream
	messageID string
}

var _ agent.EventHandler = (*sseEventHandler)(nil)

func newSSEEventHandler(stream *SSEStream) *sseEventHandler {
	return &sseEventHandler{stream: stream}
}

func (h *sseEventHandler) OnMessageStart() error {
	h.messageID = uuid.New().String()
	return h.stream.SendTextMessageStart(h.messageID, model.MessageRoleAssistant, "", false)
}

func (h *sseEventHandler) OnMessageDelta(content string) error {
	return h.stream.SendTextMessageContent(h.messageID, model.MessageRoleAssistant, content, false)
}

func (h *sseEventHandler) OnMessageEnd() error {
	return h.stream.SendTextMessageEnd(h.messageID, model.MessageRoleAssistant)
}

func (h *sseEventHandler) OnThoughtStart() error {
	h.messageID = uuid.New().String()
	return h.stream.SendTextMessageStart(h.messageID, model.MessageRoleAssistant, "", true)
}

func (h *sseEventHandler) OnThoughtDelta(content string) error {
	return h.stream.SendTextMessageContent(h.messageID, model.MessageRoleAssistant, content, true)
}

func (h *sseEventHandler) OnThoughtEnd() error {
	return h.stream.SendTextMessageEnd(h.messageID, model.MessageRoleAssistant)
}

func (h *sseEventHandler) OnToolCallStart(id, name string) error {
	return h.stream.SendToolCallStart(id, name, "")
}

func (h *sseEventHandler) OnToolCallArgs(id, args string) error {
	return h.stream.SendToolCallArgs(id, "", args)
}

func (h *sseEventHandler) OnToolCallResult(id, content string) error {
	return h.stream.SendToolCallResult(id, "", content, model.ToolCallStatusProcessing)
}

func (h *sseEventHandler) OnToolCallEnd(id string) error {
	return h.stream.SendToolCallEnd(id, "", model.ToolCallStatusSuccess)
}

func (h *sseEventHandler) OnTokenUsage(totalTokens int) error {
	logger.Infof(context.TODO(), "[OnTokenUsage] totalTokens = %d", totalTokens)
	return h.stream.SendCustomEvent("token_usage", model.TokenUsage{TotalTokens: totalTokens})
}

func (h *sseEventHandler) OnError(err error) {
	h.stream.SendRunError(err.Error(), "-1")
	logger.Errorf(context.Background(), "agent event error: %v", err)
}
