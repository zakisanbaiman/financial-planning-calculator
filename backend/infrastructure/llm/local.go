package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/financial-planning-calculator/backend/application/ports"
)

const (
	// ollamaGenerateEndpoint はOllama互換APIのエンドポイントパス
	ollamaGenerateEndpoint = "/api/generate"

	// defaultModel はデフォルトで使用するモデル名
	defaultModel = "llama3"
)

// localLLMClient はOllama互換HTTPAPIを使ったLLMクライアント実装
type localLLMClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// ollamaRequest はOllama APIリクエストの構造体
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ollamaResponse はOllama APIレスポンスの1行分の構造体
type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// NewLocalLLMClient はローカルLLMクライアントを生成する
// baseURL: Ollama等のベースURL (例: "http://localhost:11434")
func NewLocalLLMClient(baseURL string) ports.LLMClient {
	return &localLLMClient{
		baseURL:    baseURL,
		model:      defaultModel,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
	}
}

// NewLocalLLMClientWithModel はモデル名を指定してローカルLLMクライアントを生成する
func NewLocalLLMClientWithModel(baseURL, model string) ports.LLMClient {
	return &localLLMClient{
		baseURL:    baseURL,
		model:      model,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
	}
}

// StreamAnswer はプロンプトを送信し、ストリームチャンクのチャンネルを返す
func (c *localLLMClient) StreamAnswer(ctx context.Context, prompt string) (<-chan ports.StreamChunk, error) {
	if prompt == "" {
		ch := make(chan ports.StreamChunk, 1)
		go func() {
			defer close(ch)
			ch <- ports.StreamChunk{Error: errors.New("プロンプトが空です")}
		}()
		return ch, nil
	}

	reqBody := ollamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストのエンコードに失敗しました: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+ollamaGenerateEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		ch := make(chan ports.StreamChunk, 1)
		go func() {
			defer close(ch)
			ch <- ports.StreamChunk{Error: fmt.Errorf("LLMサーバーへの接続に失敗しました: %w", err)}
		}()
		return ch, nil
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		ch := make(chan ports.StreamChunk, 1)
		go func() {
			defer close(ch)
			ch <- ports.StreamChunk{Error: fmt.Errorf("LLMサーバーがエラーを返しました: status=%d", resp.StatusCode)}
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

			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var ollamaResp ollamaResponse
			if err := json.Unmarshal(line, &ollamaResp); err != nil {
				ch <- ports.StreamChunk{Error: fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)}
				return
			}

			ch <- ports.StreamChunk{
				Token: ollamaResp.Response,
				Done:  ollamaResp.Done,
			}

			if ollamaResp.Done {
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
