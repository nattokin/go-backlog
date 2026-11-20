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
//
// A Client is safe for concurrent use by multiple goroutines, as long as its
// [Doer] is. The Client and its services hold no state that changes after the
// Client is created.
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

// NewClient creates and initializes a Backlog API Client that authenticates
// with an API key.
// It requires a baseURL (e.g. "https://example.backlog.com") and an API key.
// The API key is sent in the "Backlog-API-Key" request header.
// To authenticate with an OAuth 2.0 access token, use [NewOAuthClient].
//
// The base URL must use the http or https scheme and include a host.
// With an http base URL, the API key is sent in plain text.
//
// It returns an [*InternalClientError] if the base URL or API key is invalid.
//
// Supported options:
//   - [WithDoer]
//   - [WithUserAgent]
func NewClient(baseURL, apiKey string, opts ...*ClientOption) (*Client, error) {
	httpClient, err := client.NewClient(baseURL, apiKey, innerClientOptions(opts)...)
	return newClient(httpClient, err)
}

// NewOAuthClient creates and initializes a Backlog API Client that authenticates
// with an OAuth 2.0 access token.
// It requires a baseURL (e.g. "https://example.backlog.com") and an access token.
// The access token is sent as "Authorization: Bearer <accessToken>".
//
// The Client does not obtain, refresh or track the expiry of access tokens.
// Access tokens expire after one hour; when a token expires, obtain a new one
// and create a new Client with NewOAuthClient.
// When the API rejects the access token, the returned error is an [*OAuthError].
//
// The base URL must use the http or https scheme and include a host.
// With an http base URL, the access token is sent in plain text.
//
// It returns an [*InternalClientError] if the base URL or access token is invalid.
//
// Supported options:
//   - [WithDoer]
//   - [WithUserAgent]
func NewOAuthClient(baseURL, accessToken string, opts ...*ClientOption) (*Client, error) {
	httpClient, err := client.NewOAuthClient(baseURL, accessToken, innerClientOptions(opts)...)
	return newClient(httpClient, err)
}

func innerClientOptions(opts []*ClientOption) []*client.ClientOption {
	innerOpts := make([]*client.ClientOption, len(opts))
	for i, o := range opts {
		innerOpts[i] = o.inner
	}
	return innerOpts
}

func newClient(httpClient *client.Client, err error) (*Client, error) {
	if err != nil {
		return nil, convertError(err)
	}

	c := &Client{
		httpClient: httpClient,
	}

	initServices(c)

	return c, nil
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
// If this option is not provided, [net/http.DefaultClient] is used. It has no
// timeout, so a request can block indefinitely if the server does not respond.
// Set a deadline on the context passed to each method, or pass an
// *http.Client with Timeout set:
//
//	c, err := backlog.NewClient(baseURL, apiKey,
//		backlog.WithDoer(&http.Client{Timeout: 30 * time.Second}),
//	)
//
// To use the Client from multiple goroutines, doer must be safe for
// concurrent use, as *http.Client is.
func WithDoer(doer Doer) *ClientOption {
	return &ClientOption{inner: client.WithDoer(doer)}
}

// WithUserAgent returns a ClientOption that sets the product token of the
// application, such as "myapp/1.2.0". go-backlog appends its own product
// token after it:
//
//	myapp/1.2.0 go-backlog/v0.33.0 (+https://github.com/nattokin/go-backlog)
//
// To replace the User-Agent header entirely, set the header in the [Doer]
// given with [WithDoer].
//
// [NewClient] and [NewOAuthClient] return an [*InternalClientError] if
// userAgent is empty or contains control characters.
func WithUserAgent(userAgent string) *ClientOption {
	return &ClientOption{inner: client.WithUserAgent(userAgent)}
}
