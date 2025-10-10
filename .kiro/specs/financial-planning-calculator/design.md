# 設計文書

## 概要

財務計画計算アプリケーションは、ユーザーの現在の財務状況から将来の資産推移、老後資金、緊急時資金を計算・可視化するWebアプリケーションです。シンプルで直感的なインターフェースを通じて、複雑な財務計算を分かりやすく提示します。

## アーキテクチャ

### システム構成
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   (Next.js)     │◄──►│   (Go)          │◄──►│  (PostgreSQL)   │
│                 │    │                 │    │                 │
│ - Pages/Routes  │    │ - API Routes    │    │ - User Data     │
│ - Components    │    │ - Calculations  │    │ - Financial     │
│ - State Mgmt    │    │ - Validation    │    │   Records       │
│ - Charts        │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 技術スタック
- **フロントエンド**: Next.js + TypeScript
- **状態管理**: React Context + useReducer
- **スタイリング**: Tailwind CSS
- **チャート**: Recharts
- **バックエンド**: Go + Echo
- **データベース**: PostgreSQL
- **ORM**: GORM
- **API**: RESTful API
- **認証**: JWT（将来的な拡張用）

## コンポーネントとインターフェース

### フロントエンドコンポーネント

#### 1. ページとレイアウト
- `pages/_app.tsx`: アプリケーション全体の設定
- `pages/index.tsx`: ダッシュボードページ
- `pages/input.tsx`: 財務情報入力ページ
- `pages/projection.tsx`: 資産推移表示ページ
- `pages/retirement.tsx`: 老後資金計算ページ
- `pages/goals.tsx`: 目標設定・管理ページ
- `components/Layout`: 共通レイアウトコンポーネント
- `components/Header`: ナビゲーションとタイトル
- `components/Sidebar`: メニューとナビゲーション

#### 2. 入力コンポーネント
- `FinancialInputForm`: 基本財務情報入力
- `IncomeInput`: 収入入力フィールド
- `ExpenseInput`: 支出入力フィールド
- `SavingsInput`: 貯蓄額入力フィールド
- `GoalSettingForm`: 目標設定フォーム

#### 3. 計算・表示コンポーネント
- `AssetProjectionChart`: 資産推移グラフ
- `RetirementCalculator`: 老後資金計算器
- `EmergencyFundCalculator`: 緊急資金計算器
- `ProgressTracker`: 目標進捗表示
- `SummaryDashboard`: 総合ダッシュボード

#### 4. 可視化コンポーネント
- `LineChart`: 時系列データ用
- `PieChart`: 支出内訳用
- `BarChart`: 比較データ用
- `ProgressBar`: 進捗表示用

### バックエンドAPI

#### エンドポイント設計
```
GET    /api/users/:id/profile          # ユーザー情報取得
PUT    /api/users/:id/profile          # ユーザー情報更新

GET    /api/users/:id/financial-data   # 財務データ取得
POST   /api/users/:id/financial-data   # 財務データ作成
PUT    /api/users/:id/financial-data   # 財務データ更新

POST   /api/calculations/projection    # 資産推移計算
POST   /api/calculations/retirement    # 老後資金計算
POST   /api/calculations/emergency     # 緊急資金計算

GET    /api/users/:id/goals           # 目標一覧取得
POST   /api/users/:id/goals           # 目標作成
PUT    /api/users/:id/goals/:goalId   # 目標更新
DELETE /api/users/:id/goals/:goalId   # 目標削除
```

## データモデル

### User（ユーザー）
```typescript
interface User {
  id: string;
  name: string;
  email: string;
  createdAt: Date;
  updatedAt: Date;
}
```

### FinancialProfile（財務プロフィール）
```typescript
interface FinancialProfile {
  id: string;
  userId: string;
  monthlyIncome: number;
  monthlyExpenses: ExpenseBreakdown;
  currentSavings: SavingsBreakdown;
  age: number;
  retirementAge: number;
  createdAt: Date;
  updatedAt: Date;
}

interface ExpenseBreakdown {
  housing: number;
  food: number;
  transportation: number;
  utilities: number;
  entertainment: number;
  healthcare: number;
  other: number;
}

interface SavingsBreakdown {
  cash: number;
  investments: number;
  retirement401k: number;
  other: number;
}
```

### Goal（目標）
```typescript
interface Goal {
  id: string;
  userId: string;
  type: 'savings' | 'retirement' | 'emergency' | 'custom';
  targetAmount: number;
  currentAmount: number;
  targetDate: Date;
  description: string;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}
```

### CalculationResult（計算結果）
```typescript
interface ProjectionResult {
  years: number[];
  amounts: number[];
  realAmounts: number[]; // インフレ調整後
  parameters: {
    initialAmount: number;
    monthlyContribution: number;
    annualReturn: number;
    inflationRate: number;
  };
}

interface RetirementAnalysis {
  requiredAmount: number;
  projectedAmount: number;
  shortfall: number;
  monthlyPensionIncome: number;
  yearsToRetirement: number;
  recommendedMonthlySavings: number;
}
```

## エラーハンドリング

### フロントエンド
- 入力値検証（リアルタイム）
- APIエラーの適切な表示
- ネットワークエラーハンドリング
- フォールバック表示

### バックエンド
- 入力値検証とサニタイゼーション
- データベースエラーハンドリング
- 計算エラー（ゼロ除算等）の処理
- 適切なHTTPステータスコード返却

### エラーレスポンス形式
```typescript
interface ErrorResponse {
  error: {
    code: string;
    message: string;
    details?: any;
  };
}
```

## テスト戦略

### 単体テスト
- **計算ロジック**: 複利計算、インフレ調整、年金計算の精度テスト
- **コンポーネント**: React コンポーネントの動作テスト
- **API**: エンドポイントの入出力テスト
- **バリデーション**: 入力値検証ロジックのテスト

### 統合テスト
- **API統合**: フロントエンド-バックエンド間の通信テスト
- **データフロー**: データの作成から表示までの一連の流れ
- **計算精度**: 実際のデータを使った計算結果の検証

### E2Eテスト
- **ユーザーフロー**: 典型的な使用シナリオの自動テスト
- **クロスブラウザ**: 主要ブラウザでの動作確認
- **レスポンシブ**: モバイル・デスクトップでの表示確認

### テストデータ
- **サンプルユーザー**: 様々な年齢・収入レベルのテストケース
- **境界値**: 極端な値（高収入、低収入、高齢等）でのテスト
- **エラーケース**: 無効な入力値でのエラーハンドリングテスト