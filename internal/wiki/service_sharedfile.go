package wiki

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/sharedfile"
	"github.com/nattokin/go-backlog/internal/validate"
)

// SharedFileService handles shared-file-related Backlog API calls for wiki pages.
// It delegates all HTTP operations to the shared sharedfile.Service and is
// responsible only for validation and spath construction.
type SharedFileService struct {
	base *sharedfile.Service
}

// List returns a list of shared files linked to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files-on-wiki
func (s *SharedFileService) List(ctx context.Context, wikiID int) ([]*model.SharedFile, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "sharedFiles")
	return s.base.List(ctx, spath)
}

// Link links shared files to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-wiki
func (s *SharedFileService) Link(ctx context.Context, wikiID int, fileIDs []int) ([]*model.SharedFile, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "sharedFiles")
	return s.base.Link(ctx, spath, fileIDs)
}

// Unlink removes a shared file link from the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-wiki
func (s *SharedFileService) Unlink(ctx context.Context, wikiID, fileID int) (*model.SharedFile, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}
	if err := validate.ValidateSharedFileID(fileID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "sharedFiles", strconv.Itoa(fileID))
	return s.base.Unlink(ctx, spath)
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewSharedFileService creates and returns a new wiki SharedFileService.
func NewSharedFileService(method *core.Method) *SharedFileService {
	return &SharedFileService{base: sharedfile.NewService(method)}
}
