package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/financial-planning-calculator/backend/application/ports"
)

const (
	groqAPIEndpoint  = "https://api.groq.com/openai/v1/chat/completions"
	defaultGroqModel = "llama3-8b-8192"
)

// groqClient はGroq APIを使ったLLMクライアント実装
type groqClient struct {
	apiKey     string
	model      string
	endpoint   string
	httpClient *http.Client
}

// groqRequest はGroq APIリクエストの構造体
type groqRequest struct {
	Model    string        `json:"model"`
	Messages []groqMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// groqMessage はGroq APIのメッセージ構造体
type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// groqSSEResponse はGroq APIのSSEレスポンス1行分の構造体
type groqSSEResponse struct {
	Choices []groqChoice `json:"choices"`
}

// groqChoice はGroq APIのchoice構造体
type groqChoice struct {
	Delta        groqDelta `json:"delta"`
	FinishReason *string   `json:"finish_reason"`
}

// groqDelta はGroq APIのdelta構造体
type groqDelta struct {
	Content string `json:"content"`
}

// NewGroqClient はGroq APIクライアントを生成する
// apiKey: Groq APIキー
// model: 使用するモデル名（空文字の場合はデフォルトを使用）
func NewGroqClient(apiKey, model string) ports.LLMClient {
	return NewGroqClientWithEndpoint(apiKey, model, groqAPIEndpoint)
}

// NewGroqClientWithEndpoint はエンドポイントを指定してGroq APIクライアントを生成する（主にテスト用）
func NewGroqClientWithEndpoint(apiKey, model, endpoint string) ports.LLMClient {
	if model == "" {
		model = defaultGroqModel
	}
	return &groqClient{
		apiKey:     apiKey,
		model:      model,
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
	}
}

// StreamAnswer はプロンプトを送信し、ストリームチャンクのチャンネルを返す
func (c *groqClient) StreamAnswer(ctx context.Context, prompt string) (<-chan ports.StreamChunk, error) {
	if prompt == "" {
		ch := make(chan ports.StreamChunk, 1)
		go func() {
			defer close(ch)
			ch <- ports.StreamChunk{Error: errors.New("プロンプトが空です")}
		}()
		return ch, nil
	}

	reqBody := groqRequest{
		Model: c.model,
		Messages: []groqMessage{
			{Role: "user", Content: prompt},
		},
		Stream: true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストのエンコードに失敗しました: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		ch := make(chan ports.StreamChunk, 1)
		go func() {
			defer close(ch)
			ch <- ports.StreamChunk{Error: fmt.Errorf("Groq APIサーバーへの接続に失敗しました: %w", err)}
		}()
		return ch, nil
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		ch := make(chan ports.StreamChunk, 1)
		go func() {
			defer close(ch)
			ch <- ports.StreamChunk{Error: fmt.Errorf("Groq APIサーバーがエラーを返しました: status=%d", resp.StatusCode)}
		}()
		return ch, nil
	}

	ch := make(chan ports.StreamChunk, 16)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line := scanner.Text()
			if line == "" {
				continue
			}

			// SSE形式: "data: {...}" または "data: [DONE]"
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")

			if data == "[DONE]" {
				ch <- ports.StreamChunk{Done: true}
				return
			}

			var sseResp groqSSEResponse
			if err := json.Unmarshal([]byte(data), &sseResp); err != nil {
				ch <- ports.StreamChunk{Error: fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)}
				return
			}

			if len(sseResp.Choices) == 0 {
				continue
			}

			choice := sseResp.Choices[0]
			isDone := choice.FinishReason != nil && *choice.FinishReason == "stop"

			ch <- ports.StreamChunk{
				Token: choice.Delta.Content,
				Done:  isDone,
			}

			if isDone {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- ports.StreamChunk{Error: fmt.Errorf("ストリーム読み込み中にエラーが発生しました: %w", err)}
			}
		}
	}()

	return ch, nil
}
