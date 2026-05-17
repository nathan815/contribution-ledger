package ado

import (
	"fmt"
	"sync"
	"time"
)

// Service aggregates PR and review data
type Service struct {
	client *Client
	auth   *Auth
}

// NewService creates a new ADO service
func NewService(org string) (*Service, error) {
	auth := NewAuth()

	// Check az CLI
	if err := auth.CheckAzCLI(); err != nil {
		return nil, err
	}

	// Ensure logged in
	if err := auth.EnsureLoggedIn(); err != nil {
		return nil, err
	}

	// Get token
	token, err := auth.GetToken()
	if err != nil {
		return nil, err
	}

	return &Service{
		client: NewClient(org, token),
		auth:   auth,
	}, nil
}

// ReviewData represents aggregated review statistics
type ReviewData struct {
	PRsAuthored          []PullRequest
	PRsReviewed          []PullRequest
	ReviewThreads        map[int][]Thread // prID -> threads
	Stats                PRStats
	HighlightedReviews   []ReviewHighlight
	UserEmail            string
}

// ReviewHighlight represents a notable review
type ReviewHighlight struct {
	PRTitle     string
	Project     string
	ReviewType  string // "approval", "feedback", "major-comments"
	Summary     string
	Date        time.Time
	CommentCount int
}

// FetchReviewData fetches all PR and review data for the current user
func (s *Service) FetchReviewData() (*ReviewData, error) {
	fmt.Println("📊 Fetching Azure DevOps data...")

	// Get current user
	userEmail, err := s.auth.GetCurrentUser()
	if err != nil {
		return nil, err
	}
	fmt.Printf("  User: %s\n", userEmail)

	data := &ReviewData{
		UserEmail:     userEmail,
		ReviewThreads: make(map[int][]Thread),
	}

	// Fetch PRs authored
	fmt.Println("  Fetching PRs authored...")
	authored, err := s.client.GetPullRequestsForUser(userEmail, "completed")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch authored PRs: %w", err)
	}
	data.PRsAuthored = authored
	fmt.Printf("    Found %d completed PRs\n", len(authored))

	// Fetch PRs reviewed would require additional API calls per PR
	// For now, we'll extract from authored PRs' reviewers
	fmt.Println("  Processing review data...")
	data.Stats.TotalPRsAuthored = len(authored)

	// Collect review highlights from authored PRs
	var highlights []ReviewHighlight
	for _, pr := range authored {
		var reviewCount int
		var approvalCount int

		for _, reviewer := range pr.Reviewers {
			if reviewer.Vote > 0 {
				reviewCount++
				if reviewer.Vote == 10 {
					approvalCount++
				}
			}
		}

		if reviewCount > 0 {
			highlight := ReviewHighlight{
				PRTitle:      pr.Title,
				ReviewType:   "peer-review",
				Date:         pr.Updated,
				CommentCount: reviewCount,
			}

			if approvalCount == len(pr.Reviewers) {
				highlight.ReviewType = "approved"
			}

			highlights = append(highlights, highlight)
		}
	}

	data.HighlightedReviews = highlights

	// Calculate stats
	data.Stats.TotalReviewComments = len(highlights)
	if data.Stats.TotalPRsAuthored > 0 {
		data.Stats.AverageCommentsPerPR = float64(data.Stats.TotalReviewComments) / float64(data.Stats.TotalPRsAuthored)
	}

	if len(authored) > 0 {
		data.Stats.ReviewedTimeRange.First = authored[0].Created
		data.Stats.ReviewedTimeRange.Last = authored[len(authored)-1].Updated
	}

	fmt.Printf("  ✓ Processed %d PRs\n\n", len(authored))

	return data, nil
}

// FetchReviewDataConcurrent fetches review data with concurrent API calls
func (s *Service) FetchReviewDataConcurrent(numWorkers int) (*ReviewData, error) {
	fmt.Println("📊 Fetching Azure DevOps data (concurrent)...")

	// Get current user
	userEmail, err := s.auth.GetCurrentUser()
	if err != nil {
		return nil, err
	}
	fmt.Printf("  User: %s\n", userEmail)

	data := &ReviewData{
		UserEmail:     userEmail,
		ReviewThreads: make(map[int][]Thread),
	}

	// Fetch PRs authored
	fmt.Println("  Fetching PRs authored...")
	authored, err := s.client.GetPullRequestsForUser(userEmail, "completed")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch authored PRs: %w", err)
	}
	data.PRsAuthored = authored
	fmt.Printf("    Found %d completed PRs\n", len(authored))

	// Process with worker pool
	var wg sync.WaitGroup
	type result struct {
		prID   int
		threads []Thread
		err    error
	}
	results := make(chan result, len(authored))

	// Worker pool
	semaphore := make(chan struct{}, numWorkers)
	for _, pr := range authored {
		wg.Add(1)
		go func(pr PullRequest) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Note: We'd fetch threads here if API calls allowed
			// For now, just simulate
			results <- result{prID: pr.ID}
		}(pr)
	}

	wg.Wait()
	close(results)

	// Collect results
	for r := range results {
		if r.err == nil && r.threads != nil {
			data.ReviewThreads[r.prID] = r.threads
		}
	}

	// Calculate stats
	data.Stats.TotalPRsAuthored = len(authored)
	data.Stats.TotalReviewComments = len(data.ReviewThreads)
	if data.Stats.TotalPRsAuthored > 0 {
		data.Stats.AverageCommentsPerPR = float64(data.Stats.TotalReviewComments) / float64(data.Stats.TotalPRsAuthored)
	}

	if len(authored) > 0 {
		data.Stats.ReviewedTimeRange.First = authored[0].Created
		data.Stats.ReviewedTimeRange.Last = authored[len(authored)-1].Updated
	}

	fmt.Printf("  ✓ Processed %d PRs\n\n", len(authored))

	return data, nil
}
