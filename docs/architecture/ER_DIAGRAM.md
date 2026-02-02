# ER図 (Entity-Relationship Diagram)

このドキュメントは、財務計画計算機アプリケーションのデータベーススキーマを可視化したER図です。

## データベース構造

```mermaid
erDiagram
    users ||--o{ financial_data : "has"
    users ||--o{ retirement_data : "has"
    users ||--o{ goals : "has"
    users ||--o{ refresh_tokens : "has"
    users ||--o{ webauthn_credentials : "has"
    financial_data ||--|{ expense_items : "contains"
    financial_data ||--|{ savings_items : "contains"

    users {
        varchar id PK "ユーザーID"
        varchar email UK "メールアドレス"
        varchar password_hash "パスワードハッシュ（nullable）"
        varchar provider "認証プロバイダー（local, github等）"
        varchar provider_user_id "OAuth プロバイダーのユーザーID"
        varchar name "表示名"
        text avatar_url "アバターURL"
        boolean two_factor_enabled "2FA有効フラグ"
        varchar two_factor_secret "TOTPシークレット"
        text_array two_factor_backup_codes "バックアップコード"
        boolean email_verified "メール認証済みフラグ"
        varchar email_verification_token "メール認証トークン"
        timestamp email_verification_sent_at "認証メール送信日時"
        timestamp created_at "作成日時"
        timestamp updated_at "更新日時"
    }

    financial_data {
        uuid id PK "財務データID"
        varchar user_id FK "ユーザーID"
        decimal monthly_income "月収"
        decimal investment_return "期待投資利回り（%）"
        decimal inflation_rate "インフレ率（%）"
        timestamp created_at "作成日時"
        timestamp updated_at "更新日時"
    }

    expense_items {
        uuid id PK "支出項目ID"
        uuid financial_data_id FK "財務データID"
        varchar category "支出カテゴリ"
        decimal amount "支出金額"
        text description "説明"
        timestamp created_at "作成日時"
        timestamp updated_at "更新日時"
    }

    savings_items {
        uuid id PK "貯蓄項目ID"
        uuid financial_data_id FK "財務データID"
        varchar type "貯蓄タイプ（deposit/investment/other）"
        decimal amount "金額"
        text description "説明"
        timestamp created_at "作成日時"
        timestamp updated_at "更新日時"
    }

    retirement_data {
        uuid id PK "退職データID"
        varchar user_id FK "ユーザーID"
        integer current_age "現在の年齢"
        integer retirement_age "退職予定年齢"
        integer life_expectancy "平均寿命"
        decimal monthly_retirement_expenses "退職後の月間生活費"
        decimal pension_amount "月間年金受給予定額"
        timestamp created_at "作成日時"
        timestamp updated_at "更新日時"
    }

    goals {
        uuid id PK "目標ID"
        varchar user_id FK "ユーザーID"
        varchar type "目標タイプ（savings/retirement/emergency/custom）"
        varchar title "タイトル"
        decimal target_amount "目標金額"
        date target_date "目標日"
        decimal current_amount "現在の達成額"
        decimal monthly_contribution "月間積立額"
        boolean is_active "アクティブフラグ"
        timestamp created_at "作成日時"
        timestamp updated_at "更新日時"
    }

    refresh_tokens {
        uuid id PK "リフレッシュトークンID"
        varchar user_id FK "ユーザーID"
        varchar token UK "トークン値"
        timestamp expires_at "有効期限"
        timestamp created_at "作成日時"
    }

    webauthn_credentials {
        varchar id PK "クレデンシャルID"
        varchar user_id FK "ユーザーID"
        bytea credential_id UK "WebAuthn credential ID"
        bytea public_key "公開鍵"
        varchar attestation_type "認証タイプ"
        bytea aaguid "Authenticator AAGUID"
        bigint sign_count "署名カウンター"
        boolean clone_warning "クローン警告フラグ"
        text_array transports "対応トランスポート"
        varchar name "クレデンシャル名"
        timestamp created_at "作成日時"
        timestamp updated_at "更新日時"
        timestamp last_used_at "最終使用日時"
    }
```

## テーブル説明

### users（ユーザー）
ユーザー情報を管理するマスターテーブル。認証情報（ローカル認証、OAuth、2FA、パスキー）を含みます。

### financial_data（財務データ）
ユーザーの基本的な財務情報（月収、期待投資利回り、インフレ率）を保存します。各ユーザーに対して1レコードのみ（UNIQUE制約）。

### expense_items（支出項目）
月間支出の詳細を保存します。カテゴリ別に複数の支出項目を管理可能。

### savings_items（貯蓄項目）
現在の貯蓄・投資の詳細を保存します。預金、投資、その他の3タイプに分類。

### retirement_data（退職・年金情報）
老後計画に関する情報を保存します。各ユーザーに対して1レコードのみ（UNIQUE制約）。

### goals（目標）
ユーザーの財務目標を管理します。貯蓄、退職、緊急資金、カスタムの4タイプをサポート。

### refresh_tokens（リフレッシュトークン）
JWT認証のリフレッシュトークンを管理します。セキュアなトークン更新機構を実現。

### webauthn_credentials（WebAuthn認証情報）
パスキー（生体認証等）による認証情報を管理します。複数のデバイスに対応。

## インデックス

主要なインデックスは以下の通り：

- `users`: email, provider+provider_user_id, two_factor_enabled
- `financial_data`: user_id
- `expense_items`: financial_data_id, category
- `savings_items`: financial_data_id, type
- `retirement_data`: user_id
- `goals`: user_id, type, is_active, target_date, user_id+is_active（部分インデックス）
- `refresh_tokens`: user_id, token
- `webauthn_credentials`: user_id, credential_id

## 制約

### CHECK制約
- 金額フィールド: 負の値を許可しない（一部は正の値のみ）
- 年齢フィールド: 0-120の範囲
- パーセンテージ: 適切な範囲（投資利回り: 0-100%, インフレ率: 0-50%）
- 退職年齢 > 現在の年齢
- 平均寿命 > 退職年齢
- 目標日 > 現在日

### 外部キー制約
すべての外部キーは`ON DELETE CASCADE`を設定し、親レコード削除時に関連レコードも自動削除されます。

## トリガー

すべてのテーブルに`updated_at`列の自動更新トリガーが設定されています。レコード更新時に自動的に現在日時が設定されます。
