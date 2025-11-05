package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type AppOptions struct {
	FilePath    string
	EnvPrefix   string
	Args        []string
	RequireFile bool
}

type Option func(*AppOptions)

func WithFile(path string) Option   { return func(o *AppOptions) { o.FilePath = path } }
func WithEnvPrefix(p string) Option { return func(o *AppOptions) { o.EnvPrefix = p } }
func WithArgs(a []string) Option    { return func(o *AppOptions) { o.Args = a } }

func Load(opts ...Option) (*Config, error) {
	o := &AppOptions{EnvPrefix: ""} // Remove the default prefix
	for _, fn := range opts {
		fn(o)
	}

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	flagConfig := fs.String("config", "", "path to config YAML")

	args := o.Args
	if args == nil {
		args = os.Args[1:]
	}
	_ = fs.Parse(args)

	if *flagConfig != "" {
		o.FilePath = *flagConfig
	}

	// 0) Auto-load .env (ignore if missing)
	envErr := godotenv.Load()
	if envErr != nil {
		_ = godotenv.Load(".env")
		_ = godotenv.Load("../.env")
	}

	cfg := Default()

	// 1) YAML file (this will override defaults)
	if o.FilePath != "" {
		if err := applyYAML(o.FilePath, cfg); err != nil {
			return nil, err
		}
	} else if o.RequireFile {
		return nil, fmt.Errorf("config file required but not provided")
	}

	// 2) Environment variables (this will override YAML)
	if o.EnvPrefix != "" {
		prefix := strings.TrimSuffix(o.EnvPrefix, "_")
		if prefix != "" {
			prefix += "_"
		}
		if err := env.ParseWithOptions(cfg, env.Options{Prefix: prefix}); err != nil {
			return nil, fmt.Errorf("parse env: %w", err)
		}
	} else {
		// No prefix
		if err := env.Parse(cfg); err != nil {
			return nil, fmt.Errorf("parse env: %w", err)
		}
	}

	// 3) Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func applyYAML(path string, out *Config) error {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("config file not found: %s", path)
		}
		return fmt.Errorf("read config file: %w", err)
	}
	dec := yaml.NewDecoder(strings.NewReader(string(b)))
	dec.KnownFields(true) // strict: error on unknown fields
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}
	return nil
}
