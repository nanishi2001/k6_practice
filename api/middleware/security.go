package middleware

import (
	"net/http"
)

// SecurityHeaders はOWASP推奨のセキュリティヘッダーを追加するミドルウェア
// A05:2021 - Security Misconfiguration 対策
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// XSS対策: ブラウザのXSSフィルターを有効化
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// クリックジャッキング対策: iframeでの埋め込みを禁止
		w.Header().Set("X-Frame-Options", "DENY")

		// MIMEタイプスニッフィング対策: Content-Typeを厳密に解釈
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Referrer情報の漏洩防止
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// コンテンツセキュリティポリシー: インラインスクリプト等を制限
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")

		// キャッシュ制御: 機密データのキャッシュを防止
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		w.Header().Set("Pragma", "no-cache")

		// HTTPS強制（本番環境用）
		// w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// サーバー情報の隠蔽
		w.Header().Set("X-Powered-By", "")
		w.Header().Set("Server", "")

		// Permissions-Policy: ブラウザ機能の制限
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		next.ServeHTTP(w, r)
	})
}
