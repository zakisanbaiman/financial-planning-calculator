package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/lib/pq"
)

// PostgreSQLWebAuthnCredentialRepository はPostgreSQLを使用したWebAuthn認証情報リポジトリの実装
type PostgreSQLWebAuthnCredentialRepository struct {
	db *sql.DB
}

// NewPostgreSQLWebAuthnCredentialRepository は新しいPostgreSQLWebAuthn認証情報リポジトリを作成する
func NewPostgreSQLWebAuthnCredentialRepository(db *sql.DB) repositories.WebAuthnCredentialRepository {
	return &PostgreSQLWebAuthnCredentialRepository{db: db}
}

// Save は新しいWebAuthn認証情報を保存する
func (r *PostgreSQLWebAuthnCredentialRepository) Save(ctx context.Context, credential *entities.WebAuthnCredential) error {
	query := `
		INSERT INTO webauthn_credentials (
			id, user_id, credential_id, public_key, attestation_type, aaguid, 
			sign_count, clone_warning, transports, name, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	var name *string
	if credential.Name() != "" {
		n := credential.Name()
		name = &n
	}

	_, err := r.db.ExecContext(ctx, query,
		credential.ID().String(),
		credential.UserID().String(),
		credential.CredentialID(),
		credential.PublicKey(),
		credential.AttestationType(),
		credential.AAGUID(),
		credential.SignCount(),
		credential.CloneWarning(),
		pq.Array(credential.Transports()),
		name,
		credential.CreatedAt(),
		credential.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("WebAuthn認証情報の保存に失敗しました: %w", err)
	}

	return nil
}

// FindByID は指定されたIDの認証情報を取得する
func (r *PostgreSQLWebAuthnCredentialRepository) FindByID(ctx context.Context, id entities.CredentialID) (*entities.WebAuthnCredential, error) {
	var credID, userID string
	var credentialID, publicKey, aaguid []byte
	var attestationType string
	var signCount int
	var cloneWarning bool
	var transports []string
	var name sql.NullString
	var createdAt, updatedAt sql.NullTime
	var lastUsedAt sql.NullTime

	query := `
		SELECT id, user_id, credential_id, public_key, attestation_type, aaguid, 
		       sign_count, clone_warning, transports, name, created_at, updated_at, last_used_at 
		FROM webauthn_credentials WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&credID, &userID, &credentialID, &publicKey, &attestationType, &aaguid,
		&signCount, &cloneWarning, pq.Array(&transports), &name, &createdAt, &updatedAt, &lastUsedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("WebAuthn認証情報が見つかりません: %s", id)
		}
		return nil, fmt.Errorf("WebAuthn認証情報の取得に失敗しました: %w", err)
	}

	var lastUsedAtPtr *sql.NullTime
	if lastUsedAt.Valid {
		lastUsedAtPtr = &lastUsedAt
	}

	return entities.ReconstructWebAuthnCredential(
		credID,
		userID,
		credentialID,
		publicKey,
		attestationType,
		aaguid,
		uint32(signCount),
		cloneWarning,
		transports,
		name.String,
		createdAt.Time,
		updatedAt.Time,
		getTimePtr(lastUsedAtPtr),
	)
}

// FindByCredentialID はWebAuthn credential IDから認証情報を取得する
func (r *PostgreSQLWebAuthnCredentialRepository) FindByCredentialID(ctx context.Context, credentialID []byte) (*entities.WebAuthnCredential, error) {
	var credID, userID string
	var credIDBytes, publicKey, aaguid []byte
	var attestationType string
	var signCount int
	var cloneWarning bool
	var transports []string
	var name sql.NullString
	var createdAt, updatedAt sql.NullTime
	var lastUsedAt sql.NullTime

	query := `
		SELECT id, user_id, credential_id, public_key, attestation_type, aaguid, 
		       sign_count, clone_warning, transports, name, created_at, updated_at, last_used_at 
		FROM webauthn_credentials WHERE credential_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, credentialID).Scan(
		&credID, &userID, &credIDBytes, &publicKey, &attestationType, &aaguid,
		&signCount, &cloneWarning, pq.Array(&transports), &name, &createdAt, &updatedAt, &lastUsedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("WebAuthn認証情報が見つかりません")
		}
		return nil, fmt.Errorf("WebAuthn認証情報の取得に失敗しました: %w", err)
	}

	var lastUsedAtPtr *sql.NullTime
	if lastUsedAt.Valid {
		lastUsedAtPtr = &lastUsedAt
	}

	return entities.ReconstructWebAuthnCredential(
		credID,
		userID,
		credIDBytes,
		publicKey,
		attestationType,
		aaguid,
		uint32(signCount),
		cloneWarning,
		transports,
		name.String,
		createdAt.Time,
		updatedAt.Time,
		getTimePtr(lastUsedAtPtr),
	)
}

// FindByUserID は指定されたユーザーIDの全ての認証情報を取得する
func (r *PostgreSQLWebAuthnCredentialRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.WebAuthnCredential, error) {
	query := `
		SELECT id, user_id, credential_id, public_key, attestation_type, aaguid, 
		       sign_count, clone_warning, transports, name, created_at, updated_at, last_used_at 
		FROM webauthn_credentials WHERE user_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("WebAuthn認証情報の取得に失敗しました: %w", err)
	}
	defer rows.Close()

	var credentials []*entities.WebAuthnCredential

	for rows.Next() {
		var credID, uid string
		var credentialID, publicKey, aaguid []byte
		var attestationType string
		var signCount int
		var cloneWarning bool
		var transports []string
		var name sql.NullString
		var createdAt, updatedAt sql.NullTime
		var lastUsedAt sql.NullTime

		err := rows.Scan(
			&credID, &uid, &credentialID, &publicKey, &attestationType, &aaguid,
			&signCount, &cloneWarning, pq.Array(&transports), &name, &createdAt, &updatedAt, &lastUsedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("WebAuthn認証情報のスキャンに失敗しました: %w", err)
		}

		var lastUsedAtPtr *sql.NullTime
		if lastUsedAt.Valid {
			lastUsedAtPtr = &lastUsedAt
		}

		credential, err := entities.ReconstructWebAuthnCredential(
			credID,
			uid,
			credentialID,
			publicKey,
			attestationType,
			aaguid,
			uint32(signCount),
			cloneWarning,
			transports,
			name.String,
			createdAt.Time,
			updatedAt.Time,
			getTimePtr(lastUsedAtPtr),
		)
		if err != nil {
			return nil, fmt.Errorf("WebAuthn認証情報の再構築に失敗しました: %w", err)
		}

		credentials = append(credentials, credential)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("行の反復処理に失敗しました: %w", err)
	}

	return credentials, nil
}

// Update は既存の認証情報を更新する
func (r *PostgreSQLWebAuthnCredentialRepository) Update(ctx context.Context, credential *entities.WebAuthnCredential) error {
	query := `
		UPDATE webauthn_credentials 
		SET sign_count = $1, clone_warning = $2, name = $3, updated_at = $4, last_used_at = $5
		WHERE id = $6
	`

	var name *string
	if credential.Name() != "" {
		n := credential.Name()
		name = &n
	}

	_, err := r.db.ExecContext(ctx, query,
		credential.SignCount(),
		credential.CloneWarning(),
		name,
		credential.UpdatedAt(),
		credential.LastUsedAt(),
		credential.ID().String(),
	)
	if err != nil {
		return fmt.Errorf("WebAuthn認証情報の更新に失敗しました: %w", err)
	}

	return nil
}

// Delete は指定されたIDの認証情報を削除する
func (r *PostgreSQLWebAuthnCredentialRepository) Delete(ctx context.Context, id entities.CredentialID) error {
	query := `DELETE FROM webauthn_credentials WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("WebAuthn認証情報の削除に失敗しました: %w", err)
	}

	return nil
}

// getTimePtr はsql.NullTimeから*time.Timeに変換する
func getTimePtr(nt *sql.NullTime) *time.Time {
	if nt == nil || !nt.Valid {
		return nil
	}
	t := nt.Time
	return &t
}
