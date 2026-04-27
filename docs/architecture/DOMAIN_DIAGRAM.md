# ドメイン境界の依存関係図

このドキュメントは、バックエンドの DDD 実装におけるドメイン境界・依存関係を可視化し、各層の設計意図を明確にしたものです。

## 全体依存関係マップ

### レイヤー間の依存方向

```mermaid
graph TD
    subgraph Presentation["🖥️ Presentation Layer（プレゼンテーション層）"]
        Controllers["Controllers<br/>(HTTPハンドラー)"]
        Middleware["Middleware<br/>(認証・検証)"]
    end

    subgraph Application["⚙️ Application Layer（アプリケーション層）"]
        UseCases["UseCases<br/>(ユースケース)"]
        Ports["Ports<br/>(外部サービスインターフェース)"]
    end

    subgraph Domain["🏛️ Domain Layer（ドメイン層）"]
        Aggregates["Aggregates<br/>(集約)"]
        Entities["Entities<br/>(エンティティ)"]
        ValueObjects["Value Objects<br/>(値オブジェクト)"]
        DomainServices["Domain Services<br/>(ドメインサービス)"]
        RepoInterfaces["Repository Interfaces<br/>(リポジトリIF)"]
    end

    subgraph Infrastructure["🔧 Infrastructure Layer（インフラ層）"]
        RepoPG["PostgreSQL実装<br/>(リポジトリ)"]
        RepoCache["Cacheデコレータ<br/>(リポジトリ)"]
        Redis["Redis<br/>(キャッシュ)"]
        DB["PostgreSQL<br/>(データベース)"]
        ExternalAPI["External APIs<br/>(OAuth, LLM等)"]
    end

    Controllers --> UseCases
    Middleware --> Controllers
    UseCases --> Aggregates
    UseCases --> Entities
    UseCases --> RepoInterfaces
    UseCases --> DomainServices
    UseCases --> Ports
    DomainServices --> Entities
    DomainServices --> ValueObjects
    Aggregates --> Entities
    Aggregates --> ValueObjects
    RepoInterfaces --> Entities
    RepoInterfaces --> Aggregates
    RepoPG -.->|implements| RepoInterfaces
    RepoCache -.->|implements| RepoInterfaces
    RepoCache --> RepoPG
    RepoCache --> Redis
    RepoPG --> DB
    Ports -.->|adapters in infra| ExternalAPI
```

### ドメイン層の内部依存

ドメイン層内でのサブパッケージ間の依存は一方向に保たれています。

```mermaid
graph LR
    subgraph Domain["Domain Layer"]
        direction TB
        ValueObjects["valueobjects<br/>Money / Rate / Period / Age / Email<br/>GoalType / ProgressRate"]
        Entities["entities<br/>User / FinancialProfile / Goal<br/>RetirementData / RefreshToken<br/>WebAuthnCredential"]
        Aggregates["aggregates<br/>FinancialPlan"]
        DomainServices["services<br/>FinancialCalculationService<br/>GoalRecommendationService"]
        RepoInterfaces["repositories<br/>UserRepository IF<br/>FinancialPlanRepository IF<br/>GoalRepository IF<br/>RefreshTokenRepository IF<br/>WebAuthnCredentialRepository IF"]
    end

    ValueObjects --> |no dependencies|ValueObjects
    Entities --> ValueObjects
    Aggregates --> Entities
    Aggregates --> ValueObjects
    DomainServices --> Entities
    DomainServices --> ValueObjects
    RepoInterfaces --> Entities
    RepoInterfaces --> Aggregates
```

**ルール**: 値オブジェクト → エンティティ → 集約 の順で依存し、逆方向の依存は禁止。

---

## ドメインモデルの境界

### 1. 財務計画ドメイン（Financial Planning Domain）

財務計画に関する中心的なビジネスロジックを担うドメイン。

```mermaid
graph LR
    subgraph FinancialPlanAggregate["FinancialPlan 集約"]
        FP["FinancialPlan<br/>(集約ルート)"]
        FProf["FinancialProfile<br/>(エンティティ)"]
        Goal["Goal<br/>(エンティティ ×複数)"]
        RD["RetirementData<br/>(エンティティ)"]
        EFC["EmergencyFundConfig<br/>(値オブジェクト)"]
        FP --> FProf
        FP --> Goal
        FP --> RD
        FP --> EFC
    end

    subgraph ValueObjects["値オブジェクト"]
        Money["Money<br/>(金額 + 通貨)"]
        Rate["Rate<br/>(利率)"]
        Period["Period<br/>(期間)"]
        Age["Age<br/>(年齢)"]
        GoalType["GoalType<br/>(目標種別)"]
        ProgressRate["ProgressRate<br/>(達成率)"]
    end

    FProf --> Money
    FProf --> Rate
    Goal --> Money
    Goal --> GoalType
    Goal --> ProgressRate
    RD --> Money
    RD --> Age
    RD --> Period
```

**集約の不変条件（Invariants）**:
- `FinancialPlan.AddGoal()` は目標の達成可能性を検証してから追加する
- 目標額・期日の整合性は `Goal` エンティティ内で保証する
- 集約外からは `Goal` を直接変更できない（集約ルート経由のみ）

### 2. 認証ドメイン（Authentication Domain）

ユーザー識別・認証に関するドメイン。財務計画ドメインとは独立した集約として設計。

```mermaid
graph LR
    subgraph User["User 集約"]
        U["User<br/>(エンティティ / 集約ルート)"]
        Email["Email<br/>(値オブジェクト)"]
        U --> Email
    end

    subgraph AuthTokens["認証トークン集約"]
        RT["RefreshToken<br/>(エンティティ)"]
        PRT["PasswordResetToken<br/>(エンティティ)"]
    end

    subgraph WebAuthn["WebAuthn集約"]
        WAC["WebAuthnCredential<br/>(エンティティ)"]
    end

    U -.->|userID参照| RT
    U -.->|userID参照| PRT
    U -.->|userID参照| WAC
```

**境界の意図**: 認証トークン・WebAuthn資格情報は、ライフサイクルが `User` と異なる（トークンは短命、資格情報は複数保持）ため、別の集約として分離。集約間は `UserID` 参照で結合する。

---

## Application Layer のユースケース依存

各ユースケースがどのドメインオブジェクトを利用するかを示す。

```mermaid
graph TD
    subgraph AuthGroup["認証グループ"]
        AuthUC["AuthUseCase"]
        WebAuthnUC["WebAuthnUseCase"]
    end

    subgraph FinancialGroup["財務計画グループ"]
        ManageFinancialUC["ManageFinancialDataUseCase"]
        ManageGoalsUC["ManageGoalsUseCase"]
        CalcProjectionUC["CalculateProjectionUseCase"]
        GenerateReportsUC["GenerateReportsUseCase"]
    end

    subgraph BotGroup["AIボットグループ"]
        BotUC["BotUseCase"]
    end

    AuthUC --> UserRepo["UserRepository"]
    AuthUC --> RTRepo["RefreshTokenRepository"]
    WebAuthnUC --> UserRepo
    WebAuthnUC --> WACRepo["WebAuthnCredentialRepository"]

    ManageFinancialUC --> FPRepo["FinancialPlanRepository"]
    ManageGoalsUC --> GoalRepo["GoalRepository"]
    ManageGoalsUC --> FPRepo
    CalcProjectionUC --> FPRepo
    CalcProjectionUC --> CalcSvc["FinancialCalculationService"]
    GenerateReportsUC --> FPRepo

    BotUC --> LLMPort["LLMClient Port"]
    BotUC --> FAQPort["FAQLoader Port"]
```

**設計の意図**: ユースケースはドメインサービスとリポジトリインターフェースにのみ依存する。インフラ実装（PostgreSQL、Redis、外部API）には依存しない。

---

## Infrastructure Layer の実装構造

```mermaid
graph TD
    subgraph RepositoryIF["Repository Interfaces（Domain Layer）"]
        FPRepoIF["FinancialPlanRepository"]
        GoalRepoIF["GoalRepository"]
    end

    subgraph CacheDecorators["Cache Decorators（Infrastructure Layer）"]
        CachedFP["CachedFinancialPlanRepository<br/>TTL: 5分"]
        CachedGoal["CachedGoalRepository<br/>TTL: 3分"]
    end

    subgraph PGImpl["PostgreSQL Implementations（Infrastructure Layer）"]
        PGFP["PostgreSQLFinancialPlanRepository"]
        PGGoal["PostgreSQLGoalRepository"]
    end

    CachedFP -.->|implements| FPRepoIF
    CachedGoal -.->|implements| GoalRepoIF
    PGFP -.->|implements| FPRepoIF
    PGGoal -.->|implements| GoalRepoIF

    CachedFP --> PGFP
    CachedGoal --> PGGoal

    CachedFP --> Redis["Redis\n(Cache-Aside)"]
    CachedGoal --> Redis

    PGFP --> PG["PostgreSQL"]
    PGGoal --> PG
```

**Decorator Pattern**: `CachedXxxRepository` は `XxxRepository` インターフェースを実装しながら、内部で PostgreSQL 実装をラップする。アプリケーション層はキャッシュ有無を意識しない（→ ADR-006）。

---

## 境界を越えるデータフロー

典型的なユースケース「将来予測の計算」でのデータフロー。

```mermaid
sequenceDiagram
    participant C as Controller
    participant UC as CalculateProjectionUseCase
    participant Repo as FinancialPlanRepository
    participant Cache as CachedRepository
    participant PG as PostgreSQL
    participant Svc as FinancialCalculationService
    participant Agg as FinancialPlan集約

    C->>UC: GenerateProjection(userID, years)
    UC->>Repo: FindByUserID(userID)
    Repo->>Cache: FindByUserID(userID)
    alt キャッシュヒット
        Cache-->>Repo: FinancialPlan (from Redis)
    else キャッシュミス
        Cache->>PG: FindByUserID(userID)
        PG-->>Cache: FinancialPlan
        Cache-->>Repo: FinancialPlan (saved to Redis)
    end
    Repo-->>UC: FinancialPlan
    UC->>Agg: GenerateProjection(years)
    Agg->>Svc: CalculateCompoundInterest(...)
    Svc-->>Agg: Money
    Agg-->>UC: PlanProjection
    UC-->>C: PlanProjection
```

---

## 境界の設計原則まとめ

| 原則 | 内容 | 実現方法 |
|---|---|---|
| **依存方向の統一** | 外側→内側のみ。内側は外側を知らない | Goのimportで強制 |
| **集約の自己完結性** | 集約の不変条件は集約ルートが保証する | `AddGoal()`等のメソッドで検証 |
| **集約間の結合はIDのみ** | 異なる集約はオブジェクト参照ではなくIDで結合 | `UserID`型で参照 |
| **インターフェース経由の抽象化** | インフラ実装をドメインから隠蔽 | Repository Interfaceパターン |
| **Ports & Adapters** | 外部サービス（LLM, FAQ）はポートで抽象化 | `application/ports/`パッケージ |
| **値オブジェクトの不変性** | `Money`、`Rate`等は変更不可 | Goの値型（structの値渡し） |

---

**関連ドキュメント:**
- [クラス図（各モデルの属性・メソッド詳細）](./CLASS_DIAGRAM.md)
- [ADR-007: DDDドメイン境界の設計方針](../adr/007-ddd-domain-boundary.md)
- [ADR-006: Redisキャッシュ戦略](../adr/006-redis-cache-strategy.md)

**最終更新日**: 2026-04-27
