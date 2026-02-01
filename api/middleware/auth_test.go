package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(userID int, email string, expired bool) string {
	var expiresAt time.Time
	if expired {
		expiresAt = time.Now().Add(-1 * time.Hour)
	} else {
		expiresAt = time.Now().Add(1 * time.Hour)
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JWTSecret)
	return tokenString
}

func TestAuth(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Auth(nextHandler)

	t.Run("valid token passes", func(t *testing.T) {
		token := generateTestToken(1, "test@example.com", false)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("invalid authorization format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token := generateTestToken(1, "test@example.com", true)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("returns claims from context", func(t *testing.T) {
		claims := &Claims{
			UserID: 1,
			Email:  "test@example.com",
		}
		ctx := context.WithValue(context.Background(), UserContextKey, claims)

		result := GetUserFromContext(ctx)

		if result == nil {
			t.Fatal("expected claims, got nil")
		}
		if result.UserID != 1 {
			t.Errorf("expected user_id 1, got %d", result.UserID)
		}
		if result.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", result.Email)
		}
	})

	t.Run("returns nil for empty context", func(t *testing.T) {
		result := GetUserFromContext(context.Background())

		if result != nil {
			t.Error("expected nil for empty context")
		}
	})

	t.Run("returns nil for wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserContextKey, "wrong type")

		result := GetUserFromContext(ctx)

		if result != nil {
			t.Error("expected nil for wrong type")
		}
	})
}
