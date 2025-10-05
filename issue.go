package backlog

func validateIssueIDOrKey(issueIDOrKey string) error {
	if issueIDOrKey == "" {
		return newValidationError("issueIDOrKey must not be empty")
	}
	if issueIDOrKey == "0" {
		return newValidationError("issueIDOrKey must not be '0'")
	}
	return nil
}

// IssueService handles communication with the issue-related methods of the Backlog API.
type IssueService struct {
	method *method

	Attachment *IssueAttachmentService
}
