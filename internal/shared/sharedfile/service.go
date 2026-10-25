// Package sharedfile implements shared HTTP logic for shared-file-related Backlog API endpoints.
package sharedfile

import (
	"context"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
)

// Service holds shared HTTP logic for shared-file-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *client.Method
}

func (s *Service) List(ctx context.Context, spath string) ([]*model.SharedFile, error) {
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Service) Link(ctx context.Context, spath string, fileIDs []int) ([]*model.SharedFile, error) {
	form := url.Values{}
	for _, id := range fileIDs {
		form.Add("fileId[]", strconv.Itoa(id))
	}

	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Service) Unlink(ctx context.Context, spath string) (*model.SharedFile, error) {
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.SharedFile{}
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

func NewService(method *client.Method) *Service {
	return &Service{method: method}
}
