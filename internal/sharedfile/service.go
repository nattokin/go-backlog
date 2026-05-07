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

// Service holds shared HTTP logic for shared-file-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *core.Method
}

// List returns the list of shared files at spath.
func (s *Service) List(ctx context.Context, spath string) ([]*model.SharedFile, error) {
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

// Link links shared files at spath using the given fileIDs.
func (s *Service) Link(ctx context.Context, spath string, fileIDs []int) ([]*model.SharedFile, error) {
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

// Unlink removes a shared file link at spath.
func (s *Service) Unlink(ctx context.Context, spath string) (*model.SharedFile, error) {
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

// NewService creates and returns a new shared-file Service.
func NewService(method *core.Method) *Service {
	return &Service{method: method}
}

// ──────────────────────────────────────────────────────────────
//  ProjectService
// ──────────────────────────────────────────────────────────────

// ProjectService handles communication with the project shared-file-related methods of the Backlog API.
// Kept here because GetFile is project-specific and does not fit the shared spath-agnostic pattern.
type ProjectService struct {
	method *core.Method
}

// List returns a list of shared files in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files
func (s *ProjectService) List(ctx context.Context, projectIDOrKey string) ([]*model.SharedFile, error) {
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
func (s *ProjectService) GetFile(ctx context.Context, projectIDOrKey string, sharedFileID int) (*model.FileData, error) {
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

// NewProjectService creates and returns a new shared-file ProjectService.
func NewProjectService(method *core.Method) *ProjectService {
	return &ProjectService{method: method}
}
