package usecases

import (
	"context"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/stretchr/testify/mock"
)

// mock_anything は任意の引数にマッチするtestifyのMatcherを返すヘルパー
func mock_anything() interface{} {
	return mock.MatchedBy(func(_ interface{}) bool { return true })
}

// -------------------------------------------------------------------
// MockFinancialPlanRepository
// -------------------------------------------------------------------

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

// -------------------------------------------------------------------
// MockGoalRepository
// -------------------------------------------------------------------

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

// -------------------------------------------------------------------
// MockUserRepository
// -------------------------------------------------------------------

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

// -------------------------------------------------------------------
// MockRefreshTokenRepository
// -------------------------------------------------------------------

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
