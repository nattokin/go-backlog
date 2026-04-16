package backlog

import (
	"time"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
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
//
//nolint:unused // API not implemented yet
type RepositoryService struct {
	method *core.Method
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

//nolint:unused
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
