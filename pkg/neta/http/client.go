// Package http provides HTTP request execution nodes.
package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"bento/pkg/neta"
)

// Client executes HTTP requests.
type Client struct {
	client *http.Client
}

// New creates a new HTTP client node.
func New() *Client {
	return &Client{
		client: &http.Client{},
	}
}

// Execute performs an HTTP request.
func (c *Client) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	req, err := buildRequest(ctx, params)
	if err != nil {
		return neta.Result{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return neta.Result{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return neta.Result{}, err
	}

	return neta.Result{
		Output: string(body),
	}, nil
}

// buildRequest creates an HTTP request from parameters.
func buildRequest(ctx context.Context, params map[string]interface{}) (*http.Request, error) {
	method := getStringParam(params, "method", "GET")
	url := getStringParam(params, "url", "")
	if url == "" {
		return nil, fmt.Errorf("url parameter required")
	}

	body := getStringParam(params, "body", "")
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	if err := addHeaders(req, params); err != nil {
		return nil, err
	}

	return req, nil
}

// addHeaders adds headers from params to the request.
func addHeaders(req *http.Request, params map[string]interface{}) error {
	headers, ok := params["headers"].(map[string]interface{})
	if !ok {
		return nil
	}

	for k, v := range headers {
		strVal, ok := v.(string)
		if !ok {
			return fmt.Errorf("header %s value must be string", k)
		}
		req.Header.Set(k, strVal)
	}
	return nil
}

// getStringParam extracts a string parameter with default.
func getStringParam(params map[string]interface{}, key, defaultVal string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultVal
}
