package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/financial-planning-calculator/backend/application"
	"github.com/labstack/echo/v4"
)

// sseContentType はSSEのContent-Typeヘッダー値
const sseContentType = "text/event-stream"

// BotController はBot関連のHTTPハンドラーを提供する
type BotController struct {
	useCase application.BotUseCase
}

// NewBotController はBotControllerを生成する
func NewBotController(useCase application.BotUseCase) *BotController {
	return &BotController{useCase: useCase}
}

// PostMessageRequest はBotメッセージリクエストの構造体
type PostMessageRequest struct {
	Question string `json:"question" validate:"required,max=2000"`
}

// PostMessage はBotへの質問を受け取り、SSEストリームで回答を返す
// POST /api/bot/messages
func (c *BotController) PostMessage(ctx echo.Context) error {
	// JWT認証済みチェック
	userID, ok := ctx.Get("user_id").(string)
	if !ok || userID == "" {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "認証が必要です",
		})
	}

	// リクエストのバインド・バリデーション
	var req PostMessageRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストの解析に失敗しました",
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "questionフィールドは必須です",
		})
	}

	// UseCaseからストリームを取得
	reqCtx := ctx.Request().Context()
	ch, err := c.useCase.StreamAnswer(reqCtx, req.Question)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// SSEヘッダーを設定
	w := ctx.Response()
	w.Header().Set("Content-Type", sseContentType)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	// ストリームをSSEイベントとして書き込む
	flusher, canFlush := w.Writer.(http.Flusher)

	for {
		select {
		case <-reqCtx.Done():
			return nil
		case chunk, ok := <-ch:
			if !ok {
				writeSSEEvent(w.Writer, "done", map[string]interface{}{})
				if canFlush {
					flusher.Flush()
				}
				return nil
			}
			if chunk.Err != nil {
				slog.Error("SSEストリーム中にエラーが発生しました", slog.Any("error", chunk.Err))
				writeSSEEvent(w.Writer, "error", map[string]string{
					"message": "回答の生成中にエラーが発生しました",
				})
				if canFlush {
					flusher.Flush()
				}
				return nil
			}
			if chunk.Token != "" {
				writeSSEEvent(w.Writer, "message", map[string]string{
					"token": chunk.Token,
				})
				if canFlush {
					flusher.Flush()
				}
			}
			if chunk.Done {
				writeSSEEvent(w.Writer, "done", map[string]interface{}{})
				if canFlush {
					flusher.Flush()
				}
				return nil
			}
		}
	}
}

// writeSSEEvent はSSEイベントをWriterに書き込む
func writeSSEEvent(w interface{ Write([]byte) (int, error) }, event string, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, string(dataBytes))
}
