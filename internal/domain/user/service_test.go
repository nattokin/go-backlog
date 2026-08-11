package user_test

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
	"github.com/nattokin/go-backlog/internal/domain/user"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestUserService_List(t *testing.T) {
	cases := map[string]struct {
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantLen     int
		wantErrType error
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				return mock.NewResponse(fixture.User.ListJSON), nil
			},
			wantLen: 4,
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
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
			s := user.NewService(method)
			users, err := s.List(context.Background())

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, users)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.Len(t, users, tc.wantLen)
			require.NotNil(t, users[0])
			assert.Equal(t, "admin", users[0].UserID)
			assert.Equal(t, "admin", users[0].Name)
			assert.Equal(t, "eguchi@nulab.example", users[0].MailAddress)
			assert.Equal(t, 1, users[0].RoleType)
		})
	}
}

func TestUserService_Me(t *testing.T) {
	cases := map[string]struct {
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErrType error
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself", spath)
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
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
			s := user.NewService(method)
			got, err := s.Me(context.Background())

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, "admin", got.UserID)
			assert.Equal(t, "admin", got.Name)
			assert.Equal(t, "eguchi@nulab.example", got.MailAddress)
			assert.Equal(t, 1, got.RoleType)
		})
	}
}

func TestUserService_Icon(t *testing.T) {
	cases := map[string]struct {
		id int

		mockDownloadFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantFilename           string
		wantContentType        string
	}{
		"success-id-1": {
			id: 1,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/icon", spath)
				return mock.NewBinaryResponse("avatar.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "avatar.png",
			wantContentType: "image/png",
		},
		"success-id-100": {
			id: 100,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/100/icon", spath)
				return mock.NewBinaryResponse("avatar.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "avatar.png",
			wantContentType: "image/png",
		},

		"error-validation-id-zero": {
			id:                     0,
			wantValidationErrCount: 1,
		},
		"error-validation-id-negative": {
			id:                     -1,
			wantValidationErrCount: 1,
		},

		"error-client-network": {
			id: 1,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDownloadFn != nil {
				method.Download = tc.mockDownloadFn
			}
			s := user.NewService(method)
			got, err := s.Icon(context.Background(), tc.id)

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
			assert.Equal(t, tc.wantFilename, got.Filename)
			assert.Equal(t, tc.wantContentType, got.ContentType)
			require.NotNil(t, got.Body)
			got.Body.Close()
		})
	}
}
