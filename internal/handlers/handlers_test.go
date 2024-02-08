package handlers

import (
	"github.com/FuksKS/urlshortify/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var practicumHost = "https://practicum.yandex.ru/"

func testRequest(t *testing.T, ts *httptest.Server, method,
	path, body string) (*http.Response, string) {

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
	ts := httptest.NewServer(RootHandler("localhost:8080/", "b"))
	defer ts.Close()

	type want struct {
		statusCode  int
		contentType string
		location    string
		respBody    string
	}

	stor := storage.New()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		st     storage.Storager
		want   want
	}{
		{
			name:   "simple POST test",
			method: http.MethodPost,
			path:   "/",
			st:     stor,
			body:   practicumHost,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				respBody:    defaultHost + stor.SaveShortURL(practicumHost),
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

func Test_generateShortURL(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		respBody    string
	}

	stor := storage.New()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		st     storage.Storager
		want   want
	}{
		{
			name:   "simple test",
			method: http.MethodPost,
			path:   "/",
			st:     stor,
			body:   practicumHost,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				respBody:    defaultHost + stor.SaveShortURL(practicumHost),
			},
		},
		{
			name:   "simple test",
			method: http.MethodHead,
			path:   "/",
			st:     stor,
			body:   practicumHost,
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
				respBody:    "Method not Allowed\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := generateShortURL(tt.st, "localhost:8080/")
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

func Test_getURLID(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}

	stor := storage.New()

	tests := []struct {
		name    string
		method  string
		request string
		st      storage.Storager
		want    want
	}{
		{
			name:    "simple test",
			method:  http.MethodGet,
			request: "/" + stor.SaveShortURL(practicumHost),
			st:      stor,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   practicumHost,
			},
		},
		{
			name:    "Unknown request",
			method:  http.MethodGet,
			request: "/abc",
			st:      stor,
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:    "Wrong method",
			method:  http.MethodDelete,
			request: "/abc",
			st:      stor,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				location:   "",
			},
		},
		{
			name:    "Wrong path /abc/",
			method:  http.MethodGet,
			request: "/abc/",
			st:      stor,
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:    "Wrong path /",
			method:  http.MethodGet,
			request: "/",
			st:      stor,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				location:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			h := getURLID(tt.st)
			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
