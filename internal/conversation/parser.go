package conversation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ParseFile parses a conversation from a file, auto-detecting format.
func ParseFile(path string) (*Conversation, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Try JSON first
	if strings.HasSuffix(strings.ToLower(path), ".json") {
		return ParseJSON(data)
	}

	// Try Markdown
	if strings.HasSuffix(strings.ToLower(path), ".md") {
		return ParseMarkdown(data)
	}

	// Auto-detect: try JSON, then Markdown
	conv, err := ParseJSON(data)
	if err == nil {
		return conv, nil
	}

	return ParseMarkdown(data)
}

// ParseJSON parses a conversation from JSON data.
func ParseJSON(data []byte) (*Conversation, error) {
	var conv Conversation
	if err := json.Unmarshal(data, &conv); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &conv, nil
}

// ParseMarkdown parses a conversation from Markdown format.
// Supports formats like:
//   - **User:** message
//   - **Assistant:** message
//   - User: message
//   - Assistant: message
func ParseMarkdown(data []byte) (*Conversation, error) {
	conv := &Conversation{
		Messages: []Message{},
		Metadata: make(map[string]string),
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var currentRole string
	var currentContent strings.Builder

	// Regex patterns for role detection
	boldRolePattern := regexp.MustCompile(`^\*\*([Uu]ser|[Aa]ssistant|[Ss]ystem)\*\*:\s*(.*)$`)
	simpleRolePattern := regexp.MustCompile(`^([Uu]ser|[Aa]ssistant|[Ss]ystem):\s*(.*)$`)
	titlePattern := regexp.MustCompile(`^#\s+(.+)$`)

	flushMessage := func() {
		if currentRole != "" && currentContent.Len() > 0 {
			conv.Messages = append(conv.Messages, Message{
				Role:    strings.ToLower(currentRole),
				Content: strings.TrimSpace(currentContent.String()),
			})
		}
		currentContent.Reset()
	}

	for scanner.Scan() {
		line := scanner.Text()

		// Check for title
		if matches := titlePattern.FindStringSubmatch(line); matches != nil {
			if conv.Title == "" {
				conv.Title = matches[1]
			}
			continue
		}

		// Check for bold role pattern
		if matches := boldRolePattern.FindStringSubmatch(line); matches != nil {
			flushMessage()
			currentRole = matches[1]
			if matches[2] != "" {
				currentContent.WriteString(matches[2])
			}
			continue
		}

		// Check for simple role pattern
		if matches := simpleRolePattern.FindStringSubmatch(line); matches != nil {
			flushMessage()
			currentRole = matches[1]
			if matches[2] != "" {
				currentContent.WriteString(matches[2])
			}
			continue
		}

		// Continue building current message
		if currentRole != "" {
			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(line)
		}
	}

	// Flush last message
	flushMessage()

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan markdown: %w", err)
	}

	return conv, nil
}
