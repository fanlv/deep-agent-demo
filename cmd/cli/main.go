package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/fanlv/deep-agent-demo/pkg/modelbuilder"
	"github.com/fanlv/deep-agent-demo/services/agent"
)

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

func main() {
	ctx := context.Background()

	t := time.Now()
	sessionID := fmt.Sprintf("session-%s-%06d", t.Format("20060102-150405"), t.Nanosecond()/1000)

	fmt.Printf("\n%s%s🤖 Deep Agent CLI%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", colorGray, colorReset)

	modelCfg := modelbuilder.LoadConfigFromEnv()
	if modelCfg == nil {
		log.Fatal("no model config found in environment variables")
	}

	deepAgent, err := agent.New(ctx, sessionID, modelCfg)
	if err != nil {
		log.Fatal(err)
	}

	handler := &cliEventHandler{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s%s👤 You:%s ", colorBold, colorGreen, colorReset)
		scanner.Scan()
		input := scanner.Text()
		if input == "" {
			continue
		}
		fmt.Println()

		userMessage := schema.UserMessage(input)
		if err := deepAgent.Run(ctx, []*schema.Message{userMessage}, handler); err != nil {
			log.Printf("Run error: %v", err)
		}

		fmt.Printf("\n%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", colorGray, colorReset)
	}
}

type cliEventHandler struct {
	inThought   bool
	currentArgs strings.Builder
}

var _ agent.EventHandler = (*cliEventHandler)(nil)

func (h *cliEventHandler) OnMessageStart() error {
	fmt.Printf("%s%s🤖 Assistant:%s ", colorBold, colorCyan, colorReset)
	return nil
}

func (h *cliEventHandler) OnMessageDelta(content string) error {
	fmt.Print(content)
	return nil
}

func (h *cliEventHandler) OnMessageEnd() error {
	fmt.Println()
	return nil
}

func (h *cliEventHandler) OnThoughtStart() error {
	h.inThought = true
	fmt.Printf("%s💭 Thinking...%s\n", colorGray, colorReset)
	return nil
}

func (h *cliEventHandler) OnThoughtDelta(content string) error {
	return nil
}

func (h *cliEventHandler) OnThoughtEnd() error {
	h.inThought = false
	return nil
}

func (h *cliEventHandler) OnToolCallStart(id, name string) error {
	h.currentArgs.Reset()
	fmt.Printf("%s🔧 %s%s\n", colorYellow, name, colorReset)
	return nil
}

func (h *cliEventHandler) OnToolCallArgs(id, args string) error {
	h.currentArgs.WriteString(args)
	return nil
}

func (h *cliEventHandler) OnToolCallResult(id, content string) error {
	if h.currentArgs.Len() > 0 {
		args := h.currentArgs.String()
		args = strings.ReplaceAll(args, "\n", " ")
		if len(args) > 100 {
			args = args[:100] + "..."
		}
		fmt.Printf("   %s%s%s\n", colorGray, args, colorReset)
		h.currentArgs.Reset()
	}

	lines := strings.Split(content, "\n")
	preview := lines[0]
	if len(preview) > 80 {
		preview = preview[:80] + "..."
	}
	if len(lines) > 1 {
		preview += fmt.Sprintf(" (+%d lines)", len(lines)-1)
	}
	fmt.Printf("   %s→ %s%s\n", colorGray, preview, colorReset)
	return nil
}

func (h *cliEventHandler) OnToolCallEnd(id string) error {
	return nil
}

func (h *cliEventHandler) OnTokenUsage(totalTokens int) error {
	return nil
}

func (h *cliEventHandler) OnError(err error) {
	fmt.Printf("%s❌ Error: %v%s\n", colorRed, err, colorReset)
}
