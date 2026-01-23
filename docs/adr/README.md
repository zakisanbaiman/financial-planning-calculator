# Architecture Decision Records (ADR)

このディレクトリには、プロジェクトの重要な技術的決定を記録したArchitecture Decision Records（ADR）が格納されています。

## ADRとは

ADRは、アーキテクチャ上の重要な決定を記録するための軽量なドキュメントです。各ADRには以下の情報が含まれます：

- **ステータス**: 提案中、採択済み、却下、非推奨など
- **背景**: なぜこの決定が必要になったのか
- **決定内容**: 何を決定したのか
- **理由**: なぜその決定をしたのか
- **代替案**: 他にどのような選択肢があったか
- **結果**: この決定により期待される影響

## ADR一覧

- [ADR-001: asdfからmiseへの移行](./001-migrate-to-mise.md) (2026-01-22)

## 新しいADRの作成

重要な技術的決定を行う際は、新しいADRを作成してください：

1. 連番のファイル名で作成（例: `002-description.md`）
2. 上記のテンプレートに従って記述
3. README.mdに追加

## 参考資料

- [Architecture Decision Records](https://adr.github.io/)
- [ADRのベストプラクティス](https://github.com/joelparkerhenderson/architecture-decision-record)
