package middleware

import (
	"net/http"
	"strings"
)

type CSRFConfig struct {
	AllowedOrigins []string
	// StrictMode: trueの場合、Origin/Referer/X-Requested-Withを厳密に検証
	// falseの場合、ヘッダーがない場合は許可（開発・テスト用）
	StrictMode bool
}

// CSRFProtection はCSRF攻撃を防ぐミドルウェア
// 状態変更リクエスト（POST/PUT/DELETE）に対して以下を検証:
// 1. Origin/Refererヘッダーが許可されたホストか
// 2. X-Requested-Withヘッダーが存在するか（カスタムヘッダー要求）
func CSRFProtection(config CSRFConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 安全なメソッド（GET, HEAD, OPTIONS）はスキップ
			if isSafeMethod(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			// 1. Origin または Referer ヘッダーの検証
			origin := r.Header.Get("Origin")
			referer := r.Header.Get("Referer")

			if origin == "" && referer == "" {
				if config.StrictMode {
					http.Error(w, `{"error": "CSRF validation failed: missing origin/referer"}`, http.StatusForbidden)
					return
				}
				// 非StrictMode: ヘッダーがない場合は許可（curl, k6等のテストツール用）
				next.ServeHTTP(w, r)
				return
			}

			if origin != "" && !isAllowedOrigin(origin, config.AllowedOrigins) {
				http.Error(w, `{"error": "CSRF validation failed: invalid origin"}`, http.StatusForbidden)
				return
			}

			if origin == "" && referer != "" && !isAllowedReferer(referer, config.AllowedOrigins) {
				http.Error(w, `{"error": "CSRF validation failed: invalid referer"}`, http.StatusForbidden)
				return
			}

			// 2. X-Requested-With ヘッダーの検証（XMLHttpRequest等からのリクエストを確認）
			// ブラウザからのリクエストの場合のみチェック
			if origin != "" || referer != "" {
				xRequestedWith := r.Header.Get("X-Requested-With")
				if xRequestedWith == "" {
					http.Error(w, `{"error": "CSRF validation failed: missing X-Requested-With header"}`, http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSRFProtectionSimple は簡易版（後方互換性のため）
func CSRFProtectionSimple(allowedOrigins []string) func(http.Handler) http.Handler {
	return CSRFProtection(CSRFConfig{
		AllowedOrigins: allowedOrigins,
		StrictMode:     false,
	})
}

func isSafeMethod(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodOptions
}

func isAllowedOrigin(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}
	return false
}

func isAllowedReferer(referer string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" {
			return true
		}
		// Refererはフルパスなので、プレフィックスマッチ
		if strings.HasPrefix(referer, a) {
			return true
		}
	}
	return false
}
