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
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestProjectSharedFileService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-project-key": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/files", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.SharedFile.ListJSON), nil
			},
		},
		"success-project-id": {
			projectIDOrKey: "1234",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234/files", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.SharedFile.ListJSON), nil
			},
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey: "0",
			wantErrType:    &core.ValidationError{},
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
			s := project.NewSharedFileService(method)

			files, err := s.List(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, files)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, files)
			assert.Len(t, files, len(fixture.SharedFile.List))

			for i, v := range fixture.SharedFile.List {
				assert.Equal(t, v.ID, files[i].ID)
				assert.Equal(t, v.Type, files[i].Type)
				assert.Equal(t, v.Dir, files[i].Dir)
				assert.Equal(t, v.Name, files[i].Name)
				assert.Equal(t, v.Size, files[i].Size)
			}
		})
	}
}

func TestProjectSharedFileService_GetFile(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		sharedFileID   int

		mockDownloadFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType     error
		wantFilename    string
		wantContentType string
	}{
		"success-project-key": {
			projectIDOrKey: "TEST",
			sharedFileID:   454403,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/files/454403", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("01_buz.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "01_buz.png",
			wantContentType: "image/png",
		},
		"success-project-id": {
			projectIDOrKey: "1234",
			sharedFileID:   454404,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234/files/454404", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("readme.md", "text/plain", []byte("MD")), nil
			},
			wantFilename:    "readme.md",
			wantContentType: "text/plain",
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			sharedFileID:   454403,
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey: "0",
			sharedFileID:   454403,
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-sharedFileID-zero": {
			projectIDOrKey: "TEST",
			sharedFileID:   0,
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			sharedFileID:   454403,
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
			s := project.NewSharedFileService(method)

			got, err := s.GetFile(context.Background(), tc.projectIDOrKey, tc.sharedFileID)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
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
