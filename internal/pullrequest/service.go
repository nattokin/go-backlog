package pullrequest

import (
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
)

// PullRequestService handles communication with the Pull Request-related methods of the Backlog API.
type PullRequestService struct {
	method *core.Method

	Attachment *attachment.PullRequestAttachmentService
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewPullRequestService(method *core.Method) *PullRequestService {
	return &PullRequestService{
		method:     method,
		Attachment: attachment.NewPullRequestAttachmentService(method),
	}
}
