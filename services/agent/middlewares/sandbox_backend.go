package middlewares

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk/middlewares/agentsmd"
	"github.com/cloudwego/eino/adk/middlewares/plantask"
	sandbox "github.com/deep-agent/sandbox/sdk/go"
	"github.com/deep-agent/sandbox/types/model"
)

type sandboxBackend struct {
	client sandbox.Sandbox
}

func (b *sandboxBackend) LsInfo(ctx context.Context, req *plantask.LsInfoRequest) ([]plantask.FileInfo, error) {
	resp, err := b.client.BashExec(&model.BashExecRequest{
		Command: fmt.Sprintf("[ -d %s ] || mkdir -p %s", req.Path, req.Path),
	})
	if err != nil {
		return nil, fmt.Errorf("[LsInfo]  check and create dir failed, path: %s, err: %w", req.Path, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("[LsInfo] check and create dir failed, path: %s, resp is nil", req.Path)
	}
	if resp.ExitCode != 0 {
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
		return nil, fmt.Errorf("[LsInfo] check and create dir failed, path: %s, exit_code: %d, msg: %s", req.Path, resp.ExitCode, msg)
	}

	result, err := b.client.FileList(&model.FileListRequest{
		Path: req.Path,
	})
	if err != nil {
		return nil, fmt.Errorf("[LsInfo] file list failed, path: %s, err: %w", req.Path, err)
	}

	fileInfos := make([]plantask.FileInfo, 0, len(result.Files))
	for _, f := range result.Files {
		fileInfos = append(fileInfos, plantask.FileInfo{
			Path:       f.Path,
			IsDir:      f.IsDir,
			Size:       f.Size,
			ModifiedAt: time.Unix(f.ModTimeUnix, 0).Format(time.RFC3339),
		})
	}
	return fileInfos, nil
}

func (b *sandboxBackend) Read(ctx context.Context, req *plantask.ReadRequest) (*agentsmd.FileContent, error) {
	result, err := b.client.FileRead(&model.FileReadRequest{
		File: req.FilePath,
	})
	if err != nil {
		return nil, fmt.Errorf("[Read] file read failed, path: %s, err: %w", req.FilePath, err)
	}
	return &agentsmd.FileContent{
		Content: result.Content,
	}, nil
}

func (b *sandboxBackend) Write(ctx context.Context, req *plantask.WriteRequest) error {
	err := b.client.FileWrite(&model.FileWriteRequest{
		File:    req.FilePath,
		Content: req.Content,
	})
	if err != nil {
		return fmt.Errorf("[Write] file write failed, path: %s, err: %w", req.FilePath, err)
	}
	return nil
}

func (b *sandboxBackend) Delete(ctx context.Context, req *plantask.DeleteRequest) error {
	err := b.client.FileDelete(&model.FileDeleteRequest{
		Path: req.FilePath,
	})
	if err != nil {
		return fmt.Errorf("[Delete] file delete failed, path: %s, err: %w", req.FilePath, err)
	}
	return nil
}
