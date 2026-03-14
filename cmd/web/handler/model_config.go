package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/fanlv/deep-agent-demo/types/model"
)

func (h *Handler) GetModelList(ctx context.Context, c *app.RequestContext) {
	list, err := h.modelConfig.GetProviderModelList(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	list = sanitizeProviderModelList(list)
	c.JSON(http.StatusOK, map[string]any{
		"code":                0,
		"provider_model_list": list,
	})
}

func sanitizeProviderModelList(list []*model.ProviderModelList) []*model.ProviderModelList {
	if list == nil {
		return nil
	}

	out := make([]*model.ProviderModelList, 0, len(list))
	for _, item := range list {
		if item == nil {
			out = append(out, nil)
			continue
		}

		newItem := &model.ProviderModelList{
			Provider:  item.Provider,
			ModelList: make([]*model.ModelInstance, 0, len(item.ModelList)),
		}

		for _, m := range item.ModelList {
			if m == nil {
				newItem.ModelList = append(newItem.ModelList, nil)
				continue
			}

			mCopy := *m
			if mCopy.Connection != nil {
				connCopy := *mCopy.Connection
				connCopy.APIKey = maskSecret(connCopy.APIKey)
				mCopy.Connection = &connCopy
			}
			newItem.ModelList = append(newItem.ModelList, &mCopy)
		}
		out = append(out, newItem)
	}
	return out
}

func maskSecret(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 4 {
		return "****"
	}
	if len(s) <= 8 {
		return strings.Repeat("*", len(s)-4) + s[len(s)-4:]
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}

func (h *Handler) CreateModel(ctx context.Context, c *app.RequestContext) {
	var req model.CreateModelRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "invalid request: " + err.Error(),
		})
		return
	}

	if req.DisplayName == "" || req.Connection == nil || req.Connection.Model == "" {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "display_name and connection.model are required",
		})
		return
	}

	id, err := h.modelConfig.CreateModel(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"code": 0,
		"id":   id,
	})
}

func (h *Handler) DeleteModel(ctx context.Context, c *app.RequestContext) {
	idStr := c.Query("id")
	if idStr == "" {
		var body struct {
			ID string `json:"id"`
		}
		if err := c.BindJSON(&body); err == nil {
			idStr = body.ID
		}
	}

	if idStr == "" {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "id is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "invalid id",
		})
		return
	}

	if err := h.modelConfig.DeleteModel(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"code": 0,
	})
}
