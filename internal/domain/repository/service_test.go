package repository_test

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
	"github.com/nattokin/go-backlog/internal/domain/repository"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

const (
	testProject = "PRJ"
	testRepo    = "repo1"
)

func TestService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantIDs                []int
	}{
		"success": {
			projectIDOrKey: testProject,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories", spath)
				return mock.NewResponse(fixture.Repository.ListJSON), nil
			},
			wantIDs: []int{5, 6},
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
			projectIDOrKey: testProject,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: testProject,
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
			s := repository.NewService(method)
			got, err := s.List(context.Background(), tc.projectIDOrKey)

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
			assert.Len(t, got, len(tc.wantIDs))
			for i := range got {
				assert.Equal(t, tc.wantIDs[i], got[i].ID)
			}
		})
	}
}

func TestService_One(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantID                 int
		wantName               string
	}{
		"success": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1", spath)
				return mock.NewResponse(fixture.Repository.SingleJSON), nil
			},
			wantID:   5,
			wantName: "foo",
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           testRepo,
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			repoIDOrName:           testRepo,
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         testProject,
			repoIDOrName:           "",
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-zero": {
			projectIDOrKey:         testProject,
			repoIDOrName:           "0",
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			wantValidationErrCount: 2,
		},

		"error-client-network": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
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
			s := repository.NewService(method)
			got, err := s.One(context.Background(), tc.projectIDOrKey, tc.repoIDOrName)

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
			assert.Equal(t, tc.wantID, got.ID)
			assert.Equal(t, tc.wantName, got.Name)
		})
	}
}

func TestService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeMockFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := repository.NewService(m)
			s.List(ctx, testProject) //nolint:errcheck
		}},
		{"One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := repository.NewService(m)
			s.One(ctx, testProject, testRepo) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
