package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/compose"
	"github.com/deep-agent/sandbox/sdk/go/http"
	"github.com/deep-agent/sandbox/types/model"
	"github.com/fanlv/deep-agent-demo/pkg/json"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/pkg/modelbuilder"
	"github.com/fanlv/deep-agent-demo/repository"
	"github.com/fanlv/deep-agent-demo/services/agent/chatctx"
	"github.com/fanlv/deep-agent-demo/services/agent/middlewares"
	"github.com/fanlv/deep-agent-demo/services/agent/sandbox"
)

type AgentEvent = adk.AgentEvent

type DeepAgent struct {
	Runner     *adk.Runner
	CtxManager *chatctx.ChatContextManager
	Sandbox    *sandbox.Client
	mu         sync.RWMutex
	cancel     context.CancelFunc
}

type Config struct {
	Cwd          string
	SystemPrompt string
}

type Option func(*Config)

func WithCwd(cwd string) Option {
	return func(c *Config) {
		c.Cwd = cwd
	}
}

func WithSystemPrompt(prompt string) Option {
	return func(c *Config) {
		c.SystemPrompt = prompt
	}
}

func New(ctx context.Context, sessionID string, modelCfg *modelbuilder.ModelConfig, opts ...Option) (*DeepAgent, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	sb, err := sandbox.New(sessionID, http.WithCwd(cfg.Cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to create sandbox: %w", err)
	}

	tools, err := sb.GetMCPTools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP tools: %w", err)
	}

	if modelCfg == nil {
		return nil, fmt.Errorf("model config is nil")
	}

	logger.Infof(ctx, "modelbuilder.BuildModel use modelCfg: %v", json.String(modelCfg))
	chatModel, err := modelbuilder.BuildModel(ctx, modelCfg, modelbuilder.WithLLMMaxTokens(32768))
	if err != nil {
		return nil, fmt.Errorf("failed to create chat model: %w", err)
	}

	chatContextRepo := repository.NewChatContextRepo(sb)

	handlers, err := middlewares.Init(ctx, &middlewares.InitConfig{
		ChatModel:       chatModel,
		Workspace:       sb.Ctx.Workspace,
		Sandbox:         sb.Client,
		ChatContextRepo: chatContextRepo,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init middlewares: %w", err)
	}

	systemPrompt := injectEnvPrompt(cfg.SystemPrompt, sb.Ctx)

	agent, err := deep.New(ctx, &deep.Config{
		Name:                   "DeepAgent",
		Description:            "an agent for deep task",
		ChatModel:              chatModel,
		Instruction:            systemPrompt,
		WithoutWriteTodos:      true,
		WithoutGeneralSubAgent: true,
		Handlers:               handlers,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		MaxIteration: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

	chatCtxMgr, err := chatctx.New(ctx, sb, chatContextRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat context: %w", err)
	}

	return &DeepAgent{
		Runner:     runner,
		Sandbox:    sb,
		CtxManager: chatCtxMgr,
	}, nil
}

func (d *DeepAgent) Cancel() {
	d.mu.RLock()
	cancel := d.cancel
	d.mu.RUnlock()
	if cancel != nil {
		cancel()
	}
}

func (d *DeepAgent) setCancel(cancel context.CancelFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cancel = cancel
}

func injectEnvPrompt(systemPrompt string, sandboxCtx *model.SandboxContext) string {
	now := time.Now()
	return fmt.Sprintf(`%s
<env>
Operating system: %s
Working directory: %s
Today's date: %s
Timezone: %s
Architecture: %s
</env>
`, systemPrompt,
		sandboxCtx.OS,
		sandboxCtx.Workspace,
		now.Format("2006-01-02"),
		now.Location().String(),
		sandboxCtx.Arch)
}
