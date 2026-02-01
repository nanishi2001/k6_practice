package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"k6-practice/api/handlers"
	"k6-practice/api/middleware"
	"k6-practice/api/models"
)

// 許可するOrigin（本番環境では適切に設定）
var allowedOrigins = []string{
	"http://localhost:8080",
	"http://localhost:3000",
	"http://127.0.0.1:8080",
	"http://127.0.0.1:3000",
}

// 設定値
const (
	// レート制限: 1分間に100リクエストまで
	rateLimitRequests = 100
	rateLimitWindow   = 1 * time.Minute

	// リクエストボディの最大サイズ: 1MB
	maxBodySize = 1 * 1024 * 1024
)

func main() {
	// 環境変数のチェック
	if os.Getenv("JWT_SECRET") == "" {
		log.Println("[WARNING] JWT_SECRET is not set. Using default secret. Set JWT_SECRET in production!")
	}

	// インメモリデータストアを初期化
	userStore := models.NewUserStore()

	// ハンドラーを初期化
	usersHandler := handlers.NewUsersHandler(userStore)
	authHandler := handlers.NewAuthHandler(userStore)

	// レート制限ミドルウェア
	rateLimiter := middleware.NewRateLimiter(rateLimitRequests, rateLimitWindow)

	// CSRF保護ミドルウェア
	// StrictMode: false = 開発・テスト用（curl, k6からのリクエストを許可）
	// 本番環境では StrictMode: true に設定
	csrfProtect := middleware.CSRFProtection(middleware.CSRFConfig{
		AllowedOrigins: allowedOrigins,
		StrictMode:     false,
	})

	// ボディサイズ制限ミドルウェア
	bodyLimit := middleware.BodyLimit(maxBodySize)

	// ルーティング設定
	mux := http.NewServeMux()

	// 基本エンドポイント（CSRF保護不要）
	mux.HandleFunc("/health", handlers.HealthCheck)

	// ユーザーCRUD（CSRF保護 + ボディサイズ制限あり）
	mux.Handle("/users", bodyLimit(csrfProtect(usersHandler)))
	mux.Handle("/users/", bodyLimit(csrfProtect(usersHandler)))

	// 認証エンドポイント（CSRF保護 + ボディサイズ制限あり）
	mux.Handle("/auth/login", bodyLimit(csrfProtect(authHandler)))
	mux.Handle("/auth/refresh", bodyLimit(csrfProtect(authHandler)))
	mux.Handle("/auth/me", middleware.Auth(csrfProtect(authHandler)))

	// 遅延シミュレーション（CSRF保護不要、GETのみ）
	mux.HandleFunc("/delay/", handlers.DelayHandler)
	mux.HandleFunc("/random-delay", handlers.RandomDelayHandler)
	mux.HandleFunc("/error-rate/", handlers.ErrorRateHandler)

	// ミドルウェアチェーンを構築（外側から順に適用）
	// リクエスト: Logging → SecurityHeaders → RateLimit → CORS → ルーター
	handler := middleware.Logging(
		middleware.SecurityHeaders(
			rateLimiter.Middleware(
				corsMiddleware(mux),
			),
		),
	)

	log.Println("=== API Server Starting ===")
	log.Println("")
	log.Println("Security Features (OWASP):")
	log.Println("  - A02: JWT with configurable secret (set JWT_SECRET env var)")
	log.Println("  - A03: Input validation and sanitization")
	log.Println("  - A04: Rate limiting (100 req/min), Body size limit (1MB)")
	log.Println("  - A05: Security headers (XSS, Clickjacking, MIME sniffing)")
	log.Println("  - A07: CSRF protection")
	log.Println("  - A09: Security event logging")
	log.Println("")
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
	log.Println("")
	log.Println("Listening on :8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

