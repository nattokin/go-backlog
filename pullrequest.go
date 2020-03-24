package backlog

// PullRequestService has methods for Issue.
type PullRequestService struct {
	clientMethod *clientMethod

	Attachment *PullRequestAttachmentService
}

func newPullRequestService(cm *clientMethod) *PullRequestService {
	as := &PullRequestAttachmentService{
		baseAttachmentService: &baseAttachmentService{
			clientMethod: cm,
		},
	}
	return &PullRequestService{
		clientMethod: cm,
		Attachment:   as,
	}
}
