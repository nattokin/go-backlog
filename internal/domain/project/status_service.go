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

// StatusService handles status-related Backlog API calls for a project.
type StatusService struct {
	method *client.Method
}

// List returns a list of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-status-list-of-project
func (s *StatusService) List(ctx context.Context, projectIDOrKey string) ([]*model.Status, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "statuses")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Status{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create adds a new status to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-status
func (s *StatusService) Create(ctx context.Context, projectIDOrKey, name, color string) (*model.Status, error) {
	opt := &option.OptionService{}
	nameOpt := opt.WithName(name)
	colorOpt := opt.WithColor(color)

	var ves core.ValidationErrors
	if ve := nameOpt.Check(); ve != nil {
		ves = append(ves, ve)
	}
	if ve := colorOpt.Check(); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	nameOpt.Set(form)
	colorOpt.Set(form)

	spath := path.Join("projects", projectIDOrKey, "statuses")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Status{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a status in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-status
func (s *StatusService) Update(ctx context.Context, projectIDOrKey string, statusID int, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.Status, error) {
	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamName, option.ParamColor}
	options := append([]*option.APIParamOption{opt}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if statusID < 1 {
		ves = append(ves, core.NewValidationError("statusId", "statusId must not be less than 1"))
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "statuses", strconv.Itoa(statusID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Status{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a status from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-status
func (s *StatusService) Delete(ctx context.Context, projectIDOrKey string, statusID, substituteStatusID int) (*model.Status, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if statusID < 1 {
		ves = append(ves, core.NewValidationError("statusId", "statusId must not be less than 1"))
	}
	if substituteStatusID < 1 {
		ves = append(ves, core.NewValidationError("substituteStatusId", "substituteStatusId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	form.Set("substituteStatusId", strconv.Itoa(substituteStatusID))

	spath := path.Join("projects", projectIDOrKey, "statuses", strconv.Itoa(statusID))
	resp, err := s.method.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Status{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// UpdateOrder updates the display order of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-order-of-status
func (s *StatusService) UpdateOrder(ctx context.Context, projectIDOrKey string, statusIDs []int) ([]*model.Status, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(statusIDs) == 0 {
		ves = append(ves, core.NewValidationError("statusIDs", "statusIDs must not be empty"))
	}
	for _, id := range statusIDs {
		if id < 1 {
			ves = append(ves, core.NewValidationError("statusId", "each statusId must not be less than 1"))
			break
		}
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	for _, id := range statusIDs {
		form.Add("statusId[]", strconv.Itoa(id))
	}

	spath := path.Join("projects", projectIDOrKey, "statuses", "updateDisplayOrder")
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := []*model.Status{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func NewStatusService(method *client.Method) *StatusService {
	return &StatusService{method: method}
}
