# AI Copilot統合ガイド

このドキュメントでは、GitHub Copilot、Claude Code、Cursorなど、各種AIアシスタントとRender.com MCP統合を使用する方法を説明します。

## サポートされているAIツール

### ✅ 完全対応
- **Claude Desktop** (Anthropic) - MCPネイティブサポート
- **GitHub Copilot** (VS Code) - MCP拡張機能経由
- **Cursor** - MCP統合機能

### 🔄 設定が必要
- **VS Code Copilot Chat** - カスタム設定が必要
- **その他のAIツール** - HTTP API経由での統合可能

## 各ツールの設定方法

### Claude Desktop

**最も推奨される方法** - MCPのネイティブサポート

1. **設定ファイルの場所**
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`
   - Linux: `~/.config/Claude/claude_desktop_config.json`

2. **設定内容**
   ```json
   {
     "mcpServers": {
       "render": {
         "command": "node",
         "args": [
           "/absolute/path/to/financial-planning-calculator/scripts/mcp-server-render.js"
         ],
         "env": {
           "RENDER_API_KEY": "rnd_xxxxxxxxxxxxx",
           "RENDER_OWNER_ID": "dsy-xxxxxxxxxxxxx"
         }
       }
     }
   }
   ```

3. **使用方法**
   
   Claude Desktopを再起動後、チャットで以下のように質問:
   ```
   Render.comのデプロイ状態を確認して
   ```

4. **確認方法**
   
   MCPツールが有効になっていると、Claudeの回答に「🔧 Using tools:」という表示が出ます。

### GitHub Copilot (VS Code)

**設定方法A: プロジェクト内設定（推奨）**

既にプロジェクトに`.mcp/config.json`が含まれているため、追加設定は不要です。

1. **環境変数の設定**
   ```bash
   # .zshrc または .bashrc に追加
   export RENDER_API_KEY="rnd_xxxxxxxxxxxxx"
   export RENDER_OWNER_ID="dsy-xxxxxxxxxxxxx"
   ```

2. **VS Codeで確認**
   ```
   Copilot Chat で質問:
   @workspace Render.comのデプロイエラーをチェック
   ```

**設定方法B: グローバル設定**

1. VS Codeの設定 (`settings.json`) に追加:
   ```json
   {
     "github.copilot.mcp.servers": {
       "render": {
         "command": "node",
         "args": [
           "/absolute/path/to/scripts/mcp-server-render.js"
         ],
         "env": {
           "RENDER_API_KEY": "rnd_xxxxxxxxxxxxx"
         }
       }
     }
   }
   ```

### Cursor

1. **Cursor設定を開く**
   - Cursor → Settings → Features → Model Context Protocol

2. **MCP設定を追加**
   ```json
   {
     "mcpServers": {
       "render": {
         "command": "node",
         "args": [
           "/absolute/path/to/financial-planning-calculator/scripts/mcp-server-render.js"
         ],
         "env": {
           "RENDER_API_KEY": "rnd_xxxxxxxxxxxxx"
         }
       }
     }
   }
   ```

3. **使用方法**
   
   Cursorのチャットで:
   ```
   Render.comの最新デプロイを確認
   ```

### Zed Editor

Zedは将来的にMCPをサポートする予定です。現時点では以下の方法で統合可能:

1. **カスタムLanguage Server経由**
   
   MCPサーバーをLanguage Serverとしてラップして使用

2. **外部ツール連携**
   
   コマンドパレットからスクリプトを実行

## 実際の使用例

### ケース1: デプロイエラーの自動検出

**シナリオ**: PRをマージしたが、デプロイが失敗しているかもしれない

**Claude/Copilot/Cursorで質問**:
```
Render.comのデプロイ状態を確認して、エラーがあれば詳細を教えて
```

**期待される応答**:
```
Render.comのデプロイ状態を確認しました。

📦 Services: 3つ
✅ financial-planning-frontend: 正常
❌ financial-planning-backend: エラー検出

エラー詳細:
- タイプ: npm_error
- 深刻度: high
- 内容: npm ERR! missing: @types/node@^18.0.0

推奨される対応:
1. package.jsonに以下を追加:
   "devDependencies": {
     "@types/node": "^18.0.0"
   }
2. npm installを実行
3. コミット&プッシュ
```

### ケース2: ログの詳細分析

**質問**:
```
financial-planning-backendの最新デプロイメントログを分析して、
パフォーマンスの問題がないか確認して
```

**期待される応答**:
```
デプロイメントログを分析しました。

📊 ビルド時間: 2分34秒（前回比+15秒）
⚠️ 警告が2件見つかりました:

1. メモリ使用量が閾値に近い
   - 現在: 450MB / 512MB
   - 推奨: メモリ制限を768MBに増やすか、最適化を検討

2. 依存関係のインストール時間が長い
   - npm install: 45秒
   - 推奨: package-lock.jsonをコミットして高速化

重大な問題は検出されませんでした。
```

### ケース3: 複数サービスの比較

**質問**:
```
すべてのRender.comサービスの直近3回のデプロイを比較して、
安定性の傾向を教えて
```

**期待される応答**:
```
3サービスの直近3回のデプロイを分析:

📈 financial-planning-frontend
  - 成功率: 100% (3/3)
  - 平均ビルド時間: 1分20秒
  - 評価: ✅ 非常に安定

📈 financial-planning-backend
  - 成功率: 67% (2/3)
  - 平均ビルド時間: 2分45秒
  - 失敗理由: 環境変数の設定ミス
  - 評価: ⚠️ 要注意

📈 financial-planning-db
  - 成功率: 100% (3/3)
  - 評価: ✅ 安定

推奨アクション:
バックエンドの環境変数設定を見直してください。
```

## トラブルシューティング

### MCPツールが表示されない

1. **設定ファイルのパスを確認**
   ```bash
   # 絶対パスで指定されているか確認
   cat ~/.config/Claude/claude_desktop_config.json
   ```

2. **Node.jsのバージョン確認**
   ```bash
   node --version  # 18.0.0以上が必要
   ```

3. **APIキーの確認**
   ```bash
   echo $RENDER_API_KEY
   # または設定ファイル内で直接指定されているか確認
   ```

4. **MCPサーバーの動作確認**
   ```bash
   node scripts/mcp-server-render.js
   # エラーが出ないか確認
   ```

### 「Permission denied」エラー

```bash
chmod +x scripts/mcp-server-render.js
chmod +x scripts/check-render-deployments.js
```

### 「Cannot find module」エラー

```bash
cd scripts
npm install
```

### APIキーが無効

1. Render.comでAPIキーを再生成
2. 設定ファイルまたは環境変数を更新
3. AIツールを再起動

## ベストプラクティス

### 1. セキュリティ

- ❌ APIキーをコードにハードコーディングしない
- ✅ 環境変数または設定ファイルで管理
- ✅ 設定ファイルを`.gitignore`に追加

### 2. 効率的な質問

**良い質問**:
```
Render.comの最新デプロイでエラーがないか確認し、
あれば原因と修正方法を提案して
```

**悪い質問**:
```
エラーある？
```

### 3. 定期的なメンテナンス

- 月1回: APIキーのローテーション
- 週1回: デプロイログの確認
- 随時: MCPサーバーのアップデート

### 4. チームでの活用

1. **共有APIキー**
   - チーム用のRead-Only APIキーを作成
   - チームメンバー全員が同じ設定を使用

2. **ナレッジ共有**
   - よく使う質問をドキュメント化
   - エラーパターンと解決策を共有

3. **ロールベースのアクセス**
   - 開発者: 全サービスへのアクセス
   - レビュアー: 読み取り専用

## FAQ

### Q: MCPサーバーはどのくらいのリソースを消費しますか？

A: 非常に軽量です。メモリ使用量は約10-20MB、CPU使用率は質問時のみ一時的に上がります。

### Q: オフラインでも使用できますか？

A: いいえ、Render.com APIへのアクセスが必要なため、インターネット接続が必須です。

### Q: 複数のプロジェクトで使用できますか？

A: はい。各プロジェクトごとにMCPサーバーを設定し、異なるAPIキーを使用できます。

### Q: AIツール以外からも使用できますか？

A: はい。`check-render-deployments.js`スクリプトはコマンドラインから直接実行できます。

### Q: レート制限はありますか？

A: Render.com APIのレート制限に従います（通常は十分な余裕があります）。

## 次のステップ

1. ✅ MCPサーバーをセットアップ
2. ✅ AIツールで動作確認
3. 📚 [使用例ドキュメント](./MCP_USAGE_EXAMPLES.md)を読む
4. 🚀 実際のプロジェクトで活用開始

## サポート

問題や質問がある場合は、GitHubのIssueを作成するか、
チャットで `@copilot` をメンションしてください。

---

**更新日**: 2026-01-23
