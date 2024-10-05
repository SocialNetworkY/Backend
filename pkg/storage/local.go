package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	// create folder if not exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &LocalStorage{basePath: basePath, baseURL: baseURL}, nil
}

func (s *LocalStorage) UploadImage(file io.ReadSeeker, fileName string) (string, error) {
	return s.UploadFile(file, fileName)
}

func (s *LocalStorage) UploadFile(file io.ReadSeeker, fileName string) (string, error) {
	fullPath := filepath.Join(s.basePath, fileName)
	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.baseURL, fileName), nil
}
