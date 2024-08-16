package reflexion

import (
	"time"

	"go.temporal.io/sdk/temporal"
)

// AgentResponse represents the structured output from the Reflexion Agent
type AgentResponse struct {
	Answer      string
	Reflection  string
	SearchQuery string
}

// WorkflowInput represents the input to the Reflexion Agent workflow
type WorkflowInput struct {
	Question   string
	MaxRetries int
}

// WorkflowOptions represents the configurable options for the Reflexion Agent workflow
type WorkflowOptions struct {
	TaskQueue           string
	WorkflowRunTimeout  time.Duration
	ActivityRetryPolicy *temporal.RetryPolicy
}

// SearchServiceConfig represents the configuration for the search service
type SearchServiceConfig struct {
	APIKey     string
	MaxResults int
	Timeout    time.Duration
}

// InstructorConfig represents the configuration for the instructor-ai client
type InstructorConfig struct {
	Model       string
	MaxTokens   int
	Temperature float32
}

// ReflexionAgentConfig combines all configurations for the Reflexion Agent
type ReflexionAgentConfig struct {
	WorkflowOptions     WorkflowOptions
	SearchServiceConfig SearchServiceConfig
	InstructorConfig    InstructorConfig
}

// This struct aligns with the instructor-go library's expectations
type StructuredPrompt struct {
	Question   string `json:"question"`
	Reflection string `json:"reflection,omitempty"`
	Context    string `json:"context,omitempty"`
}

// This struct aligns with the instructor-go library's response format
type StructuredResponse struct {
	Answer      string `json:"answer"`
	Reflection  string `json:"reflection"`
	SearchQuery string `json:"search_query"`
}
