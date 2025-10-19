package mocks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// NewFigmaServer creates a mock Figma API server.
// Returns URLs for test images like real Figma API.
//
// Real Figma API response format:
//
//	{
//	  "images": {
//	    "node-id": "https://figma-alpha-api.s3.us-west-2.amazonaws.com/..."
//	  }
//	}
//
// This mock returns the same structure for testing.
func NewFigmaServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check auth header (real Figma API requires X-Figma-Token)
		if r.Header.Get("X-Figma-Token") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Missing Figma token",
			}); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Simulate Figma API response
		// Real API returns: {"images": {"node-id": "https://..."}}
		response := map[string]interface{}{
			"images": map[string]interface{}{
				"test-component": "http://localhost:9999/mock-overlay.png",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}))
}
