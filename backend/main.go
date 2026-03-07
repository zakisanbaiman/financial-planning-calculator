package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/infrastructure/monitoring"
	"github.com/financial-planning-calculator/backend/infrastructure/email"
	"github.com/financial-planning-calculator/backend/infrastructure/repositories"
	"github.com/financial-planning-calculator/backend/infrastructure/web"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"

	_ "github.com/financial-planning-calculator/backend/docs"
)

// @title 財務計画計算機 API
// @version 1.0
// @description 将来の資産形成と老後の財務計画を可視化するアプリケーションのAPI
// @host localhost:8080
// @BasePath /api
func main() {
	// 設定読み込み
	cfg := config.LoadServerConfig()
	dbConfig := config.NewDatabaseConfig()

	// セキュリティ警告チェック
	checkSecurityWarnings(cfg, dbConfig)

	// 監視システムの初期化
	initMonitoring(cfg)

	// Echo インスタンス作成
	e := echo.New()

	// サーバー設定
	e.HideBanner = true
	e.Debug = cfg.Debug

	// バリデーター設定
	e.Validator = web.NewCustomValidator()

	// カスタムエラーハンドラー
	e.HTTPErrorHandler = web.CustomHTTPErrorHandler

	// ミドルウェア設定
	rateLimitStore := web.SetupMiddleware(e, cfg)

	// 依存関係の初期化
	deps := initializeDependencies()

	// コントローラーの作成
	controllers := web.NewControllers(deps)

	// ルーティング設定
	web.SetupRoutes(e, controllers, deps, rateLimitStore)

	// pprofサーバーの起動（開発環境のみ）
	if cfg.EnablePprof {
		go func() {
			log.Printf("🔍 pprof サーバーを起動: http://localhost:%s/debug/pprof/", cfg.PprofPort)
			log.Printf("   - CPU プロファイル: http://localhost:%s/debug/pprof/profile", cfg.PprofPort)
			log.Printf("   - メモリプロファイル: http://localhost:%s/debug/pprof/heap", cfg.PprofPort)
			log.Printf("   - ゴルーチン: http://localhost:%s/debug/pprof/goroutine", cfg.PprofPort)
			if err := http.ListenAndServe(":"+cfg.PprofPort, nil); err != nil {
				log.Printf("⚠️  pprof サーバーエラー: %v", err)
			}
		}()
	}

	// サーバー起動
	log.Printf("サーバーを開始します: http://localhost:%s", cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)
	log.Printf("API Base URL: http://localhost:%s/api", cfg.Port)
	log.Printf("Debug モード: %v", cfg.Debug)
	log.Printf("許可されたオリジン: %v", cfg.AllowedOrigins)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

// initMonitoring は監視システムを初期化します
func initMonitoring(cfg *config.ServerConfig) {
	// Prometheusメトリクスの初期化
	monitoring.InitPrometheus()
	log.Println("✅ Prometheusメトリクスを初期化しました")

	// エラートラッキングの初期化
	environment := "development"
	if !cfg.Debug {
		environment = "production"
	}
	monitoring.InitErrorTracker(environment)
	log.Printf("✅ エラートラッキングを初期化しました (環境: %s)", environment)
}

// initializeDependencies initializes all dependencies for the application
func initializeDependencies() *web.ServerDependencies {
	// Initialize database connection
	dbConfig := config.NewDatabaseConfig()
	db, err := config.NewDatabaseConnection(dbConfig)
	if err != nil {
		log.Fatalf("データベース接続の初期化に失敗しました: %v", err)
	}

	// Initialize repositories
	repoFactory := repositories.NewRepositoryFactory(db)

	userRepo := repoFactory.NewUserRepository()
	refreshTokenRepo := repoFactory.NewRefreshTokenRepository()
	passwordResetTokenRepo := repoFactory.NewPasswordResetTokenRepository()
	webAuthnCredentialRepo := repoFactory.NewWebAuthnCredentialRepository()
	financialPlanRepo := repoFactory.NewFinancialPlanRepository()
	goalRepo := repoFactory.NewGoalRepository()

	// Initialize domain services
	calculationService := services.NewFinancialCalculationService()
	recommendationService := services.NewGoalRecommendationService(calculationService)

	// Load server config for JWT settings
	serverCfg := config.LoadServerConfig()

	// Initialize email service
	emailService := email.NewEmailService(
		serverCfg.SMTPHost,
		serverCfg.SMTPPort,
		serverCfg.SMTPUser,
		serverCfg.SMTPPassword,
		serverCfg.SMTPFrom,
	)

	// Initialize WebAuthn
	webAuthn, err := initializeWebAuthn(serverCfg)
	if err != nil {
		log.Printf("⚠️  WebAuthn初期化に失敗しました（パスキー機能は無効）: %v", err)
	}

	return &web.ServerDependencies{
		UserRepo:                 userRepo,
		RefreshTokenRepo:         refreshTokenRepo,
		PasswordResetTokenRepo:   passwordResetTokenRepo,
		EmailService:             emailService,
		WebAuthnCredentialRepo:   webAuthnCredentialRepo,
		FinancialPlanRepo:        financialPlanRepo,
		GoalRepo:                 goalRepo,
		CalculationService:       calculationService,
		RecommendationService:    recommendationService,
		JWTSecret:                serverCfg.JWTSecret,
		JWTExpiration:            serverCfg.JWTExpiration,
		RefreshTokenExpiration:   serverCfg.RefreshTokenExpiration,
		ServerConfig:             serverCfg, // OAuth設定用 (Issue: #67)
		WebAuthn:                 webAuthn,
	}
}

// checkSecurityWarnings checks for insecure default values in production
func checkSecurityWarnings(serverCfg *config.ServerConfig, dbCfg *config.DatabaseConfig) {
	warnings := []string{}

	// Check database password
	if dbCfg.Password == "password" {
		warnings = append(warnings, "⚠️  DB_PASSWORD is set to default value 'password'. Change it in production!")
	}

	// Check temporary file secret
	if serverCfg.TempFileSecret == "change-this-secret-in-production" {
		warnings = append(warnings, "⚠️  TEMP_FILE_SECRET is set to default value. Change it in production!")
	}

	// Check JWT secret
	if serverCfg.JWTSecret == "change-this-secret-in-production" {
		warnings = append(warnings, "⚠️  JWT_SECRET is set to default value. Change it in production!")
	}

	// Check SSL mode
	if dbCfg.SSLMode == "disable" {
		warnings = append(warnings, "⚠️  DB_SSLMODE is set to 'disable'. Enable SSL in production!")
	}

	// Output warnings
	if len(warnings) > 0 {
		log.Println("==================== SECURITY WARNINGS ====================")
		for _, warning := range warnings {
			log.Println(warning)
		}
		log.Println("===========================================================")
	}
}

// initializeWebAuthn initializes WebAuthn configuration
func initializeWebAuthn(cfg *config.ServerConfig) (*webauthn.WebAuthn, error) {
	wconfig := &webauthn.Config{
		RPDisplayName: cfg.WebAuthnRPName,
		RPID:          cfg.WebAuthnRPID,
		RPOrigins:     []string{cfg.WebAuthnRPOrigin},
	}

	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ WebAuthn初期化成功")
	log.Printf("   - RP Name: %s", cfg.WebAuthnRPName)
	log.Printf("   - RP ID: %s", cfg.WebAuthnRPID)
	log.Printf("   - RP Origin: %s", cfg.WebAuthnRPOrigin)

	return webAuthn, nil
}
