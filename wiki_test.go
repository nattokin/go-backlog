package backlog_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.All(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 2)
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
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/wikis", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "56", req.PostForm.Get("projectId"))
				assert.Equal(t, "Test Wiki", req.PostForm.Get("name"))
				assert.Equal(t, "content", req.PostForm.Get("content"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MinimumJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Create(ctx, 56, "Test Wiki", "content")
				require.NoError(t, err)
				assert.Equal(t, "Minimum Wiki Page", got.Name)
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
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/wikis/34", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Delete(ctx, 34)
				require.NoError(t, err)
				assert.Equal(t, 34, got.ID)
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
