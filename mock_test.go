package backlog_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type mockDoer struct {
	do func(req *http.Request) (*http.Response, error)
}

func (d *mockDoer) Do(req *http.Request) (*http.Response, error) {
	return d.do(req)
}

// newMockDoer returns a mockDoer that always responds with HTTP 200 and the given body.
func newMockDoer(body string) *mockDoer {
	return &mockDoer{
		do: func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}, nil
		},
	}
}

// doerNoContent is a mockDoer that always responds with HTTP 204 No Content.
var doerNoContent = &mockDoer{
	do: func(_ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
	},
}

// newAuthErrorDoFunc returns a doFunc that always responds with HTTP 401 Unauthorized
// and a Backlog authentication failure error body.
// It returns a new response on each call to avoid reuse of the consumed Body reader.
func newAuthErrorDoFunc() func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
		}, nil
	}
}
