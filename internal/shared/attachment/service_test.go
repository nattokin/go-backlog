package attachment_test

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
	"github.com/nattokin/go-backlog/internal/shared/attachment"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_List(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantIDs     []int
		wantErrType error
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/attachments", spath)
				assert.Nil(t, query)
				return mock.NewResponse(fixture.Attachment.ListJSON), nil
			},
			wantIDs: []int{2, 5},
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
			s := attachment.NewService(m)
			got, err := s.List(context.Background(), "issues/TEST-1/attachments")

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, len(tc.wantIDs))
			for i, id := range tc.wantIDs {
				assert.Equal(t, id, got[i].ID)
			}
		})
	}
}

func TestService_Remove(t *testing.T) {
	cases := map[string]struct {
		mockDeleteFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID      int
		wantErrType error
	}{
		"success": {
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/attachments/8", spath)
				assert.Nil(t, form)
				return mock.NewResponse(fixture.Attachment.SingleJSON), nil
			},
			wantID: 8,
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
			s := attachment.NewService(m)
			got, err := s.Remove(context.Background(), "issues/TEST-1/attachments/8")

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

func TestService_Download(t *testing.T) {
	cases := map[string]struct {
		mockDownloadFn func(context.Context, string, url.Values) (*http.Response, error)

		wantFilename    string
		wantContentType string
		wantErrType     error
	}{
		"success": {
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/attachments/8", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("file.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "file.png",
			wantContentType: "image/png",
		},
		"error-client-network": {
			mockDownloadFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockDownloadFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Download = tc.mockDownloadFn
			s := attachment.NewService(m)
			got, err := s.Download(context.Background(), "issues/TEST-1/attachments/8")

			if tc.wantErrType != nil {
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
