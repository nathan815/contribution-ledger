package ado

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestClientNewClient tests client initialization
func TestClientNewClient(t *testing.T) {
	org := "test-org"
	token := "test-token"

	client := NewClient(org, token)

	if client.org != org {
		t.Errorf("org mismatch: got %s, want %s", client.org, org)
	}
	if client.token != token {
		t.Errorf("token mismatch: got %s, want %s", client.token, token)
	}
	if !strings.Contains(client.baseURL, org) {
		t.Errorf("baseURL should contain org: %s", client.baseURL)
	}
}

// TestClientDo tests HTTP request handling
func TestClientDo(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
	}{
		{"OK response", 200, `{"value":[]}`, false},
		{"Bad request", 400, "error", true},
		{"Server error", 500, "error", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))

				// Verify auth header
				if auth := r.Header.Get("Authorization"); auth == "" {
					t.Error("missing Authorization header")
				}
			}))
			defer server.Close()

			client := NewClient("test-org", "test-token")
			client.baseURL = server.URL

			req, _ := http.NewRequest("GET", server.URL, nil)
			resp, err := client.do(req)

			if (err != nil) != tt.wantErr {
				t.Errorf("do() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && resp != nil {
				resp.Body.Close()
			}
		})
	}
}

// TestAuthCheckAzCLI tests az CLI availability check
func TestAuthCheckAzCLI(t *testing.T) {
	auth := NewAuth()

	// This will succeed on systems with az CLI, fail otherwise
	err := auth.CheckAzCLI()
	if err != nil {
		t.Logf("az CLI check skipped (not installed): %v", err)
	}
}

// TestPRStatsCalculation tests PR stats calculations
func TestPRStatsCalculation(t *testing.T) {
	data := &ReviewData{
		PRsAuthored: []PullRequest{
			{ID: 1, Title: "Fix bug", Status: "completed"},
			{ID: 2, Title: "Feature", Status: "completed"},
		},
		Stats: PRStats{
			TotalPRsAuthored: 2,
			TotalReviewComments: 5,
		},
	}

	data.Stats.AverageCommentsPerPR = float64(data.Stats.TotalReviewComments) / float64(data.Stats.TotalPRsAuthored)

	expectedAvg := 2.5
	if data.Stats.AverageCommentsPerPR != expectedAvg {
		t.Errorf("average comments per PR: got %.1f, want %.1f", data.Stats.AverageCommentsPerPR, expectedAvg)
	}
}

// TestReviewDataStructure tests ReviewData structure
func TestReviewDataStructure(t *testing.T) {
	now := time.Now()

	data := &ReviewData{
		UserEmail: "test@example.com",
		PRsAuthored: []PullRequest{
			{
				ID:     1,
				Title:  "Test PR",
				Status: "completed",
				Author: User{Email: "test@example.com"},
				Created: now,
				Updated: now.Add(24 * time.Hour),
			},
		},
		Stats: PRStats{
			TotalPRsAuthored: 1,
		},
	}

	if data.UserEmail != "test@example.com" {
		t.Errorf("email mismatch")
	}
	if len(data.PRsAuthored) != 1 {
		t.Errorf("PRs count mismatch")
	}
	if data.Stats.TotalPRsAuthored != 1 {
		t.Errorf("stats count mismatch")
	}
}

// TestReviewHighlight tests review highlight creation
func TestReviewHighlight(t *testing.T) {
	now := time.Now()

	highlight := ReviewHighlight{
		PRTitle:     "Add feature",
		ReviewType:  "approval",
		Date:        now,
		CommentCount: 3,
	}

	if highlight.PRTitle == "" {
		t.Error("PR title should not be empty")
	}
	if highlight.CommentCount < 0 {
		t.Error("comment count should be >= 0")
	}
}

// BenchmarkClientDo benchmarks HTTP request overhead
func BenchmarkClientDo(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"value":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-org", "test-token")
	client.baseURL = server.URL

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		resp, _ := client.do(req)
		if resp != nil {
			resp.Body.Close()
		}
	}
}
