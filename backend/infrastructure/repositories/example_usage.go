package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

// ExampleUsage demonstrates how to use the repository implementations
func ExampleUsage() {
	// データベース接続を設定
	dbConfig := config.NewDatabaseConfig()
	db, err := config.NewDatabaseConnection(dbConfig)
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}
	defer db.Close()

	// リポジトリファクトリーを作成
	factory := NewRepositoryFactory(db)
	goalRepo := factory.NewGoalRepository()
	financialPlanRepo := factory.NewFinancialPlanRepository()

	ctx := context.Background()

	// サンプルユーザーIDを作成（実際のアプリケーションではユーザー認証から取得）
	userID := entities.UserID("550e8400-e29b-41d4-a716-446655440000")

	// 1. 目標を作成して保存
	fmt.Println("=== 目標の作成と保存 ===")

	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)      // 100万円
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000) // 月5万円
	targetDate := time.Now().AddDate(2, 0, 0)                 // 2年後

	goal, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		"マイホーム頭金",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		log.Printf("目標作成エラー: %v", err)
		return
	}

	// 目標を保存
	if err := goalRepo.Save(ctx, goal); err != nil {
		log.Printf("目標保存エラー: %v", err)
		return
	}
	fmt.Printf("目標を保存しました: %s\n", goal.Title())

	// 2. 目標を検索
	fmt.Println("\n=== 目標の検索 ===")

	// IDで検索
	foundGoal, err := goalRepo.FindByID(ctx, goal.ID())
	if err != nil {
		log.Printf("目標検索エラー: %v", err)
		return
	}
	fmt.Printf("IDで検索した目標: %s (目標金額: %.0f円)\n",
		foundGoal.Title(), foundGoal.TargetAmount().Amount())

	// ユーザーIDで検索
	userGoals, err := goalRepo.FindByUserID(ctx, userID)
	if err != nil {
		log.Printf("ユーザー目標検索エラー: %v", err)
		return
	}
	fmt.Printf("ユーザーの目標数: %d\n", len(userGoals))

	// アクティブな目標のみ検索
	activeGoals, err := goalRepo.FindActiveGoalsByUserID(ctx, userID)
	if err != nil {
		log.Printf("アクティブ目標検索エラー: %v", err)
		return
	}
	fmt.Printf("アクティブな目標数: %d\n", len(activeGoals))

	// 3. 目標を更新
	fmt.Println("\n=== 目標の更新 ===")

	// 現在の金額を更新
	currentAmount, _ := valueobjects.NewMoneyJPY(150000) // 15万円貯まった
	if err := goal.UpdateCurrentAmount(currentAmount); err != nil {
		log.Printf("現在金額更新エラー: %v", err)
		return
	}

	// 進捗を計算
	progress, err := goal.CalculateProgress(currentAmount)
	if err != nil {
		log.Printf("進捗計算エラー: %v", err)
		return
	}
	fmt.Printf("目標進捗: %.1f%%\n", progress.AsPercentage())

	// データベースに更新を保存
	if err := goalRepo.Update(ctx, goal); err != nil {
		log.Printf("目標更新エラー: %v", err)
		return
	}
	fmt.Println("目標を更新しました")

	// 4. 目標タイプ別の統計
	fmt.Println("\n=== 統計情報 ===")

	savingsCount, err := goalRepo.CountActiveGoalsByType(ctx, userID, entities.GoalTypeSavings)
	if err != nil {
		log.Printf("統計取得エラー: %v", err)
		return
	}
	fmt.Printf("アクティブな貯蓄目標数: %d\n", savingsCount)

	// 5. 財務計画の存在確認
	fmt.Println("\n=== 財務計画の確認 ===")

	exists, err := financialPlanRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		log.Printf("財務計画存在確認エラー: %v", err)
		return
	}

	if exists {
		fmt.Println("財務計画が存在します")

		// 財務計画を取得
		plan, err := financialPlanRepo.FindByUserID(ctx, userID)
		if err != nil {
			log.Printf("財務計画取得エラー: %v", err)
			return
		}

		fmt.Printf("財務計画の目標数: %d\n", len(plan.Goals()))
	} else {
		fmt.Println("財務計画が存在しません")
	}

	fmt.Println("\n=== 使用例完了 ===")
}

// ExampleQueryOptimization demonstrates query optimization techniques
func ExampleQueryOptimization() {
	fmt.Println("=== クエリ最適化の例 ===")

	// 1. インデックスの活用
	fmt.Println("1. インデックスが設定されているカラム:")
	fmt.Println("   - goals.user_id (ユーザー別目標検索)")
	fmt.Println("   - goals.type (目標タイプ別検索)")
	fmt.Println("   - goals.is_active (アクティブ目標検索)")
	fmt.Println("   - goals.target_date (期限別検索)")
	fmt.Println("   - goals(user_id, is_active) (複合インデックス)")

	// 2. バッチ処理の例
	fmt.Println("\n2. バッチ処理:")
	fmt.Println("   - 複数目標の一括保存時はトランザクションを使用")
	fmt.Println("   - 財務計画保存時は関連データを一括処理")

	// 3. 接続プールの活用
	fmt.Println("\n3. 接続プール:")
	fmt.Println("   - sql.DBは内部的に接続プールを管理")
	fmt.Println("   - 適切なMaxOpenConns/MaxIdleConnsの設定が重要")

	// 4. プリペアドステートメントの活用
	fmt.Println("\n4. プリペアドステートメント:")
	fmt.Println("   - 繰り返し実行されるクエリで性能向上")
	fmt.Println("   - SQLインジェクション対策にも効果的")
}

// ExampleErrorHandling demonstrates error handling patterns
func ExampleErrorHandling() {
	fmt.Println("=== エラーハンドリングの例 ===")

	fmt.Println("1. 一般的なエラーパターン:")
	fmt.Println("   - sql.ErrNoRows: データが見つからない場合")
	fmt.Println("   - 制約違反: UNIQUE制約、外部キー制約など")
	fmt.Println("   - 接続エラー: ネットワーク問題、データベース停止など")

	fmt.Println("\n2. エラーハンドリング戦略:")
	fmt.Println("   - ドメインエラーとインフラエラーの分離")
	fmt.Println("   - 適切なエラーメッセージの提供")
	fmt.Println("   - ログ出力とエラー追跡")

	fmt.Println("\n3. トランザクション管理:")
	fmt.Println("   - 複数テーブル更新時の整合性保証")
	fmt.Println("   - エラー時の自動ロールバック")
	fmt.Println("   - デッドロック対策")
}
