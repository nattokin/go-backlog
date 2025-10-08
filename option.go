package backlog

import (
	"fmt"
	"strconv"
)

// --- Request Parameter Types and Mappings ---

// formType represents the distinct fields (keys) that can be sent in a request form body.
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
	}

	return m[t]
}

// queryType represents the distinct fields (keys) that can be sent in a request URL query string.
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

// --- Query Option Definition ---

// queryOptionFunc is a function that applies a query option's value to the request's QueryParams.
type queryOptionFunc func(query *QueryParams) error

// QueryOption represents an option for a request query parameter.
// It encapsulates the query field type and the function to apply its value.
type QueryOption struct {
	t queryType       // The underlying query parameter type (e.g., queryKeyword)
	f queryOptionFunc // The function to set the parameter value
}

// validate checks if the current QueryOption is of a valid type for a specific API endpoint.
// It prevents using an option (e.g., 'keyword') on an endpoint that doesn't support it.
func (o *QueryOption) validate(validTypes []queryType) error {
	for _, valid := range validTypes {
		if o.t == valid {
			return nil
		}
	}
	// Returns a custom error indicating the option is invalid for this context.
	return newInvalidQueryOptionError(o.t, validTypes)
}

// set executes the embedded function to apply the option's value to the QueryParams struct.
func (o *QueryOption) set(query *QueryParams) error {
	return o.f(query)
}

// QueryOptionService has methods to make option for request query.
type QueryOptionService struct {
}

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

// --- Form Option Definition ---

// formOptionFunc is a function that applies a form option's value to the request's FormParams.
type formOptionFunc func(form *FormParams) error

// FormOption is option to set form parameter for a request body.
// It holds the form field type and the function to set its value.
type FormOption struct {
	t formType       // The underlying form parameter type (e.g., formPassword)
	f formOptionFunc // The function to set the parameter value
}

// validate checks if the current FormOption is of a valid type for a specific API endpoint.
// It ensures that only supported form fields are passed to the request.
func (o *FormOption) validate(validTypes []formType) error {
	for _, valid := range validTypes {
		if o.t == valid {
			return nil
		}
	}
	// Returns a custom error indicating the form option is invalid for this context.
	return newInvalidFormOptionError(o.t, validTypes)
}

// set executes the embedded function to apply the option's value to the FormParams struct.
func (o *FormOption) set(form *FormParams) error {
	return o.f(form)
}

// FormOptionService provides methods to create options for request bodies (FormOption).
type FormOptionService struct {
}

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

// ActivityOptionService provides methods to create query options for ActivityService.
// It aggregates the functionalities of QueryOptionService.
type ActivityOptionService struct {
	// Query holds all query option methods relevant to the activity context.
	Query *QueryOptionService
}

// --- Delegation Methods for Activity Query Options ---

// WithQueryActivityTypeIDs creates a query option for Activity Type I Ds.
func (s *ActivityOptionService) WithQueryActivityTypeIDs(ids []int) *QueryOption {
	return s.Query.WithActivityTypeIDs(ids)
}

// WithQueryMinID creates a query option for Min I D.
func (s *ActivityOptionService) WithQueryMinID(minID int) *QueryOption {
	return s.Query.WithMinID(minID)
}

// WithQueryMaxID creates a query option for Max I D.
func (s *ActivityOptionService) WithQueryMaxID(maxID int) *QueryOption {
	return s.Query.WithMaxID(maxID)
}

// WithQueryCount creates a query option for Count.
func (s *ActivityOptionService) WithQueryCount(count int) *QueryOption {
	return s.Query.WithCount(count)
}

// WithQueryOrder creates a query option for Order.
func (s *ActivityOptionService) WithQueryOrder(order Order) *QueryOption {
	return s.Query.WithOrder(order)
}

// ProjectOptionService provides methods to create both query and form options for ProjectService.
// It aggregates the functionalities of QueryOptionService and FormOptionService.
type ProjectOptionService struct {
	// Query holds all query option methods relevant to the project context.
	Query *QueryOptionService
	// Form holds all form option methods relevant to the project context.
	Form *FormOptionService
}

// --- Delegation Methods for Query Options (QueryOptionService) ---

// WithQueryAll returns a query option that includes archived projects in the result.
func (s *ProjectOptionService) WithQueryAll(enabled bool) *QueryOption {
	return s.Query.WithAll(enabled)
}

// WithQueryArchived returns a query option that filters projects by their archive status.
func (s *ProjectOptionService) WithQueryArchived(enabled bool) *QueryOption {
	// Assuming QueryOptionService has WithQueryArchived
	return s.Query.WithArchived(enabled)
}

// --- Delegation Methods for Form Options (FormOptionService) ---

// WithFormArchived returns a form option that sets the `archived` field for the project.
func (s *ProjectOptionService) WithFormArchived(enabled bool) *FormOption {
	return s.Form.WithArchived(enabled)
}

// WithFormChartEnabled returns a form option that sets the `chartEnabled` field for the project.
func (s *ProjectOptionService) WithFormChartEnabled(enabled bool) *FormOption {
	return s.Form.WithChartEnabled(enabled)
}

// WithFormKey returns a form option that sets the `key` field for the project.
func (s *ProjectOptionService) WithFormKey(key string) *FormOption {
	return s.Form.WithKey(key)
}

// WithFormName returns a form option that sets the `name` field for the project.
func (s *ProjectOptionService) WithFormName(name string) *FormOption {
	return s.Form.WithName(name)
}

// WithFormProjectLeaderCanEditProjectLeader returns a form option that sets the project leader edit permission.
func (s *ProjectOptionService) WithFormProjectLeaderCanEditProjectLeader(enabled bool) *FormOption {
	return s.Form.WithProjectLeaderCanEditProjectLeader(enabled)
}

// WithFormSubtaskingEnabled returns a form option that sets the `subtaskingEnabled` field for the project.
func (s *ProjectOptionService) WithFormSubtaskingEnabled(enabled bool) *FormOption {
	return s.Form.WithSubtaskingEnabled(enabled)
}

// WithFormTextFormattingRule returns a form option that sets the text formatting rule (format) for the project.
func (s *ProjectOptionService) WithFormTextFormattingRule(format Format) *FormOption {
	return s.Form.WithTextFormattingRule(format)
}

// UserOptionService provides methods to create form options for UserService.
// It aggregates the functionalities of FormOptionService.
type UserOptionService struct {
	// Form holds all form option methods relevant to the user context.
	Form *FormOptionService
}

// --- Delegation Methods for User Form Options ---

// WithFormPassword creates a form option for Password.
func (s *UserOptionService) WithFormPassword(password string) *FormOption {
	return s.Form.WithPassword(password)
}

// WithFormName creates a form option for Name.
func (s *UserOptionService) WithFormName(name string) *FormOption {
	return s.Form.WithName(name)
}

// WithFormMailAddress creates a form option for Mail Address.
func (s *UserOptionService) WithFormMailAddress(mailAddress string) *FormOption {
	return s.Form.WithMailAddress(mailAddress)
}

// WithFormRoleType creates a form option for Role Type.
func (s *UserOptionService) WithFormRoleType(roleType Role) *FormOption {
	return s.Form.WithRoleType(roleType)
}

// WikiOptionService provides methods to create both query and form options for WikiService.
// It aggregates the functionalities of QueryOptionService and FormOptionService.
type WikiOptionService struct {
	// Query holds all query option methods relevant to the wiki context.
	Query *QueryOptionService
	// Form holds all form option methods relevant to the wiki context.
	Form *FormOptionService
}

// --- Delegation Methods for Wiki Query Options ---

// WithQueryKeyword creates a query option for Keyword.
func (s *WikiOptionService) WithQueryKeyword(keyword string) *QueryOption {
	return s.Query.WithKeyword(keyword)
}

// --- Delegation Methods for Wiki Form Options ---

// WithFormName creates a form option for Name.
func (s *WikiOptionService) WithFormName(name string) *FormOption {
	return s.Form.WithName(name)
}

// WithFormContent creates a form option for Content.
func (s *WikiOptionService) WithFormContent(content string) *FormOption {
	return s.Form.WithContent(content)
}

// WithFormMailNotify creates a form option for Mail Notify.
func (s *WikiOptionService) WithFormMailNotify(enabled bool) *FormOption {
	return s.Form.WithMailNotify(enabled)
}

// hasRequiredFormOption checks if at least one form type specified in
// the requiredTypes slice is present in the given options.
// It returns true if a required type is found, otherwise false.
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
