package backlog

import (
	"fmt"
	"strconv"
)

type formType int

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

type queryType int

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

type queryOptionFunc func(query *QueryParams) error

// QueryOption represents an option for a request query.
type QueryOption struct {
	t queryType
	f queryOptionFunc
}

func (o *QueryOption) validate(validTypes []queryType) error {
	for _, valid := range validTypes {
		if o.t == valid {
			return nil
		}
	}
	return newInvalidQueryOptionError(o.t, validTypes)
}

func (o *QueryOption) set(query *QueryParams) error {
	return o.f(query)
}

func withQueryActivityTypeIDs(typeIDs []int) *QueryOption {
	return &QueryOption{queryActivityTypeIDs, func(query *QueryParams) error {
		for _, id := range typeIDs {
			if id < 1 || 26 < id {
				return newValidationError("activityTypeId must be between 1 and 26")
			}
			query.Add(queryActivityTypeIDs.Value(), strconv.Itoa(id))
		}
		return nil
	}}
}

func withQueryAll(enabeld bool) *QueryOption {
	return &QueryOption{queryAll, func(query *QueryParams) error {
		query.Set(queryAll.Value(), strconv.FormatBool(enabeld))
		return nil
	}}
}

func withQueryArchived(archived bool) *QueryOption {
	return &QueryOption{queryArchived, func(query *QueryParams) error {
		query.Set(queryArchived.Value(), strconv.FormatBool(archived))
		return nil
	}}
}

func withQueryCount(count int) *QueryOption {
	return &QueryOption{queryCount, func(query *QueryParams) error {
		if count < 1 || 100 < count {
			return newValidationError("count must be between 1 and 100")
		}
		query.Set(queryCount.Value(), strconv.Itoa(count))
		return nil
	}}
}

func withQueryKeyword(keyword string) *QueryOption {
	return &QueryOption{queryKeyword, func(query *QueryParams) error {
		query.Set(queryKeyword.Value(), keyword)
		return nil
	}}
}

func withQueryMaxID(maxID int) *QueryOption {
	return &QueryOption{queryMaxID, func(query *QueryParams) error {
		if maxID < 1 {
			return newValidationError("maxId must not be less than 1")
		}
		query.Set(queryMaxID.Value(), strconv.Itoa(maxID))
		return nil
	}}
}

func withQueryMinID(minID int) *QueryOption {
	return &QueryOption{queryMinID, func(query *QueryParams) error {
		if minID < 1 {
			return newValidationError("minId must not be less than 1")
		}
		query.Set(queryMinID.Value(), strconv.Itoa(minID))
		return nil
	}}
}

func withQueryOrder(order Order) *QueryOption {
	return &QueryOption{queryOrder, func(query *QueryParams) error {
		if order != OrderAsc && order != OrderDesc {
			msg := fmt.Sprintf("order must be only '%s' or '%s'", string(OrderAsc), string(OrderDesc))
			return newValidationError(msg)
		}
		query.Set(queryOrder.Value(), string(order))
		return nil
	}}
}

// QueryOptionService has methods to make option for request query.
type QueryOptionService struct {
}

// WithActivityTypeIDs returns option to set `activityTypeId`.
func (s *QueryOptionService) WithActivityTypeIDs(typeIDs []int) *QueryOption {
	return withQueryActivityTypeIDs(typeIDs)
}

// WithAll returns option to set `all`.
func (s *QueryOptionService) WithAll(enabeld bool) *QueryOption {
	return withQueryAll(enabeld)
}

// WithArchived returns option to set `archived`.
func (s *QueryOptionService) WithArchived(archived bool) *QueryOption {
	return withQueryArchived(archived)
}

// WithCount returns option to set `count`.
func (s *QueryOptionService) WithCount(count int) *QueryOption {
	return withQueryCount(count)
}

// WithKeyword returns option to set `keyword`.
func (s *QueryOptionService) WithKeyword(keyword string) *QueryOption {
	return withQueryKeyword(keyword)
}

// WithMaxID returns option to set `maxId`.
func (s *QueryOptionService) WithMaxID(maxID int) *QueryOption {
	return withQueryMaxID(maxID)
}

// WithMinID returns option to set `minId`.
func (s *QueryOptionService) WithMinID(minID int) *QueryOption {
	return withQueryMinID(minID)
}

// WithOrder returns option to set `order`.
func (s *QueryOptionService) WithOrder(order Order) *QueryOption {
	return withQueryOrder(order)
}

type formOptionFunc func(form *FormParams) error

// FormOption is option to set form parameter for request.
type FormOption struct {
	t formType
	f formOptionFunc
}

func (o *FormOption) validate(validTypes []formType) error {
	for _, valid := range validTypes {
		if o.t == valid {
			return nil
		}
	}
	return newInvalidFormOptionError(o.t, validTypes)
}

func (o *FormOption) set(form *FormParams) error {
	return o.f(form)
}

func withFormArchived(archived bool) *FormOption {
	return &FormOption{formArchived, func(form *FormParams) error {
		form.Set(formArchived.Value(), strconv.FormatBool(archived))
		return nil
	}}
}

func withFormChartEnabled(enabeld bool) *FormOption {
	return &FormOption{formChartEnabled, func(form *FormParams) error {
		form.Set(formChartEnabled.Value(), strconv.FormatBool(enabeld))
		return nil
	}}
}

func withFormContent(content string) *FormOption {
	return &FormOption{formContent, func(form *FormParams) error {
		if content == "" {
			return newValidationError("content must not be empty")
		}
		form.Set(formContent.Value(), content)
		return nil
	}}
}

func withFormKey(key string) *FormOption {
	return &FormOption{formKey, func(form *FormParams) error {
		if key == "" {
			return newValidationError("key must not be empty")
		}
		form.Set(formKey.Value(), key)
		return nil
	}}
}

func withFormName(name string) *FormOption {
	return &FormOption{formName, func(form *FormParams) error {
		if name == "" {
			return newValidationError("name must not be empty")
		}
		form.Set(formName.Value(), name)
		return nil
	}}
}

func withFormMailAddress(mailAddress string) *FormOption {
	// ToDo: validate mailAddress
	return &FormOption{formMailAddress, func(form *FormParams) error {
		if mailAddress == "" {
			return newValidationError("mailAddress must not be empty")
		}
		form.Set(formMailAddress.Value(), mailAddress)
		return nil
	}}
}

func withFormMailNotify(enabeld bool) *FormOption {
	return &FormOption{formMailNotify, func(form *FormParams) error {
		form.Set(formMailNotify.Value(), strconv.FormatBool(enabeld))
		return nil
	}}
}

func withFormPassword(password string) *FormOption {
	return &FormOption{formPassword, func(form *FormParams) error {
		if password == "" {
			return newValidationError("password must not be empty")
		}
		form.Set(formPassword.Value(), password)
		return nil
	}}
}

func withFormProjectLeaderCanEditProjectLeader(enabeld bool) *FormOption {
	return &FormOption{formProjectLeaderCanEditProjectLeader, func(form *FormParams) error {
		form.Set(formProjectLeaderCanEditProjectLeader.Value(), strconv.FormatBool(enabeld))
		return nil
	}}
}

func withFormRoleType(roleType Role) *FormOption {
	return &FormOption{formRoleType, func(form *FormParams) error {
		if roleType < 1 || 6 < roleType {
			return newValidationError("roleType must be between 1 and 6")
		}
		form.Set(formRoleType.Value(), strconv.Itoa(int(roleType)))
		return nil
	}}
}

func withFormSubtaskingEnabled(enabeld bool) *FormOption {
	return &FormOption{formSubtaskingEnabled, func(form *FormParams) error {
		form.Set(formSubtaskingEnabled.Value(), strconv.FormatBool(enabeld))
		return nil
	}}
}

func withFormTextFormattingRule(format Format) *FormOption {
	return &FormOption{formTextFormattingRule, func(form *FormParams) error {
		if format != FormatBacklog && format != FormatMarkdown {
			msg := fmt.Sprintf("format must be only '%s' or '%s'", string(FormatBacklog), string(FormatMarkdown))
			return newValidationError(msg)
		}
		form.Set(formTextFormattingRule.Value(), string(format))
		return nil
	}}
}

// ActivityOptionService has methods to make option for ActivityService.
type ActivityOptionService struct {
}

// WithQueryActivityTypeIDs returns option to set `activityTypeId` for user.
func (*ActivityOptionService) WithQueryActivityTypeIDs(typeIDs []int) *QueryOption {
	return withQueryActivityTypeIDs(typeIDs)
}

// WithQueryMinID returns option to set `minId` for user.
func (*ActivityOptionService) WithQueryMinID(minID int) *QueryOption {
	return withQueryMinID(minID)
}

// WithQueryMaxID returns option to set `maxId` for user.
func (*ActivityOptionService) WithQueryMaxID(maxID int) *QueryOption {
	return withQueryMaxID(maxID)
}

// WithQueryCount returns option to set `count` for user.
func (*ActivityOptionService) WithQueryCount(count int) *QueryOption {
	return withQueryCount(count)
}

// WithQueryOrder returns option to set `order` for user.
func (*ActivityOptionService) WithQueryOrder(order Order) *QueryOption {
	return withQueryOrder(order)
}

// ProjectOptionService has methods to make option for ProjectService.
type ProjectOptionService struct {
}

// WithQueryAll returns a query option that sets the `all` field in the URL query parameter.
func (*ProjectOptionService) WithQueryAll(enabeld bool) *QueryOption {
	return withQueryAll(enabeld)
}

// WithQueryArchived returns a query option that sets the `archived` field in the URL query parameter.
func (*ProjectOptionService) WithQueryArchived(archived bool) *QueryOption {
	return withQueryArchived(archived)
}

// WithFormKey returns a form option that sets the `key` field in the request body for a project.
func (*ProjectOptionService) WithFormKey(key string) *FormOption {
	return withFormKey(key)
}

// WithFormName returns a form option that sets the `name` field in the request body for a project.
func (*ProjectOptionService) WithFormName(name string) *FormOption {
	return withFormName(name)
}

// WithFormChartEnabled returns a form option that sets the `chartEnabled` field in the request body for a project.
func (*ProjectOptionService) WithFormChartEnabled(enabeld bool) *FormOption {
	return withFormChartEnabled(enabeld)
}

// WithFormSubtaskingEnabled returns a form option that sets the `subtaskingEnabled` field in the request body for a project.
func (*ProjectOptionService) WithFormSubtaskingEnabled(enabeld bool) *FormOption {
	return withFormSubtaskingEnabled(enabeld)
}

// WithFormProjectLeaderCanEditProjectLeader returns a form option that sets the `projectLeaderCanEditProjectLeader` field in the request body for a project.
func (*ProjectOptionService) WithFormProjectLeaderCanEditProjectLeader(enabeld bool) *FormOption {
	return withFormProjectLeaderCanEditProjectLeader(enabeld)
}

// WithFormTextFormattingRule returns a form option that sets the `textFormattingRule` field in the request body for a project.
func (*ProjectOptionService) WithFormTextFormattingRule(format Format) *FormOption {
	return withFormTextFormattingRule(format)
}

// WithFormArchived returns a form option that sets the `archived` field in the request body for a project.
func (*ProjectOptionService) WithFormArchived(archived bool) *FormOption {
	return withFormArchived(archived)
}

// UserOptionService has methods to make option for UserService.
type UserOptionService struct {
}

// WithFormPassword returns a form option to set the `password` field in the request body for a user.
func (*UserOptionService) WithFormPassword(password string) *FormOption {
	return withFormPassword(password)
}

// WithFormName returns a form option to set the `name` field in the request body for a user.
func (*UserOptionService) WithFormName(name string) *FormOption {
	return withFormName(name)
}

// WithFormMailAddress returns a form option to set the `mailAddress` field in the request body for a user.
func (*UserOptionService) WithFormMailAddress(mailAddress string) *FormOption {
	return withFormMailAddress(mailAddress)
}

// WithFormRoleType returns a form option that sets the `roleType` field in the request body for a user.
func (*UserOptionService) WithFormRoleType(roleType Role) *FormOption {
	return withFormRoleType(roleType)
}

// WikiOptionService provides methods to create options for WikiService.
type WikiOptionService struct {
}

// WithQueryKeyword returns a query option that sets the `keyword` field in the URL query parameter.
func (*WikiOptionService) WithQueryKeyword(keyword string) *QueryOption {
	return withQueryKeyword(keyword)
}

// WithFormName returns a form option that sets the `name` field in the request body for a wiki.
func (*WikiOptionService) WithFormName(name string) *FormOption {
	return withFormName(name)
}

// WithFormContent returns a form option that sets the `content` field in the request body for a wiki.
func (*WikiOptionService) WithFormContent(content string) *FormOption {
	return withFormContent(content)
}

// WithFormMailNotify returns a form option that sets the `mailNotify` field in the request body for a wiki (e.g., true/false).
func (*WikiOptionService) WithFormMailNotify(enabeld bool) *FormOption {
	return withFormMailNotify(enabeld)
}

// checkRequiredOptionTypes checks if at least one form type specified in
// the requiredTypes slice is present in the given options.
// It returns true if a required type is found, otherwise false.
func checkRequiredOptionTypes(options []*FormOption, requiredTypes []formType) bool {
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
