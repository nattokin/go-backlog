package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	apiVersion = "v2"
)

// ClinetError is a description of a Backlog API client error.
type ClinetError struct {
	msg string
}

func (e *ClinetError) Error() string {
	return e.msg
}

func newClientError(msg string) *ClinetError {
	return &ClinetError{msg: msg}
}

// Client is Backlog API client.
type Client struct {
	url        *url.URL
	httpClient *http.Client
	token      string
	wrapper    *wrapper

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

// QueryParams is query parameters for request.
type QueryParams struct {
	*url.Values
}

// NewQueryParams returns new QueryParams.
func NewQueryParams() *QueryParams {
	return &QueryParams{&url.Values{}}
}

// withOptions sets request query parameters from options.
func (p *QueryParams) withOptions(options []*QueryOption, validOptions ...queryType) error {
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return err
		}
	}

	for _, option := range options {
		if err := option.set(p); err != nil {
			return err
		}
	}

	return nil
}

type clientGet func(spath string, query *QueryParams) (*http.Response, error)
type clientPost func(spath string, form *FormParams) (*http.Response, error)
type clientPatch func(spath string, form *FormParams) (*http.Response, error)
type clientDelete func(spath string, form *FormParams) (*http.Response, error)
type clientUpload func(spath, fpath, fname string) (*http.Response, error)

type method struct {
	Get    clientGet
	Post   clientPost
	Patch  clientPatch
	Delete clientDelete
	Upload clientUpload
}

type wrapper struct {
	CreateFormFile func(w *multipart.Writer, fname string) (io.Writer, error)
	Copy           func(dst io.Writer, src io.Reader) error
}

func createFormFile(w *multipart.Writer, fname string) (io.Writer, error) {
	return w.CreateFormFile("file", fname)
}

func copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

// NewClient creates a new Backlog API Client.
func NewClient(baseURL, token string) (*Client, error) {
	if len(token) == 0 {
		return nil, newClientError("missing token")
	}

	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		url:        parsedURL,
		httpClient: http.DefaultClient,
		token:      token,
		wrapper: &wrapper{
			CreateFormFile: createFormFile,
			Copy:           copy,
		},
	}

	m := &method{
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
		Upload: func(spath, fpath, fname string) (*http.Response, error) {
			return c.upload(spath, fpath, fname)
		},
	}

	activityOptionService := &ActivityOptionService{}

	c.Issue = &IssueService{
		method: m,
		Attachment: &IssueAttachmentService{
			method: m,
		},
	}
	c.Project = &ProjectService{
		method: m,
		Activity: &ProjectActivityService{
			method: m,
			Option: activityOptionService,
		},
		User: &ProjectUserService{
			method: m,
		},
		Option: &ProjectOptionService{},
	}
	c.PullRequest = &PullRequestService{
		method: m,
		Attachment: &PullRequestAttachmentService{
			method: m,
		},
	}
	c.Space = &SpaceService{
		method: m,
		Activity: &SpaceActivityService{
			method: m,
			Option: activityOptionService,
		},
		Attachment: &SpaceAttachmentService{
			method: m,
		},
	}
	c.User = &UserService{
		method: m,
		Activity: &UserActivityService{
			method: m,
			Option: activityOptionService,
		},
		Option: &UserOptionService{},
	}
	c.Wiki = &WikiService{
		method: m,
		Attachment: &WikiAttachmentService{
			method: m,
		},
		Option: &WikiOptionService{},
	}

	return c, nil
}

// Creates new request.
func (c *Client) newReqest(method, spath string, query *QueryParams, body io.Reader) (*http.Request, error) {
	if spath == "" {
		return nil, errors.New("spath must not empty")
	}

	if query == nil {
		query = NewQueryParams()
	}
	query.Set("apiKey", c.token)

	u := *c.url
	u.Path = path.Join(u.Path, "api", apiVersion, spath)
	u.RawQuery = query.Encode()

	return http.NewRequest(method, u.String(), body)
}

// Do http request, and return Response.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return checkResponse(resp)
}

// Get method of http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) get(spath string, query *QueryParams) (*http.Response, error) {
	req, err := c.newReqest(http.MethodGet, spath, query, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

// Post method of http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) post(spath string, form *FormParams) (*http.Response, error) {
	if form == nil {
		form = NewFormParams()
	}
	req, err := c.newReqest(http.MethodPost, spath, nil, form.NewReader())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

// Patch method of http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) patch(spath string, form *FormParams) (*http.Response, error) {
	if form == nil {
		form = NewFormParams()
	}
	req, err := c.newReqest(http.MethodPatch, spath, nil, form.NewReader())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

// Delete method of http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) delete(spath string, form *FormParams) (*http.Response, error) {
	if form == nil {
		form = NewFormParams()
	}
	req, err := c.newReqest(http.MethodDelete, spath, nil, form.NewReader())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

// Upload file method used http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) upload(spath, fpath, fname string) (*http.Response, error) {
	if fname == "" {
		return nil, newClientError("fname is required")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.Close()

	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fw, err := c.wrapper.CreateFormFile(w, fname)
	if err != nil {
		return nil, err
	}
	if err = c.wrapper.Copy(fw, f); err != nil {
		return nil, err
	}

	req, err := c.newReqest(http.MethodPost, spath, nil, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	return c.do(req)
}

// Check HTTP status code. If it has errors, return error.
func checkResponse(r *http.Response) (*http.Response, error) {
	if sc := r.StatusCode; 200 <= sc && sc <= 299 {
		return r, nil
	}

	if r.Body == nil {
		return nil, newClientError("response body is empty")
	}
	defer r.Body.Close()

	e := &APIResponseError{}
	if err := json.NewDecoder(r.Body).Decode(e); err != nil {
		return nil, err
	}

	return nil, e
}
