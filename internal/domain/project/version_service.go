// Package version implements the Backlog Version/Milestone API service.
package project

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
)

// VersionService provides Version/Milestone API operations.
type VersionService struct {
	method *client.Method
}

// List returns versions/milestones in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-version-milestone-list
func (s *VersionService) List(ctx context.Context, projectIDOrKey string, opts ...*option.APIParamOption) ([]*model.Version, error) {
	query := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamArchived,
		option.ParamAll,
	}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, validTypes, opts...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Version{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Add adds a version/milestone to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-version-milestone
func (s *VersionService) Add(ctx context.Context, projectIDOrKey, name string, opts ...*option.APIParamOption) (*model.Version, error) {
	optSvc := &option.OptionService{}
	form := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamName,
		option.ParamDescription,
		option.ParamStartDate,
		option.ParamReleaseDueDate,
	}
	options := append([]*option.APIParamOption{optSvc.WithName(name)}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Version{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a version/milestone.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-version-milestone
func (s *VersionService) Update(ctx context.Context, projectIDOrKey string, versionID int, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.Version, error) {
	form := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamName,
		option.ParamDescription,
		option.ParamStartDate,
		option.ParamReleaseDueDate,
		option.ParamArchived,
	}
	options := append([]*option.APIParamOption{opt}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateVersionID(versionID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions", strconv.Itoa(versionID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Version{}
	if err := client.DecodeResponse(resp, &v); err != nil {
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
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewVersionService(method *client.Method) *VersionService {
	return &VersionService{method: method}
}
