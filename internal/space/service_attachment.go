package space

import (
	"context"
	"io"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// AttachmentService handles attachment-related Backlog API calls for a space.
type AttachmentService struct {
	method *core.Method
}

// Upload uploads a file to the space.
// The file name must not be empty.
// The caller is responsible for closing the reader after use.
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

func NewAttachmentService(method *core.Method) *AttachmentService {
	return &AttachmentService{method: method}
}
