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
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.Init()
	logger.Log.Info("Инит конфига прошел", zap.String("init", "config Initialize"))

	if err := logger.Init(logger.LoggerLevelINFO); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "logger Initialize"))
	}

	logger.Log.Info("Инит логгера прошел", zap.String("init", "logger Initialize"))

	logger.Log.Info("Сейчас булдет инит бд", zap.String("init", "Db Initialize"), zap.String("cfg.DBDSN", cfg.DBDSN))
	st, err := storage.New(ctx, cfg.FileStorage, cfg.DBDSN)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "set storage"))
	}

	logger.Log.Info("Инит стораджа прошел", zap.String("init", "storage Initialize"), zap.String("cfg.DBDSN", cfg.DBDSN))

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

	if err = st.Shutdown(ctx); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "storage Shutdown"))
	}

	time.Sleep(2 * time.Second)

	logger.Log.Info("Terminated. Goodbye")
}
