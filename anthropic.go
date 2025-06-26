package ai

import (
	"context"
	"errors"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicAdapter struct {
	cfg    *Config
	client *anthropic.Client
}

func NewAnthropicAdapter(cfg *Config) Adapter {
	opts := []option.RequestOption{
		option.WithAPIKey(cfg.ResolveAPIKey()),
	}

	if cfg.Timeout > 0 {
		opts = append(opts, option.WithRequestTimeout(cfg.Timeout))
	}
	client := anthropic.NewClient(opts...)

	return &AnthropicAdapter{
		cfg:    cfg,
		client: &client,
	}
}

func (a *AnthropicAdapter) Generate(ctx context.Context, req *Request) (*Response, error) {
	params := anthropicParams(a.cfg, req)

	response, err := a.client.Messages.New(ctx, params)
	if err != nil {
		return nil, err
	}

	usage := &Usage{
		PromptTokens:     int(response.Usage.InputTokens),
		CompletionTokens: int(response.Usage.OutputTokens),
		TotalTokens:      int(response.Usage.InputTokens + response.Usage.OutputTokens),
	}

	return &Response{
		Content: response.Content[0].Text,
		Usage:   usage,
		Model:   string(response.Model),
	}, nil
}

func (a *AnthropicAdapter) GenerateStream(ctx context.Context, req *Request) (<-chan StreamResponse, error) {
	params := anthropicParams(a.cfg, req)

	stream := a.client.Messages.NewStreaming(ctx, params)

	streamChan := make(chan StreamResponse)

	go func() {
		defer close(streamChan)

		for stream.Next() {
			chunk := stream.Current()
			switch chunk.Delta.Type {
			case "text":
				streamChan <- StreamResponse{
					Content: chunk.Delta.Text,
					Usage: &Usage{
						PromptTokens:     int(chunk.Usage.InputTokens),
						CompletionTokens: int(chunk.Usage.OutputTokens),
						TotalTokens:      int(chunk.Usage.InputTokens + chunk.Usage.OutputTokens),
					},
					Model: string(chunk.Message.Model),
				}
			}
		}

		if err := stream.Err(); err != nil {
			streamChan <- StreamResponse{
				Error: err,
			}
		}

		streamChan <- StreamResponse{
			Done: true,
		}
	}()

	return streamChan, nil
}

func (a *AnthropicAdapter) Validate() error {
	if a.cfg.Model != ModelSonnet3_5 &&
		a.cfg.Model != ModelSonnet3_7 &&
		a.cfg.Model != ModelSonnet4_0 {
		return errors.New("invalid model for anthropic")
	}
	return nil
}

func (a *AnthropicAdapter) GetProvider() Provider {
	return ProviderAnthropic
}

func (a *AnthropicAdapter) GetModel() string {
	return string(a.cfg.Model)
}

func anthropicParams(cfg *Config, req *Request) anthropic.MessageNewParams {
	params := anthropic.MessageNewParams{
		Model: anthropic.Model(cfg.Model),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.UserPrompt)),
		},
	}

	if req.SystemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: req.SystemPrompt},
		}
	}

	// TODO: figure out sending extra context
	// if req.Context != nil {
	// 	params.Context = req.Context
	// }

	temp := cfg.Temperature
	if req.Temperature != nil {
		temp = req.Temperature
	}
	if temp != nil {
		params.Temperature = anthropic.Float(*temp)
	}

	maxTokens := cfg.MaxTokens
	if req.MaxTokens != nil {
		maxTokens = req.MaxTokens
	}
	if maxTokens != nil {
		params.MaxTokens = *anthropic.IntPtr(*maxTokens)
	}

	return params
}
