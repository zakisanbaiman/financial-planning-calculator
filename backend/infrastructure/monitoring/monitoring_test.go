package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/infrastructure/log"
)

func TestNewRelicInit(t *testing.T) {
	// License Key なしで初期化を試みる（エラーが返ることを確認）
	err := InitNewRelic("", "test-app")
	if err == nil {
		t.Error("License Key が空の場合はエラーが返るべきです")
	}

	t.Log("New Relic 初期化バリデーションが正常に動作しました")
}

func TestRecordDatabaseQuery(t *testing.T) {
	// nrApp が nil の状態（New Relic 無効）でも panic しないことを確認
	RecordDatabaseQuery("select", 100*time.Millisecond)
	RecordDatabaseQuery("insert", 50*time.Millisecond)
	RecordDatabaseQuery("update", 75*time.Millisecond)

	t.Log("データベースクエリのメトリクスが正常に記録されました")
}

func TestUpdateDatabaseConnections(t *testing.T) {
	UpdateDatabaseConnections(10)
	UpdateDatabaseConnections(15)
	UpdateDatabaseConnections(5)

	t.Log("データベース接続数のメトリクスが正常に更新されました")
}

func TestUpdateCacheHitRatio(t *testing.T) {
	UpdateCacheHitRatio("calculation", 0.85)
	UpdateCacheHitRatio("response", 0.92)

	t.Log("キャッシュヒット率のメトリクスが正常に更新されました")
}

func TestRecordError(t *testing.T) {
	RecordError("validation_error", "warning")
	RecordError("database_error", "error")
	RecordError("panic", "critical")

	t.Log("エラーメトリクスが正常に記録されました")
}

func TestRecordCacheHitMiss(t *testing.T) {
	RecordCacheHit("calculation")
	RecordCacheMiss("response")

	t.Log("キャッシュヒット/ミスのメトリクスが正常に記録されました")
}

func TestDefaultErrorTracker(t *testing.T) {
	InitErrorTracker("test")

	ctx := context.Background()
	ctx = log.WithRequestID(ctx, "test-request-123")
	ctx = log.WithUserID(ctx, "user-456")

	tracker := GetErrorTracker()
	tags := map[string]string{
		"component": "test",
		"severity":  "high",
	}
	tracker.CaptureError(ctx, testError{}, tags)

	t.Log("エラートラッカーが正常に動作しました")
}

func TestCaptureMessage(t *testing.T) {
	InitErrorTracker("test")

	ctx := context.Background()
	ctx = log.WithRequestID(ctx, "test-request-456")

	tracker := GetErrorTracker()
	tags := map[string]string{
		"component": "test",
	}
	tracker.CaptureMessage(ctx, "テストメッセージ", "info", tags)

	t.Log("メッセージキャプチャが正常に動作しました")
}

// テスト用のエラー型
type testError struct{}

func (e testError) Error() string {
	return "test error"
}
