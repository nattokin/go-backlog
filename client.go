package backlog

import (
	"github.com/nattokin/go-backlog/internal/core"
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
	User        *UserService
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
func NewClient(baseURL, token string, opts ...*ClientOption) (*Client, error) {
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
