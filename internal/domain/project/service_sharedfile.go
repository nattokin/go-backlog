package project

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// SharedFileService handles project shared-file-related Backlog API calls.
// Kept separate from internal/sharedfile because GetFile is project-specific
// and does not fit the spath-agnostic pattern used by that package.
type SharedFileService struct {
	method *core.Method
}

// List returns a list of shared files in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files
func (s *SharedFileService) List(ctx context.Context, projectIDOrKey string) ([]*model.SharedFile, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "files")
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

// Download downloads a shared file from the project.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-file
func (s *SharedFileService) Download(ctx context.Context, projectIDOrKey string, sharedFileID int) (*model.FileData, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateSharedFileID(sharedFileID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "files", strconv.Itoa(sharedFileID))
	resp, err := s.method.Download(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	return core.DownloadResponse(resp)
}

func NewSharedFileService(method *core.Method) *SharedFileService {
	return &SharedFileService{method: method}
}
