# 実装サマリー: 状態管理とAPI通信

## 実装内容

タスク 6.3「状態管理とAPI通信」の実装が完了しました。

### 1. API型定義 (`types/api.ts`)

バックエンドAPIとの通信に必要なすべての型を定義：

- **財務データ型**: `FinancialData`, `FinancialProfile`, `RetirementData`, `EmergencyFund`
- **目標型**: `Goal`, `GoalType`
- **計算リクエスト型**: `AssetProjectionRequest`, `RetirementCalculationRequest`, 等
- **計算レスポンス型**: `AssetProjectionResponse`, `RetirementCalculationResponse`, 等
- **レポート型**: `ReportRequest`, `FinancialSummaryReport`

### 2. APIクライアント (`lib/api-client.ts`)

バックエンドAPIとの通信を担当する完全な型安全クライアント：

#### 実装されたAPI

- **財務データAPI** (`financialDataAPI`)
  - `create`: 財務データ作成
  - `get`: 財務データ取得
  - `updateProfile`: 財務プロファイル更新
  - `updateRetirement`: 退職データ更新
  - `updateEmergencyFund`: 緊急資金更新
  - `delete`: 財務データ削除

- **計算API** (`calculationsAPI`)
  - `assetProjection`: 資産推移計算
  - `retirement`: 老後資金計算
  - `emergencyFund`: 緊急資金計算
  - `goalProjection`: 目標達成計算

- **目標API** (`goalsAPI`)
  - `create`: 目標作成
  - `list`: 目標一覧取得
  - `get`: 目標取得
  - `update`: 目標更新
  - `updateProgress`: 目標進捗更新
  - `delete`: 目標削除
  - `getRecommendations`: 目標推奨事項取得
  - `analyzeFeasibility`: 目標実現可能性分析

- **レポートAPI** (`reportsAPI`)
  - `financialSummary`: 財務サマリーレポート生成
  - `getPDF`: PDFレポート取得

#### 機能

- カスタムエラークラス (`APIError`)
- 統一されたリクエストハンドリング
- 環境変数によるベースURL設定
- 適切なエラーハンドリング

### 3. 状態管理 (React Context)

#### FinancialDataContext (`lib/contexts/FinancialDataContext.tsx`)

財務データの状態管理：

- 状態: `financialData`, `loading`, `error`
- アクション: `fetchFinancialData`, `createFinancialData`, `updateProfile`, `updateRetirement`, `updateEmergencyFund`, `deleteFinancialData`
- カスタムフック: `useFinancialData()`

#### GoalsContext (`lib/contexts/GoalsContext.tsx`)

目標の状態管理：

- 状態: `goals`, `loading`, `error`
- アクション: `fetchGoals`, `createGoal`, `updateGoal`, `updateGoalProgress`, `deleteGoal`
- カスタムフック: `useGoals()`

#### CalculationsContext (`lib/contexts/CalculationsContext.tsx`)

計算結果の状態管理：

- 状態: `assetProjection`, `retirementCalculation`, `emergencyFund`, `goalProjection`, `loading`, `error`
- アクション: `calculateAssetProjection`, `calculateRetirement`, `calculateEmergencyFund`, `calculateGoalProjection`, `clearCalculations`
- カスタムフック: `useCalculations()`

### 4. 統合プロバイダー (`lib/contexts/AppProviders.tsx`)

すべてのコンテキストプロバイダーを統合し、アプリケーション全体で利用可能にします。

### 5. カスタムフック

#### useUser (`lib/hooks/useUser.ts`)

ユーザーセッション管理：

- ローカルストレージを使用した簡易実装
- 本番環境では認証システムと統合する必要あり
- 機能: `userId`, `loading`, `clearUser`

### 6. ドキュメント

- **README.md**: 使用方法とベストプラクティス
- **ExampleUsage.tsx**: 実装の参考例

### 7. レイアウト統合

`app/layout.tsx`を更新し、`AppProviders`でアプリケーション全体をラップ。

## 使用方法

### 基本的な使用例

```typescript
'use client';

import { useFinancialData, useGoals, useCalculations } from '@/lib/contexts';
import { useUser } from '@/lib/hooks';
import { useEffect } from 'react';

export default function MyPage() {
  const { userId } = useUser();
  const { financialData, fetchFinancialData } = useFinancialData();
  const { goals, fetchGoals } = useGoals();
  const { calculateAssetProjection, assetProjection } = useCalculations();

  useEffect(() => {
    if (userId) {
      fetchFinancialData(userId);
      fetchGoals(userId);
    }
  }, [userId]);

  return (
    <div>
      {/* コンポーネントの実装 */}
    </div>
  );
}
```

## 環境設定

`.env.local`ファイルを作成してAPIのベースURLを設定：

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## テスト結果

- ✅ TypeScript型チェック: 合格
- ✅ すべてのファイルがエラーなくコンパイル
- ✅ 型安全性が保証されている

## 次のステップ

この実装により、以下のタスクの実装が可能になりました：

- タスク 7: 財務データ入力機能
- タスク 8: 計算・可視化機能
- タスク 9: 目標設定・進捗管理機能
- タスク 10: レポート生成機能

各ページコンポーネントで`useFinancialData()`, `useGoals()`, `useCalculations()`フックを使用して、
バックエンドAPIとの通信と状態管理を簡単に実装できます。

## 注意事項

1. **認証**: 現在の`useUser`フックは簡易実装です。本番環境では適切な認証システムと統合してください。
2. **エラーハンドリング**: すべてのAPI呼び出しで適切なエラーハンドリングを実装してください。
3. **ローディング状態**: ユーザーエクスペリエンス向上のため、ローディング状態を適切に表示してください。
4. **環境変数**: 本番環境では適切なAPIエンドポイントを設定してください。
