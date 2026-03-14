package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/types/consts"
	"github.com/fanlv/deep-agent-demo/types/model"
)

func (h *Handler) AgentInit(ctx context.Context, c *app.RequestContext) {
	var req model.InitRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf(ctx, "Failed to parse init request: %v", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	systemPrompt, err := h.promptService.GetPrompt(ctx, consts.KeySystemPrompt)
	if err != nil {
		logger.Errorf(ctx, "Failed to load system prompt: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load system prompt"})
		return
	}

	s, err := h.sessionService.New(req.ModelID, systemPrompt)
	if err != nil {
		logger.Errorf(ctx, "Failed to create session: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}
	logger.Infof(ctx, "Agent init, sessionId: %s", s.ID)

	c.JSON(http.StatusOK, model.InitResponse{
		SessionID: s.ID,
		CreatedAt: s.CreatedAt.UnixMilli(),
	})
}
