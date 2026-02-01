package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"k6-practice/api/middleware"
	"k6-practice/api/models"
)

func setupAuthHandler() *AuthHandler {
	store := models.NewUserStore()
	return NewAuthHandler(store)
}

func TestAuthHandler_Login(t *testing.T) {
	handler := setupAuthHandler()

	t.Run("POST /auth/login with valid credentials returns tokens", func(t *testing.T) {
		body := LoginRequest{
			Email:    "alice@example.com",
			Password: "password",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp TokenResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.AccessToken == "" {
			t.Error("expected access_token to be set")
		}
		if resp.RefreshToken == "" {
			t.Error("expected refresh_token to be set")
		}
		if resp.ExpiresIn != 900 {
			t.Errorf("expected expires_in 900, got %d", resp.ExpiresIn)
		}
	})

	t.Run("POST /auth/login with invalid password returns unauthorized", func(t *testing.T) {
		body := LoginRequest{
			Email:    "alice@example.com",
			Password: "wrongpassword",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("POST /auth/login with unknown email returns unauthorized", func(t *testing.T) {
		body := LoginRequest{
			Email:    "unknown@example.com",
			Password: "password",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("POST /auth/login with invalid JSON returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("GET /auth/login returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	handler := setupAuthHandler()

	t.Run("POST /auth/refresh with valid token returns new tokens", func(t *testing.T) {
		// まずログインしてトークンを取得
		loginBody := LoginRequest{
			Email:    "alice@example.com",
			Password: "password",
		}
		jsonBody, _ := json.Marshal(loginBody)
		loginReq := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonBody))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		handler.ServeHTTP(loginRec, loginReq)

		var loginResp TokenResponse
		json.NewDecoder(loginRec.Body).Decode(&loginResp)

		// リフレッシュトークンで新しいトークンを取得
		refreshBody := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			RefreshToken: loginResp.RefreshToken,
		}
		jsonBody, _ = json.Marshal(refreshBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp TokenResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.AccessToken == "" {
			t.Error("expected access_token to be set")
		}
	})

	t.Run("POST /auth/refresh with invalid token returns unauthorized", func(t *testing.T) {
		body := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			RefreshToken: "invalid-token",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})
}

func TestAuthHandler_Me(t *testing.T) {
	handler := setupAuthHandler()

	t.Run("GET /auth/me with valid context returns user info", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)

		// コンテキストにユーザー情報を設定
		claims := &middleware.Claims{
			UserID: 1,
			Email:  "alice@example.com",
		}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp MeResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.UserID != 1 {
			t.Errorf("expected user_id 1, got %d", resp.UserID)
		}
		if resp.Email != "alice@example.com" {
			t.Errorf("expected email 'alice@example.com', got '%s'", resp.Email)
		}
	})

	t.Run("GET /auth/me without context returns unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("POST /auth/me returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestAuthHandler_NotFound(t *testing.T) {
	handler := setupAuthHandler()

	t.Run("GET /auth/unknown returns not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/unknown", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}
