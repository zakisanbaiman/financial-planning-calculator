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
