package issue_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService_Create_CustomField(t *testing.T) {
	cases := map[string]struct {
		opts       []core.RequestOption
		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"WithCustomFieldItem-multiple": {
			opts: []core.RequestOption{
				issue.WithCustomFieldItem(105, 201),
				issue.WithCustomFieldItem(105, 202),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"201", "202"}, form["customField_105[]"])
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomFieldOther": {
			opts: []core.RequestOption{issue.WithCustomFieldOther(105, "custom text")},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "custom text", form.Get("customField_105_otherValue"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Post = tc.mockPostFn

			s := issue.NewService(method)
			_, err := s.Create(context.Background(), 10, "summary", 2, 3, tc.opts...)
			require.NoError(t, err)
		})
	}
}

func TestIssueService_Update_CustomField(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option      core.RequestOption
		opts        []core.RequestOption
		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"WithCustomFieldItem": {
			option: o.WithSummary("Updated"),
			opts:   []core.RequestOption{issue.WithCustomFieldItem(203, 400)},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"400"}, form["customField_203[]"])
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomFieldOther": {
			option: o.WithSummary("Updated"),
			opts:   []core.RequestOption{issue.WithCustomFieldOther(204, "other val")},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "other val", form.Get("customField_204_otherValue"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Patch = tc.mockPatchFn

			s := issue.NewService(method)
			_, err := s.Update(context.Background(), "PRJ-1", tc.option, tc.opts...)
			require.NoError(t, err)
		})
	}
}
