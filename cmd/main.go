package main

import (
	"shortUrl/internal/config"
	"shortUrl/internal/database/db"
	"shortUrl/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log = log.With(zap.String("env", cfg.Env))

	log.Info("Init server", zap.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	_, err := db.NewDB(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize db", zap.Error(err))
	}

	enableMiddlewares(log)

}

func NewDiscardLogger() *zap.Logger {
	logger := zap.NewNop()
	defer logger.Sync()

	return logger
}

func enableMiddlewares(log *zap.Logger) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}
