# Dev Container 実装サマリー

## 概要

このドキュメントは、financial-planning-calculator プロジェクトに Dev Container を導入した実装の詳細をまとめたものです。

## 実装日

2026-02-02

## 背景と目的

### 背景

開発環境のセットアップは、新しいメンバーのオンボーディングや複数の開発マシン間での一貫性を保つために課題となることがあります。このプロジェクトでは以下の課題がありました：

- Go 1.24、Node.js 18、PostgreSQLなど複数のツールのインストールが必要
- 各開発者のマシン環境による差異
- セットアップ手順の複雑さ

### 目的

1. **簡単なセットアップ**: VS Codeで開くだけで開発環境が整う
2. **一貫性**: すべての開発者が同じ環境で作業できる
3. **効率性**: セットアップ時間を大幅に短縮
4. **再現性**: 環境の問題を減らし、トラブルシューティングを容易にする

## 実装内容

### ファイル構成

```
.devcontainer/
├── devcontainer.json       # Dev Containerのメイン設定
├── docker-compose.yml      # Dev Container用のDocker Compose設定
├── post-create.sh          # 自動セットアップスクリプト
└── README.md               # ユーザー向けガイド
```

### 1. devcontainer.json

Dev Containerのメイン設定ファイル。以下の機能を提供：

#### 基本設定
- **名前**: "Financial Planning Calculator"
- **ベース**: 既存のdocker-compose.ymlを拡張
- **ワークスペース**: `/workspace`にマウント
- **サービス**: backendコンテナを使用

#### 機能（Features）
- **common-utils**: Zsh、Oh My Zshを含む基本ユーティリティ
- **Node.js**: v18（フロントエンド開発用）
- **Go**: v1.24（バックエンド開発用）
- **Git**: 最新版

#### VS Code統合

**推奨拡張機能（自動インストール）:**

| 拡張機能 | 用途 |
|---------|------|
| golang.go | Go言語サポート |
| dbaeumer.vscode-eslint | JavaScript/TypeScript Lint |
| esbenp.prettier-vscode | コードフォーマッター |
| bradlc.vscode-tailwindcss | Tailwind CSS補完 |
| ms-azuretools.vscode-docker | Docker管理 |
| eamodio.gitlens | Git機能拡張 |
| github.copilot | AIコード補完 |
| swaggo.swaggo | Swagger/OpenAPIサポート |
| その他 | 生産性向上ツール |

**エディタ設定:**
- 保存時自動フォーマット
- ESLintの自動修正
- インポートの自動整理
- Go、TypeScript、JSON、Markdownの適切なフォーマッター設定

#### ポート転送

| ポート | サービス | 用途 |
|-------|---------|------|
| 3000 | Frontend | Next.js開発サーバー |
| 8080 | Backend | Go API サーバー |
| 5432 | PostgreSQL | データベース |
| 6060 | pprof | パフォーマンスプロファイリング |

### 2. docker-compose.yml

既存の`docker-compose.yml`を拡張するDev Container専用設定：

```yaml
services:
  backend:
    volumes:
      - ..:/workspace:cached          # プロジェクト全体をマウント
      - go_mod_cache:/go/pkg/mod      # Go modulesキャッシュ
      - devcontainer-bashhistory:/commandhistory  # コマンド履歴を永続化
    environment:
      - DEVCONTAINER=true             # Dev Container環境を示す
    command: sleep infinity           # VS Codeがアタッチするため待機
    user: root                        # VS Code管理のため
```

**特徴:**
- 既存のDocker環境と共存可能
- ワークスペース全体を`/workspace`にマウント
- コマンド履歴の永続化（再起動しても履歴が残る）
- Go modulesキャッシュの永続化（依存関係の再ダウンロードを回避）

### 3. post-create.sh

Dev Container作成後に自動実行されるセットアップスクリプト：

#### 実行内容

1. **Git設定**
   - 改行コードの設定（LF）
   - Pull戦略の設定

2. **シェル設定**
   - Bash/Zshヒストリーファイルの永続化設定

3. **Go Tools インストール**
   ```bash
   - air v1.52.3          # ホットリロード
   - golangci-lint v1.64.0 # Linter
   - swag (latest)         # Swagger生成
   ```

4. **依存関係のインストール**
   - バックエンド: `go mod download`
   - フロントエンド: `npm ci`
   - E2Eテスト: `npm ci`
   - ルート: `npm ci`

5. **Git Hooks セットアップ**
   - Huskyを使用したコミットフックの設定

6. **データベースセットアップ**
   - PostgreSQLの起動確認
   - マイグレーションの実行
   - シードデータの投入

#### エラーハンドリング
- データベースが起動していない場合は警告を表示し、後で手動実行可能
- 既にマイグレーション済みの場合はスキップ

### 4. README.md

`.devcontainer/README.md`には以下の内容を記載：
- セットアップ手順
- 使い方とコマンド一覧
- インストールされる拡張機能の説明
- カスタマイズ方法
- トラブルシューティング
- 既存のDocker環境との違い

## 技術的な考慮事項

### 1. 既存環境との共存

Dev Containerは既存のDocker開発環境と完全に共存可能です：

- 同じ`docker-compose.yml`を基盤として使用
- Dev Container専用の設定は別ファイル（`.devcontainer/docker-compose.yml`）で管理
- 既存の`make`コマンドがそのまま使用可能

### 2. パフォーマンス最適化

- **cachedマウント**: ファイルシステムの読み書き性能を最適化
- **Go modulesキャッシュ**: 依存関係の再ダウンロードを回避
- **永続化ボリューム**: コマンド履歴やキャッシュを保持

### 3. セキュリティ

- ルートユーザーで実行（VS Code Dev Containersの標準）
- 本番環境とは完全に分離
- pprofは開発環境のみで有効（既存設定を維持）

### 4. 拡張性

Dev Container設定は以下の方法で簡単にカスタマイズ可能：

```json
// .devcontainer/devcontainer.json
{
  "customizations": {
    "vscode": {
      "extensions": ["your.extension.id"],
      "settings": {
        "your.setting": "value"
      }
    }
  },
  "forwardPorts": [9000],
  "postCreateCommand": "custom-script.sh"
}
```

## 使用方法

### 基本的な使い方

1. **起動**
   ```
   VS Code でプロジェクトを開く
   → Ctrl+Shift+P (Cmd+Shift+P)
   → "Dev Containers: Reopen in Container"
   ```

2. **開発サーバー起動**
   ```bash
   make up      # バックエンド + DB
   make up-full # すべてのサービス
   ```

3. **開発**
   - コードを編集（自動フォーマット有効）
   - ホットリロードで即座に反映
   - デバッグ、テスト実行も可能

### よく使うコマンド

```bash
# 開発
make up          # サービス起動
make down        # サービス停止
make logs        # ログ表示

# テストとLint
make test        # テスト実行
make lint        # Lint実行

# データベース
make migrate     # マイグレーション
make seed        # シードデータ投入
make shell-db    # PostgreSQLシェル

# ヘルプ
make help        # すべてのコマンド
make docker-help # Docker関連コマンド
```

## メリット

### 1. 開発者体験の向上

- **セットアップ時間**: 手動セットアップ 30-60分 → Dev Container 5-10分
- **エラー削減**: 環境差異によるエラーがほぼゼロ
- **オンボーディング**: 新メンバーがすぐに開発を開始できる

### 2. 一貫性

- すべての開発者が同じツールバージョンを使用
- VS Code拡張機能も統一
- 設定の共有が容易

### 3. 生産性

- 拡張機能の自動インストール
- エディタ設定の自動適用
- ポート転送の自動設定
- ホットリロード機能の維持

### 4. 保守性

- 設定ファイルでバージョン管理
- Docker化により環境の再現が容易
- トラブルシューティングが簡単

## 制限事項と今後の改善点

### 現在の制限事項

1. **初回起動時間**: 5-10分かかる（依存関係のダウンロード）
2. **リソース使用**: Docker Desktopが必要（メモリ、CPU）
3. **VS Code専用**: 他のIDEでは使用不可

### 今後の改善案

1. **プリビルドイメージ**: GitHub Actionsでビルド済みイメージを提供
2. **マルチコンテナ**: フロントエンドとバックエンドを分離
3. **Codespaces対応**: GitHub Codespacesでの実行をサポート

## テスト結果

### 検証項目

- [x] devcontainer.jsonのJSON構文チェック
- [x] docker-compose.ymlのYAML構文チェック
- [x] post-create.shのBash構文チェック
- [x] ファイルのパーミッション設定
- [x] 既存Docker環境との共存確認

### 期待される動作

1. Dev Containerでプロジェクトを開く
2. 自動的に以下がインストールされる：
   - Go 1.24、Node.js 18
   - Go tools（air, golangci-lint, swag）
   - すべての依存関係
   - VS Code拡張機能
3. データベースのマイグレーションとシードが実行される
4. `make up`でサービスが起動する
5. ホットリロードが動作する

## 参考リンク

- [Dev Containers公式ドキュメント](https://containers.dev/)
- [VS Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
- [Dev Container Features](https://containers.dev/features)
- [Dev Container仕様](https://containers.dev/implementors/spec/)

## まとめ

Dev Containerの導入により、開発環境のセットアップが大幅に簡素化され、すべての開発者が一貫した環境で作業できるようになりました。既存のDocker環境とも共存可能で、開発者は自分の好みに応じて選択できます。

この実装は、プロジェクトの開発者体験を向上させ、新しいメンバーのオンボーディングを容易にすることを目的としています。
