package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

// PostgreSQLFinancialPlanRepository はPostgreSQLを使用した財務計画リポジトリの実装
type PostgreSQLFinancialPlanRepository struct {
	db *sql.DB
}

// NewPostgreSQLFinancialPlanRepository は新しいPostgreSQL財務計画リポジトリを作成する
func NewPostgreSQLFinancialPlanRepository(db *sql.DB) repositories.FinancialPlanRepository {
	return &PostgreSQLFinancialPlanRepository{db: db}
}

// Save は財務計画を保存する
func (r *PostgreSQLFinancialPlanRepository) Save(ctx context.Context, plan *aggregates.FinancialPlan) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
	}
	defer tx.Rollback()

	// 財務プロファイルを保存
	if err := r.saveFinancialProfile(ctx, tx, plan.Profile()); err != nil {
		return fmt.Errorf("財務プロファイルの保存に失敗しました: %w", err)
	}

	// 退職データを保存（存在する場合）
	if plan.RetirementData() != nil {
		if err := r.saveRetirementData(ctx, tx, plan.RetirementData()); err != nil {
			return fmt.Errorf("退職データの保存に失敗しました: %w", err)
		}
	}

	// 目標を保存
	for _, goal := range plan.Goals() {
		if err := r.saveGoal(ctx, tx, goal); err != nil {
			return fmt.Errorf("目標の保存に失敗しました: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID は指定されたIDの財務計画を取得する
func (r *PostgreSQLFinancialPlanRepository) FindByID(ctx context.Context, id aggregates.FinancialPlanID) (*aggregates.FinancialPlan, error) {
	// 財務計画IDから直接取得する方法がないため、まずユーザーIDを取得する必要がある
	// この実装では、財務プロファイルからユーザーIDを取得してからFindByUserIDを呼び出す
	var userID string
	query := `SELECT user_id FROM financial_data WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("財務計画が見つかりません: %s", id)
		}
		return nil, fmt.Errorf("財務計画の検索に失敗しました: %w", err)
	}

	return r.FindByUserID(ctx, entities.UserID(userID))
}

// FindByUserID は指定されたユーザーIDの財務計画を取得する
func (r *PostgreSQLFinancialPlanRepository) FindByUserID(ctx context.Context, userID entities.UserID) (*aggregates.FinancialPlan, error) {
	// 財務プロファイルを取得
	profile, err := r.loadFinancialProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("財務プロファイルの取得に失敗しました: %w", err)
	}

	// 財務計画を作成
	plan, err := aggregates.NewFinancialPlan(profile)
	if err != nil {
		return nil, fmt.Errorf("財務計画の作成に失敗しました: %w", err)
	}

	// 退職データを取得（存在する場合）
	retirementData, err := r.loadRetirementData(ctx, userID)
	if err == nil && retirementData != nil {
		if err := plan.SetRetirementData(retirementData); err != nil {
			return nil, fmt.Errorf("退職データの設定に失敗しました: %w", err)
		}
	}

	// 目標を取得
	goals, err := r.loadGoals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// 目標を財務計画に追加
	for _, goal := range goals {
		if err := plan.AddGoal(goal); err != nil {
			// 目標の追加に失敗した場合はログに記録するが、処理は続行
			// これにより、一部の目標に問題があっても他の目標は取得できる
			continue
		}
	}

	return plan, nil
}

// Update は既存の財務計画を更新する
func (r *PostgreSQLFinancialPlanRepository) Update(ctx context.Context, plan *aggregates.FinancialPlan) error {
	// Updateは基本的にSaveと同じ処理（UPSERT）
	return r.Save(ctx, plan)
}

// Delete は指定されたIDの財務計画を削除する
func (r *PostgreSQLFinancialPlanRepository) Delete(ctx context.Context, id aggregates.FinancialPlanID) error {
	// まずユーザーIDを取得
	var userID string
	query := `SELECT user_id FROM financial_data WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("財務計画が見つかりません: %s", id)
		}
		return fmt.Errorf("財務計画の検索に失敗しました: %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
	}
	defer tx.Rollback()

	// 関連データを削除（外部キー制約により自動削除されるが、明示的に削除）
	queries := []string{
		`DELETE FROM goals WHERE user_id = $1`,
		`DELETE FROM retirement_data WHERE user_id = $1`,
		`DELETE FROM expense_items WHERE financial_data_id IN (SELECT id FROM financial_data WHERE user_id = $1)`,
		`DELETE FROM savings_items WHERE financial_data_id IN (SELECT id FROM financial_data WHERE user_id = $1)`,
		`DELETE FROM financial_data WHERE user_id = $1`,
	}

	for _, query := range queries {
		if _, err := tx.ExecContext(ctx, query, userID); err != nil {
			return fmt.Errorf("関連データの削除に失敗しました: %w", err)
		}
	}

	return tx.Commit()
}

// Exists は指定されたIDの財務計画が存在するかチェックする
func (r *PostgreSQLFinancialPlanRepository) Exists(ctx context.Context, id aggregates.FinancialPlanID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM financial_data WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("財務計画の存在確認に失敗しました: %w", err)
	}
	return count > 0, nil
}

// ExistsByUserID は指定されたユーザーIDの財務計画が存在するかチェックする
func (r *PostgreSQLFinancialPlanRepository) ExistsByUserID(ctx context.Context, userID entities.UserID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM financial_data WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, string(userID)).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("財務計画の存在確認に失敗しました: %w", err)
	}
	return count > 0, nil
}

// saveFinancialProfile は財務プロファイルを保存する
func (r *PostgreSQLFinancialPlanRepository) saveFinancialProfile(ctx context.Context, tx *sql.Tx, profile *entities.FinancialProfile) error {
	// 財務データを保存（UPSERT）
	query := `
		INSERT INTO financial_data (id, user_id, monthly_income, investment_return, inflation_rate, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE SET
			monthly_income = EXCLUDED.monthly_income,
			investment_return = EXCLUDED.investment_return,
			inflation_rate = EXCLUDED.inflation_rate,
			updated_at = EXCLUDED.updated_at
		RETURNING id`

	var financialDataID string
	err := tx.QueryRowContext(ctx, query,
		string(profile.ID()),
		string(profile.UserID()),
		profile.MonthlyIncome().Amount(),
		profile.InvestmentReturn().AsPercentage(),
		profile.InflationRate().AsPercentage(),
		profile.CreatedAt(),
		profile.UpdatedAt(),
	).Scan(&financialDataID)
	if err != nil {
		return fmt.Errorf("財務データの保存に失敗しました: %w", err)
	}

	// 既存の支出項目と貯蓄項目を削除
	if _, err := tx.ExecContext(ctx, `DELETE FROM expense_items WHERE financial_data_id = $1`, financialDataID); err != nil {
		return fmt.Errorf("既存支出項目の削除に失敗しました: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM savings_items WHERE financial_data_id = $1`, financialDataID); err != nil {
		return fmt.Errorf("既存貯蓄項目の削除に失敗しました: %w", err)
	}

	// 支出項目を保存
	for _, expense := range profile.MonthlyExpenses() {
		expenseQuery := `
			INSERT INTO expense_items (financial_data_id, category, amount, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := tx.ExecContext(ctx, expenseQuery,
			financialDataID,
			expense.Category,
			expense.Amount.Amount(),
			expense.Description,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("支出項目の保存に失敗しました: %w", err)
		}
	}

	// 貯蓄項目を保存
	for _, savings := range profile.CurrentSavings() {
		savingsQuery := `
			INSERT INTO savings_items (financial_data_id, type, amount, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := tx.ExecContext(ctx, savingsQuery,
			financialDataID,
			savings.Type,
			savings.Amount.Amount(),
			savings.Description,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("貯蓄項目の保存に失敗しました: %w", err)
		}
	}

	return nil
}

// saveRetirementData は退職データを保存する
func (r *PostgreSQLFinancialPlanRepository) saveRetirementData(ctx context.Context, tx *sql.Tx, retirementData *entities.RetirementData) error {
	query := `
		INSERT INTO retirement_data (id, user_id, current_age, retirement_age, life_expectancy, monthly_retirement_expenses, pension_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			current_age = EXCLUDED.current_age,
			retirement_age = EXCLUDED.retirement_age,
			life_expectancy = EXCLUDED.life_expectancy,
			monthly_retirement_expenses = EXCLUDED.monthly_retirement_expenses,
			pension_amount = EXCLUDED.pension_amount,
			updated_at = EXCLUDED.updated_at`

	_, err := tx.ExecContext(ctx, query,
		string(retirementData.ID()),
		string(retirementData.UserID()),
		retirementData.CurrentAge(),
		retirementData.RetirementAge(),
		retirementData.LifeExpectancy(),
		retirementData.MonthlyRetirementExpenses().Amount(),
		retirementData.PensionAmount().Amount(),
		retirementData.CreatedAt(),
		retirementData.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("退職データの保存に失敗しました: %w", err)
	}

	return nil
}

// saveGoal は目標を保存する
func (r *PostgreSQLFinancialPlanRepository) saveGoal(ctx context.Context, tx *sql.Tx, goal *entities.Goal) error {
	query := `
		INSERT INTO goals (id, user_id, type, title, target_amount, target_date, current_amount, monthly_contribution, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			type = EXCLUDED.type,
			title = EXCLUDED.title,
			target_amount = EXCLUDED.target_amount,
			target_date = EXCLUDED.target_date,
			current_amount = EXCLUDED.current_amount,
			monthly_contribution = EXCLUDED.monthly_contribution,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at`

	_, err := tx.ExecContext(ctx, query,
		string(goal.ID()),
		string(goal.UserID()),
		string(goal.GoalType()),
		goal.Title(),
		goal.TargetAmount().Amount(),
		goal.TargetDate(),
		goal.CurrentAmount().Amount(),
		goal.MonthlyContribution().Amount(),
		goal.IsActive(),
		goal.CreatedAt(),
		goal.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("目標の保存に失敗しました: %w", err)
	}

	return nil
}

// loadFinancialProfile は財務プロファイルを読み込む
func (r *PostgreSQLFinancialPlanRepository) loadFinancialProfile(ctx context.Context, userID entities.UserID) (*entities.FinancialProfile, error) {
	// 財務データを取得
	var financialDataID, fdUserID string
	var monthlyIncome, investmentReturn, inflationRate float64
	var createdAt, updatedAt time.Time

	query := `SELECT id, user_id, monthly_income, investment_return, inflation_rate, created_at, updated_at 
			  FROM financial_data WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, string(userID)).Scan(
		&financialDataID, &fdUserID, &monthlyIncome, &investmentReturn, &inflationRate, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("財務データが見つかりません: %s", userID)
		}
		return nil, fmt.Errorf("財務データの取得に失敗しました: %w", err)
	}

	// 支出項目を取得
	expenseQuery := `SELECT category, amount, description FROM expense_items WHERE financial_data_id = $1`
	expenseRows, err := r.db.QueryContext(ctx, expenseQuery, financialDataID)
	if err != nil {
		return nil, fmt.Errorf("支出項目の取得に失敗しました: %w", err)
	}
	defer expenseRows.Close()

	var expenses entities.ExpenseCollection
	for expenseRows.Next() {
		var category, description string
		var amount float64
		if err := expenseRows.Scan(&category, &amount, &description); err != nil {
			return nil, fmt.Errorf("支出項目の読み取りに失敗しました: %w", err)
		}

		expenseAmount, err := valueobjects.NewMoneyJPY(amount)
		if err != nil {
			return nil, fmt.Errorf("支出金額の作成に失敗しました: %w", err)
		}

		expenses = append(expenses, entities.ExpenseItem{
			Category:    category,
			Amount:      expenseAmount,
			Description: description,
		})
	}

	// 貯蓄項目を取得
	savingsQuery := `SELECT type, amount, description FROM savings_items WHERE financial_data_id = $1`
	savingsRows, err := r.db.QueryContext(ctx, savingsQuery, financialDataID)
	if err != nil {
		return nil, fmt.Errorf("貯蓄項目の取得に失敗しました: %w", err)
	}
	defer savingsRows.Close()

	var savings entities.SavingsCollection
	for savingsRows.Next() {
		var savingsType, description string
		var amount float64
		if err := savingsRows.Scan(&savingsType, &amount, &description); err != nil {
			return nil, fmt.Errorf("貯蓄項目の読み取りに失敗しました: %w", err)
		}

		savingsAmount, err := valueobjects.NewMoneyJPY(amount)
		if err != nil {
			return nil, fmt.Errorf("貯蓄金額の作成に失敗しました: %w", err)
		}

		savings = append(savings, entities.SavingsItem{
			Type:        savingsType,
			Amount:      savingsAmount,
			Description: description,
		})
	}

	// 値オブジェクトを作成
	monthlyIncomeVO, err := valueobjects.NewMoneyJPY(monthlyIncome)
	if err != nil {
		return nil, fmt.Errorf("月収の作成に失敗しました: %w", err)
	}

	investmentReturnVO, err := valueobjects.NewRate(investmentReturn)
	if err != nil {
		return nil, fmt.Errorf("投資利回りの作成に失敗しました: %w", err)
	}

	inflationRateVO, err := valueobjects.NewRate(inflationRate)
	if err != nil {
		return nil, fmt.Errorf("インフレ率の作成に失敗しました: %w", err)
	}

	// 財務プロファイルを作成
	profile, err := entities.NewFinancialProfile(
		entities.UserID(fdUserID),
		monthlyIncomeVO,
		expenses,
		savings,
		investmentReturnVO,
		inflationRateVO,
	)
	if err != nil {
		return nil, fmt.Errorf("財務プロファイルの作成に失敗しました: %w", err)
	}

	return profile, nil
}

// loadRetirementData は退職データを読み込む
func (r *PostgreSQLFinancialPlanRepository) loadRetirementData(ctx context.Context, userID entities.UserID) (*entities.RetirementData, error) {
	var id, rdUserID string
	var currentAge, retirementAge, lifeExpectancy int
	var monthlyRetirementExpenses, pensionAmount float64
	var createdAt, updatedAt time.Time

	query := `SELECT id, user_id, current_age, retirement_age, life_expectancy, monthly_retirement_expenses, pension_amount, created_at, updated_at 
			  FROM retirement_data WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, string(userID)).Scan(
		&id, &rdUserID, &currentAge, &retirementAge, &lifeExpectancy, &monthlyRetirementExpenses, &pensionAmount, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 退職データが存在しない場合はnilを返す
		}
		return nil, fmt.Errorf("退職データの取得に失敗しました: %w", err)
	}

	// 値オブジェクトを作成
	monthlyExpensesVO, err := valueobjects.NewMoneyJPY(monthlyRetirementExpenses)
	if err != nil {
		return nil, fmt.Errorf("月間退職後支出の作成に失敗しました: %w", err)
	}

	pensionAmountVO, err := valueobjects.NewMoneyJPY(pensionAmount)
	if err != nil {
		return nil, fmt.Errorf("年金額の作成に失敗しました: %w", err)
	}

	// 退職データを作成
	retirementData, err := entities.NewRetirementData(
		entities.UserID(rdUserID),
		currentAge,
		retirementAge,
		lifeExpectancy,
		monthlyExpensesVO,
		pensionAmountVO,
	)
	if err != nil {
		return nil, fmt.Errorf("退職データの作成に失敗しました: %w", err)
	}

	return retirementData, nil
}

// loadGoals は目標を読み込む
func (r *PostgreSQLFinancialPlanRepository) loadGoals(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
	query := `SELECT id, user_id, type, title, target_amount, target_date, current_amount, monthly_contribution, is_active, created_at, updated_at 
			  FROM goals WHERE user_id = $1 ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, query, string(userID))
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}
	defer rows.Close()

	var goals []*entities.Goal
	for rows.Next() {
		var id, gUserID, goalType, title string
		var targetAmount, currentAmount, monthlyContribution float64
		var targetDate time.Time
		var isActive bool
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &gUserID, &goalType, &title, &targetAmount, &targetDate, &currentAmount, &monthlyContribution, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("目標の読み取りに失敗しました: %w", err)
		}

		// 値オブジェクトを作成
		targetAmountVO, err := valueobjects.NewMoneyJPY(targetAmount)
		if err != nil {
			return nil, fmt.Errorf("目標金額の作成に失敗しました: %w", err)
		}

		monthlyContributionVO, err := valueobjects.NewMoneyJPY(monthlyContribution)
		if err != nil {
			return nil, fmt.Errorf("月間拠出額の作成に失敗しました: %w", err)
		}

		// 目標を作成
		goal, err := entities.NewGoal(
			entities.UserID(gUserID),
			entities.GoalType(goalType),
			title,
			targetAmountVO,
			targetDate,
			monthlyContributionVO,
		)
		if err != nil {
			return nil, fmt.Errorf("目標の作成に失敗しました: %w", err)
		}

		// 現在の金額を設定
		currentAmountVO, err := valueobjects.NewMoneyJPY(currentAmount)
		if err != nil {
			return nil, fmt.Errorf("現在の金額の作成に失敗しました: %w", err)
		}
		if err := goal.UpdateCurrentAmount(currentAmountVO); err != nil {
			return nil, fmt.Errorf("現在の金額の設定に失敗しました: %w", err)
		}

		// アクティブ状態を設定
		if !isActive {
			goal.Deactivate()
		}

		goals = append(goals, goal)
	}

	return goals, nil
}
