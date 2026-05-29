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

func TestService_One(t *testing.T) {
	cases := map[string]struct {
		wikiID int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			wikiID: 112,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/112", spath)
				return mock.NewResponse(fixture.Wiki.MinimumJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-wikiID-zero": {
			wikiID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-wikiID-negative": {
			wikiID:                 -1,
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-client-network": {
			wikiID: 112,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 112,
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
			w, err := s.One(context.Background(), tc.wikiID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, w)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, w)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, w)
			assert.Equal(t, 34, w.ID)
			assert.Equal(t, "Minimum Wiki Page", w.Name)
		})
	}
}

func TestService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectID int
		name      string
		content   string
		opts      []*core.APIParamOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-minimum": {
			projectID: 1,
			name:      "Minimum Wiki Page",
			content:   "This is a minimal wiki page.",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", form.Get("projectId"))
				assert.Equal(t, "Minimum Wiki Page", form.Get("name"))
				assert.Equal(t, "This is a minimal wiki page.", form.Get("content"))
				assert.Equal(t, "", form.Get("mailNotify"))
				return mock.NewResponse(fixture.Wiki.MinimumJSON), nil
			},
		},
		"success-with-options": {
			projectID: 1,
			name:      "Wiki with options",
			content:   "content",
			opts:      []*core.APIParamOption{o.WithMailNotify(true)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "true", form.Get("mailNotify"))
				return mock.NewResponse(fixture.Wiki.MinimumJSON), nil
			},
		},

		// --- validation errors: fixed option only ---
		"error-validation-name-empty": {
			projectID:              1,
			name:                   "",
			content:                "content",
			wantValidationErrCount: 1,
		},
		"error-validation-content-empty": {
			projectID:              1,
			name:                   "Test",
			content:                "",
			wantValidationErrCount: 1,
		},
		"error-validation-fixed-options-both-empty": {
			projectID:              1,
			name:                   "",
			content:                "",
			wantValidationErrCount: 2,
		},

		// --- validation errors: argument only ---
		"error-validation-projectID-zero": {
			projectID:              0,
			name:                   "Test",
			content:                "content",
			wantValidationErrCount: 1,
		},

		// --- validation errors: all (argument + fixed option + optional opts) ---
		"error-validation-all": {
			projectID:              0,
			name:                   "",
			content:                "",
			opts:                   []*core.APIParamOption{o.WithContent("")},
			wantValidationErrCount: 3,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectID:              1,
			name:                   "Test",
			content:                "content",
			opts:                   []*core.APIParamOption{o.WithMailNotify(true), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectID:              0,
			name:                   "",
			content:                "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectID:   1,
			name:        "Test",
			content:     "content",
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectID:   0,
			name:        "",
			content:     "",
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			projectID:   1,
			name:        "Test",
			content:     "content",
			opts:        []*core.APIParamOption{mock.NewFailingSetOption(core.ParamMailNotify)},
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			projectID: 1,
			name:      "Test",
			content:   "content",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectID: 1,
			name:      "Test",
			content:   "content",
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

			s := wiki.NewService(method)
			w, err := s.Create(context.Background(), tc.projectID, tc.name, tc.content, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, w)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, w)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, w)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, w)
			assert.Equal(t, 34, w.ID)
			assert.Equal(t, "Minimum Wiki Page", w.Name)
			assert.Equal(t, "This is a minimal wiki page.", w.Content)
		})
	}
}

func TestService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		wikiID int
		option *core.APIParamOption
		opts   []*core.APIParamOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-name-only": {
			wikiID: 34,
			option: o.WithName("New Page Name"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "New Page Name", form.Get("name"))
				return mock.NewResponse(fixture.Wiki.MaximumJSON), nil
			},
		},
		"success-content-only": {
			wikiID: 34,
			option: o.WithContent("Updated content."),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "Updated content.", form.Get("content"))
				return mock.NewResponse(fixture.Wiki.MaximumJSON), nil
			},
		},
		"success-full-options": {
			wikiID: 34,
			option: o.WithName("Full Options Name"),
			opts: []*core.APIParamOption{
				o.WithContent("Full Options Content"),
				o.WithMailNotify(true),
			},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "Full Options Name", form.Get("name"))
				assert.Equal(t, "Full Options Content", form.Get("content"))
				assert.Equal(t, "true", form.Get("mailNotify"))
				return mock.NewResponse(fixture.Wiki.MaximumJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-wikiID-zero": {
			wikiID:                 0,
			option:                 o.WithName("x"),
			wantValidationErrCount: 1,
		},
		"error-validation-wikiID-negative": {
			wikiID:                 -1,
			option:                 o.WithName("x"),
			wantValidationErrCount: 1,
		},

		// --- validation errors: option-only (single invalid option) ---
		"error-validation-opt-single": {
			wikiID:                 34,
			option:                 o.WithMailNotify(true),
			opts:                   []*core.APIParamOption{o.WithContent("")},
			wantValidationErrCount: 1,
		},
		"error-validation-fixed-option-empty-name": {
			wikiID:                 34,
			option:                 o.WithName(""),
			wantValidationErrCount: 1,
		},

		// --- validation errors: option-only (multiple invalid options) ---
		"error-validation-opt-multiple": {
			wikiID:                 34,
			option:                 o.WithMailNotify(true),
			opts:                   []*core.APIParamOption{o.WithName(""), o.WithContent("")},
			wantValidationErrCount: 2,
		},

		// --- validation errors: all (argument + fixed option + optional opts) ---
		"error-validation-all": {
			wikiID:                 0,
			option:                 o.WithMailNotify(true),
			opts:                   []*core.APIParamOption{o.WithContent("")},
			wantValidationErrCount: 3,
		},

		// --- validation errors: no name/content provided ---
		"error-validation-no-name-or-content": {
			wikiID:                 34,
			option:                 o.WithMailNotify(true),
			wantValidationErrCount: 1,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			wikiID:                 34,
			option:                 o.WithName("x"),
			opts:                   []*core.APIParamOption{o.WithMailNotify(true), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			wikiID:                 0,
			option:                 o.WithName(""),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			wikiID:      34,
			option:      o.WithName("x"),
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			wikiID:      0,
			option:      o.WithName(""),
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			wikiID:      34,
			option:      o.WithName("x"),
			opts:        []*core.APIParamOption{mock.NewFailingSetOption(core.ParamMailNotify)},
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			wikiID: 34,
			option: o.WithName("Updated Name"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 34,
			option: o.WithName("Updated Name"),
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

			s := wiki.NewService(method)
			w, err := s.Update(context.Background(), tc.wikiID, tc.option, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, w)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, w)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, w)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, w)
			assert.Equal(t, 34, w.ID)
			assert.Equal(t, "Maximum Wiki Page", w.Name)
		})
	}
}

func TestService_Delete(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		wikiID int
		opts   []*core.APIParamOption

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-no-option": {
			wikiID: 1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				return mock.NewResponse(fixture.Wiki.MaximumJSON), nil
			},
		},
		"success-with-option": {
			wikiID: 34,
			opts:   []*core.APIParamOption{o.WithMailNotify(true)},
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "true", form.Get("mailNotify"))
				return mock.NewResponse(fixture.Wiki.MaximumJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-wikiID-zero": {
			wikiID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-wikiID-negative": {
			wikiID:                 -1,
			wantValidationErrCount: 1,
		},

		// --- validation errors: option only ---
		"error-option-validation-with-valid-arg": {
			wikiID:                 1,
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamMailNotify)},
			wantValidationErrCount: 1,
		},

		// --- validation errors: option + argument ---
		"error-validation-all": {
			wikiID:                 0,
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamMailNotify)},
			wantValidationErrCount: 2,
		},

		// --- fail-fast: nil option among valid values ---
		"error-nil-option-with-valid-values": {
			wikiID:                 1,
			opts:                   []*core.APIParamOption{o.WithMailNotify(true), nil},
			wantInvalidOptionError: true,
		},
		// --- fail-fast: nil option among invalid values ---
		"error-nil-option-with-invalid-values": {
			wikiID:                 0,
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			wikiID:      1,
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			wikiID:      0,
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			wikiID:      1,
			opts:        []*core.APIParamOption{mock.NewFailingSetOption(core.ParamMailNotify)},
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			wikiID: 34,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 34,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}

			s := wiki.NewService(method)
			w, err := s.Delete(context.Background(), tc.wikiID, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, w)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, w)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, w)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, w)
			assert.Equal(t, 34, w.ID)
		})
	}
}
