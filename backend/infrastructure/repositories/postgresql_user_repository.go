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

// PostgreSQLUserRepository はPostgreSQLを使用したユーザーリポジトリの実装
type PostgreSQLUserRepository struct {
	db *sql.DB
}

// NewPostgreSQLUserRepository は新しいPostgreSQLユーザーリポジトリを作成する
func NewPostgreSQLUserRepository(db *sql.DB) repositories.UserRepository {
	return &PostgreSQLUserRepository{db: db}
}

// Save は新しいユーザーを保存する
func (r *PostgreSQLUserRepository) Save(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, provider, provider_user_id, name, avatar_url, email_verified, email_verified_at, two_factor_enabled, two_factor_secret, two_factor_backup_codes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	var passwordHash *string
	if user.PasswordHash().String() != "" {
		pwdHash := user.PasswordHash().String()
		passwordHash = &pwdHash
	}

	var providerUserID *string
	if user.ProviderUserID() != "" {
		pid := user.ProviderUserID()
		providerUserID = &pid
	}

	var name *string
	if user.Name() != "" {
		n := user.Name()
		name = &n
	}

	var avatarURL *string
	if user.AvatarURL() != "" {
		au := user.AvatarURL()
		avatarURL = &au
	}

	var twoFactorSecret *string
	if user.TwoFactorSecret() != "" {
		tfs := user.TwoFactorSecret()
		twoFactorSecret = &tfs
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID().String(),
		user.Email().String(),
		passwordHash,
		string(user.Provider()),
		providerUserID,
		name,
		avatarURL,
		user.EmailVerified(),
		user.EmailVerifiedAt(),
		user.TwoFactorEnabled(),
		twoFactorSecret,
		pq.Array(user.TwoFactorBackupCodes()),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("ユーザーの保存に失敗しました: %w", err)
	}

	return nil
}

// FindByID は指定されたIDのユーザーを取得する
func (r *PostgreSQLUserRepository) FindByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	var userID, email string
	var passwordHash, provider, providerUserID, name, avatarURL, twoFactorSecret sql.NullString
	var emailVerified, twoFactorEnabled bool
	var emailVerifiedAt sql.NullTime
	var twoFactorBackupCodes []string
	var createdAt, updatedAt time.Time

	query := `SELECT id, email, password_hash, provider, provider_user_id, name, avatar_url, email_verified, email_verified_at, two_factor_enabled, two_factor_secret, two_factor_backup_codes, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&userID, &email, &passwordHash, &provider, &providerUserID, &name, &avatarURL, &emailVerified, &emailVerifiedAt, &twoFactorEnabled, &twoFactorSecret, pq.Array(&twoFactorBackupCodes), &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ユーザーが見つかりません: %s", id)
		}
		return nil, fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
	}

	var emailVerifiedAtPtr *time.Time
	if emailVerifiedAt.Valid {
		emailVerifiedAtPtr = &emailVerifiedAt.Time
	}

	return entities.ReconstructUserWithOAuth(
		userID,
		email,
		passwordHash.String,
		provider.String,
		providerUserID.String,
		name.String,
		avatarURL.String,
		emailVerified,
		emailVerifiedAtPtr,
		twoFactorEnabled,
		twoFactorSecret.String,
		twoFactorBackupCodes,
		createdAt,
		updatedAt,
	)
}

// FindByEmail はメールアドレスからユーザーを取得する
func (r *PostgreSQLUserRepository) FindByEmail(ctx context.Context, email entities.Email) (*entities.User, error) {
	var userID, emailStr string
	var passwordHash, provider, providerUserID, name, avatarURL, twoFactorSecret sql.NullString
	var emailVerified, twoFactorEnabled bool
	var emailVerifiedAt sql.NullTime
	var twoFactorBackupCodes []string
	var createdAt, updatedAt time.Time

	query := `SELECT id, email, password_hash, provider, provider_user_id, name, avatar_url, email_verified, email_verified_at, two_factor_enabled, two_factor_secret, two_factor_backup_codes, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&userID, &emailStr, &passwordHash, &provider, &providerUserID, &name, &avatarURL, &emailVerified, &emailVerifiedAt, &twoFactorEnabled, &twoFactorSecret, pq.Array(&twoFactorBackupCodes), &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ユーザーが見つかりません: %s", email)
		}
		return nil, fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
	}

	var emailVerifiedAtPtr *time.Time
	if emailVerifiedAt.Valid {
		emailVerifiedAtPtr = &emailVerifiedAt.Time
	}

	return entities.ReconstructUserWithOAuth(
		userID,
		emailStr,
		passwordHash.String,
		provider.String,
		providerUserID.String,
		name.String,
		avatarURL.String,
		emailVerified,
		emailVerifiedAtPtr,
		twoFactorEnabled,
		twoFactorSecret.String,
		twoFactorBackupCodes,
		createdAt,
		updatedAt,
	)
}

// Update は既存のユーザー情報を更新する
func (r *PostgreSQLUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, two_factor_enabled = $3, two_factor_secret = $4, two_factor_backup_codes = $5, updated_at = $6
		WHERE id = $7`

	var twoFactorSecret *string
	if user.TwoFactorSecret() != "" {
		tfs := user.TwoFactorSecret()
		twoFactorSecret = &tfs
	}

	result, err := r.db.ExecContext(ctx, query,
		user.Email().String(),
		user.PasswordHash().String(),
		user.TwoFactorEnabled(),
		twoFactorSecret,
		pq.Array(user.TwoFactorBackupCodes()),
		user.UpdatedAt(),
		user.ID().String(),
	)
	if err != nil {
		return fmt.Errorf("ユーザーの更新に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("更新結果の確認に失敗しました: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ユーザーが見つかりません: %s", user.ID())
	}

	return nil
}

// Delete は指定されたIDのユーザーを削除する
func (r *PostgreSQLUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("ユーザーの削除に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果の確認に失敗しました: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ユーザーが見つかりません: %s", id)
	}

	return nil
}

// Exists は指定されたIDのユーザーが存在するか確認する
func (r *PostgreSQLUserRepository) Exists(ctx context.Context, id entities.UserID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ユーザーの存在確認に失敗しました: %w", err)
	}

	return exists, nil
}

// ExistsByEmail はメールアドレスが既に使用されているか確認する
func (r *PostgreSQLUserRepository) ExistsByEmail(ctx context.Context, email entities.Email) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("メールアドレスの存在確認に失敗しました: %w", err)
	}

	return exists, nil
}

// FindByProviderUserID はOAuthプロバイダーのユーザーIDからユーザーを取得する
func (r *PostgreSQLUserRepository) FindByProviderUserID(ctx context.Context, provider entities.AuthProvider, providerUserID string) (*entities.User, error) {
	var userID, email string
	var passwordHash, providerStr, providerUID, name, avatarURL, twoFactorSecret sql.NullString
	var emailVerified, twoFactorEnabled bool
	var emailVerifiedAt sql.NullTime
	var twoFactorBackupCodes []string
	var createdAt, updatedAt time.Time

	query := `SELECT id, email, password_hash, provider, provider_user_id, name, avatar_url, email_verified, email_verified_at, two_factor_enabled, two_factor_secret, two_factor_backup_codes, created_at, updated_at 
			  FROM users 
			  WHERE provider = $1 AND provider_user_id = $2`
	err := r.db.QueryRowContext(ctx, query, string(provider), providerUserID).Scan(
		&userID, &email, &passwordHash, &providerStr, &providerUID, &name, &avatarURL, &emailVerified, &emailVerifiedAt, &twoFactorEnabled, &twoFactorSecret, pq.Array(&twoFactorBackupCodes), &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ユーザーが見つかりません: provider=%s, providerUserID=%s", provider, providerUserID)
		}
		return nil, fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
	}

	var emailVerifiedAtPtr *time.Time
	if emailVerifiedAt.Valid {
		emailVerifiedAtPtr = &emailVerifiedAt.Time
	}

	return entities.ReconstructUserWithOAuth(
		userID,
		email,
		passwordHash.String,
		providerStr.String,
		providerUID.String,
		name.String,
		avatarURL.String,
		emailVerified,
		emailVerifiedAtPtr,
		twoFactorEnabled,
		twoFactorSecret.String,
		twoFactorBackupCodes,
		createdAt,
		updatedAt,
	)
}
