# Render.com MCP統合 実装概要

## 実装完了日
2026-01-23

## 背景

Render.comにデプロイしているが、その際のエラーを検知できずフィードバックサイクルがうまく回っていない問題を解決するため、MCP (Model Context Protocol)を使用してAIアシスタント（Copilot、Claude、Cursor）からRender.comのデプロイ状態を監視し、エラーを自動検知・修正できる仕組みを実装。

## 実装内容

### 1. MCPサーバー実装

**ファイル**: `scripts/mcp-server-render.js`

Render.com APIと連携するMCPサーバー。以下の機能を提供:

- `list_services`: 全サービスの一覧取得
- `get_service_status`: 特定サービスの状態確認
- `list_recent_deploys`: デプロイ履歴の取得
- `get_deployment_logs`: デプロイメントログの取得
- `detect_errors`: ログからの自動エラー検出（8つのパターンに対応）

### 2. デプロイメント監視スクリプト

**ファイル**: `scripts/check-render-deployments.js`

コマンドラインおよびGitHub Actionsから実行可能なスクリプト:

- 全サービスのデプロイ状態をチェック
- エラーパターンの自動検出
- 詳細なレポート生成
- 終了コードによる成功/失敗の判定

### 3. GitHub Actions ワークフロー

**ファイル**: `.github/workflows/monitor-render-deployments.yml`

自動デプロイ監視ワークフロー:

**トリガー**:
- PR作成後（preview環境デプロイ完了後）
- 定期実行（6時間ごと）
- 手動実行

**機能**:
- デプロイエラーの自動検出
- PRへのエラー情報コメント追加
- 本番環境エラー時のIssue自動作成
- デプロイログのサマリー生成

### 4. MCP設定ファイル

**ファイル**: `.mcp/config.json`

Claude Desktop、Copilot、Cursor用のMCP設定:

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

### 5. ドキュメント

以下の包括的なドキュメントを作成:

1. **`docs/MCP_SETUP.md`** (9.2KB)
   - 詳細なセットアップ手順
   - Render.com APIキーの取得方法
   - 各種AIツールの設定方法
   - トラブルシューティング

2. **`docs/MCP_QUICK_REFERENCE.md`** (2.8KB)
   - 1分でわかるクイックガイド
   - よく使うコマンド
   - エラー対応フロー

3. **`docs/MCP_USAGE_EXAMPLES.md`** (6.7KB)
   - 5つの実践的シナリオ
   - AIへの質問例
   - チーム活用方法

4. **`docs/AI_COPILOT_INTEGRATION.md`** (8.5KB)
   - Claude Desktop設定
   - GitHub Copilot設定
   - Cursor設定
   - 実際の使用例とFAQ

5. **`scripts/README.md`** (2.3KB)
   - スクリプトの使用方法
   - 環境変数の設定
   - トラブルシューティング

### 6. その他の変更

- **README.md**: MCP統合についてのセクション追加
- **.gitignore**: MCPログファイルとnode_modulesを除外
- **scripts/package.json**: MCP SDKの依存関係定義

## エラー検出パターン

以下の8つのエラーパターンを自動検出:

| パターン | タイプ | 深刻度 |
|---------|--------|--------|
| error/failed/failure | general_error | high |
| npm ERR! | npm_error | high |
| fatal | fatal_error | critical |
| cannot find module | missing_dependency | high |
| ECONNREFUSED/ETIMEDOUT | connection_error | high |
| syntax error | syntax_error | high |
| port already in use | port_conflict | medium |
| out of memory/OOM | memory_error | critical |

## 技術スタック

- **言語**: Node.js (JavaScript/CommonJS)
- **プロトコル**: MCP (Model Context Protocol)
- **API**: Render.com REST API
- **CI/CD**: GitHub Actions
- **AIツール対応**: Claude Desktop, GitHub Copilot, Cursor

## セットアップ手順

### 1. 依存関係のインストール

```bash
cd scripts
npm install
```

### 2. 環境変数の設定

```bash
export RENDER_API_KEY="rnd_xxxxxxxxxxxxx"
export RENDER_OWNER_ID="dsy-xxxxxxxxxxxxx"  # オプション
```

### 3. 動作確認

```bash
node scripts/check-render-deployments.js
```

### 4. AIツールの設定

各AIツールの設定ファイルに`.mcp/config.json`の内容を追加:

- **Claude Desktop**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Cursor**: Settings → MCP
- **Copilot**: プロジェクト内の`.mcp/config.json`を使用

### 5. GitHub Secretsの設定

リポジトリの Settings → Secrets に以下を追加:

- `RENDER_API_KEY`: Render.com APIキー
- `RENDER_OWNER_ID`: Owner ID（オプション）

## 使用方法

### コマンドライン

```bash
# デプロイメント状態をチェック
RENDER_API_KEY=xxx node scripts/check-render-deployments.js
```

### AIアシスタント

```
Render.comのデプロイエラーをチェックして
```

### GitHub Actions

PR作成時、または定期実行で自動的に実行されます。

## 効果

1. **フィードバックサイクルの改善**
   - デプロイエラーの即座の検知
   - 自動通知によるタイムリーな対応

2. **AI連携による効率化**
   - エラー原因の自動分析
   - 修正案の即座の提示
   - 人手によるログ調査の削減

3. **チーム全体での可視化**
   - 誰でもデプロイ状態を確認可能
   - 本番環境の継続的な監視
   - 問題の早期発見

## ファイル構成

```
financial-planning-calculator/
├── .github/workflows/
│   └── monitor-render-deployments.yml   # 自動監視ワークフロー
├── .mcp/
│   └── config.json                      # MCP設定
├── docs/
│   ├── AI_COPILOT_INTEGRATION.md        # AI統合ガイド
│   ├── MCP_QUICK_REFERENCE.md           # クイックリファレンス
│   ├── MCP_SETUP.md                     # セットアップガイド
│   └── MCP_USAGE_EXAMPLES.md            # 使用例
├── scripts/
│   ├── README.md                        # スクリプト説明
│   ├── check-render-deployments.js      # 監視スクリプト
│   ├── mcp-server-render.js             # MCPサーバー
│   └── package.json                     # 依存関係
├── .gitignore                           # MCPログを除外
└── README.md                            # MCP情報追加
```

## セキュリティ考慮事項

1. **APIキーの管理**
   - 環境変数で管理
   - GitHubにはコミットしない
   - GitHub Secretsで安全に保管

2. **権限の最小化**
   - Read-OnlyのAPIキーを推奨
   - 定期的なローテーション

3. **ログの取り扱い**
   - 機密情報を含む可能性があるため注意
   - 公開リポジトリでは特に注意

## 今後の拡張案

1. **より高度なエラー分析**
   - 機械学習によるパターン学習
   - 過去のエラーからの予測

2. **自動修正機能**
   - 一般的なエラーの自動修正PR作成
   - 設定ファイルの自動更新

3. **他サービスとの連携**
   - Slack通知
   - PagerDuty連携
   - メトリクス収集

4. **パフォーマンス監視**
   - ビルド時間の追跡
   - リソース使用量の監視
   - 傾向分析

## まとめ

MCP統合により、Render.comデプロイメントの監視とエラー検知が自動化され、AIアシスタントを活用した効率的な問題解決が可能になりました。これにより、フィードバックサイクルが大幅に改善され、開発生産性の向上が期待できます。

---

**実装者**: GitHub Copilot Agent
**レビュー**: 要確認（実際のRender APIキーでのテストが必要）
**関連Issue**: renderとMCP接続する
