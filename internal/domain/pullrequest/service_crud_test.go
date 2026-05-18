package pullrequest_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/pullrequest"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_One(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantNumber  int
	}{
		"success": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/1", spath)
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo1",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-prNumber": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       0,
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := pullrequest.NewService(method)
			got, err := s.One(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantNumber, got.Number)
		})
	}
}

func TestService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		summary        string
		description    string
		base           string
		branch         string
		opts           []core.RequestOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantNumber  int
	}{
		"success-required-only": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "test description",
			base:           "main",
			branch:         "feature/foo",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests", spath)
				assert.Equal(t, "test PR", form.Get("summary"))
				assert.Equal(t, "test description", form.Get("description"))
				assert.Equal(t, "main", form.Get("base"))
				assert.Equal(t, "feature/foo", form.Get("branch"))
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"success-with-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "test description",
			base:           "main",
			branch:         "feature/foo",
			opts: []core.RequestOption{
				o.WithAssigneeID(5),
				o.WithIssueID(10),
				o.WithNotifiedUserIDs([]int{1, 2}),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "5", form.Get("assigneeId"))
				assert.Equal(t, "10", form.Get("issueId"))
				assert.Equal(t, []string{"1", "2"}, form["notifiedUserId[]"])
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-summary": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-base": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "",
			branch:         "feature/foo",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-branch": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "",
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-assigneeID": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			opts:           []core.RequestOption{o.WithAssigneeID(0)},
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}

			s := pullrequest.NewService(method)
			got, err := s.Create(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.summary, tc.description, tc.base, tc.branch, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantNumber, got.Number)
		})
	}
}

func TestService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int
		option         core.RequestOption
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantNumber  int
	}{
		"success-summary": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithSummary("Updated PR"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/1", spath)
				assert.Equal(t, "Updated PR", form.Get("summary"))
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"success-with-comment": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithSummary("Updated PR"),
			opts:           []core.RequestOption{o.WithComment("looks good")},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "Updated PR", form.Get("summary"))
				assert.Equal(t, "looks good", form.Get("comment"))
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"success-with-issueID": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithIssueID(42),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "42", form.Get("issueId"))
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithSummary("x"),
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			prNumber:       1,
			option:         o.WithSummary("x"),
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-prNumber": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       0,
			option:         o.WithSummary("x"),
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-assigneeID": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithAssigneeID(0),
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithSummary("x"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithSummary("x"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				method.Patch = tc.mockPatchFn
			}

			s := pullrequest.NewService(method)
			got, err := s.Update(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber, tc.option, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantNumber, got.Number)
		})
	}
}
