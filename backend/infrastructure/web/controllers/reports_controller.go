package controllers

import (
	"net/http"
	"strconv"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/labstack/echo/v4"
)

// ReportsController はレポート生成のコントローラー
type ReportsController struct {
	useCase usecases.GenerateReportsUseCase
}

// NewReportsController は新しいReportsControllerを作成する
func NewReportsController(useCase usecases.GenerateReportsUseCase) *ReportsController {
	return &ReportsController{
		useCase: useCase,
	}
}

// FinancialSummaryReportRequest は財務サマリーレポート生成リクエスト
type FinancialSummaryReportRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// AssetProjectionReportRequest は資産推移レポート生成リクエスト
type AssetProjectionReportRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Years  int    `json:"years" validate:"required,gte=1,lte=50"`
}

// GoalsProgressReportRequest は目標進捗レポート生成リクエスト
type GoalsProgressReportRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// RetirementPlanReportRequest は退職計画レポート生成リクエスト
type RetirementPlanReportRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// ComprehensiveReportRequest は包括的レポート生成リクエスト
type ComprehensiveReportRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Years  int    `json:"years" validate:"required,gte=1,lte=50"`
}

// ExportReportRequest はレポートエクスポートリクエスト
type ExportReportRequest struct {
	UserID     string      `json:"user_id" validate:"required"`
	ReportType string      `json:"report_type" validate:"required,oneof=financial_summary asset_projection goals_progress retirement_plan comprehensive"`
	Format     string      `json:"format" validate:"required,oneof=pdf excel csv"`
	ReportData interface{} `json:"report_data" validate:"required"`
}

// GenerateFinancialSummaryReport は財務サマリーレポートを生成する
// @Summary 財務サマリーレポート生成
// @Description 財務サマリーレポートを生成します
// @Tags reports
// @Accept json
// @Produce json
// @Param request body FinancialSummaryReportRequest true "財務サマリーレポート生成リクエスト"
// @Success 200 {object} usecases.FinancialSummaryReportOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/financial-summary [post]
func (c *ReportsController) GenerateFinancialSummaryReport(ctx echo.Context) error {
	var req FinancialSummaryReportRequest
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

	input := usecases.FinancialSummaryReportInput{
		UserID: entities.UserID(req.UserID),
	}

	output, err := c.useCase.GenerateFinancialSummaryReport(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "財務サマリーレポートの生成に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// GenerateAssetProjectionReport は資産推移レポートを生成する
// @Summary 資産推移レポート生成
// @Description 資産推移レポートを生成します
// @Tags reports
// @Accept json
// @Produce json
// @Param request body AssetProjectionReportRequest true "資産推移レポート生成リクエスト"
// @Success 200 {object} usecases.AssetProjectionReportOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/asset-projection [post]
func (c *ReportsController) GenerateAssetProjectionReport(ctx echo.Context) error {
	var req AssetProjectionReportRequest
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

	input := usecases.AssetProjectionReportInput{
		UserID: entities.UserID(req.UserID),
		Years:  req.Years,
	}

	output, err := c.useCase.GenerateAssetProjectionReport(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "資産推移レポートの生成に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// GenerateGoalsProgressReport は目標進捗レポートを生成する
// @Summary 目標進捗レポート生成
// @Description 目標進捗レポートを生成します
// @Tags reports
// @Accept json
// @Produce json
// @Param request body GoalsProgressReportRequest true "目標進捗レポート生成リクエスト"
// @Success 200 {object} usecases.GoalsProgressReportOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/goals-progress [post]
func (c *ReportsController) GenerateGoalsProgressReport(ctx echo.Context) error {
	var req GoalsProgressReportRequest
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

	input := usecases.GoalsProgressReportInput{
		UserID: entities.UserID(req.UserID),
	}

	output, err := c.useCase.GenerateGoalsProgressReport(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "目標進捗レポートの生成に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// GenerateRetirementPlanReport は退職計画レポートを生成する
// @Summary 退職計画レポート生成
// @Description 退職計画レポートを生成します
// @Tags reports
// @Accept json
// @Produce json
// @Param request body RetirementPlanReportRequest true "退職計画レポート生成リクエスト"
// @Success 200 {object} usecases.RetirementPlanReportOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/retirement-plan [post]
func (c *ReportsController) GenerateRetirementPlanReport(ctx echo.Context) error {
	var req RetirementPlanReportRequest
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

	input := usecases.RetirementPlanReportInput{
		UserID: entities.UserID(req.UserID),
	}

	output, err := c.useCase.GenerateRetirementPlanReport(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "退職計画レポートの生成に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// GenerateComprehensiveReport は包括的レポートを生成する
// @Summary 包括的レポート生成
// @Description 包括的レポートを生成します
// @Tags reports
// @Accept json
// @Produce json
// @Param request body ComprehensiveReportRequest true "包括的レポート生成リクエスト"
// @Success 200 {object} usecases.ComprehensiveReportOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/comprehensive [post]
func (c *ReportsController) GenerateComprehensiveReport(ctx echo.Context) error {
	var req ComprehensiveReportRequest
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

	input := usecases.ComprehensiveReportInput{
		UserID: entities.UserID(req.UserID),
		Years:  req.Years,
	}

	output, err := c.useCase.GenerateComprehensiveReport(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "包括的レポートの生成に失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// ExportReportToPDF はレポートをPDF形式でエクスポートする
// @Summary レポートPDFエクスポート
// @Description レポートをPDF形式でエクスポートします
// @Tags reports
// @Accept json
// @Produce json
// @Param request body ExportReportRequest true "レポートエクスポートリクエスト"
// @Success 200 {object} usecases.ExportReportOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/export [post]
func (c *ReportsController) ExportReportToPDF(ctx echo.Context) error {
	var req ExportReportRequest
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

	input := usecases.ExportReportInput{
		UserID:     entities.UserID(req.UserID),
		ReportType: req.ReportType,
		Format:     req.Format,
		ReportData: req.ReportData,
	}

	output, err := c.useCase.ExportReportToPDF(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "レポートのエクスポートに失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// GetReportPDF はPDFレポートを取得する（クエリパラメータ版）
// @Summary PDFレポート取得
// @Description PDFレポートを取得します
// @Tags reports
// @Produce json
// @Param user_id query string true "ユーザーID"
// @Param report_type query string false "レポートタイプ" Enums(financial_summary, comprehensive)
// @Param years query int false "予測年数" default(10)
// @Success 200 {object} usecases.ExportReportOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/pdf [get]
func (c *ReportsController) GetReportPDF(ctx echo.Context) error {
	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ユーザーIDは必須です",
		})
	}

	reportType := ctx.QueryParam("report_type")
	if reportType == "" {
		reportType = "comprehensive" // デフォルトは包括的レポート
	}

	yearsStr := ctx.QueryParam("years")
	years := 10 // デフォルト値
	if yearsStr != "" {
		if parsedYears, err := strconv.Atoi(yearsStr); err == nil && parsedYears > 0 && parsedYears <= 50 {
			years = parsedYears
		}
	}

	// レポートタイプに応じて適切なレポートを生成
	var reportData interface{}
	var err error

	switch reportType {
	case "financial_summary":
		input := usecases.FinancialSummaryReportInput{
			UserID: entities.UserID(userID),
		}
		output, genErr := c.useCase.GenerateFinancialSummaryReport(ctx.Request().Context(), input)
		if genErr != nil {
			return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "財務サマリーレポートの生成に失敗しました",
				Details: genErr.Error(),
			})
		}
		reportData = output.Report

	case "comprehensive":
		input := usecases.ComprehensiveReportInput{
			UserID: entities.UserID(userID),
			Years:  years,
		}
		output, genErr := c.useCase.GenerateComprehensiveReport(ctx.Request().Context(), input)
		if genErr != nil {
			return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "包括的レポートの生成に失敗しました",
				Details: genErr.Error(),
			})
		}
		reportData = output.Report

	default:
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "サポートされていないレポートタイプです",
		})
	}

	// PDFエクスポート
	exportInput := usecases.ExportReportInput{
		UserID:     entities.UserID(userID),
		ReportType: reportType,
		Format:     "pdf",
		ReportData: reportData,
	}

	output, err := c.useCase.ExportReportToPDF(ctx.Request().Context(), exportInput)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "PDFエクスポートに失敗しました",
			Details: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, output)
}

// DownloadReport はレポートファイルをダウンロードする
// @Summary レポートダウンロード
// @Description 生成されたレポートファイルをダウンロードします
// @Tags reports
// DownloadReport はトークンを使ってレポートをダウンロードする
// @Summary レポートのダウンロード
// @Description 署名付きトークンを使用してレポートファイルをダウンロードします
// @Tags reports
// @Accept json
// @Produce application/pdf
// @Param token path string true "ダウンロードトークン"
// @Success 200 {file} binary "PDFファイル"
// @Failure 400 {object} map[string]interface{} "無効なリクエスト"
// @Failure 404 {object} map[string]interface{} "ファイルが見つかりません"
// @Failure 410 {object} map[string]interface{} "ファイルの有効期限が切れています"
// @Router /reports/download/{token} [get]
func (ctrl *ReportsController) DownloadReport(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "トークンが指定されていません",
		})
	}

	// TODO: 実際の実装では、TemporaryFileStorageからファイルを取得
	// ここでは簡易的なレスポンスを返す
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ダウンロード機能は実装中です",
		"token":   token,
		"note":    "実際の実装では、PDFファイルがダウンロードされます",
	})
}
