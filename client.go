package backlog

import (
	"net/http"

	"github.com/nattokin/go-backlog/internal/core"
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
	core *core.Client

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
	coreOpts := make([]*core.ClientOption, len(opts))
	for i, o := range opts {
		coreOpts[i] = o.core
	}
	c, err := core.NewClient(baseURL, token, coreOpts...)
	if err != nil {
		return nil, convertError(err)
	}

	client := &Client{
		core: c,
	}

	initServices(client)

	return client, nil
}

// ──────────────────────────────────────────────────────────────
//  Service initialization
// ──────────────────────────────────────────────────────────────

func initServices(c *Client) {
	baseOptionService := &core.OptionService{}

	c.Issue = newIssueService(c.core.Method, baseOptionService)

	c.Project = newProjectService(c.core.Method, baseOptionService)

	c.PullRequest = newPullRequestService(c.core.Method, baseOptionService)

	c.RecentlyViewed = newRecentlyViewedService(c.core.Method, baseOptionService)

	c.Repository = newRepositoryService(c.core.Method)

	c.Space = newSpaceService(c.core.Method, baseOptionService)

	c.Star = newStarService(c.core.Method, baseOptionService)

	c.User = newUserService(c.core.Method, baseOptionService)

	c.Wiki = newWikiService(c.core.Method, baseOptionService)
}

// ──────────────────────────────────────────────────────────────
//  Client options
// ──────────────────────────────────────────────────────────────

// ClientOption defines a functional option for configuring a Client.
// It is used to change the default behavior of the Client.
type ClientOption struct {
	core *core.ClientOption
}

// WithDoer returns a ClientOption that sets the HTTP client (Doer) for the Client.
// This is useful for providing a custom *http.Client or a mock implementation during testing.
//
// If this option is not provided, http.DefaultClient is used by default.
func WithDoer(doer Doer) *ClientOption {
	return &ClientOption{core: core.WithDoer(doer)}
}
