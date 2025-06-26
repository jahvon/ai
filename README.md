# AI Adapter

<p>
    <a href="https://img.shields.io/github/v/release/jahvon/ai-adapter"><img src="https://img.shields.io/github/v/release/jahvon/ai-adapter" alt="GitHub release"></a>
    <a href="https://pkg.go.dev/github.com/jahvon/ai-adapter"><img src="https://pkg.go.dev/badge/github.com/jahvon/ai-adapter.svg" alt="Go Reference"></a>
</p>

A unified Go library for interacting with multiple AI providers through a consistent interface.
Currently only **OpenAI** and **Anthropic** are supported.

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jahvon/ai-adapter"
)

func main() {
    // Create configuration
    cfg := ai.NewConfig().
        WithProvider(ai.ProviderOpenAI).
        WithModel(ai.ModelGPT4oMini).
        WithAPIKeyEnv("OPENAI_API_KEY").
        WithMaxTokens(500)

    // Create adapter
    adapter, err := ai.NewAdapter(cfg)
    if err != nil {
        log.Fatalf("Failed to create adapter: %v", err)
    }

    // Create request
    req := &ai.Request{
        UserPrompt:   "Tell me a joke",
        SystemPrompt: "You are a helpful assistant",
    }

    // Generate response
    ctx := context.Background()
    resp, err := adapter.Generate(ctx, req)
    if err != nil {
        log.Fatalf("Failed to generate response: %v", err)
    }

    fmt.Printf("Response: %s\n", resp.Content)
    fmt.Printf("Model: %s\n", resp.Model)
    fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
}
```

### Streaming Response

```go
// Generate streaming response
stream, err := adapter.GenerateStream(ctx, req)
if err != nil {
    log.Fatalf("Failed to generate stream: %v", err)
}

// Collect all chunks into a single response
accumulated, err := ai.CollectStream(ctx, stream)
if err != nil {
    log.Fatalf("Failed to accumulate stream: %v", err)
}

fmt.Printf("Accumulated Response: %s\n", accumulated)
```
