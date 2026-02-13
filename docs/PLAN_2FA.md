# 2段階認証（TOTP）実装計画

## 概要

現在の認証システム（ローカル認証 + GitHub OAuth）に、TOTP方式の2段階認証を追加する。Google Authenticator / Authy などの認証アプリに対応し、バックアップコードによるリカバリー機能も含める。

## 現状分析

### 既存の認証アーキテクチャ
- **ローカル認証**: Email/Password → bcryptハッシュ → JWT発行
- **OAuth認証**: GitHub OAuth → ユーザー作成/取得 → JWT発行
- **トークン管理**: アクセストークン（JWT）+ リフレッシュトークン（SHA256ハッシュ化してDB保存）

### 現在のUserエンティティ（2FAフィールドなし）
```go
type User struct {
    id, email, passwordHash, provider, providerUserID,
    name, avatarURL, emailVerified, emailVerifiedAt, createdAt, updatedAt
}