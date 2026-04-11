package backlog

import "github.com/nattokin/go-backlog/internal/core"

func validateRepositoryIDOrName(repositoryIDOrName string) error {
	if repositoryIDOrName == "" {
		return core.NewValidationError("repositoryIDOrName must not be empty")
	}
	if repositoryIDOrName == "0" {
		return core.NewValidationError("repositoryIDOrName must not be '0'")
	}
	return nil
}

// RepositoryService handles communication with the repository-related methods of the Backlog API.
//
//nolint:unused // API not implemented yet
type RepositoryService struct {
	method *core.Method
}
