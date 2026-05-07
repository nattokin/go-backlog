package attachment

import (
	"context"
	"io"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// Service holds shared HTTP logic for attachment-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *core.Method
}

// List returns the list of attachments at spath.
func (s *Service) List(ctx context.Context, spath string) ([]*model.Attachment, error) {
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Attachment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Remove deletes the attachment at spath and returns the deleted attachment.
func (s *Service) Remove(ctx context.Context, spath string) (*model.Attachment, error) {
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.Attachment{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// Download streams the attachment at spath.
func (s *Service) Download(ctx context.Context, spath string) (*model.FileData, error) {
	resp, err := s.method.Download(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	return core.DownloadResponse(resp)
}

// Upload uploads a file to spath.
func (s *Service) Upload(ctx context.Context, spath string, fileName string, r io.Reader) (*model.Attachment, error) {
	resp, err := s.method.Upload(ctx, spath, fileName, r)
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

// NewService creates and returns a new attachment Service.
func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
