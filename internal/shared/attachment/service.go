// Package attachment implements shared HTTP logic for attachment-related Backlog API endpoints.
package attachment

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// Service holds shared HTTP logic for attachment-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *core.Method
}

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
// The caller is responsible for closing FileData.Body after use.
func (s *Service) Download(ctx context.Context, spath string) (*model.FileData, error) {
	resp, err := s.method.Download(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	return core.DownloadResponse(resp)
}

func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
