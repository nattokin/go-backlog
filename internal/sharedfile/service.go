package sharedfile

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service holds shared HTTP logic for shared-file-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *core.Method
}

// List returns the list of shared files at spath.
func (s *Service) List(ctx context.Context, spath string) ([]*model.SharedFile, error) {
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Link links shared files at spath using the given fileIDs.
func (s *Service) Link(ctx context.Context, spath string, fileIDs []int) ([]*model.SharedFile, error) {
	if len(fileIDs) == 0 {
		return nil, errors.New("fileIDs must not be empty")
	}

	form := url.Values{}
	for _, id := range fileIDs {
		if err := validate.ValidateSharedFileID(id); err != nil {
			return nil, err
		}
		form.Add("fileId[]", strconv.Itoa(id))
	}

	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Unlink removes a shared file link at spath.
func (s *Service) Unlink(ctx context.Context, spath string) (*model.SharedFile, error) {
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.SharedFile{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewService creates and returns a new shared-file Service.
func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
