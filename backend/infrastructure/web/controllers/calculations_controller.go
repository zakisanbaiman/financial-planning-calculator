package controllers

import (
	"net/http"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/labstack/echo/v4"
)

// CalculationsController は計算機能のコントローラー
type CalculationsController struct {
	useCase usecases.CalculateProjectionUseCase
}

// NewCalculationsController は新しいCalculationsControllerを作成する
func NewCalculationsController(useCase usecases.CalculateProjectionUseCase) *CalculationsController {
	return &CalculationsController{
		useCase: useCase,
	}
}

// AssetProjectionRequest は資産推移計算リクエスト
type AssetProjectionRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Years  int    `json:"years" validate:"required,gte=1,lte=50"`
}

// RetirementCalculationRequest は退職資金計算リクエスト
type RetirementCalculationRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// EmergencyFundCalculationRequest は緊急資金計算リクエスト
type EmergencyFundCalculationRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// ComprehensiveProjectionRequest は包括的予測計算リクエスト
type ComprehensiveProjectionRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Years  int    `json:"years" validate:"required,gte=1,lte=50"`
}

// GoalProjectionRequest は目標達成予測計算リクエスト
type GoalProjectionRequest struct {
	UserID string `json:"user_id" validate:"required"`
	GoalID string `json:"goal_id" validate:"required"`
}

// CalculateAssetProjection は資産推移を計算する
// @Summary 資産推移計算
// @Description 指定年数の資産推移を計算します
// @Tags calculations
// @Accept json
// @Produce json
// @Param request body AssetProjectionRequest true "資産推移計算リクエスト"
// @Success 200 {object} usecases.AssetProjectionOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /calculations/asset-projection [post]
func (c *CalculationsController) CalculateAssetProjection(ctx echo.Context) error {
	var req AssetProjectionRequest
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

	input := usecases.AssetProjectionInput{
		UserID: entities.UserID(req.UserID),
		Years:  req.Years,
	}

	output, err := c.useCase.CalculateAssetProjection(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "資産推移の計算に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// CalculateRetirementProjection は退職資金予測を計算する
// @Summary 退職資金計算
// @Description 退職資金の予測を計算します
// @Tags calculations
// @Accept json
// @Produce json
// @Param request body RetirementCalculationRequest true "退職資金計算リクエスト"
// @Success 200 {object} usecases.RetirementProjectionOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /calculations/retirement [post]
func (c *CalculationsController) CalculateRetirementProjection(ctx echo.Context) error {
	var req RetirementCalculationRequest
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

	input := usecases.RetirementProjectionInput{
		UserID: entities.UserID(req.UserID),
	}

	output, err := c.useCase.CalculateRetirementProjection(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "退職資金計算に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// CalculateEmergencyFundProjection は緊急資金予測を計算する
// @Summary 緊急資金計算
// @Description 緊急資金の予測を計算します
// @Tags calculations
// @Accept json
// @Produce json
// @Param request body EmergencyFundCalculationRequest true "緊急資金計算リクエスト"
// @Success 200 {object} usecases.EmergencyFundProjectionOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /calculations/emergency-fund [post]
func (c *CalculationsController) CalculateEmergencyFundProjection(ctx echo.Context) error {
	var req EmergencyFundCalculationRequest
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

	input := usecases.EmergencyFundProjectionInput{
		UserID: entities.UserID(req.UserID),
	}

	output, err := c.useCase.CalculateEmergencyFundProjection(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "緊急資金計算に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// CalculateComprehensiveProjection は包括的な財務予測を計算する
// @Summary 包括的財務予測計算
// @Description 包括的な財務予測を計算します
// @Tags calculations
// @Accept json
// @Produce json
// @Param request body ComprehensiveProjectionRequest true "包括的予測計算リクエスト"
// @Success 200 {object} usecases.ComprehensiveProjectionOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /calculations/comprehensive [post]
func (c *CalculationsController) CalculateComprehensiveProjection(ctx echo.Context) error {
	var req ComprehensiveProjectionRequest
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

	input := usecases.ComprehensiveProjectionInput{
		UserID: entities.UserID(req.UserID),
		Years:  req.Years,
	}

	output, err := c.useCase.CalculateComprehensiveProjection(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "包括的財務予測の計算に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// CalculateGoalProjection は目標達成予測を計算する
// @Summary 目標達成予測計算
// @Description 目標達成の予測を計算します
// @Tags calculations
// @Accept json
// @Produce json
// @Param request body GoalProjectionRequest true "目標達成予測計算リクエスト"
// @Success 200 {object} usecases.GoalProjectionOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /calculations/goal-projection [post]
func (c *CalculationsController) CalculateGoalProjection(ctx echo.Context) error {
	var req GoalProjectionRequest
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

	input := usecases.GoalProjectionInput{
		UserID: entities.UserID(req.UserID),
		GoalID: entities.GoalID(req.GoalID),
	}

	output, err := c.useCase.CalculateGoalProjection(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "目標達成予測の計算に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}
