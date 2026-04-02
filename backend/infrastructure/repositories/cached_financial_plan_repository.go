package repositories

import (
	"context"
	"log/slog"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	domainrepos "github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/infrastructure/monitoring"
	redisinfra "github.com/financial-planning-calculator/backend/infrastructure/redis"
)

const cacheTypeFinancialPlan = "financial_plan"

// CachedFinancialPlanRepository は Cache-Aside パターンで FinancialPlanRepository をラップするデコレータ
type CachedFinancialPlanRepository struct {
	delegate    domainrepos.FinancialPlanRepository
	redisClient redisinfra.CacheClient
}

// NewCachedFinancialPlanRepository は新しいキャッシュデコレータを作成する
func NewCachedFinancialPlanRepository(
	delegate domainrepos.FinancialPlanRepository,
	redisClient redisinfra.CacheClient,
) domainrepos.FinancialPlanRepository {
	return &CachedFinancialPlanRepository{
		delegate:    delegate,
		redisClient: redisClient,
	}
}

// FindByID はキャッシュを確認し、なければDBから取得してキャッシュに保存する
func (r *CachedFinancialPlanRepository) FindByID(ctx context.Context, id aggregates.FinancialPlanID) (*aggregates.FinancialPlan, error) {
	key := financialPlanByIDKey(string(id))

	var dto financialPlanCacheDTO
	if err := r.redisClient.GetJSON(ctx, key, &dto); err == nil {
		plan, err := financialPlanFromDTO(dto)
		if err != nil {
			// デシリアライズ失敗はキャッシュミスとして扱う
			slog.Warn("財務計画キャッシュのデシリアライズに失敗しました", slog.String("key", key), slog.Any("error", err))
		} else {
			monitoring.RecordCacheHit(cacheTypeFinancialPlan)
			return plan, nil
		}
	} else if !redisinfra.IsNil(err) {
		// redis.Nil 以外のエラーはRedis障害 → fail-open
		slog.Warn("Redis取得エラー（FindByID）、DBにフォールバック", slog.String("key", key), slog.Any("error", err))
	}

	monitoring.RecordCacheMiss(cacheTypeFinancialPlan)

	plan, err := r.delegate.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	r.setCache(ctx, financialPlanByIDKey(string(plan.ID())), plan)
	return plan, nil
}

// FindByUserID はキャッシュを確認し、なければDBから取得してキャッシュに保存する
func (r *CachedFinancialPlanRepository) FindByUserID(ctx context.Context, userID entities.UserID) (*aggregates.FinancialPlan, error) {
	key := financialPlanByUserIDKey(string(userID))

	var dto financialPlanCacheDTO
	if err := r.redisClient.GetJSON(ctx, key, &dto); err == nil {
		plan, err := financialPlanFromDTO(dto)
		if err != nil {
			slog.Warn("財務計画キャッシュのデシリアライズに失敗しました", slog.String("key", key), slog.Any("error", err))
		} else {
			monitoring.RecordCacheHit(cacheTypeFinancialPlan)
			return plan, nil
		}
	} else if !redisinfra.IsNil(err) {
		slog.Warn("Redis取得エラー（FindByUserID）、DBにフォールバック", slog.String("key", key), slog.Any("error", err))
	}

	monitoring.RecordCacheMiss(cacheTypeFinancialPlan)

	plan, err := r.delegate.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// FindByID キャッシュと FindByUserID キャッシュの両方を設定
	r.setCache(ctx, financialPlanByIDKey(string(plan.ID())), plan)
	r.setCache(ctx, key, plan)
	return plan, nil
}

// Save は委譲後にキャッシュを無効化する
func (r *CachedFinancialPlanRepository) Save(ctx context.Context, plan *aggregates.FinancialPlan) error {
	if err := r.delegate.Save(ctx, plan); err != nil {
		return err
	}
	r.invalidateCache(ctx, plan)
	return nil
}

// Update は委譲後にキャッシュを無効化する
func (r *CachedFinancialPlanRepository) Update(ctx context.Context, plan *aggregates.FinancialPlan) error {
	if err := r.delegate.Update(ctx, plan); err != nil {
		return err
	}
	r.invalidateCache(ctx, plan)
	return nil
}

// Delete は委譲後にキャッシュを無効化する
func (r *CachedFinancialPlanRepository) Delete(ctx context.Context, id aggregates.FinancialPlanID) error {
	if err := r.delegate.Delete(ctx, id); err != nil {
		return err
	}
	// DeleteはIDのみ持つため、FindByIDキャッシュのみ無効化
	if err := r.redisClient.Delete(ctx, financialPlanByIDKey(string(id))); err != nil {
		slog.Warn("財務計画キャッシュの無効化に失敗しました", slog.String("key", financialPlanByIDKey(string(id))), slog.Any("error", err))
	}
	return nil
}

// Exists は委譲するだけ（存在チェックはキャッシュ対象外）
func (r *CachedFinancialPlanRepository) Exists(ctx context.Context, id aggregates.FinancialPlanID) (bool, error) {
	return r.delegate.Exists(ctx, id)
}

// ExistsByUserID は委譲するだけ
func (r *CachedFinancialPlanRepository) ExistsByUserID(ctx context.Context, userID entities.UserID) (bool, error) {
	return r.delegate.ExistsByUserID(ctx, userID)
}

// setCache はキャッシュへの書き込みを行う（失敗はログのみ）
func (r *CachedFinancialPlanRepository) setCache(ctx context.Context, key string, plan *aggregates.FinancialPlan) {
	dto := financialPlanToDTO(plan)
	if err := r.redisClient.SetJSON(ctx, key, dto, FinancialPlanTTL); err != nil {
		slog.Warn("財務計画キャッシュへの書き込みに失敗しました", slog.String("key", key), slog.Any("error", err))
	}
}

// invalidateCache は plan に関連するキャッシュキーをすべて削除する
func (r *CachedFinancialPlanRepository) invalidateCache(ctx context.Context, plan *aggregates.FinancialPlan) {
	keys := []string{
		financialPlanByIDKey(string(plan.ID())),
		financialPlanByUserIDKey(string(plan.Profile().UserID())),
	}
	if err := r.redisClient.Delete(ctx, keys...); err != nil {
		slog.Warn("財務計画キャッシュの無効化に失敗しました", slog.Any("error", err))
	}
}
