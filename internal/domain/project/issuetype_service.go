package project

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

// IssueTypeService handles issue type-related Backlog API calls for a project.
type IssueTypeService struct {
	method *core.Method
}

// List returns a list of issue types in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-type-list
func (s *IssueTypeService) List(ctx context.Context, projectIDOrKey string) ([]*model.IssueType, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "issueTypes")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.IssueType{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create adds a new issue type to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-issue-type
func (s *IssueTypeService) Create(ctx context.Context, projectIDOrKey, name, color string, opts ...*core.APIParamOption) (*model.IssueType, error) {
	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamColor, core.ParamTemplateSummary, core.ParamTemplateDescription}
	options := append([]*core.APIParamOption{option.WithName(name), option.WithColor(color)}, opts...)

	var ves core.ValidationErrors
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		var optVes core.ValidationErrors
		if !errors.As(err, &optVes) {
			return nil, err
		}
		ves = append(ves, optVes...)
	}
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "issueTypes")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.IssueType{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates an issue type in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-issue-type
func (s *IssueTypeService) Update(ctx context.Context, projectIDOrKey string, issueTypeID int, option *core.APIParamOption, opts ...*core.APIParamOption) (*model.IssueType, error) {
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamColor, core.ParamTemplateSummary, core.ParamTemplateDescription}
	options := append([]*core.APIParamOption{option}, opts...)

	var ves core.ValidationErrors
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		var optVes core.ValidationErrors
		if !errors.As(err, &optVes) {
			return nil, err
		}
		ves = append(ves, optVes...)
	}
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if issueTypeID < 1 {
		ves = append(ves, core.NewValidationError("issueTypeId", "issueTypeId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "issueTypes", strconv.Itoa(issueTypeID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.IssueType{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes an issue type from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-type
func (s *IssueTypeService) Delete(ctx context.Context, projectIDOrKey string, issueTypeID, substituteIssueTypeID int) (*model.IssueType, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if issueTypeID < 1 {
		ves = append(ves, core.NewValidationError("issueTypeId", "issueTypeId must not be less than 1"))
	}
	if substituteIssueTypeID < 1 {
		ves = append(ves, core.NewValidationError("substituteIssueTypeId", "substituteIssueTypeId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	form.Set("substituteIssueTypeId", strconv.Itoa(substituteIssueTypeID))

	spath := path.Join("projects", projectIDOrKey, "issueTypes", strconv.Itoa(issueTypeID))
	resp, err := s.method.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.IssueType{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewIssueTypeService(method *core.Method) *IssueTypeService {
	return &IssueTypeService{method: method}
}
