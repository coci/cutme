package config

import "github.com/go-playground/validator/v10"

type HashCfg struct {
	HashSalt      string `yaml:"hash_salt" env:"HASH_ALPHABET" validate:"required"`
	HashAlphabet  string `yaml:"hash_alphabet" env:"HASH_ALPHABET" validate:"required"`
	HashMinLength int    `yaml:"hash_min_length" env:"HASH_MIN_LENGTH" validate:"required"`
}
type Config struct {
	AppName  string `yaml:"app_name"   env:"APP_NAME"   validate:"required"`
	Env      string `yaml:"env"        env:"APP_ENV"    validate:"oneof=dev test staging prod"`
	LogLevel string `yaml:"log_level"  env:"LOG_LEVEL"  validate:"oneof=debug info warn error"`
	BaseURL  string `yaml:"base_url"   env:"BASE_URL"   validate:"required,url"`

	HashCfg HashCfg `yaml:"hash_config"`
}

func Default() *Config {
	return &Config{
		AppName:  "cutme",
		Env:      "dev",
		LogLevel: "info",
		BaseURL:  "http://localhost:8080",
		HashCfg: HashCfg{
			HashSalt:      "test-env",
			HashAlphabet:  "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
			HashMinLength: 7,
		},
	}
}
func (c *Config) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}
