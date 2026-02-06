# クラス図 (Class Diagram)

このドキュメントは、財務計画計算機アプリケーションのドメインモデルとアーキテクチャ構造を可視化したクラス図です。

## Clean Architecture / DDD構造

このアプリケーションは、Clean Architecture（クリーンアーキテクチャ）とDDD（ドメイン駆動設計）の原則に基づいて設計されています。

## レイヤー構造

```mermaid
graph TB
    subgraph "Presentation Layer（プレゼンテーション層）"
        Controllers["Controllers<br/>HTTPハンドラー"]
        Middleware["Middleware<br/>認証・検証"]
    end

    subgraph "Application Layer（アプリケーション層）"
        UseCases["UseCases<br/>ビジネスロジック"]
    end

    subgraph "Domain Layer（ドメイン層）"
        Entities["Entities<br/>エンティティ"]
        ValueObjects["Value Objects<br/>値オブジェクト"]
        Aggregates["Aggregates<br/>集約"]
        Services["Domain Services<br/>ドメインサービス"]
        Repositories["Repository Interfaces<br/>リポジトリインターフェース"]
    end

    subgraph "Infrastructure Layer（インフラ層）"
        RepositoryImpl["Repository Implementations<br/>リポジトリ実装"]
        Database["Database<br/>PostgreSQL"]
        ExternalAPI["External APIs<br/>OAuth等"]
    end

    Controllers --> UseCases
    Middleware --> Controllers
    UseCases --> Services
    UseCases --> Repositories
    Services --> Entities
    Services --> ValueObjects
    Services --> Aggregates
    Repositories --> Entities
    RepositoryImpl -.implements.-> Repositories
    RepositoryImpl --> Database
    Controllers --> ExternalAPI
```

## ドメインモデル

### 集約（Aggregates）

```mermaid
classDiagram
    class FinancialPlan {
        -FinancialPlanID id
        -FinancialProfile profile
        -Goal[] goals
        -RetirementData retirementData
        -EmergencyFundConfig emergencyFund
        -Time createdAt
        -Time updatedAt
        +NewFinancialPlan(profile) FinancialPlan
        +AddGoal(goal) error
        +RemoveGoal(goalID) error
        +UpdateProfile(profile) error
        +SetRetirementData(data) error
        +GenerateProjection(years) PlanProjection
        +ValidatePlan() ValidationError[]
    }

    class PlanProjection {
        +AssetProjection[] assetProjections
        +RetirementCalculation retirementCalculation
        +EmergencyFundStatus emergencyFundStatus
        +GoalProgress[] goalProgress
    }

    class EmergencyFundConfig {
        +int targetMonths
        +Money currentFund
        +NewEmergencyFundConfig(months, fund) EmergencyFundConfig
    }

    FinancialPlan --> PlanProjection : generates
    FinancialPlan --> EmergencyFundConfig : contains
```

### エンティティ（Entities）

```mermaid
classDiagram
    class User {
        -UserID id
        -Email email
        -PasswordHash passwordHash
        -Provider provider
        -ProviderUserID providerUserID
        -Name name
        -AvatarURL avatarURL
        -TwoFactorEnabled twoFactorEnabled
        -TwoFactorSecret twoFactorSecret
        -TwoFactorBackupCodes backupCodes
        -EmailVerified emailVerified
        -Time createdAt
        -Time updatedAt
        +NewUser(email, password) User
        +EnableTwoFactor(secret) error
        +VerifyEmail() error
    }

    class FinancialProfile {
        -FinancialProfileID id
        -UserID userID
        -Money monthlyIncome
        -Rate investmentReturn
        -Rate inflationRate
        -ExpenseItems monthlyExpenses
        -SavingsItems currentSavings
        +NewFinancialProfile(income, return, inflation) FinancialProfile
        +CalculateNetSavings() Money
        +ProjectAssets(years) AssetProjection[]
        +ValidateFinancialHealth() error
    }

    class Goal {
        -GoalID id
        -UserID userID
        -GoalType type
        -string title
        -Money targetAmount
        -Date targetDate
        -Money currentAmount
        -Money monthlyContribution
        -bool isActive
        -Time createdAt
        -Time updatedAt
        +NewGoal(type, title, amount, date) Goal
        +UpdateProgress(amount) error
        +CalculateProgress(current) ProgressRate
        +IsAchievable(profile) bool
        +IsCompleted() bool
        +IsOverdue() bool
    }

    class RetirementData {
        -RetirementDataID id
        -UserID userID
        -Age currentAge
        -Age retirementAge
        -Age lifeExpectancy
        -Money monthlyRetirementExpenses
        -Money pensionAmount
        -Time createdAt
        -Time updatedAt
        +NewRetirementData(ages, expenses, pension) RetirementData
        +CalculateRetirementSufficiency(savings, netSavings, return, inflation) RetirementCalculation
        +CalculateRetirementYears() int
    }

    class RefreshToken {
        -RefreshTokenID id
        -UserID userID
        -Token token
        -Time expiresAt
        -Time createdAt
        +NewRefreshToken(userID, token, expiry) RefreshToken
        +IsExpired() bool
    }

    class WebAuthnCredential {
        -string id
        -UserID userID
        -CredentialID credentialID
        -PublicKey publicKey
        -AttestationType attestationType
        -AAGUID aaguid
        -SignCount signCount
        -bool cloneWarning
        -Transports transports
        -string name
        -Time createdAt
        -Time updatedAt
        -Time lastUsedAt
        +NewWebAuthnCredential(userID, credID, pubKey) WebAuthnCredential
        +UpdateSignCount(count) error
    }

    class AssetProjection {
        +int year
        +Money totalAssets
        +Money investmentGains
        +Money contributions
    }

    class RetirementCalculation {
        +Money requiredAmount
        +Money projectedAmount
        +Money surplus
        +bool isSufficient
    }

    User --> FinancialProfile : has
    User --> Goal : has
    User --> RetirementData : has
    User --> RefreshToken : has
    User --> WebAuthnCredential : has
    FinancialProfile --> AssetProjection : projects
    RetirementData --> RetirementCalculation : calculates
```

### 値オブジェクト（Value Objects）

```mermaid
classDiagram
    class Money {
        -float64 amount
        -Currency currency
        +NewMoneyJPY(amount) Money
        +Add(other) Money
        +Subtract(other) Money
        +MultiplyByFloat(factor) Money
        +IsPositive() bool
        +IsNegative() bool
        +Amount() float64
    }

    class Rate {
        -float64 percentage
        +NewRate(percentage) Rate
        +AsDecimal() float64
        +AsPercentage() float64
        +Validate() error
    }

    class Period {
        -int months
        +NewPeriod(years, months) Period
        +InYears() int
        +InMonths() int
        +AddMonths(months) Period
    }

    class Age {
        -int years
        +NewAge(years) Age
        +Years() int
        +IsValid() bool
    }

    class Email {
        -string address
        +NewEmail(address) Email
        +String() string
        +IsValid() bool
    }

    class GoalType {
        <<enumeration>>
        Savings
        Retirement
        Emergency
        Custom
        +String() string
    }

    class ProgressRate {
        -float64 percentage
        +NewProgressRate(current, target) ProgressRate
        +AsPercentage() float64
        +IsComplete() bool
    }
```

### ドメインサービス（Domain Services）

```mermaid
classDiagram
    class FinancialCalculationService {
        +CalculateCompoundInterest(principal, rate, years) Money
        +CalculateFutureValue(presentValue, rate, years) Money
        +CalculateMonthlyPayment(principal, rate, months) Money
        +CalculateNetWorth(assets, liabilities) Money
        +ProjectAssetGrowth(initial, contribution, return, years) AssetProjection[]
    }

    class GoalRecommendationService {
        +RecommendEmergencyFund(expenses) Money
        +RecommendRetirementSavings(age, retirement, expenses) Money
        +SuggestMonthlyContribution(goal, profile) Money
        +EvaluateGoalFeasibility(goal, profile) bool
    }

    FinancialCalculationService ..> Money : uses
    FinancialCalculationService ..> Rate : uses
    GoalRecommendationService ..> Goal : uses
    GoalRecommendationService ..> FinancialProfile : uses
```

## ユースケース（Application Layer）

```mermaid
classDiagram
    class AuthUseCase {
        -UserRepository userRepo
        -RefreshTokenRepository tokenRepo
        -JWTService jwtService
        +Register(email, password) User, error
        +Login(email, password) AuthResponse, error
        +RefreshToken(token) AuthResponse, error
        +VerifyEmail(token) error
        +EnableTwoFactor(userID) TwoFactorSetup, error
    }

    class WebAuthnUseCase {
        -UserRepository userRepo
        -WebAuthnCredentialRepository credRepo
        -WebAuthn webauthn
        +BeginRegistration(userID) CredentialCreation, error
        +FinishRegistration(userID, response) error
        +BeginLogin(email) CredentialAssertion, error
        +FinishLogin(email, response) AuthResponse, error
    }

    class ManageFinancialDataUseCase {
        -FinancialPlanRepository planRepo
        +CreateFinancialProfile(userID, data) FinancialProfile, error
        +UpdateFinancialProfile(userID, data) FinancialProfile, error
        +GetFinancialProfile(userID) FinancialProfile, error
        +AddExpenseItem(profileID, item) error
        +AddSavingsItem(profileID, item) error
    }

    class ManageGoalsUseCase {
        -GoalRepository goalRepo
        -FinancialPlanRepository planRepo
        +CreateGoal(userID, goal) Goal, error
        +UpdateGoal(goalID, data) Goal, error
        +DeleteGoal(goalID) error
        +ListGoals(userID, filter) Goal[], error
        +UpdateProgress(goalID, amount) error
    }

    class CalculateProjectionUseCase {
        -FinancialPlanRepository planRepo
        -FinancialCalculationService calcService
        +GenerateProjection(userID, years) PlanProjection, error
        +CalculateRetirementSufficiency(userID) RetirementCalculation, error
        +EvaluateEmergencyFund(userID) EmergencyFundStatus, error
    }

    class GenerateReportsUseCase {
        -FinancialPlanRepository planRepo
        -PDFGenerator pdfGen
        +GenerateFinancialReport(userID) Report, error
        +ExportToPDF(userID) []byte, error
    }

    AuthUseCase --> UserRepository
    WebAuthnUseCase --> UserRepository
    ManageFinancialDataUseCase --> FinancialPlanRepository
    ManageGoalsUseCase --> GoalRepository
    CalculateProjectionUseCase --> FinancialPlanRepository
    CalculateProjectionUseCase --> FinancialCalculationService
    GenerateReportsUseCase --> FinancialPlanRepository
```

## リポジトリ（Repository Interfaces）

```mermaid
classDiagram
    class UserRepository {
        <<interface>>
        +FindByID(id) User, error
        +FindByEmail(email) User, error
        +FindByProviderUserID(provider, id) User, error
        +Create(user) error
        +Update(user) error
        +Delete(id) error
    }

    class FinancialPlanRepository {
        <<interface>>
        +FindByUserID(userID) FinancialPlan, error
        +Create(plan) error
        +Update(plan) error
        +Delete(id) error
    }

    class GoalRepository {
        <<interface>>
        +FindByID(id) Goal, error
        +FindByUserID(userID) Goal[], error
        +FindActiveGoals(userID) Goal[], error
        +Create(goal) error
        +Update(goal) error
        +Delete(id) error
    }

    class RefreshTokenRepository {
        <<interface>>
        +FindByToken(token) RefreshToken, error
        +FindByUserID(userID) RefreshToken[], error
        +Create(token) error
        +Delete(token) error
        +DeleteExpired() error
    }

    class WebAuthnCredentialRepository {
        <<interface>>
        +FindByID(id) WebAuthnCredential, error
        +FindByUserID(userID) WebAuthnCredential[], error
        +FindByCredentialID(credID) WebAuthnCredential, error
        +Create(credential) error
        +Update(credential) error
        +Delete(id) error
    }
```

## アーキテクチャの特徴

### 依存性の方向
- 外側の層は内側の層に依存する
- 内側の層（Domain Layer）は外側の層に依存しない
- Infrastructure層はDomain層のインターフェースを実装する

### レイヤーの責務

#### Domain Layer（ドメイン層）
- ビジネスロジックの核心
- エンティティ、値オブジェクト、集約、ドメインサービス
- フレームワークやライブラリに依存しない

#### Application Layer（アプリケーション層）
- ユースケースの実装
- ドメインオブジェクトの調整
- トランザクション管理

#### Infrastructure Layer（インフラ層）
- データベースアクセス
- 外部API連携
- ファイルシステムアクセス

#### Presentation Layer（プレゼンテーション層）
- HTTPリクエスト/レスポンス処理
- 入力バリデーション
- 認証・認可

### 設計パターン

- **Repository Pattern**: データアクセスの抽象化
- **Aggregate Pattern**: 関連エンティティのグループ化
- **Value Object Pattern**: 不変な値の表現
- **Domain Service Pattern**: エンティティに属さないビジネスロジック
- **Use Case Pattern**: アプリケーション機能の明示的な表現
