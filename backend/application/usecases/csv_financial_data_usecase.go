package usecases

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
)

// CSVFinancialDataUseCase はCSVインポート・エクスポートのユースケース
type CSVFinancialDataUseCase interface {
	ExportFinancialDataToCSV(ctx context.Context, input ExportCSVInput) ([]byte, error)
	ImportFinancialDataFromCSV(ctx context.Context, input ImportCSVInput) (*ImportCSVOutput, error)
}

// ExportCSVInput はCSVエクスポートの入力
type ExportCSVInput struct {
	UserID entities.UserID
}

// ImportCSVInput はCSVインポートの入力
type ImportCSVInput struct {
	UserID  entities.UserID
	CSVData []byte
}

// ImportCSVOutput はCSVインポートの出力
type ImportCSVOutput struct {
	*FinancialDataResponse
}

type csvFinancialDataUseCaseImpl struct {
	financialPlanRepo repositories.FinancialPlanRepository
	manageUseCase     ManageFinancialDataUseCase
}

// NewCSVFinancialDataUseCase は新しいCSVFinancialDataUseCaseを生成する
func NewCSVFinancialDataUseCase(
	financialPlanRepo repositories.FinancialPlanRepository,
	manageUseCase ManageFinancialDataUseCase,
) CSVFinancialDataUseCase {
	return &csvFinancialDataUseCaseImpl{
		financialPlanRepo: financialPlanRepo,
		manageUseCase:     manageUseCase,
	}
}

// ExportFinancialDataToCSV は財務データをCSVバイト列として返す
func (uc *csvFinancialDataUseCaseImpl) ExportFinancialDataToCSV(ctx context.Context, input ExportCSVInput) ([]byte, error) {
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務データの取得に失敗しました: %w", err)
	}

	data, err := generateCSVBytes(plan.Profile())
	if err != nil {
		return nil, fmt.Errorf("CSV生成に失敗しました: %w", err)
	}
	return data, nil
}

// ImportFinancialDataFromCSV はCSVバイト列をパースして財務データを保存する
func (uc *csvFinancialDataUseCaseImpl) ImportFinancialDataFromCSV(ctx context.Context, input ImportCSVInput) (*ImportCSVOutput, error) {
	parsed, err := parseCSVBytes(input.CSVData)
	if err != nil {
		return nil, fmt.Errorf("CSVの解析に失敗しました: %w", err)
	}

	expenses := make([]ExpenseItem, len(parsed.Expenses))
	for i, e := range parsed.Expenses {
		desc := e.Description
		expenses[i] = ExpenseItem{Category: e.Category, Amount: e.Amount, Description: &desc}
	}
	savings := make([]SavingsItem, len(parsed.Savings))
	for i, s := range parsed.Savings {
		desc := s.Description
		savings[i] = SavingsItem{Type: s.Type, Amount: s.Amount, Description: &desc}
	}

	exists, err := uc.financialPlanRepo.ExistsByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務データの確認に失敗しました: %w", err)
	}

	if exists {
		output, err := uc.manageUseCase.UpdateFinancialProfile(ctx, UpdateFinancialProfileInput{
			UserID:           input.UserID,
			MonthlyIncome:    parsed.MonthlyIncome,
			MonthlyExpenses:  expenses,
			CurrentSavings:   savings,
			InvestmentReturn: parsed.InvestmentReturn,
			InflationRate:    parsed.InflationRate,
		})
		if err != nil {
			return nil, fmt.Errorf("財務プロファイルの更新に失敗しました: %w", err)
		}
		// UpdateFinancialProfileOutput は *FinancialDataResponse を埋め込んでいる
		return &ImportCSVOutput{FinancialDataResponse: output.FinancialDataResponse}, nil
	}

	_, err = uc.manageUseCase.CreateFinancialPlan(ctx, CreateFinancialPlanInput{
		UserID:           input.UserID,
		MonthlyIncome:    parsed.MonthlyIncome,
		MonthlyExpenses:  expenses,
		CurrentSavings:   savings,
		InvestmentReturn: parsed.InvestmentReturn,
		InflationRate:    parsed.InflationRate,
	})
	if err != nil {
		return nil, fmt.Errorf("財務計画の作成に失敗しました: %w", err)
	}

	return &ImportCSVOutput{FinancialDataResponse: &FinancialDataResponse{UserID: string(input.UserID)}}, nil
}

// ---- CSV生成 ----

// generateCSVBytes は財務プロファイルからCSVバイト列を生成する
//
// フォーマット:
//
//	# SECTION: PROFILE
//	field,value
//	monthly_income,300000
//	...
//
//	# SECTION: EXPENSES
//	category,amount,description
//	生活費,100000,
//	...
func generateCSVBytes(profile *entities.FinancialProfile) ([]byte, error) {
	if profile == nil {
		return nil, fmt.Errorf("財務プロファイルが存在しません")
	}

	var buf bytes.Buffer

	// PROFILE セクション
	fmt.Fprintln(&buf, "# SECTION: PROFILE")
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{"field", "value"}); err != nil {
		return nil, err
	}
	if err := w.Write([]string{"monthly_income", strconv.FormatFloat(profile.MonthlyIncome().Amount(), 'f', -1, 64)}); err != nil {
		return nil, err
	}
	if err := w.Write([]string{"investment_return", strconv.FormatFloat(profile.InvestmentReturn().AsPercentage(), 'f', -1, 64)}); err != nil {
		return nil, err
	}
	if err := w.Write([]string{"inflation_rate", strconv.FormatFloat(profile.InflationRate().AsPercentage(), 'f', -1, 64)}); err != nil {
		return nil, err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	// EXPENSES セクション
	fmt.Fprintln(&buf, "")
	fmt.Fprintln(&buf, "# SECTION: EXPENSES")
	w = csv.NewWriter(&buf)
	if err := w.Write([]string{"category", "amount", "description"}); err != nil {
		return nil, err
	}
	for _, e := range profile.MonthlyExpenses() {
		if err := w.Write([]string{e.Category, strconv.FormatFloat(e.Amount.Amount(), 'f', -1, 64), e.Description}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	// SAVINGS セクション
	fmt.Fprintln(&buf, "")
	fmt.Fprintln(&buf, "# SECTION: SAVINGS")
	w = csv.NewWriter(&buf)
	if err := w.Write([]string{"type", "amount", "description"}); err != nil {
		return nil, err
	}
	for _, s := range profile.CurrentSavings() {
		if err := w.Write([]string{s.Type, strconv.FormatFloat(s.Amount.Amount(), 'f', -1, 64), s.Description}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ---- CSVパース ----

type parsedExpense struct {
	Category    string
	Amount      float64
	Description string
}

type parsedSaving struct {
	Type        string
	Amount      float64
	Description string
}

type parsedCSVData struct {
	MonthlyIncome    float64
	InvestmentReturn float64
	InflationRate    float64
	Expenses         []parsedExpense
	Savings          []parsedSaving
}

type csvSection int

const (
	sectionNone     csvSection = iota
	sectionProfile
	sectionExpenses
	sectionSavings
)

// parseCSVBytes はCSVバイト列をパースして parsedCSVData を返す
func parseCSVBytes(data []byte) (*parsedCSVData, error) {
	result := &parsedCSVData{}
	currentSection := sectionNone
	skipHeader := false

	lines := strings.Split(string(data), "\n")
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		// セクションヘッダー行
		if strings.HasPrefix(line, "# SECTION:") {
			sectionName := strings.TrimSpace(strings.TrimPrefix(line, "# SECTION:"))
			switch sectionName {
			case "PROFILE":
				currentSection = sectionProfile
			case "EXPENSES":
				currentSection = sectionExpenses
			case "SAVINGS":
				currentSection = sectionSavings
			default:
				currentSection = sectionNone
			}
			skipHeader = true
			continue
		}

		// ヘッダー行をスキップ
		if skipHeader {
			skipHeader = false
			continue
		}

		// CSVパース（1行ずつ）
		r := csv.NewReader(strings.NewReader(line))
		r.TrimLeadingSpace = true
		fields, err := r.Read()
		if err != nil {
			continue
		}

		switch currentSection {
		case sectionProfile:
			if len(fields) < 2 {
				continue
			}
			val, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
			if err != nil {
				continue
			}
			switch strings.TrimSpace(fields[0]) {
			case "monthly_income":
				result.MonthlyIncome = val
			case "investment_return":
				result.InvestmentReturn = val
			case "inflation_rate":
				result.InflationRate = val
			}

		case sectionExpenses:
			if len(fields) < 2 {
				continue
			}
			amount, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
			if err != nil {
				continue
			}
			desc := ""
			if len(fields) >= 3 {
				desc = strings.TrimSpace(fields[2])
			}
			result.Expenses = append(result.Expenses, parsedExpense{
				Category:    strings.TrimSpace(fields[0]),
				Amount:      amount,
				Description: desc,
			})

		case sectionSavings:
			if len(fields) < 2 {
				continue
			}
			amount, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
			if err != nil {
				continue
			}
			desc := ""
			if len(fields) >= 3 {
				desc = strings.TrimSpace(fields[2])
			}
			result.Savings = append(result.Savings, parsedSaving{
				Type:        strings.TrimSpace(fields[0]),
				Amount:      amount,
				Description: desc,
			})
		}
	}

	if result.MonthlyIncome <= 0 {
		return nil, fmt.Errorf("CSVのPROFILEセクションに有効なmonthly_incomeがありません")
	}

	return result, nil
}
