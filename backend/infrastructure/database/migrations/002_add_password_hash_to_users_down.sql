-- 002_add_password_hash_to_users_down.sql
-- マイグレーションのロールバック

-- password_hashカラムを削除
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;

-- トリガーとファンクションの削除（他のテーブルでも使用している可能性があるため、注意が必要）
-- DROP TRIGGER IF EXISTS update_users_updated_at ON users;
-- DROP FUNCTION IF EXISTS update_updated_at_column();
