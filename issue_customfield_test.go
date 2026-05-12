package backlog_test

import (
	"context"
	"net/http"
	"testing"
	"time"

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
		"Create/WithCustomField_string": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "v1.2.3", req.PostForm.Get("customField_101"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Create(ctx, 10, "New issue", 2, 3,
					backlog.WithCustomField(101, "v1.2.3"),
				)
				require.NoError(t, err)
			},
		},
		"Create/WithCustomField_int": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "42", req.PostForm.Get("customField_102"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Create(ctx, 10, "New issue", 2, 3,
					backlog.WithCustomField(102, 42),
				)
				require.NoError(t, err)
			},
		},
		"Create/WithCustomField_float64": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "3.14", req.PostForm.Get("customField_103"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Create(ctx, 10, "New issue", 2, 3,
					backlog.WithCustomField(103, 3.14),
				)
				require.NoError(t, err)
			},
		},
		"Create/WithCustomField_time": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "2024-06-01", req.PostForm.Get("customField_104"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Create(ctx, 10, "New issue", 2, 3,
					backlog.WithCustomField(104, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)),
				)
				require.NoError(t, err)
			},
		},
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
		"Update/WithCustomField_string": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "new value", req.PostForm.Get("customField_201"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Update(ctx, "PRJ-1",
					c.Issue.Option.WithSummary("Updated"),
					backlog.WithCustomField(201, "new value"),
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

func TestWithCustomField_Keys(t *testing.T) {
	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"string":  {option: backlog.WithCustomField(1, "v"), wantKey: "customField_1"},
		"int":     {option: backlog.WithCustomField(2, 0), wantKey: "customField_2"},
		"float64": {option: backlog.WithCustomField(3, 0.0), wantKey: "customField_3"},
		"time":    {option: backlog.WithCustomField(4, time.Time{}), wantKey: "customField_4"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
