package client

import (
	"io"
	"net/http"
	"net/url"
)

type ClientOption struct {
	set func(config *clientConfig)
}

type clientConfig struct {
	Doer      Doer
	UserAgent *string
}

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

type HttpRequestOption struct {
	set func(config *httpRequestConfig)
}

type httpRequestConfig struct {
	Header http.Header
	Body   io.Reader
	Query  url.Values
}

func WithHeader(header http.Header) *HttpRequestOption {
	return &HttpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Header = header
		},
	}
}

func WithBody(body io.Reader) *HttpRequestOption {
	return &HttpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Body = body
		},
	}
}

func WithQuery(query url.Values) *HttpRequestOption {
	return &HttpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Query = query
		},
	}
}
