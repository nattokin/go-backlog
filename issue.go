package backlog

// IssueService has methods for Issue.
type IssueService struct {
	clientMethod *clientMethod

	Attachment *IssueAttachmentService
}

func newIssueService(cm *clientMethod) *IssueService {
	return &IssueService{
		clientMethod: cm,
		Attachment:   newIssueAttachmentService(cm),
	}
}
