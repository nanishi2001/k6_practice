package middleware

import (
	"net/http"
)

// BodyLimit はリクエストボディのサイズを制限するミドルウェア
// A04:2021 - Insecure Design 対策（DoS防止）
func BodyLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				w.Write([]byte(`{"error": "request body too large"}`))
				return
			}

			// ボディサイズを制限
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

			next.ServeHTTP(w, r)
		})
	}
}
