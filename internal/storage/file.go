package storage

import (
	"bufio"
	"encoding/json"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/google/uuid"
	"os"
	"strings"
	"sync"
)

type Producer struct {
	file      *os.File
	writer    *bufio.Writer
	fileMutex *sync.Mutex
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:      file,
		writer:    bufio.NewWriter(file),
		fileMutex: &sync.Mutex{},
	}, nil
}

func (p *Producer) WriteToFile(info models.URLInfo) error {
	info.UUID = uuid.New().String()
	data, err := json.Marshal(&info)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}
	// добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	p.fileMutex.Lock()
	err = p.writer.Flush()
	p.fileMutex.Unlock()

	return err
}

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadFromFile(shortURL string) (string, error) {
	if !c.scanner.Scan() {
		return "", c.scanner.Err()
	}

	for c.scanner.Scan() {
		line := c.scanner.Text()
		if strings.Contains(line, shortURL) {
			info := models.URLInfo{}
			err := json.Unmarshal([]byte(line), &info)
			if err != nil {
				return "", err
			}
			return info.OriginalURL, nil
		}
	}

	return "", nil
}
