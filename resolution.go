package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/resolution"
	"github.com/nattokin/go-backlog/internal/model"
)

// ──────────────────────────────────────────────────────────────
//  ResolutionService
// ──────────────────────────────────────────────────────────────

// ResolutionService handles communication with the resolution-related methods of the Backlog API.
type ResolutionService struct {
	base *resolution.Service
}

// List returns a list of resolutions.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-resolution-list
func (s *ResolutionService) List(ctx context.Context) ([]*Resolution, error) {
	v, err := s.base.List(ctx)
	return resolutionsFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newResolutionService(method *client.Method) *ResolutionService {
	return &ResolutionService{
		base: resolution.NewService(method),
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func resolutionsFromModel(ms []*model.Resolution) []*Resolution {
	if ms == nil {
		return nil
	}
	result := make([]*Resolution, len(ms))
	for i, v := range ms {
		result[i] = resolutionFromModel(v)
	}
	return result
}
