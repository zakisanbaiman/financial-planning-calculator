# 設計文書

## 概要

財務計画計算機は、ユーザーの現在の財務状況を基に将来の資産推移、老後資金、緊急時資金を計算・可視化するWebアプリケーションです。Next.js/TypeScriptをフロントエンド、Go/Echoをバックエンド、PostgreSQLをデータベースとして構築し、計算結果をインタラクティブなグラフで表示します。

## 設計思想

### 1. ユーザー中心設計 (User-Centric Design)
- **直感的なUI**: 複雑な財務計算を分かりやすいビジュアルで表現
- **段階的な情報入力**: ユーザーが圧倒されないよう、必要な情報を段階的に収集
- **即座のフィードバック**: 入力と同時に結果を表示し、変更の影響を即座に確認可能

### 2. 透明性と信頼性 (Transparency & Trust)
- **計算ロジックの明示**: どのような前提で計算しているかを明確に表示
- **データの可視化**: 数値だけでなく、グラフで直感的に理解できる表現
- **仮定の調整可能性**: インフレ率、投資利回りなどの前提条件をユーザーが調整可能

### 3. 拡張性とメンテナンス性 (Scalability & Maintainability)
- **ドメイン駆動設計**: ビジネスドメインの複雑さをコードで表現し、ドメインエキスパートとの共通言語を構築
- **クリーンアーキテクチャ**: 依存関係の方向を制御し、ビジネスロジックを外部詳細から分離
- **モジュラー設計**: 計算ロジック、UI、データ層を明確に分離
- **API First**: OpenAPI仕様による契約駆動開発
- **型安全性**: TypeScript/Goによる静的型チェック
- **テスト駆動**: 財務計算の正確性を保証する包括的なテスト

### 4. パフォーマンスと応答性 (Performance & Responsiveness)
- **リアルタイム計算**: 入力変更時の即座な再計算
- **効率的なデータ構造**: 大量の時系列データを効率的に処理
- **レスポンシブデザイン**: デスクトップ・モバイル両対応

### 5. セキュリティとプライバシー (Security & Privacy)
- **データ最小化**: 必要最小限の個人情報のみ収集
- **ローカル計算優先**: 可能な限りクライアント側で計算を実行
- **暗号化**: 機密性の高い財務データの適切な保護

### 6. 教育的価値 (Educational Value)
- **学習支援**: 財務計画の基本概念を理解できるヘルプ機能
- **シナリオ比較**: 異なる貯蓄戦略の比較機能
- **推奨事項**: データに基づく具体的なアドバイス提供

## アーキテクチャ

### システム構成

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   (Next.js/TS)  │◄──►│   (Go/Echo)     │◄──►│  (PostgreSQL)   │
│                 │    │                 │    │                 │
│ - Pages/Routes  │    │ - REST API      │    │ - User Data     │
│ - Components    │    │ - Calculations  │    │ - Financial     │
│ - State Mgmt    │    │ - Validation    │    │   Records       │
│ - Charts        │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### クリーンアーキテクチャ

バックエンドはクリーンアーキテクチャの原則に従って設計し、依存関係の方向を制御します：

```
┌─────────────────────────────────────────────────────────────┐
│                     Frameworks & Drivers                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Web (Echo)    │  │   Database      │  │   External APIs │ │
│  │   - Handlers    │  │   - PostgreSQL  │  │   - PDF Gen     │ │
│  │   - Middleware  │  │   - Migrations  │  │   - Email       │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                Interface Adapters                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Controllers   │  │   Presenters    │  │   Gateways      │ │
│  │   - HTTP        │  │   - JSON        │  │   - Repository  │ │
│  │   - Validation  │  │   - Error       │  │   - External    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                   Application Business Rules                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Use Cases     │  │   Services      │  │   Interfaces    │ │
│  │   - Calculate   │  │   - Financial   │  │   - Repository  │ │
│  │   - Manage      │  │   - Validation  │  │   - External    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                   Enterprise Business Rules                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Entities      │  │   Value Objects │  │   Domain Rules  │ │
│  │   - Financial   │  │   - Money       │  │   - Validation  │ │
│  │   - Goal        │  │   - Period      │  │   - Calculation │ │
│  │   - User        │  │   - Rate        │  │   - Business    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### ドメイン駆動設計の適用

#### ユビキタス言語 (Ubiquitous Language)

財務計画ドメインの共通言語を定義：

- **財務プロファイル (Financial Profile)**: ユーザーの現在の財務状況
- **資産推移 (Asset Projection)**: 将来の資産変化予測
- **目標 (Goal)**: 達成したい財務目標
- **緊急資金 (Emergency Fund)**: 緊急時に必要な資金
- **老後資金 (Retirement Fund)**: 退職後の生活に必要な資金
- **複利効果 (Compound Interest)**: 投資による資産増加効果
- **インフレ調整 (Inflation Adjustment)**: 物価上昇を考慮した実質価値

#### 境界づけられたコンテキスト (Bounded Context)

```
┌─────────────────────────────────────────────────────────────┐
│                    Financial Planning Context                │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Profile Mgmt  │  │   Goal Planning │  │   Calculation   │ │
│  │                 │  │                 │  │                 │ │
│  │ - FinancialData │  │ - Goal          │  │ - Projection    │ │
│  │ - Income/Expense│  │ - Progress      │  │ - Compound      │ │
│  │ - Savings       │  │ - Recommendation│  │ - Inflation     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │   Retirement    │  │   Emergency     │                   │
│  │                 │  │                 │                   │
│  │ - RetirementAge │  │ - EmergencyFund │                   │
│  │ - PensionCalc   │  │ - RiskAssess    │                   │
│  │ - LifeExpectancy│  │ - Liquidity     │                   │
│  └─────────────────┘  └─────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

### レイヤー構成

#### 1. Domain Layer (ドメイン層) - DDD適用

##### エンティティ (Entities)
```go
// domain/entities/
type FinancialProfile struct {
    id               ProfileID
    userID           UserID
    monthlyIncome    Money
    monthlyExpenses  ExpenseCollection
    currentSavings   SavingsCollection
    investmentReturn Rate
    inflationRate    Rate
    
    // ドメインメソッド
    func (fp *FinancialProfile) CalculateNetSavings() Money
    func (fp *FinancialProfile) ValidateFinancialHealth() error
    func (fp *FinancialProfile) ProjectAssets(years int) []AssetProjection
}

type Goal struct {
    id                  GoalID
    userID              UserID
    goalType            GoalType
    targetAmount        Money
    targetDate          Date
    monthlyContribution Money
    
    // ドメインメソッド
    func (g *Goal) CalculateProgress(currentAmount Money) ProgressRate
    func (g *Goal) EstimateCompletionDate(monthlySavings Money) Date
    func (g *Goal) IsAchievable(financialProfile FinancialProfile) bool
}
```

##### 値オブジェクト (Value Objects)
```go
// domain/valueobjects/
type Money struct {
    amount   float64
    currency Currency
    
    func (m Money) Add(other Money) Money
    func (m Money) Multiply(rate Rate) Money
    func (m Money) IsPositive() bool
}

type Rate struct {
    value float64
    
    func (r Rate) AsDecimal() float64
    func (r Rate) AsPercentage() float64
    func (r Rate) IsValid() bool
}

type Period struct {
    years  int
    months int
    
    func (p Period) ToMonths() int
    func (p Period) ToYears() float64
}
```

##### 集約 (Aggregates)
```go
// domain/aggregates/
type FinancialPlan struct {
    profile          FinancialProfile
    goals            []Goal
    retirementData   RetirementData
    emergencyFund    EmergencyFund
    
    // 集約ルート
    func (fp *FinancialPlan) AddGoal(goal Goal) error
    func (fp *FinancialPlan) UpdateProfile(profile FinancialProfile) error
    func (fp *FinancialPlan) GenerateProjection(years int) PlanProjection
    func (fp *FinancialPlan) ValidatePlan() []ValidationError
}
```

##### ドメインサービス (Domain Services)
```go
// domain/services/
type FinancialCalculationService struct {
    func CalculateCompoundInterest(principal Money, rate Rate, periods int) Money
    func CalculateInflationAdjustedValue(amount Money, inflationRate Rate, years int) Money
    func CalculateRetirementNeeds(monthlyExpenses Money, years int, inflationRate Rate) Money
}

type GoalRecommendationService struct {
    func RecommendMonthlySavings(goal Goal, currentSavings Money, timeRemaining Period) Money
    func SuggestGoalAdjustments(goal Goal, financialProfile FinancialProfile) []GoalAdjustment
}
```

##### リポジトリインターフェース (Repository Interfaces)
```go
// domain/repositories/
type FinancialPlanRepository interface {
    Save(ctx context.Context, plan FinancialPlan) error
    FindByUserID(ctx context.Context, userID UserID) (*FinancialPlan, error)
    Delete(ctx context.Context, planID PlanID) error
}

type GoalRepository interface {
    Save(ctx context.Context, goal Goal) error
    FindByUserID(ctx context.Context, userID UserID) ([]Goal, error)
    FindActiveGoals(ctx context.Context, userID UserID) ([]Goal, error)
}
```

#### 2. Use Cases (ユースケース層)
```go
// application/usecases/
type CalculateAssetProjectionUseCase interface {
    Execute(ctx context.Context, input AssetProjectionInput) (*AssetProjectionOutput, error)
}

type ManageFinancialDataUseCase interface {
    Create(ctx context.Context, data FinancialData) error
    Update(ctx context.Context, id EntityID, data FinancialData) error
    GetByUserID(ctx context.Context, userID UserID) (*FinancialData, error)
}
```

#### 3. Interface Adapters (インターフェースアダプター層)
```go
// infrastructure/controllers/
type FinancialController struct {
    useCase application.ManageFinancialDataUseCase
}

// infrastructure/repositories/
type PostgreSQLFinancialRepository struct {
    db *sql.DB
}

func (r *PostgreSQLFinancialRepository) Save(ctx context.Context, data domain.FinancialData) error
```

#### 4. Frameworks & Drivers (フレームワーク・ドライバー層)
```go
// infrastructure/web/
func NewEchoServer(controller *controllers.FinancialController) *echo.Echo

// infrastructure/database/
func NewPostgreSQLConnection(dsn string) (*sql.DB, error)
```

### 技術スタック

- **フロントエンド**: Next.js 14, TypeScript, Tailwind CSS, Chart.js
- **バックエンド**: Go, Echo Framework
- **データベース**: PostgreSQL
- **API仕様**: OpenAPI 3.0 (Swagger)
- **バリデーション**: Zod (フロント), Go validator (バック)
- **テスト**: Jest, React Testing Library (フロント), Go testing (バック)
- **ビルド**: Next.js build (フロント), Go build (バック)

## コンポーネントとインターフェース

### フロントエンドコンポーネント

#### 1. ページとレイアウト
- `app/layout.tsx` - ルートレイアウト
- `app/page.tsx` - ホームページ
- `app/dashboard/page.tsx` - ダッシュボード
- `components/Navigation.tsx` - ナビゲーションメニュー

#### 2. 入力フォームコンポーネント
- `FinancialInputForm.tsx` - 基本財務情報入力
- `GoalSettingForm.tsx` - 目標設定フォーム
- `RetirementForm.tsx` - 退職・年金情報入力

#### 3. 計算・表示コンポーネント
- `AssetProjectionChart.tsx` - 資産推移グラフ
- `RetirementCalculator.tsx` - 老後資金計算
- `EmergencyFundCalculator.tsx` - 緊急資金計算
- `ProgressTracker.tsx` - 目標進捗表示

#### 4. 共通コンポーネント
- `InputField.tsx` - 入力フィールド
- `Button.tsx` - ボタン
- `Modal.tsx` - モーダルダイアログ
- `LoadingSpinner.tsx` - ローディング表示

### バックエンドAPI

#### OpenAPI仕様

API設計はOpenAPI 3.0仕様で定義し、以下の利点を活用します：
- **型安全性**: フロントエンドでの型生成
- **ドキュメント自動生成**: Swagger UIでのAPI仕様確認
- **モックサーバー**: 開発初期段階でのフロントエンド開発
- **バリデーション**: リクエスト/レスポンスの自動検証

#### エンドポイント設計

```go
// ユーザー財務データ
POST   /api/financial-data     // 財務データ作成
GET    /api/financial-data     // 財務データ取得
PUT    /api/financial-data/:id // 財務データ更新
DELETE /api/financial-data/:id // 財務データ削除

// 計算API
POST   /api/calculations/asset-projection  // 資産推移計算
POST   /api/calculations/retirement        // 老後資金計算
POST   /api/calculations/emergency-fund    // 緊急資金計算

// 目標管理
POST   /api/goals           // 目標作成
GET    /api/goals           // 目標一覧取得
PUT    /api/goals/:id       // 目標更新
DELETE /api/goals/:id       // 目標削除

// レポート生成
GET    /api/reports/pdf     // PDFレポート生成

// API仕様
GET    /api/docs            // Swagger UI
GET    /api/openapi.json    // OpenAPI仕様ファイル
```

## データモデル

### 1. ユーザー財務データ (FinancialData)

```go
type FinancialData struct {
    ID               string        `json:"id" db:"id"`
    UserID           string        `json:"user_id" db:"user_id"`
    MonthlyIncome    float64       `json:"monthly_income" db:"monthly_income"`
    MonthlyExpenses  []ExpenseItem `json:"monthly_expenses" db:"monthly_expenses"`
    CurrentSavings   []SavingsItem `json:"current_savings" db:"current_savings"`
    InvestmentReturn float64       `json:"investment_return" db:"investment_return"`
    InflationRate    float64       `json:"inflation_rate" db:"inflation_rate"`
    CreatedAt        time.Time     `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
}

type ExpenseItem struct {
    Category    string  `json:"category"`
    Amount      float64 `json:"amount"`
    Description *string `json:"description,omitempty"`
}

type SavingsItem struct {
    Type        string  `json:"type" validate:"oneof=deposit investment other"`
    Amount      float64 `json:"amount"`
    Description *string `json:"description,omitempty"`
}
```

### 2. 退職・年金情報 (RetirementData)

```go
type RetirementData struct {
    ID                        string    `json:"id" db:"id"`
    UserID                    string    `json:"user_id" db:"user_id"`
    CurrentAge                int       `json:"current_age" db:"current_age"`
    RetirementAge             int       `json:"retirement_age" db:"retirement_age"`
    LifeExpectancy            int       `json:"life_expectancy" db:"life_expectancy"`
    MonthlyRetirementExpenses float64   `json:"monthly_retirement_expenses" db:"monthly_retirement_expenses"`
    PensionAmount             float64   `json:"pension_amount" db:"pension_amount"`
    CreatedAt                 time.Time `json:"created_at" db:"created_at"`
    UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}
```

### 3. 目標設定 (Goal)

```go
type Goal struct {
    ID                  string    `json:"id" db:"id"`
    UserID              string    `json:"user_id" db:"user_id"`
    Type                string    `json:"type" db:"type" validate:"oneof=savings retirement emergency custom"`
    Title               string    `json:"title" db:"title"`
    TargetAmount        float64   `json:"target_amount" db:"target_amount"`
    TargetDate          time.Time `json:"target_date" db:"target_date"`
    CurrentAmount       float64   `json:"current_amount" db:"current_amount"`
    MonthlyContribution float64   `json:"monthly_contribution" db:"monthly_contribution"`
    IsActive            bool      `json:"is_active" db:"is_active"`
    CreatedAt           time.Time `json:"created_at" db:"created_at"`
    UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
```

### 4. 計算結果 (CalculationResult)

```go
type AssetProjection struct {
    Year              int     `json:"year"`
    TotalAssets       float64 `json:"total_assets"`        // 総資産
    RealValue         float64 `json:"real_value"`          // 実質価値（インフレ調整後）
    ContributedAmount float64 `json:"contributed_amount"`  // 積立元本
    InvestmentGains   float64 `json:"investment_gains"`    // 投資収益
}

type RetirementCalculation struct {
    RequiredAmount            float64 `json:"required_amount"`             // 必要老後資金
    ProjectedAmount           float64 `json:"projected_amount"`            // 予想達成額
    Shortfall                 float64 `json:"shortfall"`                   // 不足額
    SufficiencyRate           float64 `json:"sufficiency_rate"`            // 充足率 (%)
    RecommendedMonthlySavings float64 `json:"recommended_monthly_savings"` // 推奨月間貯蓄額
}

type EmergencyFundCalculation struct {
    RequiredAmount  float64 `json:"required_amount"`   // 必要緊急資金
    CurrentAmount   float64 `json:"current_amount"`    // 現在の緊急資金
    Shortfall       float64 `json:"shortfall"`         // 不足額
    MonthsToTarget  int     `json:"months_to_target"`  // 目標達成までの月数
}
```

## エラーハンドリング

### バリデーション

```go
// Go validatorによる入力検証
type FinancialDataRequest struct {
    MonthlyIncome    float64       `json:"monthly_income" validate:"required,gt=0"`
    MonthlyExpenses  []ExpenseItem `json:"monthly_expenses" validate:"required,dive"`
    CurrentSavings   []SavingsItem `json:"current_savings" validate:"required,dive"`
    InvestmentReturn float64       `json:"investment_return" validate:"required,gte=0,lte=100"`
    InflationRate    float64       `json:"inflation_rate" validate:"required,gte=0,lte=50"`
}

type ExpenseItemRequest struct {
    Category string  `json:"category" validate:"required,min=1"`
    Amount   float64 `json:"amount" validate:"required,gt=0"`
}

type SavingsItemRequest struct {
    Type   string  `json:"type" validate:"required,oneof=deposit investment other"`
    Amount float64 `json:"amount" validate:"required,gte=0"`
}
```

### エラー処理戦略

1. **入力エラー**: フォームレベルでリアルタイム検証
2. **計算エラー**: 計算不可能な場合の代替表示
3. **ネットワークエラー**: 再試行機能とオフライン対応
4. **データエラー**: データ整合性チェックと修復提案

## テスト戦略

### 1. 単体テスト
- **計算ロジック**: 複利計算、インフレ調整、目標達成計算
- **バリデーション**: 入力値検証ロジック
- **ユーティリティ**: 日付計算、フォーマット関数

### 2. 統合テスト
- **API エンドポイント**: リクエスト/レスポンス検証
- **データベース操作**: CRUD操作の正確性
- **計算フロー**: 入力から結果表示までの一連の流れ

### 3. E2Eテスト
- **ユーザーシナリオ**: 財務データ入力から結果確認まで
- **グラフ表示**: チャートの正確な描画
- **レスポンシブ**: モバイル・デスクトップでの動作

### テストデータ

```go
// テスト用サンプルデータ
var mockFinancialData = FinancialData{
    MonthlyIncome: 400000,
    MonthlyExpenses: []ExpenseItem{
        {Category: "住居費", Amount: 120000},
        {Category: "食費", Amount: 60000},
        {Category: "交通費", Amount: 20000},
        {Category: "その他", Amount: 80000},
    },
    CurrentSavings: []SavingsItem{
        {Type: "deposit", Amount: 1000000},
        {Type: "investment", Amount: 500000},
    },
    InvestmentReturn: 5.0,
    InflationRate:    2.0,
}
```

## パフォーマンス考慮事項

### フロントエンド最適化
- **コンポーネント最適化**: React.memo, useMemo, useCallback
- **チャート最適化**: 大量データポイントの間引き表示
- **レイジーローディング**: ルート単位でのコード分割

### バックエンド最適化
- **計算キャッシュ**: 同一パラメータの計算結果キャッシュ
- **データベース最適化**: インデックス設定、クエリ最適化
- **レスポンス圧縮**: gzip圧縮によるデータ転送量削減

### セキュリティ

1. **入力サニタイゼーション**: XSS攻撃防止
2. **CORS設定**: 適切なオリジン制限
3. **レート制限**: API呼び出し頻度制限
4. **データ暗号化**: 機密財務データの暗号化保存