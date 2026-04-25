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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func newTestAttachment() *model.Attachment {
	return &model.Attachment{
		ID:   8,
		Name: "IMG0088.png",
		Size: 5563,
		Created: time.Date(
			2014,
			time.October,
			28,
			9,
			24,
			43,
			0,
			time.UTC,
		),
	}
}

func newTestAttachmentList() []*model.Attachment {
	return []*model.Attachment{
		{
			ID:   2,
			Name: "A.png",
			Size: 196186,
			Created: time.Date(
				2014,
				time.July,
				11,
				6,
				26,
				5,
				0,
				time.UTC,
			),
		},
		{
			ID:   5,
			Name: "B.png",
			Size: 201257,
			Created: time.Date(
				2014,
				time.July,
				11,
				6,
				26,
				5,
				0,
				time.UTC,
			),
		},
	}
}

func newTestAttachmentSingleList() []*model.Attachment {
	return []*model.Attachment{
		{
			ID:   2,
			Name: "A.png",
			Size: 196186,
			Created: time.Date(
				2014,
				time.September,
				11,
				6,
				26,
				5,
				0,
				time.UTC,
			),
		},
	}
}

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

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.UploadJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Upload = tc.mockUploadFn
			s := attachment.NewSpaceService(method)

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
		{"SpaceService.Upload", func(t *testing.T, m *core.Method) {
			m.Upload = func(got context.Context, _, _ string, _ io.Reader) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := attachment.NewSpaceService(m)
			s.Upload(ctx, "f", bytes.NewReader(nil)) //nolint:errcheck
		}},
		{"WikiService.Attach", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := attachment.NewWikiService(m)
			s.Attach(ctx, 1, []int{1}) //nolint:errcheck
		}},
		{"WikiService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := attachment.NewWikiService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"WikiService.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := attachment.NewWikiService(m)
			s.Remove(ctx, 1, 1) //nolint:errcheck
		}},
		{"IssueService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := attachment.NewIssueService(m)
			s.List(ctx, "TEST-1") //nolint:errcheck
		}},
		{"IssueService.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := attachment.NewIssueService(m)
			s.Remove(ctx, "TEST-1", 1) //nolint:errcheck
		}},
		{"PullRequestService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := attachment.NewPullRequestService(m)
			s.List(ctx, "TEST", "repo", 1) //nolint:errcheck
		}},
		{"PullRequestService.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := attachment.NewPullRequestService(m)
			s.Remove(ctx, "TEST", "repo", 1, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
