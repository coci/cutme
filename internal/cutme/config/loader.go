package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
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
	o := &AppOptions{EnvPrefix: "APP"}
	for _, fn := range opts {
		fn(o)
	}

	// flags (allow CLI to set/override file path & print)
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

	cfg := Default()

	// 1) YAML file
	if o.FilePath != "" {
		if err := applyYAML(o.FilePath, cfg); err != nil {
			return nil, err
		}
	} else if o.RequireFile {
		return nil, fmt.Errorf("config file required but not provided")
	}

	// 2) ENV overrides
	// Convert nested fields: we already put explicit env tags on fields.
	// Prefix applied by env package via Options.
	err := env.Parse(cfg)

	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
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
