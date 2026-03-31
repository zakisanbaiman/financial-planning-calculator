package web

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/financial-planning-calculator/backend/application"
	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/services"
	infraemail "github.com/financial-planning-calculator/backend/infrastructure/email"
	"github.com/financial-planning-calculator/backend/infrastructure/faq"
	"github.com/financial-planning-calculator/backend/infrastructure/llm"
	infrapdf "github.com/financial-planning-calculator/backend/infrastructure/pdf"
	"github.com/financial-planning-calculator/backend/infrastructure/storage"
	"github.com/financial-planning-calculator/backend/infrastructure/web/controllers"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

// ServerDependencies holds all dependencies needed for the web server
type ServerDependencies struct {
	// Repositories
	UserRepo               repositories.UserRepository
	PasswordResetTokenRepo repositories.PasswordResetTokenRepository
	// Email service
	EmailService           infraemail.EmailService
	RefreshTokenRepo       repositories.RefreshTokenRepository
	WebAuthnCredentialRepo repositories.WebAuthnCredentialRepository
	FinancialPlanRepo      repositories.FinancialPlanRepository
	GoalRepo               repositories.GoalRepository

	// Domain Services
	CalculationService    *services.FinancialCalculationService
	RecommendationService *services.GoalRecommendationService

	// Auth Config
	JWTSecret              string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration

	// Server Config (for OAuth)
	ServerConfig *config.ServerConfig

	// WebAuthn
	WebAuthn *webauthn.WebAuthn

	// AuthUseCase (ミドルウェア用、NewControllersで初期化される)
	AuthUseCase usecases.AuthUseCase

	// SkipAuth テスト用：認証をスキップする
	SkipAuth bool
}

// NewControllers creates all controller instances with their dependencies
func NewControllers(deps *ServerDependencies) (*Controllers, error) {
	// Create use cases
	authUseCase := usecases.NewAuthUseCase(
		deps.UserRepo,
		deps.RefreshTokenRepo,
		deps.PasswordResetTokenRepo,
		deps.EmailService,
		deps.JWTSecret,
		deps.JWTExpiration,
		deps.RefreshTokenExpiration,
	)

	// Store auth use case for middleware
	deps.AuthUseCase = authUseCase

	manageFinancialDataUseCase := usecases.NewManageFinancialDataUseCase(
		deps.FinancialPlanRepo,
	)

	manageGoalsUseCase := usecases.NewManageGoalsUseCase(
		deps.GoalRepo,
		deps.FinancialPlanRepo,
		deps.RecommendationService,
	)

	calculateProjectionUseCase := usecases.NewCalculateProjectionUseCase(
		deps.FinancialPlanRepo,
		deps.GoalRepo,
		deps.CalculationService,
		deps.RecommendationService,
	)

	// TemporaryFileStorage を生成
	tempFileStorage, err := storage.NewTemporaryFileStorage(
		deps.ServerConfig.TempFileDir,
		deps.ServerConfig.TempFileSecret,
		deps.ServerConfig.TempFileExpiry,
		deps.ServerConfig.CleanupInterval,
	)
	if err != nil {
		return nil, fmt.Errorf("TemporaryFileStorageの初期化に失敗しました: %w", err)
	}

	// HTMLGenerator を初期化して ReportPDFGenerator アダプターでラップする
	pdfGenerator := infrapdf.NewHTMLGeneratorAdapter()

	generateReportsUseCase := usecases.NewGenerateReportsUseCaseWithPDF(
		deps.FinancialPlanRepo,
		deps.GoalRepo,
		deps.CalculationService,
		deps.RecommendationService,
		pdfGenerator,
		tempFileStorage,
	)

	// WebAuthn use case
	var webAuthnUseCase usecases.WebAuthnUseCase
	if deps.WebAuthn != nil && deps.WebAuthnCredentialRepo != nil {
		webAuthnUseCase = usecases.NewWebAuthnUseCase(
			deps.UserRepo,
			deps.WebAuthnCredentialRepo,
			deps.RefreshTokenRepo,
			deps.WebAuthn,
			authUseCase,
			deps.JWTSecret,
			deps.JWTExpiration,
			deps.RefreshTokenExpiration,
		)
	}

	// BotController初期化
	faqLoader := faq.NewFAQLoader(deps.ServerConfig.FAQDir)
	if _, loadErr := faqLoader.Load(context.Background()); loadErr != nil {
		slog.Error("FAQの読み込みに失敗しました。BotはFAQなしで動作します。", slog.Any("error", loadErr))
	}
	var llmClient llm.LLMClient
	if deps.ServerConfig.LocalLLMModel != "" {
		llmClient = llm.NewLocalLLMClientWithModel(deps.ServerConfig.LocalLLMBaseURL, deps.ServerConfig.LocalLLMModel)
	} else {
		llmClient = llm.NewLocalLLMClient(deps.ServerConfig.LocalLLMBaseURL)
	}
	botUseCase := application.NewBotUseCase(faqLoader, llmClient)

	// Create controllers
	return &Controllers{
		Auth:          controllers.NewAuthController(authUseCase, deps.ServerConfig),
		TwoFactor:     controllers.NewTwoFactorController(authUseCase, deps.ServerConfig),
		WebAuthn:      controllers.NewWebAuthnController(webAuthnUseCase),
		FinancialData: controllers.NewFinancialDataController(manageFinancialDataUseCase),
		Calculations:  controllers.NewCalculationsController(calculateProjectionUseCase),
		Goals:         controllers.NewGoalsController(manageGoalsUseCase),
		Reports:       controllers.NewReportsController(generateReportsUseCase, tempFileStorage),
		Bot:           controllers.NewBotController(botUseCase),
	}, nil
}

// JWTAuthMiddlewareFunc returns the JWT authentication middleware
// Returns nil if SkipAuth is true (for testing)
func (deps *ServerDependencies) JWTAuthMiddlewareFunc() echo.MiddlewareFunc {
	if deps.SkipAuth {
		return nil
	}
	return JWTAuthMiddleware(deps.AuthUseCase)
}
