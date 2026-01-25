# Agent Team Content

A multi-agent team for transforming conversations into multiple content formats. This team generates blog articles, social media posts, and presentations from a single input.

## Agents

| Agent | Description |
|-------|-------------|
| `blog` | Creates engaging blog articles for Medium, Substack, and personal blogs |
| `devto` | Creates technical articles for the dev.to developer community |
| `linkedin` | Creates professional LinkedIn posts that drive engagement |
| `twitter` | Creates viral Twitter/X threads |
| `marp` | Creates Marp Markdown presentations with speaker notes |
| `revealjs` | Creates Reveal.js presentations with horizontal and vertical navigation |

## Workflow

The team uses a parallel workflow - all agents run independently on the same input conversation.

## Project Structure

```
agent-team-content/
├── specs/
│   ├── agents/           # Agent definitions (*.md with YAML frontmatter)
│   │   ├── blog.md
│   │   ├── devto.md
│   │   ├── linkedin.md
│   │   ├── twitter.md
│   │   ├── marp.md
│   │   └── revealjs.md
│   ├── teams/            # Team workflow definitions
│   │   └── content-team.json
│   └── deployments/      # Deployment configurations
│       └── local.json
├── cmd/
│   └── content/          # Runtime CLI for generating content
├── .claude/
│   └── agents/           # Generated Claude Code agents
└── plugins/
    └── kiro/
        └── agents/       # Generated Kiro CLI agents
```

## Generating Agent Files

Use [AssistantKit](https://github.com/agentplexus/assistantkit) to generate platform-specific agent files from the specs.

### Install AssistantKit

```bash
go install github.com/agentplexus/assistantkit/cmd/assistantkit@latest
```

### Generate Agents

From the repository root:

```bash
# Generate using defaults (--target=local, --output=.)
assistantkit generate agents

# Or explicitly specify options
assistantkit generate agents --specs=specs --target=local --output=.
```

This generates:

- `.claude/agents/*.md` - Claude Code agent markdown files
- `plugins/kiro/agents/*.json` - Kiro CLI agent JSON files

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--specs` | `specs` | Path to specs directory containing agents/ and deployments/ |
| `--target` | `local` | Deployment target (looks for `specs/deployments/<target>.json`) |
| `--output` | `.` | Output base directory (repo root) |

## Runtime CLI

The `content` CLI runs the agent team to generate content from conversations.

### Build

```bash
go build -o content ./cmd/content
```

### Usage

```bash
# Generate all content types from a conversation
./content generate --input=conversation.json --output=./output

# Generate specific formats only
./content generate --input=conversation.json --agents=blog,twitter

# List available agents
./content list-agents
```

### Input Format

The input can be a JSON conversation file:

```json
{
  "messages": [
    {"role": "user", "content": "Let's discuss AI agents..."},
    {"role": "assistant", "content": "AI agents are..."}
  ]
}
```

Or a Markdown file with the conversation.

## License

MIT License
