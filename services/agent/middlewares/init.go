package middlewares

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	sandbox "github.com/deep-agent/sandbox/sdk/go"
	"github.com/fanlv/deep-agent-demo/repository"
)

type InitConfig struct {
	ChatModel       model.BaseChatModel
	Workspace       string
	Sandbox         sandbox.Sandbox
	ChatContextRepo repository.ChatContextRepo
}

func Init(ctx context.Context, cfg *InitConfig) ([]adk.ChatModelAgentMiddleware, error) {
	summarizationMW, err := NewSummarizationMW(ctx, cfg.ChatModel, cfg.ChatContextRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to create summarization middleware: %w", err)
	}

	planTaskMW, err := NewPlanTaskMW(ctx, cfg.Workspace, cfg.Sandbox)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan task middleware: %w", err)
	}

	agentDocLoadMW, err := NewAgentDocLoadMW(ctx, cfg.Sandbox)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent doc load middleware: %w", err)
	}

	reductionMW, err := NewReductionMW(ctx, cfg.Workspace, cfg.Sandbox)
	if err != nil {
		return nil, fmt.Errorf("failed to create reduction middleware: %w", err)
	}

	return []adk.ChatModelAgentMiddleware{
		summarizationMW,
		planTaskMW,
		agentDocLoadMW,
		NewToolWarpMiddleware(),
		reductionMW,
	}, nil
}
