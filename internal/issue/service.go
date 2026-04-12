package issue

import (
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
)

// IssueService handles communication with the issue-related methods of the Backlog API.
type IssueService struct {
	method *core.Method

	Attachment *attachment.IssueAttachmentService
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

// NewWikiService returns a new WikiService.
func NewIssueService(method *core.Method, option *core.OptionService) *IssueService {
	return &IssueService{
		method:     method,
		Attachment: attachment.NewIssueAttachmentService(method),
	}
}
