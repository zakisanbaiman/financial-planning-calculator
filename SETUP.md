# セットアップガイド

## クイックスタート

```bash
# 1. リポジトリをクローン
git clone <repository-url>
cd financial-planning-calculator

# 2. 依存関係をインストール
make install

# 3. Git hooksをセットアップ
make setup

# 4. 開発サーバーを起動
make dev
```

これで以下にアクセスできます：

- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

## Git Hooksについて

`make setup`を実行すると、以下のGit hooksが自動的に設定されます：

### pre-commit（コミット前）

コミット前に自動的に以下を実行：

1. **Goコード**
   - `gofmt -w` でフォーマット
   - `go vet` でチェック

2. **TypeScript/JavaScript**
   - ESLintで自動修正
   - 型チェック

3. **JSON/YAML/Markdown**
   - Prettierでフォーマット

### commit-msg（コミットメッセージ）

コミットメッセージが以下の形式に従っているかチェック：

```
<type>(<scope>): <subject>
```

例：

- ✅ `feat(frontend): 新しいチャートコンポーネントを追加`
- ✅ `fix(backend): 計算ロジックのバグを修正`
- ✅ `docs: READMEを更新`
- ❌ `update code` （形式が不正）

## 使用可能なコマンド

### 開発

```bash
make dev              # フロントエンド + バックエンドを起動
make backend          # バックエンドのみ起動
make frontend         # フロントエンドのみ起動
```

### コード品質

```bash
make lint             # 全プロジェクトのLintチェック
make format           # 全プロジェクトのフォーマット
```

### テスト

```bash
make test             # ユニットテスト実行
make test-integration # 統合テスト実行
make test-e2e         # E2Eテスト実行
```

### ビルド

```bash
make build            # 全プロジェクトをビルド
make clean            # ビルド成果物を削除
```

## Git Hooksの動作確認

### テストコミット

```bash
# 1. ファイルを変更
echo "test" > test.txt

# 2. ステージング
git add test.txt

# 3. コミット（Git hooksが自動実行される）
git commit -m "test: Git hooksのテスト"
```

成功すると：

- ✅ Linterが実行される
- ✅ フォーマットが適用される
- ✅ コミットメッセージがチェックされる
- ✅ コミットが完了する

失敗すると：

- ❌ エラーメッセージが表示される
- ❌ コミットが中断される
- 🔧 エラーを修正して再度コミット

### Git Hooksの無効化（非推奨）

緊急時のみ使用：

```bash
git commit --no-verify -m "message"
```

## トラブルシューティング

### Git hooksが実行されない

```bash
# 権限を確認
ls -la .husky/

# 権限を付与
chmod +x .husky/pre-commit
chmod +x .husky/commit-msg

# 再セットアップ
make setup
```

### Linterエラーで進めない

```bash
# 自動修正を試す
make format

# それでもエラーが残る場合は手動で修正
make lint
```

### コミットメッセージエラー

正しい形式を使用：

```bash
# ❌ 間違い
git commit -m "updated files"

# ✅ 正しい
git commit -m "chore: ファイルを更新"
```

## 次のステップ

1. [CONTRIBUTING.md](CONTRIBUTING.md) - 開発ガイドライン
2. [INTEGRATION.md](INTEGRATION.md) - 統合とデプロイ
3. [frontend/PERFORMANCE.md](frontend/PERFORMANCE.md) - パフォーマンス最適化
4. [e2e/README.md](e2e/README.md) - E2Eテスト
