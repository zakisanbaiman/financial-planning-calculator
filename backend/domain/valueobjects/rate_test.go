package valueobjects

import (
	"testing"
)

func TestNewRate(t *testing.T) {
	// 正常なケース
	rate, err := NewRate(5.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if rate.AsPercentage() != 5.0 {
		t.Errorf("Expected 5.0%%, got %f%%", rate.AsPercentage())
	}
	if rate.AsDecimal() != 0.05 {
		t.Errorf("Expected 0.05, got %f", rate.AsDecimal())
	}

	// 無効なケース - 負の値
	_, err = NewRate(-1.0)
	if err == nil {
		t.Error("Expected error for negative rate")
	}

	// 無効なケース - 100%を超える値
	_, err = NewRate(101.0)
	if err == nil {
		t.Error("Expected error for rate over 100%")
	}
}

func TestNewRateFromDecimal(t *testing.T) {
	rate, err := NewRateFromDecimal(0.05)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if rate.AsPercentage() != 5.0 {
		t.Errorf("Expected 5.0%%, got %f%%", rate.AsPercentage())
	}
}

func TestRateCompoundFactor(t *testing.T) {
	rate, _ := NewRate(5.0) // 5%

	// 1年間の複利係数
	factor := rate.CompoundFactor(1)
	expected := 1.05
	if factor != expected {
		t.Errorf("Expected %f, got %f", expected, factor)
	}

	// 0年間の複利係数
	factor = rate.CompoundFactor(0)
	if factor != 1.0 {
		t.Errorf("Expected 1.0, got %f", factor)
	}
}

func TestRateMonthlyRate(t *testing.T) {
	annualRate, _ := NewRate(12.0) // 年率12%
	monthlyRate, err := annualRate.MonthlyRate()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 月利は約0.95%になるはず
	if monthlyRate.AsPercentage() < 0.9 || monthlyRate.AsPercentage() > 1.0 {
		t.Errorf("Expected monthly rate around 0.95%%, got %f%%", monthlyRate.AsPercentage())
	}
}

func TestRateIsZero(t *testing.T) {
	zeroRate, _ := NewRate(0.0)
	nonZeroRate, _ := NewRate(5.0)

	if !zeroRate.IsZero() {
		t.Error("Expected zero rate to return true for IsZero()")
	}
	if nonZeroRate.IsZero() {
		t.Error("Expected non-zero rate to return false for IsZero()")
	}
}

func TestRateString(t *testing.T) {
	rate, _ := NewRate(5.25)
	expected := "5.2500%"
	if rate.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, rate.String())
	}
}

func TestRateAdd(t *testing.T) {
	rate1, _ := NewRate(3.0)
	rate2, _ := NewRate(2.0)

	result, err := rate1.Add(rate2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.AsPercentage() != 5.0 {
		t.Errorf("Expected 5.0%%, got %f%%", result.AsPercentage())
	}

	// 100%を超える場合
	rate3, _ := NewRate(98.0)
	rate4, _ := NewRate(5.0)
	_, err = rate3.Add(rate4)
	if err == nil {
		t.Error("Expected error when sum exceeds 100%")
	}
}

func TestRateSubtract(t *testing.T) {
	rate1, _ := NewRate(8.0)
	rate2, _ := NewRate(3.0)

	result, err := rate1.Subtract(rate2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.AsPercentage() != 5.0 {
		t.Errorf("Expected 5.0%%, got %f%%", result.AsPercentage())
	}

	// 負の結果になる場合
	_, err = rate2.Subtract(rate1)
	if err == nil {
		t.Error("Expected error when result would be negative")
	}
}

func TestRateCompoundFactorEdgeCases(t *testing.T) {
	rate, _ := NewRate(5.0)

	// 負の期間（実装では0を返す）
	factor := rate.CompoundFactor(-1)
	expected := 0.0
	if factor != expected {
		t.Errorf("Expected %f, got %f", expected, factor)
	}

	// 大きな期間
	factor = rate.CompoundFactor(100)
	if factor <= 1.0 {
		t.Error("Expected compound factor to be greater than 1 for positive rate and periods")
	}
}
