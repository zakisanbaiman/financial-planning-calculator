package faq_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/financial-planning-calculator/backend/infrastructure/faq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===========================
// テストヘルパー
// ===========================

// setupFAQDir はテスト用の一時FAQディレクトリを作成し、クリーンアップ関数を返す
func setupFAQDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		err := os.WriteFile(path, []byte(content), 0o644)
		require.NoError(t, err)
	}
	return dir
}

// ===========================
// FAQLoader.Load Tests
// ===========================

func TestFAQLoader_Load(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: Markdownファイルを読み込んでFAQContentを返す", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md": "# 積立NISAとは\n\n積立NISAは非課税投資制度です。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)

		// When
		contents, err := loader.Load(ctx)

		// Then
		require.NoError(t, err)
		require.Len(t, contents, 1)
		assert.Equal(t, "nisa.md", contents[0].Filename)
		assert.Equal(t, "積立NISAとは", contents[0].Title)
		assert.Contains(t, contents[0].Body, "非課税投資制度")
	})

	t.Run("正常系: 複数のMarkdownファイルをすべて読み込む", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md":      "# 積立NISAとは\n\n積立NISAの説明です。",
			"ideco.md":     "# iDeCoとは\n\n個人型確定拠出年金の説明です。",
			"insurance.md": "# 生命保険とは\n\n生命保険の説明です。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)

		// When
		contents, err := loader.Load(ctx)

		// Then
		require.NoError(t, err)
		assert.Len(t, contents, 3)
	})

	t.Run("正常系: 非Markdownファイルは無視される", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md":   "# 積立NISAとは\n\nNISAの説明です。",
			"readme.txt": "テキストファイルは無視されます",
			"data.json": `{"key": "value"}`,
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)

		// When
		contents, err := loader.Load(ctx)

		// Then
		require.NoError(t, err)
		assert.Len(t, contents, 1)
		assert.Equal(t, "nisa.md", contents[0].Filename)
	})

	t.Run("正常系: FAQディレクトリが空の場合は空スライスを返す", func(t *testing.T) {
		// Given
		dir := setupFAQDir(t, map[string]string{})
		loader := faq.NewFAQLoader(dir)

		// When
		contents, err := loader.Load(ctx)

		// Then
		require.NoError(t, err)
		assert.Empty(t, contents)
	})

	t.Run("正常系: h1見出しがない場合はファイル名をタイトルとして使用する", func(t *testing.T) {
		// Given
		files := map[string]string{
			"overview.md": "見出しなしの本文です。\n\n詳細内容です。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)

		// When
		contents, err := loader.Load(ctx)

		// Then
		require.NoError(t, err)
		require.Len(t, contents, 1)
		assert.Equal(t, "overview", contents[0].Title)
		assert.Contains(t, contents[0].Body, "見出しなしの本文")
	})

	t.Run("正常系: FAQディレクトリ内にサブディレクトリがあってもスキップされる", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md": "# 積立NISAとは\n\nNISAの説明です。",
		}
		dir := setupFAQDir(t, files)

		// サブディレクトリを作成してMarkdownファイルを置く
		subDir := filepath.Join(dir, "subdir")
		err := os.Mkdir(subDir, 0o755)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(subDir, "sub.md"), []byte("# サブディレクトリのFAQ\n\n内容"), 0o644)
		require.NoError(t, err)

		loader := faq.NewFAQLoader(dir)

		// When
		contents, err := loader.Load(ctx)

		// Then
		require.NoError(t, err)
		assert.Len(t, contents, 1, "サブディレクトリ内のファイルは無視されること")
		assert.Equal(t, "nisa.md", contents[0].Filename)
	})

	t.Run("異常系: 存在しないディレクトリを指定した場合はエラーを返す", func(t *testing.T) {
		// Given
		loader := faq.NewFAQLoader("/nonexistent/path/to/faq")

		// When
		_, err := loader.Load(ctx)

		// Then
		require.Error(t, err)
	})

	t.Run("異常系: コンテキストがキャンセルされた場合はエラーを返す", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md": "# 積立NISAとは\n\nNISAの説明です。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 即座にキャンセル

		// When
		_, err := loader.Load(ctx)

		// Then
		require.Error(t, err)
	})
}

// ===========================
// FAQLoader.Search Tests
// ===========================

func TestFAQLoader_Search(t *testing.T) {
	t.Run("正常系: クエリに一致するFAQを返す", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md":  "# 積立NISAとは\n\n積立NISAは非課税投資制度です。",
			"ideco.md": "# iDeCoとは\n\n個人型確定拠出年金の説明です。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)
		ctx := context.Background()
		_, err := loader.Load(ctx)
		require.NoError(t, err)

		// When
		results := loader.Search("NISA")

		// Then
		assert.NotEmpty(t, results)
		found := false
		for _, r := range results {
			if r.Filename == "nisa.md" {
				found = true
				break
			}
		}
		assert.True(t, found, "nisa.mdがSearch結果に含まれていること")
	})

	t.Run("正常系: 本文に一致するキーワードでも検索できる", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md": "# 投資について\n\n積立NISAは非課税投資制度です。年間120万円まで投資できます。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)
		ctx := context.Background()
		_, err := loader.Load(ctx)
		require.NoError(t, err)

		// When
		results := loader.Search("非課税")

		// Then
		assert.NotEmpty(t, results)
	})

	t.Run("正常系: 一致するFAQがない場合は空スライスを返す", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md": "# 積立NISAとは\n\nNISAの説明です。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)
		ctx := context.Background()
		_, err := loader.Load(ctx)
		require.NoError(t, err)

		// When
		results := loader.Search("ビットコイン")

		// Then
		assert.Empty(t, results)
	})

	t.Run("正常系: 空クエリの場合はすべてのFAQを返す", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md":  "# 積立NISAとは\n\nNISAの説明です。",
			"ideco.md": "# iDeCoとは\n\niDeCoの説明です。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)
		ctx := context.Background()
		_, err := loader.Load(ctx)
		require.NoError(t, err)

		// When
		results := loader.Search("")

		// Then
		assert.Len(t, results, 2)
	})

	t.Run("正常系: 大文字小文字を区別せずに検索できる", func(t *testing.T) {
		// Given
		files := map[string]string{
			"nisa.md": "# 積立NISAとは\n\nnisaの非課税制度です。",
		}
		dir := setupFAQDir(t, files)
		loader := faq.NewFAQLoader(dir)
		ctx := context.Background()
		_, err := loader.Load(ctx)
		require.NoError(t, err)

		// When
		resultsLower := loader.Search("nisa")
		resultsUpper := loader.Search("NISA")

		// Then
		assert.Equal(t, len(resultsLower), len(resultsUpper))
	})

	t.Run("エッジケース: Load前にSearchを呼ぶと空スライスを返す", func(t *testing.T) {
		// Given
		dir := t.TempDir()
		loader := faq.NewFAQLoader(dir)

		// When
		results := loader.Search("NISA")

		// Then
		assert.Empty(t, results)
	})
}
