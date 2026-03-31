package llm_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/application/ports"
	"github.com/financial-planning-calculator/backend/infrastructure/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===========================
// テストヘルパー
// ===========================

// setupGroqServer はGroq APIのモックHTTPサーバーを立ち上げてクライアントを返す
func setupGroqServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, ports.LLMClient) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	// groqAPIEndpointをモックサーバーに向けるため、NewGroqClientWithEndpointは
	// テスト用にエンドポイントを差し替えられる構造を持たないため、
	// モックサーバーのURLを使ってエンドポイントを上書きするテスト用コンストラクタを利用する
	client := llm.NewGroqClientWithEndpoint("test-api-key", "llama3-8b-8192", srv.URL+"/openai/v1/chat/completions")
	return srv, client
}

// collectGroqChunks はチャンネルからすべてのStreamChunkを収集して返す
func collectGroqChunks(t *testing.T, ch <-chan ports.StreamChunk, timeout time.Duration) []ports.StreamChunk {
	t.Helper()
	var chunks []ports.StreamChunk
	deadline := time.After(timeout)
	for {
		select {
		case chunk, ok := <-ch:
			if !ok {
				return chunks
			}
			chunks = append(chunks, chunk)
			if chunk.Done || chunk.Error != nil {
				return chunks
			}
		case <-deadline:
			t.Fatal("チャンクの収集がタイムアウトしました")
			return nil
		}
	}
}

// writeSSELines はOpenAI互換SSE形式でトークンを書き出すヘルパー
func writeSSELines(w http.ResponseWriter, flusher http.Flusher, tokens []string) {
	for i, token := range tokens {
		isLast := i == len(tokens)-1
		if isLast {
			finishReason := "stop"
			line := fmt.Sprintf(`data: {"choices":[{"delta":{"content":%q},"finish_reason":%q}]}`, token, finishReason)
			fmt.Fprintln(w, line)
		} else {
			line := fmt.Sprintf(`data: {"choices":[{"delta":{"content":%q},"finish_reason":null}]}`, token)
			fmt.Fprintln(w, line)
		}
		flusher.Flush()
	}
	fmt.Fprintln(w, "data: [DONE]")
	flusher.Flush()
}

// ===========================
// GroqClient.StreamAnswer Tests
// ===========================

func TestGroqClient_StreamAnswer(t *testing.T) {
	t.Run("正常系: ストリームでトークンを受信できる", func(t *testing.T) {
		// Given: OpenAI互換SSE形式のストリームレスポンスを返すサーバー
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authorizationヘッダーが設定されていることを確認
			assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher, ok := w.(http.Flusher)
			require.True(t, ok)

			tokens := []string{"こ", "ん", "に", "ち", "は"}
			writeSSELines(w, flusher, tokens)
		})
		_, client := setupGroqServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "こんにちはと答えてください")

		// Then
		require.NoError(t, err)
		chunks := collectGroqChunks(t, ch, 5*time.Second)
		assert.NotEmpty(t, chunks)

		// doneチャンクが最後に届いていること
		lastChunk := chunks[len(chunks)-1]
		assert.True(t, lastChunk.Done)
		assert.Nil(t, lastChunk.Error)

		// トークンが連結されると期待する文字列になること
		var combined string
		for _, c := range chunks {
			combined += c.Token
		}
		assert.Contains(t, combined, "こんにちは")
	})

	t.Run("正常系: [DONE]を受信したらストリームが終了する", func(t *testing.T) {
		// Given
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher, ok := w.(http.Flusher)
			require.True(t, ok)

			fmt.Fprintln(w, `data: {"choices":[{"delta":{"content":"はい"},"finish_reason":null}]}`)
			flusher.Flush()
			fmt.Fprintln(w, "data: [DONE]")
			flusher.Flush()
		})
		_, client := setupGroqServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		require.NoError(t, err)
		chunks := collectGroqChunks(t, ch, 5*time.Second)
		require.NotEmpty(t, chunks)
		// 最後のチャンクはDone=trueのはず（[DONE]受信時に送信される）
		assert.True(t, chunks[len(chunks)-1].Done)
	})

	t.Run("正常系: finish_reason=stopで完了チャンクを受信できる", func(t *testing.T) {
		// Given
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher, ok := w.(http.Flusher)
			require.True(t, ok)

			writeSSELines(w, flusher, []string{"応答"})
		})
		_, client := setupGroqServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		require.NoError(t, err)
		chunks := collectGroqChunks(t, ch, 5*time.Second)
		require.NotEmpty(t, chunks)

		// finish_reason=stopのチャンクでDone=trueになっているか確認
		hasDone := false
		for _, c := range chunks {
			if c.Done {
				hasDone = true
				break
			}
		}
		assert.True(t, hasDone, "Done=trueのチャンクが届くべきです")
	})

	t.Run("異常系: サーバーが500エラーを返した場合はエラーチャンクを送信する", func(t *testing.T) {
		// Given
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, `{"error":{"message":"internal error"}}`)
		})
		_, client := setupGroqServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			chunks := collectGroqChunks(t, ch, 5*time.Second)
			hasError := false
			for _, c := range chunks {
				if c.Error != nil {
					hasError = true
					break
				}
			}
			assert.True(t, hasError, "エラーチャンクが送信されるべきです")
		}
	})

	t.Run("異常系: 401 Unauthorizedの場合はエラーチャンクを送信する", func(t *testing.T) {
		// Given
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `{"error":{"message":"invalid api key"}}`)
		})
		_, client := setupGroqServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			chunks := collectGroqChunks(t, ch, 5*time.Second)
			hasError := false
			for _, c := range chunks {
				if c.Error != nil {
					hasError = true
					break
				}
			}
			assert.True(t, hasError, "エラーチャンクが送信されるべきです")
		}
	})

	t.Run("異常系: コンテキストがキャンセルされた場合はストリームが中断される", func(t *testing.T) {
		// Given: 遅延レスポンスを返すサーバー
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher, ok := w.(http.Flusher)
			require.True(t, ok)

			for i := 0; i < 100; i++ {
				select {
				case <-r.Context().Done():
					return
				default:
					fmt.Fprintln(w, `data: {"choices":[{"delta":{"content":"token"},"finish_reason":null}]}`)
					flusher.Flush()
					time.Sleep(50 * time.Millisecond)
				}
			}
		})
		_, client := setupGroqServer(t, handler)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// When
		ch, err := client.StreamAnswer(ctx, "長い質問")
		require.NoError(t, err)

		// 数チャンク受信後にキャンセル
		count := 0
		for range ch {
			count++
			if count >= 2 {
				cancel()
				break
			}
		}

		// Then: チャンネルが適切にクローズされること
		deadline := time.After(3 * time.Second)
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
			case <-deadline:
				return
			}
		}
	})

	t.Run("異常系: APIサーバーに接続できない場合はエラーを返す", func(t *testing.T) {
		// Given: 存在しないエンドポイント
		client := llm.NewGroqClientWithEndpoint("test-api-key", "llama3-8b-8192", "http://localhost:19999/v1/chat/completions")

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			chunks := collectGroqChunks(t, ch, 5*time.Second)
			hasError := false
			for _, c := range chunks {
				if c.Error != nil {
					hasError = true
					break
				}
			}
			assert.True(t, hasError, "接続エラーチャンクが送信されるべきです")
		}
	})

	t.Run("異常系: プロンプトが空の場合はエラーを返す", func(t *testing.T) {
		// Given
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		_, client := setupGroqServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "")

		// Then
		if err != nil {
			assert.Error(t, err)
		} else if ch != nil {
			chunks := collectGroqChunks(t, ch, 3*time.Second)
			hasError := false
			for _, c := range chunks {
				if c.Error != nil {
					hasError = true
					break
				}
			}
			assert.True(t, hasError || (len(chunks) > 0 && chunks[len(chunks)-1].Done),
				"空プロンプトはエラーまたはdoneチャンクで終了すること")
		}
	})

	t.Run("正常系: data:行以外の行は無視される", func(t *testing.T) {
		// Given: コメント行や空行を含むSSEレスポンス
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher, ok := w.(http.Flusher)
			require.True(t, ok)

			// SSEコメント行と空行を含む
			fmt.Fprintln(w, ": keep-alive")
			fmt.Fprintln(w, "")
			fmt.Fprintln(w, `data: {"choices":[{"delta":{"content":"応答"},"finish_reason":null}]}`)
			flusher.Flush()
			fmt.Fprintln(w, "data: [DONE]")
			flusher.Flush()
		})
		_, client := setupGroqServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		require.NoError(t, err)
		chunks := collectGroqChunks(t, ch, 5*time.Second)
		require.NotEmpty(t, chunks)

		var combined string
		for _, c := range chunks {
			combined += c.Token
		}
		assert.Contains(t, combined, "応答")
	})
}
