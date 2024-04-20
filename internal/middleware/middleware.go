package middleware

import (
	"context"
	"errors"
	"github.com/FuksKS/urlshortify/internal/logger"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/token"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

const cookieName = "authToken"

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

// WithLogging — middleware-логер для входящих HTTP-запросов.
func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

// WithGzip - middleware поддерживающий gzip компрессию и декомпрессию
func WithGzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	})
}

type authResponseWriter struct {
	http.ResponseWriter
	UserID string
}

func (r *authResponseWriter) Write(b []byte) (int, error) {
	logger.Log.Info("WithAuth middleware. Write", zap.String("Ща запишу тело ответа ", string(b)))
	return r.ResponseWriter.Write(b)
}

func (r *authResponseWriter) WriteHeader(statusCode int) {
	authToken, err := token.MakeAuthToken(r.UserID)
	if err != nil {
		logger.Log.Error("WithAuth middleware. WriteHeader MakeAuthToken error", zap.String("error: ", err.Error()), zap.String("UserID: ", r.UserID))
	}

	cookie := http.Cookie{Name: cookieName, Value: authToken}
	http.SetCookie(r.ResponseWriter, &cookie)

	logger.Log.Info("WithAuth middleware. WriteHeader with cookie", zap.String("cookie name: ", cookie.Name), zap.String("cookie value: ", cookie.Value))

	http.SetCookie(r.ResponseWriter, &cookie)

	r.ResponseWriter.WriteHeader(statusCode)
}

// WithAuth - middleware авторизации
func WithAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenWithUser, err := r.Cookie("authToken")

		if tokenWithUser != nil {
			logger.Log.Info("WithAuth middleware. tokenWithUser != nil", zap.String(" tokenWithUser.Value", tokenWithUser.Value))
		} else {
			logger.Log.Info("WithAuth middleware. tokenWithUser == nil")
		}

		if err != nil {
			logger.Log.Info("WithAuth middleware. err from r.Cookie() != nil", zap.String("error: ", err.Error()))
		}

		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "Get cookie error", http.StatusInternalServerError)
		}

		var userID string
		if tokenWithUser != nil {
			if tokenWithUser.Value != "" {
				logger.Log.Info("WithAuth middleware. tokenWithUser.Value != ''", zap.String(" tokenWithUser.Value: ", tokenWithUser.Value))
				userID, err = token.GetUserID(tokenWithUser.Value)
				logger.Log.Info("WithAuth middleware. token.GetUserID", zap.String("userID: ", userID))
				if err != nil {
					http.Error(w, "Get userID from token error", http.StatusInternalServerError)
				}
				if userID == "" {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
				}
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}

		if userID == "" {
			userID = uuid.New().String()
			logger.Log.Info("Неизвестный юзер. Новый ID", zap.String("user_id", userID))
		} else {
			logger.Log.Info("Известный юзер. Старый ID", zap.String("user_id", userID))
		}

		aw := authResponseWriter{
			ResponseWriter: w,
			UserID:         userID,
		}

		userForContext := models.ContextKey(userID)
		ctx := context.WithValue(r.Context(), models.UserIDKey, userForContext)
		r = r.WithContext(ctx)

		h.ServeHTTP(&aw, r)
	})
}
