package space_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/domain/space"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestSpaceAttachmentService_Upload(t *testing.T) {
	cases := map[string]struct {
		fpath string

		expectError bool

		mockUploadFn func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error)
	}{
		"success": {
			fpath: "fpath",
			mockUploadFn: func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error) {
				assert.Equal(t, "space/attachment", spath)
				assert.Equal(t, "fpath", fileName)
				return mock.NewJSONResponse(fixture.Attachment.UploadJSON), nil
			},
		},

		"error-client-failure": {
			fpath:       "fpath",
			expectError: true,
			mockUploadFn: func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			fpath:       "fpath",
			expectError: true,
			mockUploadFn: func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Upload = tc.mockUploadFn
			s := space.NewAttachmentService(method)

			f, err := os.Open("../../testdata/testfile")
			require.NoError(t, err)
			defer f.Close()

			attachment, err := s.Upload(context.Background(), tc.fpath, f)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachment)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachment)

			assert.Equal(t, 1, attachment.ID)
			assert.Equal(t, "test.txt", attachment.Name)
			assert.Equal(t, 8857, attachment.Size)
		})
	}
}
