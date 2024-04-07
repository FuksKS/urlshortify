package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/google/uuid"
	"os"
)

type fileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func newFileWriter(filename string) (*fileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("storage-newFileWriter-OpenFile-err: %w", err)
	}

	return &fileWriter{file: file, writer: bufio.NewWriter(file)}, nil
}

func (f *fileWriter) Save(_ context.Context, cache map[string]models.URLInfo) error {
	for shortURL, info := range cache {
		data := models.URLInfo{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: info.OriginalURL,
			UserID:      info.UserID,
		}

		dataByte, err := json.Marshal(&data)
		if err != nil {
			return fmt.Errorf("storage-fileWriter-Save-Marshal-err: %w", err)
		}

		if _, err := f.writer.Write(dataByte); err != nil {
			return fmt.Errorf("storage-fileWriter-Save-Write-data-err: %w", err)
		}
		// добавляем перенос строки
		if err := f.writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("storage-fileWriter-Save-Write-err: %w", err)
		}

		// записываем буфер в файл
		err = f.writer.Flush()
		if err != nil {
			return fmt.Errorf("storage-fileWriter-Save-Flush-err: %w", err)
		}
	}

	return nil
}

func (f *fileWriter) SaveURLs(_ context.Context, _ []models.URLInfo) error {
	// В файл построчно не сохраняем. Метод просто для имплиментации интерфейса
	return nil
}

func (f *fileWriter) SaveOneURL(_ context.Context, _ models.URLInfo) error {
	// В файл построчно не сохраняем. Метод просто для имплиментации интерфейса
	return nil
}

type FileReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func newFileReader(filePath string) (*FileReader, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("storage-newFileReader-OpenFile-err: %w", err)
	}

	return &FileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *FileReader) ReadAll(_ context.Context) (map[string]models.URLInfo, error) {

	m := make(map[string]models.URLInfo)
	for r.scanner.Scan() {
		line := r.scanner.Text()

		info := models.URLInfo{}
		err := json.Unmarshal([]byte(line), &info)
		if err != nil {
			return nil, fmt.Errorf("storage-FileReader-Unmarshal-err: %w", err)
		}

		m[info.ShortURL] = info
	}

	return m, nil
}

func (r *FileReader) GetLongURL(_ context.Context, _ string) (models.URLInfo, error) {
	// Для имплементации. Один урл берем только из кэша
	return models.URLInfo{}, nil
}

func (r *FileReader) GetUsersURLs(_ context.Context, _ string) ([]models.URLInfo, error) {
	// Для имплементации. Забираем урлы по юзеру всегда из бд
	return nil, nil
}

func (r *FileReader) PingDB(_ context.Context) error {
	// Для имплементации
	return nil
}

func (r *FileReader) Shutdown(_ context.Context) error {
	// Для имплементации
	return nil
}
