package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nathan815/contribution-ledger/internal/config"
	"github.com/nathan815/contribution-ledger/internal/website"
)

// generateWebsite generates the portfolio website
func generateWebsite() error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("\n🌐 Generating portfolio website...")

	// Load portfolio output
	tempDir := filepath.Join(os.TempDir(), "contribution-ledger")
	outputPath := filepath.Join(tempDir, "portfolio-output.json")

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("portfolio output not found; run 'contribution-ledger summarize' first")
	}

	// Generate website in portfolio repo
	repoPath := cfg.PortfolioRepo.Path
	if repoPath == "" {
		repoPath = "."
	}

	gen := website.NewGenerator(repoPath)
	if err := gen.Generate(outputPath); err != nil {
		return fmt.Errorf("failed to generate website: %w", err)
	}

	fmt.Printf("📄 Website generated at: %s/index.html\n", repoPath)
	fmt.Println("\nYou can now:")
	fmt.Println("  1. Commit the changes: git add . && git commit -m 'Update portfolio'")
	fmt.Println("  2. Push to GitHub: git push origin main")
	fmt.Println("  3. Visit: https://<username>.github.io/contribution-ledger-data/")

	return nil
}
