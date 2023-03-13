package util

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func Post(url string, body []byte, headers map[string]string) ([]byte, error) {
	return Request("POST", url, body, headers)
}

func Get(url string, headers map[string]string) ([]byte, error) {
	return Request("GET", url, nil, headers)
}

func Request(method, url string, body []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for headerName, heanderValue := range headers {
		req.Header.Set(headerName, heanderValue)
	}

	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
