package utils

import (
	"io"
	"net/http"
	"time"
)

func HttpGet(uri string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, err
	}

	return res.Body, nil
}
