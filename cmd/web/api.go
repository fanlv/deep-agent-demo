package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/fanlv/deep-agent-demo/cmd/web/handler"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
)

func registerRoutes(s *server.Hertz, h *handler.Handler) {
	s.GET("/health", healthHandler)

	api := s.Group("/api/v1")

	agent := api.Group("/agent")
	agent.POST("/init", h.AgentInit)
	agent.POST("/run", h.AgentRun)

	sessions := api.Group("/sessions")
	sessions.GET("", h.ListSessions)
	sessions.DELETE("/:sessionId", h.DeleteSession)
	sessions.GET("/:sessionId/messages", h.GetSessionMessages)

	prompt := api.Group("/prompt")
	prompt.POST("/get", h.GetPrompt)
	prompt.POST("/save", h.SavePrompt)

	config := api.Group("/config")
	model := config.Group("/model")
	model.GET("/list", h.GetModelList)
	model.POST("/create", h.CreateModel)
	model.POST("/delete", h.DeleteModel)
}

func healthHandler(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

type httpLogConfig struct {
	skipExactPaths map[string]struct{}
	skipConfigs    []*skipConfig
}
type skipConfig struct {
	prefix string
	suffix string
}

var httpLogCfg = httpLogConfig{
	skipExactPaths: map[string]struct{}{
		"/api/v1/config/model/list": {},
		"/api/v1/sessions":          {},
	},
	skipConfigs: []*skipConfig{
		{prefix: "/api/v1/sessions/", suffix: "/messages"},
	},
}

func (cfg httpLogConfig) shouldSkip(reqPath string) bool {
	if _, ok := cfg.skipExactPaths[reqPath]; ok {
		return true
	}

	for _, sc := range cfg.skipConfigs {
		if sc == nil || sc.prefix == "" || sc.suffix == "" {
			continue
		}
		if strings.HasPrefix(reqPath, sc.prefix) && strings.HasSuffix(reqPath, sc.suffix) {
			return true
		}
	}

	return false
}

func loggerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		reqPath := string(c.Request.URI().Path())
		if httpLogCfg.shouldSkip(reqPath) {
			c.Next(ctx)
			return
		}

		start := time.Now()
		reqBody := truncate(string(c.Request.Body()), 500)
		c.Next(ctx)
		respBody := truncate(string(c.Response.Body()), 500)
		if respBody == "" {
			ct := strings.ToLower(string(c.Response.Header.ContentType()))
			if strings.Contains(ct, "text/event-stream") {
				respBody = "[SSE stream]"
			}
		}
		if reqBody == "" {
			logger.Infof(ctx, "[HTTP] %s %s %d %v\nresp=%s",
				c.Method(), reqPath, c.Response.StatusCode(), time.Since(start),
				respBody)
			return
		}
		logger.Infof(ctx, "[HTTP] %s %s %d %v\nreq=%s\nresp=%s",
			c.Method(), reqPath, c.Response.StatusCode(), time.Since(start),
			reqBody, respBody)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
