package project

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com	/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

type StatusService struct {
	method *core.Method
}

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
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *StatusService) Create(ctx context.Context, projectIDOrKey, name, color string) (*model.Status, error) {
	opt := &core.OptionService{}
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
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *StatusService) Update(ctx context.Context, projectIDOrKey string, statusID int, option *core.APIParamOption, opts ...*core.APIParamOption) (*model.Status, error) {
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamColor}
	options := append([]*core.APIParamOption{option}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if statusID < 1 {
		ves = append(ves, core.NewValidationError("statusId", "statusId must not be less than 1"))
	}
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		var optVes core.ValidationErrors
		if !errors.As(err, &optVes) {
			return nil, err
		}
		ves = append(ves, optVes...)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "statuses", strconv.Itoa(statusID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Status{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

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
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

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
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func NewStatusService(method *core.Method) *StatusService {
	return &StatusService{method: method}
}
