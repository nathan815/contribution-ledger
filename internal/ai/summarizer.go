package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Summarizer orchestrates the summarization process
type Summarizer struct {
	provider Provider
}

// NewSummarizer creates a new summarizer with the given provider
func NewSummarizer(provider Provider) *Summarizer {
	return &Summarizer{provider: provider}
}

// ProjectSummary represents a summarized project
type ProjectSummary struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Company       string       `json:"company"`
	Period        string       `json:"period"`
	Type          string       `json:"type"` // backend, frontend, devops, etc.
	Stats         ProjectStats `json:"stats"`
	AISummary     string       `json:"aiSummary"`
	Technologies  []string     `json:"technologies"`
	Impact        string       `json:"impact"`
}

// ProjectStats represents project metrics
type ProjectStats struct {
	Commits         int      `json:"commits"`
	FilesChanged    int      `json:"filesChanged"`
	LinesAdded      int      `json:"linesAdded"`
	LinesDeleted    int      `json:"linesDeleted"`
	Frequency       string   `json:"frequency"` // e.g., "5-10 commits/week"
	Languages       []string `json:"languages"`
	PrimaryLanguage string   `json:"primaryLanguage"`
}

// CodeReviewSummary represents summarized review metrics
type CodeReviewSummary struct {
	TotalReviewed            int      `json:"totalReviewed"`
	TotalComments            int      `json:"totalComments"`
	AverageCommentsPerPR     float64  `json:"averageCommentsPerPR"`
	ApprovedCount            int      `json:"approvedCount"`
	RequestedChangesCount    int      `json:"requestedChangesCount"`
	CommentedCount           int      `json:"commentedCount"`
	HighlightedReviews       []string `json:"highlightedReviews"`
}

// PortfolioOutput is the final JSON output
type PortfolioOutput struct {
	Metadata   Metadata                  `json:"metadata"`
	Projects   []ProjectSummary          `json:"projects"`
	CodeReviews CodeReviewSummary        `json:"codeReviews"`
	Timeline   map[string]TimelineEntry  `json:"timeline"`
	Skills     SkillsBreakdown           `json:"skills"`
}

// Metadata holds generation metadata
type Metadata struct {
	GeneratedAt  time.Time `json:"generatedAt"`
	Period       string    `json:"period"`
	TotalCommits int       `json:"totalCommits"`
	TotalReviews int       `json:"totalCodeReviewsLeft"`
}

// TimelineEntry represents activity in a month
type TimelineEntry struct {
	Commits     int      `json:"commits"`
	PRsReviewed int      `json:"prsReviewed"`
	Focus       []string `json:"focus"`
}

// SkillsBreakdown represents aggregated skills
type SkillsBreakdown struct {
	Languages []string `json:"languages"`
	Platforms []string `json:"platforms"`
	Practices []string `json:"practices"`
}

// SummarizeProjects summarizes all projects concurrently
func (s *Summarizer) SummarizeProjects(ctx context.Context, projects []ProjectData) ([]ProjectSummary, error) {
	fmt.Printf("🤖 Summarizing %d projects...\n", len(projects))

	results := make([]ProjectSummary, len(projects))
	errChan := make(chan error, len(projects))

	for i, proj := range projects {
		go func(idx int, p ProjectData) {
			summary, err := s.summarizeProject(ctx, p)
			if err != nil {
				errChan <- fmt.Errorf("failed to summarize project %s: %w", p.Name, err)
				return
			}
			results[idx] = summary
		}(i, proj)
	}

	// Collect errors
	for range projects {
		select {
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	fmt.Printf("✅ Summarized %d projects\n", len(results))
	return results, nil
}

// summarizeProject summarizes a single project
func (s *Summarizer) summarizeProject(ctx context.Context, proj ProjectData) (ProjectSummary, error) {
	desc := fmt.Sprintf(
		"Project: %s\nCommits: %d\nLanguages: %v\nPeriod: %s\nStats: %d files changed, +%d-%d lines",
		proj.Name,
		proj.Stats.Commits,
		proj.Stats.Languages,
		proj.Period,
		proj.Stats.FilesChanged,
		proj.Stats.LinesAdded,
		proj.Stats.LinesDeleted,
	)

	summary, err := s.provider.Summarize(ctx, desc)
	if err != nil {
		return ProjectSummary{}, err
	}

	return ProjectSummary{
		ID:            proj.ID,
		Name:          proj.Name,
		Company:       proj.Company,
		Period:        proj.Period,
		Type:          inferProjectType(proj.Stats.Languages),
		Stats:         proj.Stats,
		AISummary:     summary,
		Technologies: proj.Stats.Languages,
		Impact:        inferImpact(proj.Stats),
	}, nil
}

// SaveOutput saves the portfolio to a JSON file
func (s *Summarizer) SaveOutput(output PortfolioOutput, outputPath string) error {
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("📄 Portfolio saved to: %s\n", outputPath)
	return nil
}

// ProjectData represents raw project data to summarize
type ProjectData struct {
	ID      string
	Name    string
	Company string
	Period  string
	Stats   ProjectStats
}

// inferProjectType determines project type from languages
func inferProjectType(languages []string) string {
	if len(languages) == 0 {
		return "general"
	}

	backendLangs := map[string]bool{"Go": true, "Python": true, "Java": true, "Rust": true, "C#": true}
	frontendLangs := map[string]bool{"TypeScript": true, "JavaScript": true, "Vue": true, "React": true}
	infraLangs := map[string]bool{"Terraform": true, "HCL": true, "Shell": true, "YAML": true}

	for _, lang := range languages {
		if backendLangs[lang] {
			return "backend"
		}
		if frontendLangs[lang] {
			return "frontend"
		}
		if infraLangs[lang] {
			return "infrastructure"
		}
	}

	return "general"
}

// inferImpact generates an impact statement from stats
func inferImpact(stats ProjectStats) string {
	if stats.Commits > 500 {
		return "Led significant development effort with substantial contribution"
	}
	if stats.Commits > 200 {
		return "Drove core features and improvements"
	}
	if stats.Commits > 100 {
		return "Made meaningful contributions to project"
	}
	return "Contributed to project development"
}
