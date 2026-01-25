# Release Notes - v0.1.0

**Release Date:** January 25, 2026

## Overview

This is the initial release of **agent-team-content**, a multi-agent system for transforming conversations into multiple content formats. It orchestrates 6 specialized AI agents in a parallel workflow to produce blog articles, social media posts, and presentations from a single input.

agent-team-content integrates with:

- [AssistantKit](https://github.com/agentplexus/assistantkit) for generating platform-specific agents
- Claude API via [anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go)

## Highlights

- **6 Specialized Agents** - Content generation for multiple platforms
- **Parallel Workflow** - All agents run independently for fast generation
- **Multiple Output Formats** - Blog, technical articles, social posts, presentations
- **Claude Integration** - AI-powered content generation via Claude API
- **Cross-Platform Agents** - Generate Claude Code and Kiro agents from specs

## Installation

```bash
go install github.com/agentplexus/agent-team-content/cmd/content@v0.1.0
```

## Features

### content CLI

Command-line tool for content generation:

```bash
content generate --input conversation.json --output ./output  # Generate all formats
content generate --input conv.json --agents blog,twitter      # Specific agents
content list-agents                                           # List available agents
content version                                               # Show version
```

### Agent Team (6 Agents)

| Agent | Description |
|-------|-------------|
| blog | Creates engaging blog articles for Medium, Substack, and personal blogs |
| devto | Creates technical articles for the dev.to developer community |
| linkedin | Creates professional LinkedIn posts that drive engagement |
| twitter | Creates viral Twitter/X threads |
| marp | Creates Marp Markdown presentations with speaker notes |
| revealjs | Creates Reveal.js presentations with horizontal and vertical navigation |

### Input Formats

JSON conversation file:

```json
{
  "messages": [
    {"role": "user", "content": "Let's discuss AI agents..."},
    {"role": "assistant", "content": "AI agents are..."}
  ]
}
```

Or Markdown file with conversation text.

### Output

All agents run in parallel and write their output to the specified directory:

- `blog-article.md` - Blog article in Markdown
- `devto-article.md` - dev.to technical article
- `linkedin-post.md` - LinkedIn post
- `twitter-thread.md` - Twitter/X thread
- `presentation.marp.md` - Marp presentation
- `presentation-revealjs.html` - Reveal.js presentation
- `summary.json` - Generation metadata

## Quick Start

```bash
# Build the CLI
go build -o content ./cmd/content

# Generate all content types from a conversation
./content generate --input=examples/conversation.json --output=./output

# Generate specific formats only
./content generate --input=examples/conversation.json --agents=blog,twitter

# List available agents
./content list-agents
```

## Project Structure

```
agent-team-content/
├── cmd/content/           # CLI application
├── internal/
│   ├── agent/             # Agent implementations and orchestrator
│   ├── conversation/      # Input parsing (JSON, Markdown)
│   └── llm/               # Claude API client
├── specs/
│   ├── agents/            # Agent definitions (YAML frontmatter + Markdown)
│   ├── teams/             # Team workflow definitions
│   └── deployments/       # Deployment configurations
├── .claude/agents/        # Generated Claude Code agents
└── plugins/kiro/agents/   # Generated Kiro agents
```

## Dependencies

- Go 1.25+
- github.com/agentplexus/assistantkit v0.7.0
- github.com/anthropics/anthropic-sdk-go v1.19.0
- github.com/spf13/cobra v1.10.2

## Contributors

- John Wang (@grokify)
- Claude Opus 4.5 (Co-Author)

## Links

- [GitHub Repository](https://github.com/agentplexus/agent-team-content)
- [AssistantKit](https://github.com/agentplexus/assistantkit) - Agent generation toolkit
- [Changelog](CHANGELOG.md)
