# PDF Generator

このパッケージは財務レポートをPDF形式で生成する機能を提供します。

## 概要

PDF Generatorは、財務計画アプリケーションで生成される各種レポートをPDF形式に変換します。現在の実装では、HTML形式での出力をサポートしており、将来的に実際のPDFライブラリ（例：wkhtmltopdf、chromedp、gofpdf）と統合することができます。

## 実装されているジェネレーター

### HTMLGenerator

HTML形式でレポートを生成します。ブラウザで表示可能な形式で、印刷機能を使用してPDFに変換できます。

**サポートされているレポート:**
- 財務サマリーレポート (`GenerateFinancialSummaryPDF`)
- 包括的レポート (`GenerateComprehensivePDF`)
- 資産推移レポート (`GenerateAssetProjectionPDF`)
- 目標進捗レポート (`GenerateGoalsProgressPDF`)
- 退職計画レポート (`GenerateRetirementPlanPDF`)

### JSONGenerator

JSON形式でレポートデータを出力します。デバッグやAPI統合に便利です。

## 使用方法

```go
import (
    "github.com/financial-planning-calculator/backend/infrastructure/pdf"
    "github.com/financial-planning-calculator/backend/application/usecases"
)

// HTMLジェネレーターの作成
generator := pdf.NewHTMLGenerator()

// 財務サマリーレポートの生成
report := &usecases.FinancialSummaryReport{
    // レポートデータ
}

htmlBytes, err := generator.GenerateFinancialSummaryPDF(report)
if err != nil {
    // エラー処理
}

// HTMLをファイルに保存またはHTTPレスポンスとして返す
```

## レポートの特徴

### 財務サマリーレポート
- 財務健全性スコア（0-100点）
- 現在の財務状況（収入、支出、貯蓄、資産）
- 主要指標（貯蓄率、投資利回り、総資産）
- 推奨事項と警告

### 包括的レポート
- エグゼクティブサマリー
- 財務サマリー
- 資産推移予測
- 目標進捗状況
- アクションプラン（短期・中期・長期）

## スタイリング

生成されるHTMLには以下のスタイルが適用されます：
- プロフェッショナルなフォント（Helvetica, Arial）
- カラースキーム（青系統のプライマリカラー）
- レスポンシブなレイアウト
- 印刷最適化（ページブレーク、マージン）

## 将来の拡張

### PDF生成ライブラリの統合

実際のPDFファイルを生成するには、以下のライブラリの統合を検討できます：

1. **wkhtmltopdf** - HTMLをPDFに変換
   ```bash
   # インストール
   brew install wkhtmltopdf  # macOS
   apt-get install wkhtmltopdf  # Ubuntu
   ```

2. **chromedp** - Chromeヘッドレスブラウザを使用
   ```go
   import "github.com/chromedp/chromedp"
   ```

3. **gofpdf** - ネイティブGo PDFライブラリ
   ```go
   import "github.com/jung-kurt/gofpdf"
   ```

### チャートとグラフの追加

将来的には、以下の機能を追加できます：
- Chart.jsやD3.jsを使用したインタラクティブなグラフ
- SVGベースのチャート埋め込み
- 画像としてのグラフ生成

### カスタマイズ機能

- ユーザー定義のテンプレート
- カラースキームのカスタマイズ
- ロゴやブランディングの追加
- 多言語サポート

## テスト

```bash
# ユニットテストの実行
go test ./infrastructure/pdf/...

# 統合テストの実行
go test -tags=integration ./infrastructure/pdf/...
```

## パフォーマンス

- HTMLジェネレーター: ~10ms/レポート
- 将来のPDF生成: ~100-500ms/レポート（ライブラリに依存）

## セキュリティ

- XSS攻撃防止のため、すべてのユーザー入力をエスケープ
- ファイルパストラバーサル攻撃の防止
- 生成されたファイルの有効期限設定（24時間）
