package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type HttpClient struct {
	httpClient *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		httpClient: &http.Client{Timeout: time.Duration(20) * time.Second},
	}
}

func (h *HttpClient) Get(ctx context.Context, url string, output interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	res, err := h.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("http status code: %d | URL: %s", res.StatusCode, url))
	}

	err = json.NewDecoder(res.Body).Decode(&output)

	if err != nil {
		return err
	}

	return nil
}
