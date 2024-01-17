package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	HTTPClient ClientInterface
	logger     *zap.SugaredLogger
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func NewClient(logger *zap.SugaredLogger) *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		logger:     logger,
	}
}

func (c *Client) DoRequest(req *http.Request) (*Response, error) {
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("[REQUEST ERROR] error while performing request: %w", err)
	}

	resBody, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("[REQUEST ERROR] error while reading the response bytes of the request: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       resBody,
	}, nil
}
