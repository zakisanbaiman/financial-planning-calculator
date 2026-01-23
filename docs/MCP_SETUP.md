# Render.com MCP セットアップガイド

このガイドでは、Render.comのデプロイメントをMCP (Model Context Protocol)経由で監視し、AI assistantが自動的にエラーを検知・修正できるようにする設定方法を説明します。

## 概要

MCP (Model Context Protocol)は、AIアシスタント（Claude Code、GitHub Copilot、Cursor等）が外部システムと連携するための標準プロトコルです。このプロジェクトでは、Render.comのデプロイメント監視用MCPサーバーを提供しています。

### 提供機能

1. **デプロイメント監視**: サービスの状態とデプロイ履歴の確認
2. **エラー検知**: デプロイログから自動的にエラーパターンを検出
3. **ログ取得**: デプロイメントログの詳細な取得
4. **自動通知**: GitHub PRへのエラー通知とIssue作成

## 前提条件

- Node.js 18.0.0以上
- Render.comアカウント
- Render.com APIキー
- Claude Code、GitHub Copilot、またはCursor（いずれか）

## 1. Render.com APIキーの取得

1. [Render.com Dashboard](https://dashboard.render.com)にログイン
2. 右上のアカウントメニューから「Account Settings」を選択
3. 左メニューから「API Keys」を選択
4. 「Create API Key」をクリック
5. キーの名前を入力（例: "MCP Integration"）
6. 生成されたAPIキーをコピー（**一度しか表示されないので注意！**）

### Owner IDの取得（オプション）

Owner IDを使用すると、特定の組織のサービスのみをフィルタリングできます。

1. Render.com Dashboardで任意のサービスを開く
2. URLから Owner ID を確認: `https://dashboard.render.com/d/{OWNER_ID}/...`

## 2. 環境変数の設定

### ローカル開発用

`.env`ファイルまたはシェルの設定ファイル（`.zshrc`, `.bashrc`等）に追加:

```bash
export RENDER_API_KEY="rnd_xxxxxxxxxxxxxxxxxxxxx"
export RENDER_OWNER_ID="dsy-xxxxxxxxxxxxxxxxxxxxx"  # オプション
```

### GitHub Actions用

GitHub リポジトリの Settings → Secrets and variables → Actions で以下を設定:

- `RENDER_API_KEY`: Render.com APIキー
- `RENDER_OWNER_ID`: Owner ID（オプション）

## 3. 依存関係のインストール

```bash
cd scripts
npm install
```

これにより、MCPサーバーに必要な`@modelcontextprotocol/sdk`がインストールされます。

## 4. MCPサーバーの設定

### Claude Desktop用

`~/Library/Application Support/Claude/claude_desktop_config.json`（macOS）または
`%APPDATA%\Claude\claude_desktop_config.json`（Windows）を編集:

```json
{
  "mcpServers": {
    "render": {
      "command": "node",
      "args": [
        "/absolute/path/to/financial-planning-calculator/scripts/mcp-server-render.js"
      ],
      "env": {
        "RENDER_API_KEY": "rnd_xxxxxxxxxxxxxxxxxxxxx",
        "RENDER_OWNER_ID": "dsy-xxxxxxxxxxxxxxxxxxxxx"
      }
    }
  }
}
```

**重要**: `args`のパスは絶対パスで指定してください。

### Cursor用

Cursor Settings → Extensions → MCP で設定:

```json
{
  "render": {
    "command": "node",
    "args": [
      "/absolute/path/to/financial-planning-calculator/scripts/mcp-server-render.js"
    ],
    "env": {
      "RENDER_API_KEY": "rnd_xxxxxxxxxxxxxxxxxxxxx"
    }
  }
}
```

### GitHub Copilot用

プロジェクトルートの`.mcp/config.json`を使用（既に設定済み）:

```json
{
  "mcpServers": {
    "render": {
      "command": "node",
      "args": ["./scripts/mcp-server-render.js"],
      "env": {
        "RENDER_API_KEY": "${RENDER_API_KEY}",
        "RENDER_OWNER_ID": "${RENDER_OWNER_ID}"
      }
    }
  }
}
```

## 5. 動作確認

### コマンドラインでの確認

```bash
# デプロイメントの状態をチェック
cd scripts
RENDER_API_KEY=xxx node check-render-deployments.js
```

成功すると、以下のような出力が表示されます:

```
🔍 Checking Render.com deployments...

📦 Found 3 services

📋 Service: financial-planning-backend
   Type: web
   Status: available
   📊 Recent deployments: 3
   
   Latest Deploy:
   - ID: dep-xxxxxxxxxxxxx
   - Status: live
   - Created: 1/23/2026, 10:30:00 AM
   ✅ Deployment is live and healthy

...

✅ All deployments are healthy
```

### AI Assistantでの確認

Claude、Copilot、またはCursorで以下のように質問:

```
Render.comのデプロイメント状態を確認して
```

または

```
最新のデプロイメントにエラーがないかチェックして
```

MCPサーバーが正しく設定されていれば、AIアシスタントがRender.comのAPIを使用して情報を取得し、回答します。

## 6. 利用方法

### 基本的な使い方

MCPツールが有効になると、AIアシスタントに以下のような質問ができます:

1. **サービス一覧の確認**
   ```
   Render.comのサービス一覧を表示して
   ```

2. **特定サービスの状態確認**
   ```
   financial-planning-backendの状態を確認して
   ```

3. **デプロイ履歴の確認**
   ```
   直近5回のデプロイメント履歴を見せて
   ```

4. **エラー検知**
   ```
   最新のデプロイメントログからエラーを検出して
   ```

5. **自動修正の依頼**
   ```
   デプロイエラーを分析して、修正方法を提案して
   ```

### GitHub Actions連携

プルリクエスト作成時、またはスケジュール実行時に自動的にデプロイメントをチェックします:

- **PR作成時**: エラーがあればPRにコメントを追加
- **定期実行（6時間毎）**: 本番環境でエラーがあればIssueを作成

## 7. 利用可能なMCPツール

MCPサーバーは以下のツールを提供します:

### `list_services`

全サービスの一覧を取得

```javascript
// 使用例（AIアシスタント経由）
"Render.comのサービス一覧を表示"
```

### `get_service_status`

特定サービスの状態を取得

```javascript
// パラメータ
{
  "serviceName": "financial-planning-backend"
}
```

### `list_recent_deploys`

最近のデプロイメント一覧を取得

```javascript
// パラメータ
{
  "serviceName": "financial-planning-backend",
  "limit": 10  // オプション、デフォルト10
}
```

### `get_deployment_logs`

デプロイメントログを取得

```javascript
// パラメータ
{
  "serviceName": "financial-planning-backend",
  "deployId": "dep-xxxxx"  // オプション、省略時は最新
}
```

### `detect_errors`

ログからエラーを自動検出

```javascript
// パラメータ
{
  "serviceName": "financial-planning-backend",
  "deployId": "dep-xxxxx"  // オプション、省略時は最新
}
```

## 8. エラーパターン

MCPサーバーは以下のエラーパターンを自動検出します:

| パターン | タイプ | 深刻度 |
|---------|--------|--------|
| `error`, `failed`, `failure` | general_error | high |
| `npm ERR!` | npm_error | high |
| `fatal` | fatal_error | critical |
| `cannot find module` | missing_dependency | high |
| `ECONNREFUSED`, `ETIMEDOUT` | connection_error | high |
| `syntax error` | syntax_error | high |
| `port already in use` | port_conflict | medium |
| `out of memory`, `OOM` | memory_error | critical |

## 9. トラブルシューティング

### MCPサーバーが起動しない

1. Node.jsのバージョン確認: `node --version`（18.0.0以上が必要）
2. 依存関係のインストール: `cd scripts && npm install`
3. APIキーの確認: `echo $RENDER_API_KEY`

### APIキーエラー

```
Error: API request failed with status 401
```

- APIキーが正しいか確認
- APIキーが有効期限内か確認
- Render.comダッシュボードで新しいキーを生成

### サービスが見つからない

```
Service "xxx" not found
```

- サービス名が正確か確認（大文字小文字も区別される）
- Owner IDが正しいか確認
- `list_services`ツールで正確なサービス名を確認

### ログが取得できない

```
Could not fetch logs
```

- デプロイメントが完了しているか確認
- APIキーに適切な権限があるか確認
- しばらく待ってから再試行

## 10. セキュリティ上の注意

1. **APIキーの管理**
   - APIキーをGitにコミットしない
   - `.gitignore`に環境変数ファイルを追加
   - GitHub Secretsを使用して安全に管理

2. **権限の最小化**
   - Render.com APIキーは必要最小限の権限で作成
   - 定期的にキーをローテーション

3. **ログの取り扱い**
   - ログに機密情報が含まれる可能性があるため注意
   - 公開リポジトリでは特に注意

## 11. 次のステップ

- [GitHub Actions設定ガイド](./GITHUB_ACTIONS_RENDER_MONITORING.md)
- [エラー自動修正のベストプラクティス](./ERROR_AUTO_FIX_GUIDE.md)
- [Render.com API ドキュメント](https://api-docs.render.com/)

## サポート

質問や問題がある場合は、GitHubのIssueを作成してください。

---

**更新日**: 2026-01-23
**バージョン**: 1.0.0
