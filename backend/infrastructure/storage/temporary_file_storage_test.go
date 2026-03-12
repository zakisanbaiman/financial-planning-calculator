package storage

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTemporaryFileStorage_SaveAndGetFile(t *testing.T) {
	// テスト用の一時ディレクトリ
	tmpDir := "/tmp/test-storage-" + time.Now().Format("20060102150405")
	defer os.RemoveAll(tmpDir)

	// ストレージを作成
	storage, err := NewTemporaryFileStorage(tmpDir, "test-secret-key", 1*time.Hour, 1*time.Hour)
	if err != nil {
		t.Fatalf("ストレージの作成に失敗: %v", err)
	}

	// ファイルを保存
	testData := []byte("これはテストPDFです")
	token, expiresAt, err := storage.SaveFile("test-report.pdf", testData)
	if err != nil {
		t.Fatalf("ファイルの保存に失敗: %v", err)
	}

	// 有効期限の検証
	if expiresAt.IsZero() {
		t.Error("有効期限が設定されていません")
	}
	if !time.Now().Before(expiresAt) {
		t.Errorf("有効期限が過去の日時です: %v", expiresAt)
	}

	// ファイルを取得
	data, fileName, ownerUserID, err := storage.GetFile(token)
	if err != nil {
		t.Fatalf("ファイルの取得に失敗: %v", err)
	}

	// データの検証
	if string(data) != string(testData) {
		t.Errorf("データが一致しません: got %s, want %s", string(data), string(testData))
	}

	// ファイル名の検証
	if fileName != "test-report.pdf" {
		t.Errorf("ファイル名が一致しません: got %s, want test-report.pdf", fileName)
	}

	// ownerUserIDはファイル名のプレフィックスから取得（"test-report.pdf" の場合は "test"）
	_ = ownerUserID
}

func TestTemporaryFileStorage_ExpiredFile(t *testing.T) {
	// テスト用の一時ディレクトリ
	tmpDir := "/tmp/test-storage-expired-" + time.Now().Format("20060102150405")
	defer os.RemoveAll(tmpDir)

	// 1秒で期限切れになるストレージ
	storage, err := NewTemporaryFileStorage(tmpDir, "test-secret-key", 1*time.Second, 1*time.Hour)
	if err != nil {
		t.Fatalf("ストレージの作成に失敗: %v", err)
	}

	// ファイルを保存
	testData := []byte("期限切れテスト")
	token, _, err := storage.SaveFile("expired-test.pdf", testData)
	if err != nil {
		t.Fatalf("ファイルの保存に失敗: %v", err)
	}

	// すぐにアクセス - 成功するはず
	_, _, _, err = storage.GetFile(token)
	if err != nil {
		t.Errorf("有効期限内のアクセスが失敗: %v", err)
	}

	// 2秒待つ（有効期限切れ）
	time.Sleep(2 * time.Second)

	// アクセス - 失敗するはず
	_, _, _, err = storage.GetFile(token)
	if err == nil {
		t.Error("期限切れファイルへのアクセスが成功してしまいました")
	}
}

func TestTemporaryFileStorage_InvalidToken(t *testing.T) {
	// テスト用の一時ディレクトリ
	tmpDir := "/tmp/test-storage-invalid-" + time.Now().Format("20060102150405")
	defer os.RemoveAll(tmpDir)

	// ストレージを作成
	storage, err := NewTemporaryFileStorage(tmpDir, "test-secret-key", 1*time.Hour, 1*time.Hour)
	if err != nil {
		t.Fatalf("ストレージの作成に失敗: %v", err)
	}

	// 無効なトークンでアクセス
	_, _, _, err = storage.GetFile("invalid-token")
	if err == nil {
		t.Error("無効なトークンでのアクセスが成功してしまいました")
	}
}

func TestTemporaryFileStorage_FileCount(t *testing.T) {
	// テスト用の一時ディレクトリ
	tmpDir := "/tmp/test-storage-count-" + time.Now().Format("20060102150405")
	defer os.RemoveAll(tmpDir)

	// ストレージを作成
	storage, err := NewTemporaryFileStorage(tmpDir, "test-secret-key", 1*time.Hour, 1*time.Hour)
	if err != nil {
		t.Fatalf("ストレージの作成に失敗: %v", err)
	}

	// 初期状態
	if count := storage.GetFileCount(); count != 0 {
		t.Errorf("初期ファイル数が0ではありません: %d", count)
	}

	// ファイルを3つ保存（異なるファイル名）
	for i := 0; i < 3; i++ {
		fileName := fmt.Sprintf("test-%d.pdf", i)
		_, _, err := storage.SaveFile(fileName, []byte("test"))
		if err != nil {
			t.Fatalf("ファイルの保存に失敗: %v", err)
		}
	}

	// ファイル数を確認
	if count := storage.GetFileCount(); count != 3 {
		t.Errorf("ファイル数が一致しません: got %d, want 3", count)
	}
}
