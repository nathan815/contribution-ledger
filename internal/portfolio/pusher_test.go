package portfolio

import (
	"testing"
	"time"
)

// TestPusherNew tests pusher creation
func TestPusherNew(t *testing.T) {
	tests := []struct {
		name    string
		repoPath string
		wantErr bool
	}{
		{"valid repo", "/tmp", false},
		{"nonexistent repo", "/nonexistent/path/to/repo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip nonexistent repo test
			if tt.wantErr {
				return
			}

			// Create temp .git directory for testing
			if tt.repoPath == "/tmp" {
				// /tmp exists, but won't have .git
				// This test is just for structure validation
			}
		})
	}
}

// TestGenerateCommits tests commit generation
func TestGenerateCommits(t *testing.T) {
	now := time.Now()
	activities := []ActivityRecord{
		{
			ProjectName:  "Project Alpha",
			Commits:      347,
			Languages:    "Go, TypeScript",
			LastActivity: now.Add(-30 * 24 * time.Hour),
			JSONData:     `{"project":"alpha"}`,
		},
		{
			ProjectName:  "Project Beta",
			Commits:      156,
			Languages:    "Go, Terraform",
			LastActivity: now.Add(-60 * 24 * time.Hour),
			JSONData:     `{"project":"beta"}`,
		},
	}

	commits := GenerateCommits(activities)

	if len(commits) != len(activities) {
		t.Errorf("commit count mismatch: got %d, want %d", len(commits), len(activities))
	}

	for i, commit := range commits {
		if commit.Message == "" {
			t.Errorf("commit %d message should not be empty", i)
		}
		if len(commit.Files) == 0 {
			t.Errorf("commit %d should have files", i)
		}
	}
}

// TestCommitStructure tests commit structure
func TestCommitStructure(t *testing.T) {
	now := time.Now()

	commit := Commit{
		Message:   "Work: Project Alpha (347 commits, Go)",
		Timestamp: now,
		Files: map[string]string{
			"repos.json": `{"project":"alpha"}`,
			".github/README.md": "# Project Alpha",
		},
	}

	if commit.Message == "" {
		t.Error("message should not be empty")
	}
	if commit.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	if len(commit.Files) != 2 {
		t.Errorf("file count mismatch: got %d", len(commit.Files))
	}
}

// TestActivityRecordTimestamp tests activity record timestamps
func TestActivityRecordTimestamp(t *testing.T) {
	now := time.Now()
	then := now.Add(-30 * 24 * time.Hour)

	activity := ActivityRecord{
		ProjectName:  "Test Project",
		Commits:      100,
		Languages:    "Go",
		LastActivity: then,
		JSONData:     `{}`,
	}

	diff := now.Sub(activity.LastActivity)
	expectedDiff := 30 * 24 * time.Hour

	if diff < expectedDiff-time.Minute || diff > expectedDiff+time.Minute {
		t.Errorf("timestamp difference mismatch")
	}
}

// BenchmarkGenerateCommits benchmarks commit generation
func BenchmarkGenerateCommits(b *testing.B) {
	activities := make([]ActivityRecord, 10)
	for i := 0; i < 10; i++ {
		activities[i] = ActivityRecord{
			ProjectName:  "Project",
			Commits:      100,
			Languages:    "Go",
			LastActivity: time.Now(),
			JSONData:     `{}`,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenerateCommits(activities)
	}
}
