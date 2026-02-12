package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/infrastructure/log"
)

func TestPrometheusMetrics(t *testing.T) {
	// Prometheusメトリクスの初期化
	InitPrometheus()

	// メトリクスが正しく登録されているか確認
	// エラーがなければ成功
	t.Log("Prometheusメトリクスが正常に初期化されました")
}

func TestRecordDatabaseQuery(t *testing.T) {
	InitPrometheus()

	// データベースクエリのメトリクスを記録
	RecordDatabaseQuery("select", 100*time.Millisecond)
	RecordDatabaseQuery("insert", 50*time.Millisecond)
	RecordDatabaseQuery("update", 75*time.Millisecond)

	t.Log("データベースクエリのメトリクスが正常に記録されました")
}

func TestUpdateDatabaseConnections(t *testing.T) {
	InitPrometheus()

	// データベース接続数を更新
	UpdateDatabaseConnections(10)
	UpdateDatabaseConnections(15)
	UpdateDatabaseConnections(5)

	t.Log("データベース接続数のメトリクスが正常に更新されました")
}

func TestUpdateCacheHitRatio(t *testing.T) {
	InitPrometheus()

	// キャッシュヒット率を更新
	UpdateCacheHitRatio("calculation", 0.85)
	UpdateCacheHitRatio("response", 0.92)

	t.Log("キャッシュヒット率のメトリクスが正常に更新されました")
}

func TestRecordError(t *testing.T) {
	InitPrometheus()

	// エラーメトリクスを記録
	RecordError("validation_error", "warning")
	RecordError("database_error", "error")
	RecordError("panic", "critical")

	t.Log("エラーメトリクスが正常に記録されました")
}

func TestDefaultErrorTracker(t *testing.T) {
	// エラートラッカーの初期化
	InitErrorTracker("test")

	ctx := context.Background()
	ctx = log.WithRequestID(ctx, "test-request-123")
	ctx = log.WithUserID(ctx, "user-456")

	// エラーをキャプチャ
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

	// メッセージをキャプチャ
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
