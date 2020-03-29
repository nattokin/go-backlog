package backlog

// PullRequestService has methods for Issue.
type PullRequestService struct {
	clientMethod *clientMethod

	Attachment *PullRequestAttachmentService
}

func newPullRequestService(cm *clientMethod) *PullRequestService {
	return &PullRequestService{
		clientMethod: cm,
		Attachment:   newPullRequestAttachmentService(cm),
	}
}
