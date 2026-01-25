package conversation

import "time"

// Message represents a single message in a conversation.
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// Conversation represents a complete conversation with metadata.
type Conversation struct {
	Title    string            `json:"title,omitempty"`
	Messages []Message         `json:"messages"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ToPrompt converts the conversation to a formatted string for LLM prompts.
func (c *Conversation) ToPrompt() string {
	var result string
	if c.Title != "" {
		result = "# " + c.Title + "\n\n"
	}
	for _, msg := range c.Messages {
		result += "**" + msg.Role + ":** " + msg.Content + "\n\n"
	}
	return result
}

// Summary returns a brief summary of the conversation for context.
func (c *Conversation) Summary() string {
	if len(c.Messages) == 0 {
		return "Empty conversation"
	}
	// Return first user message as context
	for _, msg := range c.Messages {
		if msg.Role == "user" {
			content := msg.Content
			if len(content) > 200 {
				content = content[:200] + "..."
			}
			return content
		}
	}
	return c.Messages[0].Content
}
