package repository

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fanlv/deep-agent-demo/types/path"
)

type PromptRepo interface {
	Get(ctx context.Context, key string) (string, error)
	Save(ctx context.Context, key string, content string) error
}

type filePromptRepo struct {
	dir string
}

func NewPromptRepo() (PromptRepo, error) {
	dir, err := path.PromptsDir()
	if err != nil {
		return nil, err
	}
	return &filePromptRepo{dir: dir}, nil
}

func (r *filePromptRepo) Get(_ context.Context, key string) (string, error) {
	data, err := os.ReadFile(filepath.Join(r.dir, key+".md"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

func (r *filePromptRepo) Save(_ context.Context, key string, content string) error {
	if err := os.MkdirAll(r.dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(r.dir, key+".md"), []byte(content), 0644)
}
