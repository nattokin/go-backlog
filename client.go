package backlog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// ──────────────────────────────────────────────────────────────
//  Doer interface (HTTP abstraction)
// ──────────────────────────────────────────────────────────────

// Doer defines the minimal interface required to perform HTTP requests.
// It is compatible with *http.Client and allows injection of mock clients
// for unit or integration testing.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ──────────────────────────────────────────────────────────────
//  Client structure and initialization
// ──────────────────────────────────────────────────────────────

// Client represents a Backlog API client.
// It wraps an underlying HTTP Doer and provides typed services for API access.
type Client struct {
	baseURL *url.URL
	token   string
	doer    Doer
	wrapper wrapper
	method  *method

	// Service endpoints
	Issue       *IssueService
	Project     *ProjectService
	PullRequest *PullRequestService
	Space       *SpaceService
	User        *UserService
	Wiki        *WikiService
}

// ──────────────────────────────────────────────────────────────
//  HTTP method function set
// ──────────────────────────────────────────────────────────────

// method holds injected HTTP operation functions.
// Each function delegates to Client.do() but can be replaced in tests
// for fine-grained unit testing of individual services.
type method struct {
	Get    func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	Post   func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Patch  func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Delete func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Upload func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error)
}

// ──────────────────────────────────────────────────────────────
//  Client constructor
// ──────────────────────────────────────────────────────────────

// NewClient creates and initializes a Backlog API Client.
// A custom Doer (e.g., *http.Client or mock) may be provided for testing.
// If doer is nil, http.DefaultClient is used.
func NewClient(baseURL, token string, doer Doer) (*Client, error) {
	if token == "" {
		return nil, newInternalClientError("missing token")
	}
	if baseURL == "" {
		return nil, newInternalClientError("missing baseURL")
	}

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	if doer == nil {
		doer = http.DefaultClient
	}

	c := &Client{
		baseURL: u,
		doer:    doer,
		token:   token,
		wrapper: &defaultWrapper{},
	}

	// --- Inject HTTP method wrappers ------------------------------------------
	c.method = &method{
		Get: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
			return c.do(ctx, http.MethodGet, spath, withQuery(query))
		},
		Post: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
			header := http.Header{}
			header.Set("Content-Type", "application/x-www-form-urlencoded")
			if form == nil {
				form = url.Values{}
			}
			return c.do(ctx, http.MethodPost, spath, withHeader(header), withBody(strings.NewReader(form.Encode())))
		},
		Patch: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
			header := http.Header{}
			header.Set("Content-Type", "application/x-www-form-urlencoded")
			if form == nil {
				form = url.Values{}
			}
			return c.do(ctx, http.MethodPatch, spath, withHeader(header), withBody(strings.NewReader(form.Encode())))
		},
		Delete: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
			header := http.Header{}
			header.Set("Content-Type", "application/x-www-form-urlencoded")
			if form == nil {
				form = url.Values{}
			}
			return c.do(ctx, http.MethodDelete, spath, withHeader(header), withBody(strings.NewReader(form.Encode())))
		},
		Upload: func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error) {
			if fileName == "" {
				return nil, newInternalClientError("fileName is required")
			}
			var buf bytes.Buffer
			mw := c.wrapper.NewMultipartWriter(&buf)

			fw, err := mw.CreateFormFile("file", fileName)
			if err != nil {
				return nil, err
			}
			if err = c.wrapper.Copy(fw, r); err != nil {
				return nil, err
			}
			if err := mw.Close(); err != nil {
				return nil, err
			}

			header := http.Header{}
			header.Set("Content-Type", mw.FormDataContentType())

			return c.do(ctx, http.MethodPost, spath, withHeader(header), withBody(&buf))
		},
	}

	initServices(c)

	return c, nil
}

// ──────────────────────────────────────────────────────────────
//  HTTP request creation and execution
// ──────────────────────────────────────────────────────────────

// newRequest builds a new HTTP request with token-based authentication.
func (c *Client) newRequest(ctx context.Context, method, spath string, opts ...*httpRequestOption) (*http.Request, error) {
	if spath == "" {
		return nil, errors.New("spath must not be empty")
	}

	config := &httpRequestConfig{}
	for _, option := range opts {
		if option != nil {
			option.set(config)
		}
	}

	u := *c.baseURL
	u.Path = path.Join(u.Path, "api", apiVersion, spath)
	if config.Query != nil {
		u.RawQuery = config.Query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), config.Body)
	if err != nil {
		return nil, err
	}

	if config.Header != nil {
		req.Header = config.Header.Clone()
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	return req, nil
}

// do executes the given HTTP request using the injected Doer.
// All HTTP calls pass through this function, ensuring consistent error handling.
func (c *Client) do(ctx context.Context, method, spath string, opts ...*httpRequestOption) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, spath, opts...)
	if err != nil {
		return nil, err
	}

	resp, err := c.doer.Do(req)
	if err != nil {
		return nil, err
	}

	return checkResponse(resp)
}

// ──────────────────────────────────────────────────────────────
//  Response helpers
// ──────────────────────────────────────────────────────────────

// checkResponse validates an HTTP response and decodes API errors if present.
// It closes the response body in error cases to avoid leaks.
func checkResponse(r *http.Response) (*http.Response, error) {
	sc := r.StatusCode

	if 200 <= sc && sc <= 299 {
		if sc == http.StatusNoContent {
			if r.Body != nil {
				r.Body.Close()
			}
			return nil, nil
		}
		return r, nil
	}

	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}
	}()

	e := &APIResponseError{StatusCode: sc}

	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(e); err == nil {
			return nil, e
		}
	}

	return nil, e
}

// decodeResponse decodes the JSON body of resp into v and closes the body.
func decodeResponse(resp *http.Response, v any) error {
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}
