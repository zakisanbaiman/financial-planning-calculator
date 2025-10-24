package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HealthCheckHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "財務計画計算機 API サーバーが正常に動作しています")
}

func TestAPIInfoHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := APIInfoHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "財務計画計算機 API v1.0")
}

func TestPlaceholderHandlers(t *testing.T) {
	e := echo.New()

	testCases := []struct {
		name    string
		handler echo.HandlerFunc
		message string
	}{
		{"CreateFinancialData", CreateFinancialDataHandler, "財務データ作成機能は実装予定です"},
		{"GetFinancialData", GetFinancialDataHandler, "財務データ取得機能は実装予定です"},
		{"AssetProjection", AssetProjectionHandler, "資産推移計算機能は実装予定です"},
		{"RetirementCalculation", RetirementCalculationHandler, "老後資金計算機能は実装予定です"},
		{"EmergencyFund", EmergencyFundHandler, "緊急資金計算機能は実装予定です"},
		{"CreateGoal", CreateGoalHandler, "目標作成機能は実装予定です"},
		{"GetGoals", GetGoalsHandler, "目標取得機能は実装予定です"},
		{"GeneratePDFReport", GeneratePDFReportHandler, "PDFレポート生成機能は実装予定です"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := tc.handler(c)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotImplemented, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.message)
		})
	}
}
