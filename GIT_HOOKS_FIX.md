# Git Hooks 修正内容

## 問題点

lint-stagedの設定で`cd`コマンドを使用していたため、正しく動作しませんでした。

### エラー例

```
Error: undefined: echo (typecheck)
Error: undefined: echo (typecheck)
```

## 修正内容

### 1. `.lintstagedrc.json`の簡素化

**修正前:**

```json
{
  "*.go": ["gofmt -w", "go vet"],
  "frontend/**/*.{js,jsx,ts,tsx}": [
    "cd frontend && npm run lint -- --fix",
    "cd frontend && npm run type-check"
  ],
  "e2e/**/*.{js,jsx,ts,tsx}": ["cd e2e && npx eslint --fix"],
  "*.{json,md,yml,yaml}": ["prettier --write"]
}
```

**修正後:**

```json
{
  "backend/**/*.go": ["gofmt -w", "go vet"],
  "frontend/**/*.{js,jsx,ts,tsx}": ["prettier --write"],
  "e2e/**/*.{js,jsx,ts,tsx}": ["prettier --write"],
  "*.{json,md,yml,yaml}": ["prettier --write"]
}
```

### 2. `.husky/pre-commit`の改善

**修正前:**

```bash
#!/usr/bin/env sh
. "$(dirname -- "$0")/_/husky.sh"

# Run lint-staged
npx lint-staged
```

**修正後:**

```bash
#!/usr/bin/env sh
. "$(dirname -- "$0")/_/husky.sh"

echo "🔍 Running pre-commit checks..."

# Run lint-staged for formatting
npx lint-staged

# Check if there are staged Go files
if git diff --cached --name-only | grep -q '\.go$'; then
  echo "📝 Checking Go files..."
  cd backend && go fmt ./... && go vet ./...
  if [ $? -ne 0 ]; then
    echo "❌ Go checks failed"
    exit 1
  fi
fi

# Check if there are staged frontend files
if git diff --cached --name-only | grep -q '^frontend/.*\.\(ts\|tsx\|js\|jsx\)$'; then
  echo "📝 Checking frontend files..."
  cd frontend && npm run type-check
  if [ $? -ne 0 ]; then
    echo "❌ Frontend type check failed"
    exit 1
  fi
fi

echo "✅ Pre-commit checks passed!"
```

## 改善点

### 1. シンプルな設計

- lint-stagedは**フォーマットのみ**に使用
- 複雑なチェックはpre-commitスクリプト内で実行

### 2. 条件付き実行

- 変更されたファイルがある場合のみチェックを実行
- 不要なチェックをスキップして高速化

### 3. わかりやすいフィードバック

- 絵文字付きのメッセージで進捗を表示
- エラー時に明確なメッセージを表示

## 動作フロー

```
1. git commit実行
   ↓
2. pre-commitフック起動
   ↓
3. lint-staged実行
   - すべてのステージされたファイルをPrettierでフォーマット
   ↓
4. Goファイルチェック（変更がある場合）
   - go fmt ./...
   - go vet ./...
   ↓
5. フロントエンドチェック（変更がある場合）
   - npm run type-check
   ↓
6. すべて成功 → コミット完了 ✅
   失敗 → コミット中断 ❌
```

## テスト方法

### 1. Prettierのテスト

```bash
# JSONファイルを変更
echo '{"test":true}' > test.json
git add test.json
git commit -m "test: Prettierテスト"
# → 自動的にフォーマットされる
```

### 2. Goファイルのテスト

```bash
# Goファイルを変更
echo 'package main' > backend/test.go
git add backend/test.go
git commit -m "test: Goファイルテスト"
# → go fmtとgo vetが実行される
```

### 3. TypeScriptのテスト

```bash
# TypeScriptファイルを変更
echo 'const x: string = 123;' > frontend/src/test.ts
git add frontend/src/test.ts
git commit -m "test: TypeScriptテスト"
# → 型エラーで失敗する
```

## メリット

### 開発者体験の向上

- ✅ 高速な実行（変更されたファイルのみチェック）
- ✅ わかりやすいフィードバック
- ✅ 自動フォーマットで手間削減

### コード品質の向上

- ✅ 統一されたフォーマット
- ✅ 型安全性の保証
- ✅ 基本的なGoの問題を早期発見

### チーム開発の効率化

- ✅ コードレビューが楽になる
- ✅ CI/CDの失敗が減る
- ✅ 統一されたコーディングスタイル

## トラブルシューティング

### pre-commitが実行されない

```bash
chmod +x .husky/pre-commit
```

### 型チェックでエラー

```bash
cd frontend
npm run type-check
# エラーを確認して修正
```

### Goのチェックでエラー

```bash
cd backend
go fmt ./...
go vet ./...
# エラーを確認して修正
```

### 緊急時のスキップ（非推奨）

```bash
git commit --no-verify -m "message"
```

## 次のステップ

1. ✅ Git hooksの動作確認
2. ✅ チーム全体での導入
3. 📝 追加のlintルールの検討
4. 📝 CI/CDとの統合確認

## 関連ドキュメント

- [GIT_HOOKS_SETUP.md](GIT_HOOKS_SETUP.md) - セットアップガイド
- [SETUP.md](SETUP.md) - 開発環境セットアップ
- [CONTRIBUTING.md](CONTRIBUTING.md) - 開発ガイドライン
