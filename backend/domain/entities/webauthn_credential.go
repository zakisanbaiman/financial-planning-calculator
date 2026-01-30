package entities

import (
	"errors"
	"time"
)

// CredentialID はWebAuthn認証情報の一意識別子
type CredentialID string

// NewCredentialID は既存のIDからCredentialIDを生成する
func NewCredentialID(id string) (CredentialID, error) {
	if id == "" {
		return "", errors.New("クレデンシャルIDは必須です")
	}
	return CredentialID(id), nil
}

// String はCredentialIDの文字列表現を返す
func (cid CredentialID) String() string {
	return string(cid)
}

// WebAuthnCredential はパスキー（WebAuthn）認証情報エンティティ
type WebAuthnCredential struct {
	id              CredentialID
	userID          UserID
	credentialID    []byte
	publicKey       []byte
	attestationType string
	aaguid          []byte
	signCount       uint32
	cloneWarning    bool
	transports      []string
	name            string
	createdAt       time.Time
	updatedAt       time.Time
	lastUsedAt      *time.Time
}

// NewWebAuthnCredential は新しいWebAuthn認証情報を作成する
func NewWebAuthnCredential(
	id string,
	userID UserID,
	credentialID []byte,
	publicKey []byte,
	attestationType string,
	aaguid []byte,
	transports []string,
	name string,
) (*WebAuthnCredential, error) {
	// バリデーション
	credID, err := NewCredentialID(id)
	if err != nil {
		return nil, err
	}

	if len(credentialID) == 0 {
		return nil, errors.New("WebAuthn credential IDは必須です")
	}

	if len(publicKey) == 0 {
		return nil, errors.New("公開鍵は必須です")
	}

	if attestationType == "" {
		return nil, errors.New("認証タイプは必須です")
	}

	now := time.Now()

	return &WebAuthnCredential{
		id:              credID,
		userID:          userID,
		credentialID:    credentialID,
		publicKey:       publicKey,
		attestationType: attestationType,
		aaguid:          aaguid,
		signCount:       0,
		cloneWarning:    false,
		transports:      transports,
		name:            name,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

// ReconstructWebAuthnCredential はDBから取得したデータからWebAuthnCredentialを再構築する
func ReconstructWebAuthnCredential(
	id string,
	userID string,
	credentialID []byte,
	publicKey []byte,
	attestationType string,
	aaguid []byte,
	signCount uint32,
	cloneWarning bool,
	transports []string,
	name string,
	createdAt, updatedAt time.Time,
	lastUsedAt *time.Time,
) (*WebAuthnCredential, error) {
	credID, err := NewCredentialID(id)
	if err != nil {
		return nil, err
	}

	uid, err := NewUserID(userID)
	if err != nil {
		return nil, err
	}

	return &WebAuthnCredential{
		id:              credID,
		userID:          uid,
		credentialID:    credentialID,
		publicKey:       publicKey,
		attestationType: attestationType,
		aaguid:          aaguid,
		signCount:       signCount,
		cloneWarning:    cloneWarning,
		transports:      transports,
		name:            name,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
		lastUsedAt:      lastUsedAt,
	}, nil
}

// Getters

// ID はクレデンシャルIDを返す
func (wc *WebAuthnCredential) ID() CredentialID {
	return wc.id
}

// UserID はユーザーIDを返す
func (wc *WebAuthnCredential) UserID() UserID {
	return wc.userID
}

// CredentialID はWebAuthn credential IDを返す
func (wc *WebAuthnCredential) CredentialID() []byte {
	return wc.credentialID
}

// PublicKey は公開鍵を返す
func (wc *WebAuthnCredential) PublicKey() []byte {
	return wc.publicKey
}

// AttestationType は認証タイプを返す
func (wc *WebAuthnCredential) AttestationType() string {
	return wc.attestationType
}

// AAGUID はAuthenticator AAGUIDを返す
func (wc *WebAuthnCredential) AAGUID() []byte {
	return wc.aaguid
}

// SignCount は署名カウンターを返す
func (wc *WebAuthnCredential) SignCount() uint32 {
	return wc.signCount
}

// CloneWarning はクローン警告フラグを返す
func (wc *WebAuthnCredential) CloneWarning() bool {
	return wc.cloneWarning
}

// Transports は対応トランスポートを返す
func (wc *WebAuthnCredential) Transports() []string {
	return wc.transports
}

// Name はクレデンシャル名を返す
func (wc *WebAuthnCredential) Name() string {
	return wc.name
}

// CreatedAt は作成日時を返す
func (wc *WebAuthnCredential) CreatedAt() time.Time {
	return wc.createdAt
}

// UpdatedAt は更新日時を返す
func (wc *WebAuthnCredential) UpdatedAt() time.Time {
	return wc.updatedAt
}

// LastUsedAt は最終使用日時を返す
func (wc *WebAuthnCredential) LastUsedAt() *time.Time {
	return wc.lastUsedAt
}

// UpdateSignCount は署名カウンターを更新する
func (wc *WebAuthnCredential) UpdateSignCount(newCount uint32) error {
	// クローン検出：カウンターが減少している場合は警告
	if newCount < wc.signCount {
		wc.cloneWarning = true
		return errors.New("クレデンシャルのクローンが検出されました")
	}

	wc.signCount = newCount
	now := time.Now()
	wc.lastUsedAt = &now
	wc.updatedAt = now

	return nil
}

// UpdateName はクレデンシャル名を更新する
func (wc *WebAuthnCredential) UpdateName(newName string) {
	wc.name = newName
	wc.updatedAt = time.Now()
}
