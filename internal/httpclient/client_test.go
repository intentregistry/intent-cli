package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/intentregistry/intent-cli/internal/config"
)

func TestClient_Get(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		
		// Check User-Agent header
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			t.Error("Expected User-Agent header")
		}
		
		// Check telemetry header if enabled
		telemetry := r.Header.Get("X-Telemetry-Enabled")
		if telemetry == "true" {
			// Telemetry is enabled
		}
		
		// Return test response
		response := map[string]interface{}{
			"status": "ok",
			"data":   "test response",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	// Test cases
	tests := []struct {
		name      string
		config    config.Config
		path      string
		telemetry bool
		wantError bool
	}{
		{
			name: "successful GET request",
			config: config.Config{
				APIURL: server.URL,
				Token:  "test-token",
			},
			path:      "/test",
			telemetry: false,
			wantError: false,
		},
		{
			name: "successful GET request with telemetry",
			config: config.Config{
				APIURL: server.URL,
				Token:  "test-token",
			},
			path:      "/test",
			telemetry: true,
			wantError: false,
		},
		{
			name: "GET request with authentication",
			config: config.Config{
				APIURL: server.URL,
				Token:  "test-token",
			},
			path:      "/authenticated",
			telemetry: false,
			wantError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewWithOptions(tt.config, false, tt.telemetry)
			
			var result map[string]interface{}
			err := client.Get(tt.path, &result)
			
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !tt.wantError {
				if result["status"] != "ok" {
					t.Errorf("Expected status 'ok', got %v", result["status"])
				}
			}
		})
	}
}

func TestClient_PostMultipart(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		// Check content type
		contentType := r.Header.Get("Content-Type")
		if !contains(contentType, "multipart/form-data") {
			t.Errorf("Expected multipart/form-data, got %s", contentType)
		}
		
		// Return test response
		response := map[string]interface{}{
			"status": "uploaded",
			"file":   "test.tar.gz",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	// Create a temporary file for testing
	tempFile := createTempFile(t, "test content")
	defer tempFile.Close()
	
	config := config.Config{
		APIURL: server.URL,
		Token:  "test-token",
	}
	
	client := NewWithOptions(config, false, false)
	
	fields := map[string]any{
		"name":        "test-intent",
		"description": "A test intent",
	}
	
	var result map[string]interface{}
	err := client.PostMultipart("/upload", fields, "file", tempFile.Name(), &result)
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if result["status"] != "uploaded" {
		t.Errorf("Expected status 'uploaded', got %v", result["status"])
	}
}

func TestClient_RetryLogic(t *testing.T) {
	attemptCount := 0
	
	// Create a test server that fails first few requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		
		// Fail first 2 requests, succeed on 3rd
		if attemptCount <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		response := map[string]interface{}{
			"status": "success after retry",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	config := config.Config{
		APIURL: server.URL,
		Token:  "test-token",
	}
	
	client := NewWithOptions(config, true, false) // Enable debug to see retry logs
	
	var result map[string]interface{}
	err := client.Get("/test", &result)
	
	if err != nil {
		t.Errorf("Unexpected error after retries: %v", err)
	}
	
	if result["status"] != "success after retry" {
		t.Errorf("Expected 'success after retry', got %v", result["status"])
	}
	
	// Verify we actually retried
	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts (2 failures + 1 success), got %d", attemptCount)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		server   func() *httptest.Server
		wantError bool
	}{
		{
			name: "HTTP error response",
			server: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Not found"))
				}))
			},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.server()
			defer server.Close()
			
			cfg := config.Config{
				APIURL: server.URL,
				Token:  "test-token",
			}
			
			client := NewWithOptions(cfg, false, false)
			
			var result map[string]interface{}
			err := client.Get("/test", &result)
			
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func createTempFile(t *testing.T, content string) *os.File {
	file, err := os.CreateTemp("", "intent-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	
	return file
}
