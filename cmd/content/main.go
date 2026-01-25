package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/agentplexus/agent-team-content/internal/agent"
	"github.com/agentplexus/agent-team-content/internal/conversation"
	"github.com/agentplexus/agent-team-content/internal/llm"
	"github.com/spf13/cobra"
)

var (
	inputFile string
	outputDir string
	marpTheme string
	agentList string
	model     string
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "content",
		Short: "Content Agent Team - Generate content in multiple formats from conversations",
		Long: `Content Agent Team transforms conversations into multiple content formats:
- Blog articles (Medium/Substack)
- Technical articles (dev.to)
- LinkedIn posts
- Twitter/X threads
- Marp presentations
- Reveal.js presentations`,
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate content from a conversation",
		RunE:  runGenerate,
	}

	generateCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input conversation file (JSON or Markdown)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "./output", "Output directory")
	generateCmd.Flags().StringVar(&marpTheme, "theme", "", "Custom Marp theme CSS file")
	generateCmd.Flags().StringVar(&agentList, "agents", "", "Comma-separated list of agents (default: all)")
	generateCmd.Flags().StringVar(&model, "model", "claude-sonnet-4-20250514", "Claude model to use")
	if err := generateCmd.MarkFlagRequired("input"); err != nil {
		panic(err)
	}

	listCmd := &cobra.Command{
		Use:   "list-agents",
		Short: "List available agents",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available agents:")
			for _, name := range agent.ListAgents() {
				fmt.Printf("  - %s\n", name)
			}
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("content version %s\n", version)
		},
	}

	rootCmd.AddCommand(generateCmd, listCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Parse conversation
	conv, err := conversation.ParseFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse conversation: %w", err)
	}

	// Create LLM client
	cfg := llm.DefaultConfig()
	cfg.Model = model

	client, err := llm.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	// Create orchestrator
	opts := agent.Options{
		MarpTheme: marpTheme,
	}

	var orchestrator *agent.Orchestrator
	if agentList != "" {
		agents := strings.Split(agentList, ",")
		for i := range agents {
			agents[i] = strings.TrimSpace(agents[i])
		}
		orchestrator, err = agent.NewOrchestratorWithAgents(client, agents, opts)
		if err != nil {
			return fmt.Errorf("failed to create orchestrator: %w", err)
		}
	} else {
		orchestrator = agent.NewOrchestrator(client, opts)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate content
	fmt.Printf("Generating content from: %s\n", inputFile)
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Println()

	ctx := context.Background()
	startTime := time.Now()
	results := orchestrator.Generate(ctx, conv)
	duration := time.Since(startTime)

	// Write results
	var successCount, errorCount int
	summary := Summary{
		InputFile:   inputFile,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Duration:    duration.String(),
		Outputs:     make([]OutputSummary, 0, len(results)),
	}

	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("  [ERROR] %s: %v\n", result.AgentName, result.Error)
			errorCount++
			continue
		}

		outputPath := filepath.Join(outputDir, result.OutputFile)
		if err := os.WriteFile(outputPath, []byte(result.Content), 0600); err != nil {
			fmt.Printf("  [ERROR] Failed to write %s: %v\n", result.OutputFile, err)
			errorCount++
			continue
		}

		fmt.Printf("  [OK] %s -> %s\n", result.AgentName, result.OutputFile)
		successCount++

		summary.Outputs = append(summary.Outputs, OutputSummary{
			Agent: result.AgentName,
			File:  result.OutputFile,
		})
	}

	// Write summary
	summaryPath := filepath.Join(outputDir, "summary.json")
	summaryData, _ := json.MarshalIndent(summary, "", "  ")
	if err := os.WriteFile(summaryPath, summaryData, 0600); err != nil {
		fmt.Printf("  [WARN] Failed to write summary.json: %v\n", err)
	}

	fmt.Println()
	fmt.Printf("Completed in %s: %d successful, %d errors\n", duration.Round(time.Millisecond), successCount, errorCount)

	if errorCount > 0 {
		return fmt.Errorf("%d agent(s) failed", errorCount)
	}

	return nil
}

// Summary holds metadata about the generation run.
type Summary struct {
	InputFile   string          `json:"input_file"`
	GeneratedAt string          `json:"generated_at"`
	Duration    string          `json:"duration"`
	Outputs     []OutputSummary `json:"outputs"`
}

// OutputSummary describes a generated output file.
type OutputSummary struct {
	Agent string `json:"agent"`
	File  string `json:"file"`
}
