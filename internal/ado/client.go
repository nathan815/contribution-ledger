package ado

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles Azure DevOps API interactions
type Client struct {
	org      string
	baseURL  string
	token    string
	httpClient *http.Client
}

// NewClient creates a new ADO client
func NewClient(org, token string) *Client {
	return &Client{
		org:     org,
		baseURL: fmt.Sprintf("https://dev.azure.com/%s/_apis", org),
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// do makes an HTTP request with ADO auth
func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json;api-version=7.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("ADO API error %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// PullRequest represents an ADO PR
type PullRequest struct {
	ID       int    `json:"pullRequestId"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Author   User   `json:"createdBy"`
	Created  time.Time `json:"creationDate"`
	Updated  time.Time `json:"closedDate"`
	Reviewers []Reviewer `json:"reviewers"`
}

// User represents an ADO user
type User struct {
	ID      string `json:"id"`
	Display string `json:"displayName"`
	Email   string `json:"uniqueName"`
}

// Reviewer represents a PR reviewer
type Reviewer struct {
	ID       string `json:"id"`
	Display  string `json:"displayName"`
	Vote     int    `json:"vote"` // 10=approved, 5=approved w/suggestions, 0=no vote, -5=waiting, -10=rejected
	IsPending bool  `json:"isDraft"`
}

// Thread represents a code review comment thread
type Thread struct {
	ID       int       `json:"id"`
	Comments []Comment `json:"comments"`
	Status   string    `json:"status"`
}

// Comment represents a review comment
type Comment struct {
	ID    int    `json:"id"`
	Text  string `json:"content"`
	Author User  `json:"author"`
	Created time.Time `json:"publishedDate"`
}

// PRStats aggregates pull request statistics
type PRStats struct {
	TotalPRsAuthored      int
	TotalPRsReviewed      int
	TotalReviewComments   int
	AverageCommentsPerPR  float64
	ApprovedCount         int
	RequestedChangesCount int
	CommentedCount        int
	ReviewedTimeRange     struct {
		First time.Time
		Last  time.Time
	}
}

// GetPullRequestsForUser fetches PRs authored and reviewed by the user
func (c *Client) GetPullRequestsForUser(userEmail string, status string) ([]PullRequest, error) {
	// Build search criteria
	criteria := fmt.Sprintf("createdBy=\"%s\" AND status=%s", userEmail, status)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/git/pullrequests", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("searchCriteria", criteria)
	q.Add("$top", "100")
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Value []PullRequest `json:"value"`
		Count int           `json:"count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}

// GetPullRequestReviewers fetches all reviewers for a PR
func (c *Client) GetPullRequestReviewers(project, repo string, prID int) ([]Reviewer, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/git/repositories/%s/pullrequests/%d/reviewers", c.baseURL, repo, prID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Value []Reviewer `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}

// GetPullRequestThreads fetches all review comment threads for a PR
func (c *Client) GetPullRequestThreads(project, repo string, prID int) ([]Thread, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/git/repositories/%s/pullrequests/%d/threads", c.baseURL, repo, prID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Value []Thread `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}
