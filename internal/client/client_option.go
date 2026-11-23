package client

import (
	"io"
	"net/http"
	"net/url"
)

// ClientOption configures a Client when it is created.
type ClientOption struct {
	set func(config *clientConfig)
}

type clientConfig struct {
	Doer      Doer
	UserAgent *string
}

// WithDoer sets the Doer that sends requests. DefaultDoer is used when it is
// not set.
func WithDoer(doer Doer) *ClientOption {
	return &ClientOption{
		set: func(config *clientConfig) {
			config.Doer = doer
		},
	}
}

// WithUserAgent sets the product token of the application, which is placed
// before UserAgent in the User-Agent header.
func WithUserAgent(userAgent string) *ClientOption {
	return &ClientOption{
		set: func(config *clientConfig) {
			config.UserAgent = &userAgent
		},
	}
}

// HttpRequestOption configures a request built by Client.NewRequest.
type HttpRequestOption struct {
	set func(config *httpRequestConfig)
}

type httpRequestConfig struct {
	Header http.Header
	Body   io.Reader
	Query  url.Values
}

// WithHeader sets the request headers. The header is cloned, so later
// changes to it do not affect the request.
func WithHeader(header http.Header) *HttpRequestOption {
	return &HttpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Header = header
		},
	}
}

// WithBody sets the request body.
func WithBody(body io.Reader) *HttpRequestOption {
	return &HttpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Body = body
		},
	}
}

// WithQuery sets the query string of the request URL.
func WithQuery(query url.Values) *HttpRequestOption {
	return &HttpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Query = query
		},
	}
}
