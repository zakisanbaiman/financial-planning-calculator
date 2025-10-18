package valueobjects

import (
	"testing"
)

func TestNewPeriod(t *testing.T) {
	// 正常なケース
	period, err := NewPeriod(2, 6)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if period.Years() != 2 {
		t.Errorf("Expected 2 years, got %d", period.Years())
	}
	if period.Months() != 6 {
		t.Errorf("Expected 6 months, got %d", period.Months())
	}

	// 月数の正規化テスト（15ヶ月 = 1年3ヶ月）
	period, err = NewPeriod(0, 15)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if period.Years() != 1 {
		t.Errorf("Expected 1 year after normalization, got %d", period.Years())
	}
	if period.Months() != 3 {
		t.Errorf("Expected 3 months after normalization, got %d", period.Months())
	}

	// 無効なケース - 負の年数
	_, err = NewPeriod(-1, 0)
	if err == nil {
		t.Error("Expected error for negative years")
	}
}

func TestPeriodToMonths(t *testing.T) {
	period, _ := NewPeriod(2, 6)
	totalMonths := period.ToMonths()
	expected := 30 // 2年 * 12ヶ月 + 6ヶ月
	if totalMonths != expected {
		t.Errorf("Expected %d months, got %d", expected, totalMonths)
	}
}

func TestPeriodToYears(t *testing.T) {
	period, _ := NewPeriod(2, 6)
	totalYears := period.ToYears()
	expected := 2.5 // 2年 + 6ヶ月/12
	if totalYears != expected {
		t.Errorf("Expected %f years, got %f", expected, totalYears)
	}
}

func TestNewPeriodFromYearsFloat(t *testing.T) {
	period, err := NewPeriodFromYearsFloat(2.5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if period.Years() != 2 {
		t.Errorf("Expected 2 years, got %d", period.Years())
	}
	if period.Months() != 6 {
		t.Errorf("Expected 6 months, got %d", period.Months())
	}
}

func TestPeriodAdd(t *testing.T) {
	period1, _ := NewPeriod(1, 6)
	period2, _ := NewPeriod(0, 8)

	result, err := period1.Add(period2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 1年6ヶ月 + 8ヶ月 = 2年2ヶ月
	if result.Years() != 2 {
		t.Errorf("Expected 2 years, got %d", result.Years())
	}
	if result.Months() != 2 {
		t.Errorf("Expected 2 months, got %d", result.Months())
	}
}

func TestPeriodString(t *testing.T) {
	tests := []struct {
		years    int
		months   int
		expected string
	}{
		{0, 0, "0ヶ月"},
		{0, 1, "1ヶ月"},
		{0, 5, "5ヶ月"},
		{1, 0, "1年"},
		{3, 0, "3年"},
		{1, 1, "1年1ヶ月"},
		{2, 5, "2年5ヶ月"},
	}

	for _, test := range tests {
		period, _ := NewPeriod(test.years, test.months)
		result := period.String()
		if result != test.expected {
			t.Errorf("For %d years %d months, expected '%s', got '%s'",
				test.years, test.months, test.expected, result)
		}
	}
}

func TestNewPeriodFromMonths(t *testing.T) {
	period, err := NewPeriodFromMonths(30) // 30ヶ月 = 2年6ヶ月
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if period.Years() != 2 {
		t.Errorf("Expected 2 years, got %d", period.Years())
	}
	if period.Months() != 6 {
		t.Errorf("Expected 6 months, got %d", period.Months())
	}

	// 負の月数
	_, err = NewPeriodFromMonths(-5)
	if err == nil {
		t.Error("Expected error for negative months")
	}
}

func TestPeriodSubtract(t *testing.T) {
	period1, _ := NewPeriod(2, 8)
	period2, _ := NewPeriod(1, 3)

	result, err := period1.Subtract(period2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 2年8ヶ月 - 1年3ヶ月 = 1年5ヶ月
	if result.Years() != 1 {
		t.Errorf("Expected 1 year, got %d", result.Years())
	}
	if result.Months() != 5 {
		t.Errorf("Expected 5 months, got %d", result.Months())
	}

	// 負の結果になる場合
	_, err = period2.Subtract(period1)
	if err == nil {
		t.Error("Expected error when result would be negative")
	}
}

func TestPeriodIsZero(t *testing.T) {
	zeroPeriod, _ := NewPeriod(0, 0)
	nonZeroPeriod, _ := NewPeriod(1, 0)

	if !zeroPeriod.IsZero() {
		t.Error("Expected zero period to return true for IsZero()")
	}
	if nonZeroPeriod.IsZero() {
		t.Error("Expected non-zero period to return false for IsZero()")
	}
}
