package analyzer_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/loglint/internal/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func testdata() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata")
}

func TestAnalyzerDefault(t *testing.T) {
	analysistest.Run(t, testdata(), analyzer.Analyzer, "basic")
}

func TestAnalyzerAllRulesEnabled(t *testing.T) {
	cfg := analyzer.DefaultConfig()
	a := analyzer.New(cfg)
	analysistest.Run(t, testdata(), a, "basic")
}

func TestAnalyzerLowercaseOnly(t *testing.T) {
	cfg := &analyzer.Config{
		CheckLowercase:     true,
		CheckEnglishOnly:   false,
		CheckSpecialChars:  false,
		CheckSensitiveData: false,
	}
	a := analyzer.New(cfg)
	analysistest.Run(t, testdata(), a, "lowercase")
}

func TestAnalyzerEnglishOnly(t *testing.T) {
	cfg := &analyzer.Config{
		CheckLowercase:     false,
		CheckEnglishOnly:   true,
		CheckSpecialChars:  false,
		CheckSensitiveData: false,
	}
	a := analyzer.New(cfg)
	analysistest.Run(t, testdata(), a, "english")
}

func TestAnalyzerSpecialChars(t *testing.T) {
	cfg := &analyzer.Config{
		CheckLowercase:     false,
		CheckEnglishOnly:   false,
		CheckSpecialChars:  true,
		CheckSensitiveData: false,
	}
	a := analyzer.New(cfg)
	analysistest.Run(t, testdata(), a, "special")
}

func TestAnalyzerSensitiveData(t *testing.T) {
	cfg := &analyzer.Config{
		CheckLowercase:     false,
		CheckEnglishOnly:   false,
		CheckSpecialChars:  false,
		CheckSensitiveData: true,
		SensitiveKeywords:  analyzer.DefaultConfig().SensitiveKeywords,
	}
	a := analyzer.New(cfg)
	analysistest.Run(t, testdata(), a, "sensitive")
}
