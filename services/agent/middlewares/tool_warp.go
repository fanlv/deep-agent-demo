package middlewares

import (
	"context"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/fanlv/deep-agent-demo/pkg/json"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
)

type toolWarpMiddleware struct {
	adk.BaseChatModelAgentMiddleware
}

func NewToolWarpMiddleware() adk.ChatModelAgentMiddleware {
	return &toolWarpMiddleware{}
}

type toolWarpBaseModel struct {
	inner model.BaseChatModel
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	rs := []rune(s)
	if len(rs) <= max {
		return s
	}
	if max == 1 {
		return "…"
	}
	return string(rs[:max-1]) + "…"
}

func printInputMessagesIfContains(ctx context.Context, stage string, input []*schema.Message) {
	logger.Info("=========================================================")
	for idx, msg := range input {
		logStr := truncateRunes(json.String(msg), 200)

		ok := strings.Contains(msg.Content, "system-reminder")
		if ok {
			logStr = msg.Content
		}
		logger.Infof(ctx, "[%s] input(%d) :\n%s", stage, idx, logStr)
	}
	logger.Info("=========================================================\n")
}

func (w *toolWarpBaseModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	printInputMessagesIfContains(ctx, "Generate", input)
	return w.inner.Generate(ctx, input, opts...)
}

func (w *toolWarpBaseModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	printInputMessagesIfContains(ctx, "Stream", input)
	return w.inner.Stream(ctx, input, opts...)
}

func (b *toolWarpMiddleware) WrapInvokableToolCall(ctx context.Context, endpoint adk.InvokableToolCallEndpoint,
	tCtx *adk.ToolContext) (adk.InvokableToolCallEndpoint, error) {
	return func(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
		// log.Printf("Tool %s (call %s) starting with args: %s", tCtx.Name, tCtx.CallID, argumentsInJSON)

		result, err := endpoint(ctx, argumentsInJSON, opts...)

		if err != nil {
			log.Printf("Tool %s failed: %v", tCtx.Name, err)
			return err.Error(), nil
		}

		// log.Printf("Tool %s completed with result: %s", tCtx.Name, result)
		return result, nil
	}, nil

}

func (b *toolWarpMiddleware) WrapStreamableToolCall(ctx context.Context, endpoint adk.StreamableToolCallEndpoint,
	tCtx *adk.ToolContext) (adk.StreamableToolCallEndpoint, error) {
	return func(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (*schema.StreamReader[string], error) {
		log.Printf("Tool %s (call %s) starting with args: %s", tCtx.Name, tCtx.CallID, argumentsInJSON)

		result, err := endpoint(ctx, argumentsInJSON, opts...)

		if err != nil {
			// log.Printf("Tool %s failed: %v", tCtx.Name, err)
			return schema.StreamReaderFromArray([]string{err.Error()}), nil
		}

		return result, nil
	}, nil
}

func (b *toolWarpMiddleware) WrapModel(_ context.Context, m model.BaseChatModel, _ *adk.ModelContext) (model.BaseChatModel, error) {
	return &toolWarpBaseModel{inner: m}, nil
}
