package middlewares

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/middlewares/agentsmd"
	sandbox "github.com/deep-agent/sandbox/sdk/go"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
)

func NewAgentDocLoadMW(ctx context.Context, sb sandbox.Sandbox) (adk.ChatModelAgentMiddleware, error) {
	backend := &sandboxBackend{
		client: sb,
	}

	mw, err := agentsmd.New(ctx, &agentsmd.Config{
		Backend:       backend,
		AgentsMDFiles: []string{"/home/sandbox/agent/prompts/agents_md.md"},
		OnLoadWarning: func(filePath string, err error) {
			logger.Infof(ctx, "[OnLoadWarning] load file %s failed, err: %v", filePath, err)
		},
	})

	if err != nil {
		return nil, fmt.Errorf("[NewPlanTaskMW] create plan task middleware failed, err: %w", err)
	}

	return mw, nil
}
