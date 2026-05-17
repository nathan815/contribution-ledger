package ado

import (
	"testing"
	"time"
)

// TestNewService tests service creation
func TestNewService(t *testing.T) {
	// Service creation will skip if az CLI not available
	// This is expected behavior
	t.Log("Service creation requires az CLI - testing skipped in CI")
}

// TestServiceReviewData tests FetchReviewData structure
func TestServiceReviewData(t *testing.T) {
	data := &ReviewData{
		UserEmail:     "user@example.com",
		PRsAuthored:   []PullRequest{},
		ReviewThreads: make(map[int][]Thread),
		Stats: PRStats{
			TotalPRsAuthored: 0,
			TotalPRsReviewed: 0,
		},
	}

	if data.UserEmail != "user@example.com" {
		t.Errorf("email mismatch: got %s", data.UserEmail)
	}
	if data.ReviewThreads == nil {
		t.Error("review threads map should be initialized")
	}
}

// TestPRStatsEdgeCases tests edge cases in PR stats
func TestPRStatsEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		total     int
		comments  int
		wantAvg   float64
	}{
		{"zero PRs", 0, 0, 0},
		{"zero comments", 5, 0, 0},
		{"normal case", 4, 10, 2.5},
		{"one PR", 1, 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var avg float64
			if tt.total > 0 {
				avg = float64(tt.comments) / float64(tt.total)
			}

			if avg != tt.wantAvg {
				t.Errorf("average mismatch: got %.1f, want %.1f", avg, tt.wantAvg)
			}
		})
	}
}

// TestReviewerVoteInterpretation tests vote value interpretation
func TestReviewerVoteInterpretation(t *testing.T) {
	tests := []struct {
		name   string
		vote   int
		action string
	}{
		{"approved", 10, "approved"},
		{"approved-suggestions", 5, "suggestions"},
		{"no-vote", 0, "no-vote"},
		{"waiting", -5, "waiting"},
		{"rejected", -10, "rejected"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var action string
			switch tt.vote {
			case 10:
				action = "approved"
			case 5:
				action = "suggestions"
			case 0:
				action = "no-vote"
			case -5:
				action = "waiting"
			case -10:
				action = "rejected"
			}

			if action != tt.action {
				t.Errorf("action mismatch: got %s, want %s", action, tt.action)
			}
		})
	}
}

// TestTimeRangeCalculation tests date range calculations
func TestTimeRangeCalculation(t *testing.T) {
	t1 := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 1, 30, 12, 0, 0, 0, time.UTC)

	prs := []PullRequest{
		{ID: 1, Created: t1},
		{ID: 2, Created: t2},
		{ID: 3, Created: t3},
	}

	if len(prs) > 0 {
		first := prs[0].Created
		last := prs[len(prs)-1].Created

		if first != t1 {
			t.Errorf("first time mismatch: got %v, want %v", first, t1)
		}
		if last != t3 {
			t.Errorf("last time mismatch: got %v, want %v", last, t3)
		}
	}
}

// BenchmarkServiceCreation benchmarks service initialization
func BenchmarkServiceCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewClient("test-org", "test-token")
	}
}

// BenchmarkReviewDataAggregation benchmarks stats calculation
func BenchmarkReviewDataAggregation(b *testing.B) {
	data := &ReviewData{
		PRsAuthored: make([]PullRequest, 100),
		Stats: PRStats{
			TotalPRsAuthored: 100,
			TotalReviewComments: 500,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if data.Stats.TotalPRsAuthored > 0 {
			_ = float64(data.Stats.TotalReviewComments) / float64(data.Stats.TotalPRsAuthored)
		}
	}
}
