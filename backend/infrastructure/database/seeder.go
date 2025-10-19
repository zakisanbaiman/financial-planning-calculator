package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed seeds/*.sql
var seedFiles embed.FS

// Seeder handles database seeding
type Seeder struct {
	db *sql.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

// Seed loads all seed data into the database
func (s *Seeder) Seed() error {
	seedFiles, err := s.loadSeedFiles()
	if err != nil {
		return err
	}

	for _, seedFile := range seedFiles {
		log.Printf("シードファイル %s を実行中...", seedFile.Name)

		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
		}

		_, err = tx.Exec(seedFile.Content)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("シードファイル %s の実行に失敗しました: %w", seedFile.Name, err)
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("シードファイル %s のコミットに失敗しました: %w", seedFile.Name, err)
		}

		log.Printf("シードファイル %s が正常に実行されました", seedFile.Name)
	}

	log.Println("全てのシードファイルが正常に実行されました")
	return nil
}

// SeedFile represents a seed file
type SeedFile struct {
	Name    string
	Content string
}

// loadSeedFiles loads all seed files from the embedded filesystem
func (s *Seeder) loadSeedFiles() ([]*SeedFile, error) {
	var seedFileList []*SeedFile

	err := fs.WalkDir(seedFiles, "seeds", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		content, err := seedFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("シードファイル %s の読み取りに失敗しました: %w", path, err)
		}

		filename := filepath.Base(path)
		seedFile := &SeedFile{
			Name:    filename,
			Content: string(content),
		}

		seedFileList = append(seedFileList, seedFile)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("シードファイルの読み込みに失敗しました: %w", err)
	}

	// Sort seed files by name to ensure consistent execution order
	sort.Slice(seedFileList, func(i, j int) bool {
		return seedFileList[i].Name < seedFileList[j].Name
	})

	return seedFileList, nil
}
