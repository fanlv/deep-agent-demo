package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/coder/acp-go-sdk"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	logFile := setupLogFile()
	if logFile != nil {
		defer logFile.Close()
	}

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("DeepAgent ACP Agent starting...")

	handler, err := newACPHandler(logger)
	if err != nil {
		logger.Error("Failed to create ACP handler", "error", err)
		os.Exit(1)
	}
	conn := acp.NewAgentSideConnection(handler, os.Stdout, os.Stdin)
	conn.SetLogger(logger)
	handler.SetAgentConnection(conn)

	logger.Info("ACP Agent ready, waiting for client connection...")

	select {
	case <-ctx.Done():
		logger.Info("Agent shutting down...")
	case <-conn.Done():
		logger.Info("Connection closed")
	}
}

func setupLogFile() *os.File {
	exePath, err := os.Executable()
	if err != nil {
		return os.Stderr
	}

	logName := fmt.Sprintf("acp_%s.log", time.Now().Format("20060102_150405"))
	logPath := filepath.Join(filepath.Dir(exePath), logName)
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return os.Stderr
	}
	return f
}
