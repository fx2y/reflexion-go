package reflexion

import (
	"context"
	"fmt"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
	openai "github.com/sashabaranov/go-openai"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type SearchService interface {
	Search(ctx context.Context, query string) ([]SearchResult, error)
}

type SearchResult struct {
	URL     string `json:"url"`
	Content string `json:"content"`
}

type TavilySearchService struct {
	client *instructor.InstructorOpenAI
}

func NewTavilySearchService(client *instructor.InstructorOpenAI) *TavilySearchService {
	return &TavilySearchService{client: client}
}

func (s *TavilySearchService) Search(ctx context.Context, query string) ([]SearchResult, error) {
	logger := activity.GetLogger(ctx)

	type TavilySearchResponse struct {
		Results []SearchResult `json:"results"`
	}

	var response TavilySearchResponse
	_, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf(`You are a search engine. Perform a search for the following query:
                
                %s
                
                Return the results as a JSON array of objects with 'url' and 'content' fields.`, query),
			},
		},
	}, &response)

	if err != nil {
		logger.Error("Failed to perform search", "error", err)
		return nil, err
	}

	return response.Results, nil
}

func RegisterSearchActivity(worker worker.Worker, searchService SearchService) {
	worker.RegisterActivity(searchActivity)
}

func searchActivity(ctx context.Context, query string) ([]SearchResult, error) {
	searchService := ctx.Value("searchService").(SearchService)
	return searchService.Search(ctx, query)
}
