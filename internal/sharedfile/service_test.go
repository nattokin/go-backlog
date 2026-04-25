package sharedfile_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/sharedfile"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func newTestSharedFile() *model.SharedFile {
	return &model.SharedFile{
		ID:   454403,
		Type: "file",
		Dir:  "/icon/",
		Name: "01_buz.png",
		Size: 2735,
		Created: time.Date(
			2009, time.February, 27, 3, 26, 15, 0, time.UTC,
		),
		Updated: time.Date(
			2009, time.March, 3, 16, 57, 47, 0, time.UTC,
		),
	}
}

func newTestSharedFileList() []*model.SharedFile {
	return []*model.SharedFile{
		{
			ID:   454403,
			Type: "file",
			Dir:  "/icon/",
			Name: "01_buz.png",
			Size: 2735,
			Created: time.Date(
				2009, time.February, 27, 3, 26, 15, 0, time.UTC,
			),
			Updated: time.Date(
				2009, time.March, 3, 16, 57, 47, 0, time.UTC,
			),
		},
		{
			ID:   454404,
			Type: "file",
			Dir:  "/docs/",
			Name: "readme.md",
			Size: 512,
			Created: time.Date(
				2009, time.February, 27, 3, 26, 15, 0, time.UTC,
			),
			Updated: time.Date(
				2009, time.March, 3, 16, 57, 47, 0, time.UTC,
			),
		},
	}
}

func TestWikiSharedFileService_List(t *testing.T) {
	cases := map[string]struct {
		wikiID int

		expectError bool
		want        []*model.SharedFile

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1234,
			want:   newTestSharedFileList(),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/sharedFiles", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.SharedFile.ListJSON)),
					),
				}, nil
			},
		},

		"error-wikiID-zero": {
			wikiID:      0,
			expectError: true,
			mockGetFn:   mock.NewUnexpectedGetFn(t),
		},

		"error-wikiID-negative": {
			wikiID:      -1,
			expectError: true,
			mockGetFn:   mock.NewUnexpectedGetFn(t),
		},

		"error-client": {
			wikiID:      1234,
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:      1234,
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.InvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := sharedfile.NewWikiService(method)

			files, err := s.List(context.Background(), tc.wikiID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, files)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, files)
			assert.Len(t, files, len(tc.want))

			for i, w := range tc.want {
				assert.Equal(t, w.ID, files[i].ID)
				assert.Equal(t, w.Type, files[i].Type)
				assert.Equal(t, w.Dir, files[i].Dir)
				assert.Equal(t, w.Name, files[i].Name)
				assert.Equal(t, w.Size, files[i].Size)
			}
		})
	}
}

func TestWikiSharedFileService_Link(t *testing.T) {
	cases := map[string]struct {
		wikiID  int
		fileIDs []int

		expectError bool
		want        []*model.SharedFile

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success-single": {
			wikiID:  1234,
			fileIDs: []int{454403},
			want:    newTestSharedFileList()[:1],
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/sharedFiles", spath)
				assert.Equal(t, []string{"454403"}, form["fileId[]"])

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.SharedFile.SingleListJSON)),
					),
				}, nil
			},
		},

		"success-multiple": {
			wikiID:  1,
			fileIDs: []int{454403, 454404},
			want:    newTestSharedFileList(),
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.SharedFile.ListJSON)),
					),
				}, nil
			},
		},

		"error-wikiID-zero": {
			wikiID:      0,
			fileIDs:     []int{1},
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},

		"error-fileIDs-empty": {
			wikiID:      1,
			fileIDs:     []int{},
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},

		"error-fileIDs-invalid": {
			wikiID:      1,
			fileIDs:     []int{0, 1},
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},

		"error-client": {
			wikiID:  1234,
			fileIDs: []int{454403},
			expectError: true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:  1234,
			fileIDs: []int{454403},
			expectError: true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.InvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Post = tc.mockPostFn
			s := sharedfile.NewWikiService(method)

			files, err := s.Link(context.Background(), tc.wikiID, tc.fileIDs)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, files)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, files)
			assert.Len(t, files, len(tc.want))

			for i, w := range tc.want {
				assert.Equal(t, w.ID, files[i].ID)
				assert.Equal(t, w.Name, files[i].Name)
			}
		})
	}
}

func TestWikiSharedFileService_Unlink(t *testing.T) {
	cases := map[string]struct {
		wikiID int
		fileID int

		expectError bool
		want        *model.SharedFile

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1234,
			fileID: 454403,
			want:   newTestSharedFile(),
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/sharedFiles/454403", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.SharedFile.SingleJSON)),
					),
				}, nil
			},
		},

		"error-wikiID-zero": {
			wikiID:       0,
			fileID:       454403,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-wikiID-negative": {
			wikiID:       -1,
			fileID:       454403,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-fileID-zero": {
			wikiID:       1,
			fileID:       0,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-fileID-negative": {
			wikiID:       1,
			fileID:       -1,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-client": {
			wikiID: 1234,
			fileID: 454403,
			expectError: true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID: 1234,
			fileID: 454403,
			expectError: true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.InvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Delete = tc.mockDeleteFn
			s := sharedfile.NewWikiService(method)

			file, err := s.Unlink(context.Background(), tc.wikiID, tc.fileID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, file)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, file)

			assert.Equal(t, tc.want.ID, file.ID)
			assert.Equal(t, tc.want.Name, file.Name)
			assert.Equal(t, tc.want.Type, file.Type)
			assert.Equal(t, tc.want.Dir, file.Dir)
		})
	}
}

func TestWikiSharedFileService_contextPropagation(t *testing.T) {
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
		{"WikiService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := sharedfile.NewWikiService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"WikiService.Link", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := sharedfile.NewWikiService(m)
			s.Link(ctx, 1, []int{1}) //nolint:errcheck
		}},
		{"WikiService.Unlink", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := sharedfile.NewWikiService(m)
			s.Unlink(ctx, 1, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
