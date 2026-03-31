package faq

import (
	"bufio"
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/financial-planning-calculator/backend/application/ports"
)

const (
	// markdownExtension はMarkdownファイルの拡張子
	markdownExtension = ".md"
)

// faqLoader はports.FAQLoaderのファイルシステム実装
type faqLoader struct {
	dir      string
	mu       sync.RWMutex
	contents []ports.FAQContent
}

// NewFAQLoader は指定ディレクトリからFAQを読み込むローダーを作成する
func NewFAQLoader(dir string) ports.FAQLoader {
	return &faqLoader{
		dir: dir,
	}
}

// Load はFAQディレクトリからMarkdownファイルを読み込む
func (l *faqLoader) Load(ctx context.Context) ([]ports.FAQContent, error) {
	// コンテキストキャンセルチェック
	if err := ctx.Err(); err != nil {
		return nil, errors.New("コンテキストがキャンセルされました")
	}

	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return nil, err
	}

	var contents []ports.FAQContent
	for _, entry := range entries {
		// コンテキストキャンセルチェック（ファイル毎）
		if err := ctx.Err(); err != nil {
			return nil, errors.New("コンテキストがキャンセルされました")
		}

		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != markdownExtension {
			continue
		}

		content, err := l.parseMarkdownFile(filepath.Join(l.dir, name), name)
		if err != nil {
			slog.Warn("FAQファイルの読み込みに失敗しました", slog.String("file", name), slog.Any("error", err))
			continue
		}
		contents = append(contents, content)
	}

	if contents == nil {
		contents = []ports.FAQContent{}
	}

	l.mu.Lock()
	l.contents = contents
	l.mu.Unlock()

	return contents, nil
}

// Search はクエリに一致するFAQを返す
// queryが空の場合はすべてのFAQを返す
// 大文字小文字を区別しない
func (l *faqLoader) Search(query string) []ports.FAQContent {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.contents) == 0 {
		return []ports.FAQContent{}
	}

	if query == "" {
		result := make([]ports.FAQContent, len(l.contents))
		copy(result, l.contents)
		return result
	}

	queryLower := strings.ToLower(query)

	var results []ports.FAQContent
	for _, c := range l.contents {
		titleLower := strings.ToLower(c.Title)
		bodyLower := strings.ToLower(c.Body)

		if strings.Contains(titleLower, queryLower) || strings.Contains(bodyLower, queryLower) {
			results = append(results, c)
		}
	}

	if results == nil {
		return []ports.FAQContent{}
	}
	return results
}

// parseMarkdownFile はMarkdownファイルを解析してFAQContentを返す
func (l *faqLoader) parseMarkdownFile(path, filename string) (ports.FAQContent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ports.FAQContent{}, err
	}

	body := string(data)
	title := extractH1Title(body)

	if title == "" {
		// 拡張子なしのファイル名をタイトルとして使用
		title = strings.TrimSuffix(filename, markdownExtension)
	}

	return ports.FAQContent{
		Filename: filename,
		Title:    title,
		Body:     body,
	}, nil
}

// extractH1Title はMarkdownテキストからh1見出しを抽出する
func extractH1Title(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}
