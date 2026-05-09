package pullrequest

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// AttachmentService handles attachment-related Backlog API calls for pull requests.
// It delegates all HTTP operations to the shared attachment.Service and is
// responsible only for validation and spath construction.
type AttachmentService struct {
	base *attachment.Service
}

// List returns a list of attachments on the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *AttachmentService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) ([]*model.Attachment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments")
	return s.base.List(ctx, spath)
}

// Remove removes an attachment from the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *AttachmentService) Remove(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*model.Attachment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments", strconv.Itoa(attachmentID))
	return s.base.Remove(ctx, spath)
}

// Download downloads an attachment from the pull request.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/download-pull-request-attachment
func (s *AttachmentService) Download(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*model.FileData, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments", strconv.Itoa(attachmentID))
	return s.base.Download(ctx, spath)
}

func NewAttachmentService(method *core.Method) *AttachmentService {
	return &AttachmentService{base: attachment.NewService(method)}
}
