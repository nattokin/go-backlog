package wiki_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/wiki"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestWikiAttachmentService_Attach(t *testing.T) {
	cases := map[string]struct {
		wikiID        int
		attachmentIDs []int

		expectError bool
		wantIDs     []int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success-single": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			wantIDs:       []int{2},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments", spath)
				v := form
				assert.Equal(t, []string{"2"}, v["attachmentId[]"])
				return mock.NewJSONResponse(fixture.Attachment.SingleListJSON), nil
			},
		},

		"success-multiple": {
			wikiID:        1,
			attachmentIDs: []int{2, 5},
			wantIDs:       []int{2, 5},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.Attachment.ListJSON), nil
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Post = tc.mockPostFn
			s := wiki.NewAttachmentService(method)

			attachments, err := s.Attach(context.Background(), tc.wikiID, tc.attachmentIDs)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachments)

			assert.Len(t, attachments, len(tc.wantIDs))

			for i, id := range tc.wantIDs {
				assert.Equal(t, id, attachments[i].ID)
			}
		})
	}
}

func TestWikiAttachmentService_List(t *testing.T) {
	cases := map[string]struct {
		wikiID int

		expectError bool
		wantIDs     []int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID:  1234,
			wantIDs: []int{2, 5},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments", spath)
				return mock.NewJSONResponse(fixture.Attachment.ListJSON), nil
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := wiki.NewAttachmentService(method)

			attachments, err := s.List(context.Background(), tc.wikiID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachments)

			assert.Len(t, attachments, len(tc.wantIDs))

			for i, id := range tc.wantIDs {
				assert.Equal(t, id, attachments[i].ID)
			}
		})
	}
}

func TestWikiAttachmentService_Remove(t *testing.T) {
	cases := map[string]struct {
		wikiID       int
		attachmentID int

		expectError bool
		wantID      int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID:       1234,
			attachmentID: 8,
			wantID:       8,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments/8", spath)
				return mock.NewJSONResponse(fixture.Attachment.SingleJSON), nil
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Delete = tc.mockDeleteFn
			s := wiki.NewAttachmentService(method)

			attachment, err := s.Remove(context.Background(), tc.wikiID, tc.attachmentID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, attachment)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, attachment)

			assert.Equal(t, tc.wantID, attachment.ID)
		})
	}
}

func TestWikiAttachmentService_Download(t *testing.T) {
	cases := map[string]struct {
		wikiID       int
		attachmentID int

		mockDownloadFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType     error
		wantFilename    string
		wantContentType string
	}{
		"success": {
			wikiID:       34,
			attachmentID: 20,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34/attachments/20", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("doc.pdf", "application/pdf", []byte("PDF")), nil
			},
			wantFilename:    "doc.pdf",
			wantContentType: "application/pdf",
		},
		"error-validation-wikiID-zero": {
			wikiID:       0,
			attachmentID: 20,
			wantErrType:  &core.ValidationError{},
		},
		"error-validation-attachmentID-zero": {
			wikiID:       34,
			attachmentID: 0,
			wantErrType:  &core.ValidationError{},
		},
		"error-client-network": {
			wikiID:       34,
			attachmentID: 20,
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
			s := wiki.NewAttachmentService(method)

			got, err := s.Download(context.Background(), tc.wikiID, tc.attachmentID)

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
