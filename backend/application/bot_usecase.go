package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/financial-planning-calculator/backend/application/ports"
)

// BotChunk はBotからのストリームチャンクを表す
type BotChunk struct {
	Token string
	Done  bool
	Err   error
}

// BotUseCase はBot機能のユースケースインタフェース
type BotUseCase interface {
	StreamAnswer(ctx context.Context, question string) (<-chan BotChunk, error)
}

// botUseCase はBotUseCaseの実装
type botUseCase struct {
	faqLoader ports.FAQLoader
	llmClient ports.LLMClient
}

// NewBotUseCase はBotUseCaseを生成する
func NewBotUseCase(faqLoader ports.FAQLoader, llmClient ports.LLMClient) BotUseCase {
	return &botUseCase{
		faqLoader: faqLoader,
		llmClient: llmClient,
	}
}

// StreamAnswer は質問に対してFAQを検索し、LLMからの回答をストリームで返す
func (u *botUseCase) StreamAnswer(ctx context.Context, question string) (<-chan BotChunk, error) {
	if question == "" {
		ch := make(chan BotChunk, 1)
		go func() {
			defer close(ch)
			ch <- BotChunk{Err: errors.New("質問が空です")}
		}()
		return ch, nil
	}

	// FAQを検索してプロンプトを構築
	relatedFAQs := u.faqLoader.Search(question)
	prompt := buildPrompt(question, relatedFAQs)

	// LLMにストリームリクエストを送信
	llmCh, err := u.llmClient.StreamAnswer(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLMへのリクエストに失敗しました: %w", err)
	}

	// LLMのチャンクをBotChunkに変換して返す
	ch := make(chan BotChunk, 16)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case chunk, ok := <-llmCh:
				if !ok {
					return
				}
				if chunk.Error != nil {
					ch <- BotChunk{Err: chunk.Error}
					return
				}
				ch <- BotChunk{
					Token: chunk.Token,
					Done:  chunk.Done,
				}
				if chunk.Done {
					return
				}
			}
		}
	}()

	return ch, nil
}

// buildPrompt はFAQコンテキストと質問からプロンプトを構築する
func buildPrompt(question string, faqs []ports.FAQContent) string {
	var sb strings.Builder

	// サービス識別のシステム指示（常に含める）
	sb.WriteString("あなたはFinPlan（財務計画計算機）のAIアシスタントです。")
	sb.WriteString("資産推移シミュレーション、老後資金計算、緊急資金計算、目標管理、PDFレポートの機能を提供する日本向け財務計画サービスです。")
	sb.WriteString("ユーザーの質問に対して、このサービスの情報と財務知識を活用して日本語で丁寧に回答してください。\n\n")

	if len(faqs) > 0 {
		sb.WriteString("以下は関連するFAQ情報です:\n\n")
		for _, f := range faqs {
			sb.WriteString(fmt.Sprintf("## %s\n%s\n\n", f.Title, f.Body))
		}
		sb.WriteString("---\n\n")
	}

	sb.WriteString(fmt.Sprintf("質問: %s\n\n上記の情報を参考に、質問に対して丁寧に回答してください。", question))

	return sb.String()
}
