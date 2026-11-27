// Package priority implements the Backlog Priority API service.
package priority

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
)

// Service handles priority-related Backlog API calls.
type Service struct {
	method *client.Method
}

// List returns a list of priorities.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-priority-list
func (s *Service) List(ctx context.Context) ([]*model.Priority, error) {
	resp, err := s.method.Get(ctx, "priorities", nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Priority{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// NewService returns a *Service that sends requests through method.
func NewService(method *client.Method) *Service {
	return &Service{
		method: method,
	}
}
