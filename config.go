package ai

import (
	"errors"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/openai/openai-go"
)

// Provider represents different AI service providers
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
)

type Model string

const (
	ModelGPT4o     Model = Model(openai.ChatModelGPT4o)
	ModelGPT4oMini Model = Model(openai.ChatModelGPT4oMini)
	ModelGPT4Turbo Model = Model(openai.ChatModelGPT4Turbo)
	ModelSonnet3_5 Model = Model(anthropic.ModelClaude3_5SonnetLatest)
	ModelSonnet3_7 Model = Model(anthropic.ModelClaude3_7SonnetLatest)
	ModelSonnet4_0 Model = Model(anthropic.ModelClaudeSonnet4_0)
)

type Config struct {
	Provider    Provider      `yaml:"provider" json:"provider"`
	Model       Model         `yaml:"model" json:"model"`
	APIKey      string        `yaml:"apiKey,omitempty" json:"apiKey,omitempty"`
	APIKeyEnv   string        `yaml:"apiKeyEnv,omitempty" json:"apiKeyEnv,omitempty"`
	Temperature *float64      `yaml:"temperature,omitempty" json:"temperature,omitempty"`
	MaxTokens   *int64        `yaml:"maxTokens,omitempty" json:"maxTokens,omitempty"`
	Timeout     time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	if c.APIKey == "" && c.APIKeyEnv == "" {
		return errors.New("apiKey or apiKeyEnv is required")
	}

	if c.Provider == "" {
		return errors.New("provider is required")
	}

	if c.Model == "" {
		return errors.New("model is required")
	}

	return nil
}

func (c *Config) WithProvider(p Provider) *Config {
	c.Provider = p
	return c
}

func (c *Config) WithModel(m Model) *Config {
	c.Model = m
	return c
}

func (c *Config) WithAPIKey(k string) *Config {
	c.APIKey = k
	return c
}

func (c *Config) WithAPIKeyEnv(e string) *Config {
	c.APIKeyEnv = e
	return c
}

func (c *Config) WithTemperature(t float64) *Config {
	c.Temperature = &t
	return c
}

func (c *Config) WithMaxTokens(t int64) *Config {
	c.MaxTokens = &t
	return c
}

func (c *Config) WithTimeout(t time.Duration) *Config {
	c.Timeout = t
	return c
}

func (c *Config) ResolveAPIKey() string {
	if c.APIKey != "" {
		return c.APIKey
	}
	return os.Getenv(c.APIKeyEnv)
}
