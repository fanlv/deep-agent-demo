package middlewares

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/middlewares/reduction"
	sandbox "github.com/deep-agent/sandbox/sdk/go"
	"github.com/fanlv/deep-agent-demo/types/path"
)

func NewReductionMW(ctx context.Context, workspace string, sb sandbox.Sandbox) (adk.ChatModelAgentMiddleware, error) {
	backend := &sandboxBackend{
		client: sb,
	}

	rootDir := path.ReductionDir(workspace)
	return reduction.New(ctx, &reduction.Config{
		Backend:          backend,
		RootDir:          rootDir,
		ReadFileToolName: "Read",
	})
}
