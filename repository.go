package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/repository"
)

// Repository represents repository of Backlog git.
type Repository struct {
	ID           int
	ProjectID    int
	Name         string
	Description  string
	HookURL      string
	HTTPURL      string
	SSHURL       string
	DisplayOrder int
	PushedAt     time.Time
	CreatedUser  *User
	Created      time.Time
	UpdatedUser  *User
	Updated      time.Time
}

// ──────────────────────────────────────────────────────────────
//  RepositoryService
// ──────────────────────────────────────────────────────────────

// RepositoryService handles communication with the repository-related methods of the Backlog API.
type RepositoryService struct {
	base *repository.Service
}

// All returns a list of Git repositories in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-git-repositories
func (s *RepositoryService) All(ctx context.Context, projectIDOrKey string) ([]*Repository, error) {
	v, err := s.base.All(ctx, projectIDOrKey)
	return repositoriesFromModel(v), convertError(err)
}

// One returns information about a specific Git repository.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-git-repository
func (s *RepositoryService) One(ctx context.Context, projectIDOrKey string, repoIDOrName string) (*Repository, error) {
	v, err := s.base.One(ctx, projectIDOrKey, repoIDOrName)
	return repositoryFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newRepositoryService(method *core.Method) *RepositoryService {
	return &RepositoryService{
		base: repository.NewService(method),
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func repositoryFromModel(m *model.Repository) *Repository {
	if m == nil {
		return nil
	}
	return &Repository{
		ID:           m.ID,
		ProjectID:    m.ProjectID,
		Name:         m.Name,
		Description:  m.Description,
		HookURL:      m.HookURL,
		HTTPURL:      m.HTTPURL,
		SSHURL:       m.SSHURL,
		DisplayOrder: m.DisplayOrder,
		PushedAt:     m.PushedAt,
		CreatedUser:  userFromModel(m.CreatedUser),
		Created:      m.Created,
		UpdatedUser:  userFromModel(m.UpdatedUser),
		Updated:      m.Updated,
	}
}

func repositoriesFromModel(ms []*model.Repository) []*Repository {
	result := make([]*Repository, len(ms))
	for i, v := range ms {
		result[i] = repositoryFromModel(v)
	}
	return result
}
