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
