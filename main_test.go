package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareMetricsInc(t *testing.T) {
	testCases := []struct {
		name         string
		initialHits  int32
		expectedHits int32
	}{
		{
			name:         "initial hits 0",
			initialHits:  0,
			expectedHits: 1,
		},
		{
			name:         "initial hits 100",
			initialHits:  100,
			expectedHits: 101,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &apiConfig{}
			cfg.fileserverHits.Store(tt.initialHits)

			handler := cfg.middlewareMetricsInc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/app/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Result().StatusCode)
			}

			if cfg.fileserverHits.Load() != tt.expectedHits {
				t.Errorf("expected fileserverHits to be %d, got %d", tt.expectedHits, cfg.fileserverHits.Load())
			}
		})
	}
}
