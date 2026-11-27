package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/priority"
	"github.com/nattokin/go-backlog/internal/model"
)

// ──────────────────────────────────────────────────────────────
//  PriorityService
// ──────────────────────────────────────────────────────────────

// PriorityService handles communication with the priority-related methods of the Backlog API.
type PriorityService struct {
	base *priority.Service
}

// List returns a list of priorities.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-priority-list
func (s *PriorityService) List(ctx context.Context) ([]*Priority, error) {
	v, err := s.base.List(ctx)
	return prioritiesFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newPriorityService(method *client.Method) *PriorityService {
	return &PriorityService{
		base: priority.NewService(method),
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func priorityFromModel(m *model.Priority) *Priority {
	if m == nil {
		return nil
	}
	return &Priority{ID: m.ID, Name: m.Name}
}

func prioritiesFromModel(ms []*model.Priority) []*Priority {
	if ms == nil {
		return nil
	}
	result := make([]*Priority, len(ms))
	for i, v := range ms {
		result[i] = priorityFromModel(v)
	}
	return result
}
