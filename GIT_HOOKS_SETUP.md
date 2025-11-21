# Git Hooks セットアップ完了

## 🎉 実装内容

コミット時に自動的にlinterを実行するGit hooksを設定しました。

## 📦 追加されたファイル

### Git Hooks

- `.husky/pre-commit` - コミット前にlint-stagedを実行
- `.husky/commit-msg` - コミットメッセージの形式をチェック

### 設定ファイル

- `package.json` - ルートのnpm設定とスクリプト
- `.lintstagedrc.json` - lint-stagedの設定
- `.commitlintrc.json` - コミットメッセージのルール
- `.prettierrc.json` - Prettierの設定
- `.prettierignore` - Prettierの除外ファイル
- `Makefile` - 便利なmakeコマンド

### ドキュメント

- `CONTRIBUTING.md` - 開発ガイドライン
- `SETUP.md` - セットアップガイド
- `GIT_HOOKS_SETUP.md` - このファイル

## 🚀 セットアップ方法

```bash
# 1. 依存関係をインストール
npm install

# 2. Git hooksをセットアップ
npm run prepare
# または
make setup
```

## ✨ 自動実行される内容

### コミット前（pre-commit）

#### すべてのステージされたファイル

```bash
prettier --write    # フォーマット（JSON/YAML/Markdown/TS/JS）
```

#### Goファイル（変更がある場合のみ）

```bash
go fmt ./...    # フォーマット
go vet ./...    # 静的解析
```

#### フロントエンド（変更がある場合のみ）

```bash
npm run type-check    # TypeScript型チェック
```

### コミットメッセージ（commit-msg）

Conventional Commits形式をチェック：

```
<type>(<scope>): <subject>
```

**許可されるtype:**

- `feat` - 新機能
- `fix` - バグ修正
- `docs` - ドキュメント
- `style` - コードスタイル
- `refactor` - リファクタリング
- `perf` - パフォーマンス改善
- `test` - テスト
- `build` - ビルド
- `ci` - CI/CD
- `chore` - その他

## 📝 使用例

### 正しいコミット

```bash
# 新機能追加
git commit -m "feat(frontend): 資産推移チャートを追加"

# バグ修正
git commit -m "fix(backend): 計算ロジックのバグを修正"

# ドキュメント更新
git commit -m "docs: READMEにセットアップ手順を追加"

# リファクタリング
git commit -m "refactor(api): API clientを整理"
```

### 間違ったコミット（拒否される）

```bash
# ❌ typeがない
git commit -m "update code"

# ❌ 不正なtype
git commit -m "update: code changes"

# ❌ コロンがない
git commit -m "feat add feature"
```

## 🛠️ 便利なコマンド

### Makeコマンド

```bash
make help              # コマンド一覧
make install           # 依存関係インストール
make setup             # Git hooks設定
make dev               # 開発サーバー起動
make lint              # Lintチェック
make format            # コードフォーマット
make test              # テスト実行
make build             # ビルド
make clean             # クリーンアップ
```

### npmスクリプト

```bash
npm run lint           # 全プロジェクトのLint
npm run format         # 全プロジェクトのフォーマット
npm run test           # 全プロジェクトのテスト
npm run dev:frontend   # フロントエンド起動
npm run dev:backend    # バックエンド起動
```

## 🔧 トラブルシューティング

### Git hooksが実行されない

```bash
# 権限を確認
ls -la .husky/

# 権限を付与
chmod +x .husky/pre-commit
chmod +x .husky/commit-msg

# 再インストール
rm -rf node_modules
npm install
npm run prepare
```

### Linterエラーで進めない

```bash
# 自動修正を試す
make format

# 個別に確認
cd frontend && npm run lint -- --fix
cd backend && gofmt -w .
```

### コミットメッセージエラー

```bash
# 正しい形式を使用
git commit -m "feat: 新機能"
git commit -m "fix: バグ修正"
git commit -m "docs: ドキュメント更新"
```

### 緊急時（非推奨）

Git hooksをスキップ：

```bash
git commit --no-verify -m "message"
```

⚠️ **注意**: 通常は使用しないでください。CI/CDで失敗する可能性があります。

## 📊 ワークフロー

```
1. コードを変更
   ↓
2. git add .
   ↓
3. git commit -m "feat: 新機能"
   ↓
4. [pre-commit実行]
   - Linter実行
   - フォーマット適用
   - 型チェック
   ↓
5. [commit-msg実行]
   - メッセージ形式チェック
   ↓
6. コミット完了 ✅
   ↓
7. git push
   ↓
8. [GitHub Actions実行]
   - Lint
   - Test
   - Build
```

## 🎯 メリット

### 開発者

- ✅ コミット前に自動的にコードが整形される
- ✅ 型エラーを早期に発見
- ✅ 統一されたコードスタイル
- ✅ 統一されたコミットメッセージ

### チーム

- ✅ コードレビューが楽になる
- ✅ CI/CDの失敗が減る
- ✅ コードの品質が向上
- ✅ Git履歴が読みやすくなる

## 📚 関連ドキュメント

- [SETUP.md](SETUP.md) - セットアップガイド
- [CONTRIBUTING.md](CONTRIBUTING.md) - 開発ガイドライン
- [INTEGRATION.md](INTEGRATION.md) - 統合とデプロイ
- [GITHUB_ACTIONS_FIXES.md](GITHUB_ACTIONS_FIXES.md) - CI/CD修正内容

## 🔗 参考リンク

- [Husky](https://typicode.github.io/husky/)
- [lint-staged](https://github.com/okonet/lint-staged)
- [Commitlint](https://commitlint.js.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Prettier](https://prettier.io/)
