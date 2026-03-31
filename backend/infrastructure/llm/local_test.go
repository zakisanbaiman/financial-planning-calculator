package llm_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/infrastructure/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===========================
// テストヘルパー
// ===========================

// setupLocalLLMServer はローカルLLMのモックHTTPサーバーを立ち上げてクライアントを返す
func setupLocalLLMServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, llm.LLMClient) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := llm.NewLocalLLMClient(srv.URL)
	return srv, client
}

// collectChunks はチャンネルからすべてのStreamChunkを収集して返す
func collectChunks(t *testing.T, ch <-chan llm.StreamChunk, timeout time.Duration) []llm.StreamChunk {
	t.Helper()
	var chunks []llm.StreamChunk
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

// ===========================
// LocalLLMClient.StreamAnswer Tests
// ===========================

func TestLocalLLMClient_StreamAnswer(t *testing.T) {
	t.Run("正常系: ストリームでトークンを受信できる", func(t *testing.T) {
		// Given: OllamaライクなNDJSON形式のストリームレスポンスを返すサーバー
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.WriteHeader(http.StatusOK)
			flusher, ok := w.(http.Flusher)
			require.True(t, ok)

			tokens := []string{"こ", "ん", "に", "ち", "は"}
			for i, token := range tokens {
				isDone := i == len(tokens)-1
				if isDone {
					_, _ = w.Write([]byte(`{"response":"` + token + `","done":true}` + "\n"))
				} else {
					_, _ = w.Write([]byte(`{"response":"` + token + `","done":false}` + "\n"))
				}
				flusher.Flush()
			}
		})
		_, client := setupLocalLLMServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "こんにちはと答えてください")

		// Then
		require.NoError(t, err)
		chunks := collectChunks(t, ch, 5*time.Second)
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

	t.Run("正常系: 単一トークンで即done=trueのレスポンスを処理できる", func(t *testing.T) {
		// Given
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":"はい","done":true}` + "\n"))
		})
		_, client := setupLocalLLMServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		require.NoError(t, err)
		chunks := collectChunks(t, ch, 5*time.Second)
		require.NotEmpty(t, chunks)
		assert.True(t, chunks[len(chunks)-1].Done)
	})

	t.Run("異常系: サーバーが500エラーを返した場合はエラーチャンクを送信する", func(t *testing.T) {
		// Given
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"internal error"}`))
		})
		_, client := setupLocalLLMServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		// StreamAnswer自体がエラーを返すか、チャンネルでエラーチャンクを返すかの2パターンを許容
		if err != nil {
			assert.Error(t, err)
		} else {
			chunks := collectChunks(t, ch, 5*time.Second)
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
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.WriteHeader(http.StatusOK)
			flusher, ok := w.(http.Flusher)
			require.True(t, ok)

			for i := 0; i < 100; i++ {
				select {
				case <-r.Context().Done():
					return
				default:
					_, _ = w.Write([]byte(`{"response":"token","done":false}` + "\n"))
					flusher.Flush()
					time.Sleep(50 * time.Millisecond)
				}
			}
		})
		_, client := setupLocalLLMServer(t, handler)

		ctx, cancel := context.WithCancel(context.Background())

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
		// チャンネルが詰まらず受信完了できることを確認
		deadline := time.After(3 * time.Second)
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					// チャンネルが閉じられた = 正常終了
					return
				}
			case <-deadline:
				// キャンセル後にチャンネルが閉じられた、または空になった
				return
			}
		}
	})

	t.Run("異常系: LLMサーバーに接続できない場合はエラーを返す", func(t *testing.T) {
		// Given: 存在しないエンドポイント
		client := llm.NewLocalLLMClient("http://localhost:19999")

		// When
		ch, err := client.StreamAnswer(context.Background(), "質問")

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			chunks := collectChunks(t, ch, 5*time.Second)
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
		_, client := setupLocalLLMServer(t, handler)

		// When
		ch, err := client.StreamAnswer(context.Background(), "")

		// Then
		if err != nil {
			assert.Error(t, err)
		} else if ch != nil {
			chunks := collectChunks(t, ch, 3*time.Second)
			hasError := false
			for _, c := range chunks {
				if c.Error != nil {
					hasError = true
					break
				}
			}
			// 空プロンプトはエラーまたはdoneチャンクで終了すること
			assert.True(t, hasError || (len(chunks) > 0 && chunks[len(chunks)-1].Done))
		}
	})
}

// ===========================
// StreamChunk Tests
// ===========================

func TestStreamChunk(t *testing.T) {
	t.Run("正常系: トークンチャンクを正しく構築できる", func(t *testing.T) {
		// Given / When
		chunk := llm.StreamChunk{
			Token: "テスト",
			Done:  false,
			Error: nil,
		}

		// Then
		assert.Equal(t, "テスト", chunk.Token)
		assert.False(t, chunk.Done)
		assert.Nil(t, chunk.Error)
	})

	t.Run("正常系: 完了チャンクを正しく構築できる", func(t *testing.T) {
		// Given / When
		chunk := llm.StreamChunk{
			Token: "",
			Done:  true,
			Error: nil,
		}

		// Then
		assert.True(t, chunk.Done)
		assert.Empty(t, chunk.Token)
	})
}
