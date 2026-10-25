package issue

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/shared/sharedfile"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
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
	var ves validation.Errors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles")
	return s.base.List(ctx, spath)
}

// Link links shared files to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-issue
func (s *SharedFileService) Link(ctx context.Context, issueIDOrKey string, fileIDs []int) ([]*model.SharedFile, error) {
	var ves validation.Errors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(fileIDs) == 0 {
		ves = append(ves, validation.NewError("fileIDs", "fileIDs must not be empty"))
	} else {
		for _, id := range fileIDs {
			if ve := validate.ValidateSharedFileID(id); ve != nil {
				ves = append(ves, ve)
				break
			}
		}
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles")
	return s.base.Link(ctx, spath, fileIDs)
}

// Unlink removes a shared file link from the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-issue
func (s *SharedFileService) Unlink(ctx context.Context, issueIDOrKey string, fileID int) (*model.SharedFile, error) {
	var ves validation.Errors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateSharedFileID(fileID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles", strconv.Itoa(fileID))
	return s.base.Unlink(ctx, spath)
}

func NewSharedFileService(method *client.Method) *SharedFileService {
	return &SharedFileService{base: sharedfile.NewService(method)}
}
