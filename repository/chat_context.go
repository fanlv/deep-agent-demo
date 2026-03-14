package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/deep-agent/sandbox/types/model"
	"github.com/fanlv/deep-agent-demo/services/agent/sandbox"
	"github.com/fanlv/deep-agent-demo/types/path"
)

type ChatContextRepo interface {
	LoadAllMessages() ([]*schema.Message, error)
	CountMessage() (int, error)
	// AppendMessages appends the chat messages to the file.
	AppendMessages(msgs []*schema.Message) error
	// LoadSummaryMessage loads the summary message from the file.
	LoadSummaryMessage() (*SummaryMessage, error)
	// SaveSummaryMessage saves the summary message to the file.
	SaveSummaryMessage(msg *SummaryMessage) error
}

// SummaryMessage represents a summarized conversation with an index
// indicating how many original messages were summarized.
type SummaryMessage struct {
	// Index is the number of original messages that have been summarized.
	Index int `json:"index"`
	// Message is the summary content.
	Message *schema.Message `json:"message"`
}

type chatContextRepo struct {
	sandbox *sandbox.Client
}

func NewChatContextRepo(sb *sandbox.Client) ChatContextRepo {
	return &chatContextRepo{sandbox: sb}
}

func (r *chatContextRepo) LoadAllMessages() ([]*schema.Message, error) {
	filePath := path.MessagesFilePath(r.sandbox.Ctx.Workspace)
	result, err := r.sandbox.Client.FileRead(&model.FileReadRequest{
		File: filePath,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("read messages file failed: %w", err)
	}

	if result.Content == "" {
		return nil, nil
	}

	var messages []*schema.Message
	lines := strings.Split(strings.TrimSpace(result.Content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var msg schema.Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			log.Printf("[chatContextRepo] skip invalid line: %s, err: %v", line, err)
			continue
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (r *chatContextRepo) CountMessage() (int, error) {
	filePath := path.MessagesFilePath(r.sandbox.Ctx.Workspace)
	result, err := r.sandbox.Client.JSONLCountLines(&model.JSONLCountRequest{
		File: filePath,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "not found") {
			return 0, nil
		}
		return 0, fmt.Errorf("count messages failed: %w", err)
	}
	return result.Lines, nil
}

func (r *chatContextRepo) AppendMessages(msgs []*schema.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	if err := r.ensureDeepagentDir(); err != nil {
		return fmt.Errorf("ensure deepagent dir failed: %w", err)
	}

	filePath := path.MessagesFilePath(r.sandbox.Ctx.Workspace)
	JSONStrings := make([]string, 0, len(msgs))
	for _, msg := range msgs {
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("marshal message failed: %w", err)
		}
		JSONStrings = append(JSONStrings, string(jsonBytes))
	}
	if err := r.sandbox.Client.JSONLAppendLine(&model.JSONLAppendRequest{
		File:       filePath,
		JSONString: JSONStrings,
	}); err != nil {
		return fmt.Errorf("append messages failed: %w", err)
	}
	return nil
}

func (r *chatContextRepo) ensureDeepagentDir() error {
	dir := path.MetaDir(r.sandbox.Ctx.Workspace)
	_, err := r.sandbox.BashExecChecked(&model.BashExecRequest{
		Cwd:     r.sandbox.Ctx.Workspace,
		Command: fmt.Sprintf("[ -d %s ] || mkdir -p %s", dir, dir),
	})
	return err
}

func (r *chatContextRepo) LoadSummaryMessage() (*SummaryMessage, error) {
	filePath := path.SummaryFilePath(r.sandbox.Ctx.Workspace)

	result, err := r.sandbox.Client.FileRead(&model.FileReadRequest{
		File: filePath,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("read summary file failed: %w", err)
	}

	if result.Content == "" {
		return nil, nil
	}

	var msg SummaryMessage
	if err := json.Unmarshal([]byte(result.Content), &msg); err != nil {
		return nil, fmt.Errorf("unmarshal summary message failed: %w", err)
	}

	return &msg, nil
}

func (r *chatContextRepo) SaveSummaryMessage(msg *SummaryMessage) error {
	if msg == nil {
		return nil
	}

	if err := r.ensureDeepagentDir(); err != nil {
		return fmt.Errorf("ensure deepagent dir failed: %w", err)
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal summary message failed: %w", err)
	}

	filePath := path.SummaryFilePath(r.sandbox.Ctx.Workspace)
	err = r.sandbox.Client.FileWrite(&model.FileWriteRequest{
		File:    filePath,
		Content: string(data),
	})
	if err != nil {
		return fmt.Errorf("write summary file failed: %w", err)
	}

	return nil
}
