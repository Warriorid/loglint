package analyzer

type Config struct {
	CheckLowercase     bool     `yaml:"check-lowercase" mapstructure:"check-lowercase"`
	CheckEnglishOnly   bool     `yaml:"check-english-only" mapstructure:"check-english-only"`
	CheckSpecialChars  bool     `yaml:"check-special-chars" mapstructure:"check-special-chars"`
	CheckSensitiveData bool     `yaml:"check-sensitive-data" mapstructure:"check-sensitive-data"`
	SensitiveKeywords  []string `yaml:"sensitive-keywords" mapstructure:"sensitive-keywords"`
}

func DefaultConfig() *Config {
	return &Config{
		CheckLowercase:     true,
		CheckEnglishOnly:   true,
		CheckSpecialChars:  true,
		CheckSensitiveData: true,
		SensitiveKeywords:  defaultSensitiveKeywords,
	}
}

var defaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"api-key",
	"credential",
	"private_key",
	"privatekey",
	"access_key",
	"accesskey",
	"bearer",
	"jwt",
	"ssn",
	"credit_card",
	"creditcard",
}
