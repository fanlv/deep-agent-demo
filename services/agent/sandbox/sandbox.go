package sandbox

import (
	"fmt"
	"log"
	"strings"
	"time"

	sandbox "github.com/deep-agent/sandbox/sdk/go"
	httpSandbox "github.com/deep-agent/sandbox/sdk/go/http"
	"github.com/deep-agent/sandbox/types/model"
)

type Client struct {
	sessionID string
	Client    sandbox.Sandbox
	Ctx       *model.SandboxContext
	baseURL   string
}

const defaultBaseURL = "http://localhost:8080"

func (c *Client) BashExecChecked(req *model.BashExecRequest) (*model.BashExecResult, error) {
	resp, err := c.Client.BashExec(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("bash exec returned nil response")
	}
	if resp.ExitCode == 0 {
		return resp, nil
	}

	msgParts := make([]string, 0, 2)
	if s := strings.TrimSpace(resp.Output); s != "" {
		msgParts = append(msgParts, s)
	}
	if s := strings.TrimSpace(resp.Error); s != "" {
		msgParts = append(msgParts, s)
	}
	msg := strings.Join(msgParts, "\n")
	if msg == "" {
		msg = "unknown error"
	}

	return resp, fmt.Errorf("bash exec failed (exit_code=%d): %s", resp.ExitCode, msg)
}

func New(sessionID string, opts ...httpSandbox.Option) (*Client, error) {
	opts = append([]httpSandbox.Option{httpSandbox.WithTimeout(60 * time.Second)}, opts...)

	client := httpSandbox.NewClient(defaultBaseURL, sessionID, opts...)

	sandboxCtx, err := client.GetContext()
	if err != nil {
		return nil, fmt.Errorf("[newSandboxBackend] get context failed, err: %w", err)
	}

	// 确保 workspace 目录存在
	workspace := sandboxCtx.Workspace

	sb := &Client{
		sessionID: sessionID,
		Client:    client,
		Ctx:       sandboxCtx,
		baseURL:   defaultBaseURL,
	}

	resp, err := sb.BashExecChecked(&model.BashExecRequest{
		Cwd:     "/home",
		Command: fmt.Sprintf("[ -d %s ] || mkdir -p %s", workspace, workspace),
	})
	if err != nil {
		return nil, fmt.Errorf("[New] create workspace dir failed, workspace: %s (check sandbox workspace permissions/mounts), err: %w", workspace, err)
	}

	log.Printf("[New] create workspace dir success, workspace: %s, resp: %s", workspace, resp.Output)

	return sb, nil
}
