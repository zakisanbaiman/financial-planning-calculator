-- 002_add_password_hash_to_users.sql
-- ユーザーテーブルにパスワードハッシュカラムを追加してJWT認証を実装可能にする

-- password_hashカラムを追加
-- 既存のユーザーが存在する可能性を考慮してNULL許可で追加
ALTER TABLE users ADD COLUMN password_hash VARCHAR(255);

-- 既存データがない場合、またはすべてのユーザーがパスワードを設定した後に
-- NOT NULL制約を追加する場合は、以下のコメントを外して実行
-- ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;

-- インデックス: メールアドレスでの検索を高速化（ログイン時）
-- UNIQUE制約が既に存在するため、自動的にインデックスが作成されているが明示的に記載
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 更新日時自動更新トリガーの作成（まだ存在しない場合）
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- usersテーブルに更新日時トリガーを設定（まだ存在しない場合）
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
