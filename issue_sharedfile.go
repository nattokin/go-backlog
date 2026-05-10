package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
)

// IssueSharedFileService handles communication with the issue shared-file-related methods of the Backlog API.
type IssueSharedFileService struct {
	base *issue.SharedFileService
}

// List returns a list of shared files linked to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-linked-shared-files
func (s *IssueSharedFileService) List(ctx context.Context, issueIDOrKey string) ([]*SharedFile, error) {
	v, err := s.base.List(ctx, issueIDOrKey)
	return sharedFilesFromModel(v), convertError(err)
}

// Link links shared files to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-issue
func (s *IssueSharedFileService) Link(ctx context.Context, issueIDOrKey string, fileIDs []int) ([]*SharedFile, error) {
	v, err := s.base.Link(ctx, issueIDOrKey, fileIDs)
	return sharedFilesFromModel(v), convertError(err)
}

// Unlink removes a shared file link from the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-issue
func (s *IssueSharedFileService) Unlink(ctx context.Context, issueIDOrKey string, fileID int) (*SharedFile, error) {
	v, err := s.base.Unlink(ctx, issueIDOrKey, fileID)
	return sharedFileFromModel(v), convertError(err)
}

func newIssueSharedFileService(method *core.Method) *IssueSharedFileService {
	return &IssueSharedFileService{
		base: issue.NewSharedFileService(method),
	}
}
