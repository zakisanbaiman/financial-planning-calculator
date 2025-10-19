# データベーススキーマ仕様

## 概要

財務計画計算機のPostgreSQLデータベーススキーマ設計書です。

## 設計原則

1. **データ整合性**: 外部キー制約とCHECK制約による厳密なデータ検証
2. **正規化**: 第3正規形に準拠したテーブル設計
3. **パフォーマンス**: 適切なインデックス設計による高速クエリ
4. **拡張性**: 将来の機能追加に対応できる柔軟な構造
5. **監査**: 作成日時・更新日時の自動記録

## テーブル詳細

### 1. users テーブル

ユーザーの基本情報を管理します。

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**カラム説明:**
- `id`: ユーザーの一意識別子（UUID）
- `email`: ユーザーのメールアドレス（ログイン用）
- `created_at`: レコード作成日時
- `updated_at`: レコード更新日時（トリガーで自動更新）

**制約:**
- PRIMARY KEY: `id`
- UNIQUE: `email`
- NOT NULL: `email`

**インデックス:**
- `idx_users_email`: メールアドレスでの高速検索

### 2. financial_data テーブル

ユーザーの基本的な財務情報を管理します。

```sql
CREATE TABLE financial_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    monthly_income DECIMAL(15,2) NOT NULL CHECK (monthly_income >= 0),
    investment_return DECIMAL(5,2) NOT NULL CHECK (investment_return >= 0 AND investment_return <= 100),
    inflation_rate DECIMAL(5,2) NOT NULL CHECK (inflation_rate >= 0 AND inflation_rate <= 50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_user_financial_data UNIQUE (user_id)
);
```

**カラム説明:**
- `id`: 財務データの一意識別子
- `user_id`: 所有者ユーザーのID
- `monthly_income`: 月収（税込み、円）
- `investment_return`: 期待投資利回り（年率%）
- `inflation_rate`: インフレ率（年率%）

**制約:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users(id)`
- UNIQUE: `user_id`（1ユーザー1財務データ）
- CHECK: 月収は非負値
- CHECK: 投資利回りは0-100%
- CHECK: インフレ率は0-50%

### 3. expense_items テーブル

月間支出の詳細項目を管理します。

```sql
CREATE TABLE expense_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    financial_data_id UUID NOT NULL REFERENCES financial_data(id) ON DELETE CASCADE,
    category VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**カラム説明:**
- `id`: 支出項目の一意識別子
- `financial_data_id`: 関連する財務データのID
- `category`: 支出カテゴリ（住居費、食費、交通費など）
- `amount`: 支出金額（円）
- `description`: 支出の詳細説明（オプション）

**制約:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `financial_data_id` → `financial_data(id)`
- CHECK: 金額は正の値

**インデックス:**
- `idx_expense_items_financial_data_id`: 財務データIDでの検索
- `idx_expense_items_category`: カテゴリでの検索

### 4. savings_items テーブル

現在の貯蓄・投資の詳細を管理します。

```sql
CREATE TABLE savings_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    financial_data_id UUID NOT NULL REFERENCES financial_data(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('deposit', 'investment', 'other')),
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**カラム説明:**
- `id`: 貯蓄項目の一意識別子
- `financial_data_id`: 関連する財務データのID
- `type`: 貯蓄タイプ（deposit: 預金, investment: 投資, other: その他）
- `amount`: 貯蓄金額（円）
- `description`: 貯蓄の詳細説明（オプション）

**制約:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `financial_data_id` → `financial_data(id)`
- CHECK: タイプは指定された値のみ
- CHECK: 金額は非負値

**インデックス:**
- `idx_savings_items_financial_data_id`: 財務データIDでの検索
- `idx_savings_items_type`: タイプでの検索

### 5. retirement_data テーブル

退職・年金に関する情報を管理します。

```sql
CREATE TABLE retirement_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_age INTEGER NOT NULL CHECK (current_age > 0 AND current_age <= 120),
    retirement_age INTEGER NOT NULL CHECK (retirement_age > 0 AND retirement_age <= 120),
    life_expectancy INTEGER NOT NULL CHECK (life_expectancy > 0 AND life_expectancy <= 120),
    monthly_retirement_expenses DECIMAL(15,2) NOT NULL CHECK (monthly_retirement_expenses >= 0),
    pension_amount DECIMAL(15,2) NOT NULL CHECK (pension_amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_user_retirement_data UNIQUE (user_id),
    CONSTRAINT valid_retirement_age CHECK (retirement_age > current_age),
    CONSTRAINT valid_life_expectancy CHECK (life_expectancy > retirement_age)
);
```

**カラム説明:**
- `id`: 退職データの一意識別子
- `user_id`: 所有者ユーザーのID
- `current_age`: 現在の年齢
- `retirement_age`: 退職予定年齢
- `life_expectancy`: 平均寿命
- `monthly_retirement_expenses`: 退職後の月間生活費（円）
- `pension_amount`: 月間年金受給予定額（円）

**制約:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users(id)`
- UNIQUE: `user_id`（1ユーザー1退職データ）
- CHECK: 年齢は1-120の範囲
- CHECK: 退職年齢 > 現在年齢
- CHECK: 平均寿命 > 退職年齢
- CHECK: 金額は非負値

### 6. goals テーブル

ユーザーの財務目標を管理します。

```sql
CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('savings', 'retirement', 'emergency', 'custom')),
    title VARCHAR(255) NOT NULL,
    target_amount DECIMAL(15,2) NOT NULL CHECK (target_amount > 0),
    target_date DATE NOT NULL,
    current_amount DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (current_amount >= 0),
    monthly_contribution DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (monthly_contribution >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_target_date CHECK (target_date > CURRENT_DATE)
);
```

**カラム説明:**
- `id`: 目標の一意識別子
- `user_id`: 所有者ユーザーのID
- `type`: 目標タイプ（savings: 貯蓄, retirement: 退職, emergency: 緊急資金, custom: カスタム）
- `title`: 目標のタイトル
- `target_amount`: 目標金額（円）
- `target_date`: 目標達成期日
- `current_amount`: 現在の達成額（円）
- `monthly_contribution`: 月間積立額（円）
- `is_active`: 目標がアクティブかどうか

**制約:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users(id)`
- CHECK: タイプは指定された値のみ
- CHECK: 目標金額は正の値
- CHECK: 現在額・月間積立額は非負値
- CHECK: 目標期日は未来の日付

**インデックス:**
- `idx_goals_user_id`: ユーザーIDでの検索
- `idx_goals_type`: タイプでの検索
- `idx_goals_is_active`: アクティブ状態での検索
- `idx_goals_target_date`: 目標期日での検索
- `idx_goals_user_active`: ユーザーのアクティブ目標検索（複合インデックス）

## 自動更新機能

### updated_at カラムの自動更新

全テーブルで `updated_at` カラムが自動更新されるよう、PostgreSQLトリガーを設定しています。

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 各テーブルにトリガーを設定
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## データ型の選択理由

### UUID vs SERIAL

- **UUID**: グローバルに一意、分散システムに適している
- **セキュリティ**: 連番でないため推測困難
- **マイグレーション**: 他システムとの統合が容易

### DECIMAL vs FLOAT

- **DECIMAL(15,2)**: 金額計算で精度が重要
- **DECIMAL(5,2)**: パーセンテージ値（99.99%まで対応）

### TIMESTAMP WITH TIME ZONE

- **タイムゾーン対応**: グローバル展開に備えた設計
- **精度**: マイクロ秒まで記録

## パフォーマンス最適化

### インデックス戦略

1. **主キーインデックス**: 自動作成（B-tree）
2. **外部キーインデックス**: JOIN性能向上
3. **検索頻度の高いカラム**: 単一カラムインデックス
4. **複合検索**: 複合インデックス
5. **条件付きインデックス**: 部分インデックスでサイズ削減

### クエリ最適化のガイドライン

1. **適切なWHERE句**: インデックスを活用
2. **必要なカラムのみSELECT**: ネットワーク負荷軽減
3. **JOINの順序**: 小さなテーブルから結合
4. **LIMIT句の活用**: 大量データの制限

## セキュリティ考慮事項

### データ保護

1. **外部キー制約**: 参照整合性の保証
2. **CHECK制約**: 不正データの防止
3. **NOT NULL制約**: 必須データの保証
4. **UNIQUE制約**: 重複データの防止

### アクセス制御

```sql
-- アプリケーション用ユーザーの作成
CREATE USER financial_app WITH PASSWORD 'secure_password';

-- 必要最小限の権限付与
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO financial_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO financial_app;
```

## 運用・監視

### 定期メンテナンス

```sql
-- 統計情報の更新
ANALYZE;

-- インデックスの再構築（必要に応じて）
REINDEX INDEX idx_goals_user_active;

-- 不要データの削除（論理削除の場合）
DELETE FROM goals WHERE is_active = false AND updated_at < NOW() - INTERVAL '1 year';
```

### 監視クエリ

```sql
-- テーブルサイズの確認
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- インデックス使用状況
SELECT 
    indexrelname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

## 今後の拡張予定

### Phase 2 機能

1. **user_preferences テーブル**: UI設定・通知設定
2. **calculation_history テーブル**: 計算履歴の保存
3. **notifications テーブル**: 目標達成通知
4. **shared_plans テーブル**: 家族間での計画共有

### パフォーマンス改善

1. **パーティショニング**: 大量データ対応
2. **マテリアライズドビュー**: 複雑な集計の高速化
3. **読み取り専用レプリカ**: 読み取り性能向上