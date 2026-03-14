package repository

import (
	"encoding/json"
	"fmt"
	"log"

	sbmodel "github.com/deep-agent/sandbox/types/model"
	"github.com/fanlv/deep-agent-demo/services/agent/sandbox"
	"github.com/fanlv/deep-agent-demo/types/model"
	"github.com/fanlv/deep-agent-demo/types/path"
)

type SessionRepo interface {
	Save(sessionID string, meta *model.Session) error
	Load(sessionID string) (*model.Session, error)
	ListIDs() ([]string, error)
	LoadAll() ([]*model.Session, error)
}

type sessionRepo struct {
	sandbox *sandbox.Client
}

func NewSessionRepo() (SessionRepo, error) {
	sb, err := sandbox.New("")
	if err != nil {
		return nil, err
	}
	return &sessionRepo{sandbox: sb}, nil
}

func (r *sessionRepo) Save(sessionID string, meta *model.Session) error {
	sb, err := sandbox.New(sessionID)
	if err != nil {
		return fmt.Errorf("create sandbox failed: %w", err)
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal meta failed: %w", err)
	}

	metaPath := path.MetaFilePath(sb.Ctx.Workspace)
	dirPath := path.MetaDir(sb.Ctx.Workspace)

	_, err = sb.BashExecChecked(&sbmodel.BashExecRequest{
		Cwd:     sb.Ctx.Workspace,
		Command: fmt.Sprintf("[ -d %s ] || mkdir -p %s", dirPath, dirPath),
	})
	if err != nil {
		return fmt.Errorf("ensure deepagent dir failed: %w", err)
	}

	err = sb.Client.FileWrite(&sbmodel.FileWriteRequest{
		File:    metaPath,
		Content: string(data),
	})
	if err != nil {
		return fmt.Errorf("write meta file failed: %w", err)
	}

	return nil
}

func (r *sessionRepo) Load(sessionID string) (*model.Session, error) {
	metaPath := path.MetaFilePath(path.SessionDir(r.sandbox.Ctx.Workspace, sessionID))

	result, err := r.sandbox.Client.FileRead(&sbmodel.FileReadRequest{
		File: metaPath,
	})
	if err != nil {
		return nil, fmt.Errorf("read meta file failed: %w", err)
	}

	var meta model.Session
	if err := json.Unmarshal([]byte(result.Content), &meta); err != nil {
		return nil, fmt.Errorf("unmarshal meta failed: %w", err)
	}

	return &meta, nil
}

func (r *sessionRepo) ListIDs() ([]string, error) {
	workspace := path.SessionDir(r.sandbox.Ctx.Workspace, "")

	result, err := r.sandbox.Client.FileList(&sbmodel.FileListRequest{
		Path: workspace,
	})
	if err != nil {
		return nil, fmt.Errorf("list workspace failed: %w", err)
	}

	var sessionIDs []string
	for _, file := range result.Files {
		if !file.IsDir {
			continue
		}
		sessionID := file.Name
		metaPath := path.MetaFilePath(path.SessionDir(r.sandbox.Ctx.Workspace, sessionID))
		exists, err := r.sandbox.Client.FileExists(metaPath)
		if err != nil || !exists.Exists {
			continue
		}

		sessionIDs = append(sessionIDs, sessionID)
	}

	return sessionIDs, nil
}

func (r *sessionRepo) LoadAll() ([]*model.Session, error) {
	sessionIDs, err := r.ListIDs()
	if err != nil {
		return nil, err
	}

	var metas []*model.Session
	for _, sessionID := range sessionIDs {
		meta, err := r.Load(sessionID)
		if err != nil {
			log.Printf("[sessionRepo] load session %s failed: %v", sessionID, err)
			continue
		}
		if meta.Deleted {
			continue
		}
		metas = append(metas, meta)
		log.Printf("[sessionRepo] loaded session: %s, title: %s", meta.ID, meta.Title)
	}

	return metas, nil
}
