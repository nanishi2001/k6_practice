package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"k6-practice/api/middleware"
	"k6-practice/api/models"
)

type AuthHandler struct {
	store *models.UserStore
}

func NewAuthHandler(store *models.UserStore) *AuthHandler {
	return &AuthHandler{store: store}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type MeResponse struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

// トークン有効期限
const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 24 * time.Hour
	accessTokenExpiresIn = 900 // 15分（秒）
)

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/auth/")

	switch path {
	case "login":
		h.login(w, r)
	case "refresh":
		h.refresh(w, r)
	case "me":
		h.me(w, r)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// テスト用: パスワードは "password" で固定
	if req.Password != "password" {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	user := h.findUserByEmail(req.Email)
	if user == nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	tokens, err := h.generateTokenPair(user.ID, user.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return middleware.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		writeError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	tokens, err := h.generateTokenPair(claims.UserID, claims.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	json.NewEncoder(w).Encode(MeResponse{
		UserID: claims.UserID,
		Email:  claims.Email,
	})
}

// findUserByEmail はメールアドレスでユーザーを検索する
func (h *AuthHandler) findUserByEmail(email string) *models.User {
	for _, u := range h.store.List() {
		if u.Email == email {
			return u
		}
	}
	return nil
}

// generateTokenPair はアクセストークンとリフレッシュトークンを生成する
func (h *AuthHandler) generateTokenPair(userID int, email string) (*TokenResponse, error) {
	accessToken, err := generateToken(userID, email, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(userID, email, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessTokenExpiresIn,
	}, nil
}

func generateToken(userID int, email string, duration time.Duration) (string, error) {
	claims := &middleware.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JWTSecret)
}
