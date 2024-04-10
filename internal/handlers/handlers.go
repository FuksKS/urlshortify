package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/urlmaker"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func (h *URLHandler) getShorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		oneURLInfo := h.storage.GetLongURL(id)

		if oneURLInfo != "" {
			http.Redirect(w, r, oneURLInfo, http.StatusTemporaryRedirect)
			return
		}

		// урла нет в хранилище
		http.Error(w, "Unknown short URL", http.StatusBadRequest)
	}
}

func (h *URLHandler) shorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")

		longURL := string(body)
		shortURL := urlmaker.MakeShortURL(longURL)

		/*
			userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
			if !ok {
				http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
			}

		*/
		userID := "1"

		oneURLInfo := models.URLInfo{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: longURL,
			UserID:      string(userID),
		}

		err = h.storage.SaveShortURL(oneURLInfo)
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
		//ctx := r.Context()

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
		/*
			userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
			if !ok {
				http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
			}

		*/
		userID := "1"

		oneURLInfo := models.URLInfo{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: req.URL,
			UserID:      string(userID),
		}

		err = h.storage.SaveShortURL(oneURLInfo)
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
		//ctx := r.Context()

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

		/*
			userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
			if !ok {
				http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
			}

		*/
		userID := "1"

		for i := range req {
			req[i].UUID = uuid.New().String()
			req[i].ShortURL = urlmaker.MakeShortURL(req[i].OriginalURL)
			req[i].UserID = string(userID)
		}

		err = h.storage.SaveURLs(req)
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
		//ctx := r.Context()

		/*
			userID, ok := ctx.Value(models.UserIDKey).(models.ContextKey)
			if !ok {
				http.Error(w, "Get user_id from context error", http.StatusInternalServerError)
			}

		*/
		userID := "1"

		usersURLsInfo, err := h.storage.GetUsersURLs(string(userID))
		if err != nil {
			http.Error(w, "Get users URLs error", http.StatusInternalServerError)
		}

		if len(usersURLsInfo) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp := make([]models.URLInfo, 0, len(usersURLsInfo))
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

func (h *URLHandler) pingDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.db.DB.Ping(context.Background()); err != nil {
			http.Error(w, "Ping db", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
