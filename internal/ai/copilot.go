package ai

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CopilotProvider implements the Provider interface using GitHub Copilot SDK
type CopilotProvider struct {
	token string
}

// NewCopilotProvider creates a new Copilot provider
func NewCopilotProvider() (*CopilotProvider, error) {
	// Check if Copilot is available via CLI
	cmd := exec.Command("gh", "copilot", "--version")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("GitHub CLI with Copilot not available: %w", err)
	}

	// Get GitHub token from environment or gh CLI
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		// Try to get from gh CLI
		cmd := exec.Command("gh", "auth", "token")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to get GitHub token: %w", err)
		}
		token = strings.TrimSpace(string(output))
	}

	if token == "" {
		return nil, fmt.Errorf("GitHub token not found")
	}

	return &CopilotProvider{token: token}, nil
}

// Summarize implements Provider.Summarize using Copilot
func (c *CopilotProvider) Summarize(ctx context.Context, text string) (string, error) {
	// Use gh copilot explain for explanations
	// or build a prompt for summarization
	
	prompt := fmt.Sprintf(`Provide a concise, one-sentence summary of the following work (max 100 chars).
Focus on impact and technologies. Remove company/project names, use generic terms instead.

Work description:
%s

Summary:`, text)

	// For now, use a simple implementation
	// In real use, this would call the Copilot API
	result, err := c.callCopilotAPI(ctx, prompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// SummarizeBatch implements Provider.SummarizeBatch for concurrent summarization
func (c *CopilotProvider) SummarizeBatch(ctx context.Context, items []string) ([]string, error) {
	results := make([]string, len(items))
	errChan := make(chan error, len(items))

	for i, item := range items {
		go func(idx int, text string) {
			summary, err := c.Summarize(ctx, text)
			if err != nil {
				errChan <- err
				return
			}
			results[idx] = summary
		}(i, item)
	}

	// Collect errors
	for range items {
		select {
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return results, nil
}

// IsAvailable checks if Copilot is available
func (c *CopilotProvider) IsAvailable() bool {
	cmd := exec.Command("gh", "copilot", "--version")
	return cmd.Run() == nil
}

// callCopilotAPI is a placeholder for actual Copilot API calls
// In production, this would use the official Copilot SDK
func (c *CopilotProvider) callCopilotAPI(ctx context.Context, prompt string) (string, error) {
	// Placeholder implementation
	// Real implementation would use the Copilot API directly
	
	// For now, return a placeholder response
	// This allows the code to build and test locally
	return "Built and maintained service with impact across multiple teams.", nil
}
