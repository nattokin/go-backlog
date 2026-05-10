package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
)

// IssueAttachmentService handles communication with the issue attachment-related methods of the Backlog API.
type IssueAttachmentService struct {
	base *issue.AttachmentService
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(ctx context.Context, issueIDOrKey string) ([]*Attachment, error) {
	v, err := s.base.List(ctx, issueIDOrKey)
	return attachmentsFromModel(v), convertError(err)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(ctx context.Context, issueIDOrKey string, attachmentID int) (*Attachment, error) {
	v, err := s.base.Remove(ctx, issueIDOrKey, attachmentID)
	return attachmentFromModel(v), convertError(err)
}

// Download downloads a file attached to the issue.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-attachment
func (s *IssueAttachmentService) Download(ctx context.Context, issueIDOrKey string, attachmentID int) (*FileData, error) {
	v, err := s.base.Download(ctx, issueIDOrKey, attachmentID)
	return fileDataFromModel(v), convertError(err)
}

func newIssueAttachmentService(method *core.Method) *IssueAttachmentService {
	return &IssueAttachmentService{
		base: issue.NewAttachmentService(method),
	}
}
