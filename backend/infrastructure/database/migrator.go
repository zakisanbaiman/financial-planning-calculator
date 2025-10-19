package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"embed"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

// Migrator handles database migrations
type Migrator struct {
	db *sql.DB
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// createMigrationsTable creates the migrations tracking table if it doesn't exist
func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("マイグレーションテーブルの作成に失敗しました: %w", err)
	}
	return nil
}

// getAppliedMigrations returns a list of applied migrations
func (m *Migrator) getAppliedMigrations() (map[string]*Migration, error) {
	if err := m.createMigrationsTable(); err != nil {
		return nil, err
	}

	query := `SELECT version, name, applied_at FROM schema_migrations ORDER BY version`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("適用済みマイグレーションの取得に失敗しました: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]*Migration)
	for rows.Next() {
		var migration Migration
		err := rows.Scan(&migration.Version, &migration.Name, &migration.AppliedAt)
		if err != nil {
			return nil, fmt.Errorf("マイグレーション情報の読み取りに失敗しました: %w", err)
		}
		applied[migration.Version] = &migration
	}

	return applied, nil
}

// loadMigrations loads all migration files from the embedded filesystem
func (m *Migrator) loadMigrations() ([]*Migration, error) {
	var migrations []*Migration

	err := fs.WalkDir(migrationFiles, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		// Skip down migrations for now, we'll load them separately
		if strings.HasSuffix(path, "_down.sql") {
			return nil
		}

		content, err := migrationFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("マイグレーションファイル %s の読み取りに失敗しました: %w", path, err)
		}

		filename := filepath.Base(path)
		parts := strings.SplitN(filename, "_", 2)
		if len(parts) < 2 {
			return fmt.Errorf("無効なマイグレーションファイル名: %s", filename)
		}

		version := parts[0]
		name := strings.TrimSuffix(parts[1], ".sql")

		// Try to load corresponding down migration
		downPath := strings.Replace(path, ".sql", "_down.sql", 1)
		var downSQL string
		if downContent, err := migrationFiles.ReadFile(downPath); err == nil {
			downSQL = string(downContent)
		}

		migration := &Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(content),
			DownSQL: downSQL,
		}

		migrations = append(migrations, migration)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("マイグレーションファイルの読み込みに失敗しました: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// Up applies all pending migrations
func (m *Migrator) Up() error {
	migrations, err := m.loadMigrations()
	if err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if _, exists := applied[migration.Version]; exists {
			log.Printf("マイグレーション %s は既に適用済みです", migration.Version)
			continue
		}

		log.Printf("マイグレーション %s を適用中...", migration.Version)

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
		}

		// Execute migration
		_, err = tx.Exec(migration.UpSQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("マイグレーション %s の実行に失敗しました: %w", migration.Version, err)
		}

		// Record migration
		_, err = tx.Exec(
			"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
			migration.Version, migration.Name,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("マイグレーション記録の保存に失敗しました: %w", err)
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("マイグレーションのコミットに失敗しました: %w", err)
		}

		log.Printf("マイグレーション %s が正常に適用されました", migration.Version)
	}

	log.Println("全てのマイグレーションが正常に適用されました")
	return nil
}

// Down rolls back the last applied migration
func (m *Migrator) Down() error {
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	if len(applied) == 0 {
		log.Println("ロールバックするマイグレーションがありません")
		return nil
	}

	// Find the latest migration
	var latest *Migration
	for _, migration := range applied {
		if latest == nil || migration.Version > latest.Version {
			latest = migration
		}
	}

	migrations, err := m.loadMigrations()
	if err != nil {
		return err
	}

	// Find the migration with down SQL
	var targetMigration *Migration
	for _, migration := range migrations {
		if migration.Version == latest.Version {
			targetMigration = migration
			break
		}
	}

	if targetMigration == nil || targetMigration.DownSQL == "" {
		return fmt.Errorf("マイグレーション %s のロールバックスクリプトが見つかりません", latest.Version)
	}

	log.Printf("マイグレーション %s をロールバック中...", latest.Version)

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
	}

	// Execute rollback
	_, err = tx.Exec(targetMigration.DownSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("マイグレーション %s のロールバックに失敗しました: %w", latest.Version, err)
	}

	// Remove migration record
	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = $1", latest.Version)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("マイグレーション記録の削除に失敗しました: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ロールバックのコミットに失敗しました: %w", err)
	}

	log.Printf("マイグレーション %s が正常にロールバックされました", latest.Version)
	return nil
}

// Status shows the current migration status
func (m *Migrator) Status() error {
	migrations, err := m.loadMigrations()
	if err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	log.Println("マイグレーション状況:")
	log.Println("バージョン\t状態\t\t名前")
	log.Println("--------\t----\t\t----")

	for _, migration := range migrations {
		status := "未適用"
		if appliedMigration, exists := applied[migration.Version]; exists {
			status = fmt.Sprintf("適用済み (%s)", appliedMigration.AppliedAt.Format("2006-01-02 15:04:05"))
		}
		log.Printf("%s\t%s\t%s", migration.Version, status, migration.Name)
	}

	return nil
}
