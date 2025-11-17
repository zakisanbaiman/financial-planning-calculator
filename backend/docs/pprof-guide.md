# pprofガイド

## 概要

pprofはGoの標準プロファイリングツールで、アプリケーションのパフォーマンス分析に使用します。

## 有効化

開発環境では自動的に有効化されています（`ENABLE_PPROF=true`）。

```bash
# Docker環境
make up

# pprofサーバーが http://localhost:6060 で起動
```

## 基本的な使い方

### 1. ブラウザで確認

```bash
# プロファイル一覧
open http://localhost:6060/debug/pprof/

# ゴルーチン一覧（テキスト形式）
open http://localhost:6060/debug/pprof/goroutine?debug=1

# メモリ使用状況（テキスト形式）
open http://localhost:6060/debug/pprof/heap?debug=1
```

### 2. コマンドラインで分析

#### CPUプロファイル（30秒間）

```bash
# プロファイル取得
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# インタラクティブモードで分析
(pprof) top        # CPU使用率トップ10
(pprof) top -cum   # 累積時間でトップ10
(pprof) list main  # main関数の詳細
(pprof) web        # グラフ表示（graphviz必要）
```

#### メモリプロファイル

```bash
# ヒープメモリ分析
go tool pprof http://localhost:6060/debug/pprof/heap

(pprof) top        # メモリ使用量トップ10
(pprof) list <関数名>  # 特定関数の詳細
```

#### ゴルーチン

```bash
# ゴルーチン分析
go tool pprof http://localhost:6060/debug/pprof/goroutine

(pprof) top        # ゴルーチン数トップ10
```

### 3. Web UIで可視化

最も使いやすい方法です：

```bash
# graphvizをインストール（初回のみ）
brew install graphviz

# Web UIを起動（CPUプロファイル）
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30

# ブラウザで http://localhost:8081 が自動的に開きます
```

**Web UIの機能:**
- Flame Graph: 関数呼び出しの階層を視覚化
- Top: CPU/メモリ使用量のランキング
- Graph: 関数間の呼び出し関係
- Source: ソースコードレベルの分析

## 実践例

### 例1: APIエンドポイントのパフォーマンス測定

```bash
# 1. 負荷をかける
for i in {1..1000}; do
  curl http://localhost:8080/api/calculations/asset-projection \
    -H "Content-Type: application/json" \
    -d '{"monthly_income":500000,"monthly_expenses":300000,...}' &
done

# 2. CPUプロファイルを取得（30秒間）
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30

# 3. Web UIで分析
# - どの関数がCPUを使っているか
# - ボトルネックはどこか
```

### 例2: メモリリークの検出

```bash
# 1. 初期状態のメモリプロファイル
curl http://localhost:6060/debug/pprof/heap > heap_before.prof

# 2. アプリケーションを使用（負荷をかける）
# ... 操作 ...

# 3. 使用後のメモリプロファイル
curl http://localhost:6060/debug/pprof/heap > heap_after.prof

# 4. 差分を分析
go tool pprof -http=:8081 -base heap_before.prof heap_after.prof
```

### 例3: ゴルーチンリークの検出

```bash
# ゴルーチン数を確認
curl http://localhost:6060/debug/pprof/goroutine?debug=1 | head -1

# 時間をおいて再度確認
# ゴルーチン数が増え続けている場合はリークの可能性
```

## pprofコマンド一覧

インタラクティブモード内で使えるコマンド：

```
top [N]          - 上位N件を表示（デフォルト10）
top -cum [N]     - 累積時間で上位N件
list <関数名>     - 関数のソースコードと統計
web              - グラフをブラウザで表示
pdf              - PDFファイルを生成
png              - PNG画像を生成
svg              - SVG画像を生成
peek <関数名>     - 関数の呼び出し元/先を表示
help             - ヘルプ表示
quit             - 終了
```

## プロファイルの種類

| プロファイル | URL | 用途 |
|------------|-----|------|
| CPU | `/debug/pprof/profile` | CPU使用率の分析 |
| Heap | `/debug/pprof/heap` | メモリ使用量の分析 |
| Goroutine | `/debug/pprof/goroutine` | ゴルーチン数の分析 |
| Allocs | `/debug/pprof/allocs` | メモリアロケーションの分析 |
| Block | `/debug/pprof/block` | ブロッキング操作の分析 |
| Mutex | `/debug/pprof/mutex` | ミューテックス競合の分析 |
| Threadcreate | `/debug/pprof/threadcreate` | スレッド生成の分析 |

## ベストプラクティス

1. **本番環境では無効化**
   - pprofはセキュリティリスクがあるため、本番では必ず無効化
   - 環境変数 `ENABLE_PPROF=false` を設定

2. **定期的なプロファイリング**
   - 新機能追加時にパフォーマンスを確認
   - リリース前に必ずプロファイリング

3. **ベースラインの確立**
   - 正常時のプロファイルを保存
   - 問題発生時に比較

4. **負荷テストと組み合わせ**
   - 実際の負荷をかけた状態でプロファイリング
   - より正確なボトルネックを特定

## トラブルシューティング

### pprofにアクセスできない

```bash
# pprofが有効か確認
docker-compose logs backend | grep pprof

# ポートが開いているか確認
curl http://localhost:6060/debug/pprof/
```

### graphvizがない

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Windows
# https://graphviz.org/download/ からインストール
```

## 参考リンク

- [公式ドキュメント](https://pkg.go.dev/net/http/pprof)
- [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)
- [Practical Go: Real world advice for writing maintainable Go programs](https://dave.cheney.net/practical-go/presentations/qcon-china.html#_profiling)
