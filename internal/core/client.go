// Package core provides the HTTP client, request/response primitives,
// and option infrastructure used by all service packages.
package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/nattokin/go-backlog/internal/model"
)

const (
	apiVersion = "v2"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL *url.URL
	Token   string
	Doer    Doer
	Wrapper Wrapper
	Method  *Method
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
