package agent

import (
	"context"

	"github.com/agentplexus/agent-team-content/internal/conversation"
	"github.com/agentplexus/agent-team-content/internal/llm"
)

const twitterSystemPrompt = `You are a Twitter/X content creator specializing in viral threads.

Your task is to transform a conversation into a Twitter thread that:

1. Opens with a hook tweet that makes people want to read more
2. Numbers each tweet (1/, 2/, etc.)
3. Keeps each tweet under 280 characters
4. Uses clear, punchy language
5. Includes relevant insights broken into digestible pieces
6. Ends with a summary or call-to-action tweet
7. Suggests relevant hashtags for the final tweet

Format the output as a numbered thread:

1/ [Hook tweet - grab attention]

2/ [Key insight #1]

3/ [Key insight #2]

... and so on

Final tweet should encourage engagement (retweet, follow, etc.) and include 2-3 relevant hashtags.

Target: 5-10 tweets in the thread.`

const twitterUserPrompt = `Transform this conversation into a Twitter/X thread:

%s

Create an engaging thread that breaks down the key insights into tweet-sized pieces that people will want to read and share.`

// TwitterAgent creates Twitter/X threads.
type TwitterAgent struct {
	BaseAgent
}

// NewTwitterAgent creates a new Twitter thread agent.
func NewTwitterAgent(client *llm.Client) *TwitterAgent {
	return &TwitterAgent{
		BaseAgent: BaseAgent{
			name:       "twitter",
			outputFile: "twitter.md",
			client:     client,
		},
	}
}

// Generate creates a Twitter thread from the conversation.
func (a *TwitterAgent) Generate(ctx context.Context, conv *conversation.Conversation) (string, error) {
	prompt := formatPrompt(twitterUserPrompt, conv.ToPrompt())
	return a.client.Generate(ctx, twitterSystemPrompt, prompt)
}
