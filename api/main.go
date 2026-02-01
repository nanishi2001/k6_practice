package main

import (
	"log"
	"net/http"
	"strings"

	"k6-practice/api/handlers"
	"k6-practice/api/middleware"
	"k6-practice/api/models"
)

func main() {
	// インメモリデータストアを初期化
	userStore := models.NewUserStore()

	// ハンドラーを初期化
	usersHandler := handlers.NewUsersHandler(userStore)
	authHandler := handlers.NewAuthHandler(userStore)

	// ルーティング設定
	mux := http.NewServeMux()

	// 基本エンドポイント
	mux.HandleFunc("/health", handlers.HealthCheck)

	// ユーザーCRUD
	mux.Handle("/users", usersHandler)
	mux.Handle("/users/", usersHandler)

	// 認証エンドポイント
	mux.Handle("/auth/login", authHandler)
	mux.Handle("/auth/refresh", authHandler)
	mux.Handle("/auth/me", middleware.Auth(authHandler))

	// 遅延シミュレーション
	mux.HandleFunc("/delay/", handlers.DelayHandler)
	mux.HandleFunc("/random-delay", handlers.RandomDelayHandler)
	mux.HandleFunc("/error-rate/", handlers.ErrorRateHandler)

	// ミドルウェアを適用
	handler := middleware.Logging(corsMiddleware(mux))

	log.Println("Starting API server on :8080")
	log.Println("Endpoints:")
	log.Println("  GET  /health           - Health check")
	log.Println("  GET  /users            - List users")
	log.Println("  GET  /users/:id        - Get user")
	log.Println("  POST /users            - Create user")
	log.Println("  PUT  /users/:id        - Update user")
	log.Println("  DELETE /users/:id      - Delete user")
	log.Println("  POST /auth/login       - Login (get JWT)")
	log.Println("  POST /auth/refresh     - Refresh token")
	log.Println("  GET  /auth/me          - Get current user (requires auth)")
	log.Println("  GET  /delay/:ms        - Delay response by ms")
	log.Println("  GET  /random-delay     - Random delay (0-1000ms)")
	log.Println("  GET  /error-rate/:pct  - Return error at given rate")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func init() {
	// パスのトレイリングスラッシュを正規化するためのヘルパー
	_ = strings.TrimSuffix
}
