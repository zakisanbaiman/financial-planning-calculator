package main

import (
	"database/sql"
	"log"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/infrastructure/repositories"
	"github.com/financial-planning-calculator/backend/infrastructure/web"
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
	web.SetupMiddleware(e, cfg)

	// 依存関係の初期化
	deps := initializeDependencies()

	// コントローラーの作成
	controllers := web.NewControllers(deps)

	// ルーティング設定
	web.SetupRoutes(e, controllers)

	// サーバー起動
	log.Printf("サーバーを開始します: http://localhost:%s", cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)
	log.Printf("API Base URL: http://localhost:%s/api", cfg.Port)
	log.Printf("Debug モード: %v", cfg.Debug)
	log.Printf("許可されたオリジン: %v", cfg.AllowedOrigins)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

// initializeDependencies initializes all dependencies for the application
func initializeDependencies() *web.ServerDependencies {
	// Initialize repositories
	// Note: For now using nil database connection - this should be replaced with actual DB connection
	// The repository factory will handle the database connection
	var db *sql.DB // This should be initialized with actual database connection
	repoFactory := repositories.NewRepositoryFactory(db)

	financialPlanRepo := repoFactory.NewFinancialPlanRepository()
	goalRepo := repoFactory.NewGoalRepository()

	// Initialize domain services
	calculationService := services.NewFinancialCalculationService()
	recommendationService := services.NewGoalRecommendationService(calculationService)

	return &web.ServerDependencies{
		FinancialPlanRepo:     financialPlanRepo,
		GoalRepo:              goalRepo,
		CalculationService:    calculationService,
		RecommendationService: recommendationService,
	}
}
