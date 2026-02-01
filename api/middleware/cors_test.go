package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000", "https://example.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}

	handler := CORS(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("allowed origin is set", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Errorf("expected origin http://localhost:3000, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("disallowed origin is not set", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "http://evil.com")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") == "http://evil.com" {
			t.Error("expected origin to not be set for disallowed origin")
		}
	})

	t.Run("no origin header defaults to wildcard", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("expected wildcard origin, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("preflight request returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("wildcard origin allows all", func(t *testing.T) {
		wildcardCfg := CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET"},
			AllowedHeaders: []string{"Content-Type"},
		}
		wildcardHandler := CORS(wildcardCfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "http://any-origin.com")
		rec := httptest.NewRecorder()

		wildcardHandler.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != "http://any-origin.com" {
			t.Errorf("expected origin to be allowed with wildcard config")
		}
	})
}

func TestCORSSimple(t *testing.T) {
	handler := CORSSimple(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("sets CORS headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("expected wildcard origin")
		}
		if rec.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("expected methods header to be set")
		}
	})

	t.Run("preflight request returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})
}
