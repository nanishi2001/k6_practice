package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDelayHandler(t *testing.T) {
	t.Run("GET /delay/10 returns delayed response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delay/10", nil)
		rec := httptest.NewRecorder()

		DelayHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp DelayResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.DelayMs != 10 {
			t.Errorf("expected delay_ms 10, got %d", resp.DelayMs)
		}
	})

	t.Run("GET /delay/invalid returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delay/invalid", nil)
		rec := httptest.NewRecorder()

		DelayHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("GET /delay/-1 returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delay/-1", nil)
		rec := httptest.NewRecorder()

		DelayHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("GET /delay/20000 returns bad request for exceeding max", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delay/20000", nil)
		rec := httptest.NewRecorder()

		DelayHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("POST /delay/10 returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/delay/10", nil)
		rec := httptest.NewRecorder()

		DelayHandler(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestRandomDelayHandler(t *testing.T) {
	t.Run("GET /delay/random returns response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delay/random", nil)
		rec := httptest.NewRecorder()

		RandomDelayHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp DelayResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.DelayMs < 0 || resp.DelayMs > 1000 {
			t.Errorf("expected delay_ms between 0-1000, got %d", resp.DelayMs)
		}
	})

	t.Run("POST /delay/random returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/delay/random", nil)
		rec := httptest.NewRecorder()

		RandomDelayHandler(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestErrorRateHandler(t *testing.T) {
	t.Run("GET /error-rate/0 always succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error-rate/0", nil)
		rec := httptest.NewRecorder()

		ErrorRateHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp ErrorRateResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !resp.Success {
			t.Error("expected success to be true for 0% error rate")
		}
	})

	t.Run("GET /error-rate/100 always fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error-rate/100", nil)
		rec := httptest.NewRecorder()

		ErrorRateHandler(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}

		var resp ErrorRateResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Success {
			t.Error("expected success to be false for 100% error rate")
		}
	})

	t.Run("GET /error-rate/invalid returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error-rate/invalid", nil)
		rec := httptest.NewRecorder()

		ErrorRateHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("GET /error-rate/-1 returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error-rate/-1", nil)
		rec := httptest.NewRecorder()

		ErrorRateHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("GET /error-rate/101 returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error-rate/101", nil)
		rec := httptest.NewRecorder()

		ErrorRateHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("POST /error-rate/50 returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/error-rate/50", nil)
		rec := httptest.NewRecorder()

		ErrorRateHandler(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}
