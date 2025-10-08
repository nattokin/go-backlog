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
//  Core Option Types and Services
// ──────────────────────────────────────────────────────────────
//

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

// --- QueryOption -------------------------------------------------------------

// queryOptionFunc applies a query option's value to the request parameters.
type queryOptionFunc func(query *QueryParams) error

// QueryOption represents a single query parameter to be applied to a request.
type QueryOption struct {
	t queryType       // The underlying query parameter key type
	f queryOptionFunc // The function that sets the value into the query
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

// set executes the stored function to apply the option.
func (o *QueryOption) set(query *QueryParams) error {
	return o.f(query)
}

// QueryOptionService provides builders for query options.
// Each XxxOptionService selectively exposes only the valid methods.
type QueryOptionService struct{}

// WithActivityTypeIDs returns option to set `activityTypeId`.
func (s *QueryOptionService) WithActivityTypeIDs(typeIDs []int) *QueryOption {
	return &QueryOption{queryActivityTypeIDs, func(query *QueryParams) error {
		for _, id := range typeIDs {
			if id < 1 || 26 < id {
				return newValidationError("activityTypeId must be between 1 and 26")
			}
		}

		for _, id := range typeIDs {
			query.Add(queryActivityTypeIDs.Value(), strconv.Itoa(id))
		}
		return nil
	}}
}

// WithAll returns option to set `all`.
func (s *QueryOptionService) WithAll(enabeld bool) *QueryOption {
	return &QueryOption{queryAll, func(query *QueryParams) error {
		query.Set(queryAll.Value(), strconv.FormatBool(enabeld))
		return nil
	}}
}

// WithArchived returns option to set `archived`.
func (s *QueryOptionService) WithArchived(archived bool) *QueryOption {
	return &QueryOption{queryArchived, func(query *QueryParams) error {
		query.Set(queryArchived.Value(), strconv.FormatBool(archived))
		return nil
	}}
}

// WithCount returns option to set `count`.
func (s *QueryOptionService) WithCount(count int) *QueryOption {
	return &QueryOption{queryCount, func(query *QueryParams) error {
		if count < 1 || 100 < count {
			return newValidationError("count must be between 1 and 100")
		}
		query.Set(queryCount.Value(), strconv.Itoa(count))
		return nil
	}}
}

// WithKeyword returns option to set `keyword`.
func (s *QueryOptionService) WithKeyword(keyword string) *QueryOption {
	return &QueryOption{queryKeyword, func(query *QueryParams) error {
		query.Set(queryKeyword.Value(), keyword)
		return nil
	}}
}

// WithMaxID returns option to set `maxId`.
func (s *QueryOptionService) WithMaxID(maxID int) *QueryOption {
	return &QueryOption{queryMaxID, func(query *QueryParams) error {
		if maxID < 1 {
			return newValidationError("maxId must not be less than 1")
		}
		query.Set(queryMaxID.Value(), strconv.Itoa(maxID))
		return nil
	}}
}

// WithMinID returns option to set `minId`.
func (s *QueryOptionService) WithMinID(minID int) *QueryOption {
	return &QueryOption{queryMinID, func(query *QueryParams) error {
		if minID < 1 {
			return newValidationError("minId must not be less than 1")
		}
		query.Set(queryMinID.Value(), strconv.Itoa(minID))
		return nil
	}}
}

// WithOrder returns option to set `order`.
func (s *QueryOptionService) WithOrder(order Order) *QueryOption {
	return &QueryOption{queryOrder, func(query *QueryParams) error {
		if order != OrderAsc && order != OrderDesc {
			msg := fmt.Sprintf("order must be only '%s' or '%s'", string(OrderAsc), string(OrderDesc))
			return newValidationError(msg)
		}
		query.Set(queryOrder.Value(), string(order))
		return nil
	}}
}

// --- FormOption --------------------------------------------------------------

// formOptionFunc applies a form option's value to the request form body.
type formOptionFunc func(form *FormParams) error

// FormOption represents a single form field to be applied to a request body.
type FormOption struct {
	t formType       // The underlying form field type
	f formOptionFunc // The function that sets the value into the form
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

// set executes the stored function to apply the option.
func (o *FormOption) set(form *FormParams) error {
	return o.f(form)
}

// FormOptionService provides builders for form options.
// Each XxxOptionService selectively exposes only the valid subset.
type FormOptionService struct{}

// WithArchived returns a form option that sets the `archived` field (e.g., for Project).
func (*FormOptionService) WithArchived(enabled bool) *FormOption {
	return &FormOption{formArchived, func(form *FormParams) error {
		form.Set(formArchived.Value(), strconv.FormatBool(enabled))
		return nil
	}}
}

// WithChartEnabled returns a form option that sets the `chartEnabled` field (e.g., for Project).
func (*FormOptionService) WithChartEnabled(enabled bool) *FormOption {
	return &FormOption{formChartEnabled, func(form *FormParams) error {
		form.Set(formChartEnabled.Value(), strconv.FormatBool(enabled))
		return nil
	}}
}

// WithContent returns a form option that sets the `content` field (e.g., for Wiki, Comment).
func (*FormOptionService) WithContent(content string) *FormOption {
	return &FormOption{formContent, func(form *FormParams) error {
		if content == "" {
			return newValidationError("content must not be empty")
		}
		form.Set(formContent.Value(), content)
		return nil
	}}
}

// WithKey returns a form option that sets the `key` field (e.g., for Project Key).
func (*FormOptionService) WithKey(key string) *FormOption {
	return &FormOption{formKey, func(form *FormParams) error {
		if key == "" {
			return newValidationError("key must not be empty")
		}
		form.Set(formKey.Value(), key)
		return nil
	}}
}

// WithName returns a form option that sets the `name` field (e.g., for Project Name, Wiki Name).
func (*FormOptionService) WithName(name string) *FormOption {
	return &FormOption{formName, func(form *FormParams) error {
		if name == "" {
			return newValidationError("name must not be empty")
		}
		form.Set(formName.Value(), name)
		return nil
	}}
}

// WithMailAddress returns a form option that sets the `mailAddress` field (e.g., for User).
func (*FormOptionService) WithMailAddress(mailAddress string) *FormOption {
	// ToDo: validate mailAddress (Note: The validation remains as simple not-empty check)
	return &FormOption{formMailAddress, func(form *FormParams) error {
		if mailAddress == "" {
			return newValidationError("mailAddress must not be empty")
		}
		form.Set(formMailAddress.Value(), mailAddress)
		return nil
	}}
}

// WithMailNotify returns a form option that sets the `mailNotify` field (e.g., for Wiki, Issue).
func (*FormOptionService) WithMailNotify(enabled bool) *FormOption {
	return &FormOption{formMailNotify, func(form *FormParams) error {
		form.Set(formMailNotify.Value(), strconv.FormatBool(enabled))
		return nil
	}}
}

// WithPassword returns a form option that sets the `password` field (e.g., for User).
// It validates that the password meets the minimum length requirement (8 characters).
func (*FormOptionService) WithPassword(password string) *FormOption {
	return &FormOption{formPassword, func(form *FormParams) error {

		if len(password) < 8 {
			return newValidationError("password must be at least 8 characters long")
		}

		form.Set(formPassword.Value(), password)
		return nil
	}}
}

// WithProjectLeaderCanEditProjectLeader returns a form option that sets the `projectLeaderCanEditProjectLeader` field (e.g., for Project).
func (*FormOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) *FormOption {
	return &FormOption{formProjectLeaderCanEditProjectLeader, func(form *FormParams) error {
		form.Set(formProjectLeaderCanEditProjectLeader.Value(), strconv.FormatBool(enabled))
		return nil
	}}
}

// WithRoleType returns a form option that sets the `roleType` field (e.g., for User).
func (*FormOptionService) WithRoleType(roleType Role) *FormOption {
	return &FormOption{formRoleType, func(form *FormParams) error {
		if roleType < 1 || 6 < roleType {
			return newValidationError("roleType must be between 1 and 6")
		}
		form.Set(formRoleType.Value(), strconv.Itoa(int(roleType)))
		return nil
	}}
}

// WithSubtaskingEnabled returns a form option that sets the `subtaskingEnabled` field (e.g., for Project).
func (*FormOptionService) WithSubtaskingEnabled(enabled bool) *FormOption {
	return &FormOption{formSubtaskingEnabled, func(form *FormParams) error {
		form.Set(formSubtaskingEnabled.Value(), strconv.FormatBool(enabled))
		return nil
	}}
}

// WithTextFormattingRule returns a form option that sets the `textFormattingRule` field (e.g., for Project).
func (*FormOptionService) WithTextFormattingRule(format Format) *FormOption {
	return &FormOption{formTextFormattingRule, func(form *FormParams) error {
		if format != FormatBacklog && format != FormatMarkdown {
			msg := fmt.Sprintf("format must be only '%s' or '%s'", string(FormatBacklog), string(FormatMarkdown))
			return newValidationError(msg)
		}
		form.Set(formTextFormattingRule.Value(), string(format))
		return nil
	}}
}

// WithUserID creates a form option to set the user's ID (login name).
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-user
func (*FormOptionService) WithUserID(userID string) *FormOption {
	return &FormOption{formUserID, func(form *FormParams) error {
		if userID == "" {
			return newValidationError("userID must not be empty")
		}
		form.Set(formUserID.Value(), userID)
		return nil
	}}
}

// WithSendMail creates a form option to specify whether to send
// an invitation email to the newly created user.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-user
func (*FormOptionService) WithSendMail(enabled bool) *FormOption {
	return &FormOption{formSendMail, func(form *FormParams) error {
		form.Set(formSendMail.Value(), strconv.FormatBool(enabled))
		return nil
	}}
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
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/get-activities
func (s *ActivityOptionService) WithQueryActivityTypeIDs(typeIDs []int) *QueryOption {
	return s.support.query.WithActivityTypeIDs(typeIDs)
}

// WithQueryMinID returns a query option to filter activities
// that have IDs greater than or equal to the given value.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/get-activities
func (s *ActivityOptionService) WithQueryMinID(id int) *QueryOption {
	return s.support.query.WithMinID(id)
}

// WithQueryMaxID returns a query option to filter activities
// that have IDs less than or equal to the given value.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/get-activities
func (s *ActivityOptionService) WithQueryMaxID(id int) *QueryOption {
	return s.support.query.WithMaxID(id)
}

// WithQueryCount returns a query option to limit the number of activities
// returned by the API (max 100).
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/get-activities
func (s *ActivityOptionService) WithQueryCount(count int) *QueryOption {
	return s.support.query.WithCount(count)
}

// WithQueryOrder returns a query option to specify the result order (asc or desc).
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/get-activities
func (s *ActivityOptionService) WithQueryOrder(order Order) *QueryOption {
	return s.support.query.WithOrder(order)
}

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
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectOptionService) WithQueryAll(enabled bool) *QueryOption {
	return s.support.query.WithAll(enabled)
}

// WithQueryArchived returns a query option that filters projects
// by their archived status.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectOptionService) WithQueryArchived(enabled bool) *QueryOption {
	return s.support.query.WithArchived(enabled)
}

// --- Form Options ------------------------------------------------------------

// WithFormName returns a form option to set the project's name.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectOptionService) WithFormName(name string) *FormOption {
	return s.support.form.WithName(name)
}

// WithFormKey returns a form option to set the project's key.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectOptionService) WithFormKey(key string) *FormOption {
	return s.support.form.WithKey(key)
}

// WithFormArchived returns a form option to set the project's archived status.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectOptionService) WithFormArchived(enabled bool) *FormOption {
	return s.support.form.WithArchived(enabled)
}

// WithFormChartEnabled returns a form option to enable or disable charts
// for the project.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectOptionService) WithFormChartEnabled(enabled bool) *FormOption {
	return s.support.form.WithChartEnabled(enabled)
}

// WithFormSubtaskingEnabled returns a form option to enable or disable
// subtasking for the project.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectOptionService) WithFormSubtaskingEnabled(enabled bool) *FormOption {
	return s.support.form.WithSubtaskingEnabled(enabled)
}

// WithFormProjectLeaderCanEditProjectLeader returns a form option to set
// whether the project leader can edit the project leader field.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectOptionService) WithFormProjectLeaderCanEditProjectLeader(enabled bool) *FormOption {
	return s.support.form.WithProjectLeaderCanEditProjectLeader(enabled)
}

// WithFormTextFormattingRule returns a form option to set the text formatting
// rule (Backlog or Markdown) for the project.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/update-project
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
// You can add WithQueryXXX methods here when needed.

// --- Form Options ------------------------------------------------------------

// WithFormPassword returns a form option to set the user's password.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserOptionService) WithFormPassword(password string) *FormOption {
	return s.support.form.WithPassword(password)
}

// WithFormMailAddress returns a form option to set the user's mail address.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserOptionService) WithFormMailAddress(mail string) *FormOption {
	return s.support.form.WithMailAddress(mail)
}

// WithFormRoleType returns a form option to set the user's role type.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserOptionService) WithFormRoleType(role Role) *FormOption {
	return s.support.form.WithRoleType(role)
}

// WithFormName returns a form option to set the user's display name.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserOptionService) WithFormName(name string) *FormOption {
	return s.support.form.WithName(name)
}

// // WithFormUserId returns a form option to set the user's ID (login name).

// // Backlog API reference:

// // https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserOptionService) WithFormUserId(userID string) *FormOption {
	return s.support.form.WithUserID(userID)
}

// WithFormSendMail returns a form option to specify whether to send
// an invitation email to the newly created user.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserOptionService) WithFormSendMail(enabled bool) *FormOption {
	return s.support.form.WithSendMail(enabled)
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

// WithQueryKeyword returns a query option to search wiki pages by keyword.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiOptionService) WithQueryKeyword(keyword string) *QueryOption {
	return s.support.query.WithKeyword(keyword)
}

// WithFormName returns a form option to set the wiki page name.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-wiki-page
func (s *WikiOptionService) WithFormName(name string) *FormOption {
	return s.support.form.WithName(name)
}

// WithFormContent returns a form option to set the wiki page content.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/add-wiki-page
func (s *WikiOptionService) WithFormContent(content string) *FormOption {
	return s.support.form.WithContent(content)
}

// WithFormMailNotify returns a form option to enable or disable mail notifications
// for wiki updates or deletions.
//
// Backlog API reference:
//
//	https://developer.nulab.com/docs/backlog/api/2/update-wiki-page
func (s *WikiOptionService) WithFormMailNotify(enabled bool) *FormOption {
	return s.support.form.WithMailNotify(enabled)
}

//
// ──────────────────────────────────────────────────────────────
//  Internal Helpers
// ──────────────────────────────────────────────────────────────
//

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
