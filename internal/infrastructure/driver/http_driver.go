package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HttpDriver interface {
	Get(ctx context.Context, endpoint string, params any, header any) ([]byte, error)
	Post(ctx context.Context, endpoint string, reqBody any, header any) ([]byte, error)
}

type httpDriver struct {
	httpClient *http.Client
}

func NewClient(baseURL string, c *http.Client) HttpDriver {
	return &httpDriver{
		httpClient: c,
	}
}

func (c *httpDriver) Get(ctx context.Context, endpoint string, params any, header any) ([]byte, error) {
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err

	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err

	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err

	}
	return body, nil
}
