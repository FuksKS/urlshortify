package handlers

import (
	"fmt"
	"github.com/FuksKS/urlshortify/internal/storage"
	"github.com/FuksKS/urlshortify/internal/urlmaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	practicumHost   = "https://practicum.yandex.ru/"
	defaultFilePath = "/tmp/short-url-db.json"
	defaultBaseURL  = "http://localhost:8080"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path, body string) (*http.Response, string) {
	var reqBody io.Reader = nil
	if method == http.MethodPost {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	st, _ := storage.New(defaultFilePath)

	handler := New(st, defaultBaseURL)
	ts := httptest.NewServer(handler.InitRouter())
	defer ts.Close()

	type want struct {
		statusCode  int
		contentType string
		location    string
		respBody    string
	}

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		st     Storager
		want   want
	}{
		{
			name:   "simple POST test",
			method: http.MethodPost,
			path:   "/",
			st:     st,
			body:   practicumHost,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				respBody:    defaultBaseURL + "/" + urlmaker.MakeShortURL(practicumHost),
			},
		},
	}

	for _, tt := range tests {
		resp, get := testRequest(t, ts, tt.method, tt.path, tt.body)
		resp.Body.Close()

		assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		assert.Equal(t, tt.want.respBody, get)
	}
}

func Test_shorten(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		respBody    string
	}

	st, _ := storage.New(defaultFilePath)

	handler := URLHandler{
		storage: st,
		BaseURL: defaultBaseURL,
	}

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		st     Storager
		want   want
	}{
		{
			name:   "simple test",
			method: http.MethodPost,
			path:   "/",
			st:     handler.storage,
			body:   practicumHost,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				respBody:    defaultBaseURL + "/" + urlmaker.MakeShortURL(practicumHost),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := handler.shorten()
			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.respBody, string(body))
		})
	}
}

func Test_getShorten(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}

	s, _ := storage.New(defaultFilePath)
	handler := URLHandler{
		storage: s,
		BaseURL: defaultBaseURL,
	}

	tests := []struct {
		name    string
		method  string
		request string
		st      Storager
		want    want
	}{
		{
			name:    "simple test",
			method:  http.MethodGet,
			request: "/" + urlmaker.MakeShortURL(practicumHost),
			st:      handler.storage,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   practicumHost,
			},
		},
		{
			name:    "Unknown request",
			method:  http.MethodGet,
			request: "/abc",
			st:      handler.storage,
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:    "Wrong path /abc/",
			method:  http.MethodGet,
			request: "/abc/",
			st:      handler.storage,
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:    "Wrong path /",
			method:  http.MethodGet,
			request: "/",
			st:      handler.storage,
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
	}

	handler.storage.SaveShortURL(urlmaker.MakeShortURL(practicumHost), practicumHost)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			h := handler.getShorten()
			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func Test_shortenJSON(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		respBody    string
	}

	s, _ := storage.New(defaultFilePath)
	handler := URLHandler{
		storage: s,
		BaseURL: defaultBaseURL,
	}

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		st     Storager
		want   want
	}{
		{
			name:   "simple test",
			method: http.MethodPost,
			path:   "/api/shorten",
			st:     handler.storage,
			body:   fmt.Sprintf(`{"url":"%s"}`, practicumHost),
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				respBody:    fmt.Sprintf(`{"result":"%s"}`, defaultBaseURL+"/"+urlmaker.MakeShortURL(practicumHost)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := handler.shortenJSON()
			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.respBody, string(body))
		})
	}
}
