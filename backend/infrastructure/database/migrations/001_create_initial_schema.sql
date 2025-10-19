-- 001_create_initial_schema.sql
-- 財務計画アプリケーションの初期スキーマ作成

-- UUIDエクステンションを有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ユーザーテーブル（将来の拡張用）
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 財務データテーブル
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

-- 支出項目テーブル
CREATE TABLE expense_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    financial_data_id UUID NOT NULL REFERENCES financial_data(id) ON DELETE CASCADE,
    category VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 貯蓄項目テーブル
CREATE TABLE savings_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    financial_data_id UUID NOT NULL REFERENCES financial_data(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('deposit', 'investment', 'other')),
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 退職・年金情報テーブル
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

-- 目標テーブル
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

-- インデックス作成
-- ユーザーテーブル
CREATE INDEX idx_users_email ON users(email);

-- 財務データテーブル
CREATE INDEX idx_financial_data_user_id ON financial_data(user_id);

-- 支出項目テーブル
CREATE INDEX idx_expense_items_financial_data_id ON expense_items(financial_data_id);
CREATE INDEX idx_expense_items_category ON expense_items(category);

-- 貯蓄項目テーブル
CREATE INDEX idx_savings_items_financial_data_id ON savings_items(financial_data_id);
CREATE INDEX idx_savings_items_type ON savings_items(type);

-- 退職データテーブル
CREATE INDEX idx_retirement_data_user_id ON retirement_data(user_id);

-- 目標テーブル
CREATE INDEX idx_goals_user_id ON goals(user_id);
CREATE INDEX idx_goals_type ON goals(type);
CREATE INDEX idx_goals_is_active ON goals(is_active);
CREATE INDEX idx_goals_target_date ON goals(target_date);
CREATE INDEX idx_goals_user_active ON goals(user_id, is_active) WHERE is_active = true;

-- 更新日時自動更新のためのトリガー関数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 各テーブルに更新日時トリガーを設定
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_financial_data_updated_at BEFORE UPDATE ON financial_data
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_expense_items_updated_at BEFORE UPDATE ON expense_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_savings_items_updated_at BEFORE UPDATE ON savings_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_retirement_data_updated_at BEFORE UPDATE ON retirement_data
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_goals_updated_at BEFORE UPDATE ON goals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント追加
COMMENT ON TABLE users IS 'ユーザー情報テーブル';
COMMENT ON TABLE financial_data IS '財務データテーブル - ユーザーの基本的な収入・支出・投資情報';
COMMENT ON TABLE expense_items IS '支出項目テーブル - 月間支出の詳細';
COMMENT ON TABLE savings_items IS '貯蓄項目テーブル - 現在の貯蓄・投資の詳細';
COMMENT ON TABLE retirement_data IS '退職・年金情報テーブル - 老後計画に関する情報';
COMMENT ON TABLE goals IS '目標テーブル - ユーザーの財務目標';

COMMENT ON COLUMN financial_data.monthly_income IS '月収（税込み）';
COMMENT ON COLUMN financial_data.investment_return IS '期待投資利回り（年率%）';
COMMENT ON COLUMN financial_data.inflation_rate IS 'インフレ率（年率%）';
COMMENT ON COLUMN expense_items.category IS '支出カテゴリ（住居費、食費、交通費など）';
COMMENT ON COLUMN savings_items.type IS '貯蓄タイプ（deposit: 預金, investment: 投資, other: その他）';
COMMENT ON COLUMN retirement_data.current_age IS '現在の年齢';
COMMENT ON COLUMN retirement_data.retirement_age IS '退職予定年齢';
COMMENT ON COLUMN retirement_data.life_expectancy IS '平均寿命';
COMMENT ON COLUMN retirement_data.monthly_retirement_expenses IS '退職後の月間生活費';
COMMENT ON COLUMN retirement_data.pension_amount IS '月間年金受給予定額';
COMMENT ON COLUMN goals.type IS '目標タイプ（savings: 貯蓄, retirement: 退職, emergency: 緊急資金, custom: カスタム）';
COMMENT ON COLUMN goals.target_amount IS '目標金額';
COMMENT ON COLUMN goals.current_amount IS '現在の達成額';
COMMENT ON COLUMN goals.monthly_contribution IS '月間積立額';