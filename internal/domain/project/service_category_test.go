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

		wantLen     int
		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.Category.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "6",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/categories", spath)
				return mock.NewJSONResponse(fixture.Category.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
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
			s := project.NewCategoryService(method)

			categories, err := s.List(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
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

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "Bug",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories", spath)
				assert.Equal(t, "Bug", form.Get("name"))
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			name:           "Bug",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			name:           "",
			wantErrType:    &core.ValidationError{},
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
			s := project.NewCategoryService(method)

			category, err := s.Create(context.Background(), tc.projectIDOrKey, tc.name)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
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

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "Bug Fixed",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories/12", spath)
				assert.Equal(t, "Bug Fixed", form.Get("name"))
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			categoryID:     12,
			name:           "Bug Fixed",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-categoryID-zero": {
			projectIDOrKey: "TEST",
			categoryID:     0,
			name:           "Bug Fixed",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "",
			wantErrType:    &core.ValidationError{},
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
			s := project.NewCategoryService(method)

			category, err := s.Update(context.Background(), tc.projectIDOrKey, tc.categoryID, tc.name)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
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

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			categoryID:     12,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories/12", spath)
				assert.NotNil(t, form)
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			categoryID:     12,
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-categoryID-zero": {
			projectIDOrKey: "TEST",
			categoryID:     0,
			wantErrType:    &core.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
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
