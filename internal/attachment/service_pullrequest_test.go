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

func TestPullRequestAttachmentService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey     string
		repositoryIDOrName string
		prNumber           int

		expectError bool
		want        []*model.Attachment

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
						bytes.NewReader([]byte(fixture.Attachment.ListJSON)),
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
			s := attachment.NewPullRequestService(method)

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
		want        *model.Attachment

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
						bytes.NewReader([]byte(fixture.Attachment.SingleJSON)),
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
			s := attachment.NewPullRequestService(method)

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
