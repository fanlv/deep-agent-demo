package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/fanlv/deep-agent-demo/pkg/tokenizer"
	"github.com/fanlv/deep-agent-demo/services/agent"
	"github.com/fanlv/deep-agent-demo/types/model"
)

func (h *Handler) ListSessions(ctx context.Context, c *app.RequestContext) {
	sessionList := h.sessionService.List()

	sessions := make([]model.SessionInfo, 0, len(sessionList))
	for _, s := range sessionList {
		sessions = append(sessions, model.SessionInfo{
			ID:        s.ID,
			Title:     s.Title,
			CreatedAt: s.CreatedAt.UnixMilli(),
			UpdatedAt: s.UpdatedAt.UnixMilli(),
		})
	}

	c.JSON(http.StatusOK, model.ListSessionsResponse{Sessions: sessions})
}

func (h *Handler) DeleteSession(ctx context.Context, c *app.RequestContext) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "sessionId is required"})
		return
	}

	if ag, ok := h.agentService.Get(sessionID); ok {
		ag.Cancel()
	}
	h.sessionService.Delete(sessionID)
	h.agentService.Delete(sessionID)

	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetSessionMessages(ctx context.Context, c *app.RequestContext) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "sessionId is required"})
		return
	}

	s, ok := h.sessionService.Get(sessionID)
	if !ok {
		c.JSON(http.StatusNotFound, map[string]string{"error": "session not found"})
		return
	}

	modelCfg, err := h.resolveModelCfg(ctx, s.ModelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	ag, err := h.agentService.GetOrCreate(ctx, sessionID, modelCfg, agent.WithSystemPrompt(s.SystemPrompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	chatMessages := ag.CtxManager.LoadAllMessages()

	messages := make([]model.HistoryMessage, 0, len(chatMessages))

	for i, msg := range chatMessages {
		if msg.Role == schema.System {
			continue
		}

		historyMsg := model.HistoryMessage{
			ID:               fmt.Sprintf("msg_%d", i),
			Role:             model.MessageRole(msg.Role),
			Content:          msg.Content,
			ReasoningContent: msg.ReasoningContent,
		}

		if msg.ToolCallID != "" {
			historyMsg.ToolCallID = msg.ToolCallID
		}

		if len(msg.ToolCalls) > 0 {
			historyMsg.ToolCalls = make([]model.ToolCallInfo, len(msg.ToolCalls))
			for j, tc := range msg.ToolCalls {
				historyMsg.ToolCalls[j] = model.ToolCallInfo{
					ID:        tc.ID,
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				}
			}
		}

		messages = append(messages, historyMsg)
	}

	tokens := tokenizer.MessagesTokenCounter(ctx, chatMessages)

	c.JSON(http.StatusOK, model.GetMessagesResponse{
		ModelID:    s.ModelID,
		Messages:   messages,
		TokenUsage: &model.TokenUsage{TotalTokens: tokens},
	})
}
