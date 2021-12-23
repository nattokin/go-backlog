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

// RepositoryService has methods for Repository.
type RepositoryService struct {
	method *method
}
