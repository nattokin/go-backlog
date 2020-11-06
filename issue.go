package backlog

// IssueService has methods for Issue.
type IssueService struct {
	method *method

	Attachment *IssueAttachmentService
}

// IssueAttachmentService hs methods for attachment file of issue.
type IssueAttachmentService struct {
	*AttachmentService
}
