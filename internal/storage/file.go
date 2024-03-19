package storage

import (
	"bufio"
	"encoding/json"
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
		return nil, err
	}

	return &fileWriter{file: file, writer: bufio.NewWriter(file)}, nil
}

func (f *fileWriter) Save(cache map[string]string) error {
	for shortURL, longURL := range cache {
		data := models.URLInfo{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: longURL,
		}

		dataByte, err := json.Marshal(&data)
		if err != nil {
			return err
		}

		if _, err := f.writer.Write(dataByte); err != nil {
			return err
		}
		// добавляем перенос строки
		if err := f.writer.WriteByte('\n'); err != nil {
			return err
		}

		// записываем буфер в файл
		err = f.writer.Flush()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fileWriter) SaveURLs(urls []models.URLInfo) error {
	// В файл построчно не сохраняем. Метод просто для имплиментации интерфейса
	return nil
}

func (f *fileWriter) SaveOneURL(info models.URLInfo) error {
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
		return nil, err
	}

	return &FileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *FileReader) Read() (map[string]string, error) {
	/*
		if !fileConsumer.scanner.Scan() {
			fmt.Println(" !fileConsumer.scanner.Scan()")
			return nil, fileConsumer.scanner.Err()
		}

	*/

	m := make(map[string]string)
	for r.scanner.Scan() {
		line := r.scanner.Text()

		info := models.URLInfo{}
		err := json.Unmarshal([]byte(line), &info)
		if err != nil {
			return nil, err
		}

		m[info.ShortURL] = info.OriginalURL
	}

	return m, nil
}
