package valueobjects

import (
	"errors"
	"fmt"
	"math"
)

// Currency は通貨の種類を表す
type Currency string

const (
	JPY Currency = "JPY" // 日本円
	USD Currency = "USD" // 米ドル
	EUR Currency = "EUR" // ユーロ
)

// Money は通貨付きの金額を表す値オブジェクト
// 不変性を保証し、同一通貨間でのみ演算を許可する
type Money struct {
	amount   float64  // 金額（小数点以下2桁で丸められる）
	currency Currency // 通貨
}

// NewMoney は新しいMoney値オブジェクトを作成する（バリデーション付き）
func NewMoney(amount float64, currency Currency) (Money, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return Money{}, errors.New("金額にNaNや無限大は指定できません")
	}

	if currency == "" {
		return Money{}, errors.New("通貨は空にできません")
	}

	// 浮動小数点の精度問題を避けるため、小数点以下2桁で丸める
	roundedAmount := math.Round(amount*100) / 100

	return Money{
		amount:   roundedAmount,
		currency: currency,
	}, nil
}

// NewMoneyJPY は日本円のMoney値オブジェクトを作成する
func NewMoneyJPY(amount float64) (Money, error) {
	return NewMoney(amount, JPY)
}

// Amount は金額を返す
func (m Money) Amount() float64 {
	return m.amount
}

// Currency は通貨を返す
func (m Money) Currency() Currency {
	return m.currency
}

// Add は別のMoney値を加算する（同一通貨のみ）
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("異なる通貨は加算できません: %s と %s", m.currency, other.currency)
	}

	return NewMoney(m.amount+other.amount, m.currency)
}

// Subtract は別のMoney値を減算する（同一通貨のみ）
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("異なる通貨は減算できません: %s と %s", m.currency, other.currency)
	}

	return NewMoney(m.amount-other.amount, m.currency)
}

// Multiply は金額に利率を乗算する
func (m Money) Multiply(rate Rate) (Money, error) {
	return NewMoney(m.amount*rate.AsDecimal(), m.currency)
}

// MultiplyByFloat は金額に浮動小数点数を乗算する
func (m Money) MultiplyByFloat(multiplier float64) (Money, error) {
	if math.IsNaN(multiplier) || math.IsInf(multiplier, 0) {
		return Money{}, errors.New("乗数にNaNや無限大は指定できません")
	}

	return NewMoney(m.amount*multiplier, m.currency)
}

// IsPositive は金額が正の値かどうかを返す
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// IsNegative は金額が負の値かどうかを返す
func (m Money) IsNegative() bool {
	return m.amount < 0
}

// IsZero は金額がゼロかどうかを返す
func (m Money) IsZero() bool {
	return math.Abs(m.amount) < 0.01 // 1セント未満の金額はゼロとみなす
}

// GreaterThan はこの金額が他の金額より大きいかどうかを返す
func (m Money) GreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("異なる通貨は比較できません: %s と %s", m.currency, other.currency)
	}

	return m.amount > other.amount, nil
}

// LessThan はこの金額が他の金額より小さいかどうかを返す
func (m Money) LessThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("異なる通貨は比較できません: %s と %s", m.currency, other.currency)
	}

	return m.amount < other.amount, nil
}

// Equal はこの金額が他の金額と等しいかどうかを返す
func (m Money) Equal(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("異なる通貨は比較できません: %s と %s", m.currency, other.currency)
	}

	return math.Abs(m.amount-other.amount) < 0.01, nil
}

// String は金額の文字列表現を返す
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}

// Abs は金額の絶対値を返す
func (m Money) Abs() (Money, error) {
	return NewMoney(math.Abs(m.amount), m.currency)
}
