package ports

import "context"

// FAQContent はFAQドキュメントの内容を表す
type FAQContent struct {
	Filename string
	Title    string
	Body     string
}

// FAQLoader はFAQを読み込むインタフェース
type FAQLoader interface {
	Load(ctx context.Context) ([]FAQContent, error)
	Search(query string) []FAQContent
}
