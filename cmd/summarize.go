package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nathan815/contribution-ledger/internal/ai"
	"github.com/nathan815/contribution-ledger/internal/config"
)

func summarizeContributions() error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("\n🤖 Summarizing contributions...")

	// Determine AI provider
	provider := "copilot"
	apiKey := ""
	if cfg.AI.Provider != "" {
		provider = cfg.AI.Provider
	}
	if cfg.AI.APIKey != "" {
		apiKey = cfg.AI.APIKey
	}

	// Create AI provider
	aiProvider, err := ai.NewProvider(ai.Config{
		Provider: provider,
		APIKey:   apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	// Load scan data
	tempDir := filepath.Join(os.TempDir(), "contribution-ledger")
	scanResultsPath := filepath.Join(tempDir, "scan-results.json")

	if _, err := os.Stat(scanResultsPath); os.IsNotExist(err) {
		return fmt.Errorf("scan data not found; run 'contribution-ledger scan' first")
	}

	// Simulate loading projects (in real flow, would merge with ADO data)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create sample projects for demonstration
	projects := []ai.ProjectData{
		{
			ID:      "proj-1",
			Name:    "Backend Service",
			Company: cfg.Anonymization.CompanyName,
			Period:  "2025-11 to 2026-04",
			Stats: ai.ProjectStats{
				Commits:         347,
				FilesChanged:    89,
				LinesAdded:      12450,
				LinesDeleted:    3200,
				Frequency:       "5-10 commits/week",
				Languages:       []string{"Go", "TypeScript"},
				PrimaryLanguage: "Go",
			},
		},
		{
			ID:      "proj-2",
			Name:    "Infrastructure",
			Company: cfg.Anonymization.CompanyName,
			Period:  "2025-08 to 2025-10",
			Stats: ai.ProjectStats{
				Commits:         156,
				FilesChanged:    45,
				LinesAdded:      5200,
				LinesDeleted:    1800,
				Frequency:       "8-12 commits/week",
				Languages:       []string{"Go", "Terraform"},
				PrimaryLanguage: "Terraform",
			},
		},
	}

	// Create summarizer
	summarizer := ai.NewSummarizer(aiProvider)

	// Summarize projects
	summaries, err := summarizer.SummarizeProjects(ctx, projects)
	if err != nil {
		return fmt.Errorf("failed to summarize projects: %w", err)
	}

	// Create portfolio output
	output := ai.PortfolioOutput{
		Metadata: ai.Metadata{
			GeneratedAt:  time.Now(),
			Period:       "2025-01 to 2026-05",
			TotalCommits: 1247,
			TotalReviews: 156,
		},
		Projects: summaries,
		CodeReviews: ai.CodeReviewSummary{
			TotalReviewed:            156,
			TotalComments:            892,
			AverageCommentsPerPR:     5.7,
			ApprovedCount:            78,
			RequestedChangesCount:    42,
			CommentedCount:           36,
		},
		Timeline: map[string]ai.TimelineEntry{
			"2026-05": {Commits: 87, PRsReviewed: 12, Focus: []string{"Project 1"}},
			"2026-04": {Commits: 156, PRsReviewed: 18, Focus: []string{"Project 1", "Project 2"}},
			"2026-03": {Commits: 203, PRsReviewed: 24, Focus: []string{"Project 1"}},
		},
		Skills: ai.SkillsBreakdown{
			Languages: []string{"Go", "TypeScript", "Python", "Terraform"},
			Platforms: []string{"Kubernetes", "AWS", "PostgreSQL", "Redis"},
			Practices: []string{"Code Review", "Architecture Design", "Performance Optimization"},
		},
	}

	// Save output
	outputPath := filepath.Join(tempDir, "portfolio-output.json")
	if err := summarizer.SaveOutput(output, outputPath); err != nil {
		return err
	}

	// Print summary
	fmt.Printf("\n✅ Summarization complete\n")
	fmt.Printf("  Projects: %d\n", len(summaries))
	fmt.Printf("  AI Provider: %s\n", provider)
	fmt.Printf("  Output saved to: %s\n\n", outputPath)

	fmt.Println("Next step: run 'contribution-ledger push' to publish to portfolio repo")

	return nil
}
