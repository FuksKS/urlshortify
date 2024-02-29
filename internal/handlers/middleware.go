package handlers

import (
	"github.com/FuksKS/urlshortify/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// withLogging — middleware-логер для входящих HTTP-запросов.
func withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("got incoming HTTP request",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	}
}

// withGzip - middleware поддерживающий gzip компрессию и декомпрессию
func withGzip(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
			contentType := r.Header.Get("Content-Type")
			logger.Log.Info("withGzip middleware", zap.String("contentType", contentType))
			if contentType != "application/json" && contentType != "text/html" {
				h.ServeHTTP(w, r)
				return
			}
		*/

		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		clientSupportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		logger.Log.Info("withGzip middleware", zap.String("Accept-Encoding", r.Header.Get("Accept-Encoding")))
		if clientSupportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		clientSentGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if clientSentGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Add gzip compress error", http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)

	}
}
