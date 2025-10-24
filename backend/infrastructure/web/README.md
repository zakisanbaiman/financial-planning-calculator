# Web Infrastructure

このディレクトリには、Echo Webサーバーの設定とルーティングが含まれています。

## ファイル構成

- `middleware.go` - ミドルウェア設定（CORS、ログ、セキュリティ、レート制限など）
- `routes.go` - APIルーティング設定とハンドラー実装
- `routes_test.go` - ルーティングとハンドラーのテスト

## 実装されたミドルウェア

### セキュリティ
- **CORS**: フロントエンドからのアクセス許可
- **セキュリティヘッダー**: XSS保護、コンテンツタイプ検証など
- **レート制限**: API呼び出し頻度制限（100 req/sec）
- **リクエストサイズ制限**: 最大10MBまで

### パフォーマンス
- **Gzip圧縮**: レスポンスデータの圧縮
- **タイムアウト**: リクエストタイムアウト設定（30秒）

### 監視・デバッグ
- **ログ**: 詳細なリクエスト/レスポンスログ
- **リカバリー**: パニック時の自動復旧
- **リクエストID**: トレーサビリティ用のユニークID生成

### エラーハンドリング
- **統一エラーレスポンス**: 一貫したエラー形式
- **ログ出力**: エラーの詳細ログ
- **リクエストID付与**: エラー追跡用

## APIエンドポイント

### 基本エンドポイント
- `GET /health` - ヘルスチェック
- `GET /api/` - API情報
- `GET /swagger/*` - Swagger UI

### 財務データ管理
- `POST /api/financial-data` - 財務データ作成
- `GET /api/financial-data` - 財務データ取得
- `PUT /api/financial-data/:id` - 財務データ更新
- `DELETE /api/financial-data/:id` - 財務データ削除

### 計算機能
- `POST /api/calculations/asset-projection` - 資産推移計算
- `POST /api/calculations/retirement` - 老後資金計算
- `POST /api/calculations/emergency-fund` - 緊急資金計算

### 目標管理
- `POST /api/goals` - 目標作成
- `GET /api/goals` - 目標一覧取得
- `PUT /api/goals/:id` - 目標更新
- `DELETE /api/goals/:id` - 目標削除

### レポート生成
- `GET /api/reports/pdf` - PDFレポート生成

## 設定

環境変数による設定が可能です：

```bash
# サーバー設定
PORT=8080
DEBUG=false

# CORS設定
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001

# レート制限
RATE_LIMIT_RPS=100

# リクエスト設定
REQUEST_TIMEOUT=30s
MAX_REQUEST_SIZE=10M

# 圧縮設定
ENABLE_GZIP=true
GZIP_LEVEL=5

# セキュリティ設定
ENABLE_SECURE_HEADERS=true
```

## テスト実行

```bash
go test ./infrastructure/web/... -v
```

## 使用方法

```go
import (
    "github.com/financial-planning-calculator/backend/config"
    "github.com/financial-planning-calculator/backend/infrastructure/web"
    "github.com/labstack/echo/v4"
)

func main() {
    cfg := config.LoadServerConfig()
    e := echo.New()
    
    e.HTTPErrorHandler = web.CustomHTTPErrorHandler
    web.SetupMiddleware(e, cfg)
    web.SetupRoutes(e)
    
    e.Logger.Fatal(e.Start(":" + cfg.Port))
}
```

## 次のステップ

現在のハンドラーはプレースホルダー実装です。次のタスクで以下を実装予定：

1. コントローラー層の実装
2. バリデーション機能
3. 実際のビジネスロジック統合
4. エラーハンドリングの詳細化