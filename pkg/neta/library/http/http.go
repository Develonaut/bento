// Package http provides HTTP request functionality for the bento workflow system.
//
// The http neta allows you to make HTTP requests (GET, POST, PUT, DELETE, etc.)
// to external APIs and services. It supports:
//   - Custom headers (including authentication)
//   - Request timeouts
//   - JSON request/response bodies
//   - Error handling for 4xx/5xx responses
//
// Example usage:
//
//	params := map[string]interface{}{
//	    "url": "https://api.example.com/users",
//	    "method": "GET",
//	    "headers": map[string]interface{}{
//	        "Authorization": "Bearer token123",
//	    },
//	    "timeout": 30,  // seconds
//	}
//
//	result, err := httpNeta.Execute(ctx, params)
//
// The result contains:
//   - statusCode: HTTP status code (200, 404, etc.)
//   - body: Parsed JSON response body
//   - headers: Response headers
//
// Learn more about Go's net/http package: https://pkg.go.dev/net/http
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Develonaut/bento/pkg/neta"
)

const (
	// DefaultTimeout is the default HTTP request timeout in seconds.
	DefaultTimeout = 30
)

// HTTPNeta implements HTTP request functionality.
type HTTPNeta struct{}

// New creates a new HTTP neta instance.
func New() neta.Executable {
	return &HTTPNeta{}
}

// Execute performs an HTTP request based on the provided parameters.
//
// Parameters:
//   - url (string, required): The URL to request
//   - method (string, required): HTTP method (GET, POST, PUT, DELETE, etc.)
//   - headers (map[string]interface{}, optional): Custom headers
//   - body (map[string]interface{}, optional): Request body (will be JSON encoded)
//   - timeout (int, optional): Request timeout in seconds (default: 30)
//
// Returns a map with:
//   - statusCode (int): HTTP status code
//   - body (map[string]interface{}): Parsed JSON response
//   - headers (map[string]string): Response headers
func (h *HTTPNeta) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Extract URL
	url, ok := params["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url parameter is required and must be a string")
	}

	// Extract method
	method, ok := params["method"].(string)
	if !ok {
		return nil, fmt.Errorf("method parameter is required and must be a string")
	}

	// Extract timeout (optional, default 30 seconds)
	timeout := DefaultTimeout
	if t, ok := params["timeout"].(int); ok {
		timeout = t
	}

	// Prepare request body
	var reqBody io.Reader
	if body, ok := params["body"].(map[string]interface{}); ok {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default Content-Type for JSON
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	if headers, ok := params["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var responseBody map[string]interface{}
	if len(respBodyBytes) > 0 {
		if err := json.Unmarshal(respBodyBytes, &responseBody); err != nil {
			// If JSON parsing fails, return raw body as string
			responseBody = map[string]interface{}{
				"raw": string(respBodyBytes),
			}
		}
	}

	// Collect response headers
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	// Build result
	result := map[string]interface{}{
		"statusCode": resp.StatusCode,
		"body":       responseBody,
		"headers":    responseHeaders,
	}

	return result, nil
}
