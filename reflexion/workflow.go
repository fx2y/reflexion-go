package reflexion

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ReflexionAgentWorkflow(ctx workflow.Context, question string) (string, error) {
	ra := &ReflexionAgent{}
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting ReflexionAgentWorkflow", "question", question)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var initialResponse AgentResponse
	err := workflow.ExecuteActivity(ctx, ra.InitialResponderActivity, question).Get(ctx, &initialResponse)
	if err != nil {
		logger.Error("Failed to execute InitialResponderActivity", "error", err)
		return "", err
	}

	var searchResults []SearchResult
	err = workflow.ExecuteActivity(ctx, ra.SearchActivity, initialResponse.SearchQuery).Get(ctx, &searchResults)
	if err != nil {
		logger.Error("Failed to execute SearchActivity", "error", err)
		return "", err
	}

	var finalResponse AgentResponse
	err = workflow.ExecuteActivity(ctx, ra.RevisorActivity, question, initialResponse, searchResults).Get(ctx, &finalResponse)
	if err != nil {
		logger.Error("Failed to execute RevisorActivity", "error", err)
		return "", err
	}

	logger.Info("Completed ReflexionAgentWorkflow", "finalAnswer", finalResponse.Answer)
	return finalResponse.Answer, nil
}
