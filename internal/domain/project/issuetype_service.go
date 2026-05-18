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

// IssueTypeService handles issue type-related Backlog API calls for a project.
type IssueTypeService struct {
	method *core.Method
}

// List returns a list of issue types in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-type-list
func (s *IssueTypeService) List(ctx context.Context, projectIDOrKey string) ([]*model.IssueType, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
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
func (s *IssueTypeService) Create(ctx context.Context, projectIDOrKey, name, color string, opts ...core.RequestOption) (*model.IssueType, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamColor, core.ParamTemplateSummary, core.ParamTemplateDescription}
	options := append([]core.RequestOption{option.WithName(name), option.WithColor(color)}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
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
func (s *IssueTypeService) Update(ctx context.Context, projectIDOrKey string, issueTypeID int, option core.RequestOption, opts ...core.RequestOption) (*model.IssueType, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if issueTypeID < 1 {
		return nil, core.NewValidationError("issueTypeId must not be less than 1")
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamColor, core.ParamTemplateSummary, core.ParamTemplateDescription}
	options := append([]core.RequestOption{option}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
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
// substituteIssueTypeID specifies the issue type to migrate existing issues to.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-type
func (s *IssueTypeService) Delete(ctx context.Context, projectIDOrKey string, issueTypeID, substituteIssueTypeID int) (*model.IssueType, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if issueTypeID < 1 {
		return nil, core.NewValidationError("issueTypeId must not be less than 1")
	}
	if substituteIssueTypeID < 1 {
		return nil, core.NewValidationError("substituteIssueTypeId must not be less than 1")
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
