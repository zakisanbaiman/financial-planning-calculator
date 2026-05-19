# Architecture Decision Records (ADR)

このディレクトリには、プロジェクトの重要な技術的決定を記録したArchitecture Decision Records（ADR）が格納されています。

## ADRとは

ADRは、アーキテクチャ上の重要な決定を記録するための軽量なドキュメントです。各ADRには以下の情報が含まれます：

- **ステータス**: 提案中、採択済み、却下、非推奨など
- **背景**: なぜこの決定が必要になったのか
- **決定内容**: 何を決定したのか
- **重視するアーキテクチャ特性**: この決定が何を守るためのものか
- **理由**: なぜその決定をしたのか
- **代替案**: 他にどのような選択肢があったか
- **結果**: この決定により期待される影響

アーキテクチャ特性の例:

- 変更容易性（Modifiability）
- 保守性（Maintainability）
- 可用性（Availability）
- 性能・応答性（Performance / Responsiveness）
- 拡張性（Scalability）
- セキュリティ（Security）
- 運用性（Operability）
- テスト容易性（Testability）
- コスト効率（Cost Efficiency）

## ADR一覧

- [ADR-001: asdfからmiseへの移行](./001-migrate-to-mise.md) (2026-01-22)
- [ADR-002: SPAアーキテクチャの採用について](./002-spa-adoption.md) (2026-01-27)
- [ADR-003: Render.comからRailway + Neonへの移行](./003-migrate-to-railway-neon.md) (2026-03-26)
- [ADR-004: メール送信にResend HTTP APIを採用](./004-email-resend-http-api.md) (2026-03-26)
- [ADR-005: ボット応答の返却方式として SSE を採用](./005-bot-integration-sse.md) (2026-03-31)
- [ADR-006: Redisキャッシュ戦略の採用（Cache-Aside + Repositoryデコレータ）](./006-redis-cache-strategy.md) (2026-04-02)
- [ADR-007: DDDドメイン境界の設計方針](./007-ddd-domain-boundary.md) (2026-04-27)

## 新しいADRの作成

重要な技術的決定を行う際は、新しいADRを作成してください：

1. 連番のファイル名で作成（例: `002-description.md`）
2. 上記のテンプレートに従って記述
3. 「どのアーキテクチャ特性を優先した決定か」を明記する
4. 可能ならトレードオフも書く
5. README.mdに追加

## 参考資料

- [Architecture Decision Records](https://adr.github.io/)
- [ADRのベストプラクティス](https://github.com/joelparkerhenderson/architecture-decision-record)
