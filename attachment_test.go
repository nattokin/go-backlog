package backlog

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
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func newTestAttachment() *Attachment {
	return &Attachment{
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

func newTestAttachmentList() []*Attachment {
	return []*Attachment{
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

func newTestAttachmentSingleList() []*Attachment {
	return []*Attachment{
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
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataAttachmentUploadJSON))),
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
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := attachment.NewSpaceAttachmentService(&core.Method{
				Upload: tc.mockUploadFn,
			})

			f, err := os.Open("testdata/testfile")
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

func TestWikiAttachmentService_Attach(t *testing.T) {
	cases := map[string]struct {
		wikiID        int
		attachmentIDs []int

		expectError bool
		want        []*Attachment

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success-single": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			want:          newTestAttachmentSingleList(),
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments", spath)

				v := form
				assert.Equal(t, []string{"2"}, v["attachmentId[]"])

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataAttachmentSingleListJSON)),
					),
				}, nil
			},
		},

		"success-multiple": {
			wikiID:        1,
			attachmentIDs: []int{2, 5},
			want:          newTestAttachmentList(),
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataAttachmentListJSON)),
					),
				}, nil
			},
		},

		"error-wikiID-invalid": {
			wikiID:        0,
			attachmentIDs: []int{1, 2},
			expectError:   true,
			mockPostFn:    mock.NewUnexpectedPostFn(t),
		},

		"error-attachmentIDs-invalid": {
			wikiID:        1,
			attachmentIDs: []int{0, 1, 2},
			expectError:   true,
			mockPostFn:    mock.NewUnexpectedPostFn(t),
		},

		"error-attachmentIDs-empty": {
			wikiID:        1,
			attachmentIDs: []int{},
			expectError:   true,
			mockPostFn:    mock.NewUnexpectedPostFn(t),
		},

		"error-client": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			expectError:   true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			expectError:   true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataInvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := attachment.NewWikiAttachmentService(&core.Method{Post: tc.mockPostFn})

			attachments, err := s.Attach(context.Background(), tc.wikiID, tc.attachmentIDs)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachments)

			assert.Len(t, attachments, len(tc.want))

			for i, w := range tc.want {
				assert.Equal(t, w.ID, attachments[i].ID)
				assert.Equal(t, w.Name, attachments[i].Name)
				assert.Equal(t, w.Size, attachments[i].Size)
				assert.ObjectsAreEqualValues(w.Created, attachments[i].Created)
			}
		})
	}
}

func TestWikiAttachmentService_List(t *testing.T) {
	cases := map[string]struct {
		wikiID int

		expectError bool
		want        []*Attachment

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1234,
			want:   newTestAttachmentList(),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataAttachmentListJSON)),
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
						bytes.NewReader([]byte(testdataInvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := attachment.NewWikiAttachmentService(&core.Method{Get: tc.mockGetFn})

			attachments, err := s.List(context.Background(), tc.wikiID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachments)

			assert.Len(t, attachments, len(tc.want))

			for i, w := range tc.want {
				assert.Equal(t, w.ID, attachments[i].ID)
				assert.Equal(t, w.Name, attachments[i].Name)
				assert.Equal(t, w.Size, attachments[i].Size)
				assert.ObjectsAreEqualValues(w.Created, attachments[i].Created)
			}
		})
	}
}

func TestWikiAttachmentService_Remove(t *testing.T) {
	cases := map[string]struct {
		wikiID       int
		attachmentID int

		expectError bool
		want        *Attachment

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID:       1234,
			attachmentID: 8,
			want:         newTestAttachment(),
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments/8", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataAttachmentJSON)),
					),
				}, nil
			},
		},

		"error-wikiID-zero": {
			wikiID:       0,
			attachmentID: 8,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-wikiID-negative": {
			wikiID:       -1,
			attachmentID: 8,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-attachmentID-zero": {
			wikiID:       1,
			attachmentID: 0,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-attachmentID-negative": {
			wikiID:       1,
			attachmentID: -1,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-client": {
			wikiID:       1234,
			attachmentID: 8,
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:       1234,
			attachmentID: 8,
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataInvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := attachment.NewWikiAttachmentService(&core.Method{Delete: tc.mockDeleteFn})

			attachment, err := s.Remove(context.Background(), tc.wikiID, tc.attachmentID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachment)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachment)

			assert.Equal(t, tc.want.ID, attachment.ID)
			assert.Equal(t, tc.want.Name, attachment.Name)
			assert.Equal(t, tc.want.Size, attachment.Size)
			assert.ObjectsAreEqualValues(tc.want.Created, attachment.Created)
		})
	}
}

func TestIssueAttachmentService_List(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		expectError bool
		want        []*Attachment

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "1234",
			want:         newTestAttachmentList(),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1234/attachments", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataAttachmentListJSON)),
					),
				}, nil
			},
		},

		"error-invalid-issueIDOrKey": {
			issueIDOrKey: "0",
			expectError:  true,
			mockGetFn:    mock.NewUnexpectedGetFn(t),
		},

		"error-client": {
			issueIDOrKey: "1234",
			expectError:  true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			issueIDOrKey: "1234",
			expectError:  true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataInvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newIssueAttachmentService()
			s.method.Get = tc.mockGetFn

			attachments, err := s.List(context.Background(), tc.issueIDOrKey)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachments)

			assert.Len(t, attachments, len(tc.want))

			for i, w := range tc.want {
				assert.Equal(t, w.ID, attachments[i].ID)
				assert.Equal(t, w.Name, attachments[i].Name)
				assert.Equal(t, w.Size, attachments[i].Size)
				assert.Equal(t, w.Created, attachments[i].Created)
			}
		})
	}
}

func TestIssueAttachmentService_Remove(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		attachmentID int

		expectError bool
		want        *Attachment

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "1234",
			attachmentID: 8,
			want:         newTestAttachment(),
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1234/attachments/8", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataAttachmentJSON)),
					),
				}, nil
			},
		},

		"error-empty-issueKey": {
			issueIDOrKey: "",
			attachmentID: 8,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-attachmentID-zero": {
			issueIDOrKey: "test",
			attachmentID: 0,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},

		"error-client": {
			issueIDOrKey: "1234",
			attachmentID: 8,
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			issueIDOrKey: "1234",
			attachmentID: 8,
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataInvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newIssueAttachmentService()
			s.method.Delete = tc.mockDeleteFn

			attachment, err := s.Remove(context.Background(), tc.issueIDOrKey, tc.attachmentID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachment)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachment)

			assert.Equal(t, tc.want.ID, attachment.ID)
			assert.Equal(t, tc.want.Name, attachment.Name)
			assert.Equal(t, tc.want.Size, attachment.Size)
			assert.Equal(t, tc.want.Created, attachment.Created)
		})
	}
}

func TestPullRequestAttachmentService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey     string
		repositoryIDOrName string
		prNumber           int

		expectError bool
		want        []*Attachment

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "test",
			prNumber:           1234,
			want:               newTestAttachmentList(),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(
					t,
					"projects/TEST/git/repositories/test/pullRequests/1234/attachments",
					spath,
				)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataAttachmentListJSON)),
					),
				}, nil
			},
		},

		"error-invalid-project": {
			projectIDOrKey:     "0",
			repositoryIDOrName: "1",
			prNumber:           1,
			expectError:        true,
			mockGetFn:          mock.NewUnexpectedGetFn(t),
		},

		"error-invalid-repository": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "0",
			prNumber:           1,
			expectError:        true,
			mockGetFn:          mock.NewUnexpectedGetFn(t),
		},

		"error-invalid-prNumber": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "1",
			prNumber:           0,
			expectError:        true,
			mockGetFn:          mock.NewUnexpectedGetFn(t),
		},

		"error-client": {
			projectIDOrKey:     "1234",
			repositoryIDOrName: "test",
			prNumber:           10,
			expectError:        true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			projectIDOrKey:     "1234",
			repositoryIDOrName: "test",
			prNumber:           10,
			expectError:        true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataInvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newPullRequestAttachmentService()
			s.method.Get = tc.mockGetFn

			attachments, err := s.List(context.Background(),
				tc.projectIDOrKey,
				tc.repositoryIDOrName,
				tc.prNumber,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachments)

			assert.Len(t, attachments, len(tc.want))

			for i, w := range tc.want {
				assert.Equal(t, w.ID, attachments[i].ID)
				assert.Equal(t, w.Name, attachments[i].Name)
				assert.Equal(t, w.Size, attachments[i].Size)
				assert.Equal(t, w.Created, attachments[i].Created)
			}
		})
	}
}

func TestPullRequestAttachmentService_Remove(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey     string
		repositoryIDOrName string
		prNumber           int
		attachmentID       int

		expectError bool
		want        *Attachment

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "test",
			prNumber:           1234,
			attachmentID:       8,
			want:               newTestAttachment(),
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(
					t,
					"projects/TEST/git/repositories/test/pullRequests/1234/attachments/8",
					spath,
				)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataAttachmentJSON)),
					),
				}, nil
			},
		},

		"error-invalid-project": {
			projectIDOrKey:     "0",
			repositoryIDOrName: "test",
			prNumber:           1,
			attachmentID:       8,
			expectError:        true,
			mockDeleteFn:       mock.NewUnexpectedDeleteFn(t),
		},

		"error-invalid-repository": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "",
			prNumber:           1,
			attachmentID:       8,
			expectError:        true,
			mockDeleteFn:       mock.NewUnexpectedDeleteFn(t),
		},

		"error-invalid-prNumber": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "test",
			prNumber:           0,
			attachmentID:       8,
			expectError:        true,
			mockDeleteFn:       mock.NewUnexpectedDeleteFn(t),
		},

		"error-invalid-attachmentID": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "test",
			prNumber:           1,
			attachmentID:       0,
			expectError:        true,
			mockDeleteFn:       mock.NewUnexpectedDeleteFn(t),
		},

		"error-client": {
			projectIDOrKey:     "1234",
			repositoryIDOrName: "test",
			prNumber:           10,
			attachmentID:       8,
			expectError:        true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			projectIDOrKey:     "1234",
			repositoryIDOrName: "test",
			prNumber:           10,
			attachmentID:       8,
			expectError:        true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(testdataInvalidJSON)),
					),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newPullRequestAttachmentService()
			s.method.Delete = tc.mockDeleteFn

			attachment, err := s.Remove(
				context.Background(),
				tc.projectIDOrKey,
				tc.repositoryIDOrName,
				tc.prNumber,
				tc.attachmentID,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachment)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachment)

			assert.Equal(t, tc.want.ID, attachment.ID)
			assert.Equal(t, tc.want.Name, attachment.Name)
			assert.Equal(t, tc.want.Size, attachment.Size)
			assert.Equal(t, tc.want.Created, attachment.Created)
		})
	}
}

// TestAttachmentService_contextPropagation verifies that the context passed to
// each attachment service method is correctly relayed to the underlying method call.
// A sentinel value is embedded in the context and its pointer identity is
// asserted inside the mock to catch any ctx substitution (e.g. context.Background()).
func TestAttachmentService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	cases := []struct {
		name string
		call func(t *testing.T)
	}{
		{"SpaceAttachmentService.Upload", func(t *testing.T) {
			s := attachment.NewSpaceAttachmentService(&core.Method{
				Upload: func(got context.Context, _, _ string, _ io.Reader) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.Upload(ctx, "f", bytes.NewReader(nil)) //nolint:errcheck
		}},
		{"WikiAttachmentService.Attach", func(t *testing.T) {
			s := attachment.NewWikiAttachmentService(&core.Method{
				Post: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.Attach(ctx, 1, []int{1}) //nolint:errcheck
		}},
		{"WikiAttachmentService.List", func(t *testing.T) {
			s := attachment.NewWikiAttachmentService(&core.Method{
				Get: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"WikiAttachmentService.Remove", func(t *testing.T) {
			s := attachment.NewWikiAttachmentService(&core.Method{
				Delete: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.Remove(ctx, 1, 1) //nolint:errcheck
		}},
		{"IssueAttachmentService.List", func(t *testing.T) {
			s := newIssueAttachmentService()
			s.method.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s.List(ctx, "TEST-1") //nolint:errcheck
		}},
		{"IssueAttachmentService.Remove", func(t *testing.T) {
			s := newIssueAttachmentService()
			s.method.Delete = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s.Remove(ctx, "TEST-1", 1) //nolint:errcheck
		}},
		{"PullRequestAttachmentService.List", func(t *testing.T) {
			s := newPullRequestAttachmentService()
			s.method.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s.List(ctx, "TEST", "repo", 1) //nolint:errcheck
		}},
		{"PullRequestAttachmentService.Remove", func(t *testing.T) {
			s := newPullRequestAttachmentService()
			s.method.Delete = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s.Remove(ctx, "TEST", "repo", 1, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t)
		})
	}
}
