package main

import (
	"context"
	"log/slog"

	"github.com/coder/acp-go-sdk"
	"github.com/fanlv/deep-agent-demo/services/agent"
)

type acpEventHandler struct {
	conn      *acp.AgentSideConnection
	sessionID string
	logger    *slog.Logger
}

var _ agent.EventHandler = (*acpEventHandler)(nil)

func newACPEventHandler(conn *acp.AgentSideConnection, sessionID string, logger *slog.Logger) *acpEventHandler {
	return &acpEventHandler{
		conn:      conn,
		sessionID: sessionID,
		logger:    logger,
	}
}

func (h *acpEventHandler) OnMessageStart() error { return nil }

func (h *acpEventHandler) OnTokenUsage(totalTokens int) error { return nil }

func (h *acpEventHandler) OnMessageDelta(content string) error {
	return h.conn.SessionUpdate(context.Background(), acp.SessionNotification{
		SessionId: acp.SessionId(h.sessionID),
		Update:    acp.UpdateAgentMessageText(content),
	})
}

func (h *acpEventHandler) OnMessageEnd() error { return nil }

func (h *acpEventHandler) OnThoughtStart() error { return nil }

func (h *acpEventHandler) OnThoughtDelta(content string) error {
	return h.conn.SessionUpdate(context.Background(), acp.SessionNotification{
		SessionId: acp.SessionId(h.sessionID),
		Update:    acp.UpdateAgentThoughtText(content),
	})
}

func (h *acpEventHandler) OnThoughtEnd() error { return nil }

func (h *acpEventHandler) OnToolCallStart(id, name string) error {
	return h.conn.SessionUpdate(context.Background(), acp.SessionNotification{
		SessionId: acp.SessionId(h.sessionID),
		Update: acp.StartToolCall(
			acp.ToolCallId(id),
			name,
			acp.WithStartKind(getToolKind(name)),
			acp.WithStartStatus(acp.ToolCallStatusInProgress),
		),
	})
}

func (h *acpEventHandler) OnToolCallArgs(id, args string) error {
	return h.conn.SessionUpdate(context.Background(), acp.SessionNotification{
		SessionId: acp.SessionId(h.sessionID),
		Update: acp.UpdateToolCall(
			acp.ToolCallId(id),
			acp.WithUpdateStatus(acp.ToolCallStatusInProgress),
			acp.WithUpdateContent([]acp.ToolCallContent{
				acp.ToolContent(acp.TextBlock(args)),
			}),
		),
	})
}

func (h *acpEventHandler) OnToolCallResult(id, content string) error {
	return h.conn.SessionUpdate(context.Background(), acp.SessionNotification{
		SessionId: acp.SessionId(h.sessionID),
		Update: acp.UpdateToolCall(
			acp.ToolCallId(id),
			acp.WithUpdateContent([]acp.ToolCallContent{
				acp.ToolContent(acp.TextBlock(truncate(content, 500))),
			}),
		),
	})
}

func (h *acpEventHandler) OnToolCallEnd(id string) error {
	return h.conn.SessionUpdate(context.Background(), acp.SessionNotification{
		SessionId: acp.SessionId(h.sessionID),
		Update: acp.UpdateToolCall(
			acp.ToolCallId(id),
			acp.WithUpdateStatus(acp.ToolCallStatusCompleted),
		),
	})
}

func (h *acpEventHandler) OnError(err error) {
	h.logger.Error("agent event error", "error", err)
}

func getToolKind(name string) acp.ToolKind {
	switch name {
	case "Read":
		return acp.ToolKindRead
	case "Edit", "Write", "browser_type":
		return acp.ToolKindEdit
	case "Glob", "Grep", "WebSearch":
		return acp.ToolKindSearch
	case "Bash":
		return acp.ToolKindExecute
	case "WebFetch":
		return acp.ToolKindFetch
	default:
		return acp.ToolKindOther
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
