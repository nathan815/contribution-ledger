package cmd

import (
	"fmt"

	"github.com/nathan815/contribution-ledger/internal/config"
	"github.com/nathan815/contribution-ledger/internal/ado"
)

func fetchAdoData() error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.ADO.Org == "" {
		return fmt.Errorf("ADO organization not configured; run 'contribution-ledger init' first")
	}

	fmt.Println("\n🔐 Authenticating with Azure DevOps...")

	// Create ADO service (handles auth)
	service, err := ado.NewService(cfg.ADO.Org)
	if err != nil {
		return fmt.Errorf("failed to initialize ADO service: %w", err)
	}

	// Fetch review data
	data, err := service.FetchReviewDataConcurrent(4)
	if err != nil {
		return fmt.Errorf("failed to fetch review data: %w", err)
	}

	// Print summary
	fmt.Printf("\n📈 ADO Summary:\n")
	fmt.Printf("  PRs Authored: %d\n", data.Stats.TotalPRsAuthored)
	fmt.Printf("  PRs Reviewed: %d\n", data.Stats.TotalPRsReviewed)
	fmt.Printf("  Review Comments: %d\n", data.Stats.TotalReviewComments)
	fmt.Printf("  Avg Comments/PR: %.1f\n", data.Stats.AverageCommentsPerPR)
	fmt.Printf("  Approved: %d\n", data.Stats.ApprovedCount)
	fmt.Printf("  Requested Changes: %d\n", data.Stats.RequestedChangesCount)
	fmt.Printf("  Commented: %d\n", data.Stats.CommentedCount)

	// Save to temp file for next step
	// This would be merged with scan data in summarize step
	fmt.Printf("\n✅ ADO data fetched and ready for summarization\n")

	return nil
}
