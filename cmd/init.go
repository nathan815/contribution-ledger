package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nathan815/contribution-ledger/internal/config"
)

func initConfig() error {
	cfg := config.Default()
	
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("\n📊 Contribution Ledger Configuration")
	fmt.Println("=====================================\n")
	
	// ADO Organization
	fmt.Print("Azure DevOps Organization URL (e.g., dev.azure.com/your-org): ")
	org, _ := reader.ReadString('\n')
	org = strings.TrimSpace(org)
	if org != "" {
		cfg.ADO.Org = org
	}
	
	// Portfolio Repo Owner
	fmt.Print("\nGitHub Username (for portfolio repo): ")
	owner, _ := reader.ReadString('\n')
	owner = strings.TrimSpace(owner)
	if owner != "" {
		cfg.PortfolioRepo.Owner = owner
	}
	
	// Portfolio Repo Name
	fmt.Print("Portfolio Repo Name [contribution-ledger-data]: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "contribution-ledger-data"
	}
	cfg.PortfolioRepo.Name = name
	
	// Portfolio Repo Path
	fmt.Print("Portfolio Repo Path (e.g., ~/code/contribution-ledger-data): ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)
	if path != "" {
		cfg.PortfolioRepo.Path = path
	}
	
	// Scan Paths
	fmt.Print("\nDirectories to scan [~/dev]: ")
	scanPath, _ := reader.ReadString('\n')
	scanPath = strings.TrimSpace(scanPath)
	if scanPath != "" {
		cfg.ScanPaths = strings.Split(scanPath, ",")
		for i := range cfg.ScanPaths {
			cfg.ScanPaths[i] = strings.TrimSpace(cfg.ScanPaths[i])
		}
	}
	
	// Company Name
	fmt.Print("\nCompany Name (for anonymization) [Work]: ")
	company, _ := reader.ReadString('\n')
	company = strings.TrimSpace(company)
	if company != "" {
		cfg.Anonymization.CompanyName = company
	}
	
	// Save
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	path, _ = config.ConfigPath()
	fmt.Printf("\n✅ Configuration saved to %s\n", path)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run 'contribution-ledger scan' to scan your repositories")
	fmt.Println("  2. Run 'contribution-ledger ado-pull-requests' to fetch ADO data")
	fmt.Println("  3. Run 'contribution-ledger summarize' to generate AI summaries")
	fmt.Println("  4. Run 'contribution-ledger push' to publish to your portfolio repo")
	
	return nil
}
