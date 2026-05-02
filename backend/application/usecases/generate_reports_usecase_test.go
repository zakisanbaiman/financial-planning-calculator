package usecases

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===========================
// Mock: ReportPDFGenerator
// ===========================

// mockReportPDFGenerator は ReportPDFGenerator インターフェースのモック
// 実装時に usecases パッケージ内で定義される ReportPDFGenerator インターフェースに対応する
type mockReportPDFGenerator struct {
	generateFunc func(reportType string, reportData interface{}) ([]byte, error)
}

func (m *mockReportPDFGenerator) Generate(reportType string, reportData interface{}) ([]byte, error) {
	if m.generateFunc != nil {
		return m.generateFunc(reportType, reportData)
	}
	return []byte("<html>dummy pdf content</html>"), nil
}

// ===========================
// Mock: TemporaryFileStoragePort
// ===========================

// mockTemporaryFileStoragePort は TemporaryFileStoragePort インターフェースのモック
// 実装時に usecases パッケージ内で定義される TemporaryFileStoragePort インターフェースに対応する
type mockTemporaryFileStoragePort struct {
	saveFileFunc func(fileName string, data []byte) (string, time.Time, error)
	getFileFunc  func(token string) ([]byte, string, string, error)
}

func (m *mockTemporaryFileStoragePort) SaveFile(fileName string, data []byte) (string, time.Time, error) {
	if m.saveFileFunc != nil {
		return m.saveFileFunc(fileName, data)
	}
	return "test-token-abc123", time.Now().Add(24 * time.Hour), nil
}

func (m *mockTemporaryFileStoragePort) GetFile(token string) ([]byte, string, string, error) {
	if m.getFileFunc != nil {
		return m.getFileFunc(token)
	}
	return nil, "", "", errors.New("not implemented")
}

// newTestFinancialPlanWithRetirementData は退職データ付きテスト用財務計画を作成するヘルパー
func newTestFinancialPlanWithRetirementData(userID entities.UserID) *aggregates.FinancialPlan {
	plan := newTestFinancialPlan(userID)
	monthlyExpenses, _ := valueobjects.NewMoneyJPY(200000)
	pension, _ := valueobjects.NewMoneyJPY(80000)
	retirement, _ := entities.NewRetirementData(userID, 40, 65, 85, monthlyExpenses, pension)
	_ = plan.SetRetirementData(retirement)
	return plan
}

// ===========================
// GenerateFinancialSummaryReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateFinancialSummaryReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 財務サマリーレポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateFinancialSummaryReport(ctx, FinancialSummaryReportInput{
			UserID: "user-001",
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.GeneratedAt)
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: 財務計画が存在しない場合はエラー", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateFinancialSummaryReport(ctx, FinancialSummaryReportInput{
			UserID: "user-999",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "財務計画の取得に失敗しました")
		mockPlanRepo.AssertExpectations(t)
	})
}

// ===========================
// GenerateAssetProjectionReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateAssetProjectionReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 資産推移レポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateAssetProjectionReport(ctx, AssetProjectionReportInput{
			UserID: "user-001",
			Years:  10,
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: FindByUserIDのエラーを伝播する", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("db error"))

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateAssetProjectionReport(ctx, AssetProjectionReportInput{
			UserID: "user-999",
			Years:  10,
		})

		require.Error(t, err)
		mockPlanRepo.AssertExpectations(t)
	})
}

// ===========================
// GenerateGoalsProgressReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateGoalsProgressReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 目標進捗レポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return([]*entities.Goal{goal}, nil)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateGoalsProgressReport(ctx, GoalsProgressReportInput{
			UserID: "user-001",
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		mockGoalRepo.AssertExpectations(t)
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: FindByUserIDのエラーを伝播する", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("db error"))

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateGoalsProgressReport(ctx, GoalsProgressReportInput{
			UserID: "user-999",
		})

		require.Error(t, err)
		mockGoalRepo.AssertExpectations(t)
	})
}
// ===========================
// GenerateRetirementPlanReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateRetirementPlanReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 退職計画レポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlanWithRetirementData("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateRetirementPlanReport(ctx, RetirementPlanReportInput{
			UserID: "user-001",
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.GeneratedAt)
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: 財務計画が存在しない場合はエラー", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateRetirementPlanReport(ctx, RetirementPlanReportInput{
			UserID: "user-999",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "財務計画の取得に失敗しました")
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: 退職データが設定されていない場合はエラー", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001") // 退職データなし
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateRetirementPlanReport(ctx, RetirementPlanReportInput{
			UserID: "user-001",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "退職データが設定されていません")
		mockPlanRepo.AssertExpectations(t)
	})
}

// ===========================
// GenerateComprehensiveReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateComprehensiveReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 包括的レポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(nil, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateComprehensiveReport(ctx, ComprehensiveReportInput{
			UserID: "user-001",
			Years:  10,
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.GeneratedAt)
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: 財務計画が存在しない場合はエラー", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateComprehensiveReport(ctx, ComprehensiveReportInput{
			UserID: "user-999",
			Years:  10,
		})

		require.Error(t, err)
		mockPlanRepo.AssertExpectations(t)
	})
}

// ===========================
// ExportReportToPDF Tests
// ===========================

func TestGenerateReportsUseCase_ExportReportToPDF(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: PDF生成・保存が成功してトークンが返る", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)

		plan := newTestFinancialPlan(entities.UserID("user-001"))
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		pdfContent := []byte("<html>financial summary pdf</html>")
		expectedToken := "test-download-token-xyz"

		pdfGen := &mockReportPDFGenerator{
			generateFunc: func(reportType string, reportData interface{}) ([]byte, error) {
				assert.Equal(t, "financial_summary", reportType)
				return pdfContent, nil
			},
		}
		fileStorage := &mockTemporaryFileStoragePort{
			saveFileFunc: func(fileName string, data []byte) (string, time.Time, error) {
				// ファイル名にユーザーIDプレフィックスが含まれることを検証
				assert.True(t, strings.HasPrefix(fileName, "user-001_"))
				assert.Equal(t, pdfContent, data)
				return expectedToken, time.Now().Add(24 * time.Hour), nil
			},
		}

		// 新シグネチャ: NewGenerateReportsUseCaseWithPDF(planRepo, goalRepo, calcService, recService, pdfGen, fileStorage)
		uc := NewGenerateReportsUseCaseWithPDF(mockPlanRepo, mockGoalRepo, calcService, recService, pdfGen, fileStorage)
		output, err := uc.ExportReportToPDF(ctx, ExportReportInput{
			UserID:     "user-001",
			ReportType: "financial_summary",
			Format:     "pdf",
			ReportData: map[string]interface{}{"key": "value"},
		})

		require.NoError(t, err)
		require.NotNil(t, output)
		// ユースケースはDownloadTokenのみを返す（DownloadURLはControllerが構築する）
		assert.NotEmpty(t, output.DownloadToken)
		assert.Equal(t, expectedToken, output.DownloadToken)
		assert.Empty(t, output.DownloadURL)
		assert.NotEmpty(t, output.ExpiresAt)
		assert.Greater(t, output.FileSize, int64(0))
	})

	t.Run("異常系: PDF生成失敗時にエラーが返る", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)

		plan := newTestFinancialPlan(entities.UserID("user-001"))
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		pdfGen := &mockReportPDFGenerator{
			generateFunc: func(reportType string, reportData interface{}) ([]byte, error) {
				return nil, errors.New("PDF生成エンジンエラー")
			},
		}
		fileStorage := &mockTemporaryFileStoragePort{}

		uc := NewGenerateReportsUseCaseWithPDF(mockPlanRepo, mockGoalRepo, calcService, recService, pdfGen, fileStorage)
		_, err := uc.ExportReportToPDF(ctx, ExportReportInput{
			UserID:     "user-001",
			ReportType: "financial_summary",
			Format:     "pdf",
			ReportData: map[string]interface{}{"key": "value"},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "PDF")
	})

	t.Run("異常系: ストレージ保存失敗時にエラーが返る", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)

		plan := newTestFinancialPlan(entities.UserID("user-001"))
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		pdfGen := &mockReportPDFGenerator{
			generateFunc: func(reportType string, reportData interface{}) ([]byte, error) {
				return []byte("<html>pdf</html>"), nil
			},
		}
		fileStorage := &mockTemporaryFileStoragePort{
			saveFileFunc: func(fileName string, data []byte) (string, time.Time, error) {
				return "", time.Time{}, errors.New("ディスク容量不足")
			},
		}

		uc := NewGenerateReportsUseCaseWithPDF(mockPlanRepo, mockGoalRepo, calcService, recService, pdfGen, fileStorage)
		_, err := uc.ExportReportToPDF(ctx, ExportReportInput{
			UserID:     "user-001",
			ReportType: "financial_summary",
			Format:     "pdf",
			ReportData: map[string]interface{}{"key": "value"},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "保存")
	})

	t.Run("異常系: pdfGeneratorがnilの場合はエラーが返る", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)

		// pdfGeneratorなしの場合
		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.ExportReportToPDF(ctx, ExportReportInput{
			UserID:     "user-001",
			ReportType: "financial_summary",
			Format:     "pdf",
			ReportData: map[string]interface{}{"key": "value"},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "PDFジェネレーター")
	})
}