package wiki_test

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
	"github.com/nattokin/go-backlog/internal/domain/wiki"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		opts           []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-minimum": {
			projectIDOrKey: "103",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "103", query.Get("projectIdOrKey"))
				return mock.NewResponse(fixture.Wiki.ListJSON), nil
			},
		},
		"success-with-option": {
			projectIDOrKey: "PRJ_KEY",
			opts:           []*core.APIParamOption{o.WithKeyword("test")},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "PRJ_KEY", query.Get("projectIdOrKey"))
				assert.Equal(t, "test", query.Get("keyword"))
				return mock.NewResponse(fixture.Wiki.ListJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},

		// --- validation errors: option only ---
		"error-option-validation-with-valid-arg": {
			projectIDOrKey:         "PRJ",
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamKeyword)},
			wantValidationErrCount: 1,
		},

		// --- validation errors: option + argument ---
		"error-validation-all": {
			projectIDOrKey:         "",
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamKeyword)},
			wantValidationErrCount: 2,
		},

		// --- fail-fast: nil option among valid values ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "PRJ",
			opts:                   []*core.APIParamOption{o.WithKeyword("test"), nil},
			wantInvalidOptionError: true,
		},
		// --- fail-fast: nil option among validation-erroring values ---
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "PRJ",
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
			projectIDOrKey: "PRJ",
			opts:           []*core.APIParamOption{mock.NewFailingSetOption(core.ParamKeyword)},
			wantErrType:    errors.New(""),
		},
		"error-client-network": {
			projectIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "1",
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

			s := wiki.NewService(method)
			wikis, err := s.List(context.Background(), tc.projectIDOrKey, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, wikis)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, wikis)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, wikis)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, wikis)
			assert.Len(t, wikis, 2)
			assert.Equal(t, 112, wikis[0].ID)
			assert.Equal(t, "test1", wikis[0].Name)
		})
	}
}

func TestService_Count(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantCount              int
	}{
		"success-id": {
			projectIDOrKey: "103",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "103", query.Get("projectIdOrKey"))
				return mock.NewResponse(`{"count": 34}`), nil
			},
			wantCount: 34,
		},
		"success-key": {
			projectIDOrKey: "PRJ_KEY",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "PRJ_KEY", query.Get("projectIdOrKey"))
				return mock.NewResponse(`{"count": 10}`), nil
			},
			wantCount: 10,
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

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "1",
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

			s := wiki.NewService(method)
			count, err := s.Count(context.Background(), tc.projectIDOrKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Equal(t, 0, count)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Equal(t, 0, count)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}
