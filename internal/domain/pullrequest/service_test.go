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

func TestService_Count(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		opts           []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantCount              int
	}{
		"success-no-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/count", spath)
				return mock.NewResponse(`{"count":5}`), nil
			},
			wantCount: 5,
		},
		"success-with-assigneeIDs": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []*core.APIParamOption{o.WithAssigneeIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"10", "20"}, query["assigneeId[]"])
				return mock.NewResponse(`{"count":2}`), nil
			},
			wantCount: 2,
		},

		// --- validation errors: argument only ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo1",
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			wantValidationErrCount: 1,
		},

		// --- validation errors: option only ---
		"error-option-validation-with-valid-args": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamStatusIDs)},
			wantValidationErrCount: 1,
		},

		// --- validation errors: option + arguments ---
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamStatusIDs)},
			wantValidationErrCount: 3,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
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
			count, err := s.Count(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Zero(t, count)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Zero(t, count)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Zero(t, count)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}
