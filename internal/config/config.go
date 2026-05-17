package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	ADO struct {
		Org      string   `json:"org"`
		Projects []string `json:"projects"`
	} `json:"ado"`

	PortfolioRepo struct {
		Owner string `json:"owner"`
		Name  string `json:"name"`
		Path  string `json:"path"`
	} `json:"portfolioRepo"`

	ScanPaths   []string `json:"scanPaths"`
	ExcludeRepos []string `json:"excludeRepos"`

	Anonymization struct {
		ProjectNames map[string]string `json:"projectNames"`
		CompanyName  string            `json:"companyName"`
		StripFileNames bool            `json:"stripFileNames"`
	} `json:"anonymization"`
}

var (
	ErrConfigNotFound = errors.New("configuration file not found; run 'contribution-ledger init' first")
)

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".contribution-ledger", "config.json"), nil
}

// Load reads the config file
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the config file
func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// Default returns a default configuration
func Default() *Config {
	cfg := &Config{}
	cfg.ScanPaths = []string{"~/dev"}
	cfg.ExcludeRepos = []string{"node_modules", ".git", "terraform", "vendor"}
	cfg.Anonymization.CompanyName = "Work"
	cfg.Anonymization.ProjectNames = make(map[string]string)
	cfg.Anonymization.StripFileNames = true
	return cfg
}
