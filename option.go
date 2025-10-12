package backlog

import (
	"fmt"
	"strconv"
)

// optionSupport provides shared access to the query/form option builders.
// It is intended for internal composition by each XxxOptionService (e.g. WikiOptionService).
// This struct should be initialized by the Client when setting up each service.
type optionSupport struct {
	query *QueryOptionService
	form  *FormOptionService
}

//
// ──────────────────────────────────────────────────────────────
//  Option Types
// ──────────────────────────────────────────────────────────────
//

// --- QueryOption Types -------------------------------------------------------------

// queryType represents the distinct query parameter keys for Backlog API requests.
type queryType int

// Value returns the string representation of the query parameter key for the API request.
func (t queryType) Value() string {
	var m = map[queryType]string{
		queryActivityTypeIDs: "activityTypeId[]",
		queryAll:             "all",
		queryArchived:        "archived",
		queryCount:           "count",
		queryKey:             "key",
		queryKeyword:         "keyword",
		queryMaxID:           "maxId",
		queryMinID:           "minId",
		queryOrder:           "order",
	}

	return m[t]
}

// --- FormOption Types -------------------------------------------------------------

// formType represents the distinct form field keys available for Backlog API requests.
type formType int

// Value returns the string representation of the form field key for the API request.
func (t formType) Value() string {
	var m = map[formType]string{
		formArchived:                          "archived",
		formChartEnabled:                      "chartEnabled",
		formContent:                           "content",
		formKey:                               "key",
		formName:                              "name",
		formMailAddress:                       "mailAddress",
		formMailNotify:                        "mailNotify",
		formPassword:                          "password",
		formProjectLeaderCanEditProjectLeader: "projectLeaderCanEditProjectLeader",
		formRoleType:                          "roleType",
		formSubtaskingEnabled:                 "subtaskingEnabled",
		formTextFormattingRule:                "textFormattingRule",
		formUserID:                            "userId",
		formSendMail:                          "sendMail",
	}

	return m[t]
}

//
// ──────────────────────────────────────────────────────────────
//  RequestOption interface and shared function types
// ──────────────────────────────────────────────────────────────
//

// RequestOption defines a common interface for all option types (e.g., FormOption, QueryOption).
// It allows unified validation handling across different request-level options.
type RequestOption interface {
	Check() error
}

// optionCheckFunc defines a generic validation function used by all RequestOption implementations.
type optionCheckFunc func() error

// --- QueryOption -------------------------------------------------------------

// queryOptionFunc applies a query option's value to the request parameters.
type queryOptionFunc func(query *QueryParams) error

// QueryOption represents a single query parameter to be applied to a request.
type QueryOption struct {
	t         queryType       // The underlying query parameter key type
	checkFunc optionCheckFunc // The function that performs validation
	setFunc   queryOptionFunc // The function that sets the value into the query
}

// Check validates the QueryOption by executing its check function, if defined.
func (o *QueryOption) Check() error {
	if o.checkFunc != nil {
		return o.checkFunc()
	}
	return nil
}

// set executes the stored function to apply the option value into the query.
func (o *QueryOption) set(query *QueryParams) error {
	if o.setFunc == nil {
		return newValidationError("query option has no setter")
	}
	return o.setFunc(query)
}

// validate ensures that the QueryOption is allowed for the current API context.
func (o *QueryOption) validate(validTypes []queryType) error {
	for _, valid := range validTypes {
		if o.t == valid {
			return nil
		}
	}
	return newInvalidQueryOptionError(o.t, validTypes)
}

// --- FormOption --------------------------------------------------------------

// formOptionFunc applies a form option's value to the request form body.
type formOptionFunc func(form *FormParams) error

// FormOption represents a single form field to be applied to a request body.
type FormOption struct {
	t         formType        // The underlying form field type
	checkFunc optionCheckFunc // The function that performs validation
	setFunc   formOptionFunc  // The function that sets the value into the form
}

// Check validates the FormOption by executing its check function, if defined.
func (o *FormOption) Check() error {
	if o.checkFunc != nil {
		return o.checkFunc()
	}
	return nil
}

// set executes the stored function to apply the option value into the form.
func (o *FormOption) set(form *FormParams) error {
	if o.setFunc == nil {
		return newValidationError("form option has no setter")
	}
	return o.setFunc(form)
}

// validate ensures that the FormOption is allowed for the current API context.
func (o *FormOption) validate(validTypes []formType) error {
	for _, valid := range validTypes {
		if o.t == valid {
			return nil
		}
	}
	return newInvalidFormOptionError(o.t, validTypes)
}

//
// ──────────────────────────────────────────────────────────────
//  QueryOptionService
// ──────────────────────────────────────────────────────────────
//

// QueryOptionService provides builders for query options.
// Each XxxOptionService selectively exposes only the valid methods.
type QueryOptionService struct{}

// applyOptions validates and applies one or more query options to the given query parameters.
func (s *QueryOptionService) applyOptions(query *QueryParams, opts ...*QueryOption) error {
	for _, opt := range opts {
		if err := opt.Check(); err != nil {
			return err
		}
		if err := opt.set(query); err != nil {
			return err
		}
	}
	return nil
}

// --- Boolean options ------------------------------------------------------------

// WithAll returns an option to set the `all` query parameter.
func (s *QueryOptionService) WithAll(enabled bool) *QueryOption {
	return &QueryOption{
		t: queryAll,
		setFunc: func(query *QueryParams) error {
			query.Set(queryAll.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithArchived returns an option to set the `archived` query parameter.
func (s *QueryOptionService) WithArchived(archived bool) *QueryOption {
	return &QueryOption{
		t: queryArchived,
		setFunc: func(query *QueryParams) error {
			query.Set(queryArchived.Value(), strconv.FormatBool(archived))
			return nil
		},
	}
}

// --- Integer options ------------------------------------------------------------

// WithCount returns an option to set the `count` query parameter.
func (s *QueryOptionService) WithCount(count int) *QueryOption {
	return &QueryOption{
		t: queryCount,
		checkFunc: func() error {
			if count < 1 || 100 < count {
				return newValidationError("count must be between 1 and 100")
			}
			return nil
		},
		setFunc: func(query *QueryParams) error {
			query.Set(queryCount.Value(), strconv.Itoa(count))
			return nil
		},
	}
}

// WithMaxID returns an option to set the `maxId` query parameter.
func (s *QueryOptionService) WithMaxID(id int) *QueryOption {
	return &QueryOption{
		t: queryMaxID,
		checkFunc: func() error {
			if err := validateActivityID(id, "maxID"); err != nil {
				return err
			}
			return nil
		},
		setFunc: func(query *QueryParams) error {
			query.Set(queryMaxID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// WithMinID returns an option to set the `minId` query parameter.
func (s *QueryOptionService) WithMinID(id int) *QueryOption {
	return &QueryOption{
		t: queryMinID,
		checkFunc: func() error {
			if err := validateActivityID(id, "minID"); err != nil {
				return err
			}
			return nil
		},
		setFunc: func(query *QueryParams) error {
			query.Set(queryMinID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// --- String options ------------------------------------------------------------

// WithKeyword returns an option to set the `keyword` query parameter.
func (s *QueryOptionService) WithKeyword(keyword string) *QueryOption {
	return &QueryOption{
		t: queryKeyword,
		setFunc: func(query *QueryParams) error {
			query.Set(queryKeyword.Value(), keyword)
			return nil
		},
	}
}

// --- Enum or special options ----------------------------------------------------

// WithActivityTypeIDs returns an option to set multiple `activityTypeId[]` query parameters.
func (s *QueryOptionService) WithActivityTypeIDs(typeIDs []int) *QueryOption {
	return &QueryOption{
		t: queryActivityTypeIDs,
		checkFunc: func() error {
			for _, id := range typeIDs {
				if err := validateActivityID(id, "activityTypeIds"); err != nil {
					return err
				}
			}
			return nil
		},
		setFunc: func(query *QueryParams) error {
			for _, id := range typeIDs {
				query.Add(queryActivityTypeIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithOrder returns an option to set the `order` query parameter.
func (s *QueryOptionService) WithOrder(order Order) *QueryOption {
	return &QueryOption{
		t: queryOrder,
		checkFunc: func() error {
			if order != OrderAsc && order != OrderDesc {
				msg := fmt.Sprintf("order must be only '%s' or '%s'", string(OrderAsc), string(OrderDesc))
				return newValidationError(msg)
			}
			return nil
		},
		setFunc: func(query *QueryParams) error {
			query.Set(queryOrder.Value(), string(order))
			return nil
		},
	}
}

//
// ──────────────────────────────────────────────────────────────
//  FormOptionService
// ──────────────────────────────────────────────────────────────
//

// FormOptionService provides builders for form options.
// Each XxxOptionService selectively exposes only the valid subset.
type FormOptionService struct{}

// applyOptions validates and applies one or more form options to the given form.
// It returns the first error encountered from Check or set.
func (s *FormOptionService) applyOptions(form *FormParams, opts ...*FormOption) error {
	for _, opt := range opts {
		if err := opt.Check(); err != nil {
			return err
		}
		if err := opt.set(form); err != nil {
			return err
		}
	}
	return nil
}

// --- Boolean options ------------------------------------------------------------

// WithArchived returns a form option that sets the `archived` field (e.g., for Project).
func (*FormOptionService) WithArchived(enabled bool) *FormOption {
	return &FormOption{
		t: formArchived,
		setFunc: func(form *FormParams) error {
			form.Set(formArchived.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithChartEnabled returns a form option that sets the `chartEnabled` field (e.g., for Project).
func (*FormOptionService) WithChartEnabled(enabled bool) *FormOption {
	return &FormOption{
		t: formChartEnabled,
		setFunc: func(form *FormParams) error {
			form.Set(formChartEnabled.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithMailNotify returns a form option that sets the `mailNotify` field (e.g., for Wiki, Issue).
func (*FormOptionService) WithMailNotify(enabled bool) *FormOption {
	return &FormOption{
		t: formMailNotify,
		setFunc: func(form *FormParams) error {
			form.Set(formMailNotify.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithProjectLeaderCanEditProjectLeader returns a form option that sets the `projectLeaderCanEditProjectLeader` field (e.g., for Project).
func (*FormOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) *FormOption {
	return &FormOption{
		t: formProjectLeaderCanEditProjectLeader,
		setFunc: func(form *FormParams) error {
			form.Set(formProjectLeaderCanEditProjectLeader.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithSendMail creates a form option to specify whether to send
// an invitation email to the newly created user.
func (*FormOptionService) WithSendMail(enabled bool) *FormOption {
	return &FormOption{
		t: formSendMail,
		setFunc: func(form *FormParams) error {
			form.Set(formSendMail.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithSubtaskingEnabled returns a form option that sets the `subtaskingEnabled` field (e.g., for Project).
func (*FormOptionService) WithSubtaskingEnabled(enabled bool) *FormOption {
	return &FormOption{
		t: formSubtaskingEnabled,
		setFunc: func(form *FormParams) error {
			form.Set(formSubtaskingEnabled.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// --- Integer options ------------------------------------------------------------

// WithUserID creates a form option to set the user's ID (login name).
func (*FormOptionService) WithUserID(id int) *FormOption {
	return &FormOption{
		t: formUserID,
		checkFunc: func() error {
			if err := validateID(id, "userID"); err != nil {
				return err
			}
			return nil
		},
		setFunc: func(form *FormParams) error {
			form.Set(formUserID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// --- String options -------------------------------------------------------------

// WithContent returns a form option that sets the `content` field (e.g., for Wiki, Comment).
func (*FormOptionService) WithContent(content string) *FormOption {
	return &FormOption{
		t: formContent,
		checkFunc: func() error {
			if content == "" {
				return newValidationError("content must not be empty")
			}
			return nil
		},
		setFunc: func(form *FormParams) error {
			form.Set(formContent.Value(), content)
			return nil
		},
	}
}

// WithKey returns a form option that sets the `key` field (e.g., for Project Key).
func (*FormOptionService) WithKey(key string) *FormOption {
	return &FormOption{
		t: formKey,
		checkFunc: func() error {
			if key == "" {
				return newValidationError("key must not be empty")
			}
			return nil
		},
		setFunc: func(form *FormParams) error {
			form.Set(formKey.Value(), key)
			return nil
		},
	}
}

// WithMailAddress returns a form option that sets the `mailAddress` field (e.g., for User).
func (*FormOptionService) WithMailAddress(mailAddress string) *FormOption {
	// ToDo: validate mailAddress (Note: The validation remains as simple not-empty check)
	return &FormOption{
		t: formMailAddress,
		checkFunc: func() error {
			if mailAddress == "" {
				return newValidationError("mailAddress must not be empty")
			}
			return nil
		},
		setFunc: func(form *FormParams) error {
			form.Set(formMailAddress.Value(), mailAddress)
			return nil
		},
	}
}

// WithName returns a form option that sets the `name` field (e.g., for Project Name, Wiki Name).
func (*FormOptionService) WithName(name string) *FormOption {
	return &FormOption{
		t: formName,
		checkFunc: func() error {
			if name == "" {
				return newValidationError("name must not be empty")
			}
			return nil
		},
		setFunc: func(form *FormParams) error {
			form.Set(formName.Value(), name)
			return nil
		},
	}
}

// WithPassword returns a form option that sets the `password` field (e.g., for User).
// It validates that the password meets the minimum length requirement (8 characters).
func (*FormOptionService) WithPassword(password string) *FormOption {
	return &FormOption{
		t: formPassword,
		checkFunc: func() error {
			if len(password) < 8 {
				return newValidationError("password must be at least 8 characters long")
			}
			return nil
		},
		setFunc: func(form *FormParams) error {
			form.Set(formPassword.Value(), password)
			return nil
		},
	}
}

// --- Enum or Special options ----------------------------------------------------

// WithRoleType returns a form option that sets the `roleType` field (e.g., for User).
func (*FormOptionService) WithRoleType(roleType Role) *FormOption {
	return &FormOption{
		t: formRoleType,
		checkFunc: func() error {
			if roleType < 1 || 6 < roleType {
				return newValidationError("roleType must be between 1 and 6")
			}
			return nil
		},
		setFunc: func(form *FormParams) error {
			form.Set(formRoleType.Value(), strconv.Itoa(int(roleType)))
			return nil
		},
	}
}

// WithTextFormattingRule returns a form option that sets the `textFormattingRule` field (e.g., for Project).
func (*FormOptionService) WithTextFormattingRule(format Format) *FormOption {
	return &FormOption{
		t: formTextFormattingRule,
		checkFunc: func() error {
			if format != FormatBacklog && format != FormatMarkdown {
				msg := fmt.Sprintf("format must be only '%s' or '%s'", string(FormatBacklog), string(FormatMarkdown))
				return newValidationError(msg)
			}
			return nil
		},
		setFunc: func(form *FormParams) error {
			form.Set(formTextFormattingRule.Value(), string(format))
			return nil
		},
	}
}

//
// ──────────────────────────────────────────────────────────────
//  ActivityOptionService
// ──────────────────────────────────────────────────────────────
//

// ActivityOptionService provides a domain-specific set of option builders
// for operations within the ActivityService.
//
// It exposes only those options that are valid for the Backlog Activity API.
// Internally, it delegates to an optionSupport instance,
// which holds the generic QueryOptionService (Form options are not used here).
type ActivityOptionService struct {
	support *optionSupport
}

// --- Query Options -----------------------------------------------------------

// WithQueryActivityTypeIDs returns a query option to filter activities
// by one or more activity type IDs.
//
// Example: typeIds=1&typeIds=2&typeIds=3
func (s *ActivityOptionService) WithQueryActivityTypeIDs(typeIDs []int) *QueryOption {
	return s.support.query.WithActivityTypeIDs(typeIDs)
}

// WithQueryMinID returns a query option to filter activities
// that have IDs greater than or equal to the given value.
func (s *ActivityOptionService) WithQueryMinID(id int) *QueryOption {
	return s.support.query.WithMinID(id)
}

// WithQueryMaxID returns a query option to filter activities
// that have IDs less than or equal to the given value.
func (s *ActivityOptionService) WithQueryMaxID(id int) *QueryOption {
	return s.support.query.WithMaxID(id)
}

// WithQueryCount returns a query option to limit the number of activities
// returned by the API (max 100).
func (s *ActivityOptionService) WithQueryCount(count int) *QueryOption {
	return s.support.query.WithCount(count)
}

// WithQueryOrder returns a query option to specify the result order (asc or desc).
func (s *ActivityOptionService) WithQueryOrder(order Order) *QueryOption {
	return s.support.query.WithOrder(order)
}

// --- Form Options -----------------------------------------------------------

// (Backlog Activity API currently does not use form parameters extensively.)

//
// ──────────────────────────────────────────────────────────────
//  ProjectOptionService
// ──────────────────────────────────────────────────────────────
//

// ProjectOptionService provides a domain-specific set of option builders
// for operations within the ProjectService.
//
// It exposes only those options that are valid for the Backlog Project API.
// Internally, it delegates to an optionSupport instance,
// which holds the generic FormOptionService and QueryOptionService.
type ProjectOptionService struct {
	support *optionSupport
}

// --- Query Options -----------------------------------------------------------

// WithQueryAll returns a query option that includes archived projects
// in the result list.
func (s *ProjectOptionService) WithQueryAll(enabled bool) *QueryOption {
	return s.support.query.WithAll(enabled)
}

// WithQueryArchived returns a query option that filters projects
// by their archived status.
func (s *ProjectOptionService) WithQueryArchived(enabled bool) *QueryOption {
	return s.support.query.WithArchived(enabled)
}

// --- Form Options ------------------------------------------------------------

// WithFormArchived returns a form option to set the project's archived status.
func (s *ProjectOptionService) WithFormArchived(enabled bool) *FormOption {
	return s.support.form.WithArchived(enabled)
}

// WithFormChartEnabled returns a form option to enable or disable charts
// for the project.
func (s *ProjectOptionService) WithFormChartEnabled(enabled bool) *FormOption {
	return s.support.form.WithChartEnabled(enabled)
}

// WithFormKey returns a form option to set the project's key.
func (s *ProjectOptionService) WithFormKey(key string) *FormOption {
	return s.support.form.WithKey(key)
}

// WithFormName returns a form option to set the project's name.
func (s *ProjectOptionService) WithFormName(name string) *FormOption {
	return s.support.form.WithName(name)
}

// WithFormProjectLeaderCanEditProjectLeader returns a form option to set
// whether the project leader can edit the project leader field.
func (s *ProjectOptionService) WithFormProjectLeaderCanEditProjectLeader(enabled bool) *FormOption {
	return s.support.form.WithProjectLeaderCanEditProjectLeader(enabled)
}

// WithFormSubtaskingEnabled returns a form option to enable or disable
// subtasking for the project.
func (s *ProjectOptionService) WithFormSubtaskingEnabled(enabled bool) *FormOption {
	return s.support.form.WithSubtaskingEnabled(enabled)
}

// WithFormTextFormattingRule returns a form option to set the text formatting
// rule (Backlog or Markdown) for the project.
func (s *ProjectOptionService) WithFormTextFormattingRule(format Format) *FormOption {
	return s.support.form.WithTextFormattingRule(format)
}

//
// ──────────────────────────────────────────────────────────────
//  UserOptionService
// ──────────────────────────────────────────────────────────────
//

// UserOptionService provides a domain-specific set of option builders
// for operations within the UserService.
//
// It exposes only those options that are valid for the Backlog User API.
// Internally, it delegates to an optionSupport instance,
// which holds the generic FormOptionService and QueryOptionService.
type UserOptionService struct {
	support *optionSupport
}

// --- Query Options -----------------------------------------------------------

// (Backlog User API currently does not use query parameters extensively.)

// --- Form Options ------------------------------------------------------------

// WithFormMailAddress returns a form option to set the user's mail address.
func (s *UserOptionService) WithFormMailAddress(mail string) *FormOption {
	return s.support.form.WithMailAddress(mail)
}

// WithFormName returns a form option to set the user's display name.
func (s *UserOptionService) WithFormName(name string) *FormOption {
	return s.support.form.WithName(name)
}

// WithFormPassword returns a form option to set the user's password.
func (s *UserOptionService) WithFormPassword(password string) *FormOption {
	return s.support.form.WithPassword(password)
}

// WithFormRoleType returns a form option to set the user's role type.
func (s *UserOptionService) WithFormRoleType(role Role) *FormOption {
	return s.support.form.WithRoleType(role)
}

// WithFormSendMail returns a form option to specify whether to send
// an invitation email to the newly created user.
func (s *UserOptionService) WithFormSendMail(enabled bool) *FormOption {
	return s.support.form.WithSendMail(enabled)
}

// WithFormUserID returns a form option to set the user's ID (login name).
func (s *UserOptionService) WithFormUserID(id int) *FormOption {
	return s.support.form.WithUserID(id)
}

//
// ──────────────────────────────────────────────────────────────
//  WikiOptionService
// ──────────────────────────────────────────────────────────────
//

// WikiOptionService provides a domain-specific set of option builders
// for operations within the WikiService.
//
// It exposes only those options that are valid for the Backlog Wiki API.
// Internally, it delegates to an optionSupport instance,
// which holds the generic FormOptionService and QueryOptionService.
type WikiOptionService struct {
	support *optionSupport
}

// --- Query Options -----------------------------------------------------------

// WithQueryKeyword returns a query option to search wiki pages by keyword.
func (s *WikiOptionService) WithQueryKeyword(keyword string) *QueryOption {
	return s.support.query.WithKeyword(keyword)
}

// --- Form Options ------------------------------------------------------------

// WithFormContent returns a form option to set the wiki page content.
func (s *WikiOptionService) WithFormContent(content string) *FormOption {
	return s.support.form.WithContent(content)
}

// WithFormMailNotify returns a form option to enable or disable mail notifications
// for wiki updates or deletions.
func (s *WikiOptionService) WithFormMailNotify(enabled bool) *FormOption {
	return s.support.form.WithMailNotify(enabled)
}

// WithFormName returns a form option to set the wiki page name.
func (s *WikiOptionService) WithFormName(name string) *FormOption {
	return s.support.form.WithName(name)
}

//
// ──────────────────────────────────────────────────────────────
//  Internal Helpers
// ──────────────────────────────────────────────────────────────
//

func validateID(id int, key string) error {
	if id < 1 {
		return newValidationError(fmt.Sprintf("invalid %s: must not be less than 1", key))
	}

	return nil
}

// validateActivityID ensures that the given activity ID is within the valid range [1, 26].
// If the ID is out of range, it returns a validation error indicating the parameter key.
//
// This function is used internally to validate activity-related query parameters
// (e.g., "activityTypeId[]") according to Backlog API constraints.
func validateActivityID(id int, key string) error {
	if id < 1 || id > 26 {
		return newValidationError(fmt.Sprintf("invalid %s: must be between 1 and 26", key))
	}
	return nil
}

// hasRequiredFormOption checks whether the provided form options
// include at least one of the required form types.
func hasRequiredFormOption(options []*FormOption, requiredTypes []formType) bool {
	for _, opt := range options {
		optionType := opt.t
		for _, requiredType := range requiredTypes {
			if optionType == requiredType {
				return true
			}
		}
	}
	return false
}
