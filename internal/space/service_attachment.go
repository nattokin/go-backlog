package space

import (
	"context"
	"io"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// AttachmentService handles attachment-related Backlog API calls for a space.
// It delegates all HTTP operations to the shared attachment.Service.
type AttachmentService struct {
	method *core.Method
}

// Upload uploads any file to the space.
//
// The file name must not be empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/post-attachment-file
func (s *AttachmentService) Upload(ctx context.Context, fileName string, r io.Reader) (*model.Attachment, error) {
	resp, err := s.method.Upload(ctx, "space/attachment", fileName, r)
	if err != nil {
		return nil, err
	}

	v := model.Attachment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil

}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewAttachmentService creates and returns a new space AttachmentService.
func NewAttachmentService(method *core.Method) *AttachmentService {
	return &AttachmentService{method: method}
}
