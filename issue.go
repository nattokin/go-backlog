package backlog

// IssueService has methods for Issue.
type IssueService struct {
	method *method

	Attachment *IssueAttachmentService
}
