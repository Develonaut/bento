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
	"net/url"
	"os"
	"path/filepath"
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
//   - saveToFile (string, optional): Path to save response body to file (skips JSON parsing)
//   - queryParams (map[string]interface{}, optional): URL query parameters to append
//
// Returns a map with:
//   - statusCode (int): HTTP status code
//   - body (map[string]interface{}): Parsed JSON response (or empty if saveToFile is used)
//   - headers (map[string]string): Response headers
//   - filePath (string): Path where file was saved (only if saveToFile is used)
func (h *HTTPNeta) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Extract URL
	urlStr, ok := params["url"].(string)
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

	// Extract saveToFile path (optional)
	saveToFile, _ := params["saveToFile"].(string)

	// Add query parameters to URL (optional)
	if queryParams, ok := params["queryParams"].(map[string]interface{}); ok {
		var err error
		urlStr, err = addQueryParams(urlStr, queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to add query parameters: %w", err)
		}
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
	req, err := http.NewRequestWithContext(ctx, method, urlStr, reqBody)
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

	// Collect response headers
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	// If saveToFile is specified, save response to file instead of parsing JSON
	if saveToFile != "" {
		if err := saveResponseToFile(resp.Body, saveToFile); err != nil {
			return nil, fmt.Errorf("failed to save response to file: %w", err)
		}

		return map[string]interface{}{
			"statusCode": resp.StatusCode,
			"headers":    responseHeaders,
			"filePath":   saveToFile,
			"body":       map[string]interface{}{}, // Empty body when saving to file
		}, nil
	}

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

	// Build result
	result := map[string]interface{}{
		"statusCode": resp.StatusCode,
		"body":       responseBody,
		"headers":    responseHeaders,
	}

	return result, nil
}

// saveResponseToFile saves the HTTP response body to a file.
// Creates parent directories if they don't exist.
func saveResponseToFile(body io.Reader, filePath string) error {
	// Create parent directories if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	// Copy response body to file
	if _, err := io.Copy(file, body); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil
}

// addQueryParams adds query parameters to a URL.
// Returns error if URL parsing fails or if param values are unsupported types.
func addQueryParams(baseURL string, params map[string]interface{}) (string, error) {
	if len(params) == 0 {
		return baseURL, nil
	}

	// Parse the URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL %q: %w", baseURL, err)
	}

	// Get existing query params
	q := u.Query()

	// Add new params with type validation
	for key, value := range params {
		strVal := ""
		switch v := value.(type) {
		case string:
			strVal = v
		case int, int64, float64, bool:
			strVal = fmt.Sprintf("%v", v)
		case nil:
			// Skip nil values
			continue
		default:
			return "", fmt.Errorf("unsupported query param type for %q: %T", key, value)
		}
		q.Add(key, strVal)
	}

	// Set updated query params
	u.RawQuery = q.Encode()

	return u.String(), nil
}
