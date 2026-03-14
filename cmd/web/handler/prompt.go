package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/fanlv/deep-agent-demo/types/model"
)

func (h *Handler) GetPrompt(ctx context.Context, c *app.RequestContext) {
	var req model.GetPromptRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "invalid request: " + err.Error(),
		})
		return
	}

	if req.Key == "" {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "key is required",
		})
		return
	}

	content, err := h.promptService.GetPrompt(ctx, req.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.GetPromptResponse{
		Code:   0,
		Prompt: content,
	})
}

func (h *Handler) SavePrompt(ctx context.Context, c *app.RequestContext) {
	var req model.SavePromptRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "invalid request: " + err.Error(),
		})
		return
	}

	if req.Key == "" {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "key is required",
		})
		return
	}

	if err := h.promptService.SavePrompt(ctx, req.Key, req.Prompt); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.SavePromptResponse{
		Code: 0,
	})
}
