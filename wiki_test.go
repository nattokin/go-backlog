package backlog_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestWikiService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/wikis", req.URL.Path)
				assert.Equal(t, "backlog", req.URL.Query().Get("keyword"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.All(ctx, "TEST", c.Wiki.Option.WithKeyword("backlog"))
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"All/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.All(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Count": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/wikis/count", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"count": 34}`))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Count(ctx, "TEST")
				require.NoError(t, err)
				assert.Equal(t, 34, got)
			},
		},
		"Count/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No project.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Count(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/wikis/34", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.One(ctx, 34)
				require.NoError(t, err)
				assert.Equal(t, 34, got.ID)
			},
		},
		"One/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.One(ctx, 34)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/wikis", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "56", req.PostForm.Get("projectId"))
				assert.Equal(t, "Test Wiki", req.PostForm.Get("name"))
				assert.Equal(t, "content", req.PostForm.Get("content"))
				assert.Equal(t, "true", req.PostForm.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MinimumJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Create(ctx, 56, "Test Wiki", "content", c.Wiki.Option.WithMailNotify(true))
				require.NoError(t, err)
				assert.Equal(t, "Minimum Wiki Page", got.Name)
			},
		},
		"Create/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return authErrorResponse(), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Create(ctx, 56, "Test Wiki", "content")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/wikis/34", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "New Name", req.PostForm.Get("name"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Update(ctx, 34, c.Wiki.Option.WithName("New Name"))
				require.NoError(t, err)
				assert.Equal(t, 34, got.ID)
			},
		},
		"Update/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Update(ctx, 34, c.Wiki.Option.WithName("New Name"))
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/wikis/34", req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Delete(ctx, 34, c.Wiki.Option.WithMailNotify(true))
				require.NoError(t, err)
				assert.Equal(t, 34, got.ID)
			},
		},
		"Delete/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Delete(ctx, 34)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestWikiAttachmentService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"Attach": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/wikis/34/attachments", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, []string{"2", "5"}, req.PostForm["attachmentId[]"])
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Attachment.Attach(ctx, 34, []int{2, 5})
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 2, got[0].ID)
				assert.Equal(t, 5, got[1].ID)
			},
		},
		"Attach/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Attachment.Attach(ctx, 34, []int{2, 5})
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/wikis/34/attachments", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Attachment.List(ctx, 34)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 2, got[0].ID)
				assert.Equal(t, 5, got[1].ID)
			},
		},
		"List/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Attachment.List(ctx, 34)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Remove": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/wikis/34/attachments/8", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Attachment.Remove(ctx, 34, 8)
				require.NoError(t, err)
				assert.Equal(t, 8, got.ID)
			},
		},
		"Remove/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such attachment.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Attachment.Remove(ctx, 34, 8)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestWikiHistoryService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/wikis/34/history", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.WikiHistory.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.History.List(ctx, 34)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 34, got[0].PageID)
				assert.Equal(t, 2, got[0].Version)
				assert.Equal(t, "Home", got[0].Name)
				assert.Equal(t, 1, got[1].Version)
			},
		},
		"List/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.History.List(ctx, 34)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestWikiSharedFileService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/wikis/34/sharedFiles", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.SharedFile.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.SharedFile.List(ctx, 34)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 454403, got[0].ID)
				assert.Equal(t, "01_buz.png", got[0].Name)
				assert.Equal(t, 454404, got[1].ID)
				assert.Equal(t, "readme.md", got[1].Name)
			},
		},
		"List/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.SharedFile.List(ctx, 34)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Link": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/wikis/34/sharedFiles", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, []string{"454403", "454404"}, req.PostForm["fileId[]"])
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.SharedFile.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.SharedFile.Link(ctx, 34, []int{454403, 454404})
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 454403, got[0].ID)
				assert.Equal(t, 454404, got[1].ID)
			},
		},
		"Link/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.SharedFile.Link(ctx, 34, []int{454403})
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Unlink": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/wikis/34/sharedFiles/454403", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.SharedFile.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.SharedFile.Unlink(ctx, 34, 454403)
				require.NoError(t, err)
				assert.Equal(t, 454403, got.ID)
				assert.Equal(t, "01_buz.png", got.Name)
				assert.Equal(t, "/icon/", got.Dir)
			},
		},
		"Unlink/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such shared file.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.SharedFile.Unlink(ctx, 34, 454403)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestWikiStarService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/wikis/34/stars", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Star.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Star.List(ctx, 34)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 10, got[0].ID)
				assert.Equal(t, 20, got[1].ID)
			},
		},
		"List/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Star.List(ctx, 34)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Add": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/stars", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "34", req.FormValue("wikiId"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Wiki.Star.Add(ctx, 34)
				require.NoError(t, err)
			},
		},
		"Add/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return authErrorResponse(), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Wiki.Star.Add(ctx, 34)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Remove": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/stars", req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				assert.Equal(t, "42", form.Get("id"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Wiki.Star.Remove(ctx, 42)
				require.NoError(t, err)
			},
		},
		"Remove/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return authErrorResponse(), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Wiki.Star.Remove(ctx, 42)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestWikiOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.Wiki.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithKeyword": {
			option:  s.WithKeyword("backlog"),
			wantKey: core.ParamKeyword.Value(),
		},
		"WithContent": {
			option:  s.WithContent("Wiki page content"),
			wantKey: core.ParamContent.Value(),
		},
		"WithName": {
			option:  s.WithName("How to Use Backlog"),
			wantKey: core.ParamName.Value(),
		},
		"WithMailNotify": {
			option:  s.WithMailNotify(true),
			wantKey: core.ParamMailNotify.Value(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
