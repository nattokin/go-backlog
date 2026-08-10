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

		wantErrType            error
		wantValidationErrCount int
		wantNumber             int
	}{
		"success": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/1", spath)
				return mock.NewResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo1",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-prNumber-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			prNumber:               0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			wantValidationErrCount: 3,
		},

		// --- other errors ---
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
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
		opts           []*core.APIParamOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantNumber             int
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
				return mock.NewResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"success-with-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			opts: []*core.APIParamOption{
				o.WithAssigneeID(5),
				o.WithIssueID(10),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "5", form.Get("assigneeId"))
				assert.Equal(t, "10", form.Get("issueId"))
				return mock.NewResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo1",
			summary:                "test PR",
			description:            "desc",
			base:                   "main",
			branch:                 "feature/foo",
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			summary:                "test PR",
			description:            "desc",
			base:                   "main",
			branch:                 "feature/foo",
			wantValidationErrCount: 1,
		},
		"error-validation-summary-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			summary:                "",
			description:            "desc",
			base:                   "main",
			branch:                 "feature/foo",
			wantValidationErrCount: 1,
		},
		"error-validation-base-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			summary:                "test PR",
			description:            "desc",
			base:                   "",
			branch:                 "feature/foo",
			wantValidationErrCount: 1,
		},
		"error-validation-branch-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			summary:                "test PR",
			description:            "desc",
			base:                   "main",
			branch:                 "",
			wantValidationErrCount: 1,
		},
		"error-validation-opt-single": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			summary:                "test PR",
			description:            "desc",
			base:                   "main",
			branch:                 "feature/foo",
			opts:                   []*core.APIParamOption{o.WithAssigneeID(0)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			summary:                "",
			description:            "",
			base:                   "",
			branch:                 "",
			wantValidationErrCount: 5,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			summary:                "test PR",
			description:            "desc",
			base:                   "main",
			branch:                 "feature/foo",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			summary:                "",
			description:            "",
			base:                   "",
			branch:                 "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
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
		option         *core.APIParamOption
		opts           []*core.APIParamOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantNumber             int
	}{
		"success-summary": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithSummary("Updated PR"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/1", spath)
				assert.Equal(t, "Updated PR", form.Get("summary"))
				return mock.NewResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"success-with-comment": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         o.WithSummary("Updated PR"),
			opts:           []*core.APIParamOption{o.WithComment("looks good")},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "looks good", form.Get("comment"))
				return mock.NewResponse(fixture.PullRequest.SingleJSON), nil
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
				return mock.NewResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo1",
			prNumber:               1,
			option:                 o.WithSummary("x"),
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			prNumber:               1,
			option:                 o.WithSummary("x"),
			wantValidationErrCount: 1,
		},
		"error-validation-prNumber-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			prNumber:               0,
			option:                 o.WithSummary("x"),
			wantValidationErrCount: 1,
		},
		"error-validation-fixed-option": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			prNumber:               1,
			option:                 o.WithAssigneeID(0),
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			option:                 o.WithSummary(""),
			wantValidationErrCount: 4,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			prNumber:               1,
			option:                 o.WithSummary("x"),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			option:                 o.WithSummary(""),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			prNumber:       1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantNumber, got.Number)
		})
	}
}
