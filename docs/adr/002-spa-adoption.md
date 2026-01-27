# ADR 002: SPAアーキテクチャの採用について

## ステータス

採択済み (2026-01-27)

## 背景

財務計画計算機プロジェクトは、Next.js 14のApp Routerを使用して構築されており、現在SSR（Server-Side Rendering）とCSR（Client-Side Rendering）の混在アーキテクチャとなっています。プロジェクトのユーザー体験向上と開発効率を考慮し、完全なSPA（Single Page Application）化を検討することになりました。

### 現状の分析

**技術スタック:**
- フロントエンド: Next.js 14 (App Router)
- バックエンド: Go + Echo Framework (REST API)
- 全47個のTSXコンポーネント中、18個が`'use client'`ディレクティブを使用
- 認証、財務データ管理、計算機能など、ほぼすべての主要ページがクライアントサイドレンダリング

**現在のレンダリング戦略:**
- ルートレイアウト: SSR（メタデータの提供）
- ほぼすべてのページコンポーネント: CSR（`'use client'`使用）
- ナビゲーション、フォーム、チャート表示: CSR

## 決定

**Next.js 14のApp Routerを維持し、現在のハイブリッドアーキテクチャを継続する。完全なSPAへの移行は行わない。**

## 理由

### 現在のアーキテクチャを維持する理由

#### 1. **SEOとパフォーマンスの最適化**
   - 財務計画アプリは、ユーザーが検索エンジンから見つける可能性が高い
   - ランディングページ（`/`）とドキュメントページはSSRで初期表示が高速
   - メタデータ（タイトル、description）の動的生成により、SEOの柔軟性を保持
   - 初回アクセス時のTime to First Byte (TTFB)が改善される

#### 2. **段階的な静的生成（ISR）の活用可能性**
   - 将来的に静的コンテンツ（ブログ、ヘルプページ等）を追加する場合、ISRが活用できる
   - データの更新頻度に応じた最適なレンダリング戦略を選択可能

#### 3. **現在のアーキテクチャは既に最適化されている**
   - 主要な機能ページ（dashboard, calculations, goals等）は既に`'use client'`でCSR化
   - インタラクティブな機能は既にクライアント側で実行されている
   - ユーザー体験の観点では、実質的にSPAと同等の操作感を実現

#### 4. **Next.jsの強力な機能を活用**
   - 自動コード分割（Automatic Code Splitting）
   - ルートベースの遅延読み込み（Route-based Lazy Loading）
   - 画像最適化（Next.js Image Component）
   - APIルートの簡潔な実装（現在はRewritesで対応）

#### 5. **移行コストとリスク**
   - 完全SPA化のために`next export`やReact Router等への移行は、大きなリファクタリングが必要
   - Next.jsの豊富な機能（画像最適化、フォント最適化等）を放棄することになる
   - バックエンドAPIとの統合が複雑化する可能性
   - デプロイ構成の変更が必要（現在のDocker + standalone構成が最適）

#### 6. **開発者体験**
   - Next.js 14のApp Routerは、直感的なファイルベースルーティング
   - Server ComponentsとClient Componentsの明示的な分離により、保守性が高い
   - TypeScript、ESLint、Prettierとの統合が優れている

### 完全SPA化のデメリット

#### 技術的課題
1. **SEOの劣化**
   - JavaScript無効環境での表示不可
   - クローラーの実行コストが増加
   - メタタグの動的生成が困難

2. **初回ロードの遅延**
   - 全JavaScriptバンドルのダウンロードが必要
   - 初回表示までの時間（FCP: First Contentful Paint）が遅くなる

3. **ブラウザ履歴管理の複雑化**
   - React Routerなどの追加ライブラリが必要
   - ブラウザの戻る/進むボタンの挙動を独自実装

4. **Next.jsの主要機能の喪失**
   - 画像最適化、フォント最適化
   - 自動コード分割の一部機能
   - Server Actionsの活用不可

### 現在のハイブリッドアーキテクチャの利点

1. **最適なパフォーマンス**
   - 静的コンテンツはSSR/SSG
   - 動的コンテンツはCSR
   - 必要な部分だけクライアントサイドで実行

2. **柔軟性**
   - ページごとにレンダリング戦略を選択可能
   - 将来的な要件変更に対応しやすい

3. **開発効率**
   - Next.jsのエコシステムをフル活用
   - 学習コストが低い（ドキュメントが充実）

4. **ユーザー体験**
   - 初回アクセスが高速（SSR）
   - 以降のナビゲーションは瞬時（CSR）
   - プログレッシブエンハンスメント

## 代替案

### 1. **完全SPA化（Next.js + `output: 'export'`）**
   - **利点**: 完全な静的サイト生成、CDN配信が容易
   - **欠点**: SSR不可、動的ルーティング制限、画像最適化機能の制限
   - **評価**: 本プロジェクトには不適切（動的データが多い）

### 2. **React Router + Vite への移行**
   - **利点**: 完全なSPA、ビルドが高速
   - **欠点**: Next.jsの機能を全て失う、大規模なリファクタリングが必要
   - **評価**: リスクが高く、メリットが不明確

### 3. **Remix への移行**
   - **利点**: モダンなフルスタックフレームワーク
   - **欠点**: 学習コスト、エコシステムが小さい
   - **評価**: 現時点で移行する理由が不十分

### 4. **現状維持 + パフォーマンス最適化（採用）**
   - **利点**: リスクが低い、既存の利点を維持、段階的改善が可能
   - **欠点**: なし（現在のアーキテクチャが適切）
   - **評価**: 最適な選択

## 結果

### 期待される成果

1. **安定性の維持**
   - 現在の動作を維持しながら、段階的な改善が可能
   - リファクタリングリスクを回避

2. **SEOとパフォーマンスの両立**
   - 検索エンジン最適化を保持
   - 初回表示速度とインタラクティブ性の両立

3. **開発効率の維持**
   - Next.jsのベストプラクティスに従った開発
   - 豊富なドキュメントとコミュニティサポート

4. **将来の拡張性**
   - 必要に応じてSSG、ISR、PPRなどの機能を追加可能
   - 段階的な最適化戦略を実行可能

### 今後の改善項目

現在のハイブリッドアーキテクチャを維持しながら、以下の最適化を実施：

1. **コード分割の最適化**
   - 動的インポート（`next/dynamic`）の活用
   - ルートごとのバンドルサイズの監視

2. **キャッシュ戦略の改善**
   - APIレスポンスのクライアントサイドキャッシング
   - React QueryやSWRの導入検討

3. **プリフェッチの活用**
   - `<Link prefetch>`による次ページの事前読み込み
   - ユーザー行動に基づくインテリジェントなプリフェッチ

4. **パフォーマンスモニタリング**
   - Core Web Vitalsの定期的な計測
   - リアルユーザーモニタリング（RUM）の導入

5. **段階的な静的生成**
   - 頻繁に変更されないページ（About、Help等）のSSG化
   - ISRによる定期的な再生成

## 参考資料

- [Next.js 14 App Router Documentation](https://nextjs.org/docs)
- [Rendering: Server Components](https://nextjs.org/docs/app/building-your-application/rendering/server-components)
- [Rendering: Client Components](https://nextjs.org/docs/app/building-your-application/rendering/client-components)
- [SPA vs SSR vs SSG Comparison](https://web.dev/rendering-on-the-web/)
- [Next.js Performance Best Practices](https://nextjs.org/docs/app/building-your-application/optimizing)
