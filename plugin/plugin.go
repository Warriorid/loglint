package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/loglint/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type PluginSettings struct {
	CheckLowercase     bool     `json:"check-lowercase"`
	CheckEnglishOnly   bool     `json:"check-english-only"`
	CheckSpecialChars  bool     `json:"check-special-chars"`
	CheckSensitiveData bool     `json:"check-sensitive-data"`
	SensitiveKeywords  []string `json:"sensitive-keywords"`
}

type loglintPlugin struct {
	cfg *analyzer.Config
}

func (p *loglintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.New(p.cfg)}, nil
}

func (p *loglintPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func New(settings any) (register.LinterPlugin, error) {
	cfg := analyzer.DefaultConfig()

	s, err := register.DecodeSettings[PluginSettings](settings)
	if err == nil {
		cfg.CheckLowercase = s.CheckLowercase
		cfg.CheckEnglishOnly = s.CheckEnglishOnly
		cfg.CheckSpecialChars = s.CheckSpecialChars
		cfg.CheckSensitiveData = s.CheckSensitiveData
		if len(s.SensitiveKeywords) > 0 {
			cfg.SensitiveKeywords = s.SensitiveKeywords
		}
	}

	return &loglintPlugin{cfg: cfg}, nil
}
