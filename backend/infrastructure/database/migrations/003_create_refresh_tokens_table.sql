-- 003_create_refresh_tokens_table.sql
-- リフレッシュトークンテーブルを作成してJWTトークンの自動更新を実装

-- リフレッシュトークンテーブル
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- インデックス: ユーザーIDでの検索を高速化
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- インデックス: 有効期限でのクリーンアップクエリを高速化
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- インデックス: トークンハッシュでの検索を高速化（検証時）
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);

-- 更新日時自動更新トリガー
CREATE TRIGGER update_refresh_tokens_updated_at
    BEFORE UPDATE ON refresh_tokens
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- コメント追加
COMMENT ON TABLE refresh_tokens IS 'リフレッシュトークン管理テーブル。JWTアクセストークンの自動更新に使用';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'トークンのSHA-256ハッシュ値。平文トークンは保存しない';
COMMENT ON COLUMN refresh_tokens.is_revoked IS 'トークンが失効されたかどうか。ログアウト時にtrueに設定';
COMMENT ON COLUMN refresh_tokens.last_used_at IS 'トークンが最後に使用された日時。不審なアクティビティの検出に使用';
