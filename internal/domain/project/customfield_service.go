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

type CustomFieldService struct {
	method *core.Method
}

// List returns a list of custom fields in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-custom-field-list
func (s *CustomFieldService) List(ctx context.Context, projectIDOrKey string) ([]*model.CustomField, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "customFields")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.CustomField{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create adds a new custom field to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-custom-field
func (s *CustomFieldService) Create(ctx context.Context, projectIDOrKey string, fieldType int, name string, opts ...*core.APIParamOption) (*model.CustomField, error) {
	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamTypeID, core.ParamName,
		core.ParamDescription, core.ParamRequired, core.ParamApplicableIssueTypeIDs,
		core.ParamMin, core.ParamMax, core.ParamInitialValue, core.ParamUnit,
		core.ParamInitialValueType, core.ParamInitialDate, core.ParamInitialShift,
		core.ParamItems, core.ParamAllowInput, core.ParamAllowAddItem,
	}
	options := append([]*core.APIParamOption{option.WithFieldType(fieldType), option.WithName(name)}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
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

	spath := path.Join("projects", projectIDOrKey, "customFields")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.CustomField{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a custom field in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-custom-field
func (s *CustomFieldService) Update(ctx context.Context, projectIDOrKey string, customFieldID int, option *core.APIParamOption, opts ...*core.APIParamOption) (*model.CustomField, error) {
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamName, core.ParamDescription,
		core.ParamRequired, core.ParamApplicableIssueTypeIDs,
	}
	options := append([]*core.APIParamOption{option}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if customFieldID < 1 {
		ves = append(ves, core.NewValidationError("customFieldId", "customFieldId must not be less than 1"))
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

	spath := path.Join("projects", projectIDOrKey, "customFields", strconv.Itoa(customFieldID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.CustomField{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a custom field from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-custom-field
func (s *CustomFieldService) Delete(ctx context.Context, projectIDOrKey string, customFieldID int) (*model.CustomField, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if customFieldID < 1 {
		ves = append(ves, core.NewValidationError("customFieldId", "customFieldId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "customFields", strconv.Itoa(customFieldID))
	resp, err := s.method.Delete(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := model.CustomField{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// AddListItem adds a list item to a list type custom field.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-list-item-for-list-type-custom-field
func (s *CustomFieldService) AddListItem(ctx context.Context, projectIDOrKey string, customFieldID int, name string) (*model.CustomField, error) {
	option := (&core.OptionService{}).WithName(name)
	var ves core.ValidationErrors
	if ve := option.Check(); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if customFieldID < 1 {
		ves = append(ves, core.NewValidationError("customFieldId", "customFieldId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	option.Set(form)

	spath := path.Join("projects", projectIDOrKey, "customFields", strconv.Itoa(customFieldID), "items")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.CustomField{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// UpdateListItem updates a list item in a list type custom field.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-list-item-for-list-type-custom-field
func (s *CustomFieldService) UpdateListItem(ctx context.Context, projectIDOrKey string, customFieldID, itemID int, name string) (*model.CustomField, error) {
	option := (&core.OptionService{}).WithName(name)
	var ves core.ValidationErrors
	if ve := option.Check(); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if customFieldID < 1 {
		ves = append(ves, core.NewValidationError("customFieldId", "customFieldId must not be less than 1"))
	}
	if itemID < 1 {
		ves = append(ves, core.NewValidationError("itemId", "itemId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	option.Set(form)

	spath := path.Join("projects", projectIDOrKey, "customFields", strconv.Itoa(customFieldID), "items", strconv.Itoa(itemID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.CustomField{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// DeleteListItem deletes a list item from a list type custom field.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-list-item-for-list-type-custom-field
func (s *CustomFieldService) DeleteListItem(ctx context.Context, projectIDOrKey string, customFieldID, itemID int) (*model.CustomField, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if customFieldID < 1 {
		ves = append(ves, core.NewValidationError("customFieldId", "customFieldId must not be less than 1"))
	}
	if itemID < 1 {
		ves = append(ves, core.NewValidationError("itemId", "itemId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "customFields", strconv.Itoa(customFieldID), "items", strconv.Itoa(itemID))
	resp, err := s.method.Delete(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := model.CustomField{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewCustomFieldService(method *core.Method) *CustomFieldService {
	return &CustomFieldService{method: method}
}
