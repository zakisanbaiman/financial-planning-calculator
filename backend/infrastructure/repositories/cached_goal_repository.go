package repositories

import (
	"context"
	"log/slog"

	"github.com/financial-planning-calculator/backend/domain/entities"
	domainrepos "github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/infrastructure/monitoring"
	redisinfra "github.com/financial-planning-calculator/backend/infrastructure/redis"
)

const cacheTypeGoal = "goal"

// CachedGoalRepository は Cache-Aside パターンで GoalRepository をラップするデコレータ
type CachedGoalRepository struct {
	delegate    domainrepos.GoalRepository
	redisClient redisinfra.CacheClient
}

// NewCachedGoalRepository は新しいキャッシュデコレータを作成する
func NewCachedGoalRepository(
	delegate domainrepos.GoalRepository,
	redisClient redisinfra.CacheClient,
) domainrepos.GoalRepository {
	return &CachedGoalRepository{
		delegate:    delegate,
		redisClient: redisClient,
	}
}

// FindByUserID はキャッシュを確認し、なければDBから取得してキャッシュに保存する
func (r *CachedGoalRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
	key := goalsByUserIDKey(string(userID))

	var dtos []goalCacheDTO
	if err := r.redisClient.GetJSON(ctx, key, &dtos); err == nil {
		goals, err := goalsFromDTOs(dtos)
		if err != nil {
			slog.Warn("ゴールキャッシュのデシリアライズに失敗しました", slog.String("key", key), slog.Any("error", err))
		} else {
			monitoring.RecordCacheHit(cacheTypeGoal)
			return goals, nil
		}
	} else if !redisinfra.IsNil(err) {
		slog.Warn("Redis取得エラー（FindByUserID）、DBにフォールバック", slog.String("key", key), slog.Any("error", err))
	}

	monitoring.RecordCacheMiss(cacheTypeGoal)

	goals, err := r.delegate.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	r.setGoalsCache(ctx, key, goals)
	return goals, nil
}

// FindActiveGoalsByUserID はキャッシュを確認し、なければDBから取得してキャッシュに保存する
func (r *CachedGoalRepository) FindActiveGoalsByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
	key := activeGoalsByUserIDKey(string(userID))

	var dtos []goalCacheDTO
	if err := r.redisClient.GetJSON(ctx, key, &dtos); err == nil {
		goals, err := goalsFromDTOs(dtos)
		if err != nil {
			slog.Warn("アクティブゴールキャッシュのデシリアライズに失敗しました", slog.String("key", key), slog.Any("error", err))
		} else {
			monitoring.RecordCacheHit(cacheTypeGoal)
			return goals, nil
		}
	} else if !redisinfra.IsNil(err) {
		slog.Warn("Redis取得エラー（FindActiveGoalsByUserID）、DBにフォールバック", slog.String("key", key), slog.Any("error", err))
	}

	monitoring.RecordCacheMiss(cacheTypeGoal)

	goals, err := r.delegate.FindActiveGoalsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	r.setGoalsCache(ctx, key, goals)
	return goals, nil
}

// FindByID は委譲するだけ（個別取得はキャッシュ対象外）
func (r *CachedGoalRepository) FindByID(ctx context.Context, id entities.GoalID) (*entities.Goal, error) {
	return r.delegate.FindByID(ctx, id)
}

// FindByUserIDAndType は委譲するだけ（型フィルタは組み合わせ爆発のためキャッシュ対象外）
func (r *CachedGoalRepository) FindByUserIDAndType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) ([]*entities.Goal, error) {
	return r.delegate.FindByUserIDAndType(ctx, userID, goalType)
}

// Save は委譲後にユーザー単位のキャッシュを無効化する
func (r *CachedGoalRepository) Save(ctx context.Context, goal *entities.Goal) error {
	if err := r.delegate.Save(ctx, goal); err != nil {
		return err
	}
	r.invalidateUserCache(ctx, goal.UserID())
	return nil
}

// Update は委譲後にユーザー単位のキャッシュを無効化する
func (r *CachedGoalRepository) Update(ctx context.Context, goal *entities.Goal) error {
	if err := r.delegate.Update(ctx, goal); err != nil {
		return err
	}
	r.invalidateUserCache(ctx, goal.UserID())
	return nil
}

// Delete は委譲するだけ（GoalIDからUserIDが取れないため、無効化はしない）
// Note: ゴールのキャッシュTTLが短い（3分）ため、Deleteによる古いキャッシュは許容する
func (r *CachedGoalRepository) Delete(ctx context.Context, id entities.GoalID) error {
	return r.delegate.Delete(ctx, id)
}

// Exists は委譲するだけ
func (r *CachedGoalRepository) Exists(ctx context.Context, id entities.GoalID) (bool, error) {
	return r.delegate.Exists(ctx, id)
}

// CountActiveGoalsByType は委譲するだけ
func (r *CachedGoalRepository) CountActiveGoalsByType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) (int, error) {
	return r.delegate.CountActiveGoalsByType(ctx, userID, goalType)
}

// setGoalsCache はキャッシュへの書き込みを行う（失敗はログのみ）
func (r *CachedGoalRepository) setGoalsCache(ctx context.Context, key string, goals []*entities.Goal) {
	dtos := goalsToDTOs(goals)
	if err := r.redisClient.SetJSON(ctx, key, dtos, GoalTTL); err != nil {
		slog.Warn("ゴールキャッシュへの書き込みに失敗しました", slog.String("key", key), slog.Any("error", err))
	}
}

// invalidateUserCache はユーザー単位のゴールキャッシュキーをすべて削除する
func (r *CachedGoalRepository) invalidateUserCache(ctx context.Context, userID entities.UserID) {
	keys := []string{
		goalsByUserIDKey(string(userID)),
		activeGoalsByUserIDKey(string(userID)),
	}
	if err := r.redisClient.Delete(ctx, keys...); err != nil {
		slog.Warn("ゴールキャッシュの無効化に失敗しました", slog.Any("error", err))
	}
}
