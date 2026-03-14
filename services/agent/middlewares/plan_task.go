package middlewares

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/middlewares/plantask"
	sandbox "github.com/deep-agent/sandbox/sdk/go"
)

func NewPlanTaskMW(ctx context.Context, BaseDir string, sb sandbox.Sandbox) (adk.ChatModelAgentMiddleware, error) {
	backend := &sandboxBackend{
		client: sb,
	}

	mw, err := plantask.New(ctx, &plantask.Config{
		Backend: backend,
		BaseDir: fmt.Sprintf("%s/.tasks", BaseDir),
	})

	if err != nil {
		return nil, fmt.Errorf("[NewPlanTaskMW] create plan task middleware failed, err: %w", err)
	}

	return mw, nil
}
