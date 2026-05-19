package controllers

import (
	"io"
	"net/http"
	"strings"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/labstack/echo/v4"
)

// CSVFinancialDataController はCSVインポート・エクスポートのコントローラー
type CSVFinancialDataController struct {
	useCase usecases.CSVFinancialDataUseCase
}

// NewCSVFinancialDataController は新しいCSVFinancialDataControllerを作成する
func NewCSVFinancialDataController(useCase usecases.CSVFinancialDataUseCase) *CSVFinancialDataController {
	return &CSVFinancialDataController{useCase: useCase}
}

// DownloadCSV は財務データをCSVファイルとして返す
//
// GET /api/financial-data/csv?user_id={user_id}
//
// レスポンスヘッダー:
//
//	Content-Type: text/csv; charset=utf-8
//	Content-Disposition: attachment; filename="financial_data.csv"
func (c *CSVFinancialDataController) DownloadCSV(ctx echo.Context) error {
	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "ユーザーIDは必須です", nil))
	}

	data, err := c.useCase.ExportFinancialDataToCSV(ctx.Request().Context(), usecases.ExportCSVInput{
		UserID: entities.UserID(userID),
	})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "財務データの取得に失敗しました") {
			return ctx.JSON(http.StatusNotFound, NewNotFoundErrorResponse(ctx, "財務データ"))
		}
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, errMsg))
	}

	ctx.Response().Header().Set("Content-Type", "text/csv; charset=utf-8")
	ctx.Response().Header().Set("Content-Disposition", `attachment; filename="financial_data.csv"`)
	return ctx.Blob(http.StatusOK, "text/csv; charset=utf-8", data)
}

// ImportCSV はCSVファイルをアップロードして財務データを保存する
//
// POST /api/financial-data/csv/import
// Content-Type: multipart/form-data
// Form fields: file (CSV), user_id
func (c *CSVFinancialDataController) ImportCSV(ctx echo.Context) error {
	userID := ctx.FormValue("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "ユーザーIDは必須です", nil))
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "CSVファイルが必要です", err.Error()))
	}

	// 1MB 制限
	if fileHeader.Size > 1<<20 {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "ファイルサイズは1MB以下にしてください", nil))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}
	defer file.Close()

	csvData, err := io.ReadAll(file)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, err.Error()))
	}

	output, err := c.useCase.ImportFinancialDataFromCSV(ctx.Request().Context(), usecases.ImportCSVInput{
		UserID:  entities.UserID(userID),
		CSVData: csvData,
	})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "CSVの解析に失敗しました") || strings.Contains(errMsg, "有効な") {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, errMsg, nil))
		}
		return ctx.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(ctx, errMsg))
	}

	return ctx.JSON(http.StatusOK, output)
}
