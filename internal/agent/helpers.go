package agent

import "fmt"

// formatPrompt formats a prompt template with the conversation content.
func formatPrompt(template, content string) string {
	return fmt.Sprintf(template, content)
}
