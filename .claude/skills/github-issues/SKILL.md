---
name: github-issues
description: GitHub Issue作成・管理のベストプラクティス。明確なissue定義、テンプレート活用、ラベル・アサイン管理、効果的なコミュニケーション時に使用。
---

# GitHub Issue 作成ガイド

効果的で保守しやすいGitHub Issueを作成するためのガイドラインです。

## Issue作成前のチェック

### 1. Issue検索
新規作成前に、**同じ問題が既に報告されていないか確認します**

```bash
# GitHub Searchで同じキーワードを検索
# - Closed issues も含めて確認
# - タイトルと概要で重複がないか確認
```

### 2. Issue タイプの判別

| タイプ | 説明 | ラベル |
|--------|------|--------|
| **Feature Request** | 新機能・改善 | `type/feat` |
| **Bug Report** | 不具合・エラー | `type/bug` |
| **Documentation** | ドキュメント不足・改善 | `type/docs` |
| **Question** | 疑問・相談 | `type/question` |
| **Infrastructure** | インフラ・設定 | `type/infra` |

## Issue テンプレート

### 1. Bug Report テンプレート

```markdown
## 現象
[何が起きているか、簡潔に説明]

## 期待される動作
[本来どうなるべきか]

## 再現手順
1. ...
2. ...
3. ...

## 環境
- OS: macOS / Linux / Windows
- Browser: Chrome / Safari / Firefox
- Version: [バージョン番号]

## ログ・スクリーンショット
[該当するログやエラーメッセージ、スクリーンショット]

## 備考
[その他の情報、試したことなど]
```

### 2. Feature Request テンプレート

```markdown
## 提案内容
[新機能・改善内容を簡潔に説明]

## 解決する課題
[この機能が解決する具体的な問題]

## 実装案
[推奨される実装方法があれば記述]

## 優先度
- [ ] 高（すぐに必要）
- [ ] 中（今後のスプリントで対応）
- [ ] 低（いつか対応）

## 受け入れ基準
- [ ] ...基準1...
- [ ] ...基準2...
- [ ] ...基準3...
```

### 3. Task テンプレート

```markdown
## タスク概要
[実施すべきタスクの説明]

## チェックリスト
- [ ] 実装
- [ ] テスト
- [ ] ドキュメント更新
- [ ] PR作成
- [ ] コードレビュー
- [ ] デプロイ

## 関連Issue
Closes #[number] (あれば)

## 参考資料
- [ドキュメント](link)
- 関連Issue: #[number]
```

## タイトル命名規則

```
[<type>] <主語><述語> - <補足>
```

### 例

**Good:**
- `[Bug] 老後資金計算で負数が表示される`
- `[Feat] チャート拡大・縮小機能の追加`
- `[Docs] セットアップガイドの日本語化`

**Bad:**
- `バグ`
- `Fix issue`
- `チャートをなおして`

## ラベル体系

### ステータス系
- `status/backlog` - バックログ
- `status/in-progress` - 実装中
- `status/review` - レビュー待ち
- `status/blocked` - ブロック中
- `status/done` - 完了

### 優先度
- `priority/critical` - 緊急（本番障害など）
- `priority/high` - 高
- `priority/medium` - 中
- `priority/low` - 低

### タイプ
- `type/bug` - バグ・不具合
- `type/feat` - 新機能
- `type/refactor` - リファクタリング
- `type/docs` - ドキュメント
- `type/infra` - インフラ・設定
- `type/test` - テスト

### 領域
- `area/frontend` - フロントエンド
- `area/backend` - バックエンド
- `area/database` - データベース
- `area/ci-cd` - CI/CD・GitHub Actions
- `area/docs` - ドキュメント

### その他
- `good-first-issue` - 初心者向け
- `help-wanted` - 協力者求む
- `duplicate` - 重複
- `wontfix` - 修正なし

## 効果的なIssue作成のベストプラクティス

### ✅ DO

1. **明確で簡潔**
   - タイトルは一目で問題が分かるように
   - 書きすぎず、必要な情報は説明に

2. **具体例を提示**
   - 再現手順、エラーログ、スクリーンショット
   - コード例（必要な場合）

3. **受け入れ基準を定義**
   - チェックリスト形式で何が完了なのかを明確に

4. **関連Issueをリンク**
   - `Closes #123` で自動クローズ
   - `Related to #456` で関連付け

5. **適切なラベルを付与**
   - タイプ・優先度・領域を正確に
   - 複数ラベルの組み合わせOK

### ❌ DONT

1. **不明確な説明**
   - 「なんかおかしい」では修正不可
   - 具体的に何がどうなっているか

2. **過度な詳細**
   - ログは必要最小限に
   - 長すぎるのは避ける

3. **単体の長いIssue**
   - 複数の異なる問題を1つのIssueに
   - タスク分割して複数Issueに

4. **ラベルなし**
   - ラベルがないと検索・分類が困難
   - 作成時に最低限のラベルを付与

5. **期限や担当者の無視**
   - 重要なIssueは期限を設定
   - 適切に担当者をアサイン

## Gitとの連携

### Issue番号をコミットに含める

```bash
git commit -m "fix(backend): 計算エラー修正

Issue: #42
```

### ブランチ名にも含める

```bash
git checkout -b fix/calculation-error-42
```

### PullRequestでのクローズ

```markdown
# PR説明

このPRで #42 を修正します。

Fixes #42
Relates to #40, #41
```

PR作成時に `Fixes #42` と記述すると、PRマージ時にIssueが自動クローズされます。

## Issue管理のフロー

```
1. Issue作成
   ↓ (ラベル・優先度・担当者を設定)
2. In Progress
   ↓ (実装開始)
3. PR作成
   ↓ (PR作成時にIssueリンク)
4. Review
   ↓ (コードレビュー)
5. Merged + Close
   ↓ (自動でIssueもクローズ)
6. Done
```

## チェックリスト：Issue作成時

- ✅ 同じIssueが既存でないか確認
- ✅ タイプを判別している（Bug/Feat/Doc/etc）
- ✅ 明確なタイトル
- ✅ 詳細な説明（テンプレート使用）
- ✅ 必要に応じて再現手順・スクリーンショット
- ✅ 適切なラベルを3-5個選択
- ✅ 優先度を設定
- ✅ 必要に応じて担当者をアサイン
- ✅ 期限を設定（重要なタスクは）
