package main

import (
	"context"
	"github.com/FuksKS/urlshortify/internal/config"
	"github.com/FuksKS/urlshortify/internal/handlers"
	"github.com/FuksKS/urlshortify/internal/logger"
	"github.com/FuksKS/urlshortify/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.Init()

	if err := logger.Init(logger.LoggerLevelINFO); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "logger Initialize"))
	}

	st, err := storage.New(ctx, cfg.FileStorage, cfg.DBDSN)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "set storage"))
	}

	handler, err := handlers.New(st, cfg.BaseURL)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "set handler"))
	}

	logger.Log.Info("Running server", zap.String("address", cfg.HTTPAddr))

	go func() {
		if err := http.ListenAndServe(cfg.HTTPAddr, handler.InitRouter()); err != nil {
			logger.Log.Fatal(err.Error(), zap.String("event", "start server"))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	logger.Log.Info("Stop server", zap.String("address", cfg.HTTPAddr))

	if err = st.SaveCache(ctx); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "save cache to storage"))
	}

	cancel()

	logger.Log.Info("Terminated. Goodbye")
}
