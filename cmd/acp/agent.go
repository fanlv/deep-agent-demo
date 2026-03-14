package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/coder/acp-go-sdk"

	"github.com/fanlv/deep-agent-demo/pkg/json"
	"github.com/fanlv/deep-agent-demo/pkg/modelbuilder"
	"github.com/fanlv/deep-agent-demo/services/agent"
	"github.com/fanlv/deep-agent-demo/services/session"
)

type ACPHandler struct {
	conn     *acp.AgentSideConnection
	sessions session.Service
	agents   agent.Service
	modelCfg *modelbuilder.ModelConfig
	logger   *slog.Logger
}

var (
	_ acp.Agent             = (*ACPHandler)(nil)
	_ acp.AgentLoader       = (*ACPHandler)(nil)
	_ acp.AgentExperimental = (*ACPHandler)(nil)
)

func newACPHandler(logger *slog.Logger) (*ACPHandler, error) {
	sessions, err := session.NewService()
	if err != nil {
		return nil, err
	}
	modelCfg := modelbuilder.LoadConfigFromEnv()
	if modelCfg == nil {
		return nil, fmt.Errorf("no model config found in environment variables")
	}
	return &ACPHandler{
		sessions: sessions,
		agents:   agent.NewService(),
		modelCfg: modelCfg,
		logger:   logger,
	}, nil
}

func (a *ACPHandler) SetAgentConnection(conn *acp.AgentSideConnection) {
	a.conn = conn
}

func (a *ACPHandler) Initialize(ctx context.Context, params acp.InitializeRequest) (acp.InitializeResponse, error) {
	a.logger.Info("Initialize", "request", json.String(params))
	resp := acp.InitializeResponse{
		ProtocolVersion: acp.ProtocolVersionNumber,
		AgentInfo: &acp.Implementation{
			Name:    "DeepAgent",
			Version: "1.0.0",
		},
		AgentCapabilities: acp.AgentCapabilities{
			LoadSession: false,
		},
	}
	a.logger.Info("Initialize", "response", json.String(resp))
	return resp, nil
}

func (a *ACPHandler) Authenticate(ctx context.Context, params acp.AuthenticateRequest) (acp.AuthenticateResponse, error) {
	a.logger.Info("Authenticate", "request", json.String(params))
	resp := acp.AuthenticateResponse{}
	a.logger.Info("Authenticate", "response", json.String(resp))
	return resp, nil
}

func (a *ACPHandler) NewSession(ctx context.Context, params acp.NewSessionRequest) (acp.NewSessionResponse, error) {
	a.logger.Info("NewSession", "request", json.String(params))

	s, err := a.sessions.New(0, "")
	if err != nil {
		return acp.NewSessionResponse{}, err
	}
	resp := acp.NewSessionResponse{SessionId: acp.SessionId(s.ID)}
	a.logger.Info("NewSession", "response", json.String(resp))
	return resp, nil
}

func (a *ACPHandler) LoadSession(ctx context.Context, params acp.LoadSessionRequest) (acp.LoadSessionResponse, error) {
	a.logger.Info("LoadSession", "request", json.String(params))
	resp := acp.LoadSessionResponse{}
	a.logger.Info("LoadSession", "response", json.String(resp))
	return resp, nil
}

func (a *ACPHandler) SetSessionMode(ctx context.Context, params acp.SetSessionModeRequest) (acp.SetSessionModeResponse, error) {
	a.logger.Info("SetSessionMode", "request", json.String(params))
	resp := acp.SetSessionModeResponse{}
	a.logger.Info("SetSessionMode", "response", json.String(resp))
	return resp, nil
}

func (a *ACPHandler) SetSessionModel(ctx context.Context, params acp.SetSessionModelRequest) (acp.SetSessionModelResponse, error) {
	a.logger.Info("SetSessionModel", "request", json.String(params))
	resp := acp.SetSessionModelResponse{}
	a.logger.Info("SetSessionModel", "response", json.String(resp))
	return resp, nil
}

func (a *ACPHandler) Cancel(ctx context.Context, params acp.CancelNotification) error {
	a.logger.Info("Cancel", "request", json.String(params))
	ag, ok := a.agents.Get(string(params.SessionId))
	if ok && ag != nil {
		ag.Cancel()
	}
	return nil
}

func (a *ACPHandler) Prompt(ctx context.Context, params acp.PromptRequest) (acp.PromptResponse, error) {
	a.logger.Info("Prompt", "request", json.String(params))
	sid := string(params.SessionId)

	if _, ok := a.sessions.Get(sid); !ok {
		err := fmt.Errorf("session %s not found", sid)
		a.logger.Error("Prompt", "error", err)
		return acp.PromptResponse{}, err
	}

	deepAgent, err := a.agents.GetOrCreate(ctx, sid, a.modelCfg)
	if err != nil {
		a.logger.Error("Prompt", "error", err)
		return acp.PromptResponse{}, err
	}

	userMessage := a.parseUserMessage(params.Prompt)
	a.logger.Info("Prompt", "userMessage", userMessage.Content)

	handler := newACPEventHandler(a.conn, sid, a.logger)
	if err := deepAgent.Run(ctx, []*schema.Message{userMessage}, handler); err != nil {
		if ctx.Err() != nil {
			resp := acp.PromptResponse{StopReason: acp.StopReasonCancelled}
			a.logger.Info("Prompt", "response", json.String(resp), "cancelled", true)
			return resp, nil
		}
		a.logger.Error("Prompt", "error", err)
		return acp.PromptResponse{}, err
	}

	resp := acp.PromptResponse{StopReason: acp.StopReasonEndTurn}
	a.logger.Info("Prompt", "response", json.String(resp))
	return resp, nil
}

func (a *ACPHandler) parseUserMessage(blocks []acp.ContentBlock) *schema.Message {
	promptText := ""
	for _, block := range blocks {
		promptText = strings.TrimPrefix(block.Text.Text, "USER: ")
	}
	return schema.UserMessage(promptText)
}
