package backlog_test

import (
	"net/http"

	"github.com/nattokin/go-backlog"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewHTTPClientMock(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewClientMock(baseURL, token string, fn RoundTripFunc) *backlog.Client {
	c, _ := backlog.NewClient(baseURL, token)
	httpClient := NewHTTPClientMock(fn)
	c.ExportSetHTTPClient(httpClient)

	return c
}
