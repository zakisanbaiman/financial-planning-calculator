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
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// Business logic validation for asset projection
	if err := ValidateBusinessLogic(ctx,
		func() *BusinessLogicError {
			// 要件2.4: 年数の妥当性チェック
			if req.Years < 1 || req.Years > 50 {
				return CreateBusinessLogicError(
					"INVALID_PROJECTION_YEARS",
					"予測年数は1年から50年の範囲で入力してください",
					"現実的な予測期間を入力してください",
					req.Years,
					"1-50年",
				)
			}
			return nil
		},
	); err != nil {
		return err
	}

	// リクエストIDをコンテキストに追加
	reqCtx := GetRequestContextWithUserID(ctx, req.UserID)

	input := usecases.AssetProjectionInput{
		UserID: entities.UserID(req.UserID),
		Years:  req.Years,
	}

	output, err := c.useCase.CalculateAssetProjection(reqCtx, input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
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
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// リクエストIDをコンテキストに追加
	reqCtx := GetRequestContextWithUserID(ctx, req.UserID)

	input := usecases.RetirementProjectionInput{
		UserID: entities.UserID(req.UserID),
	}

	output, err := c.useCase.CalculateRetirementProjection(reqCtx, input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
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
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// リクエストIDをコンテキストに追加
	reqCtx := GetRequestContextWithUserID(ctx, req.UserID)

	input := usecases.EmergencyFundProjectionInput{
		UserID: entities.UserID(req.UserID),
	}

	output, err := c.useCase.CalculateEmergencyFundProjection(reqCtx, input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
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
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// Business logic validation
	if err := ValidateBusinessLogic(ctx,
		func() *BusinessLogicError {
			if req.Years < 1 || req.Years > 50 {
				return CreateBusinessLogicError(
					"INVALID_PROJECTION_YEARS",
					"予測年数は1年から50年の範囲で入力してください",
					"現実的な予測期間を入力してください",
					req.Years,
					"1-50年",
				)
			}
			return nil
		},
	); err != nil {
		return err
	}

	// リクエストIDをコンテキストに追加
	reqCtx := GetRequestContextWithUserID(ctx, req.UserID)

	input := usecases.ComprehensiveProjectionInput{
		UserID: entities.UserID(req.UserID),
		Years:  req.Years,
	}

	output, err := c.useCase.CalculateComprehensiveProjection(reqCtx, input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
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
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// リクエストIDをコンテキストに追加
	reqCtx := GetRequestContextWithUserID(ctx, req.UserID)

	input := usecases.GoalProjectionInput{
		UserID: entities.UserID(req.UserID),
		GoalID: entities.GoalID(req.GoalID),
	}

	output, err := c.useCase.CalculateGoalProjection(reqCtx, input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	return ctx.JSON(http.StatusOK, output)
}
