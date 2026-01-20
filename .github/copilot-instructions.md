# GitHub Copilot カスタムインストラクション

## トークン使用量の表示

**重要**: 会話の終わりや大きなタスク完了時には、必ずトークン使用量を報告すること。

フォーマット（右端に小さく表示）:
```
[Token: XX,XXX/1,000,000 (XX%)]
```

## プロジェクト固有のルール

### アーキテクチャ
- Clean Architecture / DDD（ドメイン駆動設計）を採用
- 層構造:
  - `application/usecases/` - ユースケース層
  - `domain/` - ドメイン層
  - `infrastructure/` - インフラ層（web/controllersを含む）

### 命名規則
- Controller: HTTPハンドラーをまとめた構造体（Echoでは本来Handlerだが、このプロジェクトではController）
- Repository: データ永続化の抽象化
- UseCase: ビジネスロジックの実装

### フレームワーク
- バックエンド: Echo v4
- フロントエンド: Next.js
- データベース: PostgreSQL

### コミット規約
- feat: 新機能
- fix: バグ修正
- refactor: リファクタリング
- docs: ドキュメント更新
- 必ず Issue 番号を含める（例: Issue: #66）

### ブランチ保護ルール
- **mainブランチへの直接pushは禁止**（リポジトリルールで保護）
- 変更はフィーチャーブランチを作成してPRを経由すること
- ブランチ命名例: `fix/oauth-auth-state`, `feat/new-feature`, `refactor/cleanup`
