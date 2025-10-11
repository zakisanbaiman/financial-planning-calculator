package valueobjects

import (
	"errors"
	"fmt"
	"math"
)

// Rate はパーセンテージ利率を表す値オブジェクト（例：金利、インフレ率）
type Rate struct {
	value float64 // パーセンテージで保存（例：5%の場合は5.0）
}

// NewRate は新しいRate値オブジェクトを作成する（バリデーション付き）
func NewRate(percentage float64) (Rate, error) {
	if math.IsNaN(percentage) || math.IsInf(percentage, 0) {
		return Rate{}, errors.New("利率にNaNや無限大は指定できません")
	}

	if percentage < 0 {
		return Rate{}, errors.New("利率は負の値にできません")
	}

	if percentage > 100 {
		return Rate{}, errors.New("利率は100%を超えることはできません")
	}

	// 精度のため小数点以下4桁で丸める
	roundedValue := math.Round(percentage*10000) / 10000

	return Rate{
		value: roundedValue,
	}, nil
}

// NewRateFromDecimal は小数値からRateを作成する（例：5%の場合は0.05）
func NewRateFromDecimal(decimal float64) (Rate, error) {
	return NewRate(decimal * 100)
}

// AsDecimal は利率を小数で返す（例：5%の場合は0.05）
func (r Rate) AsDecimal() float64 {
	return r.value / 100
}

// AsPercentage は利率をパーセンテージで返す（例：5%の場合は5.0）
func (r Rate) AsPercentage() float64 {
	return r.value
}

// IsValid は利率が有効かどうかを返す（非負かつ100%以下）
func (r Rate) IsValid() bool {
	return r.value >= 0 && r.value <= 100
}

// IsZero は利率がゼロかどうかを返す
func (r Rate) IsZero() bool {
	return math.Abs(r.value) < 0.0001 // 0.0001%未満の利率はゼロとみなす
}

// Add は別の利率をこの利率に加算する
func (r Rate) Add(other Rate) (Rate, error) {
	return NewRate(r.value + other.value)
}

// Subtract は別の利率をこの利率から減算する
func (r Rate) Subtract(other Rate) (Rate, error) {
	return NewRate(r.value - other.value)
}

// Multiply は利率に係数を乗算する
func (r Rate) Multiply(factor float64) (Rate, error) {
	if math.IsNaN(factor) || math.IsInf(factor, 0) {
		return Rate{}, errors.New("係数にNaNや無限大は指定できません")
	}

	if factor < 0 {
		return Rate{}, errors.New("係数は負の値にできません")
	}

	return NewRate(r.value * factor)
}

// GreaterThan はこの利率が他の利率より大きいかどうかを返す
func (r Rate) GreaterThan(other Rate) bool {
	return r.value > other.value
}

// LessThan はこの利率が他の利率より小さいかどうかを返す
func (r Rate) LessThan(other Rate) bool {
	return r.value < other.value
}

// Equal はこの利率が他の利率と等しいかどうかを返す
func (r Rate) Equal(other Rate) bool {
	return math.Abs(r.value-other.value) < 0.0001
}

// String は利率の文字列表現を返す
func (r Rate) String() string {
	return fmt.Sprintf("%.4f%%", r.value)
}

// CompoundFactor は指定された期間数に対する複利係数を計算する
// 計算式: (1 + rate)^periods
func (r Rate) CompoundFactor(periods int) float64 {
	if periods < 0 {
		return 0
	}

	if periods == 0 {
		return 1
	}

	return math.Pow(1+r.AsDecimal(), float64(periods))
}

// MonthlyRate は年利を月利に変換する
func (r Rate) MonthlyRate() (Rate, error) {
	// 年利を月利に変換: (1 + annual_rate)^(1/12) - 1
	monthlyDecimal := math.Pow(1+r.AsDecimal(), 1.0/12.0) - 1
	return NewRateFromDecimal(monthlyDecimal)
}

// AnnualRate は月利を年利に変換する
func (r Rate) AnnualRate() (Rate, error) {
	// 月利を年利に変換: (1 + monthly_rate)^12 - 1
	annualDecimal := math.Pow(1+r.AsDecimal(), 12.0) - 1
	return NewRateFromDecimal(annualDecimal)
}
