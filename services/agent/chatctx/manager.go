package chatctx

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/schema"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/repository"
	"github.com/fanlv/deep-agent-demo/services/agent/sandbox"
)

type ChatContextManager struct {
	repo     repository.ChatContextRepo
	messages []*schema.Message
	mu       sync.RWMutex
}

func New(ctx context.Context, sb *sandbox.Client, repo repository.ChatContextRepo) (*ChatContextManager, error) {
	mgr := &ChatContextManager{repo: repo}

	historyMessages, err := repo.LoadAllMessages()
	if err != nil {
		fmt.Printf("[ChatContextManager] load history messages failed: %v\n", err)
	}

	if len(historyMessages) > 0 {
		mgr.messages = historyMessages
	}

	return mgr, nil
}

func (m *ChatContextManager) LoadAllMessages() []*schema.Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.messages
}

func (m *ChatContextManager) LoadMessagesForLLM(ctx context.Context) ([]*schema.Message, error) {
	summary, err := m.repo.LoadSummaryMessage()
	if err != nil {
		return nil, fmt.Errorf("load summary message failed: %w", err)
	}

	allMsgs := m.LoadAllMessages()

	if summary == nil || summary.Message == nil {
		return allMsgs, nil
	}

	var result []*schema.Message
	result = append(result, summary.Message)

	if summary.Index < len(allMsgs) {
		result = append(result, allMsgs[summary.Index:]...)
	}

	logger.Infof(ctx, "load messages for llm, summary index=%d, msg count=%d", summary.Index, len(result))
	return result, nil
}

func (m *ChatContextManager) AppendMessages(msgs ...*schema.Message) {
	m.mu.Lock()
	m.messages = append(m.messages, msgs...)
	m.mu.Unlock()

	if err := m.repo.AppendMessages(msgs); err != nil {
		fmt.Printf("[ChatContextManager] persist messages failed: %v\n", err)
	}
}
