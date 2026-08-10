// Package version implements the Backlog Version/Milestone API service.
package project

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// VersionService provides Version/Milestone API operations.
type VersionService struct {
	method *core.Method
}

// List returns versions/milestones in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-version-milestone-list
func (s *VersionService) List(ctx context.Context, projectIDOrKey string, opts ...*core.APIParamOption) ([]*model.Version, error) {
	query := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamArchived,
		core.ParamAll,
	}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := core.MergeValidationErrors(ves, core.ApplyOptions(query, validTypes, opts...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Version{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Add adds a version/milestone to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-version-milestone
func (s *VersionService) Add(ctx context.Context, projectIDOrKey, name string, opts ...*core.APIParamOption) (*model.Version, error) {
	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamStartDate,
		core.ParamReleaseDueDate,
	}
	options := append([]*core.APIParamOption{option.WithName(name)}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := core.MergeValidationErrors(ves, core.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Version{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a version/milestone.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-version-milestone
func (s *VersionService) Update(ctx context.Context, projectIDOrKey string, versionID int, option *core.APIParamOption, opts ...*core.APIParamOption) (*model.Version, error) {
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamStartDate,
		core.ParamReleaseDueDate,
		core.ParamArchived,
	}
	options := append([]*core.APIParamOption{option}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateVersionID(versionID); ve != nil {
		ves = append(ves, ve)
	}
	if err := core.MergeValidationErrors(ves, core.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions", strconv.Itoa(versionID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Version{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a version/milestone.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-version-milestone
func (s *VersionService) Delete(ctx context.Context, projectIDOrKey string, versionID int) (*model.Version, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateVersionID(versionID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "versions", strconv.Itoa(versionID))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Version{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewVersionService(method *core.Method) *VersionService {
	return &VersionService{method: method}
}
