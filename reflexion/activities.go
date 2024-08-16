package reflexion

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"go.temporal.io/sdk/activity"
)

func (ra *ReflexionAgent) InitialResponderActivity(ctx context.Context, question string) (AgentResponse, error) {
    logger := activity.GetLogger(ctx)
    logger.Info("Starting InitialResponderActivity", "question", question)

    var response AgentResponse
    for attempt := 0; attempt < 3; attempt++ {
        _, err := ra.client.CreateChatCompletion(
            ctx,
            openai.ChatCompletionRequest{
                Model: openai.GPT3Dot5Turbo,
                Messages: []openai.ChatCompletionMessage{
                    {Role: openai.ChatMessageRoleSystem, Content: "You are a Reflexion Agent. Generate an initial answer, reflection, and search query for the given question."},
                    {Role: openai.ChatMessageRoleUser, Content: question},
                },
            },
            &response,
        )

        if err == nil {
            logger.Info("Completed InitialResponderActivity", "response", response)
            return response, nil
        }

        logger.Error("Error in InitialResponderActivity", "error", err, "attempt", attempt)
    }

    return AgentResponse{}, fmt.Errorf("failed after multiple attempts")
}
	logger := activity.GetLogger(ctx)
	logger.Info("Starting InitialResponderActivity", "question", question)

	var response AgentResponse
	for attempt := 0; attempt < 3; attempt++ {
		_, err := ra.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleSystem, Content: "You are a Reflexion Agent. Generate an initial answer, reflection, and search query for the given question."},
					{Role: openai.ChatMessageRoleUser, Content: question},
				},
			},
			&response,
		)

		if err == nil {
			logger.Info("Completed InitialResponderActivity", "response", response)
			return response, nil
		}

		logger.Error("Error in InitialResponderActivity", "error", err, "attempt", attempt)
	}

	return AgentResponse{}, fmt.Errorf("failed after multiple attempts")
}

func (ra *ReflexionAgent) SearchActivity(ctx context.Context, query string) ([]SearchResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting SearchActivity", "query", query)

	results, err := ra.searchService.Search(ctx, query)
	if err != nil {
		logger.Error("Error in SearchActivity", "error", err)
		return nil, err
	}

	logger.Info("Completed SearchActivity", "resultCount", len(results))
	return results, nil
}

func (ra *ReflexionAgent) RevisorActivity(ctx context.Context, question string, initialResponse AgentResponse, searchResults []SearchResult) (AgentResponse, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting RevisorActivity", "question", question, "initialResponse", initialResponse)

	searchContext := formatSearchResults(searchResults)

	var response AgentResponse
	for attempt := 0; attempt < 3; attempt++ {
		_, err := ra.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleSystem, Content: "You are a Reflexion Agent. Revise the initial answer based on the search results provided."},
					{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("Question: %s\nInitial Answer: %s\nInitial Reflection: %s\nSearch Results: %s",
						question, initialResponse.Answer, initialResponse.Reflection, searchContext)},
				},
			},
			&response,
		)

		if err == nil {
			logger.Info("Completed RevisorActivity", "response", response)
			return response, nil
		}

		logger.Error("Error in RevisorActivity", "error", err, "attempt", attempt)
	}

	return AgentResponse{}, fmt.Errorf("failed after multiple attempts")
}

func formatSearchResults(results []SearchResult) string {
	formatted := ""
	for _, result := range results {
		formatted += fmt.Sprintf("URL: %s\nContent: %s\n\n", result.URL, result.Content)
	}
	return formatted
}
