package backlog

// IssueService has methods for Issue.
type IssueService struct {
	clientMethod *clientMethod

	Attachment *IssueAttachmentService
}

func newIssueService(cm *clientMethod) *IssueService {
	as := &IssueAttachmentService{
		baseAttachmentService: &baseAttachmentService{
			clientMethod: cm,
		},
	}
	return &IssueService{
		clientMethod: cm,
		Attachment:   as,
	}
}
