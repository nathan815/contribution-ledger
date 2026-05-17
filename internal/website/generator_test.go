package website

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGeneratorNew tests generator creation
func TestGeneratorNew(t *testing.T) {
	gen := NewGenerator("/tmp")
	if gen == nil {
		t.Error("generator should not be nil")
	}
	if gen.outputPath != "/tmp" {
		t.Errorf("output path mismatch: got %s", gen.outputPath)
	}
}

// TestGenerateBuildHTML tests HTML building
func TestGenerateBuildHTML(t *testing.T) {
	gen := NewGenerator("/tmp")
	portfolioJSON := `{
		"metadata": {"totalCommits": 1000},
		"projects": [
			{
				"name": "Test Project",
				"aiSummary": "Test summary",
				"stats": {"commits": 500, "filesChanged": 50},
				"technologies": ["Go", "TypeScript"]
			}
		]
	}`

	html := gen.buildHTML(portfolioJSON)

	if html == "" {
		t.Error("HTML should not be empty")
	}
	if !contains(html, "<!DOCTYPE html>") {
		t.Error("HTML should have DOCTYPE")
	}
	if !contains(html, "Test Project") {
		t.Error("HTML should contain project name")
	}
	if !contains(html, portfolioJSON) {
		t.Error("HTML should contain portfolio JSON")
	}
}

// TestGenerateFile tests file generation
func TestGenerateFile(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "website-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create temp portfolio file
	portfolioPath := filepath.Join(tempDir, "portfolio.json")
	portfolioData := `{
		"metadata": {"totalCommits": 1000},
		"projects": [{"name": "Test", "aiSummary": "Summary", "stats": {"commits": 100, "filesChanged": 10}, "technologies": ["Go"]}]
	}`
	if err := os.WriteFile(portfolioPath, []byte(portfolioData), 0644); err != nil {
		t.Fatalf("failed to write portfolio file: %v", err)
	}

	// Generate website
	gen := NewGenerator(tempDir)
	err = gen.Generate(portfolioPath)
	if err != nil {
		t.Errorf("Generate() error = %v", err)
	}

	// Check file exists
	indexPath := filepath.Join(tempDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("index.html should exist")
	}

	// Read and verify content
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Errorf("failed to read index.html: %v", err)
	}

	html := string(content)
	if !contains(html, "Test") {
		t.Error("index.html should contain project data")
	}
}

// Helper function
func contains(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// BenchmarkBuildHTML benchmarks HTML generation
func BenchmarkBuildHTML(b *testing.B) {
	gen := NewGenerator("/tmp")
	portfolioJSON := `{
		"metadata": {"totalCommits": 1000},
		"projects": [
			{"name": "P1", "aiSummary": "S1", "stats": {"commits": 100, "filesChanged": 10}, "technologies": ["Go"]},
			{"name": "P2", "aiSummary": "S2", "stats": {"commits": 200, "filesChanged": 20}, "technologies": ["TypeScript"]}
		]
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.buildHTML(portfolioJSON)
	}
}
