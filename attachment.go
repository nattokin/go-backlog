package backlog

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/validate"
)

type SpaceAttachmentService = attachment.SpaceAttachmentService

type WikiAttachmentService = attachment.WikiAttachmentService

// IssueAttachmentService handles communication with the issue attachment-related methods of the Backlog API.
type IssueAttachmentService struct {
	method *core.Method
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(ctx context.Context, issueIDOrKey string) ([]*Attachment, error) {
	if err := validateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments")
	return attachment.ListAttachments(ctx, s.method, spath)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(ctx context.Context, issueIDOrKey string, attachmentID int) (*Attachment, error) {
	if err := validateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments", strconv.Itoa(attachmentID))
	return attachment.RemoveAttachment(ctx, s.method, spath)
}

// PullRequestAttachmentService handles communication with the pull request attachment-related methods of the Backlog API.
type PullRequestAttachmentService struct {
	method *core.Method
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) ([]*Attachment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validatePRNumber(prNumber); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments")
	return attachment.ListAttachments(ctx, s.method, spath)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*Attachment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validatePRNumber(prNumber); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments", strconv.Itoa(attachmentID))
	return attachment.RemoveAttachment(ctx, s.method, spath)
}
