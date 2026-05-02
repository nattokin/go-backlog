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

// ──────────────────────────────────────────────────────────────
//  CustomFieldService
// ──────────────────────────────────────────────────────────────

// CustomFieldService handles communication with the custom-field-related
// methods of the Backlog API.
type CustomFieldService struct {
	method *core.Method
}

// All returns a list of custom fields in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-custom-field-list
func (s *CustomFieldService) All(ctx context.Context, projectIDOrKey string) ([]*model.CustomField, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
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
// typeID specifies the custom field type:
//   - 1: Text
//   - 2: Sentence
//   - 3: Number
//   - 4: Date
//   - 5: Single List
//   - 6: Multiple List
//   - 7: Checkbox
//   - 8: Radio
//
// This method supports the following options:
//   - WithDescription
//   - WithRequired
//   - WithApplicableIssueTypeIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-custom-field
func (s *CustomFieldService) Create(ctx context.Context, projectIDOrKey string, typeID int, name string, opts ...core.RequestOption) (*model.CustomField, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if typeID < 1 {
		return nil, core.NewValidationError("typeId must not be less than 1")
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamTypeID, core.ParamName,
		core.ParamDescription, core.ParamRequired, core.ParamApplicableIssueTypeIDs,
	}
	options := append([]core.RequestOption{option.WithTypeID(typeID), option.WithName(name)}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
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
// This method supports the following options:
//   - WithName
//   - WithDescription
//   - WithRequired
//   - WithApplicableIssueTypeIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-custom-field
func (s *CustomFieldService) Update(ctx context.Context, projectIDOrKey string, customFieldID int, opts ...core.RequestOption) (*model.CustomField, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if customFieldID < 1 {
		return nil, core.NewValidationError("customFieldId must not be less than 1")
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamName, core.ParamDescription,
		core.ParamRequired, core.ParamApplicableIssueTypeIDs,
	}
	if err := core.ApplyOptions(form, validTypes, opts...); err != nil {
		return nil, err
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
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if customFieldID < 1 {
		return nil, core.NewValidationError("customFieldId must not be less than 1")
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
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if customFieldID < 1 {
		return nil, core.NewValidationError("customFieldId must not be less than 1")
	}

	option := (&core.OptionService{}).WithName(name)
	if err := option.Check(); err != nil {
		return nil, err
	}
	form := url.Values{}
	if err := option.Set(form); err != nil {
		return nil, err
	}

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
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if customFieldID < 1 {
		return nil, core.NewValidationError("customFieldId must not be less than 1")
	}
	if itemID < 1 {
		return nil, core.NewValidationError("itemId must not be less than 1")
	}

	option := (&core.OptionService{}).WithName(name)
	if err := option.Check(); err != nil {
		return nil, err
	}
	form := url.Values{}
	if err := option.Set(form); err != nil {
		return nil, err
	}

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
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if customFieldID < 1 {
		return nil, core.NewValidationError("customFieldId must not be less than 1")
	}
	if itemID < 1 {
		return nil, core.NewValidationError("itemId must not be less than 1")
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

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

func NewCustomFieldService(method *core.Method) *CustomFieldService {
	return &CustomFieldService{method: method}
}
