# MCP Render.com連携 クイックリファレンス

## セットアップ（1分）

```bash
# 1. 依存関係のインストール
cd scripts && npm install

# 2. 環境変数の設定
export RENDER_API_KEY="rnd_xxxxxxxxxxxxx"

# 3. 動作確認
node check-render-deployments.js
```

## よく使うコマンド

```bash
# デプロイメント状態をチェック
node scripts/check-render-deployments.js

# MCPサーバーを単体で起動（デバッグ用）
node scripts/mcp-server-render.js
```

## AIアシスタントでの使い方

### Claude / Copilot / Cursor

```
# エラーチェック
"Render.comのデプロイメントエラーをチェックして"

# 自動修正
"検出されたエラーを分析して修正方法を提案して"

# サービス状態確認
"financial-planning-backendの状態を確認"

# ログ分析
"最新のデプロイメントログを分析してエラーの原因を特定して"
```

## エラー対応フロー

```
1. エラー検知（自動）
   ↓
2. AIアシスタントに相談
   "Render.comのエラーを分析して修正案を提案して"
   ↓
3. 修正実施
   ↓
4. コミット & プッシュ
   ↓
5. 自動再デプロイ
   ↓
6. 確認
```

## GitHub Actions連携

### 自動実行タイミング

- ✅ PR作成時（preview環境デプロイ後）
- ✅ 6時間ごと（本番環境チェック）
- ✅ 手動実行可能

### PR作成時の自動コメント例

```markdown
## ⚠️ Render.com Deployment Issues Detected

エラーが検出されました。詳細:
- npm install failed
- Missing dependency: @types/node

推奨アクション:
1. package.jsonを確認
2. 依存関係を追加
3. 再デプロイ
```

## トラブルシューティング（30秒）

| 問題 | 解決方法 |
|------|---------|
| APIキーエラー | `export RENDER_API_KEY="..."` を確認 |
| サービス not found | `list_services`で正確な名前を確認 |
| 依存関係エラー | `cd scripts && npm install` |
| MCPサーバー起動失敗 | Node.js 18以上を確認 |

## 設定ファイル

```
project/
├── .mcp/
│   └── config.json          # MCP設定
├── scripts/
│   ├── package.json         # 依存関係
│   ├── mcp-server-render.js # MCPサーバー
│   └── check-render-deployments.js # チェックスクリプト
└── .github/workflows/
    └── monitor-render-deployments.yml # 自動監視
```

## リンク

- 📖 [詳細セットアップガイド](./MCP_SETUP.md)
- 🔧 [Render Dashboard](https://dashboard.render.com)
- 📚 [MCP公式ドキュメント](https://modelcontextprotocol.io/)
- 🚀 [Render API Docs](https://api-docs.render.com/)

---

**困ったら**: Issue作成 → `@copilot` にメンション
