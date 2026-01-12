# Plan: ベーシック認証の実装

全エンドポイントに対してベーシック認証を追加し、環境変数で認証情報を管理します。ヘルスチェックとSwagger UIは認証から除外します。

## ベーシック認証とは？

**HTTP Basic Authentication**は、最もシンプルな認証方式です：
- ユーザー名とパスワードで保護
- ブラウザが標準で対応（ダイアログが自動表示）
- `Authorization: Basic {base64(username:password)}` ヘッダーで送信

### 動作フロー

```
クライアント                    サーバー
    |                              |
    |---(1) GET /api/users-------->|
    |                              |
    |<--(2) 401 Unauthorized-------|
    |     WWW-Authenticate: Basic  |
    |                              |
[ログインダイアログ表示]          |
    |                              |
    |---(3) Authorization付き----->|
    |     Basic YWRtaW46cGFzcw==   |
    |                              |
    |                        [認証情報検証]
    |                              |
    |<--(4) 200 OK + データ--------|
```

### なぜこのプロジェクトに必要？

現在、APIは**誰でもアクセス可能**です。ベーシック認証を追加することで：
- ✅ 不正アクセスを防止
- ✅ デモ/ステージング環境を保護
- ✅ 本番環境での基本的なセキュリティ確保

## Steps

### Step 1: 設定を読み込む準備（Config層）
**ファイル**: `backend/config/server.go`

環境変数から認証情報を読み込む設定を追加します。

**何をするか？**
- 認証のON/OFF (`ENABLE_BASIC_AUTH`)
- ユーザー名 (`BASIC_AUTH_USERNAME`)
- パスワード (`BASIC_AUTH_PASSWORD`)

を環境変数から読み込めるようにします。

### Step 2: 認証機能を追加（Middleware層）
**ファイル**: `backend/infrastructure/web/middleware.go`

すべてのAPIリクエストに認証チェックを追加します。

**何をするか？**
- Echoの`BasicAuth`ミドルウェアを使用
- リクエストの`Authorization`ヘッダーを検証
- 正しければ通過、間違っていれば401エラー

**除外するエンドポイント**:
- `/health` - Render.comのヘルスチェック用
- `/swagger/*` - 開発時のAPI仕様確認用

### Step 3: セキュリティ警告を追加（Main層）
**ファイル**: `backend/main.go`

デフォルトパスワードを使っている場合に警告を表示します。

**何をするか？**
- 起動時にパスワードが`change-me`のままなら警告
- 本番環境での設定ミスを防ぐ

### Step 4: 環境変数を設定（Configuration Files）
**ファイル**: `backend/.env.local`, `docker-compose.yml`, `render.yaml`

各環境で認証情報を設定できるようにします。

**何をするか？**
- 開発環境: `.env.local`に追加（デフォルトOFF）
- Docker環境: `docker-compose.yml`に追加
- 本番環境: `render.yaml`に追加（機密情報は手動設定）

## Further Considerations

1. フロントエンドのAPIクライアントにBasic認証ヘッダーを追加する必要がありますか？（SPAからの認証方法）

2. E2Eテストファイルに認証ヘッダーを追加しますか？

3. 開発環境ではデフォルトで無効（`ENABLE_BASIC_AUTH=false`）、本番環境のみ有効にしますか？

## Implementation Details

### 1. ServerConfig の拡張

```go
type ServerConfig struct {
    // 既存のフィールド...
    
    // Basic認証設定
    EnableBasicAuth   bool
    BasicAuthUsername string
    BasicAuthPassword string
}

func LoadServerConfig() *ServerConfig {
    config := &ServerConfig{
        // 既存の設定...
        
        // Basic認証設定
        EnableBasicAuth:    getEnvBool("ENABLE_BASIC_AUTH", false),
        BasicAuthUsername:  getEnv("BASIC_AUTH_USERNAME", "admin"),
        BasicAuthPassword:  getEnv("BASIC_AUTH_PASSWORD", "change-me"),
    }
    return config
}
```

### 2. BasicAuth ミドルウェアの追加

```go
// Basic認証 - 全エンドポイントを保護（オプション）
if cfg.EnableBasicAuth {
    e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
        Skipper: func(c echo.Context) bool {
            // ヘルスチェックとSwaggerは除外
            return c.Path() == "/health" || strings.HasPrefix(c.Path(), "/swagger")
        },
        Validator: func(username, password string, c echo.Context) (bool, error) {
            if username == cfg.BasicAuthUsername && password == cfg.BasicAuthPassword {
                return true, nil
            }
            return false, nil
        },
        Realm: "Financial Planning Calculator API",
    }))
}
```

### 3. セキュリティ警告の追加

```go
// Basic認証の警告
if serverCfg.EnableBasicAuth && serverCfg.BasicAuthPassword == "change-me" {
    warnings = append(warnings, "⚠️  BASIC_AUTH_PASSWORD is set to default value")
}
```

### 4. 環境変数の追加

**backend/.env.local**:
```env
# Basic Authentication
ENABLE_BASIC_AUTH=false
BASIC_AUTH_USERNAME=admin
BASIC_AUTH_PASSWORD=change-me
```

**docker-compose.yml**:
```yaml
backend:
  environment:
    ENABLE_BASIC_AUTH: "true"
    BASIC_AUTH_USERNAME: admin
    BASIC_AUTH_PASSWORD: secure_password_here
```

**render.yaml**:
```yaml
services:
  - type: web
    envVars:
      - key: ENABLE_BASIC_AUTH
        value: true
      - key: BASIC_AUTH_USERNAME
        sync: false
      - key: BASIC_AUTH_PASSWORD
        sync: false
```

## Security Considerations

- 環境変数でのみ管理し、コードにハードコードしない
- デフォルト値使用時に警告を表示
- HTTPS経由でのみ使用（Render.comは自動でHTTPS対応）
- 本番環境では強力なパスワードを使用
- 定期的なパスワードローテーションを推奨
