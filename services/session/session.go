package session

import (
	"log"

	"github.com/fanlv/deep-agent-demo/repository"
	"github.com/fanlv/deep-agent-demo/types/model"
)

type Service interface {
	New(modelID int64, systemPrompt string) (*model.Session, error)
	Get(sid string) (*model.Session, bool)
	Save(s *model.Session) error
	Delete(sid string)
	List() []*model.Session
}

func NewService() (Service, error) {
	repo, err := repository.NewSessionRepo()
	if err != nil {
		return nil, err
	}

	m := &serviceImpl{
		sessions: make(map[string]*model.Session),
		repo:     repo,
	}

	if err := m.load(); err != nil {
		log.Printf("[session.Manager] load sessions failed: %v", err)
	}

	return m, nil
}
