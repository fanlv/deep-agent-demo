package session

import (
	"sync"

	"github.com/fanlv/deep-agent-demo/repository"
	"github.com/fanlv/deep-agent-demo/types/model"
)

type serviceImpl struct {
	sessions map[string]*model.Session
	mu       sync.RWMutex
	repo     repository.SessionRepo
}

func (m *serviceImpl) load() error {
	metas, err := m.repo.LoadAll()
	if err != nil {
		return err
	}

	for _, meta := range metas {
		if !meta.Deleted {
			m.store(meta.ID, meta)
		}
	}

	return nil
}

func (m *serviceImpl) New(modelID int64, systemPrompt string) (*model.Session, error) {
	s := model.NewSession()
	s.ModelID = modelID
	s.SystemPrompt = systemPrompt
	if err := m.Save(s); err != nil {
		return nil, err
	}

	m.store(s.ID, s)
	return s, nil
}

func (m *serviceImpl) store(sid string, s *model.Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[sid] = s
}

func (m *serviceImpl) Get(sid string) (*model.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[sid]
	return s, ok
}

func (m *serviceImpl) Delete(sid string) {
	m.mu.Lock()
	s, ok := m.sessions[sid]
	m.mu.Unlock()

	if ok {
		s.Deleted = true
		m.repo.Save(s.ID, s)
	}
}

func (m *serviceImpl) List() []*model.Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sessions := make([]*model.Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		if !s.Deleted {
			sessions = append(sessions, s)
		}
	}
	return sessions
}

func (m *serviceImpl) Save(s *model.Session) error {
	return m.repo.Save(s.ID, s)
}
