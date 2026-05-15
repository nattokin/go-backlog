// Package repository implements the Backlog Git Repository API service.
package repository

import (
	"context"
	"net/url"
	"path"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service handles git repository-related Backlog API calls.
type Service struct {
	method *core.Method
}

// List returns a list of Git repositories in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-git-repositories
func (s *Service) List(ctx context.Context, projectIDOrKey string) ([]*model.Repository, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories")
	resp, err := s.method.Get(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := []*model.Repository{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// One returns a specific Git repository.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-git-repository
func (s *Service) One(ctx context.Context, projectIDOrKey string, repoIDOrName string) (*model.Repository, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName)
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Repository{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
