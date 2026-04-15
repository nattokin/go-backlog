package backlog

import (
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/user"
)

// ──────────────────────────────────────────────────────────────
//  Doer interface (HTTP abstraction)
// ──────────────────────────────────────────────────────────────

type Doer = core.Doer

// ──────────────────────────────────────────────────────────────
//  Client structure and initialization
// ──────────────────────────────────────────────────────────────

// Client represents a Backlog API client.
// It wraps an underlying HTTP Doer and provides typed services for API access.
type Client struct {
	core *core.Client

	// Service endpoints
	Issue       *IssueService
	Project     *ProjectService
	PullRequest *PullRequestService
	Space       *SpaceService
	User        *user.UserService
	Wiki        *WikiService
}

// ──────────────────────────────────────────────────────────────
//  Client constructor
// ──────────────────────────────────────────────────────────────

// NewClient creates and initializes a Backlog API Client.
// It requires a baseURL and an API token.
//
// This function supports options returned by package-level functions,
// such as:
//   - WithDoer
func NewClient(baseURL, token string, opts ...*core.ClientOption) (*Client, error) {
	core, err := core.NewClient(baseURL, token, opts...)
	if err != nil {
		return nil, err
	}

	c := &Client{
		core: core,
	}

	initServices(c)

	return c, nil
}

// ──────────────────────────────────────────────────────────────
//  Service initialization
// ──────────────────────────────────────────────────────────────

func initServices(c *Client) {
	baseOptionService := &core.OptionService{}

	c.Issue = newIssueService(c.core.Method, baseOptionService)

	c.Project = newProjectService(c.core.Method, baseOptionService)

	c.PullRequest = newPullRequestService(c.core.Method)

	c.Space = newSpaceService(c.core.Method, baseOptionService)

	c.User = user.NewUserService(c.core.Method, baseOptionService)

	c.Wiki = newWikiService(c.core.Method, baseOptionService)
}

// ──────────────────────────────────────────────────────────────
//  Client options
// ──────────────────────────────────────────────────────────────

type ClientOption = core.ClientOption

// WithDoer returns a ClientOption that sets the HTTP client (Doer) for the Client.
// This is useful for providing a custom *http.Client or a mock implementation during testing.
//
// If this option is not provided, http.DefaultClient is used by default.
func WithDoer(doer Doer) *core.ClientOption {
	return core.WithDoer(doer)
}
