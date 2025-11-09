# API クライアントと状態管理

このディレクトリには、バックエンドAPIとの通信および状態管理のためのコードが含まれています。

## 構成

### API クライアント (`api-client.ts`)

バックエンドAPIとの通信を担当するクライアントライブラリです。

#### 使用例

```typescript
import { financialDataAPI, goalsAPI, calculationsAPI } from '@/lib/api-client';

// 財務データの取得
const data = await financialDataAPI.get('user_123');

// 目標の作成
const goal = await goalsAPI.create({
  user_id: 'user_123',
  type: 'savings',
  title: '新車購入',
  target_amount: 3000000,
  target_date: '2025-12-31',
  current_amount: 500000,
  monthly_contribution: 100000,
  is_active: true,
});

// 資産推移の計算
const projection = await calculationsAPI.assetProjection({
  user_id: 'user_123',
  years: 30,
  monthly_income: 400000,
  monthly_expenses: 280000,
  current_savings: 1500000,
  investment_return: 5.0,
  inflation_rate: 2.0,
});
```

#### エラーハンドリング

```typescript
import { APIError } from '@/lib/api-client';

try {
  const data = await financialDataAPI.get('user_123');
} catch (error) {
  if (error instanceof APIError) {
    console.error('API Error:', error.message, error.status);
  }
}
```

### 状態管理 (Contexts)

React Contextを使用したグローバル状態管理を提供します。

#### FinancialDataContext

財務データの状態管理を担当します。

```typescript
import { useFinancialData } from '@/lib/contexts';

function MyComponent() {
  const {
    financialData,
    loading,
    error,
    fetchFinancialData,
    updateProfile,
  } = useFinancialData();

  useEffect(() => {
    fetchFinancialData('user_123');
  }, []);

  return (
    <div>
      {loading && <p>読み込み中...</p>}
      {error && <p>エラー: {error}</p>}
      {financialData && (
        <div>
          <p>月収: {financialData.profile?.monthly_income}</p>
        </div>
      )}
    </div>
  );
}
```

#### GoalsContext

目標の状態管理を担当します。

```typescript
import { useGoals } from '@/lib/contexts';

function GoalsList() {
  const {
    goals,
    loading,
    error,
    fetchGoals,
    createGoal,
    updateGoal,
    deleteGoal,
  } = useGoals();

  useEffect(() => {
    fetchGoals('user_123');
  }, []);

  return (
    <div>
      {goals.map(goal => (
        <div key={goal.id}>
          <h3>{goal.title}</h3>
          <p>目標額: {goal.target_amount}</p>
        </div>
      ))}
    </div>
  );
}
```

#### CalculationsContext

計算結果の状態管理を担当します。

```typescript
import { useCalculations } from '@/lib/contexts';

function AssetProjectionChart() {
  const {
    assetProjection,
    loading,
    error,
    calculateAssetProjection,
  } = useCalculations();

  const handleCalculate = async () => {
    await calculateAssetProjection({
      user_id: 'user_123',
      years: 30,
      monthly_income: 400000,
      monthly_expenses: 280000,
      current_savings: 1500000,
      investment_return: 5.0,
      inflation_rate: 2.0,
    });
  };

  return (
    <div>
      <button onClick={handleCalculate}>計算する</button>
      {loading && <p>計算中...</p>}
      {assetProjection && (
        <div>
          <p>最終資産額: {assetProjection.final_amount}</p>
        </div>
      )}
    </div>
  );
}
```

### AppProviders

すべてのコンテキストプロバイダーを統合したコンポーネントです。
`layout.tsx`で使用されています。

```typescript
import { AppProviders } from '@/lib/contexts/AppProviders';

export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        <AppProviders>
          {children}
        </AppProviders>
      </body>
    </html>
  );
}
```

## カスタムフック

### useUser

ユーザーセッション管理のためのフックです。

```typescript
import { useUser } from '@/lib/hooks';

function MyComponent() {
  const { userId, loading, clearUser } = useUser();

  if (loading) return <p>読み込み中...</p>;

  return (
    <div>
      <p>ユーザーID: {userId}</p>
      <button onClick={clearUser}>ログアウト</button>
    </div>
  );
}
```

## 環境変数

`.env.local`ファイルでAPIのベースURLを設定できます：

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

デフォルトは`http://localhost:8080/api`です。

## 型定義

すべてのAPI型は`@/types/api.ts`で定義されています。

## ベストプラクティス

1. **エラーハンドリング**: すべてのAPI呼び出しでtry-catchを使用
2. **ローディング状態**: ユーザーにフィードバックを提供
3. **型安全性**: TypeScriptの型を活用
4. **再利用性**: カスタムフックでロジックを共有
5. **パフォーマンス**: useCallbackとuseMemoで最適化
