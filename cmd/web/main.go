package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/fanlv/deep-agent-demo/cmd/web/handler"
	"github.com/fanlv/deep-agent-demo/pkg/logger"
	"github.com/hertz-contrib/cors"
)

const serverPort = ":8090"

func main() {
	ctx := context.Background()

	h, err := handler.NewHandler(ctx)
	if err != nil {
		logger.Fatalf(ctx, "Failed to initialize handler: %v", err)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	s := newServer()
	registerRoutes(s, h)

	go s.Spin()
	logger.Infof(ctx, "Server is running on %s", serverPort)

	<-ctx.Done()
	logger.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		logger.Errorf(ctx, "Server shutdown error: %v", err)
	}
	logger.Info("Server stopped")
}

func newServer() *server.Hertz {
	h := server.Default(
		server.WithHostPorts(serverPort),
		server.WithExitWaitTime(5*time.Second),
	)

	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	h.Use(loggerMiddleware())

	return h
}
