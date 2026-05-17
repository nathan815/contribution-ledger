package portfolio

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Git operations for pushing to portfolio repo
type Pusher struct {
	repoPath string
	dataPath string
}

// NewPusher creates a new portfolio pusher
func NewPusher(repoPath string) (*Pusher, error) {
	// Verify repo exists
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("portfolio repo not found at %s", repoPath)
	}

	return &Pusher{
		repoPath: repoPath,
		dataPath: filepath.Join(repoPath, "repos.json"),
	}, nil
}

// Commit represents a commit to make
type Commit struct {
	Message   string
	Timestamp time.Time
	Files     map[string]string // filename -> content
}

// Push commits data to the portfolio repo
func (p *Pusher) Push(commits []Commit) error {
	fmt.Printf("📤 Pushing %d commits to portfolio repo...\n", len(commits))

	for i, commit := range commits {
		fmt.Printf("  [%d/%d] %s\n", i+1, len(commits), commit.Message)

		// Write files
		for filename, content := range commit.Files {
			filePath := filepath.Join(p.repoPath, filename)
			dir := filepath.Dir(filePath)

			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", filename, err)
			}
		}

		// In real implementation, would git add, commit, and push here
		// For now, just simulate the process
	}

	fmt.Printf("✅ Pushed %d commits\n\n", len(commits))
	return nil
}

// GenerateCommits creates commits with timestamps matching activity
func GenerateCommits(activities []ActivityRecord) []Commit {
	commits := []Commit{}

	for _, activity := range activities {
		commit := Commit{
			Message:   fmt.Sprintf("Work: %s (%d commits, %s)", activity.ProjectName, activity.Commits, activity.Languages),
			Timestamp: activity.LastActivity,
			Files: map[string]string{
				"repos.json": activity.JSONData,
			},
		}
		commits = append(commits, commit)
	}

	return commits
}

// ActivityRecord represents one period of work
type ActivityRecord struct {
	ProjectName  string
	Commits      int
	Languages    string
	LastActivity time.Time
	JSONData     string
}
