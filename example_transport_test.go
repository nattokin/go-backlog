package backlog_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

// TestMain replaces http.DefaultTransport, which the default client sends
// requests with, with exampleTransport, so that Example functions can create a
// client the same way users do (without WithDoer) and still verify their
// output. No request reaches example.backlog.com.
//
// Every test in this package's test binary shares http.DefaultTransport. Tests
// that send requests pass a mock with WithDoer and do not use it.
func TestMain(m *testing.M) {
	http.DefaultTransport = exampleTransport{}
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

	key := req.Method + " " + strings.TrimPrefix(req.URL.Path, "/api/v2")
	respond, ok := exampleRoutes[key]
	if !ok {
		return nil, fmt.Errorf("no example route for %s", key)
	}
	return respond(req)
}

// exampleRoutes maps "METHOD /path" to the response for Example functions.
// The path is relative to /api/v2, and the query string is not part of the key.
//
// IDs in the Example functions only select a route. In successful cases, pass
// the ID of the fixture that the route returns, so that the output matches
// the argument. Error cases use IDs that no fixture has.
var exampleRoutes = map[string]respondFunc{
	// Issue
	"GET /issues":          jsonBody(fixture.Issue.ListJSON),
	"POST /issues":         jsonBody(fixture.Issue.SingleJSON),
	"GET /issues/count":    jsonBody(`{"count":2}`),
	"GET /issues/PRJ-1":    jsonBody(fixture.Issue.SingleJSON),
	"PATCH /issues/PRJ-1":  jsonBody(fixture.Issue.SingleJSON),
	"DELETE /issues/PRJ-1": jsonBody(fixture.Issue.SingleJSON),

	"GET /issues/PRJ-1/comments":                  jsonBody(fixture.Comment.ListJSON),
	"POST /issues/PRJ-1/comments":                 jsonBody(fixture.Comment.SingleJSON),
	"GET /issues/PRJ-1/comments/count":            jsonBody(`{"count":2}`),
	"GET /issues/PRJ-1/comments/1":                jsonBody(fixture.Comment.SingleJSON),
	"PATCH /issues/PRJ-1/comments/1":              jsonBody(fixture.Comment.SingleJSON),
	"DELETE /issues/PRJ-1/comments/1":             jsonBody(fixture.Comment.SingleJSON),
	"GET /issues/PRJ-1/comments/1/notifications":  jsonBody(commentNotificationsJSON),
	"POST /issues/PRJ-1/comments/1/notifications": jsonBody(fixture.Comment.SingleJSON),

	"GET /issues/PRJ-1/participants": jsonBody(fixture.User.ListJSON),

	"GET /issues/TEST-1/attachments":      jsonBody(fixture.Attachment.ListJSON),
	"GET /issues/TEST-1/attachments/2":    binary("A.png"),
	"DELETE /issues/TEST-1/attachments/8": jsonBody(fixture.Attachment.SingleJSON),

	"GET /issues/TEST-1/sharedFiles":           jsonBody(fixture.SharedFile.ListJSON),
	"POST /issues/TEST-1/sharedFiles":          jsonBody(fixture.SharedFile.ListJSON),
	"DELETE /issues/TEST-1/sharedFiles/454403": jsonBody(fixture.SharedFile.SingleJSON),

	// Priority
	"GET /priorities": jsonBody(fixture.Priority.ListJSON),

	// Project
	"GET /projects":         jsonBody(fixture.Project.ListJSON),
	"POST /projects":        jsonBody(fixture.Project.SingleJSON),
	"GET /projects/TEST":    jsonBody(fixture.Project.SingleJSON),
	"PATCH /projects/TEST":  jsonBody(fixture.Project.SingleJSON),
	"DELETE /projects/TEST": jsonBody(fixture.Project.SingleJSON),

	"GET /projects/TEST/activities": jsonBody(fixture.Activity.ListJSON),

	"GET /projects/TEST/administrators":    jsonBody(fixture.User.ListJSON),
	"POST /projects/TEST/administrators":   jsonBody(fixture.User.SingleJSON),
	"DELETE /projects/TEST/administrators": jsonBody(fixture.User.SingleJSON),

	"GET /projects/TEST/categories":       jsonBody(fixture.Category.ListJSON),
	"POST /projects/TEST/categories":      jsonBody(fixture.Category.SingleJSON),
	"PATCH /projects/TEST/categories/12":  jsonBody(fixture.Category.SingleJSON),
	"DELETE /projects/TEST/categories/12": jsonBody(fixture.Category.SingleJSON),

	"GET /projects/TEST/customFields":               jsonBody(fixture.CustomField.ListJSON),
	"POST /projects/TEST/customFields":              jsonBody(fixture.CustomField.SingleJSON),
	"PATCH /projects/TEST/customFields/1":           jsonBody(fixture.CustomField.SingleJSON),
	"DELETE /projects/TEST/customFields/1":          jsonBody(fixture.CustomField.SingleJSON),
	"POST /projects/TEST/customFields/1/items":      jsonBody(fixture.CustomField.SingleJSON),
	"PATCH /projects/TEST/customFields/1/items/10":  jsonBody(fixture.CustomField.SingleJSON),
	"DELETE /projects/TEST/customFields/1/items/10": jsonBody(fixture.CustomField.SingleJSON),

	"GET /projects/TEST/diskUsage": jsonBody(fixture.Project.DiskUsageJSON),

	"GET /projects/TEST/files/1":                binary("shared.png"),
	"GET /projects/TEST/files/metadata/design/": jsonBody(fixture.SharedFile.ListJSON),

	"GET /projects/TEST/image": binary("test.png"),

	"GET /projects/TEST/issueTypes":      jsonBody(fixture.IssueType.ListJSON),
	"POST /projects/TEST/issueTypes":     jsonBody(fixture.IssueType.SingleJSON),
	"PATCH /projects/TEST/issueTypes/1":  jsonBody(fixture.IssueType.SingleJSON),
	"DELETE /projects/TEST/issueTypes/1": jsonBody(fixture.IssueType.SingleJSON),

	"GET /projects/TEST/statuses":                      jsonBody(fixture.Status.ListJSON),
	"POST /projects/TEST/statuses":                     jsonBody(fixture.Status.SingleJSON),
	"PATCH /projects/TEST/statuses/1":                  jsonBody(fixture.Status.SingleJSON),
	"DELETE /projects/TEST/statuses/1":                 jsonBody(fixture.Status.SingleJSON),
	"PATCH /projects/TEST/statuses/updateDisplayOrder": jsonBody(fixture.Status.ListJSON),

	"GET /projects/TEST/users":    jsonBody(fixture.User.ListJSON),
	"POST /projects/TEST/users":   jsonBody(fixture.User.SingleJSON),
	"DELETE /projects/TEST/users": jsonBody(fixture.User.SingleJSON),

	"GET /projects/TEST/versions":      jsonBody(fixture.Version.ListJSON),
	"POST /projects/TEST/versions":     jsonBody(fixture.Version.SingleJSON),
	"PATCH /projects/TEST/versions/1":  jsonBody(fixture.Version.SingleJSON),
	"DELETE /projects/TEST/versions/1": jsonBody(fixture.Version.SingleJSON),

	"GET /projects/TEST/webhooks":      jsonBody(fixture.Webhook.ListJSON),
	"POST /projects/TEST/webhooks":     jsonBody(fixture.Webhook.AllEventJSON),
	"GET /projects/TEST/webhooks/1":    jsonBody(fixture.Webhook.AllEventJSON),
	"PATCH /projects/TEST/webhooks/1":  jsonBody(fixture.Webhook.AllEventJSON),
	"DELETE /projects/TEST/webhooks/1": jsonBody(fixture.Webhook.AllEventJSON),

	// PullRequest
	"GET /projects/TEST/git/repositories/myrepo/pullRequests":       jsonBody(fixture.PullRequest.ListJSON),
	"POST /projects/TEST/git/repositories/myrepo/pullRequests":      jsonBody(fixture.PullRequest.SingleJSON),
	"GET /projects/TEST/git/repositories/myrepo/pullRequests/count": jsonBody(`{"count":2}`),
	"GET /projects/TEST/git/repositories/myrepo/pullRequests/1":     jsonBody(fixture.PullRequest.SingleJSON),
	"PATCH /projects/TEST/git/repositories/myrepo/pullRequests/1":   jsonBody(fixture.PullRequest.SingleJSON),

	"GET /projects/TEST/git/repositories/myrepo/pullRequests/1/attachments":      jsonBody(fixture.Attachment.ListJSON),
	"GET /projects/TEST/git/repositories/myrepo/pullRequests/1/attachments/2":    binary("A.png"),
	"DELETE /projects/TEST/git/repositories/myrepo/pullRequests/1/attachments/8": jsonBody(fixture.Attachment.SingleJSON),

	"GET /projects/TEST/git/repositories/myrepo/pullRequests/1/comments":       jsonBody(fixture.Comment.ListJSON),
	"POST /projects/TEST/git/repositories/myrepo/pullRequests/1/comments":      jsonBody(fixture.Comment.SingleJSON),
	"GET /projects/TEST/git/repositories/myrepo/pullRequests/1/comments/count": jsonBody(`{"count":2}`),
	"PATCH /projects/TEST/git/repositories/myrepo/pullRequests/1/comments/1":   jsonBody(fixture.Comment.SingleJSON),

	// RecentlyViewed
	"GET /users/myself/recentlyViewedIssues":   jsonBody(fixture.RecentlyViewed.IssueListJSON),
	"POST /users/myself/recentlyViewedIssues":  jsonBody(fixture.RecentlyViewed.IssueSingleJSON),
	"GET /users/myself/recentlyViewedProjects": jsonBody(fixture.RecentlyViewed.ProjectListJSON),
	"GET /users/myself/recentlyViewedWikis":    jsonBody(fixture.RecentlyViewed.WikiListJSON),
	"POST /users/myself/recentlyViewedWikis":   jsonBody(fixture.RecentlyViewed.WikiSingleJSON),

	// Repository
	"GET /projects/TEST/git/repositories":     jsonBody(fixture.Repository.ListJSON),
	"GET /projects/TEST/git/repositories/foo": jsonBody(fixture.Repository.SingleJSON),

	// Resolution
	"GET /resolutions": jsonBody(fixture.Resolution.ListJSON),

	// Space
	"GET /activities/3153":    jsonBody(fixture.Activity.SingleJSON),
	"GET /space":              jsonBody(fixture.Space.SpaceJSON),
	"GET /space/activities":   jsonBody(fixture.Activity.ListJSON),
	"POST /space/attachment":  jsonBody(fixture.Attachment.UploadJSON),
	"GET /space/diskUsage":    jsonBody(fixture.Space.DiskUsageJSON),
	"GET /space/image":        binary("space.png"),
	"GET /space/notification": jsonBody(fixture.Space.NotificationJSON),
	"PUT /space/notification": jsonBody(fixture.Space.NotificationJSON),

	// Star
	"POST /stars":      noContent(),
	"DELETE /stars/42": noContent(),

	// User
	"GET /users":      jsonBody(fixture.User.ListJSON),
	"POST /users":     jsonBody(fixture.User.SingleJSON),
	"GET /users/1":    jsonBody(fixture.User.SingleJSON),
	"PATCH /users/1":  jsonBody(fixture.User.SingleJSON),
	"DELETE /users/1": jsonBody(fixture.User.SingleJSON),

	"GET /users/1/activities": jsonBody(fixture.Activity.ListJSON),

	"GET /users/1/icon": binary("icon.png"),

	"GET /users/1/stars":       jsonBody(fixture.Star.ListJSON),
	"GET /users/1/stars/count": jsonBody(fixture.Star.CountJSON),

	"GET /users/myself": jsonBody(fixture.User.SingleJSON),

	// Watching
	"GET /users/1/watchings":       jsonBody(fixture.Watching.ListJSON),
	"GET /users/1/watchings/count": jsonBody(fixture.Watching.CountJSON),
	"POST /watchings":              jsonBody(fixture.Watching.SingleJSON),
	"GET /watchings/1":             jsonBody(fixture.Watching.SingleJSON),
	"PATCH /watchings/1":           jsonBody(fixture.Watching.SingleJSON),
	"DELETE /watchings/1":          jsonBody(fixture.Watching.SingleJSON),
	"POST /watchings/1/markAsRead": noContent(),

	// Wiki
	"GET /wikis":       jsonBody(fixture.Wiki.ListJSON),
	"GET /wikis/tags":  jsonBody(`[{"id": 1, "name": "test"}]`),
	"POST /wikis":      jsonBody(fixture.Wiki.MinimumJSON),
	"GET /wikis/count": jsonBody(`{"count": 5}`),
	"GET /wikis/34":    jsonBody(fixture.Wiki.MinimumJSON),
	"PATCH /wikis/34":  jsonBody(fixture.Wiki.MinimumJSON),
	"DELETE /wikis/34": jsonBody(fixture.Wiki.MinimumJSON),

	"GET /wikis/34/attachments":      jsonBody(fixture.Attachment.ListJSON),
	"POST /wikis/34/attachments":     jsonBody(fixture.Attachment.ListJSON),
	"GET /wikis/34/attachments/2":    binary("A.png"),
	"DELETE /wikis/34/attachments/8": jsonBody(fixture.Attachment.SingleJSON),

	"GET /wikis/34/history": jsonBody(fixture.WikiHistory.ListJSON),

	"GET /wikis/34/sharedFiles":           jsonBody(fixture.SharedFile.ListJSON),
	"POST /wikis/34/sharedFiles":          jsonBody(fixture.SharedFile.ListJSON),
	"DELETE /wikis/34/sharedFiles/454403": jsonBody(fixture.SharedFile.SingleJSON),

	"GET /wikis/34/stars": jsonBody(fixture.Star.ListJSON),

	// Error responses
	"GET /wikis/1": mock.NewNotFoundDoFunc(),
	"GET /wikis/2": rateLimitedDoFunc,
}

// commentNotificationsJSON is the response body of the comment notification list.
const commentNotificationsJSON = `[{"id":25,"alreadyRead":false,"reason":2,"resourceAlreadyRead":false}]`

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

// respondFunc returns the response to a request sent by an Example function.
type respondFunc = func(*http.Request) (*http.Response, error)

// jsonBody responds with HTTP 200 and the given JSON body.
func jsonBody(body string) respondFunc {
	return mock.NewDoFunc(body)
}

// binary responds with HTTP 200 and a PNG file download with the given filename.
func binary(filename string) respondFunc {
	return mock.NewBinaryDoFunc(filename, "image/png", []byte("PNG"))
}

// noContent responds with HTTP 204 No Content.
func noContent() respondFunc {
	return mock.NewNoContentDoFunc()
}
