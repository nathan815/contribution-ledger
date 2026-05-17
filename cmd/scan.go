package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/nathan815/contribution-ledger/internal/config"
	gitpkg "github.com/nathan815/contribution-ledger/internal/git"
)

func scanRepos() error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	
	fmt.Println("\n📁 Scanning repositories...")
	
	var since time.Time
	// Scan last 12 months by default
	since = time.Now().AddDate(-1, 0, 0)
	
	var allResults []gitpkg.ScanResult
	
	for _, scanPath := range cfg.ScanPaths {
		fmt.Printf("\n  Scanning: %s\n", scanPath)
		
		results, err := gitpkg.ScanDirectory(scanPath, since, cfg.ExcludeRepos)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠️  Error scanning %s: %v\n", scanPath, err)
			continue
		}
		
		allResults = append(allResults, results...)
		fmt.Printf("  ✓ Found %d repositories\n", len(results))
	}
	
	// Save results to temp file
	outputDir := filepath.Join(os.TempDir(), "contribution-ledger")
	os.MkdirAll(outputDir, 0755)
	
	outputPath := filepath.Join(outputDir, "scan-results.json")
	if err := gitpkg.SaveResults(allResults, outputPath); err != nil {
		return err
	}
	
	// Print summary
	totalCommits := 0
	totalRepos := 0
	for _, r := range allResults {
		totalCommits += r.Stats.TotalCommits
		if r.Stats.TotalCommits > 0 {
			totalRepos++
		}
	}
	
	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("  Total Repositories: %d\n", totalRepos)
	fmt.Printf("  Total Commits: %d\n", totalCommits)
	fmt.Printf("  Results saved to: %s\n\n", outputPath)
	
	return nil
}
