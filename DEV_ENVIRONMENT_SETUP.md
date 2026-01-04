# Dev環境セットアップガイド

このガイドでは、PR作成時にクラウド上で自動的に確認できるdev環境の構築方法を説明します。

## 概要

このプロジェクトでは、Render.comを使用してPR作成時に自動的にプレビュー環境をデプロイします。

### デプロイされるサービス

1. **Frontend**: Next.js アプリケーション
2. **Backend**: Go API サーバー
3. **Database**: PostgreSQL データベース

## Render.comのセットアップ

### 1. Render.comアカウントの作成

1. [Render.com](https://render.com/)にアクセス
2. GitHubアカウントでサインアップ
3. GitHubリポジトリへのアクセスを許可

### 2. Blueprintからのデプロイ

Render.comは`render.yaml`ファイルを使用して自動的にサービスをセットアップします。

1. Render.comダッシュボードで「New +」→「Blueprint」を選択
2. GitHubリポジトリ `zakisanbaiman/financial-planning-calculator` を選択
3. ブランチ: `main` を選択
4. Render.comが自動的に`render.yaml`を検出し、以下のサービスを作成:
   - PostgreSQLデータベース
   - バックエンドAPI
   - フロントエンドWebアプリ

### 3. プレビュー環境の有効化

`render.yaml`には既にプレビュー環境の設定が含まれています:

```yaml
previewsEnabled: true
previewsExpireAfterDays: 7
```

これにより:
- PRが作成されると自動的にプレビュー環境がデプロイされます
- PRごとに独立した環境が作成されます
- プレビュー環境はPRクローズ後7日で自動削除されます

### 4. 環境変数の設定

#### 自動設定される環境変数

Render.comは`render.yaml`で定義された環境変数を自動的に設定します:

- **データベース接続情報**: 自動的にバックエンドに注入されます
  - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`

#### 手動設定が必要な環境変数

フロントエンドの`NEXT_PUBLIC_API_URL`は、バックエンドのデプロイ後に手動で設定する必要があります:

1. Render.comダッシュボードで `financial-planning-frontend` サービスを開く
2. 「Environment」タブを選択
3. `NEXT_PUBLIC_API_URL` を追加:
   ```
   https://financial-planning-backend.onrender.com
   ```
   （バックエンドのURLに置き換えてください）
4. 「Save Changes」をクリック

**注**: プレビュー環境では、各PRごとに異なるURLになるため、同様に設定が必要です。

## プレビュー環境の利用方法

### PRを作成すると

1. Render.comが自動的にプレビュー環境のデプロイを開始
2. GitHub PRページにRender.comからのコメントが追加され、プレビューURLが表示される
3. デプロイ完了後、プレビューURLにアクセスして変更を確認

### プレビューURL

各PRには以下のようなURLが割り当てられます:
- Frontend: `https://financial-planning-frontend-pr-{PR番号}.onrender.com`
- Backend: `https://financial-planning-backend-pr-{PR番号}.onrender.com`

### ログの確認

1. Render.comダッシュボードにログイン
2. 該当するプレビュー環境を選択
3. 「Logs」タブでリアルタイムログを確認

## トラブルシューティング

### ビルドが失敗する場合

1. Render.comのログを確認
2. Dockerfileが正しく設定されているか確認
3. 必要な環境変数が設定されているか確認

### データベース接続エラー

1. `render.yaml`のデータベース設定を確認
2. バックエンドの環境変数が正しく設定されているか確認
3. SSL接続が有効になっているか確認（`DB_SSLMODE=require`）

### フロントエンドがバックエンドに接続できない

1. `NEXT_PUBLIC_API_URL`が正しく設定されているか確認
2. バックエンドのヘルスチェックが成功しているか確認
3. CORSの設定を確認

## ローカル開発との違い

### データベース

- ローカル: PostgreSQL 15 (Docker)
- Dev環境: PostgreSQL (Render.com managed)

### 環境変数

- ローカル: `.env.docker`
- Dev環境: `render.yaml`で定義、Render.comダッシュボードで管理

### ビルド方法

- ローカル: Docker Compose with hot reload
- Dev環境: Production builds with Dockerfile

## 代替プラットフォーム

Render.com以外にも以下のプラットフォームが利用可能です:

### Railway.app

- 長所: 高速デプロイ、シンプルな設定
- 短所: 無料プランの制限
- 設定ファイル: `railway.json`または`railway.toml`

### Fly.io

- 長所: グローバルエッジデプロイ、高パフォーマンス
- 短所: やや複雑な設定
- 設定ファイル: `fly.toml`

### Vercel (フロントエンドのみ)

- 長所: Next.jsに最適化、自動PR プレビュー
- 短所: バックエンドは別途デプロイ必要
- 設定ファイル: `vercel.json`

## 料金について

### Render.com 無料プラン

- Web Services: 750時間/月
- データベース: 90日間無料
- 制限:
  - 非アクティブ時のスリープ
  - 月間ビルド時間制限

### 有料プランへのアップグレード

本番環境では以下をお勧めします:
- Starter Plan: $7/月 (Web Service)
- Standard Plan: $25/月 (PostgreSQL)

## セキュリティ

### 環境変数

- シークレット情報はRender.comダッシュボードで管理
- `render.yaml`には機密情報を含めない
- 環境変数は暗号化されて保存

### データベース

- SSL接続を必須に設定 (`DB_SSLMODE=require`)
- データベースアクセスはプライベートネットワーク内に制限

### CORS設定

バックエンドのCORS設定でフロントエンドのURLを許可:

```go
origins := []string{
    "https://financial-planning-frontend.onrender.com",
    "https://financial-planning-frontend-pr-*.onrender.com",
}
```

## まとめ

このセットアップにより:

✅ PRごとに自動的にプレビュー環境が作成される
✅ レビュアーがブラウザで変更を確認できる
✅ 本番環境と同じ構成でテストできる
✅ PRクローズ後は自動的にクリーンアップされる

詳細な設定は`render.yaml`ファイルを参照してください。
