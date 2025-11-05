package config

import "github.com/go-playground/validator/v10"

type CassandraCfg struct {
	Hosts         []string `yaml:"hosts" env:"CASSANDRA_HOSTS" validate:"required"`
	Port          int      `yaml:"port" env:"CASSANDRA_PORT" validate:"required"`
	Username      string   `yaml:"username" env:"CASSANDRA_USERNAME" validate:"required"`
	Password      string   `yaml:"password" env:"CASSANDRA_PASSWORD" validate:"required"`
	Keyspace      string   `yaml:"keyspace" env:"CASSANDRA_KEYSPACE" validate:"required"`
	LinkTableName string   `yaml:"link_table_name" env:"CASSANDRA_LINK_TABLE" validate:"required"`
}
type RedisCfg struct {
	Host         string `yaml:"host" env:"REDIS_HOST" validate:"required"`
	Port         int    `yaml:"port" env:"REDIS_PORT" validate:"required"`
	Password     string `yaml:"password" env:"REDIS_PASSWORD"`
	DB           int    `yaml:"db" env:"REDIS_DB"`
	RedisHashKey string `yaml:"redis_hash_key" env:"REDIS_HASH_KEY" `
}
type HashCfg struct {
	HashSalt      string `yaml:"hash_salt" env:"HASH_ALPHABET" validate:"required"`
	HashAlphabet  string `yaml:"hash_alphabet" env:"HASH_ALPHABET" validate:"required"`
	HashMinLength int    `yaml:"hash_min_length" env:"HASH_MIN_LENGTH" validate:"required"`
}
type Config struct {
	AppName  string `yaml:"app_name"   env:"APP_NAME"   validate:"required"`
	Env      string `yaml:"env"        env:"APP_ENV"    validate:"oneof=dev test staging prod"`
	LogLevel string `yaml:"log_level"  env:"LOG_LEVEL"  validate:"oneof=debug info warn error"`
	BaseURL  string `yaml:"base_url"   env:"BASE_URL"   validate:"required"`

	HashCfg      HashCfg      `yaml:"hash_config"`
	RedisCfg     RedisCfg     `yaml:"redis_config"`
	CassandraCfg CassandraCfg `yaml:"cassandra_config"`
}

func Default() *Config {
	return &Config{
		AppName:  "cutme",
		Env:      "dev",
		LogLevel: "info",
		BaseURL:  "localhost:8080",
		HashCfg: HashCfg{
			HashSalt:      "test-env",
			HashAlphabet:  "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
			HashMinLength: 7,
		},
		RedisCfg: RedisCfg{
			Host:         "localhost",
			Port:         6379,
			DB:           0,
			RedisHashKey: "cutme",
		},
		CassandraCfg: CassandraCfg{
			Hosts:         []string{"localhost"},
			Port:          9042,
			Username:      "cassandra",
			Password:      "cassandra",
			LinkTableName: "links",
		},
	}
}
func (c *Config) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}
