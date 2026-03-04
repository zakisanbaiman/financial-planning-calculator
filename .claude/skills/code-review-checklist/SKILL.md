---
name: code-review-checklist
description: PRマージ前のセルフレビューチェックリスト。過去のfixコミットから抽出した再発しやすい問題を確認する時に使用。
---

# コードレビュー セルフチェックリスト

過去のリポジトリのfixコミット履歴から抽出した、**実際に発生した問題パターン**。
PRを出す前に必ずこのリストを確認すること。

---

## 🔴 [must] バックエンド (Go)

### 設定・環境変数
- [ ] ハードコーディングされた値がないか？（ポート、シークレット、パスワード等）
  - 参考: `fix: 設定値のハードコーディングの排除` (#46)
  - **対策**: 全て `config/server.go` または `config/database.go` 経由で `getEnv()` を使う
  - 危険な例: `secret := "my-secret"`, `port := "8080"` を直接コードに書く

- [ ] `go mod tidy` を実行したか？
  - 参考: `fix: run go mod tidy` (#75), `fix: Ensure go.mod and go.sum are tidy` (#33)
  - **対策**: コミット前に必ず `cd backend && go mod tidy` を実行

### テスト
- [ ] テストで `ServerConfig` が必要なミドルウェアを使うケースで `ServerConfig` をセットアップしているか？
  - 参考: `fix: Add ServerConfig to test setups for OAuth middleware`
  - **対策**: Echo のテストセットアップで `cfg := config.LoadServerConfig()` を忘れずに

- [ ] モックの `On()` と実際の呼び出しが一致しているか？
  - 参考: `fix: update mock behavior for GetFinancialPlan to handle missing data`
  - **対策**: `AssertExpectations(t)` はモックのすべての期待が呼ばれたことを検証する

### APIレスポンス
- [ ] レスポンス構造体のフィールドが `int/float64` 等のプリミティブ型か（VO型をそのまま返していないか）？
  - 参考: `fix: convert financial profile, retirement data, and emergency fund to primitive types in response`
  - **対策**: Controllerレイヤーでは必ず値オブジェクトから `.Value()` や `.Int()` で変換する

---

## 🔴 [must] フロントエンド (Next.js / TypeScript)

### 型・プロパティ名
- [ ] APIレスポンスのプロパティ名が型定義と一致しているか？
  - 参考: `fix: Replace all goal.type references with goal.goal_type` (#33系)
  - **対策**: `types/api.ts` の型定義と `src/lib/api-client.ts` のレスポンスを照合する

- [ ] 変数の重複宣言がないか？（特に同一スコープ内の `const/let`）
  - 参考: `fix: 2fa-verify/page.tsxの重複宣言エラーを修正`
  - **対策**: `npm run lint` でTypeScriptエラーを事前確認

### Next.js 固有
- [ ] `useSearchParams()` / `usePathname()` は `<Suspense>` でラップされているか？
  - 参考: `fix: Wrap useSearchParams in Suspense boundary for /auth/callback page`
  - **対策**: これらのフックを使うコンポーネントは必ずSuspenseバウンダリ内に配置

- [ ] ダークモード対応: SSR時に `localStorage` を直接参照していないか？
  - 参考: `fix: 画面遷移時の白飛びを修正`
  - **対策**: テーマ初期化は `<head>` インラインスクリプトで行い、hydrationミスマッチを防ぐ

### 認証状態
- [ ] OAuth/2FAコールバック後、認証状態の更新を **待ってから** 画面遷移しているか？
  - 参考: `fix: OAuth後の認証状態がヘッダーに反映されない問題を修正`
  - **対策**: `setAuthData()` → `await` または `useEffect` で状態確認後に `router.push()`

- [ ] URLにトークンが含まれる場合、遷移後にURLをクリアしているか？
  - 参考: `fix: OAuth後の認証状態がヘッダーに反映されない問題を修正`
  - **対策**: `router.replace()` でトークンパラメータを除去する

### null/undefined
- [ ] APIから返る値に `null` が含まれる可能性を考慮しているか？
  - 参考: `fix: handle potential null values in financial data calculations`
  - **対策**: オプショナルチェーン `?.` や `?? 0` を使う

### テスト
- [ ] Context依存のコンポーネント/フックのテストに、必要なProviderをラップしているか？
  - 参考: `fix: FinancialDataContext テストに GuestModeProvider を追加する`
  - 参考: `fix: useUser テストを AuthContext/GuestModeContext のモックに対応させる`
  - **対策**: テスト用ラッパーに `AuthProvider`, `GuestModeProvider` を正しく含める

---

## 🟡 [recommend] CI/CD

- [ ] `go mod tidy` の結果が差分なしか確認したか？
  - **コマンド**: `cd backend && go mod tidy && git diff go.mod go.sum`

- [ ] E2Eテストのタイムアウトが長すぎないか？
  - 参考: `fix: optimize E2E tests for CI by reducing timeout`
  - **基準**: playwright のタイムアウトは CI 環境に合わせて設定

- [ ] `render.yaml` の変更時、Hobby プランの制約（preview environment 不可等）を考慮しているか？
  - 参考: `fix(render): remove preview environment settings for Hobby plan`

---

## 🟡 [recommend] セキュリティ

- [ ] デフォルト値のシークレットが本番環境に流出しないよう `.env.example` に記載されているか？
- [ ] OAuth のアカウント自動リンクを実装していないか（メール一致だけでリンクするのは危険）？
  - 参考: `fix: Remove automatic OAuth account linking to prevent account takeover` (#80)

---

## チェック実行コマンド

```bash
# バックエンド
cd backend
go mod tidy
git diff --exit-code go.mod go.sum  # 差分なしなら OK
go vet ./...
golangci-lint run

# フロントエンド
cd frontend
npm run lint
npx tsc --noEmit  # 型エラー確認
npm test -- --watchAll=false  # テスト実行
```
