package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Port                string
	Debug               bool
	AllowedOrigins      []string
	CORSMaxAge          int
	RateLimitRPS        int
	RateLimitBurst      int
	AuthRateLimitRPS    int
	AuthRateLimitBurst  int
	RequestTimeout      time.Duration
	MaxRequestSize      string
	EnableGzip          bool
	GzipLevel           int
	LogFormat           string
	EnableSecureHeaders bool
	EnablePprof         bool
	PprofPort           string
	TempFileDir         string
	TempFileSecret      string
	TempFileExpiry      time.Duration
	CleanupInterval     time.Duration
	// Basic Authentication
	EnableBasicAuth     bool
	BasicAuthUsername   string
	BasicAuthPassword   string
	// JWT Authentication
	JWTSecret                string
	JWTExpiration            time.Duration
	RefreshTokenExpiration   time.Duration
	// GitHub OAuth
	GitHubClientID           string
	GitHubClientSecret       string
	GitHubCallbackURL        string
	OAuthSuccessRedirect     string
	OAuthFailureRedirect     string
	// Cookie Security
	CookieSecure             bool
	// WebAuthn Settings
	WebAuthnRPID             string // Relying Party ID (e.g., "example.com")
	WebAuthnRPName           string // Relying Party Name (e.g., "財務計画計算機")
	WebAuthnRPOrigin         string // Relying Party Origin (e.g., "https://example.com")
	// CSP
	ContentSecurityPolicy   string // Content-Security-Policy ヘッダー値（空文字の場合はヘッダーを設定しない）
	// SMTP メール設定
	SMTPHost     string // SMTP_HOST
	SMTPPort     int    // SMTP_PORT
	SMTPUser     string // SMTP_USER
	SMTPPassword string // SMTP_PASSWORD
	SMTPFrom     string // SMTP_FROM
	// フロントエンドURL（パスワードリセットURLの生成に使用）
	FrontendURL  string // FRONTEND_URL
	// Bot LLM設定
	GroqAPIKey string // GROQ_API_KEY
	GroqModel  string // GROQ_MODEL (例: "llama3-8b-8192")
	FAQDir     string // FAQ_DIR (例: "docs/faq")
}

// LoadServerConfig loads server configuration from environment variables
func LoadServerConfig() *ServerConfig {
	config := &ServerConfig{
		Port:                getEnv("PORT", "8080"),
		Debug:               getEnvBool("DEBUG", false),
		AllowedOrigins:      getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001", "https://localhost:3000", "https://localhost:3001"}),
		CORSMaxAge:          getEnvInt("CORS_MAX_AGE", 86400),
		RateLimitRPS:        getEnvInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst:      getEnvInt("RATE_LIMIT_BURST", 50),
		AuthRateLimitRPS:    getEnvInt("AUTH_RATE_LIMIT_RPS", 10),
		AuthRateLimitBurst:  getEnvInt("AUTH_RATE_LIMIT_BURST", 5),
		RequestTimeout:      getEnvDuration("REQUEST_TIMEOUT", 30*time.Second),
		MaxRequestSize:      getEnv("MAX_REQUEST_SIZE", "10M"),
		EnableGzip:          getEnvBool("ENABLE_GZIP", true),
		GzipLevel:           getEnvInt("GZIP_LEVEL", 5),
		LogFormat:           getEnv("LOG_FORMAT", "${time_rfc3339} ${method} ${uri} ${status} ${latency_human} ${bytes_in}B/${bytes_out}B ${error}\n"),
		EnableSecureHeaders: getEnvBool("ENABLE_SECURE_HEADERS", true),
		EnablePprof:         getEnvBool("ENABLE_PPROF", false),
		PprofPort:           getEnv("PPROF_PORT", "6060"),
		TempFileDir:         getEnv("TEMP_FILE_DIR", "/tmp/financial-planning-reports"),
		TempFileSecret:      getEnv("TEMP_FILE_SECRET", "change-this-secret-in-production"),
		TempFileExpiry:      getEnvDuration("TEMP_FILE_EXPIRY", 24*time.Hour),
		CleanupInterval:     getEnvDuration("CLEANUP_INTERVAL", 1*time.Hour),
		// Basic Authentication
		EnableBasicAuth:     getEnvBool("ENABLE_BASIC_AUTH", false),
		BasicAuthUsername:   getEnv("BASIC_AUTH_USERNAME", "admin"),
		BasicAuthPassword:   getEnv("BASIC_AUTH_PASSWORD", "change-me"),
		// JWT Authentication
		JWTSecret:              getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTExpiration:          getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
		RefreshTokenExpiration: getEnvDuration("REFRESH_TOKEN_EXPIRATION", 7*24*time.Hour), // 7日間
		// GitHub OAuth
		GitHubClientID:       getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret:   getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubCallbackURL:    getEnv("GITHUB_CALLBACK_URL", "http://localhost:8080/api/auth/github/callback"),
		OAuthSuccessRedirect: getEnv("OAUTH_SUCCESS_REDIRECT", "http://localhost:3000/auth/callback"),
		OAuthFailureRedirect: getEnv("OAUTH_FAILURE_REDIRECT", "http://localhost:3000/login?error=oauth_failed"),
		// Cookie Security
		CookieSecure:         getEnvBool("COOKIE_SECURE", false),
		// WebAuthn Settings
		WebAuthnRPID:         getEnv("WEBAUTHN_RP_ID", "localhost"),
		WebAuthnRPName:       getEnv("WEBAUTHN_RP_NAME", "財務計画計算機"),
		WebAuthnRPOrigin:     getEnv("WEBAUTHN_RP_ORIGIN", "http://localhost:3000"),
		// CSP: バックエンドはAPIサーバーのためHTMLを返さない厳格なポリシーをデフォルトとする
		// 本番環境では CONTENT_SECURITY_POLICY 環境変数で上書き可能
		// 開発環境では ENABLE_SECURE_HEADERS=false でヘッダー自体を無効化する
		ContentSecurityPolicy: getEnv("CONTENT_SECURITY_POLICY", "default-src 'none'; frame-ancestors 'none'; form-action 'none'"),
		// SMTP メール設定（デフォルトはResend）
		SMTPHost:     getEnv("SMTP_HOST", "smtp.resend.com"),
		SMTPPort:     getEnvInt("SMTP_PORT", 587),
		SMTPUser:     getEnv("SMTP_USER", "resend"),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@example.com"),
		// フロントエンドURL
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		// Bot LLM設定
		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
		GroqModel:  getEnv("GROQ_MODEL", "llama3-8b-8192"),
		FAQDir:     getEnv("FAQ_DIR", "docs/faq"),
	}

	return config
}

// Helper functions for environment variable parsing

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
