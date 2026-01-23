# Scripts Directory

このディレクトリには、プロジェクトで使用する各種スクリプトが含まれています。

## Render.com MCP統合スクリプト

### セットアップ

```bash
cd scripts
npm install
```

### スクリプト一覧

#### `mcp-server-render.js`

Render.com用のMCPサーバー。AIアシスタント（Claude、Copilot、Cursor）がRender.comのデプロイ情報にアクセスするためのサーバーです。

**使用方法**:
```bash
# 環境変数を設定
export RENDER_API_KEY="rnd_xxxxxxxxxxxxx"
export RENDER_OWNER_ID="dsy-xxxxxxxxxxxxx"  # オプション

# MCPサーバーを起動（通常はAIツールが自動で起動）
node mcp-server-render.js
```

**提供機能**:
- `list_services`: サービス一覧の取得
- `get_service_status`: サービス状態の確認
- `list_recent_deploys`: デプロイ履歴の取得
- `get_deployment_logs`: ログの取得
- `detect_errors`: エラーの自動検出

#### `check-render-deployments.js`

Render.comのデプロイメント状態をチェックし、エラーを検出するスクリプト。

**使用方法**:
```bash
# 環境変数を設定
export RENDER_API_KEY="rnd_xxxxxxxxxxxxx"

# デプロイメントをチェック
node check-render-deployments.js
```

**出力例**:
```
🔍 Checking Render.com deployments...

📦 Found 3 services

📋 Service: financial-planning-backend
   Type: web
   Status: available
   ✅ Deployment is live and healthy

...

✅ All deployments are healthy
```

**GitHub Actionsでの使用**:

このスクリプトは`.github/workflows/monitor-render-deployments.yml`から自動的に実行されます。

#### `summarize_failure.py`

テスト失敗の要約を生成するPythonスクリプト。

**使用方法**:
```bash
# 依存関係のインストール
pip install -r requirements.txt

# GitHub Actionsで自動実行される
```

#### `test-integration.sh`

統合テストを実行するシェルスクリプト。

**使用方法**:
```bash
./test-integration.sh
```

#### `run-ci-local.sh`

CI環境をローカルで再現するシェルスクリプト。

**使用方法**:
```bash
./run-ci-local.sh
```

## 環境変数

### Render.com関連

- `RENDER_API_KEY`: Render.com APIキー（必須）
- `RENDER_OWNER_ID`: Owner ID（オプション、特定組織のサービスのみ表示）

### 取得方法

1. [Render.com Dashboard](https://dashboard.render.com)にログイン
2. Account Settings → API Keys
3. Create API Keyで新しいキーを生成

## トラブルシューティング

### `Cannot find module '@modelcontextprotocol/sdk'`

```bash
cd scripts
npm install
```

### `RENDER_API_KEY environment variable is required`

```bash
export RENDER_API_KEY="your-api-key-here"
```

### `Permission denied`

```bash
chmod +x check-render-deployments.js
chmod +x mcp-server-render.js
chmod +x test-integration.sh
chmod +x run-ci-local.sh
```

## 関連ドキュメント

- [MCP セットアップガイド](../docs/MCP_SETUP.md)
- [MCP クイックリファレンス](../docs/MCP_QUICK_REFERENCE.md)
- [使用例](../docs/MCP_USAGE_EXAMPLES.md)
- [AI統合ガイド](../docs/AI_COPILOT_INTEGRATION.md)

## サポート

問題がある場合は、GitHubのIssueを作成してください。
