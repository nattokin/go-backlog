package backlog

import (
	"time"

	"github.com/nattokin/go-backlog/internal/core"
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
