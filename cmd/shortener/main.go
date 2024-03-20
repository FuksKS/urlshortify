package main

import (
	"context"
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
	"time"
)

func main() {
	cfg := config.Init()

	_, cancel := context.WithCancel(context.Background())

	if err := logger.Init(logger.LoggerLevelINFO); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "logger Initialize"))
	}

	var db pg.PgRepo
	var err error
	if cfg.DBDSN != "" {
		db, err = pg.NewConnect(cfg.DBDSN)
		if err != nil {
			logger.Log.Fatal(err.Error(), zap.String("init", "set db"))
		}
	}

	st, err := storage.New(db, cfg.FileStorage)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "set storage"))
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

	logger.Log.Info("Stop server", zap.String("address", cfg.HTTPAddr))
	if cfg.DBDSN == "" { // Записываем в файл только если нет бд
		if err = st.SaveCache(); err != nil {
			logger.Log.Fatal(err.Error(), zap.String("event", "save cache to storage"))
		}
	}

	time.Sleep(2 * time.Second)
	cancel()

	logger.Log.Info("Terminated. Goodbye")
}
