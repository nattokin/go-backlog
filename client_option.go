package backlog

import (
	"io"
	"net/http"
	"net/url"
)

// ──────────────────────────────────────────────────────────────
//  Client options
// ──────────────────────────────────────────────────────────────

type clientOption struct {
	set func(config *clientConfig)
}

type clientConfig struct {
	Doer Doer
}

func WithDoer(doer Doer) *clientOption {
	return &clientOption{
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

func withHeader(header http.Header) *httpRequestOption {
	return &httpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Header = header
		},
	}
}

func withBody(body io.Reader) *httpRequestOption {
	return &httpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Body = body
		},
	}
}

func withQuery(query url.Values) *httpRequestOption {
	return &httpRequestOption{
		set: func(config *httpRequestConfig) {
			config.Query = query
		},
	}
}
