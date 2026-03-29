# ADR 004: メール送信にResend HTTP APIを採用

## ステータス

採択済み (2026-03-26)

## 背景

パスワードリセット機能（Issue #199）の実装にあたり、メール送信の仕組みが必要になった。当初はSMTP（`net/smtp`）で実装していたが、本番環境（Railway）でメールが届かない問題が発生した。

## 問題

RailwayはSMTPポート（587）をブロックしているため、`net/smtp` による送信が30秒タイムアウトして失敗していた。

## 決定

**Resend（https://resend.com）のHTTP APIを使ってメール送信を行う。**

- SMTPではなくHTTPS経由でAPIを呼び出すため、ポートブロックの影響を受けない
- 既存の `SMTP_PASSWORD` 環境変数をResendのAPIキーとして流用し、設定変更を最小限に抑える

## 理由

### Resendを選んだ理由

1. **無料枠が十分**
   - 月3,000通・1日100通まで無料、クレジットカード不要
   - 個人プロジェクトの用途（パスワードリセットのみ）には十分

2. **設定が最もシンプル**
   - APIキー1つで動作
   - SMTP設定（host/port/user）が不要

3. **HTTP APIでRailwayのポート制限を回避**
   - Railway環境でSMTPポート587がブロックされていることを確認済み
   - HTTPSは常に通るため安定して動作する

### 代替案と却下理由

1. **SendGrid**
   - 利点: 実績豊富、無料枠あり（100通/日）
   - 欠点: 設定が複雑、Resendと比べて開発者体験が劣る

2. **AWS SES**
   - 利点: 大量送信に強い、安価
   - 欠点: AWS設定が必要、個人プロジェクトには過剰

3. **SMTP継続（別ポート利用）**
   - 利点: 実装変更不要
   - 欠点: Railwayでは465（SSL）も制限される可能性があり根本解決にならない

## 実装

- `backend/infrastructure/email/email_service.go` に `ResendEmailService` を実装
- `SMTP_PASSWORD` が空の場合は開発用ログ出力にフォールバック（ローカル開発での設定不要）
- `SMTP_HOST` / `SMTP_USER` のデフォルト値はResend用に設定済み（`config/server.go`）

## Railway環境変数

| 変数 | 値 | 備考 |
|---|---|---|
| `SMTP_PASSWORD` | ResendのAPIキー（`re_xxxxx`） | APIキーとして使用 |
| `SMTP_FROM` | `onboarding@resend.dev` | 独自ドメインなし。自分のメールにしか送れない制限あり |
| `FRONTEND_URL` | `https://financial-planning-frontend-production.up.railway.app` | リセットURLの生成に使用 |

## 今後の課題

- `SMTP_FROM` は現在 `onboarding@resend.dev`（テスト用）。本番ユーザーへの送信には独自ドメインの取得・DNS認証が必要
- 独自ドメイン取得後は `SMTP_FROM` を `noreply@<ドメイン>` に変更する
