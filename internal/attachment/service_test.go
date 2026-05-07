package attachment_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_Upload(t *testing.T) {
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
			s := attachment.NewService(method)

			f, err := os.Open("../../testdata/testfile")
			require.NoError(t, err)
			defer f.Close()

			attachment, err := s.Upload(context.Background(), "space/attachment", tc.fpath, f)

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

func TestService_List(t *testing.T) {
	cases := map[string]struct {
		spath       string
		expectError bool
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			spath: "issues/TEST-1/attachments",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/attachments", spath)
				return mock.NewJSONResponse(fixture.Attachment.ListJSON), nil
			},
		},
		"error-client": {
			spath:       "issues/TEST-1/attachments",
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			spath:       "issues/TEST-1/attachments",
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
			s := attachment.NewService(method)

			v, err := s.List(context.Background(), tc.spath)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
		})
	}
}

func TestService_Remove(t *testing.T) {
	cases := map[string]struct {
		spath        string
		expectError  bool
		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			spath: "issues/TEST-1/attachments/10",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/attachments/10", spath)
				return mock.NewJSONResponse(fixture.Attachment.SingleJSON), nil
			},
		},
		"error-client": {
			spath:       "issues/TEST-1/attachments/10",
			expectError: true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			spath:       "issues/TEST-1/attachments/10",
			expectError: true,
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
			s := attachment.NewService(method)

			v, err := s.Remove(context.Background(), tc.spath)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
		})
	}
}

func TestService_Download(t *testing.T) {
	cases := map[string]struct {
		spath           string
		expectError     bool
		wantFilename    string
		wantContentType string
		mockDownloadFn  func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			spath: "issues/TEST-1/attachments/10",
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/attachments/10", spath)
				return mock.NewBinaryResponse("file.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "file.png",
			wantContentType: "image/png",
		},
		"error-client": {
			spath:       "issues/TEST-1/attachments/10",
			expectError: true,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDownloadFn != nil {
				method.Download = tc.mockDownloadFn
			}
			s := attachment.NewService(method)

			got, err := s.Download(context.Background(), tc.spath)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, got)
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

func TestService_Attach(t *testing.T) {
	cases := map[string]struct {
		spath         string
		attachmentIDs []int
		expectError   bool
		mockPostFn    func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success-single": {
			spath:         "wikis/1234/attachments",
			attachmentIDs: []int{2},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments", spath)
				assert.Equal(t, []string{"2"}, form["attachmentId[]"])
				return mock.NewJSONResponse(fixture.Attachment.SingleListJSON), nil
			},
		},
		"success-multiple": {
			spath:         "wikis/1/attachments",
			attachmentIDs: []int{2, 5},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.Attachment.ListJSON), nil
			},
		},
		"error-empty-ids": {
			spath:         "wikis/1/attachments",
			attachmentIDs: []int{},
			expectError:   true,
			mockPostFn:    mock.NewUnexpectedPostFn(t),
		},
		"error-zero-id": {
			spath:         "wikis/1/attachments",
			attachmentIDs: []int{0},
			expectError:   true,
			mockPostFn:    mock.NewUnexpectedPostFn(t),
		},
		"error-client": {
			spath:         "wikis/1/attachments",
			attachmentIDs: []int{1},
			expectError:   true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			spath:         "wikis/1/attachments",
			attachmentIDs: []int{1},
			expectError:   true,
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
			s := attachment.NewService(method)

			v, err := s.Attach(context.Background(), tc.spath, tc.attachmentIDs)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
		})
	}
}

func Test_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeMockGetFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"Service.Upload", func(t *testing.T, m *core.Method) {
			m.Upload = func(got context.Context, _, _ string, _ io.Reader) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := attachment.NewService(m)
			s.Upload(ctx, "space/attachment", "f", bytes.NewReader(nil)) //nolint:errcheck
		}},
		{"Service.Attach", func(t *testing.T, m *core.Method) {
			m.Post = makeMockGetFn(t)
			s := attachment.NewService(m)
			s.Attach(ctx, "wikis/1/attachments", []int{1}) //nolint:errcheck
		}},
		{"Service.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockGetFn(t)
			s := attachment.NewService(m)
			s.List(ctx, "wikis/1/attachments") //nolint:errcheck
		}},
		{"Service.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockGetFn(t)
			s := attachment.NewService(m)
			s.Remove(ctx, "wikis/1/attachments/1") //nolint:errcheck
		}},
		{"Service.Download", func(t *testing.T, m *core.Method) {
			m.Download = makeMockGetFn(t)
			s := attachment.NewService(m)
			s.Download(ctx, "wikis/1/attachments/1") //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
