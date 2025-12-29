package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TemporaryFileStorage は一時ファイルの保存と管理を行う
type TemporaryFileStorage struct {
	baseDir         string
	secretKey       []byte
	expiryTime      time.Duration
	cleanupInterval time.Duration
	files           map[string]*FileMetadata
	mu              sync.RWMutex
}

// FileMetadata はファイルのメタデータ
type FileMetadata struct {
	FilePath  string
	FileName  string
	FileSize  int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewTemporaryFileStorage は新しいTemporaryFileStorageを作成する
func NewTemporaryFileStorage(baseDir string, secretKey string, expiryDuration time.Duration, cleanupInterval time.Duration) (*TemporaryFileStorage, error) {
	// ベースディレクトリを作成
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("ベースディレクトリの作成に失敗: %w", err)
	}

	storage := &TemporaryFileStorage{
		baseDir:         baseDir,
		secretKey:       []byte(secretKey),
		expiryTime:      expiryDuration,
		cleanupInterval: cleanupInterval,
		files:           make(map[string]*FileMetadata),
	}

	// 定期的に期限切れファイルをクリーンアップ
	go storage.startCleanupRoutine()

	return storage, nil
}

// SaveFile はファイルを保存し、署名付きトークンを返す
func (s *TemporaryFileStorage) SaveFile(fileName string, data []byte) (token string, metadata *FileMetadata, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// ファイルパスを生成
	timestamp := time.Now().Unix()
	safeFileName := fmt.Sprintf("%d_%s", timestamp, fileName)
	filePath := filepath.Join(s.baseDir, safeFileName)

	// ファイルを保存
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", nil, fmt.Errorf("ファイルの保存に失敗: %w", err)
	}

	// メタデータを作成
	now := time.Now()
	metadata = &FileMetadata{
		FilePath:  filePath,
		FileName:  fileName,
		FileSize:  int64(len(data)),
		CreatedAt: now,
		ExpiresAt: now.Add(s.expiryTime),
	}

	// 署名付きトークンを生成
	token = s.generateToken(safeFileName, metadata.ExpiresAt)

	// メタデータを保存
	s.files[token] = metadata

	return token, metadata, nil
}

// GetFile はトークンからファイルを取得する
func (s *TemporaryFileStorage) GetFile(token string) ([]byte, *FileMetadata, error) {
	s.mu.RLock()
	metadata, exists := s.files[token]
	s.mu.RUnlock()

	if !exists {
		return nil, nil, fmt.Errorf("ファイルが見つかりません")
	}

	// 期限切れチェック（現在時刻が有効期限より前でない = 期限切れ）
	if !time.Now().Before(metadata.ExpiresAt) {
		_ = s.deleteFile(token) // 削除エラーは無視（既に削除されている可能性がある）
		return nil, nil, fmt.Errorf("ファイルの有効期限が切れています")
	}

	// トークンの検証
	if !s.verifyToken(token, filepath.Base(metadata.FilePath), metadata.ExpiresAt) {
		return nil, nil, fmt.Errorf("無効なトークンです")
	}

	// ファイルを読み込み
	data, err := os.ReadFile(metadata.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	return data, metadata, nil
}

// generateToken は署名付きトークンを生成する
func (s *TemporaryFileStorage) generateToken(fileName string, expiresAt time.Time) string {
	// HMAC-SHA256で署名を生成
	message := fmt.Sprintf("%s:%d", fileName, expiresAt.Unix())
	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))

	// トークン = ファイル名:有効期限:署名
	token := fmt.Sprintf("%s:%d:%s", fileName, expiresAt.Unix(), signature)
	return token
}

// verifyToken はトークンを検証する
func (s *TemporaryFileStorage) verifyToken(token, fileName string, expiresAt time.Time) bool {
	expectedToken := s.generateToken(fileName, expiresAt)
	return hmac.Equal([]byte(token), []byte(expectedToken))
}

// deleteFile はファイルを削除する
func (s *TemporaryFileStorage) deleteFile(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	metadata, exists := s.files[token]
	if !exists {
		return nil
	}

	// ファイルを削除
	if err := os.Remove(metadata.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("ファイルの削除に失敗: %w", err)
	}

	// メタデータを削除
	delete(s.files, token)

	return nil
}

// startCleanupRoutine は定期的に期限切れファイルをクリーンアップする
func (s *TemporaryFileStorage) startCleanupRoutine() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupExpiredFiles()
	}
}

// cleanupExpiredFiles は期限切れファイルを削除する
func (s *TemporaryFileStorage) cleanupExpiredFiles() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, metadata := range s.files {
		// 現在時刻が有効期限より前でない = 期限切れ
		if !now.Before(metadata.ExpiresAt) {
			// ファイルを削除
			if err := os.Remove(metadata.FilePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("期限切れファイルの削除に失敗: %v\n", err)
			}
			// メタデータを削除
			delete(s.files, token)
		}
	}
}

// GetFileCount は保存されているファイル数を返す
func (s *TemporaryFileStorage) GetFileCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.files)
}
