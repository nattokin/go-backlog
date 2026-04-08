package backlog

import (
	"io"
	"net/http"
	"net/url"
)

// ──────────────────────────────────────────────────────────────
//  Http request otions
// ──────────────────────────────────────────────────────────────

type httpRequestOption struct {
	Set func(config *httpRequestConfig)
}

type httpRequestConfig struct {
	Header http.Header
	Body   io.Reader
	Query  url.Values
}

func withHeader(header http.Header) *httpRequestOption {
	return &httpRequestOption{
		Set: func(config *httpRequestConfig) {
			config.Header = header
		},
	}
}

func withBody(body io.Reader) *httpRequestOption {
	return &httpRequestOption{
		Set: func(config *httpRequestConfig) {
			config.Body = body
		},
	}
}

func withQuery(query url.Values) *httpRequestOption {
	return &httpRequestOption{
		Set: func(config *httpRequestConfig) {
			config.Query = query
		},
	}
}
