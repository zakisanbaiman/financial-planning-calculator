package repositories

import (
	"database/sql"

	"github.com/financial-planning-calculator/backend/domain/repositories"
)

// RepositoryFactory はリポジトリのファクトリー
type RepositoryFactory struct {
	db *sql.DB
}

// NewRepositoryFactory は新しいリポジトリファクトリーを作成する
func NewRepositoryFactory(db *sql.DB) *RepositoryFactory {
	return &RepositoryFactory{db: db}
}

// NewFinancialPlanRepository は財務計画リポジトリを作成する
func (f *RepositoryFactory) NewFinancialPlanRepository() repositories.FinancialPlanRepository {
	return NewPostgreSQLFinancialPlanRepository(f.db)
}

// NewUserRepository はユーザーリポジトリを作成する
func (f *RepositoryFactory) NewUserRepository() repositories.UserRepository {
	return NewPostgreSQLUserRepository(f.db)
}

// NewGoalRepository は目標リポジトリを作成する
func (f *RepositoryFactory) NewGoalRepository() repositories.GoalRepository {
	return NewPostgreSQLGoalRepository(f.db)
}
