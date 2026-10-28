package project

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// SharedFileService handles project shared-file-related Backlog API calls.
// Kept separate from internal/sharedfile because Download is project-specific
// and does not fit the spath-agnostic pattern used by that package.
type SharedFileService struct {
	method *client.Method
}

var sharedFileListValidTypes = []option.APIParamOptionType{
	option.ParamOrder,
	option.ParamOffset,
	option.ParamCount,
}

// List returns a list of shared files in the directory of the project.
//
// dirPath is the directory path shown in the project's file tree, e.g. "/design/".
// Use "/" to list the project root. The path is always sent with a trailing
// slash, so "/design" and "/design/" are equivalent.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files
func (s *SharedFileService) List(ctx context.Context, projectIDOrKey, dirPath string, opts ...*option.APIParamOption) ([]*model.SharedFile, error) {
	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateDirPath(dirPath); ve != nil {
		ves = append(ves, ve)
	}

	query := url.Values{}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, sharedFileListValidTypes, opts...)); err != nil {
		return nil, err
	}

	// The directory is sent in its directory form, with a trailing slash:
	// "/design/" -> files/metadata/design/, "/" -> files/metadata/.
	// Each directory name becomes its own segment; the empty segments
	// produced by leading, trailing or repeated slashes are dropped.
	segments := append([]string{"projects", projectIDOrKey, "files", "metadata"}, splitDirPath(dirPath)...)
	spath := client.PathJoin(segments...) + "/"
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.SharedFile{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Download downloads a shared file from the project.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-file
func (s *SharedFileService) Download(ctx context.Context, projectIDOrKey string, sharedFileID int) (*model.FileData, error) {
	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateSharedFileID(sharedFileID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := client.PathJoin("projects", projectIDOrKey, "files", strconv.Itoa(sharedFileID))
	resp, err := s.method.Download(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	return client.DownloadResponse(resp)
}

// splitDirPath splits dirPath into directory names, dropping empty segments.
func splitDirPath(dirPath string) []string {
	return strings.FieldsFunc(dirPath, func(r rune) bool { return r == '/' })
}

func NewSharedFileService(method *client.Method) *SharedFileService {
	return &SharedFileService{method: method}
}
