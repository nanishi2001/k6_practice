package config

import (
	"os"
	"time"
)

// Config はアプリケーション設定を保持する
type Config struct {
	Server     ServerConfig
	RateLimit  RateLimitConfig
	BodyLimit  int64
	CORS       CORSConfig
	CSRF       CSRFConfig
	JWT        JWTConfig
}

type ServerConfig struct {
	Addr string
}

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type CSRFConfig struct {
	AllowedOrigins []string
	StrictMode     bool
}

type JWTConfig struct {
	Secret []byte
}

// Load は設定を読み込む
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Addr: getEnv("SERVER_ADDR", ":8080"),
		},
		RateLimit: RateLimitConfig{
			Requests: 100,
			Window:   1 * time.Minute,
		},
		BodyLimit: 1 * 1024 * 1024, // 1MB
		CORS: CORSConfig{
			AllowedOrigins: []string{
				"http://localhost:8080",
				"http://localhost:3000",
				"http://127.0.0.1:8080",
				"http://127.0.0.1:3000",
			},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		},
		CSRF: CSRFConfig{
			AllowedOrigins: []string{
				"http://localhost:8080",
				"http://localhost:3000",
				"http://127.0.0.1:8080",
				"http://127.0.0.1:3000",
			},
			StrictMode: false, // 開発用
		},
		JWT: JWTConfig{
			Secret: getJWTSecret(),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("k6-test-secret-key-CHANGE-IN-PRODUCTION")
	}
	return []byte(secret)
}

// IsJWTSecretSet は JWT_SECRET 環境変数が設定されているか確認
func IsJWTSecretSet() bool {
	return os.Getenv("JWT_SECRET") != ""
}
