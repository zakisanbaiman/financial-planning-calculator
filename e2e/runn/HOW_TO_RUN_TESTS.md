# 財務データ登録と計算テストの実行方法

このガイドでは、`financial-data-registration-and-calculation.yml` テストを実行する方法を説明します。

## 前提条件

### 1. runn のインストール

```bash
# Go を使用してインストール
go install github.com/k1LoW/runn/cmd/runn@latest

# または Homebrew を使用（macOS）
brew install k1LoW/tap/runn

# または直接ダウンロード
# https://github.com/k1LoW/runn/releases から最新版をダウンロード
```

### 2. バックエンドサーバーとデータベースの起動

Dockerを使用する場合：

```bash
# プロジェクトルートディレクトリで実行
cd /path/to/financial-planning-calculator

# 初回セットアップ（ビルド、マイグレーション、シード）
make dev-setup

# または既に環境がセットアップ済みの場合
make up
```

ローカル環境で実行する場合：

```bash
# データベースが起動していることを確認
# PostgreSQL が localhost:5432 で実行されている必要があります

# バックエンドサーバーを起動
cd backend
go run main.go
```

### 3. 環境変数の設定

```bash
# API URL を設定（デフォルトは http://localhost:8080/api）
export API_URL=http://localhost:8080/api

# ランID を設定（テストユーザーIDの一意性を保証）
export RUN_ID=$(date +%s)
```

## テストの実行

### 基本的な実行

```bash
cd e2e/runn
runn run financial-data-registration-and-calculation.yml
```

### 詳細な出力付きで実行

```bash
runn run --verbose financial-data-registration-and-calculation.yml
```

### HTMLレポートを生成

```bash
runn run -o result.html financial-data-registration-and-calculation.yml
```

### JSON形式で結果を出力

```bash
runn run --format json financial-data-registration-and-calculation.yml
```

## テストが検証する内容

このテストは以下の一連のフローを検証します：

1. **財務データの作成**
   - 月収: 500,000円
   - 月間支出: 住居費、食費、光熱費、通信費、交通費
   - 現在の貯蓄: 預金 2,000,000円、投資 1,000,000円
   - 投資利回り: 5.0%
   - インフレ率: 2.0%

2. **データの取得と検証**
   - 作成したデータが正しく保存されているか確認
   - 全てのフィールドが期待値と一致するか検証

3. **退職データの追加**
   - 退職年齢: 65歳
   - 退職後月間支出: 250,000円
   - 年金受給額: 150,000円

4. **緊急資金設定の追加**
   - 目標月数: 6ヶ月
   - 現在額: 500,000円

5. **計算機能のテスト**
   - 資産推移計算（10年間予測）
   - 退職資金計算
   - 緊急資金計算
   - 包括的予測計算（30年間予測）

6. **データの永続性確認**
   - 計算実行後もデータが保持されているか確認

7. **クリーンアップ**
   - テストデータの削除
   - 削除の確認

## トラブルシューティング

### エラー: Connection refused

バックエンドサーバーが起動していることを確認してください：

```bash
# ヘルスチェックエンドポイントにアクセス
curl http://localhost:8080/health
```

正常に動作している場合、以下のようなレスポンスが返ります：

```json
{
  "status": "ok",
  "message": "財務計画計算機 API サーバーが正常に動作しています",
  "timestamp": "2024-01-05T12:00:00Z",
  "version": "1.0.0"
}
```

### エラー: Database connection failed

データベースが起動していることを確認してください：

```bash
# Docker を使用している場合
docker ps | grep postgres

# ローカル環境の場合
psql -U postgres -d financial_planning -c "SELECT 1"
```

### エラー: 404 Not Found

- API_URL 環境変数が正しく設定されているか確認
- エンドポイントパスが正しいか確認（`/api/financial-data` など）

### テストの一部が失敗する場合

1. **データが既に存在する場合**
   - 別の RUN_ID を使用してテストを再実行
   - または既存のテストデータを手動で削除

2. **計算結果の検証が失敗する場合**
   - バックエンドの計算ロジックが最新の実装と一致するか確認
   - テストのアサーション条件を確認

## 参考リンク

- [runn ドキュメント](https://github.com/k1LoW/runn)
- [プロジェクト README](../../README.md)
- [Docker セットアップガイド](../../DOCKER_SETUP.md)
