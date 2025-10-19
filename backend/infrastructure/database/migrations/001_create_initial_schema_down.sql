-- 001_create_initial_schema_down.sql
-- 初期スキーマのロールバック

-- トリガーを削除
DROP TRIGGER IF EXISTS update_goals_updated_at ON goals;
DROP TRIGGER IF EXISTS update_retirement_data_updated_at ON retirement_data;
DROP TRIGGER IF EXISTS update_savings_items_updated_at ON savings_items;
DROP TRIGGER IF EXISTS update_expense_items_updated_at ON expense_items;
DROP TRIGGER IF EXISTS update_financial_data_updated_at ON financial_data;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- トリガー関数を削除
DROP FUNCTION IF EXISTS update_updated_at_column();

-- テーブルを削除（外部キー制約の順序を考慮）
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS retirement_data;
DROP TABLE IF EXISTS savings_items;
DROP TABLE IF EXISTS expense_items;
DROP TABLE IF EXISTS financial_data;
DROP TABLE IF EXISTS users;

-- エクステンションを削除（他で使用されていない場合のみ）
-- DROP EXTENSION IF EXISTS "uuid-ossp";