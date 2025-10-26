package controllers

import (
	"fmt"
	"net/http"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/labstack/echo/v4"
)

// GoalsController は目標管理のコントローラー
type GoalsController struct {
	useCase usecases.ManageGoalsUseCase
}

// NewGoalsController は新しいGoalsControllerを作成する
func NewGoalsController(useCase usecases.ManageGoalsUseCase) *GoalsController {
	return &GoalsController{
		useCase: useCase,
	}
}

// CreateGoalRequest は目標作成リクエスト
type CreateGoalRequest struct {
	UserID              string  `json:"user_id" validate:"required"`
	GoalType            string  `json:"goal_type" validate:"required,oneof=savings retirement emergency custom"`
	Title               string  `json:"title" validate:"required,min=1,max=100"`
	TargetAmount        float64 `json:"target_amount" validate:"required,gt=0"`
	TargetDate          string  `json:"target_date" validate:"required"` // RFC3339 format
	CurrentAmount       float64 `json:"current_amount" validate:"gte=0"`
	MonthlyContribution float64 `json:"monthly_contribution" validate:"gte=0"`
	Description         *string `json:"description,omitempty"`
}

// UpdateGoalRequest は目標更新リクエスト
type UpdateGoalRequest struct {
	Title               *string  `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	TargetAmount        *float64 `json:"target_amount,omitempty" validate:"omitempty,gt=0"`
	TargetDate          *string  `json:"target_date,omitempty"` // RFC3339 format
	MonthlyContribution *float64 `json:"monthly_contribution,omitempty" validate:"omitempty,gte=0"`
	Description         *string  `json:"description,omitempty"`
	IsActive            *bool    `json:"is_active,omitempty"`
}

// UpdateGoalProgressRequest は目標進捗更新リクエスト
type UpdateGoalProgressRequest struct {
	CurrentAmount float64 `json:"current_amount" validate:"required,gte=0"`
	Note          *string `json:"note,omitempty"`
}

// GetGoalsQueryParams は目標一覧取得のクエリパラメータ
type GetGoalsQueryParams struct {
	UserID     string `query:"user_id" validate:"required"`
	GoalType   string `query:"goal_type,omitempty"`
	ActiveOnly bool   `query:"active_only,omitempty"`
}

// CreateGoal は新しい目標を作成する
// @Summary 目標作成
// @Description 新しい財務目標を作成します
// @Tags goals
// @Accept json
// @Produce json
// @Param request body CreateGoalRequest true "目標作成リクエスト"
// @Success 201 {object} usecases.CreateGoalOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /goals [post]
func (c *GoalsController) CreateGoal(ctx echo.Context) error {
	var req CreateGoalRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// Business logic validation for goals
	if err := ValidateBusinessLogic(ctx,
		func() *BusinessLogicError {
			// 要件6.4: 目標金額の妥当性チェック
			if req.TargetAmount <= 0 {
				return CreateBusinessLogicError(
					"INVALID_TARGET_AMOUNT",
					"目標金額は0より大きい値を入力してください",
					"達成したい具体的な金額を入力してください",
					req.TargetAmount,
					"正の数値",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 現在の金額の妥当性チェック
			if req.CurrentAmount < 0 {
				return CreateBusinessLogicError(
					"INVALID_CURRENT_AMOUNT",
					"現在の金額は0以上の値を入力してください",
					"現在の達成状況を正確に入力してください",
					req.CurrentAmount,
					"0以上の数値",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 月間積立額の妥当性チェック
			if req.MonthlyContribution < 0 {
				return CreateBusinessLogicError(
					"INVALID_MONTHLY_CONTRIBUTION",
					"月間積立額は0以上の値を入力してください",
					"毎月積み立て可能な金額を入力してください",
					req.MonthlyContribution,
					"0以上の数値",
				)
			}
			return nil
		},
		func() *BusinessLogicError {
			// 目標達成の実現可能性チェック
			if req.CurrentAmount > req.TargetAmount {
				return CreateBusinessLogicError(
					"GOAL_ALREADY_ACHIEVED",
					"現在の金額が目標金額を上回っています",
					"目標金額を見直すか、新しい目標を設定してください",
					req.CurrentAmount,
					fmt.Sprintf("%.0f以下", req.TargetAmount),
				)
			}
			return nil
		},
	); err != nil {
		return err
	}

	input := usecases.CreateGoalInput{
		UserID:              entities.UserID(req.UserID),
		GoalType:            req.GoalType,
		Title:               req.Title,
		TargetAmount:        req.TargetAmount,
		TargetDate:          req.TargetDate,
		CurrentAmount:       req.CurrentAmount,
		MonthlyContribution: req.MonthlyContribution,
		Description:         req.Description,
	}

	output, err := c.useCase.CreateGoal(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.JSON(http.StatusCreated, output)
}

// GetGoals は目標一覧を取得する
// @Summary 目標一覧取得
// @Description ユーザーの目標一覧を取得します
// @Tags goals
// @Produce json
// @Param user_id query string true "ユーザーID"
// @Param goal_type query string false "目標タイプ"
// @Param active_only query bool false "アクティブな目標のみ"
// @Success 200 {object} usecases.GetGoalsByUserOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /goals [get]
func (c *GoalsController) GetGoals(ctx echo.Context) error {
	var params GetGoalsQueryParams
	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "クエリパラメータの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&params); err != nil {
		return err // Validator already returns proper error response
	}

	input := usecases.GetGoalsByUserInput{
		UserID:     entities.UserID(params.UserID),
		ActiveOnly: params.ActiveOnly,
	}

	// 目標タイプが指定されている場合は設定
	if params.GoalType != "" {
		goalType := entities.GoalType(params.GoalType)
		if goalType.IsValid() {
			input.GoalType = &goalType
		} else {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "無効な目標タイプです", map[string]string{
				"valid_types": "savings, retirement, emergency, custom",
			}))
		}
	}

	output, err := c.useCase.GetGoalsByUser(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.JSON(http.StatusOK, output)
}

// GetGoal は特定の目標を取得する
// @Summary 目標取得
// @Description 特定の目標を取得します
// @Tags goals
// @Produce json
// @Param id path string true "目標ID"
// @Param user_id query string true "ユーザーID"
// @Success 200 {object} usecases.GetGoalOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /goals/{id} [get]
func (c *GoalsController) GetGoal(ctx echo.Context) error {
	goalID := ctx.Param("id")
	if goalID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "目標IDは必須です",
		})
	}

	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ユーザーIDは必須です",
		})
	}

	input := usecases.GetGoalInput{
		GoalID: entities.GoalID(goalID),
		UserID: entities.UserID(userID),
	}

	output, err := c.useCase.GetGoal(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "目標が見つかりません",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// UpdateGoal は目標を更新する
// @Summary 目標更新
// @Description 目標を更新します
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "目標ID"
// @Param user_id query string true "ユーザーID"
// @Param request body UpdateGoalRequest true "目標更新リクエスト"
// @Success 200 {object} usecases.UpdateGoalOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /goals/{id} [put]
func (c *GoalsController) UpdateGoal(ctx echo.Context) error {
	goalID := ctx.Param("id")
	if goalID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "目標IDは必須です",
		})
	}

	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ユーザーIDは必須です",
		})
	}

	var req UpdateGoalRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "リクエストの解析に失敗しました",
			Details: err.Error(),
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "入力値が無効です",
			Details: err.Error(),
		})
	}

	input := usecases.UpdateGoalInput{
		GoalID:              entities.GoalID(goalID),
		UserID:              entities.UserID(userID),
		Title:               req.Title,
		TargetAmount:        req.TargetAmount,
		TargetDate:          req.TargetDate,
		MonthlyContribution: req.MonthlyContribution,
		Description:         req.Description,
		IsActive:            req.IsActive,
	}

	output, err := c.useCase.UpdateGoal(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "目標の更新に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// UpdateGoalProgress は目標の進捗を更新する
// @Summary 目標進捗更新
// @Description 目標の進捗を更新します
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "目標ID"
// @Param user_id query string true "ユーザーID"
// @Param request body UpdateGoalProgressRequest true "目標進捗更新リクエスト"
// @Success 200 {object} usecases.UpdateGoalProgressOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /goals/{id}/progress [put]
func (c *GoalsController) UpdateGoalProgress(ctx echo.Context) error {
	goalID := ctx.Param("id")
	if goalID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "目標IDは必須です",
		})
	}

	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ユーザーIDは必須です",
		})
	}

	var req UpdateGoalProgressRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "リクエストの解析に失敗しました",
			Details: err.Error(),
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "入力値が無効です",
			Details: err.Error(),
		})
	}

	input := usecases.UpdateGoalProgressInput{
		GoalID:        entities.GoalID(goalID),
		UserID:        entities.UserID(userID),
		CurrentAmount: req.CurrentAmount,
		Note:          req.Note,
	}

	output, err := c.useCase.UpdateGoalProgress(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "目標進捗の更新に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// DeleteGoal は目標を削除する
// @Summary 目標削除
// @Description 目標を削除します
// @Tags goals
// @Param id path string true "目標ID"
// @Param user_id query string true "ユーザーID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /goals/{id} [delete]
func (c *GoalsController) DeleteGoal(ctx echo.Context) error {
	goalID := ctx.Param("id")
	if goalID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "目標IDは必須です",
		})
	}

	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ユーザーIDは必須です",
		})
	}

	input := usecases.DeleteGoalInput{
		GoalID: entities.GoalID(goalID),
		UserID: entities.UserID(userID),
	}

	err := c.useCase.DeleteGoal(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "目標の削除に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// GetGoalRecommendations は目標の推奨事項を取得する
// @Summary 目標推奨事項取得
// @Description 目標の推奨事項を取得します
// @Tags goals
// @Produce json
// @Param id path string true "目標ID"
// @Param user_id query string true "ユーザーID"
// @Success 200 {object} usecases.GetGoalRecommendationsOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /goals/{id}/recommendations [get]
func (c *GoalsController) GetGoalRecommendations(ctx echo.Context) error {
	goalID := ctx.Param("id")
	if goalID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "目標IDは必須です",
		})
	}

	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ユーザーIDは必須です",
		})
	}

	input := usecases.GetGoalRecommendationsInput{
		GoalID: entities.GoalID(goalID),
		UserID: entities.UserID(userID),
	}

	output, err := c.useCase.GetGoalRecommendations(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "目標推奨事項の取得に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// AnalyzeGoalFeasibility は目標の実現可能性を分析する
// @Summary 目標実現可能性分析
// @Description 目標の実現可能性を分析します
// @Tags goals
// @Produce json
// @Param id path string true "目標ID"
// @Param user_id query string true "ユーザーID"
// @Success 200 {object} usecases.AnalyzeGoalFeasibilityOutput
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /goals/{id}/feasibility [get]
func (c *GoalsController) AnalyzeGoalFeasibility(ctx echo.Context) error {
	goalID := ctx.Param("id")
	if goalID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "目標IDは必須です",
		})
	}

	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ユーザーIDは必須です",
		})
	}

	input := usecases.AnalyzeGoalFeasibilityInput{
		GoalID: entities.GoalID(goalID),
		UserID: entities.UserID(userID),
	}

	output, err := c.useCase.AnalyzeGoalFeasibility(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "目標実現可能性の分析に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}
