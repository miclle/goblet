package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fox-gonic/fox"
)

func TestRegisteredRoutes(t *testing.T) {
	router := fox.New()
	New(nil).RegisterRoutes(router)

	for _, tt := range []struct {
		path   string
		status int
		json   bool
	}{
		{path: "/health", status: http.StatusOK},
		{path: "/api/v1/hello", status: http.StatusOK, json: true},
		{path: "/api/v1/missing", status: http.StatusNotFound, json: true},
	} {
		t.Run(tt.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tt.status, response.Body.String())
			}
			if tt.json {
				if !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
					t.Fatalf("expected JSON response, got %s", response.Header().Get("Content-Type"))
				}
				var body map[string]any
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if tt.path == "/api/v1/hello" && body["message"] != "Hello from Goblet!" {
					t.Fatalf("unexpected greeting: %v", body)
				}
			} else if response.Body.String() != "ok" {
				t.Fatalf("unexpected health response: %q", response.Body.String())
			}
		})
	}
}
