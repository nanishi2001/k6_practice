package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"k6-practice/api/middleware"
	"k6-practice/api/models"
)

type UsersHandler struct {
	store *models.UserStore
}

func NewUsersHandler(store *models.UserStore) *UsersHandler {
	return &UsersHandler{store: store}
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/users")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.list(w, r)
		case http.MethodPost:
			h.create(w, r)
		default:
			methodNotAllowed(w)
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.get(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		methodNotAllowed(w)
	}
}

func (h *UsersHandler) list(w http.ResponseWriter, r *http.Request) {
	users := h.store.List()
	json.NewEncoder(w).Encode(users)
}

func (h *UsersHandler) get(w http.ResponseWriter, r *http.Request, id int) {
	user := h.store.Get(id)
	if user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UsersHandler) create(w http.ResponseWriter, r *http.Request) {
	req, ok := h.parseAndValidateUserRequest(w, r)
	if !ok {
		return
	}

	user := h.store.Create(req.Name, req.Email)
	writeJSON(w, http.StatusCreated, user)
}

func (h *UsersHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	if !middleware.ValidateID(id) {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	req, ok := h.parseAndValidateUserRequest(w, r)
	if !ok {
		return
	}

	user := h.store.Update(id, req.Name, req.Email)
	if user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UsersHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if !h.store.Delete(id) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseAndValidateUserRequest はリクエストボディをパースしてバリデーションを行う
func (h *UsersHandler) parseAndValidateUserRequest(w http.ResponseWriter, r *http.Request) (*CreateUserRequest, bool) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.LogSecurityEvent(middleware.EventInvalidInput, r, "invalid JSON body")
		writeError(w, http.StatusBadRequest, "invalid request body")
		return nil, false
	}

	// 入力のサニタイズ
	req.Name = middleware.SanitizeString(req.Name)
	req.Email = middleware.SanitizeString(req.Email)

	// 必須チェック
	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email are required")
		return nil, false
	}

	// Name検証
	if !middleware.ValidateName(req.Name) {
		middleware.LogSecurityEvent(middleware.EventInvalidInput, r, "invalid name format")
		writeError(w, http.StatusBadRequest, "invalid name format")
		return nil, false
	}

	// Email検証
	if !middleware.ValidateEmail(req.Email) {
		middleware.LogSecurityEvent(middleware.EventInvalidInput, r, "invalid email format")
		writeError(w, http.StatusBadRequest, "invalid email format")
		return nil, false
	}

	return &req, true
}
