# runnを使用したAPIテスト

このディレクトリには`runn`を使用したAPI安定性テストが含まれています。

## 必要な準備

### 1. runnのインストール

```bash
# macOSの場合
brew install runnme/tap/runn

# またはGitHubから直接ダウンロード
# https://github.com/k1LoW/runn/releases
```

### 2. 環境変数の設定

```bash
export API_URL=http://localhost:8080/api
export RUN_ID=$(date +%s)
```

## テストファイル

### 1. `user-onboarding-flow.yml`
ユーザーオンボーディングの完全なフロー（CRUD操作）をテストします：
- 財務データの作成・取得・更新
- ゴール（目標）の作成・取得・更新・削除
- データの永続化確認
- 並行ゴール作成（5個）

**実行方法：**
```bash
runn run user-onboarding-flow.yml
```

### 2. `api-stability.yml`
API安定性・信頼性を集中的にテストします：
- 10個の同時ゴール作成リクエスト（並行処理テスト）
- データ一貫性の検証（5回の読み込み）
- 複数回の更新テスト
- エラーハンドリング（無効なID、バリデーション）
- ビジネスロジック検証（負数、過去の日付）
- 削除と削除確認

**実行方法：**
```bash
runn run api-stability.yml
```

## すべてのテストを実行

```bash
# 環境変数を設定してから実行
export API_URL=http://localhost:8080/api
export RUN_ID=$(date +%s)

# すべてのYMLテストを実行
runn run --glob "*.yml"

# または個別に実行
runn run user-onboarding-flow.yml && runn run api-stability.yml
```

## レポート出力

```bash
# HTML形式のレポートを生成
runn run -o result.html user-onboarding-flow.yml

# JSON形式で結果を出力
runn run --format json user-onboarding-flow.yml
```

## トラブルシューティング

### APIが応答しない場合
```bash
# APIサーバーが起動しているか確認
curl http://localhost:8080/api/health

# バックエンドを起動
cd ../../backend
make run
```

### runnコマンドが見つからない場合
```bash
# インストール確認
which runn

# パスが通っていない場合は絶対パスで実行
/usr/local/bin/runn run user-onboarding-flow.yml
```

## テスト結果の解釈

- ✓（チェックマーク）: テスト成功
- ✗（バツ印）: テスト失敗
- Status: HTTPステータスコード
- Assert: 検証条件の結果

各テストステップが成功したら、APIは安定して動作しています。

## 参考リンク

- [runnドキュメント](https://k1low.dev/runn/)
- [runnGitHub](https://github.com/k1LoW/runn)
