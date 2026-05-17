package git

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitStats struct {
	Date         time.Time `json:"date"`
	Language     string    `json:"language"`
	FilesChanged int       `json:"filesChanged"`
	Additions    int       `json:"additions"`
	Deletions    int       `json:"deletions"`
	Message      string    `json:"message"`
}

type ScanResult struct {
	RepoName string         `json:"repoName"`
	Path     string         `json:"path"`
	Commits  []CommitStats  `json:"commits"`
	Stats    RepoStats      `json:"stats"`
}

type RepoStats struct {
	TotalCommits int               `json:"totalCommits"`
	Languages    map[string]int    `json:"languages"`
	FilesChanged int               `json:"filesChanged"`
	Additions    int               `json:"additions"`
	Deletions    int               `json:"deletions"`
	TotalAdded   int               `json:"totalAdded"`
	TotalDeleted int               `json:"totalDeleted"`
	DateRange    struct {
		First time.Time `json:"first"`
		Last  time.Time `json:"last"`
	} `json:"dateRange"`
}

// DetectLanguage detects the primary language of a repository
func DetectLanguage(repoPath string) string {
	// Simple heuristic: look for common file extensions
	extensions := make(map[string]int)
	
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		
		// Skip hidden files and common exclusions
		name := info.Name()
		if strings.HasPrefix(name, ".") {
			return nil
		}
		
		ext := filepath.Ext(name)
		if ext != "" {
			extensions[ext]++
		}
		return nil
	})
	
	if err != nil {
		return "unknown"
	}
	
	// Map extensions to languages
	extToLang := map[string]string{
		".go":    "Go",
		".rs":    "Rust",
		".py":    "Python",
		".ts":    "TypeScript",
		".js":    "JavaScript",
		".java":  "Java",
		".kt":    "Kotlin",
		".cs":    "C#",
		".cpp":   "C++",
		".c":     "C",
		".rb":    "Ruby",
		".php":   "PHP",
		".sh":    "Shell",
		".tf":    "Terraform",
	}
	
	// Find most common extension
	maxCount := 0
	lang := "unknown"
	for ext, count := range extensions {
		if l, ok := extToLang[ext]; ok && count > maxCount {
			maxCount = count
			lang = l
		}
	}
	
	return lang
}

// ScanRepository scans a git repository and extracts commit statistics
func ScanRepository(repoPath string, since time.Time) (*ScanResult, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo: %w", err)
	}
	
	result := &ScanResult{
		RepoName: filepath.Base(repoPath),
		Path:     repoPath,
		Commits:  []CommitStats{},
		Stats: RepoStats{
			Languages: make(map[string]int),
		},
	}
	
	// Get all commits
	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}
	
	lang := DetectLanguage(repoPath)
	
	err = iter.ForEach(func(c *object.Commit) error {
		if c.Author.When.Before(since) {
			return nil
		}
		
		result.Commits = append(result.Commits, CommitStats{
			Date:     c.Author.When,
			Language: lang,
			Message:  strings.TrimSpace(c.Message),
		})
		
		result.Stats.TotalCommits++
		result.Stats.Languages[lang]++
		result.Stats.Additions += 10  // Placeholder
		result.Stats.Deletions += 3   // Placeholder
		
		if result.Stats.DateRange.First.IsZero() {
			result.Stats.DateRange.First = c.Author.When
		}
		result.Stats.DateRange.Last = c.Author.When
		
		return nil
	})
	
	return result, err
}

// ScanDirectory scans all git repositories in a directory
func ScanDirectory(dir string, since time.Time, excludeRepos []string) ([]ScanResult, error) {
	var results []ScanResult
	
	// Expand ~ in path
	if strings.HasPrefix(dir, "~") {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, dir[1:])
	}
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		// Skip excluded
		shouldExclude := false
		for _, exclude := range excludeRepos {
			if entry.Name() == exclude {
				shouldExclude = true
				break
			}
		}
		if shouldExclude {
			continue
		}
		
		repoPath := filepath.Join(dir, entry.Name())
		
		// Check if it's a git repo
		if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
			result, err := ScanRepository(repoPath, since)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to scan %s: %v\n", repoPath, err)
				continue
			}
			results = append(results, *result)
		}
	}
	
	return results, nil
}

// SaveResults saves scan results to JSON file
func SaveResults(results []ScanResult, outputPath string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	
	return os.WriteFile(outputPath, data, 0644)
}

// LoadResults loads scan results from JSON file
func LoadResults(inputPath string) ([]ScanResult, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	var results []ScanResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal results: %w", err)
	}
	
	return results, nil
}
