package backlog

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		fpath     string
		expectErr bool

		mockUploadFn func(spath, fileName string, r io.Reader) (*http.Response, error)
	}{
		"success": {
			fpath: "fpath",
			mockUploadFn: func(spath, fileName string, r io.Reader) (*http.Response, error) {
				assert.Equal(t, "space/attachment", spath)
				assert.Equal(t, "fpath", fileName)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataAttachmentUploadJSON))),
				}, nil
			},
		},

		"error-client-failure": {
			fpath:     "fpath",
			expectErr: true,
			mockUploadFn: func(spath, fileName string, r io.Reader) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			fpath:     "fpath",
			expectErr: true,
			mockUploadFn: func(spath, fileName string, r io.Reader) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newSpaceAttachmentService()
			s.method.Upload = tc.mockUploadFn

			f, err := os.Open("testdata/testfile")
			require.NoError(t, err)
			defer f.Close()

			attachment, err := s.Upload(tc.fpath, f)

			if tc.expectErr {
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

		expectErr bool
		want      []*Attachment

		mockPostFn func(spath string, form *FormParams) (*http.Response, error)
	}{
		"success-single": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			want:          newTestAttachmentSingleList(),
			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments", spath)

				v := *form.Values
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
			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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
			expectErr:     true,
			mockPostFn:    newUnexpectedPostFn(t, "invalid wikiID"),
		},

		"error-attachmentIDs-invalid": {
			wikiID:        1,
			attachmentIDs: []int{0, 1, 2},
			expectErr:     true,
			mockPostFn:    newUnexpectedPostFn(t, "invalid attachmentIDs"),
		},

		"error-attachmentIDs-empty": {
			wikiID:        1,
			attachmentIDs: []int{},
			expectErr:     true,
			mockPostFn:    newUnexpectedPostFn(t, "empty attachmentIDs"),
		},

		"error-client": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			expectErr:     true,
			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			expectErr:     true,
			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiAttachmentService()
			s.method.Post = tc.mockPostFn

			attachments, err := s.Attach(tc.wikiID, tc.attachmentIDs)

			if tc.expectErr {
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

		expectErr bool
		want      []*Attachment

		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)
	}{
		"success": {
			wikiID: 1234,
			want:   newTestAttachmentList(),
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
			wikiID:    0,
			expectErr: true,
			mockGetFn: newUnexpectedGetFn(t, "invalid wikiID"),
		},

		"error-wikiID-negative": {
			wikiID:    -1,
			expectErr: true,
			mockGetFn: newUnexpectedGetFn(t, "invalid wikiID"),
		},

		"error-client": {
			wikiID:    1234,
			expectErr: true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:    1234,
			expectErr: true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiAttachmentService()
			s.method.Get = tc.mockGetFn

			attachments, err := s.List(tc.wikiID)

			if tc.expectErr {
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

		expectErr bool
		want      *Attachment

		mockDeleteFn func(spath string, form *FormParams) (*http.Response, error)
	}{
		"success": {
			wikiID:       1234,
			attachmentID: 8,
			want:         newTestAttachment(),
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
			expectErr:    true,
			mockDeleteFn: newUnexpectedDeleteFn(t, "invalid wikiID"),
		},

		"error-wikiID-negative": {
			wikiID:       -1,
			attachmentID: 8,
			expectErr:    true,
			mockDeleteFn: newUnexpectedDeleteFn(t, "invalid wikiID"),
		},

		"error-attachmentID-zero": {
			wikiID:       1,
			attachmentID: 0,
			expectErr:    true,
			mockDeleteFn: newUnexpectedDeleteFn(t, "invalid attachmentID"),
		},

		"error-attachmentID-negative": {
			wikiID:       1,
			attachmentID: -1,
			expectErr:    true,
			mockDeleteFn: newUnexpectedDeleteFn(t, "invalid attachmentID"),
		},

		"error-client": {
			wikiID:       1234,
			attachmentID: 8,
			expectErr:    true,
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:       1234,
			attachmentID: 8,
			expectErr:    true,
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiAttachmentService()
			s.method.Delete = tc.mockDeleteFn

			attachment, err := s.Remove(tc.wikiID, tc.attachmentID)

			if tc.expectErr {
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

		expectErr bool
		want      []*Attachment

		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "1234",
			want:         newTestAttachmentList(),
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
			expectErr:    true,
			mockGetFn:    newUnexpectedGetFn(t, "invalid issueIDOrKey"),
		},

		"error-client": {
			issueIDOrKey: "1234",
			expectErr:    true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			issueIDOrKey: "1234",
			expectErr:    true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newIssueAttachmentService()
			s.method.Get = tc.mockGetFn

			attachments, err := s.List(tc.issueIDOrKey)

			if tc.expectErr {
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

		expectErr bool
		want      *Attachment

		mockDeleteFn func(spath string, form *FormParams) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "1234",
			attachmentID: 8,
			want:         newTestAttachment(),
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
			expectErr:    true,
			mockDeleteFn: newUnexpectedDeleteFn(t, "invalid issueIDOrKey"),
		},

		"error-attachmentID-zero": {
			issueIDOrKey: "test",
			attachmentID: 0,
			expectErr:    true,
			mockDeleteFn: newUnexpectedDeleteFn(t, "invalid attachmentID"),
		},

		"error-client": {
			issueIDOrKey: "1234",
			attachmentID: 8,
			expectErr:    true,
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			issueIDOrKey: "1234",
			attachmentID: 8,
			expectErr:    true,
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newIssueAttachmentService()
			s.method.Delete = tc.mockDeleteFn

			attachment, err := s.Remove(tc.issueIDOrKey, tc.attachmentID)

			if tc.expectErr {
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

		expectErr bool
		want      []*Attachment

		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "test",
			prNumber:           1234,
			want:               newTestAttachmentList(),
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
			expectErr:          true,
			mockGetFn:          newUnexpectedGetFn(t, "invalid projectIDOrKey"),
		},

		"error-invalid-repository": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "0",
			prNumber:           1,
			expectErr:          true,
			mockGetFn:          newUnexpectedGetFn(t, "invalid repositoryIDOrName"),
		},

		"error-invalid-prNumber": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "1",
			prNumber:           0,
			expectErr:          true,
			mockGetFn:          newUnexpectedGetFn(t, "invalid prNumber"),
		},

		"error-client": {
			projectIDOrKey:     "1234",
			repositoryIDOrName: "test",
			prNumber:           10,
			expectErr:          true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			projectIDOrKey:     "1234",
			repositoryIDOrName: "test",
			prNumber:           10,
			expectErr:          true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newPullRequestAttachmentService()
			s.method.Get = tc.mockGetFn

			attachments, err := s.List(
				tc.projectIDOrKey,
				tc.repositoryIDOrName,
				tc.prNumber,
			)

			if tc.expectErr {
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

		expectErr bool
		want      *Attachment

		mockDeleteFn func(spath string, form *FormParams) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "test",
			prNumber:           1234,
			attachmentID:       8,
			want:               newTestAttachment(),
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
			expectErr:          true,
			mockDeleteFn:       newUnexpectedDeleteFn(t, "invalid projectIDOrKey"),
		},

		"error-invalid-repository": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "",
			prNumber:           1,
			attachmentID:       8,
			expectErr:          true,
			mockDeleteFn:       newUnexpectedDeleteFn(t, "invalid repositoryIDOrName"),
		},

		"error-invalid-prNumber": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "test",
			prNumber:           0,
			attachmentID:       8,
			expectErr:          true,
			mockDeleteFn:       newUnexpectedDeleteFn(t, "invalid prNumber"),
		},

		"error-invalid-attachmentID": {
			projectIDOrKey:     "1",
			repositoryIDOrName: "test",
			prNumber:           1,
			attachmentID:       0,
			expectErr:          true,
			mockDeleteFn:       newUnexpectedDeleteFn(t, "invalid attachmentID"),
		},

		"error-client": {
			projectIDOrKey:     "1234",
			repositoryIDOrName: "test",
			prNumber:           10,
			attachmentID:       8,
			expectErr:          true,
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			projectIDOrKey:     "1234",
			repositoryIDOrName: "test",
			prNumber:           10,
			attachmentID:       8,
			expectErr:          true,
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newPullRequestAttachmentService()
			s.method.Delete = tc.mockDeleteFn

			attachment, err := s.Remove(
				tc.projectIDOrKey,
				tc.repositoryIDOrName,
				tc.prNumber,
				tc.attachmentID,
			)

			if tc.expectErr {
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
