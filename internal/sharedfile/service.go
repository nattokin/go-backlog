package sharedfile

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// ──────────────────────────────────────────────────────────────
//  IssueService
// ──────────────────────────────────────────────────────────────

// IssueService handles communication with the issue shared-file-related methods of the Backlog API.
type IssueService struct {
	method *core.Method
}

// List returns a list of shared files linked to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-linked-shared-files
func (s *IssueService) List(ctx context.Context, issueIDOrKey string) ([]*model.SharedFile, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Link links shared files to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-issue
func (s *IssueService) Link(ctx context.Context, issueIDOrKey string, fileIDs []int) ([]*model.SharedFile, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if len(fileIDs) == 0 {
		return nil, errors.New("fileIDs must not be empty")
	}

	form := url.Values{}
	for _, id := range fileIDs {
		if err := validate.ValidateSharedFileID(id); err != nil {
			return nil, err
		}
		form.Add("fileId[]", strconv.Itoa(id))
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Unlink removes a shared file link from the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-issue
func (s *IssueService) Unlink(ctx context.Context, issueIDOrKey string, fileID int) (*model.SharedFile, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateSharedFileID(fileID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "sharedFiles", strconv.Itoa(fileID))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.SharedFile{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// NewIssueService creates and returns a new shared-file IssueService.
func NewIssueService(method *core.Method) *IssueService {
	return &IssueService{method: method}
}

// ──────────────────────────────────────────────────────────────
//  WikiService
// ──────────────────────────────────────────────────────────────

// WikiService handles communication with the wiki shared-file-related methods of the Backlog API.
type WikiService struct {
	method *core.Method
}

// List returns a list of shared files linked to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files-on-wiki
func (s *WikiService) List(ctx context.Context, wikiID int) ([]*model.SharedFile, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "sharedFiles")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Link links shared files to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-wiki
func (s *WikiService) Link(ctx context.Context, wikiID int, fileIDs []int) ([]*model.SharedFile, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}
	if len(fileIDs) == 0 {
		return nil, errors.New("fileIDs must not be empty")
	}

	form := url.Values{}
	for _, id := range fileIDs {
		if err := validate.ValidateSharedFileID(id); err != nil {
			return nil, err
		}
		form.Add("fileId[]", strconv.Itoa(id))
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "sharedFiles")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Unlink removes a shared file link from the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-wiki
func (s *WikiService) Unlink(ctx context.Context, wikiID, fileID int) (*model.SharedFile, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}
	if err := validate.ValidateSharedFileID(fileID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "sharedFiles", strconv.Itoa(fileID))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.SharedFile{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewWikiService creates and returns a new shared-file WikiService.
func NewWikiService(method *core.Method) *WikiService {
	return &WikiService{method: method}
}
