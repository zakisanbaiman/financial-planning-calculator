-- 006_add_two_factor_auth.sql
-- 2段階認証（2FA）機能の追加

ALTER TABLE users ADD COLUMN two_factor_enabled BOOLEAN DEFAULT false NOT NULL;
ALTER TABLE users ADD COLUMN two_factor_secret VARCHAR(255);
ALTER TABLE users ADD COLUMN two_factor_backup_codes TEXT[];

-- 2FA有効ユーザーを高速検索するためのインデックス
CREATE INDEX idx_users_two_factor_enabled ON users(two_factor_enabled);

-- コメント追加
COMMENT ON COLUMN users.two_factor_enabled IS '2段階認証が有効かどうか';
COMMENT ON COLUMN users.two_factor_secret IS 'TOTP用のシークレット（アプリケーション側で暗号化した値を保存）';
COMMENT ON COLUMN users.two_factor_backup_codes IS 'リカバリー用バックアップコード（ハッシュ化して保存）';
