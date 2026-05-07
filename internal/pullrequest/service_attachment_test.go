package pullrequest_test

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
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/pullrequest"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func newPRTestAttachment() *model.Attachment {
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

func newPRTestAttachmentList() []*model.Attachment {
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
			want:               newPRTestAttachmentList(),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(
					t,
					"projects/TEST/git/repositories/test/pullRequests/1234/attachments",
					spath,
				)
				return mock.NewJSONResponse(fixture.Attachment.ListJSON), nil
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := pullrequest.NewAttachmentService(method)

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
			want:               newPRTestAttachment(),
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(
					t,
					"projects/TEST/git/repositories/test/pullRequests/1234/attachments/8",
					spath,
				)
				return mock.NewJSONResponse(fixture.Attachment.SingleJSON), nil
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Delete = tc.mockDeleteFn
			s := pullrequest.NewAttachmentService(method)

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

func TestPullRequestAttachmentService_Download(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey     string
		repositoryIDOrName string
		prNumber           int
		attachmentID       int

		mockDownloadFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType     error
		wantFilename    string
		wantContentType string
	}{
		"success": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "repo1",
			prNumber:           5,
			attachmentID:       30,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/git/repositories/repo1/pullRequests/5/attachments/30", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("patch.diff", "text/plain", []byte("DIFF")), nil
			},
			wantFilename:    "patch.diff",
			wantContentType: "text/plain",
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:     "",
			repositoryIDOrName: "repo1",
			prNumber:           5,
			attachmentID:       30,
			wantErrType:        &core.ValidationError{},
		},
		"error-validation-repositoryIDOrName-empty": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "",
			prNumber:           5,
			attachmentID:       30,
			wantErrType:        &core.ValidationError{},
		},
		"error-validation-prNumber-zero": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "repo1",
			prNumber:           0,
			attachmentID:       30,
			wantErrType:        &core.ValidationError{},
		},
		"error-validation-attachmentID-zero": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "repo1",
			prNumber:           5,
			attachmentID:       0,
			wantErrType:        &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey:     "TEST",
			repositoryIDOrName: "repo1",
			prNumber:           5,
			attachmentID:       30,
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
			s := pullrequest.NewAttachmentService(method)

			got, err := s.Download(context.Background(), tc.projectIDOrKey, tc.repositoryIDOrName, tc.prNumber, tc.attachmentID)

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
