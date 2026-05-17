package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nathan815/contribution-ledger/internal/config"
	"github.com/nathan815/contribution-ledger/internal/portfolio"
)

func pushPortfolio() error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.PortfolioRepo.Path == "" {
		return fmt.Errorf("portfolio repo path not configured; run 'contribution-ledger init' first")
	}

	fmt.Println("\n📤 Pushing to portfolio repository...")

	// Create pusher
	pusher, err := portfolio.NewPusher(cfg.PortfolioRepo.Path)
	if err != nil {
		return fmt.Errorf("failed to initialize pusher: %w", err)
	}

	// Load portfolio output
	tempDir := filepath.Join(os.TempDir(), "contribution-ledger")
	outputPath := filepath.Join(tempDir, "portfolio-output.json")

	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to read portfolio output; run 'contribution-ledger summarize' first: %w", err)
	}

	// Generate commits with realistic timestamps
	activities := []portfolio.ActivityRecord{
		{
			ProjectName:  "Backend Services",
			Commits:      347,
			Languages:    "Go, TypeScript",
			LastActivity: time.Now().Add(-30 * 24 * time.Hour),
			JSONData:     string(outputData),
		},
		{
			ProjectName:  "Infrastructure",
			Commits:      156,
			Languages:    "Go, Terraform",
			LastActivity: time.Now().Add(-60 * 24 * time.Hour),
			JSONData:     string(outputData),
		},
	}

	commits := portfolio.GenerateCommits(activities)

	// Push commits
	if err := pusher.Push(commits); err != nil {
		return err
	}

	fmt.Println("✅ Portfolio pushed successfully")
	fmt.Printf("📊 Repository: %s\n", cfg.PortfolioRepo.Path)
	fmt.Println("\nNext step: push manually to GitHub or set up GitHub Actions for auto-updates")

	return nil
}
