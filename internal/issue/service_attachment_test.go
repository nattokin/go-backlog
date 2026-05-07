package issue_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func newIssueTestAttachment() *model.Attachment {
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

func newIssueTestAttachmentList() []*model.Attachment {
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

func TestIssueAttachmentService_List(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		expectError bool
		want        []*model.Attachment

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "1234",
			want:         newIssueTestAttachmentList(),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1234/attachments", spath)
				return mock.NewJSONResponse(fixture.Attachment.ListJSON), nil
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := issue.NewAttachmentService(method)

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
		want        *model.Attachment

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "1234",
			attachmentID: 8,
			want:         newIssueTestAttachment(),
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1234/attachments/8", spath)
				return mock.NewJSONResponse(fixture.Attachment.SingleJSON), nil
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Delete = tc.mockDeleteFn
			s := issue.NewAttachmentService(method)

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

func TestIssueAttachmentService_Download(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		attachmentID int

		mockDownloadFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType     error
		wantFilename    string
		wantContentType string
	}{
		"success": {
			issueIDOrKey: "TEST-1",
			attachmentID: 10,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/attachments/10", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("file.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "file.png",
			wantContentType: "image/png",
		},
		"success-issue-id": {
			issueIDOrKey: "123",
			attachmentID: 5,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/123/attachments/5", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("doc.pdf", "application/pdf", []byte("PDF")), nil
			},
			wantFilename:    "doc.pdf",
			wantContentType: "application/pdf",
		},
		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey: "",
			attachmentID: 10,
			wantErrType:  &core.ValidationError{},
		},
		"error-validation-attachmentID-zero": {
			issueIDOrKey: "TEST-1",
			attachmentID: 0,
			wantErrType:  &core.ValidationError{},
		},
		"error-client-network": {
			issueIDOrKey: "TEST-1",
			attachmentID: 10,
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
			s := issue.NewAttachmentService(method)

			got, err := s.Download(context.Background(), tc.issueIDOrKey, tc.attachmentID)

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
