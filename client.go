package backlog

import (
	"net/http"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/option"
)

// ──────────────────────────────────────────────────────────────
//  Doer interface (HTTP abstraction)
// ──────────────────────────────────────────────────────────────

// Doer defines the minimal interface required to perform HTTP requests.
// It is compatible with *http.Client and allows injection of mock clients
// for unit or integration testing.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ──────────────────────────────────────────────────────────────
//  Client structure and initialization
// ──────────────────────────────────────────────────────────────

// Client represents a Backlog API client.
// It wraps an underlying HTTP Doer and provides typed services for API access.
type Client struct {
	httpClient *client.Client

	// Issue provides access to issue-related API endpoints.
	Issue *IssueService
	// Project provides access to project-related API endpoints.
	Project *ProjectService
	// PullRequest provides access to pull request-related API endpoints.
	PullRequest *PullRequestService
	// RecentlyViewed provides access to recently viewed resource endpoints.
	RecentlyViewed *RecentlyViewedService
	// Repository provides access to Git repository endpoints.
	Repository *RepositoryService
	// Space provides access to space-related API endpoints.
	Space *SpaceService
	// Star provides access to star-related API endpoints.
	Star *StarService
	// User provides access to user-related API endpoints.
	User *UserService
	// Wiki provides access to wiki-related API endpoints.
	Wiki *WikiService
}

// ──────────────────────────────────────────────────────────────
//  Client constructor
// ──────────────────────────────────────────────────────────────

// NewClient creates and initializes a Backlog API Client.
// It requires a baseURL (e.g. "https://example.backlog.com") and an API token.
//
// It returns an [*InternalClientError] if the base URL or token is invalid.
//
// Supported options:
//   - [WithDoer]
func NewClient(baseURL, token string, opts ...*ClientOption) (*Client, error) {
	clientOpts := make([]*client.ClientOption, len(opts))
	for i, o := range opts {
		clientOpts[i] = o.inner
	}
	c, err := client.NewClient(baseURL, token, clientOpts...)
	if err != nil {
		return nil, convertError(err)
	}

	bc := &Client{
		httpClient: c,
	}

	initServices(bc)

	return bc, nil
}

// ──────────────────────────────────────────────────────────────
//  Service initialization
// ──────────────────────────────────────────────────────────────

func initServices(c *Client) {
	baseOptionService := &option.OptionService{}

	c.Issue = newIssueService(c.httpClient.Method, baseOptionService)

	c.Project = newProjectService(c.httpClient.Method, baseOptionService)

	c.PullRequest = newPullRequestService(c.httpClient.Method, baseOptionService)

	c.RecentlyViewed = newRecentlyViewedService(c.httpClient.Method, baseOptionService)

	c.Repository = newRepositoryService(c.httpClient.Method)

	c.Space = newSpaceService(c.httpClient.Method, baseOptionService)

	c.Star = newStarService(c.httpClient.Method, baseOptionService)

	c.User = newUserService(c.httpClient.Method, baseOptionService)

	c.Wiki = newWikiService(c.httpClient.Method, baseOptionService)
}

// ──────────────────────────────────────────────────────────────
//  Client options
// ──────────────────────────────────────────────────────────────

// ClientOption defines a functional option for configuring a Client.
// It is used to change the default behavior of the Client.
type ClientOption struct {
	inner *client.ClientOption
}

// WithDoer returns a ClientOption that sets the HTTP client (Doer) for the Client.
// This is useful for providing a custom *http.Client or a mock implementation during testing.
//
// If this option is not provided, http.DefaultClient is used by default.
func WithDoer(doer Doer) *ClientOption {
	return &ClientOption{inner: client.WithDoer(doer)}
}
