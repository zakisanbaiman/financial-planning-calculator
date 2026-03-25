package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// EmailService はメール送信サービスのインターフェース
type EmailService interface {
	SendPasswordResetEmail(ctx context.Context, toEmail, resetURL string) error
}

// LogEmailService は開発用のメールサービス（stdoutにログ出力）
type LogEmailService struct{}

// NewLogEmailService は開発用メールサービスを作成する
func NewLogEmailService() EmailService {
	return &LogEmailService{}
}

// SendPasswordResetEmail はリセットURLをログに出力する（開発用）
func (s *LogEmailService) SendPasswordResetEmail(_ context.Context, toEmail, resetURL string) error {
	slog.Info("パスワードリセットメール（開発モード）",
		"to", toEmail,
		"reset_url", resetURL,
	)
	return nil
}

// ResendEmailService はResend HTTP APIを使ったメールサービス
type ResendEmailService struct {
	apiKey string
	from   string
}

// NewResendEmailService はResendメールサービスを作成する
func NewResendEmailService(apiKey, from string) EmailService {
	return &ResendEmailService{
		apiKey: apiKey,
		from:   from,
	}
}

// SendPasswordResetEmail はResend APIでパスワードリセットメールを送信する
func (s *ResendEmailService) SendPasswordResetEmail(ctx context.Context, toEmail, resetURL string) error {
	body := fmt.Sprintf(`パスワードリセットのリクエストを受け付けました。

以下のURLをクリックしてパスワードをリセットしてください（有効期限: 30分）:

%s

このメールに心当たりがない場合は無視してください。
`, resetURL)

	payload := map[string]any{
		"from":    s.from,
		"to":      []string{toEmail},
		"subject": "パスワードリセットのご案内",
		"text":    body,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("リクエストの生成に失敗しました: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("メール送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Resend APIエラー: status=%d", resp.StatusCode)
	}

	return nil
}

// NewEmailService はAPI Key設定に基づいてメールサービスを作成する
// SMTP_PASSWORDをResend APIキーとして使用する
func NewEmailService(host string, port int, user, password, from string) EmailService {
	if password == "" {
		slog.Warn("SMTP_PASSWORDが未設定のため開発用メールサービス（ログ出力）を使用します")
		return NewLogEmailService()
	}
	return NewResendEmailService(password, from)
}
