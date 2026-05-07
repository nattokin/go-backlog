package space

import (
	"context"
	"io"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// AttachmentService handles attachment-related Backlog API calls for a space.
// It delegates all HTTP operations to the shared attachment.Service.
type AttachmentService struct {
	svc *attachment.Service
}

// Upload uploads any file to the space.
//
// The file name must not be empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/post-attachment-file
func (s *AttachmentService) Upload(ctx context.Context, fileName string, r io.Reader) (*model.Attachment, error) {
	return s.svc.Upload(ctx, "space/attachment", fileName, r)
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewAttachmentService creates and returns a new space AttachmentService.
func NewAttachmentService(method *core.Method) *AttachmentService {
	return &AttachmentService{svc: attachment.NewService(method)}
}
