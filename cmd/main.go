package main

import (
	"log"
	"net/http"

	"github.com/Wayru-Network/network-services/apps/gateway/internal/infra"
	"github.com/Wayru-Network/network-services/apps/gateway/internal/server"
	"github.com/Wayru-Network/network-services/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	env, err := infra.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.Init(env.AppEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	srv, err := server.NewServer(env)
	if err != nil {
		logger.Error("Failed to create server", zap.Error(err))
	}

	// Start server
	logger.Info("Starting server on port", zap.Int("port", env.Port))
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("HTTP server error", zap.Error(err))
	}
}
