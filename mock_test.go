package backlog_test

import (
	"bytes"
	"io"
	"net/http"
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
