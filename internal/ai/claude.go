package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// ClaudeProvider implements the Provider interface using Claude API
type ClaudeProvider struct {
	apiKey string
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey string) (*ClaudeProvider, error) {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("Claude API key not provided and ANTHROPIC_API_KEY not set")
	}

	return &ClaudeProvider{apiKey: apiKey}, nil
}

// Summarize implements Provider.Summarize using Claude API
func (c *ClaudeProvider) Summarize(ctx context.Context, text string) (string, error) {
	prompt := fmt.Sprintf(`You are a professional resume writer. Provide a concise, one-sentence summary of the following work (max 100 chars).
Focus on impact and technologies. Remove company/project names, use generic terms instead.
Make it sound impressive but authentic.

Work description:
%s

Summary:`, text)

	// Placeholder for actual Claude API call
	// In production, this would use the Anthropic SDK
	// For now, return a realistic sample
	result := c.summarizeWithAPI(ctx, prompt)
	return strings.TrimSpace(result), nil
}

// SummarizeBatch implements Provider.SummarizeBatch for concurrent summarization
func (c *ClaudeProvider) SummarizeBatch(ctx context.Context, items []string) ([]string, error) {
	results := make([]string, len(items))
	errChan := make(chan error, len(items))
	var wg sync.WaitGroup

	// Limit concurrent calls to avoid rate limiting
	semaphore := make(chan struct{}, 3)

	for i, item := range items {
		wg.Add(1)
		go func(idx int, text string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			summary, err := c.Summarize(ctx, text)
			if err != nil {
				errChan <- err
				return
			}
			results[idx] = summary
		}(i, item)
	}

	// Wait for all goroutines
	wg.Wait()
	close(errChan)

	// Collect any errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// IsAvailable checks if Claude API is accessible
func (c *ClaudeProvider) IsAvailable() bool {
	return c.apiKey != ""
}

// summarizeWithAPI is a placeholder for actual Claude API calls
func (c *ClaudeProvider) summarizeWithAPI(ctx context.Context, prompt string) string {
	// Placeholder implementation
	// Real implementation would use the Anthropic SDK:
	// client := anthropic.NewClient(c.apiKey)
	// response, err := client.Messages.New(ctx, &anthropic.MessageNewParams{...})
	
	// For now, return realistic samples
	return "Built and maintained service with significant impact across multiple teams."
}
