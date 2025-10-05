package backlog

func validateRepositoryIDOrName(repositoryIDOrName string) error {
	if repositoryIDOrName == "" {
		return newValidationError("repositoryIDOrName must not be empty")
	}
	if repositoryIDOrName == "0" {
		return newValidationError("repositoryIDOrName must not be '0'")
	}
	return nil
}

// RepositoryService handles communication with the repository-related methods of the Backlog API.
type RepositoryService struct {
	method *method
}
