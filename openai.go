package ai

import (
	"context"
	"errors"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIAdapter struct {
	cfg    *Config
	client *openai.Client
}

func NewOpenAIAdapter(cfg *Config) Adapter {
	opts := []option.RequestOption{
		option.WithAPIKey(cfg.ResolveAPIKey()),
	}

	if cfg.Timeout > 0 {
		opts = append(opts, option.WithRequestTimeout(cfg.Timeout))
	}
	client := openai.NewClient(opts...)

	return &OpenAIAdapter{
		cfg:    cfg,
		client: &client,
	}
}

func (a *OpenAIAdapter) Generate(ctx context.Context, req *Request) (*Response, error) {
	params := openAIParams(a.cfg, req)

	chatCompletion, err := a.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}

	usage := &Usage{
		PromptTokens:     int(chatCompletion.Usage.PromptTokens),
		CompletionTokens: int(chatCompletion.Usage.CompletionTokens),
		TotalTokens:      int(chatCompletion.Usage.TotalTokens),
	}

	return &Response{
		Content: chatCompletion.Choices[0].Message.Content,
		Usage:   usage,
		Model:   chatCompletion.Model,
	}, nil
}

func (a *OpenAIAdapter) GenerateStream(ctx context.Context, req *Request) (<-chan StreamResponse, error) {
	params := openAIParams(a.cfg, req)

	stream := a.client.Chat.Completions.NewStreaming(ctx, params)

	streamChan := make(chan StreamResponse)

	go func() {
		defer close(streamChan)

		for stream.Next() {
			chunk := stream.Current()

			streamChan <- StreamResponse{
				Content: chunk.Choices[0].Delta.Content,
				Usage: &Usage{
					PromptTokens:     int(chunk.Usage.PromptTokens),
					CompletionTokens: int(chunk.Usage.CompletionTokens),
					TotalTokens:      int(chunk.Usage.TotalTokens),
				},
				Model: chunk.Model,
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

func (a *OpenAIAdapter) Validate() error {
	if a.cfg.Model != ModelGPT4o &&
		a.cfg.Model != ModelGPT4oMini &&
		a.cfg.Model != ModelGPT4Turbo {
		return errors.New("invalid model for openai")
	}
	return nil
}

func (a *OpenAIAdapter) GetProvider() Provider {
	return ProviderOpenAI
}

func (a *OpenAIAdapter) GetModel() string {
	return string(a.cfg.Model)
}

func openAIParams(cfg *Config, req *Request) openai.ChatCompletionNewParams {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(req.UserPrompt),
		},
		Model: openai.ChatModel(cfg.Model),
	}

	if req.SystemPrompt != "" {
		params.Messages = append(params.Messages, openai.SystemMessage(req.SystemPrompt))
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
		params.Temperature = openai.Float(*temp)
	}

	maxTokens := cfg.MaxTokens
	if req.MaxTokens != nil {
		maxTokens = req.MaxTokens
	}
	if maxTokens != nil {
		params.MaxTokens = openai.Int(*maxTokens)
	}

	return params
}
