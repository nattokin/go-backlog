package backlog_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService_CustomField(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"Create/WithCustomFieldItem": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, []string{"201", "202"}, req.PostForm["customField_105[]"])
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Create(ctx, 10, "New issue", 2, 3,
					c.Issue.Option.WithCustomFieldItem(105, 201),
					c.Issue.Option.WithCustomFieldItem(105, 202),
				)
				require.NoError(t, err)
			},
		},
		"Create/WithCustomFieldOther": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "custom text", req.PostForm.Get("customField_105_otherValue"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Create(ctx, 10, "New issue", 2, 3,
					c.Issue.Option.WithCustomFieldOther(105, "custom text"),
				)
				require.NoError(t, err)
			},
		},
		"Update/WithCustomFieldItem": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, []string{"301"}, req.PostForm["customField_202[]"])
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Update(ctx, "PRJ-1",
					c.Issue.Option.WithSummary("Updated"),
					c.Issue.Option.WithCustomFieldItem(202, 301),
				)
				require.NoError(t, err)
			},
		},
		"Update/WithCustomFieldOther": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "other val", req.PostForm.Get("customField_203_otherValue"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Update(ctx, "PRJ-1",
					c.Issue.Option.WithSummary("Updated"),
					c.Issue.Option.WithCustomFieldOther(203, "other val"),
				)
				require.NoError(t, err)
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

func TestIssueOptionService_CustomField(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.Issue.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithCustomFieldItem":  {option: s.WithCustomFieldItem(1, 10), wantKey: "customField_1[]"},
		"WithCustomFieldOther": {option: s.WithCustomFieldOther(1, "other"), wantKey: "customField_1_otherValue"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
