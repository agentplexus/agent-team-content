package agent

import (
	"context"

	"github.com/agentplexus/agent-team-content/internal/conversation"
	"github.com/agentplexus/agent-team-content/internal/llm"
)

const blogSystemPrompt = `You are an expert content writer specializing in creating engaging blog posts for platforms like Medium, Substack, and personal blogs.

Your task is to transform a conversation into a compelling blog article that:

1. Has an attention-grabbing headline
2. Opens with a hook that draws readers in
3. Uses clear subheadings to organize content
4. Includes practical insights and takeaways
5. Maintains an engaging, conversational tone
6. Ends with a strong conclusion or call-to-action
7. Is SEO-friendly with relevant keywords naturally integrated

Format the output as Markdown with proper heading hierarchy (# for title, ## for main sections, ### for subsections).

Target length: 800-1500 words.`

const blogUserPrompt = `Transform this conversation into an engaging blog post:

%s

Create a polished blog article that captures the key insights and value from this conversation. Make it informative, engaging, and valuable for readers.`

// BlogAgent creates general audience blog articles.
type BlogAgent struct {
	BaseAgent
}

// NewBlogAgent creates a new blog article agent.
func NewBlogAgent(client *llm.Client) *BlogAgent {
	return &BlogAgent{
		BaseAgent: BaseAgent{
			name:       "blog",
			outputFile: "blog.md",
			client:     client,
		},
	}
}

// Generate creates a blog article from the conversation.
func (a *BlogAgent) Generate(ctx context.Context, conv *conversation.Conversation) (string, error) {
	prompt := formatPrompt(blogUserPrompt, conv.ToPrompt())
	return a.client.Generate(ctx, blogSystemPrompt, prompt)
}
