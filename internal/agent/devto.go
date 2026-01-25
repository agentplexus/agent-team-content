package agent

import (
	"context"

	"github.com/agentplexus/agent-team-content/internal/conversation"
	"github.com/agentplexus/agent-team-content/internal/llm"
)

const devtoSystemPrompt = `You are a technical writer creating articles for dev.to, a community of software developers.

Your task is to transform a conversation into a technical blog post that:

1. Starts with dev.to frontmatter (title, published, description, tags, cover_image placeholder)
2. Has a clear, descriptive title
3. Includes code examples where relevant (properly formatted with language hints)
4. Explains technical concepts clearly
5. Uses practical examples developers can relate to
6. Includes relevant diagrams or architecture descriptions when helpful
7. Ends with next steps or resources for further learning

Format the output as Markdown with dev.to frontmatter at the top:

---
title: "Your Title Here"
published: false
description: "Brief description for SEO"
tags: tag1, tag2, tag3, tag4
cover_image: https://dev.to/placeholder.png
---

Target length: 1000-2000 words.
Focus on technical accuracy and practical value for developers.`

const devtoUserPrompt = `Transform this conversation into a technical dev.to article:

%s

Create a developer-focused article that provides technical insights, code examples where appropriate, and practical value for the dev.to community.`

// DevToAgent creates technical articles for dev.to.
type DevToAgent struct {
	BaseAgent
}

// NewDevToAgent creates a new dev.to article agent.
func NewDevToAgent(client *llm.Client) *DevToAgent {
	return &DevToAgent{
		BaseAgent: BaseAgent{
			name:       "devto",
			outputFile: "devto.md",
			client:     client,
		},
	}
}

// Generate creates a dev.to article from the conversation.
func (a *DevToAgent) Generate(ctx context.Context, conv *conversation.Conversation) (string, error) {
	prompt := formatPrompt(devtoUserPrompt, conv.ToPrompt())
	return a.client.Generate(ctx, devtoSystemPrompt, prompt)
}
