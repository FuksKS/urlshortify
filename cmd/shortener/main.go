package main

import (
	"github.com/FuksKS/urlshortify/internal/config"
	"github.com/FuksKS/urlshortify/internal/handlers"
	"github.com/FuksKS/urlshortify/internal/logger"
	"github.com/FuksKS/urlshortify/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

const defaultLoggerLevel = "INFO"

func main() {
	if err := logger.Init(defaultLoggerLevel); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "logger Initialize"))
	}
	st := storage.New()
	cfg := config.Init()

	fileProducer, err := storage.NewProducer(cfg.FileStorage)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "set file storage producer"))
	}
	fileConsumer, err := storage.NewConsumer(cfg.FileStorage)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "set file storage consumer"))
	}

	handler := handlers.New(st, fileProducer, fileConsumer, cfg.HTTPAddr, cfg.HTTPAddr)

	logger.Log.Info("Running server", zap.String("address", cfg.HTTPAddr))

	if err := http.ListenAndServe(handler.HTTPAddr, handler.InitRouter()); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "start server"))
	}
}
