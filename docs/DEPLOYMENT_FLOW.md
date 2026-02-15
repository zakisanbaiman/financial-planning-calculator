# デプロイフロー

このドキュメントでは、財務計画計算機の開発からデプロイまでの全体フローを説明します。

## 概要

```
[ローカル開発] → [GitHub] → [Render.com]
     │              │              │
     │              ├─ PR作成 ───→ プレビュー環境
     │              │
     │              └─ mainマージ → 本番環境
     │
     └─ docker-compose で動作確認
```

## 現在のデプロイ構成

| 環境 | プラットフォーム | トリガー |
|------|----------------|---------|
| ローカル | Docker Compose | 手動 (`make dev`) |
| プレビュー | Render.com | PR 作成・更新時 |
| 本番 | Render.com | main ブランチへのマージ |

---

## 1. ローカル開発

### 起動方法

```bash
# 開発環境起動（Backend + DB）
make dev

# フロントエンドも含める場合
make dev-frontend
```

### 構成

```yaml
# docker-compose.yml
services:
  postgres:    # PostgreSQL (port 5432)
  backend:     # Go/Echo API (port 8080)
  frontend:    # Next.js (port 3000) ※profiles: [frontend]
```

---

## 2. GitHub → Render 連携

### 連携の仕組み

```
[GitHub リポジトリ]
       │
       ├─ OAuth 認可 ──────────────────┐
       │                               ↓
       └─ Webhook ────────────→ [Render.com]
          (push イベント通知)          │
                                       ↓
                               変更ファイル一覧を取得
                               (GitHub API 経由)
                                       ↓
                               buildFilter と照合
                                       ↓
                               該当サービスをデプロイ
```

### 必要な権限（OAuth）

Render を GitHub に接続すると、以下の権限を付与:

| 権限 | 用途 |
|------|------|
| リポジトリ読み取り | コードを clone してビルド |
| Webhook 登録 | push イベントを受け取る |
| コミット/差分参照 | どのファイルが変更されたか確認 |

**確認方法**: GitHub → Settings → Applications → Authorized OAuth Apps

---

## 3. 本番デプロイ（main マージ時）

### フロー

```
[PR を main にマージ]
       ↓
[Render が Webhook で検知]
       ↓
[変更ファイルをチェック]
       ↓
[buildFilter に一致するサービスのみ再ビルド]
       ↓
[デプロイ完了]
```

### buildFilter による選択的デプロイ

`render.yaml` の `buildFilter` 設定により、変更されたファイルに応じて必要なサービスのみデプロイされます。

```yaml
# Backend サービス
buildFilter:
  paths:
    - backend/**

# Frontend サービス
buildFilter:
  paths:
    - frontend/**
```

| 変更されたファイル | デプロイ対象 |
|------------------|-------------|
| `backend/main.go` | Backend のみ |
| `frontend/src/app/page.tsx` | Frontend のみ |
| 両方に変更 | Backend + Frontend |
| `README.md` のみ | デプロイなし |

### render.yaml の構成

```yaml
databases:
  - name: financial-planning-db
    plan: free
    region: oregon

services:
  # Backend API
  - type: web
    name: financial-planning-backend
    env: docker
    dockerfilePath: ./backend/Dockerfile
    healthCheckPath: /health
    envVars:
      - key: DB_HOST
        fromDatabase:
          name: financial-planning-db
          property: host
      # ... その他のDB接続情報
    buildFilter:
      paths:
        - backend/**

  # Frontend
  - type: web
    name: financial-planning-frontend
    env: docker
    dockerfilePath: ./frontend/Dockerfile.prod
    envVars:
      - key: NEXT_PUBLIC_API_URL
        sync: false  # Render ダッシュボードで手動設定
    buildFilter:
      paths:
        - frontend/**
```

---

## 4. PR プレビュー環境

PR を作成すると、独立したプレビュー環境が自動デプロイされます。

### URL 形式

```
Frontend: https://financial-planning-frontend-pr-{PR番号}.onrender.com
Backend:  https://financial-planning-backend-pr-{PR番号}.onrender.com
```

### 特徴

| 項目 | 内容 |
|------|------|
| 独立した DB | PR ごとに PostgreSQL インスタンス |
| 自動デプロイ | PR 作成・更新時 |
| 有効期限 | PR クローズ後 7 日間 |
| 無料プラン制限 | 15分非アクティブでスリープ |

詳細: [プレビュー環境クイックリファレンス](./PREVIEW_ENVIRONMENT_QUICK_REF.md)

---

## 5. 環境変数管理

### Backend

| 変数 | 設定方法 | 用途 |
|------|---------|------|
| `DB_*` | `fromDatabase` で自動注入 | DB 接続情報 |
| `JWT_SECRET` | Render ダッシュボードで手動設定 | JWT 署名 |
| `DB_SSLMODE` | `render.yaml` で固定値 | SSL 接続 |

### Frontend

| 変数 | 設定方法 | 用途 |
|------|---------|------|
| `NEXT_PUBLIC_API_URL` | Render ダッシュボードで手動設定 | API エンドポイント |
| `NODE_ENV` | `render.yaml` で固定値 | 実行環境 |

**`sync: false`** = Render ダッシュボードで手動設定が必要

---

## 6. デプロイ監視

### GitHub Actions による自動監視

`.github/workflows/monitor-render-deployments.yml` により:

1. デプロイ状況を定期監視
2. 失敗時に PR へコメント追加
3. 重大な失敗時に Issue 作成

### MCP 連携

AI アシスタント（Claude、Copilot、Cursor）が Render のデプロイ状態を直接監視できます。

詳細: [MCP セットアップガイド](./MCP_SETUP.md)

---

## 7. 将来の AWS 移行案（参考）

現在は Render.com 無料プランで運用していますが、商用化時には AWS への移行を検討できます。

### 比較表

| 観点 | Render (現状) | AWS (将来) |
|------|--------------|-----------|
| 月額費用 | $0 | $30〜250+ |
| 可用性 | Single Region | Multi-AZ 対応 |
| スケーリング | 手動 | Auto Scaling |
| 監視 | 基本的なログ | CloudWatch + アラート |

### AWS 構成案（商用化時）

```
[CloudFront + WAF]
       ↓
     [ALB]
    ↙    ↘
[ECS Fargate]  [ECS Fargate]
  Backend        Frontend
    ↓   ↘
    ↓    [ElastiCache Redis]
    ↓
[RDS PostgreSQL Multi-AZ]
```

### 段階的移行戦略

```
[現在] Render 無料 ($0)
  ↓ ユーザー増加
[次] Render 有料 ($7-25)
  ↓ 売上が立った段階
[将来] AWS 本番構成 ($150-250)
```

**推奨**: PMF（Product-Market Fit）検証まで Render 無料プランを継続

---

## 関連ドキュメント

- [render.yaml](../render.yaml) - Render デプロイ設定
- [docker-compose.yml](../docker-compose.yml) - ローカル開発環境
- [プレビュー環境クイックリファレンス](./PREVIEW_ENVIRONMENT_QUICK_REF.md)
- [MCP セットアップガイド](./MCP_SETUP.md)
- [Docker セットアップ](./DOCKER_SETUP.md)
