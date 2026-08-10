package wiki_test

import (
	"context"
	"encoding/json"
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

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantIDs                []int
	}{
		"success-single": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			wantIDs:       []int{2},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments", spath)
				v := form
				assert.Equal(t, []string{"2"}, v["attachmentId[]"])
				return mock.NewResponse(fixture.Attachment.SingleListJSON), nil
			},
		},

		"success-multiple": {
			wikiID:        1,
			attachmentIDs: []int{2, 5},
			wantIDs:       []int{2, 5},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.Attachment.ListJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-wikiID-invalid": {
			wikiID:                 0,
			attachmentIDs:          []int{1, 2},
			wantValidationErrCount: 1,
		},
		"error-validation-attachmentIDs-invalid": {
			wikiID:                 1,
			attachmentIDs:          []int{0, 1, 2},
			wantValidationErrCount: 1,
		},
		"error-validation-attachmentIDs-empty": {
			wikiID:                 1,
			attachmentIDs:          []int{},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			wikiID:                 0,
			attachmentIDs:          []int{},
			wantValidationErrCount: 2,
		},

		// --- other errors ---
		"error-client-network": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID:        1234,
			attachmentIDs: []int{2},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := wiki.NewAttachmentService(method)

			attachments, err := s.Attach(context.Background(), tc.wikiID, tc.attachmentIDs)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantIDs                []int
	}{
		"success": {
			wikiID:  1234,
			wantIDs: []int{2, 5},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments", spath)
				return mock.NewResponse(fixture.Attachment.ListJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-wikiID-zero": {
			wikiID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-wikiID-negative": {
			wikiID:                 -1,
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-client-network": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := wiki.NewAttachmentService(method)

			attachments, err := s.List(context.Background(), tc.wikiID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, attachments)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantID                 int
	}{
		"success": {
			wikiID:       1234,
			attachmentID: 8,
			wantID:       8,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/attachments/8", spath)
				return mock.NewResponse(fixture.Attachment.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-wikiID-zero": {
			wikiID:                 0,
			attachmentID:           8,
			wantValidationErrCount: 1,
		},
		"error-validation-wikiID-negative": {
			wikiID:                 -1,
			attachmentID:           8,
			wantValidationErrCount: 1,
		},
		"error-validation-attachmentID-zero": {
			wikiID:                 1,
			attachmentID:           0,
			wantValidationErrCount: 1,
		},
		"error-validation-attachmentID-negative": {
			wikiID:                 1,
			attachmentID:           -1,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			wikiID:                 0,
			attachmentID:           0,
			wantValidationErrCount: 2,
		},

		// --- other errors ---
		"error-client-network": {
			wikiID:       1234,
			attachmentID: 8,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID:       1234,
			attachmentID: 8,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := wiki.NewAttachmentService(method)

			attachment, err := s.Remove(context.Background(), tc.wikiID, tc.attachmentID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, attachment)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, attachment)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

		wantErrType            error
		wantValidationErrCount int
		wantFilename           string
		wantContentType        string
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

		// --- validation errors ---
		"error-validation-wikiID-zero": {
			wikiID:                 0,
			attachmentID:           20,
			wantValidationErrCount: 1,
		},
		"error-validation-attachmentID-zero": {
			wikiID:                 34,
			attachmentID:           0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			wikiID:                 0,
			attachmentID:           0,
			wantValidationErrCount: 2,
		},

		// --- other errors ---
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

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
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
