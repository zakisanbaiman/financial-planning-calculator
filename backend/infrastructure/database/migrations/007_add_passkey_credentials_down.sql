-- 007_add_passkey_credentials_down.sql
-- パスキー（WebAuthn）認証機能のロールバック

DROP TRIGGER IF EXISTS update_webauthn_credentials_updated_at ON webauthn_credentials;
DROP INDEX IF EXISTS idx_webauthn_credentials_credential_id;
DROP INDEX IF EXISTS idx_webauthn_credentials_user_id;
DROP TABLE IF EXISTS webauthn_credentials;
