package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type ImgurStorage struct {
	clientID string
}

func NewImgurStorage(clientID string) *ImgurStorage {
	return &ImgurStorage{clientID: clientID}
}

func (s *ImgurStorage) UploadImage(file io.ReadSeeker, fileName string) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("image", fileName)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return "", err
	}
	w.Close()

	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Client-ID "+s.clientID)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		if link, ok := data["link"].(string); ok {
			return link, nil
		}
	}

	return "", fmt.Errorf("failed to upload image")
}
