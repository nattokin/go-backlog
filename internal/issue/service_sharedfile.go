package issue

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/sharedfile"
	"github.com/nattokin/go-backlog/internal/validate"
)

// SharedFileService handles shared-file-related Backlog API calls for issues.
// It delegates all HTTP operations to the shared sharedfile.Service and is
// responsible only for validation and spath construction.
type SharedFileService struct {
	base *sharedfile.Service
}

// List returns a list of shared files linked to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-linked-shared-files
func (s *SharedFileService) List(ctx context.Context, issueIDOrKey string) ([]*model.SharedFile, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles")
	return s.base.List(ctx, spath)
}

// Link links shared files to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-issue
func (s *SharedFileService) Link(ctx context.Context, issueIDOrKey string, fileIDs []int) ([]*model.SharedFile, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles")
	return s.base.Link(ctx, spath, fileIDs)
}

// Unlink removes a shared file link from the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-issue
func (s *SharedFileService) Unlink(ctx context.Context, issueIDOrKey string, fileID int) (*model.SharedFile, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateSharedFileID(fileID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles", strconv.Itoa(fileID))
	return s.base.Unlink(ctx, spath)
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewSharedFileService creates and returns a new issue SharedFileService.
func NewSharedFileService(method *core.Method) *SharedFileService {
	return &SharedFileService{base: sharedfile.NewService(method)}
}
