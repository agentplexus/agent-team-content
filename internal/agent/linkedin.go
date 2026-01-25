package agent

import (
	"context"

	"github.com/agentplexus/agent-team-content/internal/conversation"
	"github.com/agentplexus/agent-team-content/internal/llm"
)

const linkedinSystemPrompt = `You are a LinkedIn content strategist creating professional posts that drive engagement.

Your task is to transform a conversation into a LinkedIn post that:

1. Opens with a strong hook (first line is crucial - it shows in the preview)
2. Uses short paragraphs and line breaks for readability
3. Includes relevant emojis sparingly for visual appeal
4. Provides value or insight to a professional audience
5. Ends with a question or call-to-action to encourage engagement
6. Uses relevant hashtags (3-5 at the end)

Format guidelines:
- Maximum 1300 characters (LinkedIn's limit before "see more")
- Use line breaks between paragraphs
- Start with attention-grabbing first line
- Include 3-5 relevant hashtags at the end

Tone: Professional but approachable, thought-leadership oriented.`

const linkedinUserPrompt = `Transform this conversation into a LinkedIn post:

%s

Create a compelling LinkedIn post that shares key insights professionally and encourages engagement.`

// LinkedInAgent creates LinkedIn posts.
type LinkedInAgent struct {
	BaseAgent
}

// NewLinkedInAgent creates a new LinkedIn post agent.
func NewLinkedInAgent(client *llm.Client) *LinkedInAgent {
	return &LinkedInAgent{
		BaseAgent: BaseAgent{
			name:       "linkedin",
			outputFile: "linkedin.md",
			client:     client,
		},
	}
}

// Generate creates a LinkedIn post from the conversation.
func (a *LinkedInAgent) Generate(ctx context.Context, conv *conversation.Conversation) (string, error) {
	prompt := formatPrompt(linkedinUserPrompt, conv.ToPrompt())
	return a.client.Generate(ctx, linkedinSystemPrompt, prompt)
}
