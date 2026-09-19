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

	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/validation"
)

func TestSharedFileService_List(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		dirPath        string
		options        []*option.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success-project-key": {
			projectIDOrKey: "TEST",
			dirPath:        "/design/",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/files/metadata/design", spath)
				return mock.NewResponse(fixture.SharedFile.ListJSON), nil
			},
		},
		"success-project-id": {
			projectIDOrKey: "1234",
			dirPath:        "/design/",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234/files/metadata/design", spath)
				return mock.NewResponse(fixture.SharedFile.ListJSON), nil
			},
		},
		"success-nested-dirPath": {
			projectIDOrKey: "TEST",
			dirPath:        "/PressRelease/20091130/",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/files/metadata/PressRelease/20091130", spath)
				return mock.NewResponse(fixture.SharedFile.ListJSON), nil
			},
		},
		"success-root-dirPath": {
			projectIDOrKey: "TEST",
			dirPath:        "/",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/files/metadata", spath)
				return mock.NewResponse(fixture.SharedFile.ListJSON), nil
			},
		},
		"success-with-options": {
			projectIDOrKey: "TEST",
			dirPath:        "/design/",
			options: []*option.APIParamOption{
				o.WithOrder("asc"),
				o.WithOffset(10),
				o.WithSharedFileCount(1000),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/files/metadata/design", spath)
				assert.Equal(t, "asc", query.Get("order"))
				assert.Equal(t, "10", query.Get("offset"))
				assert.Equal(t, "1000", query.Get("count"))
				return mock.NewResponse(fixture.SharedFile.ListJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			dirPath:                "/design/",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			dirPath:                "/design/",
			wantValidationErrCount: 1,
		},
		"error-validation-dirPath-empty": {
			projectIDOrKey:         "TEST",
			dirPath:                "",
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			dirPath:                "",
			wantValidationErrCount: 2,
		},
		"error-validation-option": {
			projectIDOrKey:         "TEST",
			dirPath:                "/design/",
			options:                []*option.APIParamOption{o.WithSharedFileCount(1001)},
			wantValidationErrCount: 1,
		},
		"error-validation-invalid-option-type": {
			projectIDOrKey:         "TEST",
			dirPath:                "/design/",
			options:                []*option.APIParamOption{o.WithKeyword("foo")},
			wantErrType:            &option.InvalidOptionError{},
			wantValidationErrCount: 0,
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			dirPath:        "/design/",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			dirPath:        "/design/",
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
			s := project.NewSharedFileService(method)
			files, err := s.List(context.Background(), tc.projectIDOrKey, tc.dirPath, tc.options...)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, files)
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, files)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, files)
			assert.Len(t, files, len(fixture.SharedFile.List))
		})
	}
}

func TestSharedFileService_Download(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		sharedFileID   int

		mockDownloadFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantFilename           string
		wantContentType        string
	}{
		"success-project-key": {
			projectIDOrKey: "TEST",
			sharedFileID:   454403,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/files/454403", spath)
				return mock.NewBinaryResponse("01_buz.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "01_buz.png",
			wantContentType: "image/png",
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			sharedFileID:           454403,
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			sharedFileID:           454403,
			wantValidationErrCount: 1,
		},
		"error-validation-sharedFileID-zero": {
			projectIDOrKey:         "TEST",
			sharedFileID:           0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			sharedFileID:           0,
			wantValidationErrCount: 2,
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
			got, err := s.Download(context.Background(), tc.projectIDOrKey, tc.sharedFileID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves validation.Errors
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
