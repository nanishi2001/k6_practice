package middleware

import (
	"log"
	"net/http"
	"time"
)

// SecurityEvent はセキュリティイベントの種類
type SecurityEvent string

const (
	EventAuthSuccess     SecurityEvent = "AUTH_SUCCESS"
	EventAuthFailure     SecurityEvent = "AUTH_FAILURE"
	EventRateLimitHit    SecurityEvent = "RATE_LIMIT_HIT"
	EventCSRFBlocked     SecurityEvent = "CSRF_BLOCKED"
	EventInvalidInput    SecurityEvent = "INVALID_INPUT"
	EventUnauthorized    SecurityEvent = "UNAUTHORIZED"
	EventSuspiciousInput SecurityEvent = "SUSPICIOUS_INPUT"
)

// SecurityLogger はセキュリティイベントをログに記録
// A09:2021 - Security Logging and Monitoring Failures 対策
type SecurityLogger struct {
	enabled bool
}

func NewSecurityLogger(enabled bool) *SecurityLogger {
	return &SecurityLogger{enabled: enabled}
}

func (sl *SecurityLogger) Log(event SecurityEvent, r *http.Request, details string) {
	if !sl.enabled {
		return
	}

	ip := getClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	log.Printf("[SECURITY] event=%s ip=%s method=%s path=%s user_agent=%q details=%q timestamp=%s",
		event,
		ip,
		r.Method,
		r.URL.Path,
		userAgent,
		details,
		time.Now().Format(time.RFC3339),
	)
}

// Global security logger instance
var SecLogger = NewSecurityLogger(true)

// LogSecurityEvent はセキュリティイベントをログに記録するヘルパー
func LogSecurityEvent(event SecurityEvent, r *http.Request, details string) {
	SecLogger.Log(event, r, details)
}
