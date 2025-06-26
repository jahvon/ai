package ai_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jahvon/ai-adapter"
)

type providerTestCase struct {
	name        string
	provider    ai.Provider
	model       ai.Model
	apiKeyEnv   string
	skipMessage string
}

func TestProviders_E2E(t *testing.T) {
	testCases := map[string]providerTestCase{
		"OpenAI": {
			name:        "OpenAI",
			provider:    ai.ProviderOpenAI,
			model:       ai.ModelGPT4oMini,
			apiKeyEnv:   "OPENAI_API_KEY",
			skipMessage: "OPENAI_API_KEY environment variable not set",
		},
		"Anthropic": {
			name:        "Anthropic",
			provider:    ai.ProviderAnthropic,
			model:       ai.ModelSonnet3_5,
			apiKeyEnv:   "ANTHROPIC_API_KEY",
			skipMessage: "ANTHROPIC_API_KEY environment variable not set",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			apiKey := os.Getenv(tc.apiKeyEnv)
			if apiKey == "" {
				t.Skip(tc.skipMessage)
			}

			cfg := ai.NewConfig().
				WithProvider(tc.provider).
				WithModel(tc.model).
				WithAPIKeyEnv(tc.apiKeyEnv).
				WithMaxTokens(int64(500))

			adapter, err := ai.NewAdapter(cfg)
			if err != nil {
				t.Fatalf("Failed to create %s adapter: %v", tc.name, err)
			}

			req := &ai.Request{
				UserPrompt:   "Tell me a joke",
				SystemPrompt: "You are a pirate stranded at sea",
			}

			ctx := context.Background()
			resp, err := adapter.Generate(ctx, req)
			if err != nil {
				t.Fatalf("Failed to generate response: %v", err)
			}

			fmt.Printf("%s Response:\n%s\n\n", tc.name, resp.Content)
			fmt.Printf("Model:\n%s\n\n", resp.Model)
			fmt.Printf("Usage:\n%+v\n\n", resp.Usage)

			if resp.Content == "" {
				t.Error("Expected non-empty response content")
			}

			if resp.Usage == nil {
				t.Error("Expected usage information")
			}

			stream, err := adapter.GenerateStream(ctx, req)
			if err != nil {
				t.Fatalf("Failed to generate stream: %v", err)
			}

			accumulated, err := ai.CollectStream(ctx, stream)
			if err != nil {
				t.Fatalf("Failed to accumulate stream: %v", err)
			}

			fmt.Printf("Accumulated Response:\n%s\n\n", accumulated)
			if accumulated == "" {
				t.Error("Expected non-empty response content")
			}
		})
	}
}
