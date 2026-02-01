package middleware

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validator は入力検証ユーティリティ
// A03:2021 - Injection 対策

// ValidateEmail はメールアドレスの形式を検証
func ValidateEmail(email string) bool {
	if len(email) > 254 {
		return false
	}
	// RFC 5322に基づく簡易的なメールアドレス検証
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateName はユーザー名を検証
func ValidateName(name string) bool {
	// 長さチェック
	if len(name) < 1 || len(name) > 100 {
		return false
	}
	// UTF-8として有効かチェック
	if !utf8.ValidString(name) {
		return false
	}
	// 制御文字を含まないかチェック
	for _, r := range name {
		if r < 32 && r != '\t' && r != '\n' {
			return false
		}
	}
	return true
}

// SanitizeString は文字列から危険な文字を除去
func SanitizeString(s string) string {
	// HTMLタグを除去
	s = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(s, "")
	// 制御文字を除去（タブ、改行は許可）
	s = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`).ReplaceAllString(s, "")
	// 前後の空白を除去
	s = strings.TrimSpace(s)
	return s
}

// ValidatePassword はパスワードの強度を検証
// A07:2021 - Identification and Authentication Failures 対策
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "password must be at least 8 characters"
	}
	if len(password) > 128 {
		return false, "password must be at most 128 characters"
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, c := range password {
		switch {
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case 'a' <= c && c <= 'z':
			hasLower = true
		case '0' <= c && c <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "password must contain at least one number"
	}
	if !hasSpecial {
		return false, "password must contain at least one special character"
	}

	return true, ""
}

// ValidateID はIDが正の整数かを検証
func ValidateID(id int) bool {
	return id > 0
}
