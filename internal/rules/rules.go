package rules

import (
	"strings"
	"unicode"
)

func CheckLowercase(msg string) bool {
	for _, r := range msg {
		if unicode.IsLetter(r) {
			return unicode.IsLower(r)
		}
	}
	return true
}

func CheckEnglishOnly(msg string) bool {
	for _, r := range msg {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func CheckNoSpecialChars(msg string) bool {
	runes := []rune(msg)
	for i, r := range runes {
		if isEmoji(r) {
			return false
		}
		if isSpecialChar(r) {
			return false
		}
		if i > 0 && isPunct(r) && isPunct(runes[i-1]) {
			return false
		}
	}
	return true
}

func isPunct(r rune) bool {
	return r == '.' || r == '!' || r == '?' || r == ','
}

func isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x1F300 && r <= 0x1F5FF) ||
		(r >= 0x1F680 && r <= 0x1F6FF) ||
		(r >= 0x1F1E0 && r <= 0x1F1FF) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF) ||
		(r >= 0xFE00 && r <= 0xFE0F) ||
		(r >= 0x1F900 && r <= 0x1F9FF) ||
		(r >= 0x1FA00 && r <= 0x1FA6F) ||
		(r >= 0x1FA70 && r <= 0x1FAFF) ||
		(r >= 0x231A && r <= 0x231B) ||
		(r >= 0x23E9 && r <= 0x23F3) ||
		(r >= 0x25AA && r <= 0x25AB) ||
		(r >= 0x25B6 && r <= 0x25C0) ||
		r == 0x200D
}

func isSpecialChar(r rune) bool {
	if r > unicode.MaxASCII {
		return false
	}
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 -_/.:,()[]{}'\"`#@$%^&*=+<>|~\\"
	return !strings.ContainsRune(allowed, r) && r != '\t' && r != '\n'
}

func CheckNoSensitiveData(msg string, keywords []string) bool {
	lower := strings.ToLower(msg)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return false
		}
	}
	return true
}
