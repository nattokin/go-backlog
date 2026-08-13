package issue

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/shared/attachment"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// AttachmentService handles attachment-related Backlog API calls for issues.
// It delegates all HTTP operations to the shared attachment.Service and is
// responsible only for validation and spath construction.
type AttachmentService struct {
	base *attachment.Service
}

// List returns a list of attachments on the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *AttachmentService) List(ctx context.Context, issueIDOrKey string) ([]*model.Attachment, error) {
	var ves validation.Errors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "attachments")
	return s.base.List(ctx, spath)
}

// Remove removes an attachment from the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *AttachmentService) Remove(ctx context.Context, issueIDOrKey string, attachmentID int) (*model.Attachment, error) {
	var ves validation.Errors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateAttachmentID(attachmentID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "attachments", strconv.Itoa(attachmentID))
	return s.base.Remove(ctx, spath)
}

// Download downloads an attachment from the issue.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-attachment
func (s *AttachmentService) Download(ctx context.Context, issueIDOrKey string, attachmentID int) (*model.FileData, error) {
	var ves validation.Errors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateAttachmentID(attachmentID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "attachments", strconv.Itoa(attachmentID))
	return s.base.Download(ctx, spath)
}

func NewAttachmentService(method *client.Method) *AttachmentService {
	return &AttachmentService{base: attachment.NewService(method)}
}
