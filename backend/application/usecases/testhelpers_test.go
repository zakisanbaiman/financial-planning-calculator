package usecases

import (
"context"
"time"

"github.com/financial-planning-calculator/backend/domain/aggregates"
"github.com/financial-planning-calculator/backend/domain/entities"
"github.com/financial-planning-calculator/backend/domain/valueobjects"
"github.com/stretchr/testify/mock"
)

// ─── MockFinancialPlanRepository ───────────────────────────────────────────

type MockFinancialPlanRepository struct {
mock.Mock
}

func (m *MockFinancialPlanRepository) Save(ctx context.Context, plan *aggregates.FinancialPlan) error {
args := m.Called(ctx, plan)
return args.Error(0)
}

func (m *MockFinancialPlanRepository) FindByID(ctx context.Context, id aggregates.FinancialPlanID) (*aggregates.FinancialPlan, error) {
args := m.Called(ctx, id)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*aggregates.FinancialPlan), args.Error(1)
}

func (m *MockFinancialPlanRepository) FindByUserID(ctx context.Context, userID entities.UserID) (*aggregates.FinancialPlan, error) {
args := m.Called(ctx, userID)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*aggregates.FinancialPlan), args.Error(1)
}

func (m *MockFinancialPlanRepository) Update(ctx context.Context, plan *aggregates.FinancialPlan) error {
args := m.Called(ctx, plan)
return args.Error(0)
}

func (m *MockFinancialPlanRepository) Delete(ctx context.Context, id aggregates.FinancialPlanID) error {
args := m.Called(ctx, id)
return args.Error(0)
}

func (m *MockFinancialPlanRepository) Exists(ctx context.Context, id aggregates.FinancialPlanID) (bool, error) {
args := m.Called(ctx, id)
return args.Bool(0), args.Error(1)
}

func (m *MockFinancialPlanRepository) ExistsByUserID(ctx context.Context, userID entities.UserID) (bool, error) {
args := m.Called(ctx, userID)
return args.Bool(0), args.Error(1)
}

// ─── MockGoalRepository ────────────────────────────────────────────────────

type MockGoalRepository struct {
mock.Mock
}

func (m *MockGoalRepository) Save(ctx context.Context, goal *entities.Goal) error {
args := m.Called(ctx, goal)
return args.Error(0)
}

func (m *MockGoalRepository) FindByID(ctx context.Context, id entities.GoalID) (*entities.Goal, error) {
args := m.Called(ctx, id)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*entities.Goal), args.Error(1)
}

func (m *MockGoalRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
args := m.Called(ctx, userID)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).([]*entities.Goal), args.Error(1)
}

func (m *MockGoalRepository) FindActiveGoalsByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
args := m.Called(ctx, userID)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).([]*entities.Goal), args.Error(1)
}

func (m *MockGoalRepository) FindByUserIDAndType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) ([]*entities.Goal, error) {
args := m.Called(ctx, userID, goalType)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).([]*entities.Goal), args.Error(1)
}

func (m *MockGoalRepository) Update(ctx context.Context, goal *entities.Goal) error {
args := m.Called(ctx, goal)
return args.Error(0)
}

func (m *MockGoalRepository) Delete(ctx context.Context, id entities.GoalID) error {
args := m.Called(ctx, id)
return args.Error(0)
}

func (m *MockGoalRepository) Exists(ctx context.Context, id entities.GoalID) (bool, error) {
args := m.Called(ctx, id)
return args.Bool(0), args.Error(1)
}

func (m *MockGoalRepository) CountActiveGoalsByType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) (int, error) {
args := m.Called(ctx, userID, goalType)
return args.Int(0), args.Error(1)
}

// ─── MockUserRepository ────────────────────────────────────────────────────

type MockUserRepository struct {
mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *entities.User) error {
args := m.Called(ctx, user)
return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
args := m.Called(ctx, id)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email entities.Email) (*entities.User, error) {
args := m.Called(ctx, email)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
args := m.Called(ctx, user)
return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id entities.UserID) error {
args := m.Called(ctx, id)
return args.Error(0)
}

func (m *MockUserRepository) Exists(ctx context.Context, id entities.UserID) (bool, error) {
args := m.Called(ctx, id)
return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email entities.Email) (bool, error) {
args := m.Called(ctx, email)
return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindByProviderUserID(ctx context.Context, provider entities.AuthProvider, providerUserID string) (*entities.User, error) {
args := m.Called(ctx, provider, providerUserID)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*entities.User), args.Error(1)
}

// ─── MockRefreshTokenRepository ────────────────────────────────────────────

type MockRefreshTokenRepository struct {
mock.Mock
}

func (m *MockRefreshTokenRepository) Save(ctx context.Context, token *entities.RefreshToken) error {
args := m.Called(ctx, token)
return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
args := m.Called(ctx, tokenHash)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*entities.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.RefreshToken, error) {
args := m.Called(ctx, userID)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).([]*entities.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Update(ctx context.Context, token *entities.RefreshToken) error {
args := m.Called(ctx, token)
return args.Error(0)
}

func (m *MockRefreshTokenRepository) Delete(ctx context.Context, id entities.RefreshTokenID) error {
args := m.Called(ctx, id)
return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID entities.UserID) error {
args := m.Called(ctx, userID)
return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
args := m.Called(ctx)
return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeByUserID(ctx context.Context, userID entities.UserID) error {
args := m.Called(ctx, userID)
return args.Error(0)
}

// ─── MockWebAuthnCredentialRepository ──────────────────────────────────────

type MockWebAuthnCredentialRepository struct {
mock.Mock
}

func (m *MockWebAuthnCredentialRepository) Save(ctx context.Context, credential *entities.WebAuthnCredential) error {
args := m.Called(ctx, credential)
return args.Error(0)
}

func (m *MockWebAuthnCredentialRepository) FindByID(ctx context.Context, id entities.CredentialID) (*entities.WebAuthnCredential, error) {
args := m.Called(ctx, id)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*entities.WebAuthnCredential), args.Error(1)
}

func (m *MockWebAuthnCredentialRepository) FindByCredentialID(ctx context.Context, credentialID []byte) (*entities.WebAuthnCredential, error) {
args := m.Called(ctx, credentialID)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).(*entities.WebAuthnCredential), args.Error(1)
}

func (m *MockWebAuthnCredentialRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.WebAuthnCredential, error) {
args := m.Called(ctx, userID)
if args.Get(0) == nil {
return nil, args.Error(1)
}
return args.Get(0).([]*entities.WebAuthnCredential), args.Error(1)
}

func (m *MockWebAuthnCredentialRepository) Update(ctx context.Context, credential *entities.WebAuthnCredential) error {
args := m.Called(ctx, credential)
return args.Error(0)
}

func (m *MockWebAuthnCredentialRepository) Delete(ctx context.Context, id entities.CredentialID) error {
args := m.Called(ctx, id)
return args.Error(0)
}

// ─── テスト用ヘルパー関数 ─────────────────────────────────────────────────

// mustCreateMoneyUsecase は指定金額のMoneyを作成する（エラー時パニック）
func mustCreateMoneyUsecase(amount float64) valueobjects.Money {
m, err := valueobjects.NewMoneyJPY(amount)
if err != nil {
panic(err)
}
return m
}

// newTestFinancialPlan はテスト用の財務計画を作成するヘルパー
func newTestFinancialPlan(userID string) *aggregates.FinancialPlan {
monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)
expenses := entities.ExpenseCollection{
{Category: "住居費", Amount: mustCreateMoneyUsecase(120000)},
{Category: "食費", Amount: mustCreateMoneyUsecase(60000)},
}
savings := entities.SavingsCollection{
{Type: "deposit", Amount: mustCreateMoneyUsecase(1000000)},
}
investmentReturn, _ := valueobjects.NewRate(5.0)
inflationRate, _ := valueobjects.NewRate(2.0)

profile, err := entities.NewFinancialProfile(
entities.UserID(userID),
monthlyIncome,
expenses,
savings,
investmentReturn,
inflationRate,
)
if err != nil {
panic(err)
}

plan, err := aggregates.NewFinancialPlan(profile)
if err != nil {
panic(err)
}
return plan
}

// newTestFinancialPlanWithGoal はテスト用の目標入り財務計画を作成するヘルパー
func newTestFinancialPlanWithGoal(userID string, goal *entities.Goal) *aggregates.FinancialPlan {
	plan := newTestFinancialPlan(userID)
	if err := plan.AddGoal(goal); err != nil {
		panic(err)
	}
	return plan
}

// newTestFinancialPlanWithRetirementData はテスト用の退職データ入り財務計画を作成するヘルパー
func newTestFinancialPlanWithRetirementData(userID string) *aggregates.FinancialPlan {
	plan := newTestFinancialPlan(userID)
	monthlyExpenses := mustCreateMoneyUsecase(200000)
	pension := mustCreateMoneyUsecase(80000)
	retirementData, err := entities.NewRetirementData(
		entities.UserID(userID),
		35,  // 現在年齢
		65,  // 退職年齢
		85,  // 平均寿命
		monthlyExpenses,
		pension,
	)
	if err != nil {
		panic(err)
	}
	if err := plan.SetRetirementData(retirementData); err != nil {
		panic(err)
	}
	return plan
}
func newTestGoal(userID string, goalType entities.GoalType) *entities.Goal {
targetAmount := mustCreateMoneyUsecase(1000000)
monthlyContribution := mustCreateMoneyUsecase(50000)
targetDate := time.Now().AddDate(2, 0, 0)

goal, err := entities.NewGoal(
entities.UserID(userID),
goalType,
"テスト目標",
targetAmount,
targetDate,
monthlyContribution,
)
if err != nil {
panic(err)
}
return goal
}
