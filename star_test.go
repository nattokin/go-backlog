package backlog_test

import (
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
)

func TestStarService_Add(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc  func(req *http.Request) (*http.Response, error)
		call    func(t *testing.T, c *backlog.Client)
		wantErr bool
	}{
		"success-issueId": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/stars", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "1", req.FormValue("issueId"))
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Star.Add(ctx, c.Star.Option.WithIssueID(1))
				require.NoError(t, err)
			},
		},
		"success-commentId": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "5", req.FormValue("commentId"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Star.Add(ctx, c.Star.Option.WithCommentID(5))
				require.NoError(t, err)
			},
		},
		"success-wikiPageId": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "10", req.FormValue("wikiId"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Star.Add(ctx, c.Star.Option.WithWikiPageID(10))
				require.NoError(t, err)
			},
		},
		"success-pullRequestId": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "3", req.FormValue("pullRequestId"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Star.Add(ctx, c.Star.Option.WithPullRequestID(3))
				require.NoError(t, err)
			},
		},
		"success-pullRequestCommentId": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "7", req.FormValue("pullRequestCommentId"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Star.Add(ctx, c.Star.Option.WithPullRequestCommentID(7))
				require.NoError(t, err)
			},
		},
		"error-api": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Star.Add(ctx, c.Star.Option.WithIssueID(1))
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var doer *mockDoer
			if tc.doFunc != nil {
				doer = &mockDoer{do: tc.doFunc}
			} else {
				doer = &mockDoer{do: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("should not be called")
				}}
			}

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(doer))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestStarService_Remove(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		starID  int
		doFunc  func(req *http.Request) (*http.Response, error)
		wantErr bool
	}{
		"success": {
			starID: 42,
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
		},
		"error-invalid-id": {
			starID:  0,
			wantErr: true,
		},
		"error-api": {
			starID: 1,
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var doer *mockDoer
			if tc.doFunc != nil {
				doer = &mockDoer{do: tc.doFunc}
			} else {
				doer = &mockDoer{do: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("should not be called")
				}}
			}

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(doer))
			require.NoError(t, err)

			err = c.Star.Remove(ctx, tc.starID)

			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
