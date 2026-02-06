# システム構成図 (System Architecture Diagram)

このドキュメントは、財務計画計算機アプリケーションのシステム全体の構成とネットワーク構造を可視化した図です。

## システム全体構成

```mermaid
graph TB
    subgraph "クライアント"
        Browser["Webブラウザ"]
        PasskeyDevice["パスキーデバイス<br/>（生体認証等）"]
    end

    subgraph "フロントエンド層"
        NextJS["Next.js 14<br/>TypeScript<br/>Port: 3000"]
        NextJSFeatures["- App Router<br/>- React Hook Form<br/>- Chart.js<br/>- Tailwind CSS"]
    end

    subgraph "バックエンド層"
        API["Go API Server<br/>Echo Framework<br/>Port: 8080"]
        APIFeatures["- RESTful API<br/>- JWT認証<br/>- OpenAPI/Swagger<br/>- pprof（開発環境）"]
    end

    subgraph "データベース層"
        PostgreSQL["PostgreSQL 15<br/>Port: 5432"]
        DBFeatures["- UUIDエクステンション<br/>- トリガー<br/>- インデックス最適化"]
    end

    subgraph "外部サービス"
        GitHub["GitHub OAuth<br/>認証プロバイダー"]
        WebAuthn["WebAuthn/FIDO2<br/>パスキー認証"]
    end

    subgraph "開発ツール"
        Swagger["Swagger UI<br/>localhost:8080/swagger"]
        pprof["pprof<br/>localhost:6060/debug/pprof"]
    end

    Browser --> NextJS
    PasskeyDevice -.-> Browser
    NextJS --> API
    API --> PostgreSQL
    API --> GitHub
    API --> WebAuthn
    Browser --> Swagger
    Browser --> pprof

    style NextJS fill:#61dafb
    style API fill:#00add8
    style PostgreSQL fill:#336791
    style GitHub fill:#181717
    style WebAuthn fill:#3c790a
```

## Docker環境構成

```mermaid
graph TB
    subgraph "Docker Network: financial_planning_network"
        subgraph "frontend Container"
            FrontendApp["Next.js App<br/>Port: 3000"]
            FrontendVolume["Volume: ./frontend:/app"]
        end

        subgraph "backend Container"
            BackendApp["Go API<br/>Port: 8080, 6060"]
            Air["Air<br/>ホットリロード"]
            BackendVolume["Volume: ./backend:/app"]
        end

        subgraph "postgres Container"
            DB["PostgreSQL 15<br/>Port: 5432"]
            DBVolume["Volume: postgres_data"]
            InitScripts["Init Scripts<br/>./backend/infrastructure/database/init"]
        end

        subgraph "db-tools Container"
            DBTools["Migration & Seed Tools"]
            DBToolsVolume["Volume: ./backend:/app"]
        end
    end

    FrontendApp --> BackendApp
    BackendApp --> DB
    DBTools --> DB

    Host["ホストマシン<br/>localhost"]
    Host -- "3000:3000" --> FrontendApp
    Host -- "8080:8080" --> BackendApp
    Host -- "6060:6060" --> BackendApp
    Host -- "5432:5432" --> DB

    style FrontendApp fill:#61dafb
    style BackendApp fill:#00add8
    style DB fill:#336791
```

## デプロイメント構成（本番環境）

```mermaid
graph TB
    subgraph "Render.com"
        subgraph "Web Service"
            Frontend["Next.js<br/>Static/SSR"]
        end

        subgraph "Backend Service"
            Backend["Go API Server"]
        end

        subgraph "Database Service"
            RenderDB["PostgreSQL<br/>Managed Database"]
        end
    end

    Internet["インターネット"]
    
    Internet --> Frontend
    Frontend --> Backend
    Backend --> RenderDB

    subgraph "GitHub"
        Repository["GitHubリポジトリ"]
        Actions["GitHub Actions<br/>CI/CD"]
    end

    Repository --> Actions
    Actions --> Frontend
    Actions --> Backend

    style Frontend fill:#61dafb
    style Backend fill:#00add8
    style RenderDB fill:#336791
```

## ネットワークフロー

### ユーザー登録・ログインフロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Browser as ブラウザ
    participant Frontend as Next.js
    participant Backend as Go API
    participant DB as PostgreSQL
    participant OAuth as GitHub OAuth

    Note over User,OAuth: ローカル認証フロー
    User->>Browser: 登録フォーム入力
    Browser->>Frontend: POST /register
    Frontend->>Backend: POST /api/auth/register
    Backend->>DB: INSERT users
    DB-->>Backend: OK
    Backend-->>Frontend: {user, token}
    Frontend-->>Browser: ログイン成功
    Browser-->>User: ダッシュボード表示

    Note over User,OAuth: OAuth認証フロー
    User->>Browser: "GitHubでログイン"
    Browser->>Frontend: Click
    Frontend->>Backend: GET /api/auth/github
    Backend->>OAuth: 認証リクエスト
    OAuth-->>Browser: リダイレクト
    Browser->>OAuth: ユーザー認証
    OAuth-->>Backend: Callback + Code
    Backend->>OAuth: トークン交換
    OAuth-->>Backend: Access Token + User Info
    Backend->>DB: UPSERT users
    DB-->>Backend: OK
    Backend-->>Frontend: リダイレクト + Token
    Frontend-->>Browser: ログイン成功
    Browser-->>User: ダッシュボード表示

    Note over User,OAuth: パスキー登録フロー
    User->>Browser: パスキー登録開始
    Browser->>Frontend: POST /passkey/register/begin
    Frontend->>Backend: POST /api/webauthn/register/begin
    Backend->>Backend: Generate Challenge
    Backend-->>Frontend: {challenge, options}
    Frontend-->>Browser: navigator.credentials.create()
    Browser-->>User: 生体認証プロンプト
    User->>Browser: 認証実行
    Browser-->>Frontend: Credential
    Frontend->>Backend: POST /api/webauthn/register/finish
    Backend->>DB: INSERT webauthn_credentials
    DB-->>Backend: OK
    Backend-->>Frontend: Success
    Frontend-->>Browser: 登録完了
    Browser-->>User: 完了メッセージ
```

### 財務計画計算フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as Next.js
    participant Backend as Go API
    participant DB as PostgreSQL
    participant Calculator as 計算サービス

    User->>Frontend: 財務データ入力
    Frontend->>Backend: POST /api/financial-data
    Backend->>DB: INSERT/UPDATE financial_data
    Backend->>DB: INSERT expense_items
    Backend->>DB: INSERT savings_items
    DB-->>Backend: OK

    User->>Frontend: 将来予測リクエスト
    Frontend->>Backend: GET /api/calculations/projection?years=30
    Backend->>DB: SELECT financial_data, goals, retirement_data
    DB-->>Backend: Data
    Backend->>Calculator: GenerateProjection()
    Calculator->>Calculator: 複利計算
    Calculator->>Calculator: 資産推移予測
    Calculator->>Calculator: 目標進捗評価
    Calculator-->>Backend: PlanProjection
    Backend-->>Frontend: JSON Response
    Frontend->>Frontend: Chart.js描画
    Frontend-->>User: グラフ表示
```

### レポート生成フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Frontend as Next.js
    participant Backend as Go API
    participant DB as PostgreSQL
    participant PDF as PDFジェネレーター

    User->>Frontend: PDFレポート生成リクエスト
    Frontend->>Backend: POST /api/reports/generate
    Backend->>DB: SELECT 全財務データ
    DB-->>Backend: Data
    Backend->>PDF: GeneratePDF()
    PDF->>PDF: HTML生成
    PDF->>PDF: PDFレンダリング
    PDF-->>Backend: PDF Bytes
    Backend-->>Frontend: application/pdf
    Frontend-->>User: PDFダウンロード
```

## セキュリティ構成

```mermaid
graph TB
    subgraph "セキュリティ層"
        HTTPS["HTTPS/TLS<br/>通信暗号化"]
        JWT["JWT認証<br/>トークンベース"]
        CORS["CORS設定<br/>オリジン制御"]
        RateLimit["レート制限<br/>DDoS対策"]
        InputValidation["入力検証<br/>Zod + Echo Validator"]
    end

    subgraph "認証・認可"
        LocalAuth["ローカル認証<br/>bcrypt"]
        OAuthAuth["OAuth 2.0<br/>GitHub"]
        TwoFactorAuth["2FA<br/>TOTP"]
        PasskeyAuth["パスキー<br/>WebAuthn/FIDO2"]
    end

    subgraph "データ保護"
        PasswordHash["パスワードハッシュ化<br/>bcrypt"]
        TokenEncryption["トークン暗号化"]
        DBEncryption["DB通信暗号化<br/>SSL/TLS"]
        SecretManagement["シークレット管理<br/>環境変数"]
    end

    HTTPS --> JWT
    JWT --> CORS
    CORS --> RateLimit
    RateLimit --> InputValidation

    LocalAuth --> PasswordHash
    OAuthAuth --> TokenEncryption
    TwoFactorAuth --> TokenEncryption
    PasskeyAuth --> DBEncryption

    style HTTPS fill:#3c790a
    style JWT fill:#3c790a
    style LocalAuth fill:#3c790a
    style PasswordHash fill:#3c790a
```

## 技術スタック詳細

### フロントエンド
| 技術 | バージョン | 用途 |
|------|-----------|------|
| Next.js | 14 | Reactフレームワーク、App Router |
| TypeScript | 5.x | 型安全な開発 |
| Tailwind CSS | 3.x | ユーティリティファーストCSS |
| Chart.js | 4.x | データ可視化 |
| React Hook Form | 7.x | フォーム管理 |
| Zod | 3.x | スキーマ検証 |
| Axios | 1.x | HTTPクライアント |

### バックエンド
| 技術 | バージョン | 用途 |
|------|-----------|------|
| Go | 1.24.0 | プログラミング言語 |
| Echo | 4.x | Webフレームワーク |
| GORM | 1.x | ORM（検討中） |
| go-webauthn | 0.x | WebAuthn実装 |
| jwt-go | 5.x | JWT認証 |
| bcrypt | - | パスワードハッシュ化 |
| golang-migrate | 4.x | DBマイグレーション |
| swaggo | 1.x | OpenAPI/Swagger生成 |
| pprof | - | パフォーマンスプロファイリング |

### データベース
| 技術 | バージョン | 用途 |
|------|-----------|------|
| PostgreSQL | 15 | リレーショナルデータベース |
| uuid-ossp | - | UUID生成エクステンション |

### 開発ツール
| 技術 | バージョン | 用途 |
|------|-----------|------|
| Docker | - | コンテナ化 |
| Docker Compose | - | マルチコンテナ管理 |
| Air | - | Goホットリロード |
| golangci-lint | 1.64+ | Goリンター |
| ESLint | - | JavaScriptリンター |
| Prettier | - | コードフォーマッター |

## ポート一覧

| サービス | ポート | 説明 |
|---------|-------|------|
| フロントエンド | 3000 | Next.js開発サーバー |
| バックエンドAPI | 8080 | Go Echo サーバー |
| pprof | 6060 | パフォーマンスプロファイリング（開発環境のみ） |
| PostgreSQL | 5432 | データベース |
| Swagger UI | 8080/swagger | API仕様書UI |

## 環境変数

### バックエンド
```bash
# データベース
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=financial_planning
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# OAuth
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_CALLBACK_URL=http://localhost:8080/api/auth/github/callback

# アプリケーション
GIN_MODE=debug  # 注: 現在は使用されていない（レガシー設定）
ENABLE_PPROF=true
PPROF_PORT=6060
```

### フロントエンド
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
WATCHPACK_POLLING=true
```

## ヘルスチェック

### バックエンド
- **エンドポイント**: `GET /health`
- **間隔**: 30秒
- **タイムアウト**: 10秒
- **再試行**: 3回

### PostgreSQL
- **コマンド**: `pg_isready -U postgres -d financial_planning`
- **間隔**: 10秒
- **タイムアウト**: 5秒
- **再試行**: 5回

## スケーリング戦略

### 水平スケーリング
- フロントエンド: Render.comの自動スケーリング
- バックエンド: 複数インスタンスへの負荷分散
- データベース: リードレプリカの追加（将来）

### 垂直スケーリング
- より大きなコンテナサイズへの移行
- データベースリソースの増強

## 監視とロギング

### ログ
- アプリケーションログ: stdout/stderr
- アクセスログ: Echo middleware
- エラーログ: 構造化ログ（JSON）

### メトリクス（将来実装予定）
- pprof: パフォーマンスメトリクス
- Prometheus: メトリクス収集
- Grafana: 可視化

## バックアップ戦略

### データベース
- 自動バックアップ: 日次
- ポイントインタイムリカバリ: 7日間
- スナップショット: 手動実行可能

### アプリケーション
- GitHubリポジトリ: ソースコードのバックアップ
- 環境変数: セキュアな管理（Render.com）
