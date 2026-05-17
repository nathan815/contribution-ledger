package ai

import (
	"context"
	"fmt"
)

// Provider defines the interface for AI summarization
type Provider interface {
	Summarize(ctx context.Context, text string) (string, error)
	SummarizeBatch(ctx context.Context, items []string) ([]string, error)
	IsAvailable() bool
}

// Config holds AI provider configuration
type Config struct {
	Provider string // "copilot" or "claude"
	APIKey   string // Optional, for non-Copilot providers
}

// NewProvider creates a new AI provider based on config
func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case "copilot", "":
		// Try Copilot first (uses GitHub token from environment)
		provider, err := NewCopilotProvider()
		if err == nil {
			return provider, nil
		}
		// Fall back to Claude if Copilot unavailable
		fmt.Println("⚠️  Copilot SDK not available, falling back to Claude")
		return NewClaudeProvider(cfg.APIKey)

	case "claude":
		return NewClaudeProvider(cfg.APIKey)

	default:
		return nil, fmt.Errorf("unknown AI provider: %s", cfg.Provider)
	}
}

// TextContent wraps summarization requests
type TextContent struct {
	Title   string
	Content string
}

// SummaryResult holds a summary with metadata
type SummaryResult struct {
	Original string
	Summary  string
	Error    error
}

// PromptTemplate contains reusable prompt templates
type PromptTemplate struct {
	ProjectSummary  string
	ReviewHighlight string
	ImpactStatement string
}

// DefaultPrompts returns standard prompt templates
func DefaultPrompts() PromptTemplate {
	return PromptTemplate{
		ProjectSummary: `Summarize this project work in one concise sentence (max 100 chars). 
Focus on impact and technology. Don't mention company names or internal details.
Commits: {{.Commits}} | Languages: {{.Languages}} | Period: {{.Period}}
{{.Description}}
Summary:`,

		ReviewHighlight: `Create a brief, professional summary of this code review (max 50 chars).
Focus on the type and impact of feedback.
PRs Reviewed: {{.PRCount}} | Comments: {{.Comments}}
Feedback Types: {{.FeedbackTypes}}
Summary:`,

		ImpactStatement: `Generate a one-line impact statement for this accomplishment (max 80 chars).
Make it quantifiable and impressive without being boastful.
Details: {{.Details}}
Impact:`,
	}
}
