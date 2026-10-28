// Package client provides the HTTP client and request/response primitives
// used by all service packages.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/nattokin/go-backlog/internal/model"
)

const (
	apiVersion = "v2"
)

// Error represents one of Backlog API response errors.
type Error struct {
	Message  string `json:"message,omitempty"`
	Code     int    `json:"code,omitempty"`
	MoreInfo string `json:"moreInfo,omitempty"`
}

func (e *Error) Error() string {
	msg := fmt.Sprintf("Message:%s, Code:%d", e.Message, e.Code)

	if e.MoreInfo == "" {
		return msg
	}

	return msg + ", MoreInfo:" + e.MoreInfo
}

// APIResponseError represents Error Response of Backlog API.
type APIResponseError struct {
	StatusCode      int              `json:"-"` // HTTP status code (4xx or 5xx)
	Errors          []*Error         `json:"errors,omitempty"`
	RateLimitStatus *RateLimitStatus `json:"-"` // nil when the rate limit headers are missing or invalid
}

// RateLimitStatus represents the rate limit status reported in the
// X-RateLimit-* response headers. Reset is in UTC epoch seconds.
type RateLimitStatus struct {
	Limit     int
	Remaining int
	Reset     int64
}

func (e *APIResponseError) Error() string {
	msgs := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}

	return fmt.Sprintf("Status Code:%d\n%s", e.StatusCode, strings.Join(msgs, "\n"))
}

// OAuthError represents an OAuth 2.0 error reported in the WWW-Authenticate
// response header (RFC 6750), such as an invalid or expired access token.
// Response holds the underlying API error response.
type OAuthError struct {
	Code        string
	Description string
	Response    *APIResponseError
}

func (e *OAuthError) Error() string {
	msg := fmt.Sprintf("Status Code:%d\nOAuth Error:%s", e.Response.StatusCode, e.Code)

	if e.Description == "" {
		return msg
	}

	return msg + ", Description:" + e.Description
}

// Unwrap returns the underlying *APIResponseError.
func (e *OAuthError) Unwrap() error {
	return e.Response
}

// InternalClientError represents client-side configuration or usage errors.
// It is distinct from API-level errors and indicates issues like a missing
// credential or a malformed base URL.
type InternalClientError struct {
	msg string
}

func (e *InternalClientError) Error() string {
	return e.msg
}

func NewInternalClientError(msg string) *InternalClientError {
	return &InternalClientError{msg: msg}
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// CredentialType identifies how a credential is sent to the Backlog API.
type CredentialType int

const (
	// CredentialAPIKey is an API key sent in the Backlog-API-Key header.
	CredentialAPIKey CredentialType = iota + 1
	// CredentialAccessToken is an OAuth 2.0 access token sent in the Authorization header.
	CredentialAccessToken
)

// Credential holds the value used to authenticate requests.
type Credential struct {
	Type  CredentialType
	Value string
}

type Client struct {
	BaseURL    *url.URL
	Credential Credential
	Doer       Doer
	Wrapper    Wrapper
	Method     *Method
}

// Method holds injected HTTP operation functions.
// Each function delegates to Client.Do() but can be replaced in tests
// for fine-grained unit testing of individual services.
type Method struct {
	Get      func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	Post     func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Patch    func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Put      func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Delete   func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	Upload   func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error)
	Download func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
}

// NewClient creates a Client that authenticates with an API key.
func NewClient(baseURL, apiKey string, opts ...*ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, NewInternalClientError("missing API key")
	}
	return newClient(baseURL, Credential{Type: CredentialAPIKey, Value: apiKey}, opts...)
}

// NewOAuthClient creates a Client that authenticates with an OAuth 2.0 access token.
func NewOAuthClient(baseURL, accessToken string, opts ...*ClientOption) (*Client, error) {
	if accessToken == "" {
		return nil, NewInternalClientError("missing access token")
	}
	return newClient(baseURL, Credential{Type: CredentialAccessToken, Value: accessToken}, opts...)
}

func newClient(baseURL string, credential Credential, opts ...*ClientOption) (*Client, error) {
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
		BaseURL:    u,
		Credential: credential,
		Doer:       config.Doer,
		Wrapper:    &DefaultWrapper{},
	}

	c.Method = &Method{
		Get:      c.Get,
		Post:     c.Post,
		Patch:    c.Patch,
		Put:      c.Put,
		Delete:   c.Delete,
		Upload:   c.Upload,
		Download: c.Download,
	}

	return c, nil
}

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
	// path.Join drops a trailing slash; keep it for endpoints whose path
	// parameter is a directory (e.g. the shared file list).
	if strings.HasSuffix(spath, "/") {
		u.Path += "/"
	}
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
	switch c.Credential.Type {
	case CredentialAPIKey:
		req.Header.Set("Backlog-API-Key", c.Credential.Value)
	case CredentialAccessToken:
		req.Header.Set("Authorization", "Bearer "+c.Credential.Value)
	}

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

func (c *Client) Put(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	if form == nil {
		form = url.Values{}
	}
	return c.Do(ctx, http.MethodPut, spath, WithHeader(header), WithBody(strings.NewReader(form.Encode())))
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

func (c *Client) Download(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, spath, WithQuery(query))
}

// CheckResponse validates an HTTP response and decodes API errors if present.
// When the WWW-Authenticate header carries a Bearer error (RFC 6750), the
// returned error is an *OAuthError wrapping the *APIResponseError.
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

	e := &APIResponseError{StatusCode: sc, RateLimitStatus: parseRateLimitStatus(r.Header)}

	if r.Body != nil {
		// A body that is not a JSON error response leaves Errors empty.
		_ = json.NewDecoder(r.Body).Decode(e)
	}

	if code, description, ok := parseBearerError(r.Header); ok {
		return nil, &OAuthError{Code: code, Description: description, Response: e}
	}

	return nil, e
}

var (
	// bearerChallenge matches a WWW-Authenticate value using the Bearer scheme,
	// allowing a stray leading quote.
	bearerChallenge = regexp.MustCompile(`(?i)^"?bearer\s`)
	// bearerParam matches an error or error_description attribute. The value is
	// captured as a quoted-string body (group 2) or a token (group 3).
	bearerParam = regexp.MustCompile(`(?i)[\s,](error|error_description)\s*=\s*(?:"((?:[^"\\]|\\.)*)"|([^\s,"]+))`)
	// quotedPair matches a backslash escape in a quoted-string.
	quotedPair = regexp.MustCompile(`\\(.)`)
)

// parseBearerError extracts the error and error_description attributes of a
// Bearer challenge in the WWW-Authenticate header (RFC 6750 Section 3).
// ok is false when no Bearer challenge with an error attribute is found.
func parseBearerError(header http.Header) (code, description string, ok bool) {
	for _, v := range header.Values("WWW-Authenticate") {
		v = strings.TrimSpace(v)
		if !bearerChallenge.MatchString(v) {
			continue
		}

		params := map[string]string{}
		for _, m := range bearerParam.FindAllStringSubmatch(v, -1) {
			params[strings.ToLower(m[1])] = quotedPair.ReplaceAllString(m[2], "$1") + m[3]
		}

		if params["error"] != "" {
			return params["error"], params["error_description"], true
		}
	}

	return "", "", false
}

// parseRateLimitStatus parses the X-RateLimit-Limit, X-RateLimit-Remaining
// and X-RateLimit-Reset headers. It returns nil unless all three are present
// and are non-negative integers.
func parseRateLimitStatus(header http.Header) *RateLimitStatus {
	limit, err := strconv.Atoi(strings.TrimSpace(header.Get("X-RateLimit-Limit")))
	if err != nil || limit < 0 {
		return nil
	}

	remaining, err := strconv.Atoi(strings.TrimSpace(header.Get("X-RateLimit-Remaining")))
	if err != nil || remaining < 0 {
		return nil
	}

	reset, err := strconv.ParseInt(strings.TrimSpace(header.Get("X-RateLimit-Reset")), 10, 64)
	if err != nil || reset < 0 {
		return nil
	}

	return &RateLimitStatus{Limit: limit, Remaining: remaining, Reset: reset}
}

// DecodeResponse decodes the JSON body of resp into v and closes the body.
// A nil resp (as returned by CheckResponse for 204 No Content) or a nil body
// has nothing to decode, so v is left untouched and nil is returned.
func DecodeResponse(resp *http.Response, v any) error {
	if resp == nil || resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

// DownloadResponse extracts FileData from a binary HTTP response.
// It parses the filename from Content-Disposition and the media type from Content-Type.
// The caller is responsible for closing FileData.Body.
func DownloadResponse(resp *http.Response) (*model.FileData, error) {
	filename := ""
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		_, params, err := mime.ParseMediaType(cd)
		if err == nil {
			filename = params["filename"]
		}
	}

	contentType := ""
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		mediaType, _, err := mime.ParseMediaType(ct)
		if err == nil {
			contentType = mediaType
		}
	}

	return &model.FileData{
		Body:        resp.Body,
		Filename:    filename,
		ContentType: contentType,
	}, nil
}
