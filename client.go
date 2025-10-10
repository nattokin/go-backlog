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

// InternalClientError is an error type for client-side issues (e.g., missing token, URL parsing failure).
type InternalClientError struct {
	msg string
}

func (e *InternalClientError) Error() string {
	return e.msg
}

func newInternalClientError(msg string) *InternalClientError {
	return &InternalClientError{msg: msg}
}

// Client is a Backlog API client.
type Client struct {
	url        *url.URL
	httpClient *http.Client
	token      string
	wrapper    wrapper
	method     *method

	Issue       *IssueService
	Project     *ProjectService
	PullRequest *PullRequestService
	Space       *SpaceService
	User        *UserService
	Wiki        *WikiService
}

// FormParams wraps url.Values.
type FormParams struct {
	*url.Values
}

// NewFormParams returns new FormParams.
func NewFormParams() *FormParams {
	return &FormParams{&url.Values{}}
}

// NewReader converts FormParams to io.Reader.
func (p *FormParams) NewReader() io.Reader {
	return strings.NewReader(p.Encode())
}

// QueryParams represents query parameters for a request.
type QueryParams struct {
	*url.Values
}

// NewQueryParams returns new QueryParams.
func NewQueryParams() *QueryParams {
	return &QueryParams{&url.Values{}}
}

type clientGet func(spath string, query *QueryParams) (*http.Response, error)
type clientPost func(spath string, form *FormParams) (*http.Response, error)
type clientPatch func(spath string, form *FormParams) (*http.Response, error)
type clientDelete func(spath string, form *FormParams) (*http.Response, error)
type clientUpload func(spath, fileName string, r io.Reader) (*http.Response, error)

type method struct {
	Get    clientGet
	Post   clientPost
	Patch  clientPatch
	Delete clientDelete
	Upload clientUpload
}

type wrapper interface {
	Copy(dst io.Writer, src io.Reader) error
	NewMultipartWriter(w io.Writer) multipartWriter
}

type multipartWriter interface {
	CreateFormFile(fieldname, filename string) (io.Writer, error)
	FormDataContentType() string
	Close() error
}

//--------------------------------------
// Default implementations
//--------------------------------------

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

// NewClient creates a new Backlog API Client.
func NewClient(baseURL, token string) (*Client, error) {
	if len(token) == 0 {
		return nil, newInternalClientError("missing token")
	}

	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		url:        parsedURL,
		httpClient: http.DefaultClient,
		token:      token,
		wrapper:    &defaultWrapper{},
	}

	c.method = &method{
		Get: func(spath string, query *QueryParams) (*http.Response, error) {
			return c.get(spath, query)
		},
		Post: func(spath string, form *FormParams) (*http.Response, error) {
			return c.post(spath, form)
		},
		Patch: func(spath string, form *FormParams) (*http.Response, error) {
			return c.patch(spath, form)
		},
		Delete: func(spath string, form *FormParams) (*http.Response, error) {
			return c.delete(spath, form)
		},
		Upload: func(spath, fileName string, r io.Reader) (*http.Response, error) {
			return c.upload(spath, fileName, r)
		},
	}

	optionSupport := &optionSupport{
		query: &QueryOptionService{},
		form:  &FormOptionService{},
	}

	activityOptionService := &ActivityOptionService{
		support: optionSupport,
	}

	c.Issue = &IssueService{
		method: c.method,
		Attachment: &IssueAttachmentService{
			method: c.method,
		},
	}
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
	c.PullRequest = &PullRequestService{
		method: c.method,
		Attachment: &PullRequestAttachmentService{
			method: c.method,
		},
	}
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

// newRequest creates a new HTTP request for the Backlog API.
func (c *Client) newRequest(method, spath string, header http.Header, body io.Reader, query *QueryParams) (*http.Request, error) {
	if spath == "" {
		return nil, errors.New("spath must not be empty")
	}

	if query == nil {
		query = NewQueryParams()
	}
	query.Set("apiKey", c.token)

	u := *c.url
	u.Path = path.Join(u.Path, "api", apiVersion, spath)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header = header
	return req, nil
}

// do performs the HTTP request and returns the response.
func (c *Client) do(method, spath string, header http.Header, body io.Reader, query *QueryParams) (*http.Response, error) {
	req, err := c.newRequest(method, spath, header, body, query)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return checkResponse(resp)
}

// get performs a GET request to the Backlog API.
func (c *Client) get(spath string, query *QueryParams) (*http.Response, error) {
	return c.do(http.MethodGet, spath, nil, nil, query)
}

// post performs a POST request to the Backlog API.
func (c *Client) post(spath string, form *FormParams) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	if form == nil {
		form = NewFormParams()
	}

	return c.do(http.MethodPost, spath, header, form.NewReader(), nil)
}

// patch performs a PATCH request to the Backlog API.
func (c *Client) patch(spath string, form *FormParams) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	if form == nil {
		form = NewFormParams()
	}

	return c.do(http.MethodPatch, spath, header, form.NewReader(), nil)
}

// delete performs a DELETE request to the Backlog API.
func (c *Client) delete(spath string, form *FormParams) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	if form == nil {
		form = NewFormParams()
	}

	return c.do(http.MethodDelete, spath, header, form.NewReader(), nil)
}

// upload performs a POST request to upload a file to the Backlog API.
func (c *Client) upload(spath, fileName string, r io.Reader) (*http.Response, error) {
	if fileName == "" {
		return nil, newInternalClientError("fname is required")
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
}

// checkResponse checks the HTTP status code. If it indicates an error, it returns an API error.
func checkResponse(r *http.Response) (*http.Response, error) {
	sc := r.StatusCode

	// Check for success status codes (2xx)
	if 200 <= sc && sc <= 299 {
		// Handle 204 No Content
		if sc == http.StatusNoContent {
			if r.Body != nil {
				r.Body.Close()
			}
			return nil, nil
		}
		// Return successful response
		return r, nil
	}

	// Handle error response (4xx/5xx)
	defer func() {
		// Ensure the response body is closed
		if r.Body != nil {
			r.Body.Close()
		}
	}()

	e := &APIResponseError{
		StatusCode: sc,
	}

	// Attempt to decode error details from body
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(e); err == nil {
			// Successfully decoded API error
			return nil, e
		}
	}

	// If decoding fails (invalid JSON or empty body), return the APIResponseError
	// containing only the StatusCode, as the error is from the API service.
	return nil, e
}
