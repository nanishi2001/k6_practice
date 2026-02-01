package middleware

import (
	"net/http"
	"strings"
)

// CORSConfig はCORSミドルウェアの設定
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// CORS はCross-Origin Resource Sharingを処理するミドルウェア
func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	methods := strings.Join(cfg.AllowedMethods, ", ")
	headers := strings.Join(cfg.AllowedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Originがある場合、許可リストをチェック
			if origin != "" {
				allowed := false
				for _, o := range cfg.AllowedOrigins {
					if o == "*" || o == origin {
						allowed = true
						break
					}
				}
				if allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			} else {
				// Originがない場合（同一オリジンまたはcurl等）は*を設定
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)

			// プリフライトリクエスト
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSSimple は簡易版CORS（後方互換性用）
func CORSSimple(next http.Handler) http.Handler {
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
