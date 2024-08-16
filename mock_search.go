package main

import (
	"context"

	"github.com/fx2y/reflexion-go/reflexion"
)

type MockSearchService struct{}

func (m *MockSearchService) Search(ctx context.Context, query string) ([]reflexion.SearchResult, error) {
	return []reflexion.SearchResult{
		{
			URL:     "https://example.com",
			Content: "This is a mock search result for query: " + query,
		},
	}, nil
}
