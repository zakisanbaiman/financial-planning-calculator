package controllers

import (
	"net/http"
	"strings"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/labstack/echo/v4"
)

// FinancialDataController は財務データ管理のコントローラー
type FinancialDataController struct {
	useCase usecases.ManageFinancialDataUseCase
}

// NewFinancialDataController は新しいFinancialDataControllerを作成する
func NewFinancialDataController(useCase usecases.ManageFinancialDataUseCase) *FinancialDataController {
	return &FinancialDataController{
		useCase: useCase,
	}
}

// CreateFinancialDataRequest は財務データ作成リクエスト
type CreateFinancialDataRequest struct {
	UserID                     string               `json:"user_id" validate:"required"`
	MonthlyIncome              float64              `json:"monthly_income" validate:"required,gt=0"`
	MonthlyExpenses            []ExpenseItemRequest `json:"monthly_expenses" validate:"required,dive"`
	CurrentSavings             []SavingsItemRequest `json:"current_savings" validate:"required,dive"`
	InvestmentReturn           float64              `json:"investment_return" validate:"required,gte=0,lte=100"`
	InflationRate              float64              `json:"inflation_rate" validate:"required,gte=0,lte=50"`
	RetirementAge              *int                 `json:"retirement_age,omitempty" validate:"omitempty,gte=50,lte=100"`
	MonthlyRetirementExpenses  *float64             `json:"monthly_retirement_expenses,omitempty" validate:"omitempty,gt=0"`
	PensionAmount              *float64             `json:"pension_amount,omitempty" validate:"omitempty,gte=0"`
	EmergencyFundTargetMonths  *int                 `json:"emergency_fund_target_months,omitempty" validate:"omitempty,gte=1,lte=24"`
	EmergencyFundCurrentAmount *float64             `json:"emergency_fund_current_amount,omitempty" validate:"omitempty,gte=0"`
}

// ExpenseItemRequest は支出項目リクエスト
type ExpenseItemRequest struct {
	Category    string  `json:"category" validate:"required,min=1"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description *string `json:"description,omitempty"`
}

// SavingsItemRequest は貯蓄項目リクエスト
type SavingsItemRequest struct {
	Type        string  `json:"type" validate:"required,oneof=deposit investment other"`
	Amount      float64 `json:"amount" validate:"required,gte=0"`
	Description *string `json:"description,omitempty"`
}

// UpdateFinancialProfileRequest は財務プロファイル更新リクエスト
type UpdateFinancialProfileRequest struct {
	MonthlyIncome    float64              `json:"monthly_income" validate:"required,gt=0"`
	MonthlyExpenses  []ExpenseItemRequest `json:"monthly_expenses" validate:"required,dive"`
	CurrentSavings   []SavingsItemRequest `json:"current_savings" validate:"required,dive"`
	InvestmentReturn float64              `json:"investment_return" validate:"required,gte=0,lte=100"`
	InflationRate    float64              `json:"inflation_rate" validate:"required,gte=0,lte=50"`
}

// UpdateRetirementDataRequest は退職データ更新リクエスト
type UpdateRetirementDataRequest struct {
	RetirementAge             int     `json:"retirement_age" validate:"required,gte=50,lte=100"`
	MonthlyRetirementExpenses float64 `json:"monthly_retirement_expenses" validate:"required,gt=0"`
	PensionAmount             float64 `json:"pension_amount" validate:"required,gte=0"`
}

// UpdateEmergencyFundRequest は緊急資金更新リクエスト
type UpdateEmergencyFundRequest struct {
	TargetMonths  int     `json:"target_months" validate:"required,gte=1,lte=24"`
	CurrentAmount float64 `json:"current_amount" validate:"required,gte=0"`
}

// CreateFinancialData は財務データを作成する
// @Summary 財務データ作成
// @Description 新しい財務計画を作成します
// @Tags financial-data
// @Accept json
// @Produce json
// @Param request body CreateFinancialDataRequest true "財務データ作成リクエスト"
// @Success 201 {object} usecases.CreateFinancialPlanOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /financial-data [post]
func (c *FinancialDataController) CreateFinancialData(ctx echo.Context) error {
	var req CreateFinancialDataRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// Business logic validation
	if err := ValidateBusinessLogic(ctx,
		func() *BusinessLogicError {
			// 要件1.4: 入力値が無効（負の値など）の場合のエラー
			if req.MonthlyIncome <= 0 {
				return CreateBusinessLogicError(
					"INVALID_MONTHLY_INCOME",
					"月収は0より大きい値を入力してください",
					"正の数値を入力してください",
					req.MonthlyIncome,
					"正の数値",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 投資利回りの妥当性チェック
			if req.InvestmentReturn < 0 || req.InvestmentReturn > 100 {
				return CreateBusinessLogicError(
					"INVALID_INVESTMENT_RETURN",
					"投資利回りは0%から100%の範囲で入力してください",
					"現実的な投資利回り（例：3-7%）を入力してください",
					req.InvestmentReturn,
					"0-100%",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// インフレ率の妥当性チェック
			if req.InflationRate < 0 || req.InflationRate > 50 {
				return CreateBusinessLogicError(
					"INVALID_INFLATION_RATE",
					"インフレ率は0%から50%の範囲で入力してください",
					"現実的なインフレ率（例：1-3%）を入力してください",
					req.InflationRate,
					"0-50%",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 支出項目の妥当性チェック
			totalExpenses := 0.0
			for _, expense := range req.MonthlyExpenses {
				if expense.Amount <= 0 {
					return CreateBusinessLogicError(
						"INVALID_EXPENSE_AMOUNT",
						"支出金額は0より大きい値を入力してください",
						"正の数値を入力してください",
						expense.Amount,
						"正の数値",
					)
				}
				totalExpenses += expense.Amount
			}

			// 要件2.4: 貯蓄額が月間支出を下回る場合の警告
			totalSavings := 0.0
			for _, saving := range req.CurrentSavings {
				totalSavings += saving.Amount
			}

			monthlySavings := req.MonthlyIncome - totalExpenses
			if monthlySavings < 0 {
				return CreateBusinessLogicError(
					"INSUFFICIENT_SAVINGS",
					"月間支出が月収を上回っています",
					"支出を見直すか、収入を増やすことを検討してください",
					monthlySavings,
					"正の数値",
				)
			}

			return nil
		},
	); err != nil {
		return err
	}

	// リクエストをユースケース入力に変換
	input := usecases.CreateFinancialPlanInput{
		UserID:                     entities.UserID(req.UserID),
		MonthlyIncome:              req.MonthlyIncome,
		MonthlyExpenses:            convertExpenseItems(req.MonthlyExpenses),
		CurrentSavings:             convertSavingsItems(req.CurrentSavings),
		InvestmentReturn:           req.InvestmentReturn,
		InflationRate:              req.InflationRate,
		RetirementAge:              req.RetirementAge,
		MonthlyRetirementExpenses:  req.MonthlyRetirementExpenses,
		PensionAmount:              req.PensionAmount,
		EmergencyFundTargetMonths:  req.EmergencyFundTargetMonths,
		EmergencyFundCurrentAmount: req.EmergencyFundCurrentAmount,
	}

	output, err := c.useCase.CreateFinancialPlan(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.JSON(http.StatusCreated, output)
}

// GetFinancialData は財務データを取得する
// @Summary 財務データ取得
// @Description ユーザーの財務計画を取得します
// @Tags financial-data
// @Produce json
// @Param user_id query string true "ユーザーID"
// @Success 200 {object} usecases.GetFinancialPlanOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /financial-data [get]
func (c *FinancialDataController) GetFinancialData(ctx echo.Context) error {
	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "ユーザーIDは必須です", nil))
	}

	input := usecases.GetFinancialPlanInput{
		UserID: entities.UserID(userID),
	}

	output, err := c.useCase.GetFinancialPlan(ctx.Request().Context(), input)
	if err != nil {
		// 404 for not found, 500 for other errors
		// Check for various forms of "financial data not found" error messages
		errMsg := err.Error()
		if strings.Contains(errMsg, "財務データが見つかりません") || 
		   strings.Contains(errMsg, "財務プロファイルの取得に失敗しました") {
			return ctx.JSON(http.StatusNotFound, NewNotFoundErrorResponse(ctx, "財務データ"))
		}
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.JSON(http.StatusOK, output)
}

// UpdateFinancialProfile は財務プロファイルを更新する
// @Summary 財務プロファイル更新
// @Description 財務プロファイルを更新します
// @Tags financial-data
// @Accept json
// @Produce json
// @Param user_id path string true "ユーザーID"
// @Param request body UpdateFinancialProfileRequest true "財務プロファイル更新リクエスト"
// @Success 200 {object} usecases.UpdateFinancialProfileOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /financial-data/{user_id}/profile [put]
func (c *FinancialDataController) UpdateFinancialProfile(ctx echo.Context) error {
	userID := ctx.Param("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "ユーザーIDは必須です", nil))
	}

	var req UpdateFinancialProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// Business logic validation
	if err := ValidateBusinessLogic(ctx,
		func() *BusinessLogicError {
			// 要件1.4: 入力値が無効（負の値など）の場合のエラー
			if req.MonthlyIncome <= 0 {
				return CreateBusinessLogicError(
					"INVALID_MONTHLY_INCOME",
					"月収は0より大きい値を入力してください",
					"正の数値を入力してください",
					req.MonthlyIncome,
					"正の数値",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 要件2.4: 貯蓄額が月間支出を下回る場合の警告
			totalExpenses := 0.0
			for _, expense := range req.MonthlyExpenses {
				totalExpenses += expense.Amount
			}

			monthlySavings := req.MonthlyIncome - totalExpenses
			if monthlySavings < 0 {
				return CreateBusinessLogicError(
					"INSUFFICIENT_SAVINGS",
					"月間支出が月収を上回っています",
					"支出を見直すか、収入を増やすことを検討してください",
					monthlySavings,
					"正の数値",
				)
			}

			return nil
		},
	); err != nil {
		return err
	}

	input := usecases.UpdateFinancialProfileInput{
		UserID:           entities.UserID(userID),
		MonthlyIncome:    req.MonthlyIncome,
		MonthlyExpenses:  convertExpenseItems(req.MonthlyExpenses),
		CurrentSavings:   convertSavingsItems(req.CurrentSavings),
		InvestmentReturn: req.InvestmentReturn,
		InflationRate:    req.InflationRate,
	}

	output, err := c.useCase.UpdateFinancialProfile(ctx.Request().Context(), input)
	if err != nil {
		// If underlying error indicates missing financial data, return 404
		if strings.Contains(err.Error(), "財務データが見つかりません") || strings.Contains(err.Error(), "財務計画の取得に失敗しました") || strings.Contains(err.Error(), "財務プロファイルの取得に失敗しました") {
			return ctx.JSON(http.StatusNotFound, NewNotFoundErrorResponse(ctx, "財務データ"))
		}
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.JSON(http.StatusOK, output)
}

// UpdateRetirementData は退職データを更新する
// @Summary 退職データ更新
// @Description 退職データを更新します
// @Tags financial-data
// @Accept json
// @Produce json
// @Param user_id path string true "ユーザーID"
// @Param request body UpdateRetirementDataRequest true "退職データ更新リクエスト"
// @Success 200 {object} usecases.UpdateRetirementDataOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /financial-data/{user_id}/retirement [put]
func (c *FinancialDataController) UpdateRetirementData(ctx echo.Context) error {
	userID := ctx.Param("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "ユーザーIDは必須です", nil))
	}

	var req UpdateRetirementDataRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// Business logic validation for retirement data
	if err := ValidateBusinessLogic(ctx,
		func() *BusinessLogicError {
			// 要件3.4: 退職年齢の妥当性チェック
			if req.RetirementAge < 50 || req.RetirementAge > 100 {
				return CreateBusinessLogicError(
					"INVALID_RETIREMENT_AGE",
					"退職年齢は50歳から100歳の範囲で入力してください",
					"現実的な退職年齢を入力してください",
					req.RetirementAge,
					"50-100歳",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 老後生活費の妥当性チェック
			if req.MonthlyRetirementExpenses <= 0 {
				return CreateBusinessLogicError(
					"INVALID_RETIREMENT_EXPENSES",
					"老後月間生活費は0より大きい値を入力してください",
					"現実的な生活費を入力してください",
					req.MonthlyRetirementExpenses,
					"正の数値",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 年金受給額の妥当性チェック
			if req.PensionAmount < 0 {
				return CreateBusinessLogicError(
					"INVALID_PENSION_AMOUNT",
					"年金受給額は0以上の値を入力してください",
					"予想される年金受給額を入力してください",
					req.PensionAmount,
					"0以上の数値",
				)
			}
			return nil
		},
	); err != nil {
		return err
	}

	input := usecases.UpdateRetirementDataInput{
		UserID:                    entities.UserID(userID),
		RetirementAge:             req.RetirementAge,
		MonthlyRetirementExpenses: req.MonthlyRetirementExpenses,
		PensionAmount:             req.PensionAmount,
	}

	output, err := c.useCase.UpdateRetirementData(ctx.Request().Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "財務データが見つかりません") || strings.Contains(err.Error(), "財務計画の取得に失敗しました") || strings.Contains(err.Error(), "財務プロファイルの取得に失敗しました") {
			return ctx.JSON(http.StatusNotFound, NewNotFoundErrorResponse(ctx, "財務データ"))
		}
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.JSON(http.StatusOK, output)
}

// UpdateEmergencyFund は緊急資金設定を更新する
// @Summary 緊急資金設定更新
// @Description 緊急資金設定を更新します
// @Tags financial-data
// @Accept json
// @Produce json
// @Param user_id path string true "ユーザーID"
// @Param request body UpdateEmergencyFundRequest true "緊急資金設定更新リクエスト"
// @Success 200 {object} usecases.UpdateEmergencyFundOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /financial-data/{user_id}/emergency-fund [put]
func (c *FinancialDataController) UpdateEmergencyFund(ctx echo.Context) error {
	userID := ctx.Param("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "ユーザーIDは必須です", nil))
	}

	var req UpdateEmergencyFundRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// Business logic validation for emergency fund
	if err := ValidateBusinessLogic(ctx,
		func() *BusinessLogicError {
			// 要件4.4: 緊急資金目標月数の妥当性チェック
			if req.TargetMonths < 1 || req.TargetMonths > 24 {
				return CreateBusinessLogicError(
					"INVALID_TARGET_MONTHS",
					"緊急資金目標月数は1ヶ月から24ヶ月の範囲で入力してください",
					"一般的には3-6ヶ月分の生活費が推奨されます",
					req.TargetMonths,
					"1-24ヶ月",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 現在の緊急資金額の妥当性チェック
			if req.CurrentAmount < 0 {
				return CreateBusinessLogicError(
					"INVALID_CURRENT_AMOUNT",
					"現在の緊急資金額は0以上の値を入力してください",
					"現在保有している緊急資金の金額を入力してください",
					req.CurrentAmount,
					"0以上の数値",
				)
			}
			return nil
		},
	); err != nil {
		return err
	}

	input := usecases.UpdateEmergencyFundInput{
		UserID:        entities.UserID(userID),
		TargetMonths:  req.TargetMonths,
		CurrentAmount: req.CurrentAmount,
	}

	output, err := c.useCase.UpdateEmergencyFund(ctx.Request().Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "財務データが見つかりません") || strings.Contains(err.Error(), "財務計画の取得に失敗しました") || strings.Contains(err.Error(), "財務プロファイルの取得に失敗しました") {
			return ctx.JSON(http.StatusNotFound, NewNotFoundErrorResponse(ctx, "財務データ"))
		}
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.JSON(http.StatusOK, output)
}

// DeleteFinancialData は財務データを削除する
// @Summary 財務データ削除
// @Description 財務計画を削除します
// @Tags financial-data
// @Param user_id path string true "ユーザーID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /financial-data/{user_id} [delete]
func (c *FinancialDataController) DeleteFinancialData(ctx echo.Context) error {
	userID := ctx.Param("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "ユーザーIDは必須です", nil))
	}

	input := usecases.DeleteFinancialPlanInput{
		UserID: entities.UserID(userID),
	}

	err := c.useCase.DeleteFinancialPlan(ctx.Request().Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "財務データが見つかりません") || strings.Contains(err.Error(), "財務計画の取得に失敗しました") || strings.Contains(err.Error(), "財務プロファイルの取得に失敗しました") {
			return ctx.JSON(http.StatusNotFound, NewNotFoundErrorResponse(ctx, "財務データ"))
		}
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.NoContent(http.StatusNoContent)
}

// convertExpenseItems はExpenseItemRequestをusecases.ExpenseItemに変換する
func convertExpenseItems(items []ExpenseItemRequest) []usecases.ExpenseItem {
	result := make([]usecases.ExpenseItem, len(items))
	for i, item := range items {
		result[i] = usecases.ExpenseItem{
			Category:    item.Category,
			Amount:      item.Amount,
			Description: item.Description,
		}
	}
	return result
}

// convertSavingsItems はSavingsItemRequestをusecases.SavingsItemに変換する
func convertSavingsItems(items []SavingsItemRequest) []usecases.SavingsItem {
	result := make([]usecases.SavingsItem, len(items))
	for i, item := range items {
		result[i] = usecases.SavingsItem{
			Type:        item.Type,
			Amount:      item.Amount,
			Description: item.Description,
		}
	}
	return result
}
