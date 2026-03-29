# ADR 003: Render.comからRailway + Neonへの移行

## ステータス

採択済み (2026-03-26)

## 背景

本プロジェクトはRender.comでバックエンド・フロントエンドをホスティングし、PostgreSQLもRender上で運用していた。しかし運用コストや機能面を検討した結果、プラットフォームの移行を行った。

## 決定

- **バックエンド・フロントエンド**: Render.com → **Railway** に移行
- **データベース**: Render PostgreSQL → **Neon（サーバーレスPostgreSQL）** に移行

## 理由

### Railwayの利点

1. **開発者体験の向上**
   - CLIが直感的で操作しやすい（`railway logs`, `railway variable set`等）
   - GitHubへのpushで自動デプロイ
   - サービス間のプライベートネットワークが使いやすい

2. **コスト効率**
   - 使用量ベースの課金でアイドル時のコストが低い
   - 無料枠でも実用的に使える

### Neonの利点

1. **サーバーレスPostgreSQL**
   - アイドル時は自動スリープしコストを削減
   - 接続プーリングが組み込み済み（`-pooler` エンドポイント）

2. **ブランチ機能**
   - DBのブランチを切って開発・本番を分離可能
   - PR環境ごとのDB分離も将来的に可能

## 代替案

1. **Render.comを継続使用**
   - 利点: 変更不要
   - 欠点: コスト・機能面でRailwayに劣る

2. **Fly.io**
   - 利点: エッジデプロイが可能
   - 欠点: 設定が複雑、学習コストが高い

3. **AWS/GCP**
   - 利点: スケーラビリティが高い
   - 欠点: 個人プロジェクトには過剰、運用コストが高い

## 結果

- デプロイ・ログ確認・環境変数管理がRailway CLIで完結
- NeonのサーバーレスによりDB維持コストを削減
- **注意**: `DEPLOYMENT_FLOW.md` はRender時代の記述のため参考にしないこと

## Railwayサービス構成

| サービス名 | 役割 |
|---|---|
| `financial-planning-backend` | Go/Echoバックエンド |
| `Redis-62U3` | Redisキャッシュ |
| フロントエンドサービス | Next.jsフロントエンド |

フロントエンドURL: `https://financial-planning-frontend-production.up.railway.app`
