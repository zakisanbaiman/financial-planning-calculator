-- 007_add_passkey_credentials.sql
-- パスキー（WebAuthn）認証機能の追加

-- WebAuthn認証情報テーブル
CREATE TABLE webauthn_credentials (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_id BYTEA NOT NULL UNIQUE,
    public_key BYTEA NOT NULL,
    attestation_type VARCHAR(50) NOT NULL,
    aaguid BYTEA NOT NULL,
    sign_count BIGINT NOT NULL DEFAULT 0 CHECK (sign_count >= 0),
    clone_warning BOOLEAN NOT NULL DEFAULT false,
    transports TEXT[],
    name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE
);

-- インデックス作成
CREATE INDEX idx_webauthn_credentials_user_id ON webauthn_credentials(user_id);
CREATE INDEX idx_webauthn_credentials_credential_id ON webauthn_credentials(credential_id);

-- 更新日時自動更新のトリガーを設定
CREATE TRIGGER update_webauthn_credentials_updated_at BEFORE UPDATE ON webauthn_credentials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント追加
COMMENT ON TABLE webauthn_credentials IS 'WebAuthn（パスキー）認証情報テーブル';
COMMENT ON COLUMN webauthn_credentials.id IS 'クレデンシャルの一意識別子';
COMMENT ON COLUMN webauthn_credentials.user_id IS 'ユーザーID';
COMMENT ON COLUMN webauthn_credentials.credential_id IS 'WebAuthn credential ID（バイナリ）';
COMMENT ON COLUMN webauthn_credentials.public_key IS '公開鍵（バイナリ）';
COMMENT ON COLUMN webauthn_credentials.attestation_type IS '認証タイプ（none, indirect, direct）';
COMMENT ON COLUMN webauthn_credentials.aaguid IS 'Authenticator AAGUID';
COMMENT ON COLUMN webauthn_credentials.sign_count IS '署名カウンター（クローン検出用）';
COMMENT ON COLUMN webauthn_credentials.clone_warning IS 'クローン警告フラグ';
COMMENT ON COLUMN webauthn_credentials.transports IS '対応トランスポート（usb, nfc, ble, internal）';
COMMENT ON COLUMN webauthn_credentials.name IS 'クレデンシャルの名前（ユーザーが設定）';
COMMENT ON COLUMN webauthn_credentials.last_used_at IS '最終使用日時';
