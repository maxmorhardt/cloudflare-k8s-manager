package main

import (
	"log/slog"
	"os"
	"github.com/maxmorhardt/cloudflare-k8s-manager/internal/k8s"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("Application start")
	k8s.Run(logger)
}