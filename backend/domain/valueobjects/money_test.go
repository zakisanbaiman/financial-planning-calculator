package valueobjects

import (
	"math"
	"testing"
)

func TestNewMoney(t *testing.T) {
	// 正常なケース
	money, err := NewMoney(1000.50, JPY)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if money.Amount() != 1000.50 {
		t.Errorf("Expected amount 1000.50, got %f", money.Amount())
	}
	if money.Currency() != JPY {
		t.Errorf("Expected currency JPY, got %s", money.Currency())
	}

	// 無効なケース - NaN
	_, err = NewMoney(math.NaN(), JPY)
	if err == nil {
		t.Error("Expected error for NaN amount")
	}

	// 無効なケース - 空の通貨
	_, err = NewMoney(1000, "")
	if err == nil {
		t.Error("Expected error for empty currency")
	}
}

func TestMoneyAdd(t *testing.T) {
	money1, _ := NewMoney(1000, JPY)
	money2, _ := NewMoney(500, JPY)

	result, err := money1.Add(money2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Amount() != 1500 {
		t.Errorf("Expected 1500, got %f", result.Amount())
	}

	// 異なる通貨での加算
	moneyUSD, _ := NewMoney(100, USD)
	_, err = money1.Add(moneyUSD)
	if err == nil {
		t.Error("Expected error when adding different currencies")
	}
}

func TestMoneyIsPositive(t *testing.T) {
	positive, _ := NewMoney(100, JPY)
	negative, _ := NewMoney(-100, JPY)
	zero, _ := NewMoney(0, JPY)

	if !positive.IsPositive() {
		t.Error("Expected positive money to return true for IsPositive()")
	}
	if negative.IsPositive() {
		t.Error("Expected negative money to return false for IsPositive()")
	}
	if zero.IsPositive() {
		t.Error("Expected zero money to return false for IsPositive()")
	}
}

func TestMoneySubtract(t *testing.T) {
	money1, _ := NewMoney(1000, JPY)
	money2, _ := NewMoney(300, JPY)

	result, err := money1.Subtract(money2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Amount() != 700 {
		t.Errorf("Expected 700, got %f", result.Amount())
	}

	// 異なる通貨での減算
	moneyUSD, _ := NewMoney(100, USD)
	_, err = money1.Subtract(moneyUSD)
	if err == nil {
		t.Error("Expected error when subtracting different currencies")
	}
}

func TestMoneyMultiply(t *testing.T) {
	money, _ := NewMoney(1000, JPY)
	rate, _ := NewRate(5.0) // 5%

	result, err := money.Multiply(rate)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Amount() != 50 { // 1000 * 0.05 = 50
		t.Errorf("Expected 50, got %f", result.Amount())
	}
}

func TestMoneyMultiplyByFloat(t *testing.T) {
	money, _ := NewMoney(1000, JPY)

	result, err := money.MultiplyByFloat(1.5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Amount() != 1500 {
		t.Errorf("Expected 1500, got %f", result.Amount())
	}

	// NaNでの乗算
	_, err = money.MultiplyByFloat(math.NaN())
	if err == nil {
		t.Error("Expected error for NaN multiplier")
	}

	// 無限大での乗算
	_, err = money.MultiplyByFloat(math.Inf(1))
	if err == nil {
		t.Error("Expected error for infinite multiplier")
	}
}

func TestMoneyIsNegative(t *testing.T) {
	positive, _ := NewMoney(100, JPY)
	negative, _ := NewMoney(-100, JPY)
	zero, _ := NewMoney(0, JPY)

	if positive.IsNegative() {
		t.Error("Expected positive money to return false for IsNegative()")
	}
	if !negative.IsNegative() {
		t.Error("Expected negative money to return true for IsNegative()")
	}
	if zero.IsNegative() {
		t.Error("Expected zero money to return false for IsNegative()")
	}
}

func TestMoneyIsZero(t *testing.T) {
	zero, _ := NewMoney(0, JPY)
	almostZero, _ := NewMoney(0.005, JPY) // 0.5セント（1セント未満なのでゼロとみなされる）
	notZero, _ := NewMoney(0.02, JPY)     // 2セント

	if !zero.IsZero() {
		t.Error("Expected zero money to return true for IsZero()")
	}
	if !almostZero.IsZero() {
		t.Errorf("Expected almost zero money to return true for IsZero(), amount: %f", almostZero.Amount())
	}
	if notZero.IsZero() {
		t.Error("Expected non-zero money to return false for IsZero()")
	}
}

func TestMoneyComparisons(t *testing.T) {
	money1, _ := NewMoney(1000, JPY)
	money2, _ := NewMoney(500, JPY)
	money3, _ := NewMoney(1000, JPY)
	moneyUSD, _ := NewMoney(1000, USD)

	// GreaterThan
	greater, err := money1.GreaterThan(money2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !greater {
		t.Error("Expected 1000 > 500 to be true")
	}

	// LessThan
	less, err := money2.LessThan(money1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !less {
		t.Error("Expected 500 < 1000 to be true")
	}

	// Equal
	equal, err := money1.Equal(money3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !equal {
		t.Error("Expected 1000 == 1000 to be true")
	}

	// 異なる通貨での比較
	_, err = money1.GreaterThan(moneyUSD)
	if err == nil {
		t.Error("Expected error when comparing different currencies")
	}
}

func TestMoneyString(t *testing.T) {
	money, _ := NewMoney(1234.56, JPY)
	expected := "1234.56 JPY"
	if money.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, money.String())
	}
}

func TestMoneyAbs(t *testing.T) {
	positive, _ := NewMoney(100, JPY)
	negative, _ := NewMoney(-100, JPY)

	absPositive, err := positive.Abs()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if absPositive.Amount() != 100 {
		t.Errorf("Expected 100, got %f", absPositive.Amount())
	}

	absNegative, err := negative.Abs()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if absNegative.Amount() != 100 {
		t.Errorf("Expected 100, got %f", absNegative.Amount())
	}
}

func TestNewMoneyJPY(t *testing.T) {
	money, err := NewMoneyJPY(1000)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if money.Currency() != JPY {
		t.Errorf("Expected JPY currency, got %s", money.Currency())
	}
	if money.Amount() != 1000 {
		t.Errorf("Expected 1000, got %f", money.Amount())
	}
}
