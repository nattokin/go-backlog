package project

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// SharedFileService handles communication with the project shared-file-related methods of the Backlog API.
// Kept here because GetFile is project-specific and does not fit the shared spath-agnostic pattern.
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

// GetFile downloads a shared file from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-file
func (s *SharedFileService) GetFile(ctx context.Context, projectIDOrKey string, sharedFileID int) (*model.FileData, error) {
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

// NewSharedFileService creates and returns a new project SharedFileService.
func NewSharedFileService(method *core.Method) *SharedFileService {
	return &SharedFileService{method: method}
}
