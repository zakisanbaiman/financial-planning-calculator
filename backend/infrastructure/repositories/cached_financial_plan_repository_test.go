package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	redisinfra "github.com/financial-planning-calculator/backend/infrastructure/redis"
	goredis "github.com/redis/go-redis/v9"
)

// --- モック: FinancialPlanRepository ---

type mockFinancialPlanRepository struct {
	findByIDFunc     func(ctx context.Context, id aggregates.FinancialPlanID) (*aggregates.FinancialPlan, error)
	findByUserIDFunc func(ctx context.Context, userID entities.UserID) (*aggregates.FinancialPlan, error)
	saveFunc         func(ctx context.Context, plan *aggregates.FinancialPlan) error
	updateFunc       func(ctx context.Context, plan *aggregates.FinancialPlan) error
	deleteFunc       func(ctx context.Context, id aggregates.FinancialPlanID) error
	existsFunc       func(ctx context.Context, id aggregates.FinancialPlanID) (bool, error)
	existsByUserFunc func(ctx context.Context, userID entities.UserID) (bool, error)
	callCount        map[string]int
}

func newMockFinancialPlanRepo() *mockFinancialPlanRepository {
	return &mockFinancialPlanRepository{callCount: make(map[string]int)}
}

func (m *mockFinancialPlanRepository) FindByID(ctx context.Context, id aggregates.FinancialPlanID) (*aggregates.FinancialPlan, error) {
	m.callCount["FindByID"]++
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFinancialPlanRepository) FindByUserID(ctx context.Context, userID entities.UserID) (*aggregates.FinancialPlan, error) {
	m.callCount["FindByUserID"]++
	if m.findByUserIDFunc != nil {
		return m.findByUserIDFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFinancialPlanRepository) Save(ctx context.Context, plan *aggregates.FinancialPlan) error {
	m.callCount["Save"]++
	if m.saveFunc != nil {
		return m.saveFunc(ctx, plan)
	}
	return nil
}

func (m *mockFinancialPlanRepository) Update(ctx context.Context, plan *aggregates.FinancialPlan) error {
	m.callCount["Update"]++
	if m.updateFunc != nil {
		return m.updateFunc(ctx, plan)
	}
	return nil
}

func (m *mockFinancialPlanRepository) Delete(ctx context.Context, id aggregates.FinancialPlanID) error {
	m.callCount["Delete"]++
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockFinancialPlanRepository) Exists(ctx context.Context, id aggregates.FinancialPlanID) (bool, error) {
	m.callCount["Exists"]++
	if m.existsFunc != nil {
		return m.existsFunc(ctx, id)
	}
	return false, nil
}

func (m *mockFinancialPlanRepository) ExistsByUserID(ctx context.Context, userID entities.UserID) (bool, error) {
	m.callCount["ExistsByUserID"]++
	if m.existsByUserFunc != nil {
		return m.existsByUserFunc(ctx, userID)
	}
	return false, nil
}

// --- モック: CacheClient ---

type mockCacheClient struct {
	getJSONFunc        func(ctx context.Context, key string, dest any) error
	setJSONFunc        func(ctx context.Context, key string, value any, ttl time.Duration) error
	deleteFunc         func(ctx context.Context, keys ...string) error
	deleteByPatternFunc func(ctx context.Context, pattern string) error
	callCount          map[string]int
}

func newMockCacheClient() *mockCacheClient {
	return &mockCacheClient{callCount: make(map[string]int)}
}

func (m *mockCacheClient) GetJSON(ctx context.Context, key string, dest any) error {
	m.callCount["GetJSON"]++
	if m.getJSONFunc != nil {
		return m.getJSONFunc(ctx, key, dest)
	}
	return goredis.Nil // デフォルトはキャッシュミス
}

func (m *mockCacheClient) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	m.callCount["SetJSON"]++
	if m.setJSONFunc != nil {
		return m.setJSONFunc(ctx, key, value, ttl)
	}
	return nil
}

func (m *mockCacheClient) Delete(ctx context.Context, keys ...string) error {
	m.callCount["Delete"]++
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, keys...)
	}
	return nil
}

func (m *mockCacheClient) DeleteByPattern(ctx context.Context, pattern string) error {
	m.callCount["DeleteByPattern"]++
	if m.deleteByPatternFunc != nil {
		return m.deleteByPatternFunc(ctx, pattern)
	}
	return nil
}

// --- テスト用ヘルパー ---

func createTestPlanForCache(t *testing.T, userID entities.UserID) *aggregates.FinancialPlan {
	t.Helper()

	income, err := valueobjects.NewMoneyJPY(400000)
	if err != nil {
		t.Fatalf("月収の作成に失敗: %v", err)
	}

	investReturn, err := valueobjects.NewRate(5.0)
	if err != nil {
		t.Fatalf("投資利回りの作成に失敗: %v", err)
	}

	inflRate, err := valueobjects.NewRate(2.0)
	if err != nil {
		t.Fatalf("インフレ率の作成に失敗: %v", err)
	}

	profile, err := entities.NewFinancialProfile(
		userID,
		income,
		entities.ExpenseCollection{},
		entities.SavingsCollection{},
		investReturn,
		inflRate,
	)
	if err != nil {
		t.Fatalf("財務プロファイルの作成に失敗: %v", err)
	}

	plan, err := aggregates.NewFinancialPlan(profile)
	if err != nil {
		t.Fatalf("財務計画の作成に失敗: %v", err)
	}

	return plan
}

// --- テスト ---

func TestCachedFinancialPlanRepository_FindByUserID_CacheHit(t *testing.T) {
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	plan := createTestPlanForCache(t, userID)

	// キャッシュにDTOを事前に書き込む
	dto := financialPlanToDTO(plan)

	mockRepo := newMockFinancialPlanRepo()
	mockCache := newMockCacheClient()
	mockCache.getJSONFunc = func(ctx context.Context, key string, dest any) error {
		if p, ok := dest.(*financialPlanCacheDTO); ok {
			*p = dto
		}
		return nil // キャッシュヒット
	}

	repo := NewCachedFinancialPlanRepository(mockRepo, mockCache)

	result, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("FindByUserID エラー: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	// DBが呼ばれないことを確認
	if mockRepo.callCount["FindByUserID"] != 0 {
		t.Errorf("キャッシュヒット時にDBが呼ばれました（呼び出し回数: %d）", mockRepo.callCount["FindByUserID"])
	}
}

func TestCachedFinancialPlanRepository_FindByUserID_CacheMiss(t *testing.T) {
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	plan := createTestPlanForCache(t, userID)

	mockRepo := newMockFinancialPlanRepo()
	mockRepo.findByUserIDFunc = func(ctx context.Context, uid entities.UserID) (*aggregates.FinancialPlan, error) {
		return plan, nil
	}

	mockCache := newMockCacheClient()
	// デフォルトで GetJSON は redis.Nil を返す（キャッシュミス）

	repo := NewCachedFinancialPlanRepository(mockRepo, mockCache)

	result, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("FindByUserID エラー: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	// DBが呼ばれることを確認
	if mockRepo.callCount["FindByUserID"] != 1 {
		t.Errorf("キャッシュミス時にDBが呼ばれませんでした（呼び出し回数: %d）", mockRepo.callCount["FindByUserID"])
	}

	// キャッシュへの書き込みが行われることを確認（FindByID + FindByUserID の2キー）
	if mockCache.callCount["SetJSON"] < 2 {
		t.Errorf("キャッシュへの書き込みが不足しています（呼び出し回数: %d）", mockCache.callCount["SetJSON"])
	}
}

func TestCachedFinancialPlanRepository_FindByUserID_RedisFailure_Fallback(t *testing.T) {
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	plan := createTestPlanForCache(t, userID)

	mockRepo := newMockFinancialPlanRepo()
	mockRepo.findByUserIDFunc = func(ctx context.Context, uid entities.UserID) (*aggregates.FinancialPlan, error) {
		return plan, nil
	}

	mockCache := newMockCacheClient()
	mockCache.getJSONFunc = func(ctx context.Context, key string, dest any) error {
		return errors.New("redis: connection refused") // Redis障害
	}

	repo := NewCachedFinancialPlanRepository(mockRepo, mockCache)

	result, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Redis障害時にエラーが返されました（fail-openのはず）: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	// fail-open: DBにフォールバックされることを確認
	if mockRepo.callCount["FindByUserID"] != 1 {
		t.Errorf("Redis障害時にDBにフォールバックしませんでした（呼び出し回数: %d）", mockRepo.callCount["FindByUserID"])
	}
}

func TestCachedFinancialPlanRepository_Save_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	plan := createTestPlanForCache(t, userID)

	mockRepo := newMockFinancialPlanRepo()
	deletedKeys := []string{}
	mockCache := newMockCacheClient()
	mockCache.deleteFunc = func(ctx context.Context, keys ...string) error {
		deletedKeys = append(deletedKeys, keys...)
		return nil
	}

	repo := NewCachedFinancialPlanRepository(mockRepo, mockCache)

	if err := repo.Save(ctx, plan); err != nil {
		t.Fatalf("Save エラー: %v", err)
	}

	// キャッシュが削除されることを確認
	if mockCache.callCount["Delete"] == 0 {
		t.Error("Save後にキャッシュが削除されませんでした")
	}

	// 正しいキーが削除されることを確認
	expectedByID := financialPlanByIDKey(string(plan.ID()))
	expectedByUserID := financialPlanByUserIDKey(string(plan.Profile().UserID()))

	hasIDKey := false
	hasUserIDKey := false
	for _, k := range deletedKeys {
		if k == expectedByID {
			hasIDKey = true
		}
		if k == expectedByUserID {
			hasUserIDKey = true
		}
	}
	if !hasIDKey {
		t.Errorf("IDキャッシュが削除されませんでした: %s", expectedByID)
	}
	if !hasUserIDKey {
		t.Errorf("UserIDキャッシュが削除されませんでした: %s", expectedByUserID)
	}
}

func TestCachedFinancialPlanRepository_DTORoundTrip(t *testing.T) {
	userID := entities.UserID("test-user-id")
	plan := createTestPlanForCache(t, userID)

	dto := financialPlanToDTO(plan)
	restored, err := financialPlanFromDTO(dto)
	if err != nil {
		t.Fatalf("DTO復元エラー: %v", err)
	}

	if string(restored.ID()) != string(plan.ID()) {
		t.Errorf("IDが一致しません: got %s, want %s", restored.ID(), plan.ID())
	}
	if string(restored.Profile().UserID()) != string(plan.Profile().UserID()) {
		t.Errorf("UserIDが一致しません: got %s, want %s", restored.Profile().UserID(), plan.Profile().UserID())
	}
	if restored.Profile().MonthlyIncome().Amount() != plan.Profile().MonthlyIncome().Amount() {
		t.Errorf("月収が一致しません: got %f, want %f", restored.Profile().MonthlyIncome().Amount(), plan.Profile().MonthlyIncome().Amount())
	}
}

// IsNil は redis.Nil エラーかどうかを判定するヘルパー（テストでインポートせずに使用）
func isNilError(err error) bool {
	return redisinfra.IsNil(err)
}
