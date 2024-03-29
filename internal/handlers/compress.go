package handlers

import (
	"compress/gzip"
	"github.com/FuksKS/urlshortify/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	contentType := c.w.Header().Get("Content-Type")
	logger.Log.Info("withGzip middleware", zap.String("contentType", contentType))

	if contentType == "application/json" || contentType == "text/html" {
		size, err := c.zw.Write(p)
		c.zw.Close()
		return size, err
	}

	return c.w.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		contentType := c.w.Header().Get("Content-Type")
		if contentType == "application/json" || contentType == "text/html" {
			c.w.Header().Set("Content-Encoding", "gzip")
		}
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	// обнулим буфер gzip чтоб не дописывалось ничего в ответ
	c.zw.Reset(nil)
	return nil //return c.zw.Close() не работает т.к. запаникуем
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
