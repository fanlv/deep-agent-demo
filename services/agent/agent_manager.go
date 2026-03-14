package agent

import (
	"context"
	"sync"

	"github.com/fanlv/deep-agent-demo/pkg/modelbuilder"
)

type Service interface {
	GetOrCreate(ctx context.Context, sessionID string, modelCfg *modelbuilder.ModelConfig, opts ...Option) (*DeepAgent, error)
	Get(sessionID string) (*DeepAgent, bool)
	Delete(sessionID string)
	List() []*DeepAgent
}

type service struct {
	agents map[string]*DeepAgent
	mu     sync.RWMutex
}

func NewService() Service {
	return &service{
		agents: make(map[string]*DeepAgent),
	}
}

func (s *service) GetOrCreate(ctx context.Context, sessionID string, modelCfg *modelbuilder.ModelConfig, opts ...Option) (*DeepAgent, error) {
	s.mu.RLock()
	ag, ok := s.agents[sessionID]
	s.mu.RUnlock()
	if ok {
		return ag, nil
	}

	newAgent, err := New(ctx, sessionID, modelCfg, opts...)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if ag, ok := s.agents[sessionID]; ok {
		return ag, nil
	}
	s.agents[sessionID] = newAgent
	return newAgent, nil
}

func (s *service) Get(sessionID string) (*DeepAgent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ag, ok := s.agents[sessionID]
	return ag, ok
}

func (s *service) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.agents, sessionID)
}

func (s *service) List() []*DeepAgent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	agents := make([]*DeepAgent, 0, len(s.agents))
	for _, ag := range s.agents {
		agents = append(agents, ag)
	}
	return agents
}
