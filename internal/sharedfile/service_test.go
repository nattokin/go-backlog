package sharedfile_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/sharedfile"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

// ──────────────────────────────────────────────────────────────
//  IssueService tests
// ──────────────────────────────────────────────────────────────

func TestIssueSharedFileService_List(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		expectError bool

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "TEST-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/sharedFiles", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.SharedFile.ListJSON)),
					),
				}, nil
			},
		},

		"success-numeric-id": {
			issueIDOrKey: "1234",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1234/sharedFiles", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.SharedFile.ListJSON)),
					),
				}, nil
			},
		},

		"error-issueIDOrKey-empty": {
			issueIDOrKey: "",
			expectError:  true,
			mockGetFn:    mock.NewUnexpectedGetFn(t),
		},

		"error-issueIDOrKey-zero": {
			issueIDOrKey: "0",
			expectError:  true,
			mockGetFn:    mock.NewUnexpectedGetFn(t),
		},

		"error-client": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
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
			s := sharedfile.NewIssueService(method)

			files, err := s.List(context.Background(), tc.issueIDOrKey)

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
				assert.Equal(t, w.Type, files[i].Type)
				assert.Equal(t, w.Dir, files[i].Dir)
				assert.Equal(t, w.Name, files[i].Name)
				assert.Equal(t, w.Size, files[i].Size)
			}
		})
	}
}

func TestIssueSharedFileService_Link(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		fileIDs      []int

		expectError bool

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success-single": {
			issueIDOrKey: "TEST-1",
			fileIDs:      []int{454403},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/sharedFiles", spath)
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
			issueIDOrKey: "TEST-1",
			fileIDs:      []int{454403, 454404},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.SharedFile.ListJSON)),
					),
				}, nil
			},
		},

		"error-issueIDOrKey-empty": {
			issueIDOrKey: "",
			fileIDs:      []int{1},
			expectError:  true,
			mockPostFn:   mock.NewUnexpectedPostFn(t),
		},

		"error-fileIDs-empty": {
			issueIDOrKey: "TEST-1",
			fileIDs:      []int{},
			expectError:  true,
			mockPostFn:   mock.NewUnexpectedPostFn(t),
		},

		"error-fileIDs-invalid": {
			issueIDOrKey: "TEST-1",
			fileIDs:      []int{0, 1},
			expectError:  true,
			mockPostFn:   mock.NewUnexpectedPostFn(t),
		},

		"error-client": {
			issueIDOrKey: "TEST-1",
			fileIDs:      []int{454403},
			expectError:  true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			issueIDOrKey: "TEST-1",
			fileIDs:      []int{454403},
			expectError:  true,
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
			s := sharedfile.NewIssueService(method)

			files, err := s.Link(context.Background(), tc.issueIDOrKey, tc.fileIDs)

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

func TestIssueSharedFileService_Unlink(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		fileID       int

		expectError bool

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "TEST-1",
			fileID:       454403,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/sharedFiles/454403", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.SharedFile.SingleJSON)),
					),
				}, nil
			},
		},

		"error-issueIDOrKey-empty": {
			issueIDOrKey: "",
			fileID:       454403,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-issueIDOrKey-zero": {
			issueIDOrKey: "0",
			fileID:       454403,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-fileID-zero": {
			issueIDOrKey: "TEST-1",
			fileID:       0,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-fileID-negative": {
			issueIDOrKey: "TEST-1",
			fileID:       -1,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-client": {
			issueIDOrKey: "TEST-1",
			fileID:       454403,
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			issueIDOrKey: "TEST-1",
			fileID:       454403,
			expectError:  true,
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
			s := sharedfile.NewIssueService(method)

			file, err := s.Unlink(context.Background(), tc.issueIDOrKey, tc.fileID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, file)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, file)

			assert.Equal(t, fixture.SharedFile.Single.ID, file.ID)
			assert.Equal(t, fixture.SharedFile.Single.Name, file.Name)
			assert.Equal(t, fixture.SharedFile.Single.Type, file.Type)
			assert.Equal(t, fixture.SharedFile.Single.Dir, file.Dir)
		})
	}
}

// ──────────────────────────────────────────────────────────────
//  WikiService tests
// ──────────────────────────────────────────────────────────────

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
			assert.Len(t, files, len(fixture.SharedFile.List))

			for i, w := range fixture.SharedFile.List {
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

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success-single": {
			wikiID:  1234,
			fileIDs: []int{454403},
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

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1234,
			fileID: 454403,
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
			wikiID:      1234,
			fileID:      454403,
			expectError: true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:      1234,
			fileID:      454403,
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

			assert.Equal(t, fixture.SharedFile.Single.ID, file.ID)
			assert.Equal(t, fixture.SharedFile.Single.Name, file.Name)
			assert.Equal(t, fixture.SharedFile.Single.Type, file.Type)
			assert.Equal(t, fixture.SharedFile.Single.Dir, file.Dir)
		})
	}
}

func Test_contextPropagation(t *testing.T) {
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
		{"IssueService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := sharedfile.NewIssueService(m)
			s.List(ctx, "TEST-1") //nolint:errcheck
		}},
		{"IssueService.Link", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := sharedfile.NewIssueService(m)
			s.Link(ctx, "TEST-1", []int{1}) //nolint:errcheck
		}},
		{"IssueService.Unlink", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := sharedfile.NewIssueService(m)
			s.Unlink(ctx, "TEST-1", 1) //nolint:errcheck
		}},
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
