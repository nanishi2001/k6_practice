package middleware

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@mail.example.com", true},
		{"valid email with plus", "test+tag@example.com", true},
		{"invalid email no at", "testexample.com", false},
		{"invalid email no domain", "test@", false},
		{"invalid email no local", "@example.com", false},
		{"empty email", "", false},
		{"too long email", strings.Repeat("a", 250) + "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			if result != tt.expected {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid name", "John Doe", true},
		{"valid name japanese", "山田太郎", true},
		{"valid name with numbers", "User123", true},
		{"empty name", "", false},
		{"too long name", strings.Repeat("a", 101), false},
		{"name with control char", "Test\x00Name", false},
		{"name with tab", "Test\tName", true},
		{"name with newline", "Test\nName", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateName(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal string", "hello world", "hello world"},
		{"string with html", "<script>alert('xss')</script>", "alert('xss')"},
		{"string with spaces", "  hello  ", "hello"},
		{"string with control chars", "hello\x00world", "helloworld"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectValid bool
	}{
		{"valid password", "Password1!", true},
		{"valid complex password", "MyP@ssw0rd123", true},
		{"too short", "Pass1!", false},
		{"no uppercase", "password1!", false},
		{"no lowercase", "PASSWORD1!", false},
		{"no number", "Password!!", false},
		{"no special char", "Password123", false},
		{"empty password", "", false},
		{"too long", strings.Repeat("Aa1!", 50), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := ValidatePassword(tt.password)
			if valid != tt.expectValid {
				t.Errorf("ValidatePassword(%q) valid = %v, want %v", tt.password, valid, tt.expectValid)
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		expected bool
	}{
		{"positive id", 1, true},
		{"large id", 999999, true},
		{"zero id", 0, false},
		{"negative id", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateID(tt.id)
			if result != tt.expected {
				t.Errorf("ValidateID(%d) = %v, want %v", tt.id, result, tt.expected)
			}
		})
	}
}
