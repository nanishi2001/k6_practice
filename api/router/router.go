package router

import (
	"net/http"

	"k6-practice/api/config"
	"k6-practice/api/handlers"
	"k6-practice/api/middleware"
	"k6-practice/api/models"
)

// Router はアプリケーションのルーターを構築
type Router struct {
	cfg       *config.Config
	userStore *models.UserStore
}

// New は新しいRouterを作成
func New(cfg *config.Config, userStore *models.UserStore) *Router {
	return &Router{
		cfg:       cfg,
		userStore: userStore,
	}
}

// Build はHTTPハンドラーを構築
func (r *Router) Build() http.Handler {
	mux := http.NewServeMux()

	// ハンドラー初期化
	usersHandler := handlers.NewUsersHandler(r.userStore)
	authHandler := handlers.NewAuthHandler(r.userStore)

	// ミドルウェア
	rateLimiter := middleware.NewRateLimiter(
		r.cfg.RateLimit.Requests,
		r.cfg.RateLimit.Window,
	)

	csrfProtect := middleware.CSRFProtection(middleware.CSRFConfig{
		AllowedOrigins: r.cfg.CSRF.AllowedOrigins,
		StrictMode:     r.cfg.CSRF.StrictMode,
	})

	bodyLimit := middleware.BodyLimit(r.cfg.BodyLimit)

	cors := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins: r.cfg.CORS.AllowedOrigins,
		AllowedMethods: r.cfg.CORS.AllowedMethods,
		AllowedHeaders: r.cfg.CORS.AllowedHeaders,
	})

	// ミドルウェアチェーン定義
	// 保護付きエンドポイント用
	protected := middleware.NewChain(bodyLimit, csrfProtect)
	// 認証必須エンドポイント用
	authenticated := middleware.NewChain(middleware.Auth, csrfProtect)

	// ルート登録
	// ヘルスチェック（ミドルウェアなし）
	mux.HandleFunc("/health", handlers.HealthCheck)

	// ユーザーCRUD
	mux.Handle("/users", protected.Then(usersHandler))
	mux.Handle("/users/", protected.Then(usersHandler))

	// 認証エンドポイント
	mux.Handle("/auth/login", protected.Then(authHandler))
	mux.Handle("/auth/refresh", protected.Then(authHandler))
	mux.Handle("/auth/me", authenticated.Then(authHandler))

	// 遅延・エラーシミュレーション（CSRF保護不要）
	mux.HandleFunc("/delay/", handlers.DelayHandler)
	mux.HandleFunc("/random-delay", handlers.RandomDelayHandler)
	mux.HandleFunc("/error-rate/", handlers.ErrorRateHandler)

	// グローバルミドルウェアチェーン
	global := middleware.NewChain(
		middleware.Logging,
		middleware.SecurityHeaders,
		rateLimiter.Middleware,
		cors,
	)

	return global.Then(mux)
}
