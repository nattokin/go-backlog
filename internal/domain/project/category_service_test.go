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

func TestCategoryService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantLen                int
		wantErrType            error
		wantValidationErrCount int
	}{
		"success-key": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories", spath)
				return mock.NewResponse(fixture.Category.ListJSON), nil
			},
			wantLen: 2,
		},
		"success-id": {
			projectIDOrKey: "6",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/categories", spath)
				return mock.NewResponse(fixture.Category.ListJSON), nil
			},
			wantLen: 2,
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
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
			s := project.NewCategoryService(method)
			categories, err := s.List(context.Background(), tc.projectIDOrKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, categories)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, categories)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, categories)
			assert.Len(t, categories, tc.wantLen)
		})
	}
}

func TestCategoryService_Create(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		name           string

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories", spath)
				assert.Equal(t, "Bug", form.Get("name"))
				return mock.NewResponse(fixture.Category.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			name:                   "Bug",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			name:                   "Bug",
			wantValidationErrCount: 1,
		},

		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			name:                   "",
			wantValidationErrCount: 1,
		},

		"error-validation-all": {
			projectIDOrKey:         "",
			name:                   "",
			wantValidationErrCount: 2,
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "Bug",
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
			s := project.NewCategoryService(method)
			category, err := s.Create(context.Background(), tc.projectIDOrKey, tc.name)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, category)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, category)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, category)
			assert.Equal(t, 12, category.ID)
			assert.Equal(t, "Bug", category.Name)
		})
	}
}

func TestCategoryService_Update(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		categoryID     int
		name           string

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "Bug Fixed",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories/12", spath)
				assert.Equal(t, "Bug Fixed", form.Get("name"))
				return mock.NewResponse(fixture.Category.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			categoryID:             12,
			name:                   "Bug Fixed",
			wantValidationErrCount: 1,
		},
		"error-validation-categoryID-zero": {
			projectIDOrKey:         "TEST",
			categoryID:             0,
			name:                   "Bug Fixed",
			wantValidationErrCount: 1,
		},

		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			categoryID:             12,
			name:                   "",
			wantValidationErrCount: 1,
		},

		"error-validation-all": {
			projectIDOrKey:         "",
			categoryID:             0,
			name:                   "",
			wantValidationErrCount: 3,
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "Bug Fixed",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "Bug Fixed",
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
			s := project.NewCategoryService(method)
			category, err := s.Update(context.Background(), tc.projectIDOrKey, tc.categoryID, tc.name)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, category)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, category)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, category)
			assert.Equal(t, 12, category.ID)
		})
	}
}

func TestCategoryService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		categoryID     int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories/12", spath)
				return mock.NewResponse(fixture.Category.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			categoryID:             12,
			wantValidationErrCount: 1,
		},
		"error-validation-categoryID-zero": {
			projectIDOrKey:         "TEST",
			categoryID:             0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			categoryID:             0,
			wantValidationErrCount: 2,
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			categoryID:     12,
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
			s := project.NewCategoryService(method)
			category, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.categoryID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, category)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, category)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, category)
			assert.Equal(t, 12, category.ID)
			assert.Equal(t, "Bug", category.Name)
		})
	}
}
