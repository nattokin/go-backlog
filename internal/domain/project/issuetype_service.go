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

type IssueTypeService struct {
	method *client.Method
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
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create adds a new issue type to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-issue-type
func (s *IssueTypeService) Create(ctx context.Context, projectIDOrKey, name, color string, opts ...*option.APIParamOption) (*model.IssueType, error) {
	optSvc := &option.OptionService{}
	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamName, option.ParamColor, option.ParamTemplateSummary, option.ParamTemplateDescription}
	options := append([]*option.APIParamOption{optSvc.WithName(name), optSvc.WithColor(color)}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "issueTypes")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.IssueType{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates an issue type in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-issue-type
func (s *IssueTypeService) Update(ctx context.Context, projectIDOrKey string, issueTypeID int, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.IssueType, error) {
	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamName, option.ParamColor, option.ParamTemplateSummary, option.ParamTemplateDescription}
	options := append([]*option.APIParamOption{opt}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if issueTypeID < 1 {
		ves = append(ves, core.NewValidationError("issueTypeId", "issueTypeId must not be less than 1"))
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "issueTypes", strconv.Itoa(issueTypeID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.IssueType{}
	if err := client.DecodeResponse(resp, &v); err != nil {
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
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewIssueTypeService(method *client.Method) *IssueTypeService {
	return &IssueTypeService{method: method}
}
