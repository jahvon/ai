package ai_test

import (
	"context"
	"log"
	"os"
	"testing"

	ai "github.com/jahvon/ai"
)

func TestOpenAIAdapter_E2E(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY environment variable not set")
	}

	cfg := ai.NewConfig().
		WithProvider(ai.ProviderOpenAI).
		WithModel(ai.ModelGPT4oMini).
		WithAPIKeyEnv("OPENAI_API_KEY")

	adapter, err := ai.NewAdapter(cfg)
	if err != nil {
		t.Fatalf("Failed to create OpenAI adapter: %v", err)
	}

	req := &ai.Request{
		UserPrompt: "Say hello in exactly 3 words.",
	}

	ctx := context.Background()
	resp, err := adapter.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Failed to generate response: %v", err)
	}

	log.Printf("OpenAI Response: %s", resp.Content)
	log.Printf("Model: %s", resp.Model)
	log.Printf("Usage: %+v", resp.Usage)

	if resp.Content == "" {
		t.Error("Expected non-empty response content")
	}

	if resp.Usage == nil {
		t.Error("Expected usage information")
	}
}

func TestAnthropicAdapter_E2E(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable not set")
	}

	cfg := ai.NewConfig().
		WithProvider(ai.ProviderAnthropic).
		WithModel(ai.ModelSonnet3_5).
		WithAPIKeyEnv("ANTHROPIC_API_KEY")

	adapter, err := ai.NewAdapter(cfg)
	if err != nil {
		t.Fatalf("Failed to create Anthropic adapter: %v", err)
	}

	req := &ai.Request{
		UserPrompt: "Say hello in exactly 3 words.",
	}

	ctx := context.Background()
	resp, err := adapter.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Failed to generate response: %v", err)
	}

	log.Printf("Anthropic Response: %s", resp.Content)
	log.Printf("Model: %s", resp.Model)
	log.Printf("Usage: %+v", resp.Usage)

	if resp.Content == "" {
		t.Error("Expected non-empty response content")
	}

	if resp.Usage == nil {
		t.Error("Expected usage information")
	}
}
