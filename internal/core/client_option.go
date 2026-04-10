package core

import (
	"io"
	"net/http"
	"net/url"
)

// ──────────────────────────────────────────────────────────────
//  Client options
// ──────────────────────────────────────────────────────────────

// ClientOption defines a functional option for configuring a Client.
// It is used to change the default behavior of the Client.
type ClientOption struct {
	set func(config *clientConfig)
}

// clientConfig holds the internal configuration settings for the Client.
type clientConfig struct {
	// Doer is the HTTP client used to make requests.
	Doer Doer
}

// WithDoer returns a ClientOption that sets the HTTP client (Doer) for the Client.
// This is useful for providing a custom *http.Client or a mock implementation during testing.
//
// If this option is not provided, http.DefaultClient is used by default.
func WithDoer(doer Doer) *ClientOption {
	return &ClientOption{
		set: func(config *clientConfig) {
			config.Doer = doer
		},
	}
}

// ──────────────────────────────────────────────────────────────
//  Http request otions
// ──────────────────────────────────────────────────────────────

type httpRequestOption struct {
	set func(config *httpRequestConfig)
}

type httpRequestConfig struct {
	Header http.Header
	Body   io.Reader
	Query  url.Values
}

func WithHeader(header http.Header) *httpRequestOption {
	return &httpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Header = header
		},
	}
}

func WithBody(body io.Reader) *httpRequestOption {
	return &httpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Body = body
		},
	}
}

func WithQuery(query url.Values) *httpRequestOption {
	return &httpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Query = query
		},
	}
}
