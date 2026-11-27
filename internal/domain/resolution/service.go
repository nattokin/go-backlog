// Package resolution implements the Backlog Resolution API service.
package resolution

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
)

// Service handles resolution-related Backlog API calls.
type Service struct {
	method *client.Method
}

// List returns a list of resolutions.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-resolution-list
func (s *Service) List(ctx context.Context) ([]*model.Resolution, error) {
	resp, err := s.method.Get(ctx, "resolutions", nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Resolution{}
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
