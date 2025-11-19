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
	storage, err := NewTemporaryFileStorage(tmpDir, "test-secret-key", 1*time.Hour)
	if err != nil {
		t.Fatalf("ストレージの作成に失敗: %v", err)
	}

	// ファイルを保存
	testData := []byte("これはテストPDFです")
	token, metadata, err := storage.SaveFile("test-report.pdf", testData)
	if err != nil {
		t.Fatalf("ファイルの保存に失敗: %v", err)
	}

	// メタデータの検証
	if metadata.FileName != "test-report.pdf" {
		t.Errorf("ファイル名が一致しません: got %s, want test-report.pdf", metadata.FileName)
	}
	if metadata.FileSize != int64(len(testData)) {
		t.Errorf("ファイルサイズが一致しません: got %d, want %d", metadata.FileSize, len(testData))
	}

	// ファイルを取得
	data, retrievedMetadata, err := storage.GetFile(token)
	if err != nil {
		t.Fatalf("ファイルの取得に失敗: %v", err)
	}

	// データの検証
	if string(data) != string(testData) {
		t.Errorf("データが一致しません: got %s, want %s", string(data), string(testData))
	}

	// メタデータの検証
	if retrievedMetadata.FileName != metadata.FileName {
		t.Errorf("ファイル名が一致しません")
	}
}

func TestTemporaryFileStorage_ExpiredFile(t *testing.T) {
	// テスト用の一時ディレクトリ
	tmpDir := "/tmp/test-storage-expired-" + time.Now().Format("20060102150405")
	defer os.RemoveAll(tmpDir)

	// 1秒で期限切れになるストレージ
	storage, err := NewTemporaryFileStorage(tmpDir, "test-secret-key", 1*time.Second)
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
	_, _, err = storage.GetFile(token)
	if err != nil {
		t.Errorf("有効期限内のアクセスが失敗: %v", err)
	}

	// 2秒待つ（有効期限切れ）
	time.Sleep(2 * time.Second)

	// アクセス - 失敗するはず
	_, _, err = storage.GetFile(token)
	if err == nil {
		t.Error("期限切れファイルへのアクセスが成功してしまいました")
	}
}

func TestTemporaryFileStorage_InvalidToken(t *testing.T) {
	// テスト用の一時ディレクトリ
	tmpDir := "/tmp/test-storage-invalid-" + time.Now().Format("20060102150405")
	defer os.RemoveAll(tmpDir)

	// ストレージを作成
	storage, err := NewTemporaryFileStorage(tmpDir, "test-secret-key", 1*time.Hour)
	if err != nil {
		t.Fatalf("ストレージの作成に失敗: %v", err)
	}

	// 無効なトークンでアクセス
	_, _, err = storage.GetFile("invalid-token")
	if err == nil {
		t.Error("無効なトークンでのアクセスが成功してしまいました")
	}
}

func TestTemporaryFileStorage_FileCount(t *testing.T) {
	// テスト用の一時ディレクトリ
	tmpDir := "/tmp/test-storage-count-" + time.Now().Format("20060102150405")
	defer os.RemoveAll(tmpDir)

	// ストレージを作成
	storage, err := NewTemporaryFileStorage(tmpDir, "test-secret-key", 1*time.Hour)
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
