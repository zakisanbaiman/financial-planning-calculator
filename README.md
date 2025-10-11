# 財務計画計算機

将来の資産形成と老後の財務計画を可視化するWebアプリケーション

## 概要

このアプリケーションは、ユーザーが現在の収入、支出、貯蓄状況を入力することで、将来の資産推移、老後資金、緊急時の備えを計算し、安心できる財務計画を立てられるようにします。

## 技術スタック

### フロントエンド
- Next.js 14
- TypeScript
- Tailwind CSS
- Chart.js
- React Hook Form + Zod

### バックエンド
- Go
- Echo Framework
- PostgreSQL
- OpenAPI/Swagger

## プロジェクト構造

```
financial-planning-calculator/
├── frontend/          # Next.jsフロントエンド
│   ├── src/
│   │   ├── app/       # App Router
│   │   ├── components/
│   │   ├── lib/
│   │   └── types/
│   ├── package.json
│   ├── next.config.js
│   └── tailwind.config.js
├── backend/           # Goバックエンド
│   ├── config/        # 設定
│   ├── docs/          # OpenAPI仕様
│   ├── go.mod
│   └── main.go
└── README.md
```

## セットアップ

### 前提条件
- Node.js 18+
- Go 1.21+
- PostgreSQL 13+

### フロントエンド

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

フロントエンドは http://localhost:3000 で起動します。

### バックエンド

```bash
cd backend
go mod tidy
cp .env.example .env
go run main.go
```

バックエンドは http://localhost:8080 で起動します。

### API仕様

Swagger UIは http://localhost:8080/swagger/index.html で確認できます。

## 開発

### フロントエンド開発
- `npm run dev` - 開発サーバー起動
- `npm run build` - プロダクションビルド
- `npm run lint` - ESLint実行
- `npm run type-check` - TypeScript型チェック

### バックエンド開発
- `go run main.go` - サーバー起動
- `go test ./...` - テスト実行
- `go mod tidy` - 依存関係整理

## 機能

### 実装予定機能
- [ ] 財務データ入力・管理
- [ ] 資産推移シミュレーション
- [ ] 老後資金計算
- [ ] 緊急資金計算
- [ ] 目標設定・進捗管理
- [ ] データ可視化（グラフ・チャート）
- [ ] PDFレポート生成

## ライセンス

MIT License