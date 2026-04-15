package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/pullrequest"
)

// ──────────────────────────────────────────────────────────────
//  PullRequestService
// ──────────────────────────────────────────────────────────────

// PullRequestService handles communication with the pull request-related methods of the Backlog API.
type PullRequestService struct {
	base *pullrequest.PullRequestService

	Attachment *PullRequestAttachmentService
}

// ──────────────────────────────────────────────────────────────
//  PullRequestAttachmentService
// ──────────────────────────────────────────────────────────────

// PullRequestAttachmentService handles communication with the pull request attachment-related methods of the Backlog API.
type PullRequestAttachmentService struct {
	base *attachment.PullRequestAttachmentService
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) ([]*model.Attachment, error) {
	return s.base.List(ctx, projectIDOrKey, repositoryIDOrName, prNumber)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*model.Attachment, error) {
	return s.base.Remove(ctx, projectIDOrKey, repositoryIDOrName, prNumber, attachmentID)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newPullRequestService(method *core.Method) *PullRequestService {
	return &PullRequestService{
		base:       pullrequest.NewPullRequestService(method),
		Attachment: newPullRequestAttachmentService(method),
	}
}

func newPullRequestAttachmentService(method *core.Method) *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		base: attachment.NewPullRequestAttachmentService(method),
	}
}
