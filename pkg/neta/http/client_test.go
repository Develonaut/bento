package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Execute(t *testing.T) {
	tests := []struct {
		name       string
		params     map[string]interface{}
		serverResp string
		serverCode int
		wantErr    bool
	}{
		{
			name: "successful GET request",
			params: map[string]interface{}{
				"method": "GET",
				"url":    "TEST_SERVER_URL",
			},
			serverResp: `{"status":"ok"}`,
			serverCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "successful POST with body",
			params: map[string]interface{}{
				"method": "POST",
				"url":    "TEST_SERVER_URL",
				"body":   `{"name":"test"}`,
			},
			serverResp: `{"created":true}`,
			serverCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "with custom headers",
			params: map[string]interface{}{
				"method": "GET",
				"url":    "TEST_SERVER_URL",
				"headers": map[string]interface{}{
					"Accept":        "application/json",
					"Authorization": "Bearer token123",
				},
			},
			serverResp: `{"data":"test"}`,
			serverCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "missing url parameter",
			params: map[string]interface{}{
				"method": "GET",
			},
			wantErr: true,
		},
		{
			name: "default method is GET",
			params: map[string]interface{}{
				"url": "TEST_SERVER_URL",
			},
			serverResp: `{"default":"get"}`,
			serverCode: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverCode)
				w.Write([]byte(tt.serverResp))
			}))
			defer server.Close()

			// Replace TEST_SERVER_URL with actual server URL
			if url, ok := tt.params["url"].(string); ok && url == "TEST_SERVER_URL" {
				tt.params["url"] = server.URL
			}

			client := New()
			result, err := client.Execute(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error: %v", err)
				return
			}

			if result.Output == nil {
				t.Error("Execute() output is nil")
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		name       string
		params     map[string]interface{}
		wantMethod string
		wantErr    bool
	}{
		{
			name: "valid GET request",
			params: map[string]interface{}{
				"method": "GET",
				"url":    "https://example.com",
			},
			wantMethod: "GET",
			wantErr:    false,
		},
		{
			name: "missing url",
			params: map[string]interface{}{
				"method": "GET",
			},
			wantErr: true,
		},
		{
			name: "invalid header type",
			params: map[string]interface{}{
				"url": "https://example.com",
				"headers": map[string]interface{}{
					"Accept": 123, // Should be string
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := buildRequest(context.Background(), tt.params)
			if tt.wantErr {
				if err == nil {
					t.Error("buildRequest() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("buildRequest() unexpected error: %v", err)
				return
			}

			if tt.wantMethod != "" && req.Method != tt.wantMethod {
				t.Errorf("buildRequest() method = %s, want %s", req.Method, tt.wantMethod)
			}
		})
	}
}

func TestGetStringParam(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		key    string
		def    string
		want   string
	}{
		{
			name:   "existing string value",
			params: map[string]interface{}{"key": "value"},
			key:    "key",
			def:    "default",
			want:   "value",
		},
		{
			name:   "missing key returns default",
			params: map[string]interface{}{},
			key:    "missing",
			def:    "default",
			want:   "default",
		},
		{
			name:   "wrong type returns default",
			params: map[string]interface{}{"key": 123},
			key:    "key",
			def:    "default",
			want:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringParam(tt.params, tt.key, tt.def)
			if got != tt.want {
				t.Errorf("getStringParam() = %s, want %s", got, tt.want)
			}
		})
	}
}
