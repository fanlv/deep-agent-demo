package handler

import (
	"context"

	"github.com/fanlv/deep-agent-demo/services/agent"
	"github.com/fanlv/deep-agent-demo/services/config"
	"github.com/fanlv/deep-agent-demo/services/prompt"
	"github.com/fanlv/deep-agent-demo/services/session"
)

type Handler struct {
	sessionService session.Service
	agentService   agent.Service
	modelConfig    config.ModelConfigService
	promptService  prompt.Service
}

func NewHandler(ctx context.Context) (*Handler, error) {
	ss, err := session.NewService()
	if err != nil {
		return nil, err
	}

	ps, err := prompt.NewService()
	if err != nil {
		return nil, err
	}

	return &Handler{
		sessionService: ss,
		agentService:   agent.NewService(),
		modelConfig:    config.NewModelConfigService(ctx),
		promptService:  ps,
	}, nil
}
