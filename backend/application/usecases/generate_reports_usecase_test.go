package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestGenerateReportsUseCase(
	planRepo *MockFinancialPlanRepository,
	goalRepo *MockGoalRepository,
) GenerateReportsUseCase {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)
	return NewGenerateReportsUseCase(planRepo, goalRepo, calcService, recommendService)
}

func TestGenerateFinancialSummaryReport(t *testing.T) {
	tests := []struct {
		name        string
		input       FinancialSummaryReportInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 財務サマリーレポート生成",
			input: FinancialSummaryReportInput{
				UserID: "user-001",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-001")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: FinancialSummaryReportInput{
				UserID: "user-999",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestGenerateReportsUseCase(planRepo, goalRepo)
			output, err := uc.GenerateFinancialSummaryReport(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.GeneratedAt)
				assert.Equal(t, tt.input.UserID, output.Report.UserID)
			}
			planRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateAssetProjectionReport(t *testing.T) {
	tests := []struct {
		name        string
		input       AssetProjectionReportInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 資産推移レポート生成",
			input: AssetProjectionReportInput{
				UserID: "user-001",
				Years:  10,
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-001")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: AssetProjectionReportInput{
				UserID: "user-999",
				Years:  10,
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestGenerateReportsUseCase(planRepo, goalRepo)
			output, err := uc.GenerateAssetProjectionReport(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.GeneratedAt)
				assert.Equal(t, tt.input.Years, output.Report.ProjectionYears)
			}
			planRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateGoalsProgressReport(t *testing.T) {
	tests := []struct {
		name        string
		input       GoalsProgressReportInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 目標進捗レポート生成（目標なし）",
			input: GoalsProgressReportInput{
				UserID: "user-001",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-001")
				gr.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return([]*entities.Goal{}, nil)
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "正常系: 目標進捗レポート生成（目標あり）",
			input: GoalsProgressReportInput{
				UserID: "user-002",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				goals := []*entities.Goal{
					newTestGoal("user-002", entities.GoalTypeSavings),
					newTestGoal("user-002", entities.GoalTypeCustom),
				}
				plan := newTestFinancialPlan("user-002")
				gr.On("FindByUserID", mock.Anything, entities.UserID("user-002")).Return(goals, nil)
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-002")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: リポジトリエラー",
			input: GoalsProgressReportInput{
				UserID: "user-999",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				gr.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("DBエラー"))
			},
			expectError: true,
			errContains: "目標の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestGenerateReportsUseCase(planRepo, goalRepo)
			output, err := uc.GenerateGoalsProgressReport(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.GeneratedAt)
			}
			goalRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateRetirementPlanReport(t *testing.T) {
	tests := []struct {
		name        string
		input       RetirementPlanReportInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 退職計画レポート生成",
			input: RetirementPlanReportInput{
				UserID: "user-001",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlanWithRetirementData("user-001")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: RetirementPlanReportInput{
				UserID: "user-999",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
		{
			name: "異常系: 退職データが未設定",
			input: RetirementPlanReportInput{
				UserID: "user-002",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-002")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-002")).Return(plan, nil)
			},
			expectError: true,
			errContains: "退職データが設定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestGenerateReportsUseCase(planRepo, goalRepo)
			output, err := uc.GenerateRetirementPlanReport(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.GeneratedAt)
			}
			planRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateComprehensiveReport(t *testing.T) {
	tests := []struct {
		name        string
		input       ComprehensiveReportInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 包括的レポート生成",
			input: ComprehensiveReportInput{
				UserID: "user-001",
				Years:  5,
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-001")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
				gr.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return([]*entities.Goal{}, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: ComprehensiveReportInput{
				UserID: "user-999",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestGenerateReportsUseCase(planRepo, goalRepo)
			output, err := uc.GenerateComprehensiveReport(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.GeneratedAt)
			}
			planRepo.AssertExpectations(t)
			goalRepo.AssertExpectations(t)
		})
	}
}

func TestExportReportToPDF(t *testing.T) {
	tests := []struct {
		name           string
		input          ExportReportInput
		setupMocks     func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError    bool
		expectFileName string
	}{
		{
			name: "正常系: PDFエクスポート（財務サマリー）",
			input: ExportReportInput{
				UserID:     "user-001",
				ReportType: "financial_summary",
				Format:     "pdf",
			},
			setupMocks:     func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {},
			expectError:    false,
			expectFileName: "financial_summary",
		},
		{
			name: "正常系: PDFエクスポート（包括的レポート）",
			input: ExportReportInput{
				UserID:     "user-001",
				ReportType: "comprehensive",
				Format:     "pdf",
			},
			setupMocks:     func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {},
			expectError:    false,
			expectFileName: "comprehensive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestGenerateReportsUseCase(planRepo, goalRepo)
			output, err := uc.ExportReportToPDF(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.DownloadURL)
				assert.NotEmpty(t, output.FileName)
				assert.Contains(t, output.FileName, tt.expectFileName)
				assert.Positive(t, output.FileSize)
			}
			planRepo.AssertExpectations(t)
		})
	}
}

// newTestPlanWithProfile creates a plan with specific financial parameters for scenario testing
func newTestPlanForReport(userID string, income, expenses, savings, investReturn float64) *aggregates.FinancialPlan {
	inc, _ := valueobjects.NewMoneyJPY(income)
	expCol := entities.ExpenseCollection{
		{Category: "生活費", Amount: mustCreateMoneyUsecase(expenses)},
	}
	savCol := entities.SavingsCollection{}
	if savings > 0 {
		savCol = entities.SavingsCollection{
			{Type: "deposit", Amount: mustCreateMoneyUsecase(savings)},
		}
	}
	ir, _ := valueobjects.NewRate(investReturn)
	inflation, _ := valueobjects.NewRate(2.0)
	profile, err := entities.NewFinancialProfile(entities.UserID(userID), inc, expCol, savCol, ir, inflation)
	if err != nil {
		panic(err)
	}
	plan, err := aggregates.NewFinancialPlan(profile)
	if err != nil {
		panic(err)
	}
	return plan
}

// TestGenerateFinancialSummaryReport_Scenarios tests different financial health levels
func TestGenerateFinancialSummaryReport_Scenarios(t *testing.T) {

	tests := []struct {
		name     string
		userID   string
		income   float64
		expenses float64
		savings  float64
		irr      float64
	}{
		{
			name:     "高貯蓄率・高投資利回り（excellent）",
			userID:   "report-user-1",
			income:   600000,
			expenses: 250000, // 貯蓄率 ~58%
			savings:  3600000,
			irr:      6.0,
		},
		{
			name:     "中程度貯蓄率（good）",
			userID:   "report-user-2",
			income:   500000,
			expenses: 380000, // 貯蓄率 ~24%
			savings:  1200000,
			irr:      4.0,
		},
		{
			name:     "低貯蓄率・低投資利回り（fair/poor）",
			userID:   "report-user-3",
			income:   300000,
			expenses: 290000, // 貯蓄率 ~3%
			savings:  100000,
			irr:      1.5,
		},
		{
			name:     "支出=収入（poor）",
			userID:   "report-user-4",
			income:   300000,
			expenses: 300000, // 貯蓄率 0%
			savings:  0,
			irr:      0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			plan := newTestPlanForReport(tt.userID, tt.income, tt.expenses, tt.savings, tt.irr)
			planRepo.On("FindByUserID", mock.Anything, entities.UserID(tt.userID)).Return(plan, nil)

			uc := newTestGenerateReportsUseCase(planRepo, goalRepo)
			output, err := uc.GenerateFinancialSummaryReport(context.Background(), FinancialSummaryReportInput{
				UserID: entities.UserID(tt.userID),
			})

			require.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.GeneratedAt)
			planRepo.AssertExpectations(t)
		})
	}
}
