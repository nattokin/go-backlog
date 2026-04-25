package attachment_test

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

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueAttachmentService_List(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		expectError bool
		want        []*model.Attachment

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
						bytes.NewReader([]byte(fixture.Attachment.ListJSON)),
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
			s := attachment.NewIssueService(method)

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
			want:         newTestAttachment(),
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1234/attachments/8", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(fixture.Attachment.SingleJSON)),
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
			s := attachment.NewIssueService(method)

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
