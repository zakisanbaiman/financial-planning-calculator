# Repository Implementation

このディレクトリには、財務計画アプリケーションのリポジトリ実装が含まれています。

## 概要

リポジトリパターンを使用して、ドメインエンティティの永続化を抽象化しています。PostgreSQLを使用した具象実装を提供し、将来的に他のデータベースへの切り替えも可能な設計になっています。

## アーキテクチャ

```
domain/repositories/          # リポジトリインターフェース
├── financial_plan_repository.go
└── goal_repository.go

infrastructure/repositories/  # 具象実装
├── postgresql_financial_plan_repository.go
├── postgresql_goal_repository.go
├── repository_factory.go
└── example_usage.go
```

## 実装されたリポジトリ

### 1. FinancialPlanRepository

財務計画の永続化を担当するリポジトリです。

**主要メソッド:**
- `Save(ctx, plan)` - 財務計画を保存
- `FindByUserID(ctx, userID)` - ユーザーIDで財務計画を取得
- `Update(ctx, plan)` - 財務計画を更新
- `Delete(ctx, id)` - 財務計画を削除
- `ExistsByUserID(ctx, userID)` - 財務計画の存在確認

**特徴:**
- 財務プロファイル、退職データ、目標を含む複合エンティティの管理
- トランザクションを使用した整合性保証
- 関連データの一括保存・取得

### 2. GoalRepository

目標の永続化を担当するリポジトリです。

**主要メソッド:**
- `Save(ctx, goal)` - 目標を保存
- `FindByID(ctx, id)` - IDで目標を取得
- `FindByUserID(ctx, userID)` - ユーザーの全目標を取得
- `FindActiveGoalsByUserID(ctx, userID)` - アクティブな目標のみ取得
- `FindByUserIDAndType(ctx, userID, type)` - タイプ別目標取得
- `Update(ctx, goal)` - 目標を更新
- `Delete(ctx, id)` - 目標を削除
- `CountActiveGoalsByType(ctx, userID, type)` - タイプ別アクティブ目標数

**特徴:**
- 効率的な検索のためのインデックス活用
- 目標タイプ別の統計機能
- アクティブ/非アクティブ状態の管理

## 使用方法

### 基本的な使用例

```go
package main

import (
    "context"
    "log"
    
    "github.com/financial-planning-calculator/backend/config"
    "github.com/financial-planning-calculator/backend/infrastructure/repositories"
)

func main() {
    // データベース接続
    dbConfig := config.NewDatabaseConfig()
    db, err := config.NewDatabaseConnection(dbConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // リポジトリファクトリー作成
    factory := repositories.NewRepositoryFactory(db)
    goalRepo := factory.NewGoalRepository()
    
    ctx := context.Background()
    
    // 目標を検索
    goals, err := goalRepo.FindByUserID(ctx, userID)
    if err != nil {
        log.Fatal(err)
    }
    
    // 処理...
}
```

### 財務計画の操作

```go
// 財務計画リポジトリを取得
planRepo := factory.NewFinancialPlanRepository()

// 財務計画を作成・保存
plan, err := aggregates.NewFinancialPlan(profile)
if err != nil {
    return err
}

err = planRepo.Save(ctx, plan)
if err != nil {
    return err
}

// 財務計画を取得
plan, err = planRepo.FindByUserID(ctx, userID)
if err != nil {
    return err
}
```

## データベーススキーマ

### 主要テーブル

1. **financial_data** - 財務データ
   - 月収、投資利回り、インフレ率
   - 支出項目、貯蓄項目との関連

2. **goals** - 目標
   - 目標タイプ、金額、期日
   - 進捗管理、アクティブ状態

3. **retirement_data** - 退職データ
   - 退職年齢、生活費、年金額

### インデックス最適化

効率的な検索のために以下のインデックスが設定されています：

```sql
-- 基本インデックス
CREATE INDEX idx_goals_user_id ON goals(user_id);
CREATE INDEX idx_goals_type ON goals(type);
CREATE INDEX idx_goals_is_active ON goals(is_active);
CREATE INDEX idx_goals_target_date ON goals(target_date);

-- 複合インデックス
CREATE INDEX idx_goals_user_active ON goals(user_id, is_active) 
WHERE is_active = true;
```

## エラーハンドリング

### 一般的なエラーパターン

1. **データが見つからない場合**
   ```go
   if err == sql.ErrNoRows {
       return nil, fmt.Errorf("目標が見つかりません: %s", id)
   }
   ```

2. **制約違反**
   ```go
   if pqErr, ok := err.(*pq.Error); ok {
       switch pqErr.Code {
       case "23505": // UNIQUE制約違反
           return fmt.Errorf("重複データです")
       case "23503": // 外部キー制約違反
           return fmt.Errorf("関連データが存在しません")
       }
   }
   ```

3. **トランザクション管理**
   ```go
   tx, err := r.db.BeginTx(ctx, nil)
   if err != nil {
       return err
   }
   defer tx.Rollback() // エラー時は自動ロールバック
   
   // 処理...
   
   return tx.Commit()
   ```

## パフォーマンス最適化

### 1. 接続プール設定

```go
db.SetMaxOpenConns(25)    // 最大接続数
db.SetMaxIdleConns(5)     // アイドル接続数
db.SetConnMaxLifetime(5 * time.Minute) // 接続の最大生存時間
```

### 2. プリペアドステートメント

繰り返し実行されるクエリではプリペアドステートメントを使用：

```go
stmt, err := db.Prepare("SELECT * FROM goals WHERE user_id = $1")
defer stmt.Close()

rows, err := stmt.QueryContext(ctx, userID)
```

### 3. バッチ処理

複数レコードの処理時はトランザクションでまとめて実行：

```go
tx, err := db.BeginTx(ctx, nil)
for _, goal := range goals {
    _, err = tx.ExecContext(ctx, query, goal.params...)
}
tx.Commit()
```

## テスト

### テスト実行

```bash
# 全テスト実行
go test ./infrastructure/repositories -v

# 短縮モード（データベース接続なし）
go test ./infrastructure/repositories -short
```

### テストデータベース

テストでは実際のPostgreSQLデータベースを使用します。CI環境では自動的にスキップされます。

## 今後の拡張

1. **キャッシュ層の追加**
   - Redis等を使用した読み取り性能向上

2. **読み取り専用レプリカ対応**
   - 読み取りクエリの負荷分散

3. **メトリクス収集**
   - クエリ実行時間、エラー率の監視

4. **他データベース対応**
   - MySQL、SQLite等への対応

## 関連ドキュメント

- [データベーススキーマ](../../docs/database_schema.md)
- [ドメインモデル](../../domain/README.md)
- [API仕様](../../docs/swagger.yaml)