# プレビュー環境クイックリファレンス

## 🚀 自動デプロイ

PRを作成すると、以下が自動的に実行されます：

1. **GitHub Actions** が起動し、PRにプレビューURLをコメント
2. **Render.com** が自動的に以下をデプロイ:
   - PostgreSQLデータベース（独立したインスタンス）
   - バックエンドAPI
   - フロントエンドWebアプリ

## 📋 プレビューURL形式

```
Frontend: https://financial-planning-frontend-pr-{PR番号}.onrender.com
Backend:  https://financial-planning-backend-pr-{PR番号}.onrender.com
API Docs: https://financial-planning-backend-pr-{PR番号}.onrender.com/swagger/index.html
```

## ⚙️ 初回セットアップ（リポジトリオーナーのみ）

### 1. Render.comアカウント作成

```bash
1. https://render.com にアクセス
2. GitHubアカウントでサインアップ
3. リポジトリへのアクセスを許可
```

### 2. Blueprintデプロイ

```bash
1. Render.comダッシュボード → "New +" → "Blueprint"
2. リポジトリ選択: zakisanbaiman/financial-planning-calculator
3. ブランチ: main
4. "Apply" をクリック
```

### 3. 環境変数設定

フロントエンドサービスに以下を手動設定:

```
NEXT_PUBLIC_API_URL=https://financial-planning-backend.onrender.com
```

### 4. プレビュー環境有効化

Render.comダッシュボードで:

```
Settings → Pull Request Previews → Enable
```

## 🔍 トラブルシューティング

### デプロイが失敗する

```bash
# Render.comダッシュボードでログ確認
1. 該当サービスを選択
2. "Logs" タブを開く
3. エラーメッセージを確認
```

### 環境変数が反映されない

```bash
# 環境変数を再確認
1. サービスの "Environment" タブ
2. 変数を確認・更新
3. "Save Changes" 後に自動再デプロイ
```

### プレビューが作成されない

```bash
# 以下を確認:
1. Render.comでPull Request Previewsが有効か
2. render.yamlが正しくコミットされているか
3. GitHub Actionsが正常に動作しているか
```

## 📝 制限事項

### 無料プラン

- **自動スリープ**: 15分間非アクティブでスリープ
- **起動時間**: 初回アクセス時に数秒かかる
- **ビルド時間**: 月間制限あり
- **保存期間**: PRクローズ後7日間

### 対処法

- 本番環境には有料プラン推奨
- 開発中はローカル環境を優先使用

## 🔗 便利なリンク

- [セットアップガイド](./SETUP.md)
- [Render.com Dashboard](https://dashboard.render.com)
- [Render.com Docs](https://render.com/docs)

## 💡 ベストプラクティス

1. **ローカルテスト優先**: プレビュー前にローカルで動作確認
2. **小さなPR**: 大きな変更は複数PRに分割
3. **環境変数確認**: デプロイ前に必要な変数が設定されているか確認
4. **ログ監視**: デプロイ後はログを確認して問題がないか確認
