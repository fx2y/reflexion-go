package reflexion

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// TracedReflexionAgent wraps ReflexionAgent with tracing capabilities
type TracedReflexionAgent struct {
	*ReflexionAgent
	tracer trace.Tracer
	logger log.Logger
}

// NewTracedReflexionAgent creates a new TracedReflexionAgent
func NewTracedReflexionAgent(agent *ReflexionAgent, logger log.Logger) *TracedReflexionAgent {
	return &TracedReflexionAgent{
		ReflexionAgent: agent,
		tracer:         otel.Tracer("reflexion-agent"),
		logger:         logger,
	}
}

// TracedInitialResponderActivity wraps InitialResponderActivity with tracing
func (tra *TracedReflexionAgent) TracedInitialResponderActivity(ctx context.Context, question string) (AgentResponse, error) {
	ctx, span := tra.tracer.Start(ctx, "InitialResponderActivity")
	defer span.End()

	span.SetAttributes(attribute.String("question", question))
	tra.logger.Info("Starting InitialResponderActivity", "question", question)

	response, err := tra.ReflexionAgent.InitialResponderActivity(ctx, question)
	if err != nil {
		span.RecordError(err)
		tra.logger.Error("InitialResponderActivity failed", "error", err)
		return AgentResponse{}, err
	}

	span.SetAttributes(
		attribute.String("answer", response.Answer),
		attribute.String("reflection", response.Reflection),
		attribute.String("search_query", response.SearchQuery),
	)
	tra.logger.Info("InitialResponderActivity completed", "response", response)

	return response, nil
}

// TracedSearchActivity wraps SearchActivity with tracing
func (tra *TracedReflexionAgent) TracedSearchActivity(ctx context.Context, query string) ([]SearchResult, error) {
	ctx, span := tra.tracer.Start(ctx, "SearchActivity")
	defer span.End()

	span.SetAttributes(attribute.String("query", query))
	tra.logger.Info("Starting SearchActivity", "query", query)

	results, err := tra.ReflexionAgent.SearchActivity(ctx, query)
	if err != nil {
		span.RecordError(err)
		tra.logger.Error("SearchActivity failed", "error", err)
		return nil, err
	}

	span.SetAttributes(attribute.Int("result_count", len(results)))
	tra.logger.Info("SearchActivity completed", "result_count", len(results))

	return results, nil
}

// TracedRevisorActivity wraps RevisorActivity with tracing
func (tra *TracedReflexionAgent) TracedRevisorActivity(ctx context.Context, question string, initialResponse AgentResponse, searchResults []SearchResult) (AgentResponse, error) {
	ctx, span := tra.tracer.Start(ctx, "RevisorActivity")
	defer span.End()

	span.SetAttributes(
		attribute.String("question", question),
		attribute.String("initial_answer", initialResponse.Answer),
		attribute.Int("search_result_count", len(searchResults)),
	)
	tra.logger.Info("Starting RevisorActivity", "question", question, "initial_answer", initialResponse.Answer)

	response, err := tra.ReflexionAgent.RevisorActivity(ctx, question, initialResponse, searchResults)
	if err != nil {
		span.RecordError(err)
		tra.logger.Error("RevisorActivity failed", "error", err)
		return AgentResponse{}, err
	}

	span.SetAttributes(
		attribute.String("revised_answer", response.Answer),
		attribute.String("revised_reflection", response.Reflection),
		attribute.String("revised_search_query", response.SearchQuery),
	)
	tra.logger.Info("RevisorActivity completed", "response", response)

	return response, nil
}

// TracedReflexionAgentWorkflow wraps ReflexionAgentWorkflow with tracing
func TracedReflexionAgentWorkflow(ctx workflow.Context, question string) (string, error) {
	logger := workflow.GetLogger(ctx)
	// tracer := workflow.GetActivityOptions(ctx).StartToCloseTimeout // Unused variable

	logger.Info("Starting ReflexionAgentWorkflow", "question", question)

	var initialResponse AgentResponse
	err := workflow.ExecuteActivity(ctx, "TracedInitialResponderActivity", question).Get(ctx, &initialResponse)
	if err != nil {
		logger.Error("InitialResponderActivity failed", "error", err)
		return "", err
	}

	var searchResults []SearchResult
	err = workflow.ExecuteActivity(ctx, "TracedSearchActivity", initialResponse.SearchQuery).Get(ctx, &searchResults)
	if err != nil {
		logger.Error("SearchActivity failed", "error", err)
		return "", err
	}

	var finalResponse AgentResponse
	err = workflow.ExecuteActivity(ctx, "TracedRevisorActivity", question, initialResponse, searchResults).Get(ctx, &finalResponse)
	if err != nil {
		logger.Error("RevisorActivity failed", "error", err)
		return "", err
	}

	logger.Info("ReflexionAgentWorkflow completed", "final_answer", finalResponse.Answer)

	return finalResponse.Answer, nil
}

// RegisterTracedWorkflowAndActivities registers the traced workflow and activities
func (tra *TracedReflexionAgent) RegisterTracedWorkflowAndActivities(worker worker.Worker) {
	worker.RegisterWorkflow(TracedReflexionAgentWorkflow)
	worker.RegisterActivity(tra.TracedInitialResponderActivity)
	worker.RegisterActivity(tra.TracedSearchActivity)
	worker.RegisterActivity(tra.TracedRevisorActivity)
}
