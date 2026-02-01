package server

import (
	"log"
	"net/http"

	"k6-practice/api/config"
	"k6-practice/api/models"
	"k6-practice/api/router"
)

// Server はHTTPサーバーを表す
type Server struct {
	cfg       *config.Config
	handler   http.Handler
	userStore *models.UserStore
}

// New は新しいServerを作成
func New(cfg *config.Config) *Server {
	userStore := models.NewUserStore()

	r := router.New(cfg, userStore)
	handler := r.Build()

	return &Server{
		cfg:       cfg,
		handler:   handler,
		userStore: userStore,
	}
}

// Run はサーバーを起動
func (s *Server) Run() error {
	s.printStartupInfo()
	return http.ListenAndServe(s.cfg.Server.Addr, s.handler)
}

func (s *Server) printStartupInfo() {
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
	log.Printf("Listening on %s\n", s.cfg.Server.Addr)
}
