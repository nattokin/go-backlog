package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	apiVersion = "v2"
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

type Client struct {
	BaseURL *url.URL
	Token   string
	Doer    Doer
	Wrapper Wrapper
	Method  *Method
}

// ──────────────────────────────────────────────────────────────
//  HTTP Method function set
// ──────────────────────────────────────────────────────────────

// Method holds injected HTTP operation functions.
// Each function delegates to Client.do() but can be replaced in tests
// for fine-grained unit testing of individual services.
type Method struct {
	Get    func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	Post   func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Patch  func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Delete func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Upload func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error)
}

// ──────────────────────────────────────────────────────────────
//  Client constructor
// ──────────────────────────────────────────────────────────────

func NewClient(baseURL, token string, opts ...*ClientOption) (*Client, error) {
	if token == "" {
		return nil, NewInternalClientError("missing token")
	}
	if baseURL == "" {
		return nil, NewInternalClientError("missing baseURL")
	}

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	config := &clientConfig{}
	for _, option := range opts {
		if option != nil {
			option.set(config)
		}
	}

	if config.Doer == nil {
		config.Doer = http.DefaultClient
	}

	c := &Client{
		BaseURL: u,
		Doer:    config.Doer,
		Token:   token,
		Wrapper: &DefaultWrapper{},
	}

	// --- Inject HTTP Method Wrappers ------------------------------------------
	c.Method = &Method{
		Get:    c.Get,
		Post:   c.Post,
		Patch:  c.Patch,
		Delete: c.Delete,
		Upload: c.Upload,
	}

	return c, nil
}

// ──────────────────────────────────────────────────────────────
//  HTTP request creation and execution
// ──────────────────────────────────────────────────────────────

func (c *Client) NewRequest(ctx context.Context, Method, spath string, opts ...*HttpRequestOption) (*http.Request, error) {
	if spath == "" {
		return nil, errors.New("spath must not be empty")
	}

	config := &httpRequestConfig{}
	for _, option := range opts {
		if option != nil {
			option.set(config)
		}
	}

	u := *c.BaseURL
	u.Path = path.Join(u.Path, "api", apiVersion, spath)
	if config.Query != nil {
		u.RawQuery = config.Query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, Method, u.String(), config.Body)
	if err != nil {
		return nil, err
	}

	if config.Header != nil {
		req.Header = config.Header.Clone()
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	return req, nil
}

// Do executes the given HTTP request using the injected Doer.
// All HTTP calls pass through this function, ensuring consistent error handling.
func (c *Client) Do(ctx context.Context, Method, spath string, opts ...*HttpRequestOption) (*http.Response, error) {
	req, err := c.NewRequest(ctx, Method, spath, opts...)
	if err != nil {
		return nil, err
	}

	resp, err := c.Doer.Do(req)
	if err != nil {
		return nil, err
	}

	return CheckResponse(resp)
}

func (c *Client) Get(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, spath, WithQuery(query))
}

func (c *Client) Post(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	if form == nil {
		form = url.Values{}
	}
	return c.Do(ctx, http.MethodPost, spath, WithHeader(header), WithBody(strings.NewReader(form.Encode())))
}

func (c *Client) Patch(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	if form == nil {
		form = url.Values{}
	}
	return c.Do(ctx, http.MethodPatch, spath, WithHeader(header), WithBody(strings.NewReader(form.Encode())))
}

func (c *Client) Delete(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	if form == nil {
		form = url.Values{}
	}
	return c.Do(ctx, http.MethodDelete, spath, WithHeader(header), WithBody(strings.NewReader(form.Encode())))
}

func (c *Client) Upload(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error) {
	if fileName == "" {
		return nil, NewInternalClientError("fileName is required")
	}
	var buf bytes.Buffer
	mw := c.Wrapper.NewMultipartWriter(&buf)

	fw, err := mw.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	if err = c.Wrapper.Copy(fw, r); err != nil {
		return nil, err
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Set("Content-Type", mw.FormDataContentType())

	return c.Do(ctx, http.MethodPost, spath, WithHeader(header), WithBody(&buf))
}

// ──────────────────────────────────────────────────────────────
//  Response helpers
// ──────────────────────────────────────────────────────────────

// CheckResponse validates an HTTP response and decodes API errors if present.
// It closes the response body in error cases to avoid leaks.
func CheckResponse(r *http.Response) (*http.Response, error) {
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

// DecodeResponse decodes the JSON body of resp into v and closes the body.
func DecodeResponse(resp *http.Response, v any) error {
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

// Error represents one of Backlog API response errors.
type Error struct {
	// Message is the detailed error message from the API.
	Message  string `json:"message,omitempty"`
	Code     int    `json:"code,omitempty"`
	MoreInfo string `json:"moreInfo,omitempty"`
}

// Error returns the API error message.
func (e *Error) Error() string {
	msg := fmt.Sprintf("Message:%s, Code:%d", e.Message, e.Code)

	if e.MoreInfo == "" {
		return msg
	}

	return msg + ", MoreInfo:" + e.MoreInfo
}

// InternalClientError represents client-side configuration or usage errors.
// It is distinct from API-level errors and indicates issues like missing Token
// or malformed base URL.
type InternalClientError struct {
	msg string
}

func (e *InternalClientError) Error() string {
	return e.msg
}

func NewInternalClientError(msg string) *InternalClientError {
	return &InternalClientError{msg: msg}
}

// APIResponseError represents Error Response of Backlog API.
type APIResponseError struct {
	StatusCode int      `json:"-"` // HTTP status code (4xx or 5xx)
	Errors     []*Error `json:"errors,omitempty"`
}

// Error returns all error messages in APIResponseError.
func (e *APIResponseError) Error() string {
	msgs := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}

	return fmt.Sprintf("Status Code:%d\n%s", e.StatusCode, strings.Join(msgs, "\n"))
}
