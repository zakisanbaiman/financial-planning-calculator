package ports

import "context"

// StreamChunk はLLMからのストリームチャンクを表す
type StreamChunk struct {
	Token string
	Done  bool
	Error error
}

// LLMClient はLLMクライアントのインタフェース
type LLMClient interface {
	StreamAnswer(ctx context.Context, prompt string) (<-chan StreamChunk, error)
}
