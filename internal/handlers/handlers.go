package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/logger"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/urlmaker"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func (h *URLHandler) getShorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := chi.URLParam(r, "id")

		oneURLInfo, err := h.storage.GetLongURL(ctx, id)
		if err != nil {
			http.Error(w, "GetLongURL error", http.StatusInternalServerError)
		}

		if oneURLInfo.IsDeleted {
			w.WriteHeader(http.StatusGone)
			return
		}

		if oneURLInfo.OriginalURL != "" {
			http.Redirect(w, r, oneURLInfo.OriginalURL, http.StatusTemporaryRedirect)
			return
		}

		// урла нет в хранилище
		http.Error(w, "Unknown short URL", http.StatusBadRequest)
	}
}

func (h *URLHandler) shorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")

		longURL := string(body)
		shortURL := urlmaker.MakeShortURL(longURL)

		userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
		if !ok {
			http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
		}

		oneURLInfo := models.URLInfo{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: longURL,
			UserID:      string(userID),
		}

		logger.Log.Info("Сохранение в базу из shorten()", zap.String("user_id", string(userID)), zap.String("short URL", shortURL), zap.String("original URL", longURL))
		err = h.storage.SaveShortURL(ctx, oneURLInfo)
		if err != nil && errors.Is(err, models.ErrAffectNoRows) {
			w.WriteHeader(http.StatusConflict)
		} else if err != nil {
			http.Error(w, "SaveShortURL error", http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		resp := fmt.Sprintf("%s/%s", h.BaseURL, shortURL)
		w.Write([]byte(resp))
	}
}

func (h *URLHandler) shortenJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		req := models.ShortenReq{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Unmarshal body error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		shortURL := urlmaker.MakeShortURL(req.URL)
		userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
		if !ok {
			http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
		}

		logger.Log.Info("Сохранение в базу из shortenJSON()", zap.String("user_id", string(userID)), zap.String("short URL", shortURL), zap.String("original URL", req.URL))
		oneURLInfo := models.URLInfo{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: req.URL,
			UserID:      string(userID),
		}

		err = h.storage.SaveShortURL(ctx, oneURLInfo)
		if err != nil && errors.Is(err, models.ErrAffectNoRows) {
			w.WriteHeader(http.StatusConflict)
		} else if err != nil {
			http.Error(w, "SaveShortURL error", http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		fullHost := fmt.Sprintf("%s/%s", h.BaseURL, shortURL)
		resp := models.ShortenResp{Result: fullHost}
		respB, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(respB)
	}
}

func (h *URLHandler) shortenBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		var req []models.URLInfo
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Unmarshal body error", http.StatusInternalServerError)
			return
		}

		userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
		if !ok {
			http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
		}

		for i := range req {
			req[i].ShortURL = urlmaker.MakeShortURL(req[i].OriginalURL)
			req[i].UserID = string(userID)

			logger.Log.Info("Сохранение в базу из shortenBatch()", zap.String("user_id", string(userID)), zap.String("short URL", req[i].ShortURL), zap.String("original URL", req[i].OriginalURL))
		}

		err = h.storage.SaveURLs(ctx, req)
		if err != nil {
			http.Error(w, "SaveURLs error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		resp := make([]models.URLInfo, len(req))
		for i := range req {
			resp[i].UUID = req[i].UUID
			resp[i].ShortURL = fmt.Sprintf("%s/%s", h.BaseURL, req[i].ShortURL)
		}

		respB, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(respB)
	}
}

func (h *URLHandler) getUsersShorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Предположим, что user_id сохранен в контексте под ключом "user_id"
		userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
		if !ok {
			http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
		}

		logger.Log.Info("Достаем из база с getUsersShorten()", zap.String("user_id", string(userID)))

		usersURLsInfo, err := h.storage.GetUsersURLs(ctx, string(userID))
		if err != nil {
			http.Error(w, "Get users URLs error", http.StatusInternalServerError)
		}

		if len(usersURLsInfo) == 0 {
			logger.Log.Info("Достали из базы с getUsersShorten() 0 урлов", zap.String("user_id", string(userID)))
			w.WriteHeader(http.StatusNoContent)
			return
		}
		logger.Log.Info("Достали из базы с getUsersShorten() нулевой короткий урл", zap.String("user_id", string(userID)), zap.String("short_url", usersURLsInfo[0].ShortURL))

		resp := make([]models.URLInfo, len(usersURLsInfo))
		for i := range usersURLsInfo {
			resp[i].OriginalURL = usersURLsInfo[i].OriginalURL
			resp[i].ShortURL = fmt.Sprintf("%s/%s", h.BaseURL, usersURLsInfo[i].ShortURL)
		}

		respB, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(respB)
	}
}

func (h *URLHandler) deleteShortenBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
		if !ok {
			http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		var shortURLs []string
		err = json.Unmarshal(body, &shortURLs)
		if err != nil {
			http.Error(w, "Unmarshal body error", http.StatusInternalServerError)
			return
		}

		// отправим сообщение в очередь на сохранение
		h.DeleteURLChan <- models.DeleteURLs{
			URLs:   shortURLs,
			UserID: string(userID),
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// flushMessages постоянно дропает несколько отправленных урлов в хранилище с определённым интервалом
func (h *URLHandler) flushDeleteURLs() {
	// будем дропать урлы, накопленные за последние 10 секунд
	ticker := time.NewTicker(10 * time.Second)

	var deleteURLs []models.DeleteURLs

	for {
		select {
		case urls := <-h.DeleteURLChan:
			// добавим урлы в слайс для последующего сохранения
			deleteURLs = append(deleteURLs, urls)
		case <-ticker.C:
			// подождём, пока придёт хотя бы одно сообщение
			if len(deleteURLs) == 0 {
				continue
			}
			// дропнем все пришедшие урлы одновременно
			err := h.storage.DeleteURLs(context.Background(), deleteURLs)
			if err != nil {
				logger.Log.Error("cannot DeleteURLs", zap.Error(err))
				// не будем стирать сообщения, попробуем отправить их чуть позже
				continue
			}
			// сотрём успешно дропнутые урлы
			deleteURLs = nil
		}
	}
}

func (h *URLHandler) pingDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := h.storage.PingDB(ctx); err != nil {
			http.Error(w, "Ping db", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
