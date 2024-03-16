package main

import (
	"github.com/FuksKS/urlshortify/internal/config"
	"github.com/FuksKS/urlshortify/internal/handlers"
	"github.com/FuksKS/urlshortify/internal/logger"
	"github.com/FuksKS/urlshortify/internal/pg"
	"github.com/FuksKS/urlshortify/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Init()

	if err := logger.Init(logger.LoggerLevelINFO); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "logger Initialize"))
	}

	st, err := storage.New(cfg.FileStorage)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "set storage"))
	}

	db, err := pg.Connect(cfg.DbDSN)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "set db"))
	}

	handler := handlers.New(st, db, cfg.HTTPAddr, cfg.HTTPAddr)

	logger.Log.Info("Running server", zap.String("address", cfg.HTTPAddr))

	go func() {
		if err := http.ListenAndServe(handler.HTTPAddr, handler.InitRouter()); err != nil {
			logger.Log.Fatal(err.Error(), zap.String("event", "start server"))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	if err = st.SaveCache(); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "save cache to storage"))
	}
}
