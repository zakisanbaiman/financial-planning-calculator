package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	goredis "github.com/redis/go-redis/v9"
)

// --- モック: GoalRepository ---

type mockGoalRepository struct {
	findByIDFunc           func(ctx context.Context, id entities.GoalID) (*entities.Goal, error)
	findByUserIDFunc       func(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error)
	findActiveByUserIDFunc func(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error)
	findByTypeFunc         func(ctx context.Context, userID entities.UserID, goalType entities.GoalType) ([]*entities.Goal, error)
	saveFunc               func(ctx context.Context, goal *entities.Goal) error
	updateFunc             func(ctx context.Context, goal *entities.Goal) error
	deleteFunc             func(ctx context.Context, id entities.GoalID) error
	existsFunc             func(ctx context.Context, id entities.GoalID) (bool, error)
	countActiveFunc        func(ctx context.Context, userID entities.UserID, goalType entities.GoalType) (int, error)
	callCount              map[string]int
}

func newMockGoalRepo() *mockGoalRepository {
	return &mockGoalRepository{callCount: make(map[string]int)}
}

func (m *mockGoalRepository) FindByID(ctx context.Context, id entities.GoalID) (*entities.Goal, error) {
	m.callCount["FindByID"]++
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockGoalRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
	m.callCount["FindByUserID"]++
	if m.findByUserIDFunc != nil {
		return m.findByUserIDFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockGoalRepository) FindActiveGoalsByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
	m.callCount["FindActiveGoalsByUserID"]++
	if m.findActiveByUserIDFunc != nil {
		return m.findActiveByUserIDFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockGoalRepository) FindByUserIDAndType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) ([]*entities.Goal, error) {
	m.callCount["FindByUserIDAndType"]++
	if m.findByTypeFunc != nil {
		return m.findByTypeFunc(ctx, userID, goalType)
	}
	return nil, nil
}

func (m *mockGoalRepository) Save(ctx context.Context, goal *entities.Goal) error {
	m.callCount["Save"]++
	if m.saveFunc != nil {
		return m.saveFunc(ctx, goal)
	}
	return nil
}

func (m *mockGoalRepository) Update(ctx context.Context, goal *entities.Goal) error {
	m.callCount["Update"]++
	if m.updateFunc != nil {
		return m.updateFunc(ctx, goal)
	}
	return nil
}

func (m *mockGoalRepository) Delete(ctx context.Context, id entities.GoalID) error {
	m.callCount["Delete"]++
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockGoalRepository) Exists(ctx context.Context, id entities.GoalID) (bool, error) {
	m.callCount["Exists"]++
	if m.existsFunc != nil {
		return m.existsFunc(ctx, id)
	}
	return false, nil
}

func (m *mockGoalRepository) CountActiveGoalsByType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) (int, error) {
	m.callCount["CountActiveGoalsByType"]++
	if m.countActiveFunc != nil {
		return m.countActiveFunc(ctx, userID, goalType)
	}
	return 0, nil
}

// --- テスト用ヘルパー ---

func createTestGoal(t *testing.T, userID entities.UserID) *entities.Goal {
	t.Helper()

	targetAmount, err := valueobjects.NewMoneyJPY(1000000)
	if err != nil {
		t.Fatalf("目標金額の作成に失敗: %v", err)
	}
	monthlyContrib, err := valueobjects.NewMoneyJPY(50000)
	if err != nil {
		t.Fatalf("月間拠出額の作成に失敗: %v", err)
	}

	goal, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		"テスト貯蓄目標",
		targetAmount,
		time.Now().AddDate(2, 0, 0),
		monthlyContrib,
	)
	if err != nil {
		t.Fatalf("ゴールの作成に失敗: %v", err)
	}
	return goal
}

// --- テスト ---

func TestCachedGoalRepository_FindByUserID_CacheHit(t *testing.T) {
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	goal := createTestGoal(t, userID)
	goals := []*entities.Goal{goal}

	dtos := goalsToDTOs(goals)

	mockRepo := newMockGoalRepo()
	mockCache := newMockCacheClient()
	mockCache.getJSONFunc = func(ctx context.Context, key string, dest any) error {
		if p, ok := dest.(*[]goalCacheDTO); ok {
			*p = dtos
		}
		return nil // キャッシュヒット
	}

	repo := NewCachedGoalRepository(mockRepo, mockCache)

	result, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("FindByUserID エラー: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("ゴール数が一致しません: got %d, want 1", len(result))
	}

	// DBが呼ばれないことを確認
	if mockRepo.callCount["FindByUserID"] != 0 {
		t.Errorf("キャッシュヒット時にDBが呼ばれました（呼び出し回数: %d）", mockRepo.callCount["FindByUserID"])
	}
}

func TestCachedGoalRepository_FindByUserID_CacheMiss(t *testing.T) {
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	goal := createTestGoal(t, userID)
	goals := []*entities.Goal{goal}

	mockRepo := newMockGoalRepo()
	mockRepo.findByUserIDFunc = func(ctx context.Context, uid entities.UserID) ([]*entities.Goal, error) {
		return goals, nil
	}

	mockCache := newMockCacheClient()
	// デフォルト: GetJSON は redis.Nil（キャッシュミス）

	repo := NewCachedGoalRepository(mockRepo, mockCache)

	result, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("FindByUserID エラー: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("ゴール数が一致しません: got %d, want 1", len(result))
	}

	if mockRepo.callCount["FindByUserID"] != 1 {
		t.Errorf("キャッシュミス時にDBが呼ばれませんでした（呼び出し回数: %d）", mockRepo.callCount["FindByUserID"])
	}
	if mockCache.callCount["SetJSON"] != 1 {
		t.Errorf("キャッシュへの書き込み回数が不正: got %d, want 1", mockCache.callCount["SetJSON"])
	}
}

func TestCachedGoalRepository_FindByUserID_RedisFailure_Fallback(t *testing.T) {
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	goal := createTestGoal(t, userID)
	goals := []*entities.Goal{goal}

	mockRepo := newMockGoalRepo()
	mockRepo.findByUserIDFunc = func(ctx context.Context, uid entities.UserID) ([]*entities.Goal, error) {
		return goals, nil
	}

	mockCache := newMockCacheClient()
	mockCache.getJSONFunc = func(ctx context.Context, key string, dest any) error {
		return errors.New("redis: connection timeout")
	}

	repo := NewCachedGoalRepository(mockRepo, mockCache)

	result, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Redis障害時にエラーが返されました（fail-openのはず）: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("結果数が不正: got %d, want 1", len(result))
	}

	if mockRepo.callCount["FindByUserID"] != 1 {
		t.Errorf("Redis障害時にDBにフォールバックしませんでした（呼び出し回数: %d）", mockRepo.callCount["FindByUserID"])
	}
}

func TestCachedGoalRepository_Save_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	goal := createTestGoal(t, userID)

	mockRepo := newMockGoalRepo()
	deletedKeys := []string{}
	mockCache := newMockCacheClient()
	mockCache.deleteFunc = func(ctx context.Context, keys ...string) error {
		deletedKeys = append(deletedKeys, keys...)
		return nil
	}

	repo := NewCachedGoalRepository(mockRepo, mockCache)

	if err := repo.Save(ctx, goal); err != nil {
		t.Fatalf("Save エラー: %v", err)
	}

	if mockCache.callCount["Delete"] == 0 {
		t.Error("Save後にキャッシュが削除されませんでした")
	}

	// FindByUserID と FindActiveGoalsByUserID の両キーが削除されることを確認
	expectedByUserID := goalsByUserIDKey(string(userID))
	expectedActive := activeGoalsByUserIDKey(string(userID))

	hasByUserID := false
	hasActive := false
	for _, k := range deletedKeys {
		if k == expectedByUserID {
			hasByUserID = true
		}
		if k == expectedActive {
			hasActive = true
		}
	}
	if !hasByUserID {
		t.Errorf("FindByUserIDキャッシュが削除されませんでした: %s", expectedByUserID)
	}
	if !hasActive {
		t.Errorf("FindActiveGoalsByUserIDキャッシュが削除されませんでした: %s", expectedActive)
	}
}

func TestCachedGoalRepository_GoalDTORoundTrip(t *testing.T) {
	userID := entities.UserID("test-user-id")
	original := createTestGoal(t, userID)

	dto := goalToDTO(original)
	restored, err := goalFromDTO(dto)
	if err != nil {
		t.Fatalf("DTO復元エラー: %v", err)
	}

	if string(restored.ID()) != string(original.ID()) {
		t.Errorf("IDが一致しません: got %s, want %s", restored.ID(), original.ID())
	}
	if restored.Title() != original.Title() {
		t.Errorf("タイトルが一致しません: got %s, want %s", restored.Title(), original.Title())
	}
	if restored.TargetAmount().Amount() != original.TargetAmount().Amount() {
		t.Errorf("目標金額が一致しません: got %f, want %f", restored.TargetAmount().Amount(), original.TargetAmount().Amount())
	}
	if restored.IsActive() != original.IsActive() {
		t.Errorf("IsActiveが一致しません: got %v, want %v", restored.IsActive(), original.IsActive())
	}
}

// redis.Nil 定数が使われることを確認するテスト
func TestCacheMissDetection(t *testing.T) {
	if !isNilError(goredis.Nil) {
		t.Error("redis.Nil が IsNil で検出されませんでした")
	}
	if isNilError(errors.New("some other error")) {
		t.Error("一般エラーが IsNil で誤検出されました")
	}
}
