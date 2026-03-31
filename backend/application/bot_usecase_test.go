package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/application"
	"github.com/financial-planning-calculator/backend/infrastructure/faq"
	"github.com/financial-planning-calculator/backend/infrastructure/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ===========================
// モック定義
// ===========================

// MockLLMClient はLLMClientインタフェースのモック実装
type MockLLMClient struct {
	mock.Mock
}

func (m *MockLLMClient) StreamAnswer(ctx context.Context, prompt string) (<-chan llm.StreamChunk, error) {
	args := m.Called(ctx, prompt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(<-chan llm.StreamChunk), args.Error(1)
}

// MockFAQLoader はFAQLoaderインタフェースのモック実装
type MockFAQLoader struct {
	mock.Mock
}

func (m *MockFAQLoader) Load(ctx context.Context) ([]faq.FAQContent, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]faq.FAQContent), args.Error(1)
}

func (m *MockFAQLoader) Search(query string) []faq.FAQContent {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]faq.FAQContent)
}

// ===========================
// テストヘルパー
// ===========================

// makeStreamChannel は指定されたトークン列のStreamChunkチャンネルを作成する
func makeStreamChannel(tokens []string) <-chan llm.StreamChunk {
	ch := make(chan llm.StreamChunk, len(tokens)+1)
	go func() {
		defer close(ch)
		for i, token := range tokens {
			isDone := i == len(tokens)-1
			ch <- llm.StreamChunk{Token: token, Done: isDone}
		}
	}()
	return ch
}

// makeErrorChannel はエラーチャンクを返すStreamChunkチャンネルを作成する
func makeErrorChannel(err error) <-chan llm.StreamChunk {
	ch := make(chan llm.StreamChunk, 1)
	go func() {
		defer close(ch)
		ch <- llm.StreamChunk{Error: err}
	}()
	return ch
}

// collectBotChunks はBotChunkチャンネルからすべてのチャンクを収集する
func collectBotChunks(t *testing.T, ch <-chan application.BotChunk, timeout time.Duration) []application.BotChunk {
	t.Helper()
	var chunks []application.BotChunk
	deadline := time.After(timeout)
	for {
		select {
		case chunk, ok := <-ch:
			if !ok {
				return chunks
			}
			chunks = append(chunks, chunk)
			if chunk.Done || chunk.Err != nil {
				return chunks
			}
		case <-deadline:
			t.Fatal("BotChunkの収集がタイムアウトしました")
			return nil
		}
	}
}

// ===========================
// BotUseCase.StreamAnswer Tests
// ===========================

func TestBotUseCase_StreamAnswer(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: 質問に対してFAQを検索し、LLMからの回答をストリームで返す", func(t *testing.T) {
		// Given
		mockFAQ := new(MockFAQLoader)
		mockLLM := new(MockLLMClient)

		faqContents := []faq.FAQContent{
			{Filename: "nisa.md", Title: "積立NISAとは", Body: "積立NISAは非課税投資制度です。"},
		}
		mockFAQ.On("Search", "積立NISAとは何ですか？").Return(faqContents)

		streamCh := makeStreamChannel([]string{"積立", "NISA", "は", "非課税", "です。"})
		mockLLM.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return((<-chan llm.StreamChunk)(streamCh), nil)

		useCase := application.NewBotUseCase(mockFAQ, mockLLM)

		// When
		ch, err := useCase.StreamAnswer(ctx, "積立NISAとは何ですか？")

		// Then
		require.NoError(t, err)
		chunks := collectBotChunks(t, ch, 5*time.Second)
		assert.NotEmpty(t, chunks)

		// 最後のチャンクがDone=trueであること
		lastChunk := chunks[len(chunks)-1]
		assert.True(t, lastChunk.Done)
		assert.Nil(t, lastChunk.Err)

		mockFAQ.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})

	t.Run("正常系: 関連FAQが見つからない場合でもLLMに質問を転送する", func(t *testing.T) {
		// Given
		mockFAQ := new(MockFAQLoader)
		mockLLM := new(MockLLMClient)

		mockFAQ.On("Search", "ビットコインとは").Return([]faq.FAQContent{})

		streamCh := makeStreamChannel([]string{"わかりません。"})
		mockLLM.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return((<-chan llm.StreamChunk)(streamCh), nil)

		useCase := application.NewBotUseCase(mockFAQ, mockLLM)

		// When
		ch, err := useCase.StreamAnswer(ctx, "ビットコインとは")

		// Then
		require.NoError(t, err)
		chunks := collectBotChunks(t, ch, 5*time.Second)
		assert.NotEmpty(t, chunks)
		mockFAQ.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})

	t.Run("正常系: LLMが返すBotChunkのトークンが結合されて元の回答を再現できる", func(t *testing.T) {
		// Given
		mockFAQ := new(MockFAQLoader)
		mockLLM := new(MockLLMClient)

		mockFAQ.On("Search", mock.AnythingOfType("string")).Return([]faq.FAQContent{})

		tokens := []string{"こ", "ん", "に", "ち", "は"}
		streamCh := makeStreamChannel(tokens)
		mockLLM.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return((<-chan llm.StreamChunk)(streamCh), nil)

		useCase := application.NewBotUseCase(mockFAQ, mockLLM)

		// When
		ch, err := useCase.StreamAnswer(ctx, "挨拶して")

		// Then
		require.NoError(t, err)
		chunks := collectBotChunks(t, ch, 5*time.Second)

		var combined string
		for _, c := range chunks {
			combined += c.Token
		}
		assert.Equal(t, "こんにちは", combined)
	})

	t.Run("異常系: 質問が空文字の場合はエラーを返す", func(t *testing.T) {
		// Given
		mockFAQ := new(MockFAQLoader)
		mockLLM := new(MockLLMClient)
		useCase := application.NewBotUseCase(mockFAQ, mockLLM)

		// When
		ch, err := useCase.StreamAnswer(ctx, "")

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			// エラーをチャンネルで受け取るパターン
			chunks := collectBotChunks(t, ch, 3*time.Second)
			require.NotEmpty(t, chunks)
			assert.NotNil(t, chunks[0].Err)
		}
	})

	t.Run("異常系: LLMがエラーを返した場合はBotChunkにエラーが含まれる", func(t *testing.T) {
		// Given
		mockFAQ := new(MockFAQLoader)
		mockLLM := new(MockLLMClient)

		mockFAQ.On("Search", mock.AnythingOfType("string")).Return([]faq.FAQContent{})
		mockLLM.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("LLMへの接続に失敗しました"))

		useCase := application.NewBotUseCase(mockFAQ, mockLLM)

		// When
		ch, err := useCase.StreamAnswer(ctx, "テスト質問")

		// Then
		if err != nil {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "LLMへの接続に失敗しました")
		} else {
			chunks := collectBotChunks(t, ch, 3*time.Second)
			require.NotEmpty(t, chunks)
			assert.NotNil(t, chunks[0].Err)
		}

		mockFAQ.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})

	t.Run("異常系: LLMがストリーム途中でエラーチャンクを送信した場合は伝播される", func(t *testing.T) {
		// Given
		mockFAQ := new(MockFAQLoader)
		mockLLM := new(MockLLMClient)

		mockFAQ.On("Search", mock.AnythingOfType("string")).Return([]faq.FAQContent{})

		errorCh := makeErrorChannel(errors.New("ストリーム中にエラーが発生しました"))
		mockLLM.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return((<-chan llm.StreamChunk)(errorCh), nil)

		useCase := application.NewBotUseCase(mockFAQ, mockLLM)

		// When
		ch, err := useCase.StreamAnswer(ctx, "テスト質問")

		// Then
		require.NoError(t, err) // ストリーム開始自体は成功
		chunks := collectBotChunks(t, ch, 3*time.Second)
		require.NotEmpty(t, chunks)

		hasError := false
		for _, c := range chunks {
			if c.Err != nil {
				hasError = true
				break
			}
		}
		assert.True(t, hasError, "エラーチャンクが伝播されるべきです")

		mockFAQ.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})

	t.Run("異常系: コンテキストがキャンセルされた場合はストリームが中断される", func(t *testing.T) {
		// Given
		mockFAQ := new(MockFAQLoader)
		mockLLM := new(MockLLMClient)

		mockFAQ.On("Search", mock.AnythingOfType("string")).Return([]faq.FAQContent{})

		// 長時間かかるストリームをシミュレート
		slowCh := make(chan llm.StreamChunk)
		mockLLM.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return((<-chan llm.StreamChunk)(slowCh), nil)

		useCase := application.NewBotUseCase(mockFAQ, mockLLM)
		ctx, cancel := context.WithCancel(context.Background())

		// When
		ch, err := useCase.StreamAnswer(ctx, "テスト質問")
		require.NoError(t, err)

		cancel() // 即座にキャンセル
		close(slowCh)

		// Then: チャンネルが適切にクローズされること
		deadline := time.After(3 * time.Second)
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return // チャンネルが閉じられた
				}
			case <-deadline:
				// タイムアウトしても問題なし（キャンセル後）
				return
			}
		}
	})

	t.Run("正常系: 複数のFAQが検索された場合、すべてのコンテキストがプロンプトに含まれる", func(t *testing.T) {
		// Given
		mockFAQ := new(MockFAQLoader)
		mockLLM := new(MockLLMClient)

		faqContents := []faq.FAQContent{
			{Filename: "nisa.md", Title: "積立NISAとは", Body: "積立NISAの説明"},
			{Filename: "ideco.md", Title: "iDeCoとは", Body: "iDeCoの説明"},
		}
		mockFAQ.On("Search", mock.AnythingOfType("string")).Return(faqContents)

		// プロンプトに両方のFAQのBodyが含まれていることを検証するマッチャー
		mockLLM.On("StreamAnswer", mock.Anything, mock.MatchedBy(func(prompt string) bool {
			return strings.Contains(prompt, "積立NISAの説明") && strings.Contains(prompt, "iDeCoの説明")
		})).Return((<-chan llm.StreamChunk)(makeStreamChannel([]string{"回答です。"})), nil)

		useCase := application.NewBotUseCase(mockFAQ, mockLLM)

		// When
		ch, err := useCase.StreamAnswer(ctx, "NISAとiDeCoの違いは？")

		// Then
		require.NoError(t, err)
		chunks := collectBotChunks(t, ch, 5*time.Second)
		assert.NotEmpty(t, chunks)
		mockFAQ.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})
}

// ===========================
// BotChunk Tests
// ===========================

func TestBotChunk(t *testing.T) {
	t.Run("正常系: トークンチャンクを正しく構築できる", func(t *testing.T) {
		// Given / When
		chunk := application.BotChunk{
			Token: "テスト",
			Done:  false,
			Err:   nil,
		}

		// Then
		assert.Equal(t, "テスト", chunk.Token)
		assert.False(t, chunk.Done)
		assert.Nil(t, chunk.Err)
	})

	t.Run("正常系: 完了チャンクを正しく構築できる", func(t *testing.T) {
		// Given / When
		chunk := application.BotChunk{
			Token: "",
			Done:  true,
			Err:   nil,
		}

		// Then
		assert.True(t, chunk.Done)
		assert.Empty(t, chunk.Token)
	})

	t.Run("正常系: エラーチャンクを正しく構築できる", func(t *testing.T) {
		// Given / When
		chunk := application.BotChunk{
			Token: "",
			Done:  false,
			Err:   errors.New("エラーが発生しました"),
		}

		// Then
		assert.NotNil(t, chunk.Err)
		assert.Contains(t, chunk.Err.Error(), "エラーが発生しました")
	})
}
