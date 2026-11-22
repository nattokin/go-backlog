package sharedfile_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/shared/sharedfile"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_List(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantFirstID int
		wantErrType error
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/sharedFiles", spath)
				assert.Nil(t, query)
				return mock.NewResponse(fixture.SharedFile.ListJSON), nil
			},
			wantFirstID: 454403,
		},
		"error-client-network": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Get = tc.mockGetFn
			s := sharedfile.NewService(m)
			got, err := s.List(context.Background(), "issues/TEST-1/sharedFiles")

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, got)
			assert.Equal(t, tc.wantFirstID, got[0].ID)
		})
	}
}

func TestService_Link(t *testing.T) {
	cases := map[string]struct {
		fileIDs []int

		mockPostFn func(context.Context, string, url.Values) (*http.Response, error)

		wantFirstID int
		wantErrType error
	}{
		"success": {
			fileIDs: []int{454403, 54},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/sharedFiles", spath)
				assert.Equal(t, []string{"454403", "54"}, form["fileId[]"])
				return mock.NewResponse(fixture.SharedFile.ListJSON), nil
			},
			wantFirstID: 454403,
		},
		"error-client-network": {
			fileIDs: []int{1},
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			fileIDs: []int{1},
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			fileIDs: []int{1},
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Post = tc.mockPostFn
			s := sharedfile.NewService(m)
			got, err := s.Link(context.Background(), "issues/TEST-1/sharedFiles", tc.fileIDs)

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, got)
			assert.Equal(t, tc.wantFirstID, got[0].ID)
		})
	}
}

func TestService_Unlink(t *testing.T) {
	cases := map[string]struct {
		mockDeleteFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID      int
		wantErrType error
	}{
		"success": {
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/sharedFiles/454403", spath)
				assert.Nil(t, form)
				return mock.NewResponse(fixture.SharedFile.SingleJSON), nil
			},
			wantID: 454403,
		},
		"error-client-network": {
			mockDeleteFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockDeleteFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockDeleteFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Delete = tc.mockDeleteFn
			s := sharedfile.NewService(m)
			got, err := s.Unlink(context.Background(), "issues/TEST-1/sharedFiles/454403")

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
