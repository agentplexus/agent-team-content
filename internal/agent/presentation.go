package agent

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/agentplexus/agent-team-content/internal/conversation"
	"github.com/agentplexus/agent-team-content/internal/llm"
)

const marpSystemPrompt = `You are a presentation designer creating Marp Markdown slides.

Your task is to transform a conversation into a Marp presentation that:

1. Starts with Marp frontmatter (marp: true, theme, paginate, etc.)
2. Has a clear title slide
3. Organizes content into logical slides (use --- as slide separators)
4. Uses bullet points for clarity
5. Includes speaker notes where helpful (using <!-- --> comments)
6. Keeps slides focused - one main idea per slide
7. Uses proper heading hierarchy within slides

Format:
---
marp: true
theme: %s
paginate: true
header: ''
footer: ''
---

# Presentation Title

---

## Slide Title

- Key point 1
- Key point 2

<!-- Speaker notes go here -->

---

... continue with more slides

Target: 8-15 slides.`

const marpUserPrompt = `Transform this conversation into a Marp presentation:

%s

Create a clear, well-organized presentation that communicates the key points effectively.`

const revealjsSystemPrompt = `You are a presentation designer creating Reveal.js Markdown slides.

Your task is to transform a conversation into a Reveal.js presentation that:

1. Uses --- as horizontal slide separators
2. Uses -- as vertical slide separators (for drill-down content)
3. Has a clear title slide
4. Organizes content logically
5. Includes speaker notes using Note: syntax
6. Keeps slides focused
7. Uses Markdown formatting effectively

Format:
# Presentation Title

---

## Section Title

Content here

Note: Speaker notes for this slide

---

## Another Section

--

### Subsection (vertical slide)

More detail here

... continue with more slides

Target: 8-15 slides.`

const revealjsUserPrompt = `Transform this conversation into a Reveal.js presentation:

%s

Create a clear, well-organized presentation that communicates the key points effectively with good use of horizontal and vertical slide organization.`

// MarpAgent creates Marp presentations.
type MarpAgent struct {
	BaseAgent
	theme string
}

// NewMarpAgent creates a new Marp presentation agent.
func NewMarpAgent(client *llm.Client, themePath string) *MarpAgent {
	theme := "default"
	if themePath != "" {
		// Read custom theme or use the path as theme name
		if _, err := os.Stat(themePath); err == nil {
			// File exists, embed it
			themeContent, err := os.ReadFile(themePath)
			if err == nil {
				theme = fmt.Sprintf("custom\nstyle: |\n%s", indentCSS(string(themeContent)))
			}
		} else {
			// Use as theme name
			theme = themePath
		}
	}

	return &MarpAgent{
		BaseAgent: BaseAgent{
			name:       "marp",
			outputFile: "marp.md",
			client:     client,
		},
		theme: theme,
	}
}

// Generate creates a Marp presentation from the conversation.
func (a *MarpAgent) Generate(ctx context.Context, conv *conversation.Conversation) (string, error) {
	systemPrompt := fmt.Sprintf(marpSystemPrompt, a.theme)
	prompt := formatPrompt(marpUserPrompt, conv.ToPrompt())
	return a.client.Generate(ctx, systemPrompt, prompt)
}

// RevealJSAgent creates Reveal.js presentations.
type RevealJSAgent struct {
	BaseAgent
}

// NewRevealJSAgent creates a new Reveal.js presentation agent.
func NewRevealJSAgent(client *llm.Client) *RevealJSAgent {
	return &RevealJSAgent{
		BaseAgent: BaseAgent{
			name:       "revealjs",
			outputFile: "revealjs.md",
			client:     client,
		},
	}
}

// Generate creates a Reveal.js presentation from the conversation.
func (a *RevealJSAgent) Generate(ctx context.Context, conv *conversation.Conversation) (string, error) {
	prompt := formatPrompt(revealjsUserPrompt, conv.ToPrompt())
	return a.client.Generate(ctx, revealjsSystemPrompt, prompt)
}

// indentCSS adds proper indentation for YAML embedding.
func indentCSS(css string) string {
	lines := strings.Split(css, "\n")
	var indented []string
	for _, line := range lines {
		indented = append(indented, "    "+line)
	}
	return strings.Join(indented, "\n")
}
