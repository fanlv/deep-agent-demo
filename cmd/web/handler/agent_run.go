package handler

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/pkg/modelbuilder"
	"github.com/fanlv/deep-agent-demo/services/agent"
	"github.com/fanlv/deep-agent-demo/types/model"
	"github.com/google/uuid"
)

func (h *Handler) AgentRun(ctx context.Context, c *app.RequestContext) {
	logger.Infof(ctx, "[AgentRun] request received")
	var req model.RunAgentRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf(ctx, "[AgentRun] Failed to parse request: %v", err)
		c.JSON(consts.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	logger.Infof(ctx, "[AgentRun] sessionId: %s, modelId: %d, messages: %d", req.SessionID, req.ModelID, len(req.Messages))

	s, ok := h.sessionService.Get(req.SessionID)
	if !ok {
		logger.Errorf(ctx, "[AgentRun] Session not found: %s", req.SessionID)
		c.JSON(consts.StatusBadRequest, map[string]string{"error": "Session not found, please call /init first"})
		return
	}

	if req.ModelID != 0 && req.ModelID != s.ModelID {
		s.ModelID = req.ModelID
		if err := h.sessionService.Save(s); err != nil {
			logger.Errorf(ctx, "Failed to save session model_id: %v", err)
		}
	}

	runID := uuid.New().String()
	logger.Infof(ctx, "[AgentRun] starting SSE stream, sessionId: %s, runId: %s", req.SessionID, runID)

	c.SetStatusCode(consts.StatusOK)
	w := sse.NewWriter(c)
	stream := NewSSEStream(ctx, w, req.SessionID, runID)
	defer stream.Close()

	if err := stream.SendRunStarted(); err != nil {
		logger.Errorf(ctx, "[AgentRun] Failed to send RUN_STARTED: %v", err)
		return
	}
	logger.Infof(ctx, "[AgentRun] RUN_STARTED sent, resolving model config")

	modelCfg, err := h.resolveModelCfg(ctx, s.ModelID)
	if err != nil {
		logger.Errorf(ctx, "[AgentRun] Failed to resolve model config: %v", err)
		stream.SendRunError(err.Error(), "INIT_ERROR")
		return
	}

	logger.Infof(ctx, "[AgentRun] model resolved, getting agent")
	deepAgent, err := h.agentService.GetOrCreate(ctx, req.SessionID, modelCfg,
		agent.WithSystemPrompt(s.SystemPrompt))
	if err != nil {
		logger.Errorf(ctx, "[AgentRun] Failed to initialize agent: %v", err)
		stream.SendRunError(err.Error(), "INIT_ERROR")
		return
	}

	logger.Infof(ctx, "[AgentRun] agent ready, running")
	userMessages := h.parseUserMessages(ctx, s, &req)

	handler := newSSEEventHandler(stream)
	if err := deepAgent.Run(ctx, userMessages, handler); err != nil {
		if ctx.Err() != nil {
			logger.Infof(ctx, "[AgentRun] cancelled")
		} else {
			logger.Errorf(ctx, "[AgentRun] error: %v", err)
			stream.SendRunError(err.Error(), "AGENT_ERROR")
		}
		return
	}

	logger.Infof(ctx, "[AgentRun] finished")
	if err := stream.SendRunFinished(); err != nil {
		logger.Errorf(ctx, "[AgentRun] Failed to send RUN_FINISHED: %v", err)
	}
}

func (h *Handler) parseUserMessages(ctx context.Context, s *model.Session, req *model.RunAgentRequest) []*schema.Message {
	messages := make([]*schema.Message, 0, len(req.Messages))
	for _, msg := range req.Messages {
		switch msg.Role {
		case string(schema.User):
			messages = append(messages, schema.UserMessage(msg.Content))
			h.tryUpdateSessionTitleFromUserContent(ctx, s, msg.Content)
		default:
			logger.Errorf(ctx, "Unexpected client request message role: role=%s id=%s type=%s", msg.Role, msg.ID, msg.Type)
		}
	}
	return messages
}

func (h *Handler) tryUpdateSessionTitleFromUserContent(ctx context.Context, s *model.Session, content string) {
	if s.Title != "New Chat" && s.Title != "" {
		return
	}

	title := []rune(content)
	if len(title) > 30 {
		title = append(title[:50], []rune("...")...)
	}
	s.Title = string(title)
	if err := h.sessionService.Save(s); err != nil {
		logger.Errorf(ctx, "Failed to save session: %v", err)
	}
}

func (h *Handler) resolveModelCfg(ctx context.Context, modelID int64) (*modelbuilder.ModelConfig, error) {
	if modelID == 0 {
		return nil, fmt.Errorf("model_id is required")
	}
	if h.modelConfig == nil {
		return nil, fmt.Errorf("model config is nil")
	}

	inst, err := h.modelConfig.GetModelByID(ctx, modelID)
	if err != nil {
		return nil, err
	}

	return &modelbuilder.ModelConfig{
		ModelClass:   inst.ModelClass,
		Connection:   inst.Connection,
		ThinkingType: inst.ThinkingType,
	}, nil
}
