package path

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// metaDir is the directory name for deepagent data
	metaDir  = ".meta"
	agentDir = "agent"

	// metaFile is the filename for session metadata
	metaFile = "meta.json"

	// messagesFile is the filename for chat messages
	messagesFile = "messages.jsonl"

	// summaryFile is the filename for summary message
	summaryFile = "summary.json"

	// modelsFile is the filename for model configurations
	modelsFile = "models.json"
)

func ReductionDir(workspace string) string {
	return filepath.Join(workspace, "reduction")
}

// SessionDir returns the workspace path for a session in sandbox
func SessionDir(workspace string, sessionID string) string {
	return filepath.Join(workspace, sessionID)
}

// MetaDir returns the .deepagent directory path within a workspace
func MetaDir(sessionDir string) string {
	return filepath.Join(sessionDir, metaDir)
}

// MetaFilePath returns the meta.json file path within a workspace
func MetaFilePath(sessionDir string) string {
	return filepath.Join(sessionDir, metaDir, metaFile)
}

// MessagesFilePath returns the messages.jsonl file path within a workspace
func MessagesFilePath(sessionDir string) string {
	return filepath.Join(sessionDir, metaDir, messagesFile)
}

// SummaryFilePath returns the summary.json file path within a workspace
func SummaryFilePath(sessionDir string) string {
	return filepath.Join(sessionDir, metaDir, summaryFile)
}

func AgentDir() (string, error) {
	if ws := os.Getenv("LOCAL_MEMORY"); ws != "" {
		return filepath.Join(ws, agentDir), nil
	}

	return "", fmt.Errorf("LOCAL_MEMORY env var is not set")
}

// ModelsConfigFile returns the models.json file path in user's home .deepagent directory
func ModelsConfigFile() (string, error) {
	dir, err := AgentDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, modelsFile), nil
}

// PromptsDir returns the prompts directory path within AgentDir
func PromptsDir() (string, error) {
	dir, err := AgentDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "prompts"), nil
}
