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
	// Star
	"POST /api/v2/stars":      mock.NewNoContentDoFunc(),
	"DELETE /api/v2/stars/42": mock.NewNoContentDoFunc(),

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
