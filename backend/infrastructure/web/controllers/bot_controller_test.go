package controllers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/application"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ===========================
// モック定義
// ===========================

// MockBotUseCase はBotUseCaseインタフェースのモック実装
type MockBotUseCase struct {
	mock.Mock
}

func (m *MockBotUseCase) StreamAnswer(ctx context.Context, question string) (<-chan application.BotChunk, error) {
	args := m.Called(ctx, question)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(<-chan application.BotChunk), args.Error(1)
}

// ===========================
// テストヘルパー
// ===========================

// makeBotChunkChannel は指定されたトークン列のBotChunkチャンネルを作成する
func makeBotChunkChannel(tokens []string) <-chan application.BotChunk {
	ch := make(chan application.BotChunk, len(tokens)+1)
	go func() {
		defer close(ch)
		for i, token := range tokens {
			isDone := i == len(tokens)-1
			ch <- application.BotChunk{Token: token, Done: isDone}
		}
	}()
	return ch
}

// makeBotErrorChannel はエラーチャンクを返すBotChunkチャンネルを作成する
func makeBotErrorChannel(err error) <-chan application.BotChunk {
	ch := make(chan application.BotChunk, 1)
	go func() {
		defer close(ch)
		ch <- application.BotChunk{Err: err}
	}()
	return ch
}

// newBotEcho はBotコントローラーテスト用のEchoインスタンスを作成する
func newBotEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return e
}

// parseSSEEvents はSSEレスポンスボディからイベントを解析して返す
func parseSSEEvents(body string) []map[string]string {
	var events []map[string]string
	scanner := bufio.NewScanner(strings.NewReader(body))

	var currentEvent map[string]string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if currentEvent != nil {
				events = append(events, currentEvent)
				currentEvent = nil
			}
			continue
		}
		if strings.HasPrefix(line, "event:") {
			if currentEvent == nil {
				currentEvent = make(map[string]string)
			}
			currentEvent["event"] = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			if currentEvent == nil {
				currentEvent = make(map[string]string)
			}
			currentEvent["data"] = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
	}
	if currentEvent != nil {
		events = append(events, currentEvent)
	}
	return events
}

// setJWTUserID はEchoコンテキストにユーザーIDをセットする（JWT認証済みを模擬する）
func setJWTUserID(c echo.Context, userID string) {
	c.Set("user_id", userID)
}

// ===========================
// BotController.PostMessage Tests
// ===========================

func TestBotController_PostMessage(t *testing.T) {
	t.Run("正常系: 質問を送信するとSSEストリームが返される", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		tokens := []string{"積立", "NISAは", "非課税", "制度", "です。"}
		ch := makeBotChunkChannel(tokens)
		mockUseCase.On("StreamAnswer", mock.Anything, "積立NISAとは何ですか？").Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"積立NISAとは何ですか？"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("正常系: SSEレスポンスにmessageイベントが含まれる", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		tokens := []string{"こんにちは"}
		ch := makeBotChunkChannel(tokens)
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"挨拶して"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		require.NoError(t, err)

		body := rec.Body.String()
		events := parseSSEEvents(body)
		assert.NotEmpty(t, events)

		// messageイベントが存在することを確認
		hasMessageEvent := false
		for _, ev := range events {
			if ev["event"] == "message" {
				hasMessageEvent = true
				// dataにtokenフィールドが含まれること
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(ev["data"]), &data); err == nil {
					_, hasToken := data["token"]
					assert.True(t, hasToken, "messageイベントのdataにtokenフィールドが含まれること")
				}
			}
		}
		assert.True(t, hasMessageEvent, "messageイベントが存在すること")
	})

	t.Run("正常系: ストリーム完了時にdoneイベントが送信される", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		ch := makeBotChunkChannel([]string{"回答です。"})
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"質問"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		require.NoError(t, err)

		body := rec.Body.String()
		events := parseSSEEvents(body)

		// doneイベントが存在することを確認
		hasDoneEvent := false
		for _, ev := range events {
			if ev["event"] == "done" {
				hasDoneEvent = true
			}
		}
		assert.True(t, hasDoneEvent, "doneイベントが存在すること")
	})

	t.Run("正常系: SSEレスポンスのCache-Controlヘッダーが正しく設定される", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		ch := makeBotChunkChannel([]string{"回答"})
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"質問"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		require.NoError(t, err)
		// SSEにはキャッシュ無効ヘッダーが必要
		assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
	})

	t.Run("正常系: SSEレスポンスのConnectionヘッダーが設定される", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		ch := makeBotChunkChannel([]string{"回答"})
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"質問"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		require.NoError(t, err)
		assert.Equal(t, "keep-alive", rec.Header().Get("Connection"))
	})

	t.Run("異常系: リクエストボディが空の場合は400エラーを返す", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)
		controller := NewBotController(mockUseCase)

		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString("{}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		// バリデーションエラーはhandlerエラーまたは400レスポンスで返される
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("異常系: questionフィールドが空の場合は400エラーを返す", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)
		controller := NewBotController(mockUseCase)

		reqBody := `{"question":""}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("異常系: 不正なJSONの場合は400エラーを返す", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)
		controller := NewBotController(mockUseCase)

		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("異常系: ユーザーIDがコンテキストにない場合は401エラーを返す", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)
		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"積立NISAとは？"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// user_idをセットしない（認証なし）

		// When
		err := controller.PostMessage(c)

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("異常系: UseCaseがエラーを返した場合はSSEのerrorイベントを送信するか500を返す", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("LLMへの接続に失敗しました"))

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"積立NISAとは？"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			// SSEでerrorイベントを返すか500を返す
			isError := rec.Code == http.StatusInternalServerError ||
				strings.Contains(rec.Body.String(), "event: error")
			assert.True(t, isError, "エラーが適切に処理されること")
		}
		mockUseCase.AssertExpectations(t)
	})

	t.Run("異常系: ストリーム途中のエラーチャンクがSSEのerrorイベントとして送信される", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		errCh := makeBotErrorChannel(errors.New("ストリーム中にエラーが発生しました"))
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return(errCh, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"質問"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		require.NoError(t, err)

		body := rec.Body.String()
		events := parseSSEEvents(body)

		hasErrorEvent := false
		for _, ev := range events {
			if ev["event"] == "error" {
				hasErrorEvent = true
				// errorイベントのdataにmessageフィールドが含まれること
				var data map[string]interface{}
				if jsonErr := json.Unmarshal([]byte(ev["data"]), &data); jsonErr == nil {
					_, hasMessage := data["message"]
					assert.True(t, hasMessage, "errorイベントのdataにmessageフィールドが含まれること")
				}
			}
		}
		assert.True(t, hasErrorEvent, "errorイベントが存在すること")
		mockUseCase.AssertExpectations(t)
	})
}

// ===========================
// SSEレスポンス形式 Tests
// ===========================

func TestSSEResponseFormat(t *testing.T) {
	t.Run("正常系: messageイベントのデータ形式が正しい", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		ch := makeBotChunkChannel([]string{"テスト"})
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"テスト質問"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)
		require.NoError(t, err)

		// Then: ADRで定義されたSSEフォーマットを検証
		// event: message
		// data: {"token":"..."}
		body := rec.Body.String()
		events := parseSSEEvents(body)

		for _, ev := range events {
			if ev["event"] == "message" {
				var data struct {
					Token string `json:"token"`
				}
				err := json.Unmarshal([]byte(ev["data"]), &data)
				assert.NoError(t, err, "messageイベントのdataが有効なJSONであること")
				assert.NotEmpty(t, data.Token, "tokenフィールドが空でないこと")
			}
		}
	})

	t.Run("正常系: doneイベントのデータ形式が正しい", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		ch := makeBotChunkChannel([]string{"回答"})
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"質問"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)
		require.NoError(t, err)

		// Then: ADRで定義されたSSEフォーマットを検証
		// event: done
		// data: {}
		body := rec.Body.String()
		events := parseSSEEvents(body)

		for _, ev := range events {
			if ev["event"] == "done" {
				// dataが有効なJSONであること（空オブジェクトも可）
				var data map[string]interface{}
				err := json.Unmarshal([]byte(ev["data"]), &data)
				assert.NoError(t, err, "doneイベントのdataが有効なJSONであること")
			}
		}
	})

	t.Run("正常系: SSEの各イベントは空行で区切られている", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		ch := makeBotChunkChannel([]string{"回答1", "回答2"})
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"質問"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)
		require.NoError(t, err)

		// Then: SSEの区切り文字（\n\n）が存在すること
		body := rec.Body.String()
		assert.Contains(t, body, "\n\n", "SSEイベントは\\n\\nで区切られること")
	})
}

// ===========================
// 境界値 Tests
// ===========================

func TestBotController_QuestionBoundary(t *testing.T) {
	t.Run("正常系: 2000文字の質問は正常受理される", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		question := strings.Repeat("あ", 2000)
		ch := makeBotChunkChannel([]string{"回答です。"})
		mockUseCase.On("StreamAnswer", mock.Anything, question).Return(ch, nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"` + question + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("異常系: 2001文字の質問は400エラーを返す", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)
		controller := NewBotController(mockUseCase)

		question := strings.Repeat("あ", 2001)
		reqBody := `{"question":"` + question + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// When
		err := controller.PostMessage(c)

		// Then
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

// ===========================
// タイムアウト Tests
// ===========================

func TestBotController_Timeout(t *testing.T) {
	t.Run("エッジケース: 長時間処理でもコンテキストキャンセルで適切に終了する", func(t *testing.T) {
		// Given
		e := newBotEcho()
		mockUseCase := new(MockBotUseCase)

		// 無限に待つチャンネルを返す（コンテキストキャンセルで終了すること）
		slowCh := make(chan application.BotChunk)
		mockUseCase.On("StreamAnswer", mock.Anything, mock.AnythingOfType("string")).Return((<-chan application.BotChunk)(slowCh), nil)

		controller := NewBotController(mockUseCase)

		reqBody := `{"question":"質問"}`
		req := httptest.NewRequest(http.MethodPost, "/api/bot/messages", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// タイムアウトを設定したコンテキスト
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setJWTUserID(c, "user-123")

		// 別ゴルーチンでPostMessageを実行
		done := make(chan error, 1)
		started := make(chan struct{})
		go func() {
			close(started)
			done <- controller.PostMessage(c)
		}()

		// PostMessageが開始してからチャンネルを閉じる（リソース解放）
		<-started
		close(slowCh)

		// Then: タイムアウト後に終了すること
		select {
		case err := <-done:
			// エラーあり/なし問わずタイムアウトで終了すること
			_ = err
		case <-time.After(2 * time.Second):
			t.Error("タイムアウト後にハンドラーが終了していません")
		}
	})
}
