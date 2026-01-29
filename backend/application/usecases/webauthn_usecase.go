package usecases

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// WebAuthnUseCase はWebAuthn（パスキー）認証関連のユースケース
type WebAuthnUseCase interface {
	// BeginRegistration はパスキー登録を開始する
	BeginRegistration(ctx context.Context, userID string) (*BeginRegistrationOutput, error)

	// FinishRegistration はパスキー登録を完了する
	FinishRegistration(ctx context.Context, input FinishRegistrationInput) error

	// BeginLogin はパスキーログインを開始する
	BeginLogin(ctx context.Context, input BeginLoginInput) (*BeginLoginOutput, error)

	// FinishLogin はパスキーログインを完了する
	FinishLogin(ctx context.Context, input FinishLoginInput) (*LoginOutput, error)

	// ListCredentials はユーザーの全パスキーを取得する
	ListCredentials(ctx context.Context, userID string) ([]*CredentialInfo, error)

	// DeleteCredential はパスキーを削除する
	DeleteCredential(ctx context.Context, userID string, credentialID string) error

	// RenameCredential はパスキーの名前を変更する
	RenameCredential(ctx context.Context, userID string, credentialID string, newName string) error
}

// BeginRegistrationOutput はパスキー登録開始の出力
type BeginRegistrationOutput struct {
	PublicKeyOptions string `json:"publicKey"` // JSON形式のCredentialCreationOptions
	SessionData      string `json:"sessionData"` // セッションデータ（次のステップで使用）
}

// FinishRegistrationInput はパスキー登録完了の入力
type FinishRegistrationInput struct {
	UserID          string `json:"user_id"`
	CredentialName  string `json:"credential_name"`
	Response        string `json:"response"` // JSON形式のAuthenticatorAttestationResponse
	SessionData     string `json:"session_data"`
}

// BeginLoginInput はパスキーログイン開始の入力
type BeginLoginInput struct {
	Email string `json:"email"` // オプション：ユーザー特定用
}

// BeginLoginOutput はパスキーログイン開始の出力
type BeginLoginOutput struct {
	PublicKeyOptions string `json:"publicKey"` // JSON形式のCredentialRequestOptions
	SessionData      string `json:"sessionData"` // セッションデータ（次のステップで使用）
}

// FinishLoginInput はパスキーログイン完了の入力
type FinishLoginInput struct {
	Response    string `json:"response"` // JSON形式のAuthenticatorAssertionResponse
	SessionData string `json:"session_data"`
}

// CredentialInfo はパスキー情報
type CredentialInfo struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	CreatedAt  string  `json:"created_at"`
	LastUsedAt *string `json:"last_used_at,omitempty"`
}

// webAuthnUseCaseImpl はWebAuthnUseCaseの実装
type webAuthnUseCaseImpl struct {
	userRepo             repositories.UserRepository
	credentialRepo       repositories.WebAuthnCredentialRepository
	refreshTokenRepo     repositories.RefreshTokenRepository
	webAuthn             *webauthn.WebAuthn
	authUseCase          AuthUseCase
	jwtSecret            string
	jwtExpiration        time.Duration
	refreshTokenExpiration time.Duration
}

// NewWebAuthnUseCase は新しいWebAuthnUseCaseを作成する
func NewWebAuthnUseCase(
	userRepo repositories.UserRepository,
	credentialRepo repositories.WebAuthnCredentialRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	webAuthn *webauthn.WebAuthn,
	authUseCase AuthUseCase,
	jwtSecret string,
	jwtExpiration time.Duration,
	refreshTokenExpiration time.Duration,
) WebAuthnUseCase {
	return &webAuthnUseCaseImpl{
		userRepo:               userRepo,
		credentialRepo:         credentialRepo,
		refreshTokenRepo:       refreshTokenRepo,
		webAuthn:               webAuthn,
		authUseCase:            authUseCase,
		jwtSecret:              jwtSecret,
		jwtExpiration:          jwtExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

// webAuthnUser はWebAuthnのUserインターフェースを実装する
type webAuthnUser struct {
	id          []byte
	name        string
	displayName string
	credentials []webauthn.Credential
}

func (u *webAuthnUser) WebAuthnID() []byte {
	return u.id
}

func (u *webAuthnUser) WebAuthnName() string {
	return u.name
}

func (u *webAuthnUser) WebAuthnDisplayName() string {
	return u.displayName
}

func (u *webAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func (u *webAuthnUser) WebAuthnIcon() string {
	return ""
}

// BeginRegistration はパスキー登録を開始する
func (uc *webAuthnUseCaseImpl) BeginRegistration(ctx context.Context, userID string) (*BeginRegistrationOutput, error) {
	// ユーザーを取得
	uid, err := entities.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("無効なユーザーID: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// 既存のクレデンシャルを取得
	existingCreds, err := uc.credentialRepo.FindByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("既存のクレデンシャル取得に失敗: %w", err)
	}

	// WebAuthn用のクレデンシャルに変換
	webAuthnCreds := make([]webauthn.Credential, 0, len(existingCreds))
	for _, cred := range existingCreds {
		webAuthnCreds = append(webAuthnCreds, webauthn.Credential{
			ID:              cred.CredentialID(),
			PublicKey:       cred.PublicKey(),
			AttestationType: cred.AttestationType(),
			Transport:       convertTransports(cred.Transports()),
			Flags: webauthn.CredentialFlags{
				UserPresent:    true,
				UserVerified:   true,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:    cred.AAGUID(),
				SignCount: cred.SignCount(),
			},
		})
	}

	// WebAuthnユーザーを作成
	webUser := &webAuthnUser{
		id:          []byte(userID),
		name:        user.Email().String(),
		displayName: user.Name(),
		credentials: webAuthnCreds,
	}

	// 登録セッションを開始
	options, sessionData, err := uc.webAuthn.BeginRegistration(webUser)
	if err != nil {
		return nil, fmt.Errorf("パスキー登録の開始に失敗: %w", err)
	}

	// JSONにエンコード
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("オプションのエンコードに失敗: %w", err)
	}

	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		return nil, fmt.Errorf("セッションデータのエンコードに失敗: %w", err)
	}

	return &BeginRegistrationOutput{
		PublicKeyOptions: string(optionsJSON),
		SessionData:      base64.StdEncoding.EncodeToString(sessionDataJSON),
	}, nil
}

// FinishRegistration はパスキー登録を完了する
func (uc *webAuthnUseCaseImpl) FinishRegistration(ctx context.Context, input FinishRegistrationInput) error {
	// ユーザーを取得
	uid, err := entities.NewUserID(input.UserID)
	if err != nil {
		return fmt.Errorf("無効なユーザーID: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		return fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// セッションデータをデコード
	sessionDataBytes, err := base64.StdEncoding.DecodeString(input.SessionData)
	if err != nil {
		return fmt.Errorf("セッションデータのデコードに失敗: %w", err)
	}

	var sessionData webauthn.SessionData
	if err := json.Unmarshal(sessionDataBytes, &sessionData); err != nil {
		return fmt.Errorf("セッションデータのパースに失敗: %w", err)
	}

	// WebAuthnユーザーを作成
	webUser := &webAuthnUser{
		id:          []byte(input.UserID),
		name:        user.Email().String(),
		displayName: user.Name(),
	}

	// レスポンスをパース
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader([]byte(input.Response)))
	if err != nil {
		return fmt.Errorf("認証レスポンスのパースに失敗: %w", err)
	}

	// 登録を完了
	credential, err := uc.webAuthn.CreateCredential(webUser, sessionData, parsedResponse)
	if err != nil {
		return fmt.Errorf("クレデンシャルの作成に失敗: %w", err)
	}

	// データベースに保存
	credID := uuid.New().String()
	webAuthnCred, err := entities.NewWebAuthnCredential(
		credID,
		uid,
		credential.ID,
		credential.PublicKey,
		credential.AttestationType,
		credential.Authenticator.AAGUID,
		convertTransportsToStrings(credential.Transport),
		input.CredentialName,
	)
	if err != nil {
		return fmt.Errorf("クレデンシャルエンティティの作成に失敗: %w", err)
	}

	if err := uc.credentialRepo.Save(ctx, webAuthnCred); err != nil {
		return fmt.Errorf("クレデンシャルの保存に失敗: %w", err)
	}

	return nil
}

// BeginLogin はパスキーログインを開始する
func (uc *webAuthnUseCaseImpl) BeginLogin(ctx context.Context, input BeginLoginInput) (*BeginLoginOutput, error) {
	// ログインセッションを開始（ユーザーレス認証をサポート）
	options, sessionData, err := uc.webAuthn.BeginDiscoverableLogin()
	if err != nil {
		return nil, fmt.Errorf("パスキーログインの開始に失敗: %w", err)
	}

	// JSONにエンコード
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("オプションのエンコードに失敗: %w", err)
	}

	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		return nil, fmt.Errorf("セッションデータのエンコードに失敗: %w", err)
	}

	return &BeginLoginOutput{
		PublicKeyOptions: string(optionsJSON),
		SessionData:      base64.StdEncoding.EncodeToString(sessionDataJSON),
	}, nil
}

// FinishLogin はパスキーログインを完了する
func (uc *webAuthnUseCaseImpl) FinishLogin(ctx context.Context, input FinishLoginInput) (*LoginOutput, error) {
	// セッションデータをデコード
	sessionDataBytes, err := base64.StdEncoding.DecodeString(input.SessionData)
	if err != nil {
		return nil, fmt.Errorf("セッションデータのデコードに失敗: %w", err)
	}

	var sessionData webauthn.SessionData
	if err := json.Unmarshal(sessionDataBytes, &sessionData); err != nil {
		return nil, fmt.Errorf("セッションデータのパースに失敗: %w", err)
	}

	// レスポンスをパース
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader([]byte(input.Response)))
	if err != nil {
		return nil, fmt.Errorf("認証レスポンスのパースに失敗: %w", err)
	}

	// クレデンシャルを取得
	credential, err := uc.credentialRepo.FindByCredentialID(ctx, parsedResponse.RawID)
	if err != nil {
		return nil, fmt.Errorf("クレデンシャルが見つかりません: %w", err)
	}

	// ユーザーを取得
	user, err := uc.userRepo.FindByID(ctx, credential.UserID())
	if err != nil {
		return nil, fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// WebAuthnユーザーを作成
	webUser := &webAuthnUser{
		id:          []byte(credential.UserID().String()),
		name:        user.Email().String(),
		displayName: user.Name(),
		credentials: []webauthn.Credential{
			{
				ID:              credential.CredentialID(),
				PublicKey:       credential.PublicKey(),
				AttestationType: credential.AttestationType(),
				Transport:       convertTransports(credential.Transports()),
				Flags: webauthn.CredentialFlags{
					UserPresent:  true,
					UserVerified: true,
				},
				Authenticator: webauthn.Authenticator{
					AAGUID:    credential.AAGUID(),
					SignCount: credential.SignCount(),
				},
			},
		},
	}

	// ログイン検証
	validatedCredential, err := uc.webAuthn.ValidateLogin(webUser, sessionData, parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("パスキー認証に失敗: %w", err)
	}

	// サインカウントを更新
	if err := credential.UpdateSignCount(validatedCredential.Authenticator.SignCount); err != nil {
		return nil, fmt.Errorf("サインカウントの更新に失敗: %w", err)
	}

	if err := uc.credentialRepo.Update(ctx, credential); err != nil {
		return nil, fmt.Errorf("クレデンシャルの更新に失敗: %w", err)
	}

	// JWTトークンとリフレッシュトークンを生成
	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("トークンの生成に失敗: %w", err)
	}

	refreshToken, err := uc.generateRefreshToken(ctx, user.ID())
	if err != nil {
		return nil, fmt.Errorf("リフレッシュトークンの生成に失敗: %w", err)
	}

	return &LoginOutput{
		UserID:       user.ID().String(),
		Email:        user.Email().String(),
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
	}, nil
}

// ListCredentials はユーザーの全パスキーを取得する
func (uc *webAuthnUseCaseImpl) ListCredentials(ctx context.Context, userID string) ([]*CredentialInfo, error) {
	uid, err := entities.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("無効なユーザーID: %w", err)
	}

	credentials, err := uc.credentialRepo.FindByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("クレデンシャルの取得に失敗: %w", err)
	}

	result := make([]*CredentialInfo, 0, len(credentials))
	for _, cred := range credentials {
		var lastUsedAt *string
		if cred.LastUsedAt() != nil {
			lu := cred.LastUsedAt().Format(time.RFC3339)
			lastUsedAt = &lu
		}

		result = append(result, &CredentialInfo{
			ID:         cred.ID().String(),
			Name:       cred.Name(),
			CreatedAt:  cred.CreatedAt().Format(time.RFC3339),
			LastUsedAt: lastUsedAt,
		})
	}

	return result, nil
}

// DeleteCredential はパスキーを削除する
func (uc *webAuthnUseCaseImpl) DeleteCredential(ctx context.Context, userID string, credentialID string) error {
	uid, err := entities.NewUserID(userID)
	if err != nil {
		return fmt.Errorf("無効なユーザーID: %w", err)
	}

	cid, err := entities.NewCredentialID(credentialID)
	if err != nil {
		return fmt.Errorf("無効なクレデンシャルID: %w", err)
	}

	// クレデンシャルを取得して、所有者を確認
	credential, err := uc.credentialRepo.FindByID(ctx, cid)
	if err != nil {
		return fmt.Errorf("クレデンシャルが見つかりません: %w", err)
	}

	if credential.UserID() != uid {
		return fmt.Errorf("このクレデンシャルの所有者ではありません")
	}

	if err := uc.credentialRepo.Delete(ctx, cid); err != nil {
		return fmt.Errorf("クレデンシャルの削除に失敗: %w", err)
	}

	return nil
}

// RenameCredential はパスキーの名前を変更する
func (uc *webAuthnUseCaseImpl) RenameCredential(ctx context.Context, userID string, credentialID string, newName string) error {
	uid, err := entities.NewUserID(userID)
	if err != nil {
		return fmt.Errorf("無効なユーザーID: %w", err)
	}

	cid, err := entities.NewCredentialID(credentialID)
	if err != nil {
		return fmt.Errorf("無効なクレデンシャルID: %w", err)
	}

	// クレデンシャルを取得して、所有者を確認
	credential, err := uc.credentialRepo.FindByID(ctx, cid)
	if err != nil {
		return fmt.Errorf("クレデンシャルが見つかりません: %w", err)
	}

	if credential.UserID() != uid {
		return fmt.Errorf("このクレデンシャルの所有者ではありません")
	}

	credential.UpdateName(newName)

	if err := uc.credentialRepo.Update(ctx, credential); err != nil {
		return fmt.Errorf("クレデンシャルの更新に失敗: %w", err)
	}

	return nil
}

// convertTransports はstring配列をprotocol.AuthenticatorTransport配列に変換する
func convertTransports(transports []string) []protocol.AuthenticatorTransport {
	result := make([]protocol.AuthenticatorTransport, 0, len(transports))
	for _, t := range transports {
		result = append(result, protocol.AuthenticatorTransport(t))
	}
	return result
}

// convertTransportsToStrings はprotocol.AuthenticatorTransport配列をstring配列に変換する
func convertTransportsToStrings(transports []protocol.AuthenticatorTransport) []string {
	result := make([]string, 0, len(transports))
	for _, t := range transports {
		result = append(result, string(t))
	}
	return result
}

// generateToken はJWTトークンを生成する
func (uc *webAuthnUseCaseImpl) generateToken(user *entities.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(uc.jwtExpiration)

	claims := TokenClaims{
		UserID: user.ID().String(),
		Email:  user.Email().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// generateRefreshToken はリフレッシュトークンを生成する
func (uc *webAuthnUseCaseImpl) generateRefreshToken(ctx context.Context, userID entities.UserID) (string, error) {
	// 有効期限を設定
	expiresAt := time.Now().Add(uc.refreshTokenExpiration)

	// リフレッシュトークンエンティティを作成
	refreshTokenEntity, rawToken, err := entities.NewRefreshToken(userID, expiresAt)
	if err != nil {
		return "", fmt.Errorf("リフレッシュトークンエンティティの作成に失敗しました: %w", err)
	}

	// データベースに保存
	if err := uc.refreshTokenRepo.Save(ctx, refreshTokenEntity); err != nil {
		return "", fmt.Errorf("リフレッシュトークンの保存に失敗しました: %w", err)
	}

	// 元のトークン（ハッシュ化前）を返す
	return rawToken, nil
}
