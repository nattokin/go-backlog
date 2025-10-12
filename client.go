package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	apiVersion = "v2"
)

// ──────────────────────────────────────────────────────────────
//  Error types
// ──────────────────────────────────────────────────────────────

// InternalClientError represents client-side configuration or usage errors.
// It is distinct from API-level errors and indicates issues like missing token
// or malformed base URL.
type InternalClientError struct {
	msg string
}

func (e *InternalClientError) Error() string {
	return e.msg
}

func newInternalClientError(msg string) *InternalClientError {
	return &InternalClientError{msg: msg}
}

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
	doer    Doer
	token   string
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
	Get    func(spath string, query *QueryParams) (*http.Response, error)
	Post   func(spath string, form *FormParams) (*http.Response, error)
	Patch  func(spath string, form *FormParams) (*http.Response, error)
	Delete func(spath string, form *FormParams) (*http.Response, error)
	Upload func(spath, fileName string, r io.Reader) (*http.Response, error)
}

// ──────────────────────────────────────────────────────────────
//  Parameter types
// ──────────────────────────────────────────────────────────────

// FormParams represents form-encoded parameters used in Backlog API POST,
// PATCH, and DELETE requests. It wraps url.Values to provide helper methods
// for constructing form data and converting it to io.Reader for HTTP bodies.
//
// Typical usage:
//
//	form := NewFormParams()
//	form.Add("name", "ProjectX")
//	resp, err := c.method.Post("projects", form)
type FormParams struct {
	*url.Values
}

// NewFormParams creates and returns a new, initialized FormParams instance.
// It is primarily used to prepare form data for POST, PATCH, and DELETE requests.
func NewFormParams() *FormParams {
	return &FormParams{&url.Values{}}
}

// NewReader returns an io.Reader containing the URL-encoded form data.
// This is used to provide the request body to an HTTP client.
func (p *FormParams) NewReader() io.Reader {
	return strings.NewReader(p.Encode())
}

// QueryParams represents query string parameters used in Backlog API GET requests.
// It wraps url.Values to simplify parameter creation and encoding.
//
// Typical usage:
//
//	query := NewQueryParams()
//	query.Add("projectId", "123")
//	resp, err := c.method.Get("issues", query)
type QueryParams struct {
	*url.Values
}

// NewQueryParams creates and returns a new, initialized QueryParams instance.
// It is typically used to construct URL query strings for GET requests.
func NewQueryParams() *QueryParams {
	return &QueryParams{&url.Values{}}
}

// ──────────────────────────────────────────────────────────────
//  Wrapper interface for I/O abstractions
// ──────────────────────────────────────────────────────────────

type wrapper interface {
	Copy(dst io.Writer, src io.Reader) error
	NewMultipartWriter(w io.Writer) multipartWriter
}

type multipartWriter interface {
	CreateFormFile(fieldname, filename string) (io.Writer, error)
	FormDataContentType() string
	Close() error
}

// ──────────────────────────────────────────────────────────────
//  Default wrapper implementations
// ──────────────────────────────────────────────────────────────

type defaultWrapper struct{}

func (*defaultWrapper) Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

func (*defaultWrapper) NewMultipartWriter(w io.Writer) multipartWriter {
	return &defaultMultipartWriter{multipart.NewWriter(w)}
}

type defaultMultipartWriter struct {
	*multipart.Writer
}

func (mw *defaultMultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	return mw.Writer.CreateFormFile(fieldname, filename)
}

func (mw *defaultMultipartWriter) FormDataContentType() string {
	return mw.Writer.FormDataContentType()
}

func (mw *defaultMultipartWriter) Close() error {
	return mw.Writer.Close()
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
		Get: func(spath string, query *QueryParams) (*http.Response, error) {
			return c.do(http.MethodGet, spath, nil, nil, query)
		},
		Post: func(spath string, form *FormParams) (*http.Response, error) {
			header := http.Header{}
			header.Set("Content-Type", "application/x-www-form-urlencoded")
			if form == nil {
				form = NewFormParams()
			}
			return c.do(http.MethodPost, spath, header, form.NewReader(), nil)
		},
		Patch: func(spath string, form *FormParams) (*http.Response, error) {
			header := http.Header{}
			header.Set("Content-Type", "application/x-www-form-urlencoded")
			if form == nil {
				form = NewFormParams()
			}
			return c.do(http.MethodPatch, spath, header, form.NewReader(), nil)
		},
		Delete: func(spath string, form *FormParams) (*http.Response, error) {
			header := http.Header{}
			header.Set("Content-Type", "application/x-www-form-urlencoded")
			if form == nil {
				form = NewFormParams()
			}
			return c.do(http.MethodDelete, spath, header, form.NewReader(), nil)
		},
		Upload: func(spath, fileName string, r io.Reader) (*http.Response, error) {
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

			return c.do(http.MethodPost, spath, header, &buf, nil)
		},
	}

	// ──────────────────────────────────────────────────────────────
	//  Service initialization
	// ──────────────────────────────────────────────────────────────

	// --- Initialize shared option services --------------------------------------
	// Option services provide reusable form and query parameter builders
	// used across multiple Backlog API services.
	optionSupport := &optionSupport{
		query: &QueryOptionService{},
		form:  &FormOptionService{},
	}

	// ActivityOptionService wraps shared optionSupport to be reused
	// by activity-related services such as ProjectActivityService or SpaceActivityService.
	activityOptionService := &ActivityOptionService{
		support: optionSupport,
	}

	// --- Initialize IssueService -------------------------------------------------
	// Provides methods for issue management and file attachment operations.
	c.Issue = &IssueService{
		method: c.method,
		Attachment: &IssueAttachmentService{
			method: c.method,
		},
	}

	// --- Initialize ProjectService ----------------------------------------------
	// Includes sub-services for project activities, users, and project options.
	c.Project = &ProjectService{
		method: c.method,
		Activity: &ProjectActivityService{
			method: c.method,
			Option: activityOptionService,
		},
		User: &ProjectUserService{
			method: c.method,
		},
		Option: &ProjectOptionService{
			support: optionSupport,
		},
	}

	// --- Initialize PullRequestService ------------------------------------------
	// Handles pull request operations and related file attachments.
	c.PullRequest = &PullRequestService{
		method: c.method,
		Attachment: &PullRequestAttachmentService{
			method: c.method,
		},
	}

	// --- Initialize SpaceService -------------------------------------------------
	// Provides access to space-level APIs including activities and attachments.
	c.Space = &SpaceService{
		method: c.method,
		Activity: &SpaceActivityService{
			method: c.method,
			Option: activityOptionService,
		},
		Attachment: &SpaceAttachmentService{
			method: c.method,
		},
	}

	// --- Initialize UserService --------------------------------------------------
	// Provides APIs related to user activities and user option settings.
	c.User = &UserService{
		method: c.method,
		Activity: &UserActivityService{
			method: c.method,
			Option: activityOptionService,
		},
		Option: &UserOptionService{
			support: optionSupport,
		},
	}

	// --- Initialize WikiService --------------------------------------------------
	// Provides wiki page APIs, including file attachments and option configurations.
	c.Wiki = &WikiService{
		method: c.method,
		Attachment: &WikiAttachmentService{
			method: c.method,
		},
		Option: &WikiOptionService{
			support: optionSupport,
		},
	}

	return c, nil
}

// ──────────────────────────────────────────────────────────────
//  HTTP request creation and execution
// ──────────────────────────────────────────────────────────────

// newRequest builds a new HTTP request with token-based authentication.
func (c *Client) newRequest(method, spath string, header http.Header, body io.Reader, query *QueryParams) (*http.Request, error) {
	if spath == "" {
		return nil, errors.New("spath must not be empty")
	}

	if query == nil {
		query = NewQueryParams()
	}
	query.Set("apiKey", c.token)

	u := *c.baseURL
	u.Path = path.Join(u.Path, "api", apiVersion, spath)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header = header
	return req, nil
}

// do executes the given HTTP request using the injected Doer.
// All HTTP calls pass through this function, ensuring consistent error handling.
func (c *Client) do(method, spath string, header http.Header, body io.Reader, query *QueryParams) (*http.Response, error) {
	req, err := c.newRequest(method, spath, header, body, query)
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
//  Response validation
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
