package storage

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSStorage struct {
	client     *storage.Client
	bucketName string
}

func NewGCSStorage(ctx context.Context, bucketName, credentialsFile string) (*GCSStorage, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}
	return &GCSStorage{client: client, bucketName: bucketName}, nil
}

func (s *GCSStorage) UploadImage(file io.ReadSeeker, fileName string) (string, error) {
	return s.UploadFile(file, fileName)
}

func (s *GCSStorage) UploadFile(file io.ReadSeeker, fileName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	wc := s.client.Bucket(s.bucketName).Object(fileName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	return wc.Attrs().MediaLink, nil
}
