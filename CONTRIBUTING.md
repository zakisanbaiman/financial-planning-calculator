# Contributing Guide

## 開発環境のセットアップ

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd financial-planning-calculator
```

### 2. 依存関係のインストール

```bash
make install
```

または個別に：

```bash
npm install                    # ルートの依存関係
cd frontend && npm install     # フロントエンド
cd e2e && npm install          # E2Eテスト
cd backend && go mod download  # バックエンド
```

### 3. Git Hooksのセットアップ

```bash
make setup
```

これにより、以下のGit hooksが設定されます：

- **pre-commit**: コミット前にlinterを実行
- **commit-msg**: コミットメッセージの形式をチェック

## コーディング規約

### コミットメッセージ

Conventional Commits形式を使用します：

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Type（必須）

- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響しない変更（空白、フォーマット等）
- `refactor`: バグ修正や機能追加ではないコード変更
- `perf`: パフォーマンス改善
- `test`: テストの追加や修正
- `build`: ビルドシステムや外部依存関係の変更
- `ci`: CI設定ファイルやスクリプトの変更
- `chore`: その他の変更

#### 例

```bash
feat(frontend): 資産推移チャートにズーム機能を追加

ユーザーが特定の期間にフォーカスできるように、
チャートにズーム機能を実装しました。

Closes #123
```

### コードスタイル

#### Go（バックエンド）

- `gofmt`でフォーマット
- `go vet`でチェック
- 標準的なGo命名規則に従う

```bash
make format     # 自動フォーマット
make lint       # Lintチェック
```

#### TypeScript/JavaScript（フロントエンド）

- ESLintルールに従う
- Prettierでフォーマット
- Next.jsのベストプラクティスに従う

```bash
cd frontend
npm run lint -- --fix    # Lint + 自動修正
npm run type-check       # 型チェック
```

## 開発ワークフロー

### 1. ブランチの作成

```bash
git checkout -b feature/your-feature-name
# または
git checkout -b fix/your-bug-fix
```

### 2. 開発

```bash
make dev    # 開発サーバー起動
```

### 3. テスト

```bash
make test              # 全テスト実行
make test-integration  # 統合テスト
make test-e2e         # E2Eテスト
```

### 4. コミット

Git hooksが自動的に実行されます：

```bash
git add .
git commit -m "feat: 新機能の追加"
```

コミット時に以下が自動実行されます：

- ✅ Linterチェック
- ✅ 型チェック
- ✅ コミットメッセージ形式チェック

### 5. プッシュ

```bash
git push origin feature/your-feature-name
```

### 6. Pull Request作成

GitHub上でPull Requestを作成します。

## テスト

### ユニットテスト

```bash
# バックエンド
cd backend && go test ./...

# フロントエンド
cd frontend && npm test
```

### E2Eテスト

```bash
cd e2e
npm test
```

### 統合テスト

```bash
./scripts/test-integration.sh
```

## トラブルシューティング

### Git hooksが動作しない

```bash
chmod +x .husky/pre-commit
chmod +x .husky/commit-msg
```

### Linterエラー

```bash
make format    # 自動修正を試す
make lint      # エラー確認
```

### 依存関係の問題

```bash
make clean
make install
```

## 便利なコマンド

```bash
make help              # 利用可能なコマンド一覧
make install           # 依存関係インストール
make setup             # Git hooks設定
make dev               # 開発サーバー起動
make lint              # Lintチェック
make format            # コードフォーマット
make test              # テスト実行
make build             # ビルド
make clean             # クリーンアップ
```

## 質問やサポート

- GitHub Issuesで質問を投稿
- ドキュメントを確認：
  - [INTEGRATION.md](INTEGRATION.md)
  - [frontend/PERFORMANCE.md](frontend/PERFORMANCE.md)
  - [e2e/README.md](e2e/README.md)
