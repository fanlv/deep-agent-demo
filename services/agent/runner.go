package agent

import (
	"context"
	"io"

	"github.com/cloudwego/eino/schema"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/pkg/tokenizer"
)

type TokenUsageHandler interface {
	OnTokenUsage(totalTokens int) error
}
type MessageHandler interface {
	OnMessageStart() error
	OnMessageDelta(content string) error
	OnMessageEnd() error
}

type ThoughtHandler interface {
	OnThoughtStart() error
	OnThoughtDelta(content string) error
	OnThoughtEnd() error
}

type ToolCallHandler interface {
	OnToolCallStart(id, name string) error
	OnToolCallArgs(id, args string) error
	OnToolCallResult(id, content string) error
	OnToolCallEnd(id string) error
}

type EventHandler interface {
	MessageHandler
	ThoughtHandler
	ToolCallHandler
	TokenUsageHandler
	OnError(err error)
}

func (d *DeepAgent) Run(ctx context.Context, userMessages []*schema.Message, handler EventHandler) error {
	d.Cancel()

	runCtx, cancel := context.WithCancel(ctx)
	d.setCancel(cancel)
	defer d.setCancel(nil)

	d.CtxManager.AppendMessages(userMessages...)

	chatMessages, err := d.CtxManager.LoadMessagesForLLM(ctx)
	if err != nil {
		return err
	}

	iter := d.Runner.Run(runCtx, chatMessages)

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		history := d.handleEvent(runCtx, handler, event)
		if len(history) > 0 {
			d.CtxManager.AppendMessages(history...)
		}

		chatMessages, err := d.CtxManager.LoadMessagesForLLM(ctx)
		if err != nil {
			logger.Warnf(ctx, "[Run] LoadMessagesForLLM  error = %v ", err)
		} else {
			tokens := tokenizer.MessagesTokenCounter(ctx, chatMessages)
			handler.OnTokenUsage(tokens)
		}
	}

	return nil
}

func (d *DeepAgent) handleEvent(ctx context.Context, handler EventHandler, event *AgentEvent) []*schema.Message {
	if event.Err != nil {
		handler.OnError(event.Err)
		return nil
	}

	if event.Output == nil || event.Output.MessageOutput == nil {
		return nil
	}

	output := event.Output.MessageOutput
	if output.Message != nil {
		return d.handleDirectMessage(ctx, handler, output.Message)
	}

	if output.MessageStream != nil {
		return d.handleMessageStream(ctx, handler, output.MessageStream)
	}

	return nil
}

func (d *DeepAgent) handleDirectMessage(ctx context.Context, handler EventHandler, m *schema.Message) []*schema.Message {
	if len(m.Content) > 0 {
		if m.Role == schema.Tool {
			handler.OnToolCallResult(m.ToolCallID, m.Content)
			handler.OnToolCallEnd(m.ToolCallID)
		} else {
			handler.OnMessageStart()
			handler.OnMessageDelta(m.Content)
			handler.OnMessageEnd()
		}
	}

	for _, tc := range m.ToolCalls {
		handler.OnToolCallStart(tc.ID, tc.Function.Name)
		if tc.Function.Arguments != "" {
			handler.OnToolCallArgs(tc.ID, tc.Function.Arguments)
		}
	}

	return []*schema.Message{m}
}

func (d *DeepAgent) handleMessageStream(ctx context.Context, handler EventHandler, stream *schema.StreamReader[*schema.Message]) []*schema.Message {
	chunks := []*schema.Message{}
	state := &streamState{}

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			handler.OnError(err)
			continue
		}
		if chunk == nil {
			continue
		}

		chunks = append(chunks, chunk)
		d.processStreamChunk(ctx, handler, chunk, state)
	}

	state.finalize(handler)

	if len(chunks) == 0 {
		return nil
	}

	msg, err := schema.ConcatMessages(chunks)
	if err != nil {
		handler.OnError(err)
		return nil
	}

	return []*schema.Message{msg}
}

type streamState struct {
	inMessage bool
	inThought bool

	openToolResultID string

	toolCallIDs     map[int]string
	toolCallNames   map[int]string
	toolCallStarted map[int]bool
}

func (s *streamState) finalize(handler EventHandler) {
	if s.inMessage {
		handler.OnMessageEnd()
	}
	if s.inThought {
		handler.OnThoughtEnd()
	}
	if s.openToolResultID != "" {
		handler.OnToolCallEnd(s.openToolResultID)
	}
}

func (d *DeepAgent) processStreamChunk(ctx context.Context, handler EventHandler, chunk *schema.Message, state *streamState) {
	if chunk == nil {
		return
	}

	if chunk.ReasoningContent != "" {
		if !state.inThought {
			if state.inMessage {
				handler.OnMessageEnd()
				state.inMessage = false
			}
			handler.OnThoughtStart()
			state.inThought = true
		}
		handler.OnThoughtDelta(chunk.ReasoningContent)
	}

	if chunk.Role == schema.Tool {
		if state.openToolResultID == "" {
			state.openToolResultID = chunk.ToolCallID
		} else if chunk.ToolCallID != "" && chunk.ToolCallID != state.openToolResultID {
			handler.OnToolCallEnd(state.openToolResultID)
			state.openToolResultID = chunk.ToolCallID
		}

		if content := streamChunkText(chunk); content != "" {
			handler.OnToolCallResult(chunk.ToolCallID, content)
		}
	} else if chunk.Content != "" {
		if !state.inMessage {
			if state.inThought {
				handler.OnThoughtEnd()
				state.inThought = false
			}
			handler.OnMessageStart()
			state.inMessage = true
		}
		handler.OnMessageDelta(chunk.Content)
	}

	// stream tool call
	for i, tc := range chunk.ToolCalls {
		key := toolCallKey(tc, i)
		if state.toolCallIDs == nil {
			state.toolCallIDs = make(map[int]string)
			state.toolCallNames = make(map[int]string)
			state.toolCallStarted = make(map[int]bool)
		}

		if tc.ID != "" {
			state.toolCallIDs[key] = tc.ID
		}
		if tc.Function.Name != "" {
			state.toolCallNames[key] = tc.Function.Name
		}

		id := state.toolCallIDs[key]
		name := state.toolCallNames[key]

		if !state.toolCallStarted[key] && id != "" && name != "" {
			if state.inMessage {
				handler.OnMessageEnd()
				state.inMessage = false
			}
			if state.inThought {
				handler.OnThoughtEnd()
				state.inThought = false
			}
			handler.OnToolCallStart(id, name)
			state.toolCallStarted[key] = true
		}

		if tc.Function.Arguments != "" && id != "" {
			handler.OnToolCallArgs(id, tc.Function.Arguments)
		}
	}
}

func toolCallKey(tc schema.ToolCall, pos int) int {
	if tc.Index != nil {
		return *tc.Index
	}
	return pos
}

func streamChunkText(m *schema.Message) string {
	if m == nil {
		return ""
	}
	if m.Content != "" {
		return m.Content
	}
	if len(m.UserInputMultiContent) == 0 {
		return ""
	}
	out := ""
	for _, part := range m.UserInputMultiContent {
		if part.Text != "" {
			out += part.Text
		}
	}
	return out
}
