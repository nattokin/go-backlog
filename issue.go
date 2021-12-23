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

// IssueService has methods for Issue.
type IssueService struct {
	method *method

	Attachment *IssueAttachmentService
}
