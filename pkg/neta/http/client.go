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

	body, err := readResponseBody(resp)
	if err != nil {
		return neta.Result{}, err
	}

	// Check for HTTP error status codes
	if resp.StatusCode >= 400 {
		return neta.Result{}, fmt.Errorf("HTTP %d: %s - %s", resp.StatusCode, resp.Status, body)
	}

	return neta.Result{Output: body}, nil
}

// readResponseBody reads and returns the response body as a string.
func readResponseBody(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// buildRequest creates an HTTP request from parameters.
func buildRequest(ctx context.Context, params map[string]interface{}) (*http.Request, error) {
	method := neta.GetStringParam(params, "method", "GET")
	url := neta.GetStringParam(params, "url", "")
	if url == "" {
		return nil, fmt.Errorf("url parameter required")
	}

	body := neta.GetStringParam(params, "body", "")
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
