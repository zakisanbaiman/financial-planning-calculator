# expires_at実装ガイド

## 概要

`expires_at`は一時的なリソース（ファイル、URL、トークンなど）に有効期限を設定する機能です。
このプロジェクトでは、PDFレポートの一時ダウンロードURLに実装しています。

## アーキテクチャ

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ 1. POST /api/reports/export
       ▼
┌─────────────────────────────┐
│  ReportsController          │
│  - ExportReportToPDF()      │
└──────┬──────────────────────┘
       │ 2. Generate PDF
       ▼
┌─────────────────────────────┐
│  GenerateReportsUseCase     │
│  - ExportReportToPDF()      │
└──────┬──────────────────────┘
       │ 3. Save with expiry
       ▼
┌─────────────────────────────┐
│  TemporaryFileStorage       │
│  - SaveFile()               │
│  - generateToken()          │
│  - cleanupExpiredFiles()    │
└──────┬──────────────────────┘
       │ 4. Return signed URL
       ▼
┌─────────────────────────────┐
│  Response                   │
│  {                          │
│    "download_url": "...",   │
│    "expires_at": "2024-..." │
│  }                          │
└─────────────────────────────┘
```

## 実装の詳細

### 1. TemporaryFileStorage

一時ファイルの保存と管理を行うサービスです。

**主な機能:**
- ファイルの保存と署名付きトークンの生成
- HMAC-SHA256による改ざん防止
- 自動クリーンアップ（1時間ごと）
- スレッドセーフな実装

**コード例:**
```go
// ストレージの初期化
storage, err := storage.NewTemporaryFileStorage(
    "/tmp/reports",           // 保存先ディレクトリ
    "secret-key",             // 署名用の秘密鍵
    24 * time.Hour,           // 有効期限（24時間）
)

// ファイルの保存
token, metadata, err := storage.SaveFile("report.pdf", pdfData)

// トークンの例:
// "1234567890_report.pdf:1700000000:a1b2c3d4e5f6..."
//  ↑ファイル名      ↑有効期限    ↑HMAC署名
```

### 2. トークンの構造

```
{filename}:{expiry_unix}:{hmac_signature}
```

- **filename**: タイムスタンプ付きファイル名
- **expiry_unix**: Unix時間での有効期限
- **hmac_signature**: HMAC-SHA256署名（改ざん防止）

### 3. セキュリティ対策

#### 改ざん防止
```go
func (s *TemporaryFileStorage) generateToken(fileName string, expiresAt time.Time) string {
    message := fmt.Sprintf("%s:%d", fileName, expiresAt.Unix())
    h := hmac.New(sha256.New, s.secretKey)
    h.Write([]byte(message))
    signature := hex.EncodeToString(h.Sum(nil))
    return fmt.Sprintf("%s:%d:%s", fileName, expiresAt.Unix(), signature)
}
```

#### 検証
```go
func (s *TemporaryFileStorage) verifyToken(token, fileName string, expiresAt time.Time) bool {
    expectedToken := s.generateToken(fileName, expiresAt)
    return hmac.Equal([]byte(token), []byte(expectedToken))
}
```

### 4. 自動クリーンアップ

期限切れファイルを1時間ごとに自動削除します。

```go
func (s *TemporaryFileStorage) startCleanupRoutine() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        s.cleanupExpiredFiles()
    }
}
```

## 使用例

### レポートのエクスポート

```bash
# 1. レポートをエクスポート
curl -X POST http://localhost:8080/api/reports/export \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "report_type": "comprehensive",
    "format": "pdf"
  }'

# レスポンス
{
  "file_name": "comprehensive_report_20241116_123456.pdf",
  "file_size": 524288,
  "download_url": "/api/reports/download/1234567890_report.pdf:1700000000:a1b2c3d4...",
  "expires_at": "2024-11-17T12:34:56Z"
}
```

### ダウンロード

```bash
# 2. トークンを使ってダウンロード
curl -O http://localhost:8080/api/reports/download/1234567890_report.pdf:1700000000:a1b2c3d4...
```

### 期限切れの場合

```bash
# 24時間後にアクセス
curl http://localhost:8080/api/reports/download/expired_token

# レスポンス
{
  "error": "expired",
  "message": "ファイルの有効期限が切れています"
}
```

## 設定

### 環境変数

```bash
# docker-compose.yml または .env
TEMP_FILE_DIR=/tmp/financial-planning-reports
TEMP_FILE_SECRET=your-secret-key-here-change-in-production
TEMP_FILE_EXPIRY=24h
```

### 本番環境での推奨設定

```bash
# 強力な秘密鍵を使用
TEMP_FILE_SECRET=$(openssl rand -hex 32)

# 短めの有効期限
TEMP_FILE_EXPIRY=1h

# 専用のストレージディレクトリ
TEMP_FILE_DIR=/var/lib/app/temp-files
```

## ベストプラクティス

### 1. 秘密鍵の管理

```go
// ❌ ハードコードしない
secretKey := "my-secret-key"

// ✅ 環境変数から読み込む
secretKey := os.Getenv("TEMP_FILE_SECRET")
if secretKey == "" {
    log.Fatal("TEMP_FILE_SECRET is required")
}
```

### 2. 有効期限の設定

```go
// 用途に応じて適切な期限を設定
var expiry time.Duration
switch reportType {
case "quick_summary":
    expiry = 1 * time.Hour  // 簡易レポートは短め
case "comprehensive":
    expiry = 24 * time.Hour // 詳細レポートは長め
case "sensitive":
    expiry = 15 * time.Minute // 機密情報は超短め
}
```

### 3. エラーハンドリング

```go
data, metadata, err := storage.GetFile(token)
if err != nil {
    if strings.Contains(err.Error(), "有効期限") {
        return c.JSON(http.StatusGone, map[string]interface{}{
            "error": "expired",
            "message": "ファイルの有効期限が切れています",
        })
    }
    if strings.Contains(err.Error(), "見つかりません") {
        return c.JSON(http.StatusNotFound, map[string]interface{}{
            "error": "not_found",
            "message": "ファイルが見つかりません",
        })
    }
    return c.JSON(http.StatusInternalServerError, map[string]interface{}{
        "error": "internal_error",
        "message": "ファイルの取得に失敗しました",
    })
}
```

### 4. ログ記録

```go
log.Printf("ファイル保存: user=%s, file=%s, expires=%s",
    userID, fileName, expiresAt.Format(time.RFC3339))

log.Printf("ファイルダウンロード: token=%s, remaining=%s",
    token, time.Until(expiresAt))

log.Printf("期限切れファイル削除: count=%d", deletedCount)
```

## 実際のプロダクションでの応用例

### 1. AWS S3の署名付きURL

```go
// S3の署名付きURLも同じ概念
presignedURL, err := s3Client.PresignGetObject(&s3.GetObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("report.pdf"),
}, func(opts *s3.PresignOptions) {
    opts.Expires = 24 * time.Hour // expires_at相当
})
```

### 2. JWTトークン

```go
// JWTのexpクレームも同じ概念
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": "user123",
    "exp":     time.Now().Add(24 * time.Hour).Unix(), // expires_at
})
```

### 3. セッション管理

```go
// セッションの有効期限
session := &Session{
    ID:        generateID(),
    UserID:    userID,
    CreatedAt: time.Now(),
    ExpiresAt: time.Now().Add(30 * time.Minute), // expires_at
}
```

## テスト

### ユニットテスト例

```go
func TestTemporaryFileStorage_ExpiresAt(t *testing.T) {
    storage, _ := NewTemporaryFileStorage("/tmp/test", "secret", 1*time.Second)
    
    // ファイルを保存
    token, metadata, _ := storage.SaveFile("test.pdf", []byte("test"))
    
    // すぐにアクセス - 成功
    data, _, err := storage.GetFile(token)
    assert.NoError(t, err)
    assert.Equal(t, []byte("test"), data)
    
    // 2秒待つ（有効期限切れ）
    time.Sleep(2 * time.Second)
    
    // アクセス - 失敗
    _, _, err = storage.GetFile(token)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "有効期限")
}
```

## まとめ

`expires_at`の実装により：

✅ **セキュリティ向上**
- 一時的なアクセスのみ許可
- 改ざん防止（HMAC署名）
- 自動クリーンアップ

✅ **リソース管理**
- ディスク容量の節約
- 古いファイルの自動削除

✅ **ユーザー体験**
- 明確な有効期限の提示
- 適切なエラーメッセージ

この実装パターンは、ファイルダウンロード以外にも、API トークン、セッション管理、キャッシュなど、様々な場面で応用できます。
