package backlog_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

// TestMain replaces the transport of http.DefaultClient with exampleTransport,
// so that Example functions can create a client the same way users do (without
// WithDoer) and still verify their output. No request reaches
// example.backlog.com.
//
// Every test in this package's test binary shares http.DefaultClient. Tests
// that send requests pass a mock with WithDoer and do not use it.
func TestMain(m *testing.M) {
	http.DefaultClient.Transport = exampleTransport{}
	os.Exit(m.Run())
}

// exampleTransport is an http.RoundTripper that returns the response
// registered in exampleRoutes for the request method and path. A request to
// an unregistered route fails, and never falls back to the network.
type exampleTransport struct{}

func (exampleTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		defer req.Body.Close()
	}

	key := req.Method + " " + req.URL.Path
	respond, ok := exampleRoutes[key]
	if !ok {
		return nil, fmt.Errorf("no example route for %s", key)
	}
	return respond(req)
}

// exampleRoutes maps "METHOD /api/v2/path" to the response for Example
// functions. The query string is not part of the key.
//
// IDs in the Example functions only select a route. In successful cases, pass
// the ID of the fixture that the route returns, so that the output matches
// the argument. Error cases use IDs that no fixture has.
var exampleRoutes = map[string]func(*http.Request) (*http.Response, error){
	// Issue
	"GET /api/v2/issues":                                 mock.NewDoFunc(fixture.Issue.ListJSON),
	"POST /api/v2/issues":                                mock.NewDoFunc(fixture.Issue.SingleJSON),
	"GET /api/v2/issues/count":                           mock.NewDoFunc(`{"count":2}`),
	"GET /api/v2/issues/PRJ-1":                           mock.NewDoFunc(fixture.Issue.SingleJSON),
	"PATCH /api/v2/issues/PRJ-1":                         mock.NewDoFunc(fixture.Issue.SingleJSON),
	"DELETE /api/v2/issues/PRJ-1":                        mock.NewDoFunc(fixture.Issue.SingleJSON),
	"GET /api/v2/issues/PRJ-1/comments":                  mock.NewDoFunc(fixture.Comment.ListJSON),
	"POST /api/v2/issues/PRJ-1/comments":                 mock.NewDoFunc(fixture.Comment.SingleJSON),
	"GET /api/v2/issues/PRJ-1/comments/count":            mock.NewDoFunc(`{"count":2}`),
	"GET /api/v2/issues/PRJ-1/comments/1":                mock.NewDoFunc(fixture.Comment.SingleJSON),
	"PATCH /api/v2/issues/PRJ-1/comments/1":              mock.NewDoFunc(fixture.Comment.SingleJSON),
	"DELETE /api/v2/issues/PRJ-1/comments/1":             mock.NewDoFunc(fixture.Comment.SingleJSON),
	"GET /api/v2/issues/PRJ-1/comments/1/notifications":  mock.NewDoFunc(`[{"id":25,"alreadyRead":false,"reason":2,"resourceAlreadyRead":false}]`),
	"POST /api/v2/issues/PRJ-1/comments/1/notifications": mock.NewDoFunc(fixture.Comment.SingleJSON),
	"GET /api/v2/issues/PRJ-1/participants":              mock.NewDoFunc(fixture.User.ListJSON),
	"GET /api/v2/issues/TEST-1/attachments":              mock.NewDoFunc(fixture.Attachment.ListJSON),
	"GET /api/v2/issues/TEST-1/attachments/2":            mock.NewBinaryDoFunc("A.png", "image/png", []byte("PNG")),
	"DELETE /api/v2/issues/TEST-1/attachments/8":         mock.NewDoFunc(fixture.Attachment.SingleJSON),
	"GET /api/v2/issues/TEST-1/sharedFiles":              mock.NewDoFunc(fixture.SharedFile.ListJSON),
	"POST /api/v2/issues/TEST-1/sharedFiles":             mock.NewDoFunc(fixture.SharedFile.ListJSON),
	"DELETE /api/v2/issues/TEST-1/sharedFiles/454403":    mock.NewDoFunc(fixture.SharedFile.SingleJSON),

	// PullRequest
	"GET /api/v2/projects/TEST/git/repositories/myrepo/pullRequests":                    mock.NewDoFunc(fixture.PullRequest.ListJSON),
	"POST /api/v2/projects/TEST/git/repositories/myrepo/pullRequests":                   mock.NewDoFunc(fixture.PullRequest.SingleJSON),
	"GET /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/count":              mock.NewDoFunc(`{"count":2}`),
	"GET /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1":                  mock.NewDoFunc(fixture.PullRequest.SingleJSON),
	"PATCH /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1":                mock.NewDoFunc(fixture.PullRequest.SingleJSON),
	"GET /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1/attachments":      mock.NewDoFunc(fixture.Attachment.ListJSON),
	"GET /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1/attachments/2":    mock.NewBinaryDoFunc("A.png", "image/png", []byte("PNG")),
	"DELETE /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1/attachments/8": mock.NewDoFunc(fixture.Attachment.SingleJSON),
	"GET /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1/comments":         mock.NewDoFunc(fixture.Comment.ListJSON),
	"POST /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1/comments":        mock.NewDoFunc(fixture.Comment.SingleJSON),
	"GET /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1/comments/count":   mock.NewDoFunc(`{"count":2}`),
	"PATCH /api/v2/projects/TEST/git/repositories/myrepo/pullRequests/1/comments/1":     mock.NewDoFunc(fixture.Comment.SingleJSON),

	// RecentlyViewed
	"GET /api/v2/users/myself/recentlyViewedIssues":   mock.NewDoFunc(fixture.RecentlyViewed.IssueListJSON),
	"POST /api/v2/users/myself/recentlyViewedIssues":  mock.NewDoFunc(fixture.RecentlyViewed.IssueSingleJSON),
	"GET /api/v2/users/myself/recentlyViewedProjects": mock.NewDoFunc(fixture.RecentlyViewed.ProjectListJSON),
	"GET /api/v2/users/myself/recentlyViewedWikis":    mock.NewDoFunc(fixture.RecentlyViewed.WikiListJSON),
	"POST /api/v2/users/myself/recentlyViewedWikis":   mock.NewDoFunc(fixture.RecentlyViewed.WikiSingleJSON),

	// Repository
	"GET /api/v2/projects/TEST/git/repositories":     mock.NewDoFunc(fixture.Repository.ListJSON),
	"GET /api/v2/projects/TEST/git/repositories/foo": mock.NewDoFunc(fixture.Repository.SingleJSON),

	// Space
	"GET /api/v2/activities/3153":    mock.NewDoFunc(fixture.Activity.SingleJSON),
	"GET /api/v2/space":              mock.NewDoFunc(fixture.Space.SpaceJSON),
	"GET /api/v2/space/activities":   mock.NewDoFunc(fixture.Activity.ListJSON),
	"POST /api/v2/space/attachment":  mock.NewDoFunc(fixture.Attachment.UploadJSON),
	"GET /api/v2/space/diskUsage":    mock.NewDoFunc(fixture.Space.DiskUsageJSON),
	"GET /api/v2/space/notification": mock.NewDoFunc(fixture.Space.NotificationJSON),
	"PUT /api/v2/space/notification": mock.NewDoFunc(fixture.Space.NotificationJSON),

	// Star
	"POST /api/v2/stars":      mock.NewNoContentDoFunc(),
	"DELETE /api/v2/stars/42": mock.NewNoContentDoFunc(),

	// User
	"GET /api/v2/users":               mock.NewDoFunc(fixture.User.ListJSON),
	"POST /api/v2/users":              mock.NewDoFunc(fixture.User.SingleJSON),
	"GET /api/v2/users/1":             mock.NewDoFunc(fixture.User.SingleJSON),
	"PATCH /api/v2/users/1":           mock.NewDoFunc(fixture.User.SingleJSON),
	"DELETE /api/v2/users/1":          mock.NewDoFunc(fixture.User.SingleJSON),
	"GET /api/v2/users/1/activities":  mock.NewDoFunc(fixture.Activity.ListJSON),
	"GET /api/v2/users/1/icon":        mock.NewBinaryDoFunc("icon.png", "image/png", []byte("PNG")),
	"GET /api/v2/users/1/stars":       mock.NewDoFunc(fixture.Star.ListJSON),
	"GET /api/v2/users/1/stars/count": mock.NewDoFunc(fixture.Star.CountJSON),
	"GET /api/v2/users/myself":        mock.NewDoFunc(fixture.User.SingleJSON),

	// Wiki
	"GET /api/v2/wikis":                          mock.NewDoFunc(fixture.Wiki.ListJSON),
	"POST /api/v2/wikis":                         mock.NewDoFunc(fixture.Wiki.MinimumJSON),
	"GET /api/v2/wikis/count":                    mock.NewDoFunc(`{"count": 5}`),
	"GET /api/v2/wikis/34":                       mock.NewDoFunc(fixture.Wiki.MinimumJSON),
	"PATCH /api/v2/wikis/34":                     mock.NewDoFunc(fixture.Wiki.MinimumJSON),
	"DELETE /api/v2/wikis/34":                    mock.NewDoFunc(fixture.Wiki.MinimumJSON),
	"GET /api/v2/wikis/34/attachments":           mock.NewDoFunc(fixture.Attachment.ListJSON),
	"POST /api/v2/wikis/34/attachments":          mock.NewDoFunc(fixture.Attachment.ListJSON),
	"GET /api/v2/wikis/34/attachments/2":         mock.NewBinaryDoFunc("A.png", "image/png", []byte("PNG")),
	"DELETE /api/v2/wikis/34/attachments/8":      mock.NewDoFunc(fixture.Attachment.SingleJSON),
	"GET /api/v2/wikis/34/history":               mock.NewDoFunc(fixture.WikiHistory.ListJSON),
	"GET /api/v2/wikis/34/sharedFiles":           mock.NewDoFunc(fixture.SharedFile.ListJSON),
	"POST /api/v2/wikis/34/sharedFiles":          mock.NewDoFunc(fixture.SharedFile.ListJSON),
	"DELETE /api/v2/wikis/34/sharedFiles/454403": mock.NewDoFunc(fixture.SharedFile.SingleJSON),
	"GET /api/v2/wikis/34/stars":                 mock.NewDoFunc(fixture.Star.ListJSON),

	// Error responses
	"GET /api/v2/wikis/1": mock.NewNotFoundDoFunc(),
	"GET /api/v2/wikis/2": rateLimitedDoFunc,
}

// rateLimitedDoFunc responds with HTTP 429 Too Many Requests and the
// X-RateLimit-* headers.
func rateLimitedDoFunc(_ *http.Request) (*http.Response, error) {
	resp := mock.NewErrorResponse(http.StatusTooManyRequests, `{"errors":[]}`)
	resp.Header = http.Header{}
	resp.Header.Set("X-RateLimit-Limit", "600")
	resp.Header.Set("X-RateLimit-Remaining", "0")
	resp.Header.Set("X-RateLimit-Reset", "1735689600")
	return resp, nil
}
