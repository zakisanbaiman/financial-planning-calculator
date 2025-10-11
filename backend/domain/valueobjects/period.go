package valueobjects

import (
	"errors"
	"fmt"
	"math"
)

// Period は年と月で表される期間を表す値オブジェクト
type Period struct {
	years  int // 年数
	months int // 月数（0-11の範囲で正規化される）
}

// NewPeriod は新しいPeriod値オブジェクトを作成する（バリデーション付き）
func NewPeriod(years, months int) (Period, error) {
	if years < 0 {
		return Period{}, errors.New("年数は負の値にできません")
	}

	if months < 0 {
		return Period{}, errors.New("月数は負の値にできません")
	}

	// 月数を正規化（余分な月数を年数に変換）
	normalizedYears := years + (months / 12)
	normalizedMonths := months % 12

	return Period{
		years:  normalizedYears,
		months: normalizedMonths,
	}, nil
}

// NewPeriodFromYears は年数のみからPeriodを作成する
func NewPeriodFromYears(years int) (Period, error) {
	return NewPeriod(years, 0)
}

// NewPeriodFromMonths は月数のみからPeriodを作成する
func NewPeriodFromMonths(months int) (Period, error) {
	return NewPeriod(0, months)
}

// NewPeriodFromYearsFloat は小数年数からPeriodを作成する
func NewPeriodFromYearsFloat(years float64) (Period, error) {
	if years < 0 {
		return Period{}, errors.New("年数は負の値にできません")
	}

	if math.IsNaN(years) || math.IsInf(years, 0) {
		return Period{}, errors.New("年数にNaNや無限大は指定できません")
	}

	totalMonths := int(math.Round(years * 12))
	return NewPeriodFromMonths(totalMonths)
}

// Years は年数部分を返す
func (p Period) Years() int {
	return p.years
}

// Months は月数部分を返す
func (p Period) Months() int {
	return p.months
}

// ToMonths は期間全体を月数に変換する
func (p Period) ToMonths() int {
	return p.years*12 + p.months
}

// ToYears は期間全体を年数（浮動小数点）に変換する
func (p Period) ToYears() float64 {
	return float64(p.years) + float64(p.months)/12.0
}

// Add は別の期間をこの期間に加算する
func (p Period) Add(other Period) (Period, error) {
	return NewPeriod(p.years+other.years, p.months+other.months)
}

// Subtract は別の期間をこの期間から減算する
func (p Period) Subtract(other Period) (Period, error) {
	totalMonths := p.ToMonths() - other.ToMonths()
	if totalMonths < 0 {
		return Period{}, errors.New("結果の期間は負の値にできません")
	}

	return NewPeriodFromMonths(totalMonths)
}

// Multiply は期間に係数を乗算する
func (p Period) Multiply(factor float64) (Period, error) {
	if factor < 0 {
		return Period{}, errors.New("係数は負の値にできません")
	}

	if math.IsNaN(factor) || math.IsInf(factor, 0) {
		return Period{}, errors.New("係数にNaNや無限大は指定できません")
	}

	totalMonths := float64(p.ToMonths()) * factor
	return NewPeriodFromMonths(int(math.Round(totalMonths)))
}

// IsZero は期間がゼロかどうかを返す
func (p Period) IsZero() bool {
	return p.years == 0 && p.months == 0
}

// IsPositive は期間が正の値かどうかを返す
func (p Period) IsPositive() bool {
	return p.years > 0 || p.months > 0
}

// GreaterThan はこの期間が他の期間より大きいかどうかを返す
func (p Period) GreaterThan(other Period) bool {
	return p.ToMonths() > other.ToMonths()
}

// LessThan はこの期間が他の期間より小さいかどうかを返す
func (p Period) LessThan(other Period) bool {
	return p.ToMonths() < other.ToMonths()
}

// Equal はこの期間が他の期間と等しいかどうかを返す
func (p Period) Equal(other Period) bool {
	return p.years == other.years && p.months == other.months
}

// String は期間の文字列表現を返す
func (p Period) String() string {
	if p.years == 0 && p.months == 0 {
		return "0ヶ月"
	}

	if p.years == 0 {
		return fmt.Sprintf("%dヶ月", p.months)
	}

	if p.months == 0 {
		return fmt.Sprintf("%d年", p.years)
	}

	return fmt.Sprintf("%d年%dヶ月", p.years, p.months)
}

// AddMonths は指定された月数を期間に加算する
func (p Period) AddMonths(months int) (Period, error) {
	if months < 0 {
		return Period{}, errors.New("加算する月数は負の値にできません")
	}

	return NewPeriod(p.years, p.months+months)
}

// AddYears は指定された年数を期間に加算する
func (p Period) AddYears(years int) (Period, error) {
	if years < 0 {
		return Period{}, errors.New("加算する年数は負の値にできません")
	}

	return NewPeriod(p.years+years, p.months)
}

// RemainingMonthsInYear は現在の年を完了するまでの残り月数を返す
func (p Period) RemainingMonthsInYear() int {
	if p.months == 0 {
		return 0
	}
	return 12 - p.months
}
