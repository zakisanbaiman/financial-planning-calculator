# Render.com MCP連携の使用例

このドキュメントでは、Render.com MCP統合の実際の使用例を紹介します。

## シナリオ1: デプロイエラーの検出と修正

### 状況
PRをマージしたが、Render.comでのデプロイが失敗している。

### 手順

1. **エラー検知（自動）**
   
   GitHub Actionsが自動的にエラーを検出し、PRにコメント:
   
   ```markdown
   ⚠️ Render.com Deployment Issues Detected
   
   Error: npm install failed
   - Cannot find module '@types/node'
   ```

2. **AIアシスタントに相談**
   
   Claude/Copilot/Cursorで質問:
   ```
   Render.comのデプロイエラーを分析して、修正方法を提案してください
   ```

3. **AIの回答例**
   
   ```
   デプロイログを分析した結果、以下の問題が見つかりました:
   
   問題: @types/node パッケージが見つからない
   原因: package.jsonに依存関係が記載されていない
   
   修正方法:
   1. package.jsonに以下を追加:
      "devDependencies": {
        "@types/node": "^18.0.0"
      }
   
   2. ローカルで確認:
      npm install
      npm run type-check
   
   3. コミット & プッシュ:
      git add package.json
      git commit -m "fix: Add missing @types/node dependency"
      git push
   ```

4. **修正実施と検証**
   
   修正後、自動的に再デプロイが開始され、成功を確認。

## シナリオ2: 本番環境の定期監視

### 状況
本番環境で突然エラーが発生。

### 手順

1. **自動検出（6時間ごと）**
   
   GitHub Actionsが定期チェックを実行し、エラーを検出するとIssueを自動作成:
   
   ```markdown
   🚨 Render.com Production Deployment Errors Detected
   
   Error: Database connection timeout
   Severity: Critical
   ```

2. **即座に対応**
   
   Issueの通知を受けた開発者がAIアシスタントに相談:
   ```
   本番環境のRenderデプロイでデータベース接続エラーが発生しています。
   原因を特定して修正してください。
   ```

3. **AIによる診断と修正**
   
   AIがログを分析し、環境変数の設定ミスを発見し、修正案を提示。

## シナリオ3: プレビュー環境の確認

### 状況
新機能をPRで作成し、プレビュー環境で動作確認したい。

### 手順

1. **PR作成**
   ```bash
   git checkout -b feature/new-api
   # ... 変更を加える
   git push origin feature/new-api
   ```

2. **プレビュー環境の自動デプロイ**
   
   GitHub Actionsがプレビュー環境をデプロイし、URLをコメント:
   ```
   🚀 Preview Environment
   - Frontend: https://financial-planning-frontend-pr-123.onrender.com
   - Backend: https://financial-planning-backend-pr-123.onrender.com
   ```

3. **デプロイ状態の確認（AIアシスタント経由）**
   
   ```
   PR #123のプレビュー環境のデプロイ状態を確認して
   ```
   
   AIの回答:
   ```
   プレビュー環境のデプロイ状態:
   ✅ Frontend: 正常にデプロイ済み
   ✅ Backend: 正常にデプロイ済み
   
   利用可能なURL:
   - https://financial-planning-frontend-pr-123.onrender.com
   
   すべてのサービスが正常に稼働しています。
   ```

## シナリオ4: ローカルでのデバッグ

### 状況
ローカルでRender.comのデプロイ状態を確認したい。

### 手順

1. **コマンドラインで確認**
   ```bash
   cd scripts
   RENDER_API_KEY=xxx node check-render-deployments.js
   ```

2. **詳細なログ分析**
   ```bash
   # 特定のサービスのログを取得（MCPサーバー経由）
   # AIアシスタントで以下のように質問:
   financial-planning-backendの最新デプロイメントログを取得して、
   エラーがあれば詳細を教えて
   ```

3. **エラーパターンの学習**
   
   AIが過去のエラーから学習し、同じエラーを事前に防ぐ提案:
   ```
   過去のデプロイログを分析した結果、
   以下のパターンでエラーが頻発しています:
   
   1. npm install時のメモリ不足
      → package.jsonのdependenciesを最小化
   
   2. 環境変数の設定ミス
      → .env.exampleを最新に保つ
   
   3. データベースマイグレーションの失敗
      → start.shでマイグレーションを確実に実行
   ```

## シナリオ5: チーム全体での利用

### 状況
チームメンバー全員がRenderのデプロイ状態を把握したい。

### セットアップ

1. **チーム用のRender APIキーを作成**
   ```
   Render.com → Account Settings → API Keys
   キー名: "Team MCP Integration (Read Only)"
   ```

2. **GitHub Secretsに追加**
   ```
   Repository Settings → Secrets → Actions
   RENDER_API_KEY: (生成したキー)
   ```

3. **各メンバーがMCPを設定**
   
   各自のClaude/Copilot/Cursorに同じAPIキーを設定することで、
   全員が同じデプロイ情報にアクセス可能。

4. **チーム内での活用例**
   ```
   開発者A: 「最新のデプロイ、問題なさそう？」
   開発者B: 「AIに確認したら、バックエンドでwarningが3件あるって」
   開発者A: 「詳細見てみる」→ AIに詳細を依頼
   開発者C: 「修正PR出すね」
   ```

## よく使うAIへの質問例

### デプロイ確認
```
Render.comの全サービスの状態を確認して
```

### エラー検出
```
最新のデプロイでエラーがないかチェックして
```

### 詳細分析
```
financial-planning-backendの直近3回のデプロイを比較して、
パフォーマンスが悪化していないか確認して
```

### 自動修正
```
検出されたエラーを分析して、コードの修正案を提示して。
可能であれば自動修正もお願いします。
```

### プロアクティブな提案
```
過去のデプロイログを分析して、
今後発生しそうな問題を予測して対策を提案して
```

## まとめ

MCP統合により:
- ✅ デプロイエラーを即座に検知
- ✅ AIが自動的に原因を分析
- ✅ 修正案を即座に提示
- ✅ チーム全体でデプロイ状態を共有
- ✅ 本番環境の継続的な監視

これにより、**フィードバックサイクルが大幅に改善**され、問題の早期発見・早期解決が可能になります。

---

**関連ドキュメント**:
- [MCPセットアップガイド](./MCP_SETUP.md)
- [クイックリファレンス](./MCP_QUICK_REFERENCE.md)
