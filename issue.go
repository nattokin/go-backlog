package backlog

import "github.com/nattokin/go-backlog/internal/core"

func validateIssueIDOrKey(issueIDOrKey string) error {
	if issueIDOrKey == "" {
		return core.NewValidationError("issueIDOrKey must not be empty")
	}
	if issueIDOrKey == "0" {
		return core.NewValidationError("issueIDOrKey must not be '0'")
	}
	return nil
}

// IssueService handles communication with the issue-related methods of the Backlog API.
type IssueService struct {
	method *core.Method

	Attachment *IssueAttachmentService
}
