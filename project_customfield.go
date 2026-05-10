package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/project"
)

// ──────────────────────────────────────────────────────────────
//  ProjectCustomFieldService
// ──────────────────────────────────────────────────────────────

// ProjectCustomFieldService handles communication with the project custom-field-related
// methods of the Backlog API.
type ProjectCustomFieldService struct {
	base   *project.CustomFieldService
	Option *ProjectCustomFieldOptionService
}

// All returns a list of custom fields in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-custom-field-list
func (s *ProjectCustomFieldService) All(ctx context.Context, projectIDOrKey string) ([]*CustomField, error) {
	v, err := s.base.All(ctx, projectIDOrKey)
	return customFieldsFromModel(v), convertError(err)
}

// Create adds a new custom field to a project.
//
// fieldType specifies the custom field type. Use the CustomFieldType constants:
//   - CustomFieldTypeText
//   - CustomFieldTypeSentence
//   - CustomFieldTypeNumber
//   - CustomFieldTypeDate
//   - CustomFieldTypeSingleList
//   - CustomFieldTypeMultipleList
//   - CustomFieldTypeCheckbox
//   - CustomFieldTypeRadio
//
// This method supports options returned by methods in "*Client.Project.CustomField.Option",
// such as:
//   - WithDescription
//   - WithRequired
//   - WithApplicableIssueTypeIDs
//   - WithMin, WithMax, WithInitialValue, WithUnit (Number type)
//   - WithInitialValueType, WithInitialDate, WithInitialShift (Date type)
//   - WithItems, WithAllowInput, WithAllowAddItem (List types)
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-custom-field
func (s *ProjectCustomFieldService) Create(ctx context.Context, projectIDOrKey string, fieldType CustomFieldType, name string, opts ...RequestOption) (*CustomField, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, int(fieldType), name, toCoreOptions(opts)...)
	return customFieldFromModel(v), convertError(err)
}

// Update updates a custom field in a project.
//
// At least one option must be provided. This method supports options returned
// by methods in "*Client.Project.CustomField.Option", such as:
//   - WithName
//   - WithDescription
//   - WithRequired
//   - WithApplicableIssueTypeIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-custom-field
func (s *ProjectCustomFieldService) Update(ctx context.Context, projectIDOrKey string, customFieldID int, opt RequestOption, opts ...RequestOption) (*CustomField, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, customFieldID, opt, toCoreOptions(opts)...)
	return customFieldFromModel(v), convertError(err)
}

// Delete deletes a custom field from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-custom-field
func (s *ProjectCustomFieldService) Delete(ctx context.Context, projectIDOrKey string, customFieldID int) (*CustomField, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, customFieldID)
	return customFieldFromModel(v), convertError(err)
}

// AddListItem adds a list item to a list type custom field.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-list-item-for-list-type-custom-field
func (s *ProjectCustomFieldService) AddListItem(ctx context.Context, projectIDOrKey string, customFieldID int, name string) (*CustomField, error) {
	v, err := s.base.AddListItem(ctx, projectIDOrKey, customFieldID, name)
	return customFieldFromModel(v), convertError(err)
}

// UpdateListItem updates a list item in a list type custom field.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-list-item-for-list-type-custom-field
func (s *ProjectCustomFieldService) UpdateListItem(ctx context.Context, projectIDOrKey string, customFieldID, itemID int, name string) (*CustomField, error) {
	v, err := s.base.UpdateListItem(ctx, projectIDOrKey, customFieldID, itemID, name)
	return customFieldFromModel(v), convertError(err)
}

// DeleteListItem deletes a list item from a list type custom field.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-list-item-for-list-type-custom-field
func (s *ProjectCustomFieldService) DeleteListItem(ctx context.Context, projectIDOrKey string, customFieldID, itemID int) (*CustomField, error) {
	v, err := s.base.DeleteListItem(ctx, projectIDOrKey, customFieldID, itemID)
	return customFieldFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectCustomFieldOptionService
// ──────────────────────────────────────────────────────────────

// ProjectCustomFieldOptionService provides a domain-specific set of option builders
// for operations within the ProjectCustomFieldService.
type ProjectCustomFieldOptionService struct {
	base *core.OptionService
}

// WithName sets the custom field name.
func (s *ProjectCustomFieldOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// WithDescription sets the custom field description.
func (s *ProjectCustomFieldOptionService) WithDescription(description string) RequestOption {
	return s.base.WithDescription(description)
}

// WithRequired sets whether the custom field is required.
func (s *ProjectCustomFieldOptionService) WithRequired(required bool) RequestOption {
	return s.base.WithRequired(required)
}

// WithApplicableIssueTypeIDs sets the issue type IDs to which the custom field applies.
func (s *ProjectCustomFieldOptionService) WithApplicableIssueTypeIDs(ids []int) RequestOption {
	return s.base.WithApplicableIssueTypeIDs(ids)
}

// WithMin sets the minimum value for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithMin(min float64) RequestOption {
	return s.base.WithMin(min)
}

// WithMax sets the maximum value for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithMax(max float64) RequestOption {
	return s.base.WithMax(max)
}

// WithInitialValue sets the initial value for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithInitialValue(value float64) RequestOption {
	return s.base.WithInitialValue(value)
}

// WithUnit sets the unit string for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithUnit(unit string) RequestOption {
	return s.base.WithUnit(unit)
}

// WithInitialValueType sets the initial value type for a date type custom field.
// 1: today, 2: today+N days, 3: specified date.
func (s *ProjectCustomFieldOptionService) WithInitialValueType(initialValueType int) RequestOption {
	return s.base.WithInitialValueType(initialValueType)
}

// WithInitialDate sets the initial date ("yyyy-MM-dd") for a date type custom field
// when InitialValueType is 3.
func (s *ProjectCustomFieldOptionService) WithInitialDate(date string) RequestOption {
	return s.base.WithInitialDate(date)
}

// WithInitialShift sets the number of days to shift from today for a date type custom field
// when InitialValueType is 2.
func (s *ProjectCustomFieldOptionService) WithInitialShift(days int) RequestOption {
	return s.base.WithInitialShift(days)
}

// WithItems sets the list items for a list type custom field.
func (s *ProjectCustomFieldOptionService) WithItems(items []string) RequestOption {
	return s.base.WithItems(items)
}

// WithAllowInput sets whether free-text input is allowed for a list type custom field.
func (s *ProjectCustomFieldOptionService) WithAllowInput(enabled bool) RequestOption {
	return s.base.WithAllowInput(enabled)
}

// WithAllowAddItem sets whether users can add new items to a list type custom field.
func (s *ProjectCustomFieldOptionService) WithAllowAddItem(enabled bool) RequestOption {
	return s.base.WithAllowAddItem(enabled)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectCustomFieldService(method *core.Method, option *core.OptionService) *ProjectCustomFieldService {
	return &ProjectCustomFieldService{
		base:   project.NewCustomFieldService(method),
		Option: newProjectCustomFieldOptionService(option),
	}
}

func newProjectCustomFieldOptionService(option *core.OptionService) *ProjectCustomFieldOptionService {
	return &ProjectCustomFieldOptionService{base: option}
}
