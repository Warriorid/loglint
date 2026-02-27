package rules_test

import (
	"testing"

	"github.com/loglint/internal/rules"
)

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name  string
		msg   string
		valid bool
	}{
		{"empty string", "", true},
		{"lowercase start", "starting server on port 8080", true},
		{"uppercase start", "Starting server on port 8080", false},
		{"uppercase error", "Failed to connect to database", false},
		{"lowercase error", "failed to connect to database", true},
		{"digit start", "123 servers", true},
		{"space start", " starting", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.CheckLowercase(tt.msg)
			if got != tt.valid {
				t.Errorf("CheckLowercase(%q) = %v, want %v", tt.msg, got, tt.valid)
			}
		})
	}
}

func TestCheckEnglishOnly(t *testing.T) {
	tests := []struct {
		name  string
		msg   string
		valid bool
	}{
		{"english only", "starting server", true},
		{"russian text", "запуск сервера", false},
		{"mixed", "server запуск", false},
		{"digits and symbols", "server 123 ok!", true},
		{"empty", "", true},
		{"arabic", "مرحبا", false},
		{"chinese", "你好", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.CheckEnglishOnly(tt.msg)
			if got != tt.valid {
				t.Errorf("CheckEnglishOnly(%q) = %v, want %v", tt.msg, got, tt.valid)
			}
		})
	}
}

func TestCheckNoSpecialChars(t *testing.T) {
	tests := []struct {
		name  string
		msg   string
		valid bool
	}{
		{"plain text", "server started", true},
		{"with colon", "server started: ok", true},
		{"with comma", "host, port", true},
		{"emoji rocket", "server started! 🚀", false},
		{"multiple exclamation", "connection failed!!!", false},
		{"ellipsis dots", "something went wrong...", false},
		{"with hyphen", "failed to connect", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.CheckNoSpecialChars(tt.msg)
			if got != tt.valid {
				t.Errorf("CheckNoSpecialChars(%q) = %v, want %v", tt.msg, got, tt.valid)
			}
		})
	}
}

func TestCheckNoSensitiveData(t *testing.T) {
	keywords := []string{
		"password", "passwd", "secret", "token", "api_key", "apikey",
		"credential", "private_key", "access_key", "bearer", "jwt",
	}
	tests := []struct {
		name  string
		msg   string
		valid bool
	}{
		{"clean message", "user authenticated successfully", true},
		{"contains password", "user password: secret123", false},
		{"contains token", "token validated", false},
		{"contains api_key", "api_key=mykey", false},
		{"contains secret", "my secret value", false},
		{"contains bearer", "bearer token set", false},
		{"contains jwt", "jwt decoded", false},
		{"api request", "api request completed", true},
		{"empty", "", true},
		{"case insensitive password", "user Password detected", false},
		{"case insensitive TOKEN", "TOKEN value", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.CheckNoSensitiveData(tt.msg, keywords)
			if got != tt.valid {
				t.Errorf("CheckNoSensitiveData(%q) = %v, want %v", tt.msg, got, tt.valid)
			}
		})
	}
}
