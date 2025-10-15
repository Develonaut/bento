// Package schemas provides validation schemas for built-in node types.
package schemas

import (
	"fmt"
	"strings"
)

// HTTPSchema validates HTTP node parameters.
type HTTPSchema struct{}

// NewHTTPSchema creates an HTTP schema validator.
func NewHTTPSchema() *HTTPSchema {
	return &HTTPSchema{}
}

// Validate checks HTTP node parameters.
func (s *HTTPSchema) Validate(params map[string]interface{}) error {
	var errs ValidationErrors

	errs = appendURLErrors(errs, params)
	errs = appendMethodErrors(errs, params)
	errs = appendHeaderErrors(errs, params)

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// appendURLErrors validates the url parameter.
func appendURLErrors(errs ValidationErrors, params map[string]interface{}) ValidationErrors {
	url, ok := params["url"].(string)
	if !ok || url == "" {
		errs = append(errs, ValidationError{
			Field:   "url",
			Message: "is required and must be a non-empty string",
		})
	}
	return errs
}

// appendMethodErrors validates the method parameter.
func appendMethodErrors(errs ValidationErrors, params map[string]interface{}) ValidationErrors {
	if method, ok := params["method"].(string); ok {
		if !isValidHTTPMethod(method) {
			errs = append(errs, ValidationError{
				Field:   "method",
				Message: "must be one of: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS",
			})
		}
	}
	return errs
}

// appendHeaderErrors validates the headers parameter.
func appendHeaderErrors(errs ValidationErrors, params map[string]interface{}) ValidationErrors {
	if headers, ok := params["headers"]; ok {
		if err := validateHeaders(headers); err != nil {
			errs = append(errs, ValidationError{
				Field:   "headers",
				Message: err.Error(),
			})
		}
	}
	return errs
}

// Fields returns field definitions for form generation.
func (s *HTTPSchema) Fields() []Field {
	return []Field{
		{
			Name:        "url",
			Type:        FieldString,
			Required:    true,
			Description: "HTTP(S) URL to request",
		},
		{
			Name:        "method",
			Type:        FieldString,
			Required:    false,
			Description: "HTTP method",
			Default:     "GET",
			Enum:        []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
		},
		{
			Name:        "headers",
			Type:        FieldMap,
			Required:    false,
			Description: "HTTP headers as key-value pairs",
		},
		{
			Name:        "body",
			Type:        FieldString,
			Required:    false,
			Description: "Request body (for POST/PUT/PATCH)",
		},
	}
}

// isValidHTTPMethod checks if method is a valid HTTP verb.
func isValidHTTPMethod(method string) bool {
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	upper := strings.ToUpper(method)
	for _, valid := range validMethods {
		if valid == upper {
			return true
		}
	}
	return false
}

// validateHeaders checks if headers is a valid map structure.
func validateHeaders(headers interface{}) error {
	headersMap, ok := headers.(map[string]interface{})
	if !ok {
		return fmt.Errorf("must be a map of string to string")
	}

	for k, v := range headersMap {
		if _, ok := v.(string); !ok {
			return fmt.Errorf("header %q value must be a string", k)
		}
	}

	return nil
}
