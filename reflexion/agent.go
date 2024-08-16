package reflexion

import (
	"os"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
	openai "github.com/sashabaranov/go-openai"
	"go.temporal.io/sdk/worker"
)

// ReflexionAgent represents the main agent structure
type ReflexionAgent struct {
	client        *instructor.InstructorOpenAI
	searchService SearchService
}

// NewReflexionAgent creates a new instance of ReflexionAgent
func NewReflexionAgent(searchService SearchService) *ReflexionAgent {
	client := instructor.FromOpenAI(
		openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		instructor.WithMode(instructor.ModeJSON),
		instructor.WithMaxRetries(3),
	)
	return &ReflexionAgent{
		client:        client,
		searchService: searchService,
	}
}

// RegisterWorkflowAndActivities registers the workflow and activities with the Temporal worker
func (ra *ReflexionAgent) RegisterWorkflowAndActivities(w worker.Worker) {
	w.RegisterWorkflow(ReflexionAgentWorkflow)
	w.RegisterActivity(ra.InitialResponderActivity)
	w.RegisterActivity(ra.SearchActivity)
	w.RegisterActivity(ra.RevisorActivity)
}
