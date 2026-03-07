package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
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

// SMTPEmailService はSMTPを使ったメールサービス
type SMTPEmailService struct {
	host     string
	port     int
	user     string
	password string
	from     string
}

// NewSMTPEmailService はSMTPメールサービスを作成する
func NewSMTPEmailService(host string, port int, user, password, from string) EmailService {
	return &SMTPEmailService{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
	}
}

// SendPasswordResetEmail はパスワードリセットメールをSMTPで送信する
func (s *SMTPEmailService) SendPasswordResetEmail(_ context.Context, toEmail, resetURL string) error {
	subject := "パスワードリセットのご案内"
	body := fmt.Sprintf(`パスワードリセットのリクエストを受け付けました。

以下のURLをクリックしてパスワードをリセットしてください（有効期限: 30分）:

%s

このメールに心当たりがない場合は無視してください。
`, resetURL)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.from, toEmail, subject, body,
	))

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	if err := smtp.SendMail(addr, auth, s.from, []string{toEmail}, msg); err != nil {
		return fmt.Errorf("メール送信に失敗しました: %w", err)
	}
	return nil
}

// NewEmailService はSMTP設定に基づいてメールサービスを作成する
// SMTP設定がない場合はログ出力のフォールバックを使用する
func NewEmailService(host string, port int, user, password, from string) EmailService {
	if host == "" {
		slog.Warn("SMTP設定がないため開発用メールサービス（ログ出力）を使用します")
		return NewLogEmailService()
	}
	return NewSMTPEmailService(host, port, user, password, from)
}
