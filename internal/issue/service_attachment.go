package issue

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// AttachmentService handles attachment-related Backlog API calls for issues.
// It delegates all HTTP operations to the shared attachment.Service and is
// responsible only for validation and spath construction.
type AttachmentService struct {
	svc *attachment.Service
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *AttachmentService) List(ctx context.Context, issueIDOrKey string) ([]*model.Attachment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments")
	return s.svc.List(ctx, spath)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *AttachmentService) Remove(ctx context.Context, issueIDOrKey string, attachmentID int) (*model.Attachment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments", strconv.Itoa(attachmentID))
	return s.svc.Remove(ctx, spath)
}

// Download downloads a file attached to the issue.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-attachment
func (s *AttachmentService) Download(ctx context.Context, issueIDOrKey string, attachmentID int) (*model.FileData, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments", strconv.Itoa(attachmentID))
	return s.svc.Download(ctx, spath)
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewAttachmentService creates and returns a new issue AttachmentService.
func NewAttachmentService(method *core.Method) *AttachmentService {
	return &AttachmentService{svc: attachment.NewService(method)}
}
