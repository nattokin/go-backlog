package wiki_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/wiki"
)

func TestWikiSharedFileService_List(t *testing.T) {
	cases := map[string]struct {
		wikiID int

		expectError bool

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/sharedFiles", spath)
				return mock.NewJSONResponse(fixture.SharedFile.ListJSON), nil
			},
		},

		"error-wikiID-zero": {
			wikiID:      0,
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := wiki.NewSharedFileService(method)

			files, err := s.List(context.Background(), tc.wikiID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, files)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, files)
			assert.Len(t, files, len(fixture.SharedFile.List))

			for i, w := range fixture.SharedFile.List {
				assert.Equal(t, w.ID, files[i].ID)
			}
		})
	}
}

func TestWikiSharedFileService_Link(t *testing.T) {
	cases := map[string]struct {
		wikiID  int
		fileIDs []int

		expectError bool

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success-single": {
			wikiID:  1234,
			fileIDs: []int{454403},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/sharedFiles", spath)
				assert.Equal(t, []string{"454403"}, form["fileId[]"])
				return mock.NewJSONResponse(fixture.SharedFile.SingleListJSON), nil
			},
		},

		"success-multiple": {
			wikiID:  1,
			fileIDs: []int{454403, 454404},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.SharedFile.ListJSON), nil
			},
		},

		"error-wikiID-zero": {
			wikiID:      0,
			fileIDs:     []int{1},
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},

		"error-fileIDs-empty": {
			wikiID:      1234,
			fileIDs:     []int{},
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},

		"error-fileIDs-invalid": {
			wikiID:      1234,
			fileIDs:     []int{0, 1},
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},

		"error-client": {
			wikiID:      1234,
			fileIDs:     []int{454403},
			expectError: true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:      1234,
			fileIDs:     []int{454403},
			expectError: true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Post = tc.mockPostFn
			s := wiki.NewSharedFileService(method)

			files, err := s.Link(context.Background(), tc.wikiID, tc.fileIDs)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, files)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, files)
			assert.NotEmpty(t, files)

			for _, f := range files {
				assert.Positive(t, f.ID)
				assert.NotEmpty(t, f.Name)
			}
		})
	}
}

func TestWikiSharedFileService_Unlink(t *testing.T) {
	cases := map[string]struct {
		wikiID int
		fileID int

		expectError bool
		wantID      int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1234,
			fileID: 454403,
			wantID: fixture.SharedFile.Single.ID,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/sharedFiles/454403", spath)
				return mock.NewJSONResponse(fixture.SharedFile.SingleJSON), nil
			},
		},

		"error-wikiID-zero": {
			wikiID:       0,
			fileID:       454403,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-fileID-zero": {
			wikiID:       1234,
			fileID:       0,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-fileID-negative": {
			wikiID:       1234,
			fileID:       -1,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-client": {
			wikiID:       1234,
			fileID:       454403,
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:       1234,
			fileID:       454403,
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Delete = tc.mockDeleteFn
			s := wiki.NewSharedFileService(method)

			file, err := s.Unlink(context.Background(), tc.wikiID, tc.fileID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, file)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, file)
			assert.Equal(t, tc.wantID, file.ID)
		})
	}
}
