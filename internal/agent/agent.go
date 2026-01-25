package agent

import (
	"context"

	"github.com/agentplexus/agent-team-content/internal/conversation"
	"github.com/agentplexus/agent-team-content/internal/llm"
)

// Agent represents a content generation agent.
type Agent interface {
	// Name returns the agent's identifier.
	Name() string

	// OutputFile returns the default output filename.
	OutputFile() string

	// Generate creates content from the conversation.
	Generate(ctx context.Context, conv *conversation.Conversation) (string, error)
}

// BaseAgent provides common functionality for agents.
type BaseAgent struct {
	name       string
	outputFile string
	client     *llm.Client
}

// Name returns the agent's identifier.
func (a *BaseAgent) Name() string {
	return a.name
}

// OutputFile returns the default output filename.
func (a *BaseAgent) OutputFile() string {
	return a.outputFile
}

// Result holds the output from an agent.
type Result struct {
	AgentName  string
	OutputFile string
	Content    string
	Error      error
}

// Options holds configuration for agent creation.
type Options struct {
	MarpTheme string // Path to custom Marp theme CSS
}
