package ai

import (
	"context"
	"strings"
	"sync"
)

// StreamAccumulator accumulates streaming responses
type StreamAccumulator struct {
	mu      sync.RWMutex
	content strings.Builder
	chunks  []StreamResponse
	done    bool
	err     error

	// Optional callbacks
	onChunk    func(chunk StreamResponse)
	onComplete func(content string)
	onError    func(error)
}

func NewStreamAccumulator() *StreamAccumulator {
	return &StreamAccumulator{
		chunks: make([]StreamResponse, 0),
	}
}

func (sa *StreamAccumulator) WithChunkCallback(fn func(chunk StreamResponse)) *StreamAccumulator {
	sa.onChunk = fn
	return sa
}

func (sa *StreamAccumulator) WithCompleteCallback(fn func(content string)) *StreamAccumulator {
	sa.onComplete = fn
	return sa
}

func (sa *StreamAccumulator) WithErrorCallback(fn func(error)) *StreamAccumulator {
	sa.onError = fn
	return sa
}

// Accumulate processes a stream channel and accumulates results
func (sa *StreamAccumulator) Accumulate(ctx context.Context, stream <-chan StreamResponse) {
	for {
		select {
		case <-ctx.Done():
			sa.mu.Lock()
			sa.err = ctx.Err()
			sa.done = true
			sa.mu.Unlock()
			if sa.onError != nil {
				sa.onError(ctx.Err())
			}
			return
		case chunk, ok := <-stream:
			if !ok {
				sa.mu.Lock()
				sa.done = true
				content := sa.content.String()
				sa.mu.Unlock()

				if sa.onComplete != nil {
					sa.onComplete(content)
				}
				return
			} else if chunk.Error != nil {
				sa.mu.Lock()
				sa.err = chunk.Error
				sa.done = true
				sa.mu.Unlock()

				if sa.onError != nil {
					sa.onError(chunk.Error)
				}
				return
			} else if chunk.Done {
				sa.mu.Lock()
				sa.done = true
				sa.content.WriteString(chunk.Content)
				sa.chunks = append(sa.chunks, chunk)
				content := sa.content.String()
				sa.mu.Unlock()

				if sa.onChunk != nil {
					sa.onChunk(chunk)
				}
				if sa.onComplete != nil {
					sa.onComplete(content)
				}
				return
			} else {
				sa.mu.Lock()
				sa.content.WriteString(chunk.Content)
				sa.chunks = append(sa.chunks, chunk)
				sa.mu.Unlock()

				if sa.onChunk != nil {
					sa.onChunk(chunk)
				}
			}
		}
	}
}

func (sa *StreamAccumulator) GetContent() string {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	return sa.content.String()
}
func (sa *StreamAccumulator) GetChunks() []StreamResponse {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	chunks := make([]StreamResponse, len(sa.chunks))
	copy(chunks, sa.chunks)
	return chunks
}

func (sa *StreamAccumulator) IsDone() bool {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	return sa.done
}

func (sa *StreamAccumulator) GetError() error {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	return sa.err
}

func CollectStream(ctx context.Context, stream <-chan StreamResponse) (string, error) {
	accumulator := NewStreamAccumulator()

	done := make(chan struct{})
	go func() {
		defer close(done)
		accumulator.Accumulate(ctx, stream)
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-done:
		if err := accumulator.GetError(); err != nil {
			return "", err
		}
		return accumulator.GetContent(), nil
	}
}
