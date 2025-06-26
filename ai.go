package ai

import (
	"context"
	"errors"
)

// Request represents a standardized AI request
type Request struct {
	SystemPrompt string `yaml:"systemPrompt,omitempty" json:"systemPrompt,omitempty"`
	UserPrompt   string `yaml:"userPrompt" json:"userPrompt"`
	Stream       bool   `yaml:"stream,omitempty" json:"stream,omitempty"`

	// Context and additional data
	Context map[string]interface{} `yaml:"context,omitempty" json:"context,omitempty"`

	// Request-specific overrides
	Temperature *float64 `yaml:"temperature,omitempty" json:"temperature,omitempty"`
	MaxTokens   *int64   `yaml:"maxTokens,omitempty" json:"maxTokens,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}

// Response represents a standardized AI response
type Response struct {
	Content string `json:"content"`
	Usage   *Usage `json:"usage,omitempty"`
	Model   string `json:"model"`
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	Content string `json:"content"`
	Usage   *Usage `json:"usage,omitempty"`
	Model   string `json:"model"`
	Done    bool   `json:"done"`
	Error   error  `json:"error,omitempty"`
}

// Adapter is the main interface for AI interactions
type Adapter interface {
	// Generate sends a request and returns a complete response
	Generate(ctx context.Context, req *Request) (*Response, error)

	// GenerateStream sends a request and returns a stream of response chunks
	GenerateStream(ctx context.Context, req *Request) (<-chan StreamResponse, error)

	// Validate checks if the adapter configuration is valid
	Validate() error

	// GetProvider returns the provider type
	GetProvider() Provider

	// GetModel returns the configured model
	GetModel() string
}

func NewAdapter(cfg *Config) (Adapter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var adapter Adapter
	switch cfg.Provider {
	case ProviderOpenAI:
		adapter = NewOpenAIAdapter(cfg)
	case ProviderAnthropic:
		adapter = NewAnthropicAdapter(cfg)
	default:
		return nil, errors.New("invalid provider")
	}

	if err := adapter.Validate(); err != nil {
		return nil, err
	}

	return adapter, nil
}
