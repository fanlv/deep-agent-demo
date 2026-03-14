package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/types/model"
)

type SSEStream struct {
	ctx       context.Context
	writer    *sse.Writer
	sessionID string
	runID     string

	respPreview     []byte
	respPreviewMax  int
	respPayloadSize int
}

func NewSSEStream(ctx context.Context, w *sse.Writer, sessionID, runID string) *SSEStream {
	return &SSEStream{
		ctx:            ctx,
		writer:         w,
		sessionID:      sessionID,
		runID:          runID,
		respPreviewMax: 500,
	}
}

func (s *SSEStream) sendEvent(event any) error {
	data, err := sonic.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	s.respPayloadSize += len(data)
	if len(s.respPreview) < s.respPreviewMax {
		remain := s.respPreviewMax - len(s.respPreview)
		if remain > len(data) {
			remain = len(data)
		}
		s.respPreview = append(s.respPreview, data[:remain]...)
	}
	return s.writer.WriteEvent("", "message", data)
}

func (s *SSEStream) Close() {
	s.writer.WriteEvent("", "message", []byte("[DONE]"))
	s.writer.Close()
	if s.ctx != nil {
		preview := string(s.respPreview)
		if preview != "" {
			logger.Infof(s.ctx, "[AgentRun] SSE resp preview sessionId=%s runId=%s bytes=%d preview=%s",
				s.sessionID, s.runID, s.respPayloadSize, preview)
		} else {
			logger.Infof(s.ctx, "[AgentRun] SSE resp preview sessionId=%s runId=%s bytes=%d preview=%s",
				s.sessionID, s.runID, s.respPayloadSize, "[empty]")
		}
	}
}

func (s *SSEStream) baseEvent(eventType model.EventType) model.BaseEvent {
	return model.BaseEvent{
		Type:      eventType,
		SessionID: s.sessionID,
		RunID:     s.runID,
		Timestamp: time.Now().UnixMilli(),
	}
}

func (s *SSEStream) SendRunStarted() error {
	return s.sendEvent(model.RunStartedEvent{
		BaseEvent: s.baseEvent(model.EventTypeRunStarted),
	})
}

func (s *SSEStream) SendRunFinished() error {
	return s.sendEvent(model.RunFinishedEvent{
		BaseEvent: s.baseEvent(model.EventTypeRunFinished),
	})
}

func (s *SSEStream) SendRunError(message, code string) error {
	return s.sendEvent(model.RunErrorEvent{
		BaseEvent: s.baseEvent(model.EventTypeRunError),
		Message:   message,
		Code:      code,
	})
}

func (s *SSEStream) SendTextMessageStart(messageID string, role model.MessageRole, name string, isThinking bool) error {
	event := model.TextMessageStartEvent{
		BaseEvent: s.baseEvent(model.EventTypeTextMessageStart),
		MessageID: messageID,
		Role:      role,
		Name:      name,
	}
	if isThinking {
		event.External = map[string]any{"isThinking": true}
	}
	return s.sendEvent(event)
}

func (s *SSEStream) SendTextMessageContent(messageID string, role model.MessageRole, delta string, isThinking bool) error {
	event := model.TextMessageContentEvent{
		BaseEvent: s.baseEvent(model.EventTypeTextMessageContent),
		MessageID: messageID,
		Role:      role,
		Delta:     delta,
	}
	if isThinking {
		event.External = map[string]any{"isThinking": true}
	}
	return s.sendEvent(event)
}

func (s *SSEStream) SendTextMessageEnd(messageID string, role model.MessageRole) error {
	return s.sendEvent(model.TextMessageEndEvent{
		BaseEvent: s.baseEvent(model.EventTypeTextMessageEnd),
		MessageID: messageID,
		Role:      role,
	})
}

func (s *SSEStream) SendCustomEvent(name string, value any) error {
	return s.sendEvent(model.CustomEvent{
		BaseEvent: s.baseEvent(model.EventTypeCustom),
		Name:      name,
		Value:     value,
	})
}

func (s *SSEStream) SendToolCallStart(toolCallID, toolCallName, parentMessageID string) error {
	return s.sendEvent(model.ToolCallStartEvent{
		BaseEvent:       s.baseEvent(model.EventTypeToolCallStart),
		ToolCallID:      toolCallID,
		ToolCallName:    toolCallName,
		ParentMessageID: parentMessageID,
		ToolCallStatus:  model.ToolCallStatusProcessing,
	})
}

func (s *SSEStream) SendToolCallArgs(toolCallID, toolCallName, delta string) error {
	return s.sendEvent(model.ToolCallArgsEvent{
		BaseEvent:      s.baseEvent(model.EventTypeToolCallArgs),
		ToolCallID:     toolCallID,
		ToolCallName:   toolCallName,
		Delta:          delta,
		ToolCallStatus: model.ToolCallStatusProcessing,
	})
}

func (s *SSEStream) SendToolCallResult(toolCallID, toolCallName, delta string, status model.ToolCallStatus) error {
	return s.sendEvent(model.ToolCallResultEvent{
		BaseEvent:      s.baseEvent(model.EventTypeToolCallResult),
		ToolCallID:     toolCallID,
		ToolCallName:   toolCallName,
		Delta:          delta,
		ToolCallStatus: status,
	})
}

func (s *SSEStream) SendToolCallEnd(toolCallID, toolCallName string, status model.ToolCallStatus) error {
	return s.sendEvent(model.ToolCallEndEvent{
		BaseEvent:      s.baseEvent(model.EventTypeToolCallEnd),
		ToolCallID:     toolCallID,
		ToolCallName:   toolCallName,
		ToolCallStatus: status,
	})
}
