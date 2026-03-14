package prompt

import (
	"context"
	"fmt"

	"github.com/fanlv/deep-agent-demo/repository"
)

type Service interface {
	GetPrompt(ctx context.Context, key string) (string, error)
	SavePrompt(ctx context.Context, key string, content string) error
}

type serviceImpl struct {
	repo repository.PromptRepo
}

func NewService() (Service, error) {
	repo, err := repository.NewPromptRepo()
	if err != nil {
		return nil, fmt.Errorf("init prompt repo failed: %w", err)
	}
	return &serviceImpl{repo: repo}, nil
}

func (s *serviceImpl) GetPrompt(ctx context.Context, key string) (string, error) {
	result, err := s.repo.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if result == "" {
		return "You are a helpful assistant. Use the instructions below and the tools available to you to assist the user.", nil
	}

	return result, nil
}

func (s *serviceImpl) SavePrompt(ctx context.Context, key string, content string) error {
	return s.repo.Save(ctx, key, content)
}
