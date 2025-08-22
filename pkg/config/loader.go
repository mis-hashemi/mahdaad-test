package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	defaultDelimiter    = "."
	defaultSeparator    = "__"
	defaultYamlFilePath = "./config.yml"
)

type Config any

// Loader struct is responsible for loading configuration from various sources.
type Loader[T Config] struct {
	koanf *koanf.Koanf
	opts  Option
}

// Option contains configuration options for the loader.
type Option struct {
	Prefix    string
	Delimiter string
	Separator string
}

// NewConfigLoader creates a new instance of Loader.
func NewConfigLoader[T Config](opts Option) *Loader[T] {
	if opts.Delimiter == "" {
		opts.Delimiter = defaultDelimiter
	}
	if opts.Separator == "" {
		opts.Separator = defaultSeparator
	}

	return &Loader[T]{
		koanf: koanf.New(opts.Delimiter),
		opts:  opts,
	}
}

type EnvCallbackFunc func(string) string

func (cl *Loader[T]) defaultCallbackEnv(source string) string {
	base := strings.ToLower(strings.TrimPrefix(source, cl.opts.Prefix))
	return strings.ReplaceAll(base, cl.opts.Separator, cl.opts.Delimiter)
}

// WithEnvProvider adds an environment variable provider to the loader.
func (cl *Loader[T]) WithEnvProvider() *Loader[T] {
	if err := cl.koanf.Load(env.Provider(cl.opts.Prefix, cl.opts.Delimiter, cl.defaultCallbackEnv), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}
	return cl
}

// WithYamlProvider adds a YAML file provider to the loader.
func (cl *Loader[T]) WithYamlProvider(filePath string) *Loader[T] {
	if filePath == "" {
		filePath = defaultYamlFilePath
	}
	if err := cl.koanf.Load(file.Provider(filePath), yaml.Parser()); err != nil {
		log.Printf("error loading yaml config: %s", err)
	}
	return cl
}

// WithDefaultProvider adds the default provider for structs to the loader.
func (cl *Loader[T]) WithDefaultProvider(defaultConfig T) *Loader[T] {
	if err := cl.koanf.Load(structs.Provider(defaultConfig, "koanf"), nil); err != nil {
		log.Printf("error loading default config: %s", err)
	}
	return cl
}

// Load loads the configuration into the provided Config struct.
func (cl *Loader[T]) Load(config *T) error {
	return cl.koanf.Unmarshal("", config)
}

// MustLoad loads config and panics on error (useful for startup).
func (cl *Loader[T]) MustLoad(config *T) {
	if err := cl.Load(config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
}
