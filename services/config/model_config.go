package config

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/fanlv/deep-agent-demo/pkg/modelbuilder"
	"github.com/fanlv/deep-agent-demo/repository"
	"github.com/fanlv/deep-agent-demo/types/model"
)

type ModelConfigService interface {
	GetProviderModelList(ctx context.Context) ([]*model.ProviderModelList, error)
	CreateModel(ctx context.Context, req *model.CreateModelRequest) (int64, error)
	DeleteModel(ctx context.Context, id int64) error
	GetModelByID(ctx context.Context, id int64) (*model.ModelInstance, error)
	GetOnlineModelList(ctx context.Context) ([]*model.ModelInstance, error)
}

type modelConfigServiceImpl struct {
	repo repository.ModelConfigRepo
}

func NewModelConfigService(ctx context.Context) ModelConfigService {
	repo, err := repository.NewModelConfigRepo()
	if err != nil {
		logger.Errorf(ctx, "[ModelConfigService] failed to init model config repo: %v", err)
		return nil
	}

	return &modelConfigServiceImpl{repo: repo}
}

func (s *modelConfigServiceImpl) GetProviderModelList(ctx context.Context) ([]*model.ProviderModelList, error) {
	models, err := s.repo.Load(ctx)
	if err != nil {
		return nil, err
	}
	if models == nil {
		models = []*model.ModelInstance{}
	}

	result := make([]*model.ProviderModelList, 0, len(model.DefaultProviders))
	for _, p := range model.DefaultProviders {
		provider := p
		pml := &model.ProviderModelList{
			Provider:  &provider,
			ModelList: make([]*model.ModelInstance, 0),
		}

		for _, m := range models {
			if m.DeletedAt == 0 && m.ModelClass == p.ModelClass {
				pml.ModelList = append(pml.ModelList, m)
			}
		}
		result = append(result, pml)
	}
	return result, nil
}

func (s *modelConfigServiceImpl) CreateModel(ctx context.Context, req *model.CreateModelRequest) (int64, error) {
	if err := s.validateModelConfig(ctx, req); err != nil {
		return 0, err
	}

	inst := &model.ModelInstance{
		ModelClass:      req.ModelClass,
		DisplayName:     req.DisplayName,
		Connection:      req.Connection,
		ThinkingType:    req.ThinkingType,
		EnableBase64URL: req.EnableBase64URL,
		Status:          1,
	}

	if err := s.repo.Save(ctx, inst); err != nil {
		return 0, fmt.Errorf("save config failed: %w", err)
	}
	return inst.ID, nil
}

func (s *modelConfigServiceImpl) validateModelConfig(ctx context.Context, req *model.CreateModelRequest) error {
	if req == nil {
		return fmt.Errorf("invalid request: nil")
	}
	if req.Connection == nil {
		return fmt.Errorf("invalid request: connection is nil")
	}
	if req.Connection.Model == "" {
		return fmt.Errorf("invalid request: connection.model is empty")
	}

	validateCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cfg := &modelbuilder.ModelConfig{
		ModelClass:   req.ModelClass,
		Connection:   req.Connection,
		ThinkingType: req.ThinkingType,
	}

	chatModel, err := modelbuilder.BuildModel(validateCtx, cfg)
	if err != nil {
		return fmt.Errorf("build model failed: %w", err)
	}

	respMsg, err := chatModel.Generate(validateCtx, []*schema.Message{
		schema.SystemMessage("Just answer with a number, no explanation."),
		schema.UserMessage("1+1=?"),
	})
	if err != nil {
		return fmt.Errorf("generate failed: %w", err)
	}
	if respMsg == nil {
		return fmt.Errorf("generate failed: empty response")
	}
	return nil
}

func (s *modelConfigServiceImpl) DeleteModel(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return s.repo.Save(ctx, &model.ModelInstance{ID: id, DeletedAt: now, UpdatedAt: now})
}

func (s *modelConfigServiceImpl) GetModelByID(ctx context.Context, id int64) (*model.ModelInstance, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *modelConfigServiceImpl) GetOnlineModelList(ctx context.Context) ([]*model.ModelInstance, error) {
	models, err := s.repo.Load(ctx)
	if err != nil {
		return nil, err
	}
	if models == nil {
		models = []*model.ModelInstance{}
	}

	result := make([]*model.ModelInstance, 0)
	for _, m := range models {
		if m.DeletedAt == 0 && m.Status == 1 {
			result = append(result, m)
		}
	}
	return result, nil
}
