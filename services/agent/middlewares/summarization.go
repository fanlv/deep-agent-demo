package middlewares

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/middlewares/summarization"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/repository"
)

func NewSummarizationMW(ctx context.Context, chatModel model.BaseChatModel, repo repository.ChatContextRepo) (adk.ChatModelAgentMiddleware, error) {
	mw, err := summarization.New(ctx, &summarization.Config{
		Model: chatModel,
		Trigger: &summarization.TriggerCondition{
			ContextTokens: 190000,
		},
		EmitInternalEvents: true,
		Finalize: func(ctx context.Context, originalMessages []adk.Message, summary adk.Message) ([]adk.Message, error) {
			// Load all persisted messages to get the accurate index,
			// since in-memory message count may differ after multiple summarizations.
			count, err := repo.CountMessage()
			if err != nil {
				logger.Errorf(ctx, "load messages for summary index failed: %v", err)
			}

			err = repo.SaveSummaryMessage(&repository.SummaryMessage{
				Index:   count,
				Message: summary,
			})
			if err != nil {
				logger.Errorf(ctx, "save summary message failed: %v", err)
			} else {
				logger.Infof(ctx, "saved summary message, index=%d", count)
			}

			// Preserve default behavior: system messages + summary
			var result []adk.Message
			for _, msg := range originalMessages {
				if msg.Role == schema.System {
					result = append(result, msg)
				}
			}
			result = append(result, summary)
			return result, nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create summarization middleware: %w", err)
	}
	return mw, nil
}
