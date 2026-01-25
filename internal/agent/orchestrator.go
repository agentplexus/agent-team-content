package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/agentplexus/agent-team-content/internal/conversation"
	"github.com/agentplexus/agent-team-content/internal/llm"
)

// Orchestrator coordinates multiple agents to generate content.
type Orchestrator struct {
	client  *llm.Client
	agents  []Agent
	options Options
}

// NewOrchestrator creates a new orchestrator with the specified agents.
func NewOrchestrator(client *llm.Client, opts Options) *Orchestrator {
	agents := []Agent{
		NewBlogAgent(client),
		NewDevToAgent(client),
		NewLinkedInAgent(client),
		NewTwitterAgent(client),
		NewMarpAgent(client, opts.MarpTheme),
		NewRevealJSAgent(client),
	}

	return &Orchestrator{
		client:  client,
		agents:  agents,
		options: opts,
	}
}

// NewOrchestratorWithAgents creates an orchestrator with specific agents.
func NewOrchestratorWithAgents(client *llm.Client, agentNames []string, opts Options) (*Orchestrator, error) {
	agentMap := map[string]func() Agent{
		"blog":     func() Agent { return NewBlogAgent(client) },
		"devto":    func() Agent { return NewDevToAgent(client) },
		"linkedin": func() Agent { return NewLinkedInAgent(client) },
		"twitter":  func() Agent { return NewTwitterAgent(client) },
		"marp":     func() Agent { return NewMarpAgent(client, opts.MarpTheme) },
		"revealjs": func() Agent { return NewRevealJSAgent(client) },
	}

	var agents []Agent
	for _, name := range agentNames {
		createFn, ok := agentMap[name]
		if !ok {
			return nil, fmt.Errorf("unknown agent: %s", name)
		}
		agents = append(agents, createFn())
	}

	return &Orchestrator{
		client:  client,
		agents:  agents,
		options: opts,
	}, nil
}

// Generate runs all agents concurrently and collects results.
func (o *Orchestrator) Generate(ctx context.Context, conv *conversation.Conversation) []Result {
	var (
		results []Result
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	for _, agent := range o.agents {
		wg.Add(1)
		go func(a Agent) {
			defer wg.Done()

			content, err := a.Generate(ctx, conv)
			result := Result{
				AgentName:  a.Name(),
				OutputFile: a.OutputFile(),
				Content:    content,
				Error:      err,
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(agent)
	}

	wg.Wait()
	return results
}

// ListAgents returns the names of all available agents.
func ListAgents() []string {
	return []string{"blog", "devto", "linkedin", "twitter", "marp", "revealjs"}
}
