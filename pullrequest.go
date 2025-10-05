package backlog

func validatePRNumber(prNumber int) error {
	if prNumber < 1 {
		return newValidationError("prNumber must not be less than 1")
	}
	return nil
}

// PullRequestService handles communication with the Pull Request-related methods of the Backlog API.
type PullRequestService struct {
	method *method

	Attachment *PullRequestAttachmentService
}
