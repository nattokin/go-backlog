package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
)

// ──────────────────────────────────────────────────────────────
//  CustomField models
// ──────────────────────────────────────────────────────────────

// CustomField represents a custom field defined in the project.
//
// The type-specific settings are set only for the matching custom field type
// and hold their zero value otherwise:
//   - Min, Max, InitialValue, Unit: number type
//   - MinDate, MaxDate, InitialDate: date type
//   - AllowInput: list types
//
// Min, Max and InitialValue are pointers because 0 is a valid setting.
// MinDate, MaxDate and InitialDate.Date hold midnight of the date in the space's time zone,
// expressed in UTC, in the same way as the StartDate field of [Version].
type CustomField struct {
	ID                     int
	TypeID                 int
	Name                   string
	Description            string
	Required               bool
	ApplicableIssueTypeIDs []int
	AllowAddItem           bool
	Items                  []*CustomFieldItem

	Min          *float64
	Max          *float64
	InitialValue *float64
	Unit         string
	MinDate      Timestamp
	MaxDate      Timestamp
	InitialDate  *CustomFieldInitialDate
	AllowInput   bool
}

// CustomFieldInitialDate represents how the initial value of a date type [CustomField] is decided.
//
// InitialValueType is 1 for today, 2 for today plus Shift days, and 3 for Date,
// matching the value set by WithInitialValueType.
type CustomFieldInitialDate struct {
	InitialValueType int
	Shift            int
	Date             Timestamp
}

// CustomFieldItem represents one item in a list type [CustomField].
type CustomFieldItem struct {
	ID           int
	Name         string
	DisplayOrder int
}

// ──────────────────────────────────────────────────────────────
//  ProjectCustomFieldService
// ──────────────────────────────────────────────────────────────

// ProjectCustomFieldService handles communication with the project custom-field-related
// methods of the Backlog API.
type ProjectCustomFieldService struct {
	base *project.CustomFieldService
	// Option provides the options accepted by [ProjectCustomFieldService] methods.
	Option *ProjectCustomFieldOptionService
}

// List returns a list of custom fields in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-custom-field-list
func (s *ProjectCustomFieldService) List(ctx context.Context, projectIDOrKey string) ([]*CustomField, error) {
	v, err := s.base.List(ctx, projectIDOrKey)
	return customFieldsFromModel(v), convertError(err)
}

// Create adds a new custom field to a project.
//
// fieldType specifies the custom field type. Use the CustomFieldType constants:
//   - [CustomFieldTypeText]
//   - [CustomFieldTypeSentence]
//   - [CustomFieldTypeNumber]
//   - [CustomFieldTypeDate]
//   - [CustomFieldTypeSingleList]
//   - [CustomFieldTypeMultipleList]
//   - [CustomFieldTypeCheckbox]
//   - [CustomFieldTypeRadio]
//
// This method supports options returned by methods in [ProjectCustomFieldOptionService],
// such as:
//   - [ProjectCustomFieldOptionService.WithDescription]
//   - [ProjectCustomFieldOptionService.WithRequired]
//   - [ProjectCustomFieldOptionService.WithApplicableIssueTypeIDs]
//   - [ProjectCustomFieldOptionService.WithMin] (Number type)
//   - [ProjectCustomFieldOptionService.WithMax] (Number type)
//   - [ProjectCustomFieldOptionService.WithInitialValue] (Number type)
//   - [ProjectCustomFieldOptionService.WithUnit] (Number type)
//   - [ProjectCustomFieldOptionService.WithInitialValueType] (Date type)
//   - [ProjectCustomFieldOptionService.WithInitialDate] (Date type)
//   - [ProjectCustomFieldOptionService.WithInitialShift] (Date type)
//   - [ProjectCustomFieldOptionService.WithMinDate] (Date type)
//   - [ProjectCustomFieldOptionService.WithMaxDate] (Date type)
//   - [ProjectCustomFieldOptionService.WithItems] (List types)
//   - [ProjectCustomFieldOptionService.WithAllowInput] (List types)
//   - [ProjectCustomFieldOptionService.WithAllowAddItem] (List types)
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-custom-field
func (s *ProjectCustomFieldService) Create(ctx context.Context, projectIDOrKey string, fieldType CustomFieldType, name string, opts ...RequestOption) (*CustomField, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, int(fieldType), name, toInnerOptions(opts)...)
	return customFieldFromModel(v), convertError(err)
}

// Update updates a custom field in a project.
//
// At least one option must be provided. This method supports options returned
// by methods in [ProjectCustomFieldOptionService], such as:
//   - [ProjectCustomFieldOptionService.WithName]
//   - [ProjectCustomFieldOptionService.WithDescription]
//   - [ProjectCustomFieldOptionService.WithRequired]
//   - [ProjectCustomFieldOptionService.WithApplicableIssueTypeIDs]
//   - [ProjectCustomFieldOptionService.WithMin] (Number type)
//   - [ProjectCustomFieldOptionService.WithMax] (Number type)
//   - [ProjectCustomFieldOptionService.WithInitialValue] (Number type)
//   - [ProjectCustomFieldOptionService.WithUnit] (Number type)
//   - [ProjectCustomFieldOptionService.WithInitialValueType] (Date type)
//   - [ProjectCustomFieldOptionService.WithInitialDate] (Date type)
//   - [ProjectCustomFieldOptionService.WithInitialShift] (Date type)
//   - [ProjectCustomFieldOptionService.WithMinDate] (Date type)
//   - [ProjectCustomFieldOptionService.WithMaxDate] (Date type)
//   - [ProjectCustomFieldOptionService.WithItems] (List types)
//   - [ProjectCustomFieldOptionService.WithAllowInput] (List types)
//   - [ProjectCustomFieldOptionService.WithAllowAddItem] (List types)
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-custom-field
func (s *ProjectCustomFieldService) Update(ctx context.Context, projectIDOrKey string, customFieldID int, option RequestOption, opts ...RequestOption) (*CustomField, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, customFieldID, toInnerOption(option), toInnerOptions(opts)...)
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
	base *option.OptionService
}

// WithName sets the custom field name.
func (s *ProjectCustomFieldOptionService) WithName(name string) RequestOption {
	return &requestOption{opt: s.base.WithName(name)}
}

// WithDescription sets the custom field description.
func (s *ProjectCustomFieldOptionService) WithDescription(description string) RequestOption {
	return &requestOption{opt: s.base.WithDescription(description)}
}

// WithRequired sets whether the custom field is required.
func (s *ProjectCustomFieldOptionService) WithRequired(required bool) RequestOption {
	return &requestOption{opt: s.base.WithRequired(required)}
}

// WithApplicableIssueTypeIDs sets the issue type IDs to which the custom field applies.
func (s *ProjectCustomFieldOptionService) WithApplicableIssueTypeIDs(ids []int) RequestOption {
	return &requestOption{opt: s.base.WithApplicableIssueTypeIDs(ids)}
}

// WithMin sets the minimum value for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithMin(min float64) RequestOption {
	return &requestOption{opt: s.base.WithMin(min)}
}

// WithMax sets the maximum value for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithMax(max float64) RequestOption {
	return &requestOption{opt: s.base.WithMax(max)}
}

// WithInitialValue sets the initial value for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithInitialValue(value float64) RequestOption {
	return &requestOption{opt: s.base.WithInitialValue(value)}
}

// WithUnit sets the unit string for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithUnit(unit string) RequestOption {
	return &requestOption{opt: s.base.WithUnit(unit)}
}

// WithInitialValueType sets the initial value type for a date type custom field.
// 1: today, 2: today+N days, 3: specified date.
func (s *ProjectCustomFieldOptionService) WithInitialValueType(initialValueType int) RequestOption {
	return &requestOption{opt: s.base.WithInitialValueType(initialValueType)}
}

// WithInitialDate sets the initial date ("yyyy-MM-dd") for a date type custom field
// when InitialValueType is 3.
func (s *ProjectCustomFieldOptionService) WithInitialDate(date string) RequestOption {
	return &requestOption{opt: s.base.WithInitialDate(date)}
}

// WithInitialShift sets the number of days to shift from today for a date type custom field
// when InitialValueType is 2.
func (s *ProjectCustomFieldOptionService) WithInitialShift(days int) RequestOption {
	return &requestOption{opt: s.base.WithInitialShift(days)}
}

// WithMinDate sets the minimum date ("yyyy-MM-dd") for a date type custom field.
// Use WithMin for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithMinDate(date string) RequestOption {
	return &requestOption{opt: s.base.WithMinDate(date)}
}

// WithMaxDate sets the maximum date ("yyyy-MM-dd") for a date type custom field.
// Use WithMax for a number type custom field.
func (s *ProjectCustomFieldOptionService) WithMaxDate(date string) RequestOption {
	return &requestOption{opt: s.base.WithMaxDate(date)}
}

// WithItems sets the list items for a list type custom field.
func (s *ProjectCustomFieldOptionService) WithItems(items []string) RequestOption {
	return &requestOption{opt: s.base.WithItems(items)}
}

// WithAllowInput sets whether free-text input is allowed for a list type custom field.
func (s *ProjectCustomFieldOptionService) WithAllowInput(enabled bool) RequestOption {
	return &requestOption{opt: s.base.WithAllowInput(enabled)}
}

// WithAllowAddItem sets whether users can add new items to a list type custom field.
func (s *ProjectCustomFieldOptionService) WithAllowAddItem(enabled bool) RequestOption {
	return &requestOption{opt: s.base.WithAllowAddItem(enabled)}
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectCustomFieldService(method *client.Method, option *option.OptionService) *ProjectCustomFieldService {
	return &ProjectCustomFieldService{
		base:   project.NewCustomFieldService(method),
		Option: newProjectCustomFieldOptionService(option),
	}
}

func newProjectCustomFieldOptionService(option *option.OptionService) *ProjectCustomFieldOptionService {
	return &ProjectCustomFieldOptionService{base: option}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func customFieldItemFromModel(m *model.CustomFieldItem) *CustomFieldItem {
	if m == nil {
		return nil
	}
	return &CustomFieldItem{
		ID:           m.ID,
		Name:         m.Name,
		DisplayOrder: m.DisplayOrder,
	}
}

func customFieldFromModel(m *model.CustomField) *CustomField {
	if m == nil {
		return nil
	}
	items := make([]*CustomFieldItem, len(m.Items))
	for i, v := range m.Items {
		items[i] = customFieldItemFromModel(v)
	}
	return &CustomField{
		ID:                     m.ID,
		TypeID:                 m.TypeID,
		Name:                   m.Name,
		Description:            m.Description,
		Required:               m.Required,
		ApplicableIssueTypeIDs: m.ApplicableIssueTypeIDs,
		AllowAddItem:           m.AllowAddItem,
		Items:                  items,

		Min:          numberFromCustomFieldBound(m.Min),
		Max:          numberFromCustomFieldBound(m.Max),
		InitialValue: m.InitialValue,
		Unit:         m.Unit,
		MinDate:      dateFromCustomFieldBound(m.Min),
		MaxDate:      dateFromCustomFieldBound(m.Max),
		InitialDate:  customFieldInitialDateFromModel(m.InitialDate),
		AllowInput:   m.AllowInput,
	}
}

func numberFromCustomFieldBound(b *model.CustomFieldBound) *float64 {
	if b == nil {
		return nil
	}
	return b.Number
}

func dateFromCustomFieldBound(b *model.CustomFieldBound) Timestamp {
	if b == nil || b.Date == nil {
		return Timestamp{}
	}
	return Timestamp{*b.Date}
}

func customFieldInitialDateFromModel(m *model.CustomFieldInitialDate) *CustomFieldInitialDate {
	if m == nil {
		return nil
	}
	return &CustomFieldInitialDate{
		InitialValueType: m.ID,
		Shift:            m.Shift,
		Date:             Timestamp{m.Date},
	}
}

func customFieldsFromModel(m []*model.CustomField) []*CustomField {
	if m == nil {
		return nil
	}
	result := make([]*CustomField, len(m))
	for i, v := range m {
		result[i] = customFieldFromModel(v)
	}
	return result
}
