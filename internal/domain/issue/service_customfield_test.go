package issue_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService_Create_CustomField(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		opts       []core.RequestOption
		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"WithCustomField-string": {
			opts: []core.RequestOption{core.WithCustomField(101, "v1.2.3")},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "v1.2.3", form.Get("customField_101"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomField-int": {
			opts: []core.RequestOption{core.WithCustomField(102, 42)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "42", form.Get("customField_102"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomField-float64": {
			opts: []core.RequestOption{core.WithCustomField(103, 3.14)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "3.14", form.Get("customField_103"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomField-time": {
			opts: []core.RequestOption{core.WithCustomField(104, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "2024-06-01", form.Get("customField_104"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomFieldItem-multiple": {
			opts: []core.RequestOption{
				core.WithCustomFieldItem(105, 201),
				core.WithCustomFieldItem(105, 202),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"201", "202"}, form["customField_105[]"])
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomFieldOther": {
			opts: []core.RequestOption{core.WithCustomFieldOther(105, "custom text")},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "custom text", form.Get("customField_105_otherValue"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"mixed-regular-and-custom": {
			opts: []core.RequestOption{
				o.WithDescription("desc"),
				core.WithCustomField(200, "text val"),
				core.WithCustomFieldItem(201, 300),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "desc", form.Get("description"))
				assert.Equal(t, "text val", form.Get("customField_200"))
				assert.Equal(t, []string{"300"}, form["customField_201[]"])
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
		"WithCustomField-as-primary": {
			option: core.WithCustomField(201, "updated"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "updated", form.Get("customField_201"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomField-in-opts": {
			option: o.WithSummary("Updated"),
			opts:   []core.RequestOption{core.WithCustomField(202, "val")},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "Updated", form.Get("summary"))
				assert.Equal(t, "val", form.Get("customField_202"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomFieldItem": {
			option: o.WithSummary("Updated"),
			opts:   []core.RequestOption{core.WithCustomFieldItem(203, 400)},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"400"}, form["customField_203[]"])
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"WithCustomFieldOther": {
			option: o.WithSummary("Updated"),
			opts:   []core.RequestOption{core.WithCustomFieldOther(204, "other val")},
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
