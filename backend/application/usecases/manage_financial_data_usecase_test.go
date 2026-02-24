package usecases

import (
"context"
"errors"
"testing"

"github.com/financial-planning-calculator/backend/domain/entities"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/mock"
"github.com/stretchr/testify/require"
)

func TestCreateFinancialPlan(t *testing.T) {
tests := []struct {
name        string
input       CreateFinancialPlanInput
setupMock   func(*MockFinancialPlanRepository)
expectError bool
errContains string
}{
{
name: "正常系: 財務計画作成",
input: CreateFinancialPlanInput{
UserID:           "user-001",
MonthlyIncome:    400000,
MonthlyExpenses:  []ExpenseItem{{Category: "住居費", Amount: 100000}},
CurrentSavings:   []SavingsItem{{Type: "deposit", Amount: 500000}},
InvestmentReturn: 5.0,
InflationRate:    2.0,
},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("ExistsByUserID", mock.Anything, entities.UserID("user-001")).Return(false, nil)
m.On("Save", mock.Anything, mock.AnythingOfType("*aggregates.FinancialPlan")).Return(nil)
},
expectError: false,
},
{
name: "異常系: 既存財務計画が存在する",
input: CreateFinancialPlanInput{
UserID:           "user-002",
MonthlyIncome:    400000,
MonthlyExpenses:  []ExpenseItem{{Category: "住居費", Amount: 100000}},
CurrentSavings:   []SavingsItem{{Type: "deposit", Amount: 500000}},
InvestmentReturn: 5.0,
InflationRate:    2.0,
},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("ExistsByUserID", mock.Anything, entities.UserID("user-002")).Return(true, nil)
},
expectError: true,
errContains: "ユーザーの財務計画は既に存在します",
},
{
name: "異常系: ExistsByUserID失敗",
input: CreateFinancialPlanInput{
UserID:           "user-003",
MonthlyIncome:    400000,
MonthlyExpenses:  []ExpenseItem{{Category: "住居費", Amount: 100000}},
CurrentSavings:   []SavingsItem{{Type: "deposit", Amount: 500000}},
InvestmentReturn: 5.0,
InflationRate:    2.0,
},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("ExistsByUserID", mock.Anything, entities.UserID("user-003")).Return(false, errors.New("DBエラー"))
},
expectError: true,
errContains: "既存財務計画の確認に失敗しました",
},
{
name: "異常系: 無効な貯蓄タイプ",
input: CreateFinancialPlanInput{
UserID:           "user-004",
MonthlyIncome:    400000,
MonthlyExpenses:  []ExpenseItem{{Category: "住居費", Amount: 100000}},
CurrentSavings:   []SavingsItem{{Type: "invalid_type", Amount: 500000}},
InvestmentReturn: 5.0,
InflationRate:    2.0,
},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("ExistsByUserID", mock.Anything, entities.UserID("user-004")).Return(false, nil)
},
expectError: true,
errContains: "財務プロファイルの作成に失敗しました",
},
{
name: "正常系: 退職データ付き",
input: CreateFinancialPlanInput{
UserID:                    "user-005",
MonthlyIncome:             500000,
MonthlyExpenses:           []ExpenseItem{{Category: "住居費", Amount: 150000}},
CurrentSavings:            []SavingsItem{{Type: "investment", Amount: 2000000}},
InvestmentReturn:          5.0,
InflationRate:             2.0,
RetirementAge:             func() *int { v := 65; return &v }(),
MonthlyRetirementExpenses: func() *float64 { v := 250000.0; return &v }(),
PensionAmount:             func() *float64 { v := 100000.0; return &v }(),
},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("ExistsByUserID", mock.Anything, entities.UserID("user-005")).Return(false, nil)
m.On("Save", mock.Anything, mock.AnythingOfType("*aggregates.FinancialPlan")).Return(nil)
},
expectError: false,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockRepo := new(MockFinancialPlanRepository)
tt.setupMock(mockRepo)
uc := NewManageFinancialDataUseCase(mockRepo)

output, err := uc.CreateFinancialPlan(context.Background(), tt.input)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
} else {
require.NoError(t, err)
assert.NotNil(t, output)
assert.Equal(t, string(tt.input.UserID), string(output.UserID))
assert.NotEmpty(t, output.PlanID)
}
mockRepo.AssertExpectations(t)
})
}
}

func TestGetFinancialPlan(t *testing.T) {
tests := []struct {
name        string
userID      string
setupMock   func(*MockFinancialPlanRepository)
expectError bool
errContains string
}{
{
name:   "正常系: 財務計画取得",
userID: "user-001",
setupMock: func(m *MockFinancialPlanRepository) {
plan := newTestFinancialPlan("user-001")
m.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
},
expectError: false,
},
{
name:   "異常系: 財務計画が存在しない",
userID: "user-002",
setupMock: func(m *MockFinancialPlanRepository) {
m.On("FindByUserID", mock.Anything, entities.UserID("user-002")).Return(nil, errors.New("財務データが見つかりません"))
},
expectError: true,
errContains: "財務計画の取得に失敗しました",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockRepo := new(MockFinancialPlanRepository)
tt.setupMock(mockRepo)
uc := NewManageFinancialDataUseCase(mockRepo)

output, err := uc.GetFinancialPlan(context.Background(), GetFinancialPlanInput{UserID: entities.UserID(tt.userID)})

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
} else {
require.NoError(t, err)
assert.NotNil(t, output)
assert.NotNil(t, output.Plan)
}
mockRepo.AssertExpectations(t)
})
}
}

func TestUpdateFinancialProfile(t *testing.T) {
tests := []struct {
name        string
input       UpdateFinancialProfileInput
setupMock   func(*MockFinancialPlanRepository)
expectError bool
errContains string
}{
{
name: "正常系: 財務プロファイル更新",
input: UpdateFinancialProfileInput{
UserID:           "user-001",
MonthlyIncome:    500000,
MonthlyExpenses:  []ExpenseItem{{Category: "住居費", Amount: 130000}},
CurrentSavings:   []SavingsItem{{Type: "deposit", Amount: 2000000}},
InvestmentReturn: 6.0,
InflationRate:    2.0,
},
setupMock: func(m *MockFinancialPlanRepository) {
plan := newTestFinancialPlan("user-001")
m.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
m.On("Update", mock.Anything, mock.AnythingOfType("*aggregates.FinancialPlan")).Return(nil)
},
expectError: false,
},
{
name: "異常系: 財務計画が存在しない",
input: UpdateFinancialProfileInput{
UserID:           "user-999",
MonthlyIncome:    400000,
MonthlyExpenses:  []ExpenseItem{{Category: "住居費", Amount: 100000}},
CurrentSavings:   []SavingsItem{{Type: "deposit", Amount: 500000}},
InvestmentReturn: 5.0,
InflationRate:    2.0,
},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("FindByUserID", mock.Anything, entities.UserID("user-999")).Return(nil, errors.New("財務データが見つかりません"))
},
expectError: true,
errContains: "財務計画の取得に失敗しました",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockRepo := new(MockFinancialPlanRepository)
tt.setupMock(mockRepo)
uc := NewManageFinancialDataUseCase(mockRepo)

output, err := uc.UpdateFinancialProfile(context.Background(), tt.input)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
} else {
require.NoError(t, err)
assert.NotNil(t, output)
}
mockRepo.AssertExpectations(t)
})
}
}

func TestUpdateRetirementData(t *testing.T) {
tests := []struct {
name        string
input       UpdateRetirementDataInput
setupMock   func(*MockFinancialPlanRepository)
expectError bool
errContains string
}{
{
name: "正常系: 退職データ更新",
input: UpdateRetirementDataInput{
UserID:                    "user-001",
RetirementAge:             65,
MonthlyRetirementExpenses: 250000,
PensionAmount:             100000,
},
setupMock: func(m *MockFinancialPlanRepository) {
plan := newTestFinancialPlan("user-001")
m.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
m.On("Update", mock.Anything, mock.AnythingOfType("*aggregates.FinancialPlan")).Return(nil)
},
expectError: false,
},
{
name: "異常系: 財務計画が存在しない",
input: UpdateRetirementDataInput{
UserID:                    "user-999",
RetirementAge:             65,
MonthlyRetirementExpenses: 250000,
PensionAmount:             100000,
},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("FindByUserID", mock.Anything, entities.UserID("user-999")).Return(nil, errors.New("財務データが見つかりません"))
},
expectError: true,
errContains: "財務計画の取得に失敗しました",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockRepo := new(MockFinancialPlanRepository)
tt.setupMock(mockRepo)
uc := NewManageFinancialDataUseCase(mockRepo)

output, err := uc.UpdateRetirementData(context.Background(), tt.input)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
} else {
require.NoError(t, err)
assert.NotNil(t, output)
}
mockRepo.AssertExpectations(t)
})
}
}

func TestUpdateEmergencyFund(t *testing.T) {
tests := []struct {
name        string
input       UpdateEmergencyFundInput
setupMock   func(*MockFinancialPlanRepository)
expectError bool
errContains string
}{
{
name: "正常系: 緊急資金設定更新",
input: UpdateEmergencyFundInput{
UserID:        "user-001",
TargetMonths:  6,
CurrentAmount: 500000,
},
setupMock: func(m *MockFinancialPlanRepository) {
plan := newTestFinancialPlan("user-001")
m.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
m.On("Update", mock.Anything, mock.AnythingOfType("*aggregates.FinancialPlan")).Return(nil)
},
expectError: false,
},
{
name: "異常系: 財務計画が存在しない",
input: UpdateEmergencyFundInput{
UserID:        "user-999",
TargetMonths:  6,
CurrentAmount: 500000,
},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("FindByUserID", mock.Anything, entities.UserID("user-999")).Return(nil, errors.New("財務データが見つかりません"))
},
expectError: true,
errContains: "財務計画の取得に失敗しました",
},
{
name: "異常系: 無効な目標月数（負の値）",
input: UpdateEmergencyFundInput{
UserID:        "user-001",
TargetMonths:  -1,
CurrentAmount: 500000,
},
setupMock: func(m *MockFinancialPlanRepository) {
plan := newTestFinancialPlan("user-001")
m.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
},
expectError: true,
errContains: "緊急資金設定の作成に失敗しました",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockRepo := new(MockFinancialPlanRepository)
tt.setupMock(mockRepo)
uc := NewManageFinancialDataUseCase(mockRepo)

output, err := uc.UpdateEmergencyFund(context.Background(), tt.input)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
} else {
require.NoError(t, err)
assert.NotNil(t, output)
}
mockRepo.AssertExpectations(t)
})
}
}

func TestDeleteFinancialPlan(t *testing.T) {
tests := []struct {
name        string
input       DeleteFinancialPlanInput
setupMock   func(*MockFinancialPlanRepository)
expectError bool
errContains string
}{
{
name:  "正常系: 財務計画削除",
input: DeleteFinancialPlanInput{UserID: "user-001"},
setupMock: func(m *MockFinancialPlanRepository) {
plan := newTestFinancialPlan("user-001")
m.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
m.On("Delete", mock.Anything, plan.ID()).Return(nil)
},
expectError: false,
},
{
name:  "異常系: 財務計画が存在しない",
input: DeleteFinancialPlanInput{UserID: "user-999"},
setupMock: func(m *MockFinancialPlanRepository) {
m.On("FindByUserID", mock.Anything, entities.UserID("user-999")).Return(nil, errors.New("財務データが見つかりません"))
},
expectError: true,
errContains: "財務計画の取得に失敗しました",
},
{
name:  "異常系: 削除処理失敗",
input: DeleteFinancialPlanInput{UserID: "user-001"},
setupMock: func(m *MockFinancialPlanRepository) {
plan := newTestFinancialPlan("user-001")
m.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
m.On("Delete", mock.Anything, plan.ID()).Return(errors.New("削除エラー"))
},
expectError: true,
errContains: "財務計画の削除に失敗しました",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockRepo := new(MockFinancialPlanRepository)
tt.setupMock(mockRepo)
uc := NewManageFinancialDataUseCase(mockRepo)

err := uc.DeleteFinancialPlan(context.Background(), tt.input)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
} else {
require.NoError(t, err)
}
mockRepo.AssertExpectations(t)
})
}
}


func TestCreateFinancialPlan_EmergencyFund(t *testing.T) {
targetMonths := 6
currentAmount := 600000.0

t.Run("正常系: 緊急資金設定あり", func(t *testing.T) {
mockRepo := new(MockFinancialPlanRepository)
mockRepo.On("ExistsByUserID", mock.Anything, entities.UserID("user-ef")).Return(false, nil)
mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*aggregates.FinancialPlan")).Return(nil)

uc := NewManageFinancialDataUseCase(mockRepo)
out, err := uc.CreateFinancialPlan(context.Background(), CreateFinancialPlanInput{
UserID:                     "user-ef",
MonthlyIncome:              400000,
MonthlyExpenses:            []ExpenseItem{{Category: "食費", Amount: 100000}},
CurrentSavings:             []SavingsItem{{Type: "deposit", Amount: 1000000}},
InvestmentReturn:           3.0,
InflationRate:              2.0,
EmergencyFundTargetMonths:  &targetMonths,
EmergencyFundCurrentAmount: &currentAmount,
})

require.NoError(t, err)
assert.NotNil(t, out)
mockRepo.AssertExpectations(t)
})
}
