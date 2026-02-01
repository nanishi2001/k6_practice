package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"k6-practice/api/models"
)

func setupUsersHandler() *UsersHandler {
	store := models.NewUserStore()
	return NewUsersHandler(store)
}

func TestUsersHandler_List(t *testing.T) {
	handler := setupUsersHandler()

	t.Run("GET /users returns user list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var users []*models.User
		if err := json.NewDecoder(rec.Body).Decode(&users); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		// 初期データは3件
		if len(users) != 3 {
			t.Errorf("expected 3 users, got %d", len(users))
		}
	})
}

func TestUsersHandler_Get(t *testing.T) {
	handler := setupUsersHandler()

	t.Run("GET /users/1 returns user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var user models.User
		if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if user.ID != 1 {
			t.Errorf("expected user ID 1, got %d", user.ID)
		}
	})

	t.Run("GET /users/999 returns not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("GET /users/invalid returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestUsersHandler_Create(t *testing.T) {
	t.Run("POST /users creates user", func(t *testing.T) {
		handler := setupUsersHandler()

		body := CreateUserRequest{
			Name:  TestNewUserName,
			Email: TestNewUserEmail,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		var user models.User
		if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if user.Name != TestNewUserName {
			t.Errorf("expected name '%s', got '%s'", TestNewUserName, user.Name)
		}
		if user.Email != TestNewUserEmail {
			t.Errorf("expected email '%s', got '%s'", TestNewUserEmail, user.Email)
		}
	})

	t.Run("POST /users with missing fields returns bad request", func(t *testing.T) {
		handler := setupUsersHandler()

		body := CreateUserRequest{
			Name: TestNewUserName,
			// Email missing
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("POST /users with invalid JSON returns bad request", func(t *testing.T) {
		handler := setupUsersHandler()

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("POST /users with invalid email returns bad request", func(t *testing.T) {
		handler := setupUsersHandler()

		body := CreateUserRequest{
			Name:  TestNewUserName,
			Email: "invalid-email",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestUsersHandler_Update(t *testing.T) {
	t.Run("PUT /users/1 updates user", func(t *testing.T) {
		handler := setupUsersHandler()

		body := CreateUserRequest{
			Name:  TestUpdatedUserName,
			Email: TestUpdatedUserEmail,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var user models.User
		if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if user.Name != TestUpdatedUserName {
			t.Errorf("expected name '%s', got '%s'", TestUpdatedUserName, user.Name)
		}
	})

	t.Run("PUT /users/999 returns not found", func(t *testing.T) {
		handler := setupUsersHandler()

		body := CreateUserRequest{
			Name:  TestUpdatedUserName,
			Email: TestUpdatedUserEmail,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestUsersHandler_Delete(t *testing.T) {
	t.Run("DELETE /users/1 deletes user", func(t *testing.T) {
		handler := setupUsersHandler()

		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, rec.Code)
		}

		// 削除後に取得できないことを確認
		req = httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d after delete, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("DELETE /users/999 returns not found", func(t *testing.T) {
		handler := setupUsersHandler()

		req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestUsersHandler_MethodNotAllowed(t *testing.T) {
	handler := setupUsersHandler()

	t.Run("PATCH /users returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("PATCH /users/1 returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users/1", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}
