package project_test

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
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestActivityService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		opts           []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-no-option": {
			projectIDOrKey: "TEST",
			opts:           []*core.APIParamOption{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/activities", spath)
				return mock.NewResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withActivityTypeIDs": {
			projectIDOrKey: "TEST",
			opts: []*core.APIParamOption{
				o.WithActivityTypeIDs([]int{1}),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1"}, query["activityTypeId[]"])
				return mock.NewResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withMinID": {
			projectIDOrKey: "TEST",
			opts: []*core.APIParamOption{
				o.WithMinID(1),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("minId"))
				return mock.NewResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withMaxID": {
			projectIDOrKey: "TEST",
			opts: []*core.APIParamOption{
				o.WithMaxID(1),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("maxId"))
				return mock.NewResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withCount": {
			projectIDOrKey: "TEST",
			opts: []*core.APIParamOption{
				o.WithCount(1),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("count"))
				return mock.NewResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withOrder": {
			projectIDOrKey: "TEST",
			opts: []*core.APIParamOption{
				o.WithOrder("asc"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-multiple-options": {
			projectIDOrKey: "TEST",
			opts: []*core.APIParamOption{
				o.WithActivityTypeIDs([]int{1, 2}),
				o.WithMinID(1),
				o.WithMaxID(26),
				o.WithCount(20),
				o.WithOrder("asc"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1", "2"}, query["activityTypeId[]"])
				assert.Equal(t, "1", query.Get("minId"))
				assert.Equal(t, "26", query.Get("maxId"))
				assert.Equal(t, "20", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewResponse(fixture.Activity.ListJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},
		"error-validation-opt-single": {
			projectIDOrKey:         "TEST",
			opts:                   []*core.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			opts:                   []*core.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 2,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			opts:                   []*core.APIParamOption{o.WithCount(10), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			opts:           []*core.APIParamOption{mock.NewFailingSetOption(core.ParamCount)},
			wantErrType:    errors.New(""),
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
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
			s := project.NewActivityService(method)
			got, err := s.List(context.Background(), tc.projectIDOrKey, tc.opts...)

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
		})
	}
}
