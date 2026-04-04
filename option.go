package backlog

import (
	"fmt"
	"net/url"
	"strconv"
)

// optionRegistry provides shared access to the option builder.
// It is intended for internal composition by each XxxOptionService.
type optionRegistry struct {
	option *OptionService
}

//
// ──────────────────────────────────────────────────────────────
//  Option Types
// ──────────────────────────────────────────────────────────────
//

// requestOptionType is a constraint for InvalidOptionError.
// Restricted to queryType and formType.
type requestOptionType interface {
	queryType | formType
	Value() string
}

// --- QueryOption Types -------------------------------------------------------------

// queryType represents the distinct query parameter keys for Backlog API requests.
type queryType string

// Value returns the string representation of the query parameter key for the API request.
func (t queryType) Value() string {
	return string(t)
}

// --- FormOption Types -------------------------------------------------------------

// formType represents the distinct form field keys available for Backlog API requests.
type formType string

// Value returns the string representation of the form field key for the API request.
func (t formType) Value() string {
	return string(t)
}

//
// ──────────────────────────────────────────────────────────────
//  RequestOption interface
// ──────────────────────────────────────────────────────────────
//

// RequestOption defines a common interface for all option types.
// It allows unified validation and application handling across different request-level options.
// Callers can implement this interface to provide custom options (e.g. for mocking in tests).
type RequestOption interface {
	Check() error
	Set(url.Values) error
}

//
// ──────────────────────────────────────────────────────────────
//  apiOption — unified internal option type
// ──────────────────────────────────────────────────────────────
//

// apiOption is the single internal representation of a request option.
// It replaces the former QueryOption and FormOption types.
// The t field is either a queryType or formType, used to validate allowed options.
type apiOption struct {
	t         any            // queryType or formType
	checkFunc func() error   // optional validation
	setFunc   func(url.Values) error // applies the value to query/form
}

// Check validates the option by executing its checkFunc, if defined.
func (o *apiOption) Check() error {
	if o.checkFunc != nil {
		return o.checkFunc()
	}
	return nil
}

// Set applies the option value to the given url.Values.
func (o *apiOption) Set(v url.Values) error {
	if o.setFunc == nil {
		return newValidationError("option has no setter")
	}
	return o.setFunc(v)
}

// validateQueryType ensures the option is an allowed query parameter type.
func (o *apiOption) validateQueryType(validTypes []queryType) error {
	t, ok := o.t.(queryType)
	if !ok {
		// formType passed where queryType expected
		var zero queryType
		return newInvalidOptionError(zero, validTypes)
	}
	for _, valid := range validTypes {
		if t == valid {
			return nil
		}
	}
	return newInvalidOptionError(t, validTypes)
}

// validateFormType ensures the option is an allowed form field type.
func (o *apiOption) validateFormType(validTypes []formType) error {
	t, ok := o.t.(formType)
	if !ok {
		// queryType passed where formType expected
		var zero formType
		return newInvalidOptionError(zero, validTypes)
	}
	for _, valid := range validTypes {
		if t == valid {
			return nil
		}
	}
	return newInvalidOptionError(t, validTypes)
}

//
// ──────────────────────────────────────────────────────────────
//  OptionService — unified builder
// ──────────────────────────────────────────────────────────────
//

// OptionService provides builders for both query and form options.
// Each XxxOptionService selectively exposes only the valid methods.
type OptionService struct{}

// applyQueryOptions validates and applies query options to the given url.Values.
func (s *OptionService) applyQueryOptions(query url.Values, validTypes []queryType, opts ...RequestOption) error {
	for _, opt := range opts {
		ao, ok := opt.(*apiOption)
		if !ok {
			// user-provided custom implementation: just Check + Set
			if err := opt.Check(); err != nil {
				return err
			}
			if err := opt.Set(query); err != nil {
				return err
			}
			continue
		}
		if err := ao.validateQueryType(validTypes); err != nil {
			return err
		}
		if err := ao.Check(); err != nil {
			return err
		}
		if err := ao.Set(query); err != nil {
			return err
		}
	}
	return nil
}

// applyFormOptions validates and applies form options to the given url.Values.
func (s *OptionService) applyFormOptions(form url.Values, validTypes []formType, opts ...RequestOption) error {
	for _, opt := range opts {
		ao, ok := opt.(*apiOption)
		if !ok {
			if err := opt.Check(); err != nil {
				return err
			}
			if err := opt.Set(form); err != nil {
				return err
			}
			continue
		}
		if err := ao.validateFormType(validTypes); err != nil {
			return err
		}
		if err := ao.Check(); err != nil {
			return err
		}
		if err := ao.Set(form); err != nil {
			return err
		}
	}
	return nil
}

// --- Boolean options ------------------------------------------------------------

// WithAll returns an option to set the `all` query parameter.
func (s *OptionService) WithAll(enabled bool) RequestOption {
	return &apiOption{
		t: queryAll,
		setFunc: func(query url.Values) error {
			query.Set(queryAll.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithArchived returns an option to set the `archived` parameter (query or form).
func (s *OptionService) WithArchived(enabled bool) RequestOption {
	// formArchived and queryArchived share the same string value "archived",
	// so we use formArchived as the canonical type here and accept both in services.
	return &apiOption{
		t: formArchived,
		setFunc: func(v url.Values) error {
			v.Set(formArchived.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithChartEnabled returns a form option that sets the `chartEnabled` field.
func (s *OptionService) WithChartEnabled(enabled bool) RequestOption {
	return &apiOption{
		t: formChartEnabled,
		setFunc: func(form url.Values) error {
			form.Set(formChartEnabled.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithMailNotify returns a form option that sets the `mailNotify` field.
func (s *OptionService) WithMailNotify(enabled bool) RequestOption {
	return &apiOption{
		t: formMailNotify,
		setFunc: func(form url.Values) error {
			form.Set(formMailNotify.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithProjectLeaderCanEditProjectLeader returns a form option.
func (s *OptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return &apiOption{
		t: formProjectLeaderCanEditProjectLeader,
		setFunc: func(form url.Values) error {
			form.Set(formProjectLeaderCanEditProjectLeader.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithSendMail returns a form option to specify whether to send an invitation email.
func (s *OptionService) WithSendMail(enabled bool) RequestOption {
	return &apiOption{
		t: formSendMail,
		setFunc: func(form url.Values) error {
			form.Set(formSendMail.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithSubtaskingEnabled returns a form option that sets the `subtaskingEnabled` field.
func (s *OptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return &apiOption{
		t: formSubtaskingEnabled,
		setFunc: func(form url.Values) error {
			form.Set(formSubtaskingEnabled.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// --- Integer options ------------------------------------------------------------

// WithCount returns an option to set the `count` query parameter.
func (s *OptionService) WithCount(count int) RequestOption {
	return &apiOption{
		t: queryCount,
		checkFunc: func() error {
			if count < 1 || 100 < count {
				return newValidationError("count must be between 1 and 100")
			}
			return nil
		},
		setFunc: func(query url.Values) error {
			query.Set(queryCount.Value(), strconv.Itoa(count))
			return nil
		},
	}
}

// WithMaxID returns an option to set the `maxId` query parameter.
func (s *OptionService) WithMaxID(id int) RequestOption {
	return &apiOption{
		t: queryMaxID,
		checkFunc: func() error {
			return validateActivityID(id, "maxID")
		},
		setFunc: func(query url.Values) error {
			query.Set(queryMaxID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// WithMinID returns an option to set the `minId` query parameter.
func (s *OptionService) WithMinID(id int) RequestOption {
	return &apiOption{
		t: queryMinID,
		checkFunc: func() error {
			return validateActivityID(id, "minID")
		},
		setFunc: func(query url.Values) error {
			query.Set(queryMinID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// WithUserID returns a form option to set the user's ID.
func (s *OptionService) WithUserID(id int) RequestOption {
	return &apiOption{
		t: formUserID,
		checkFunc: func() error {
			return validateID(id, "userID")
		},
		setFunc: func(form url.Values) error {
			form.Set(formUserID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// --- String options ------------------------------------------------------------

// WithContent returns a form option that sets the `content` field.
func (s *OptionService) WithContent(content string) RequestOption {
	return &apiOption{
		t: formContent,
		checkFunc: func() error {
			if content == "" {
				return newValidationError("content must not be empty")
			}
			return nil
		},
		setFunc: func(form url.Values) error {
			form.Set(formContent.Value(), content)
			return nil
		},
	}
}

// WithKey returns a form option that sets the `key` field.
func (s *OptionService) WithKey(key string) RequestOption {
	return &apiOption{
		t: formKey,
		checkFunc: func() error {
			if key == "" {
				return newValidationError("key must not be empty")
			}
			return nil
		},
		setFunc: func(form url.Values) error {
			form.Set(formKey.Value(), key)
			return nil
		},
	}
}

// WithKeyword returns an option to set the `keyword` query parameter.
func (s *OptionService) WithKeyword(keyword string) RequestOption {
	return &apiOption{
		t: queryKeyword,
		setFunc: func(query url.Values) error {
			query.Set(queryKeyword.Value(), keyword)
			return nil
		},
	}
}

// WithMailAddress returns a form option that sets the `mailAddress` field.
func (s *OptionService) WithMailAddress(mailAddress string) RequestOption {
	// ToDo: validate mailAddress (Note: The validation remains as simple not-empty check)
	return &apiOption{
		t: formMailAddress,
		checkFunc: func() error {
			if mailAddress == "" {
				return newValidationError("mailAddress must not be empty")
			}
			return nil
		},
		setFunc: func(form url.Values) error {
			form.Set(formMailAddress.Value(), mailAddress)
			return nil
		},
	}
}

// WithName returns a form option that sets the `name` field.
func (s *OptionService) WithName(name string) RequestOption {
	return &apiOption{
		t: formName,
		checkFunc: func() error {
			if name == "" {
				return newValidationError("name must not be empty")
			}
			return nil
		},
		setFunc: func(form url.Values) error {
			form.Set(formName.Value(), name)
			return nil
		},
	}
}

// WithPassword returns a form option that sets the `password` field.
func (s *OptionService) WithPassword(password string) RequestOption {
	return &apiOption{
		t: formPassword,
		checkFunc: func() error {
			if len(password) < 8 {
				return newValidationError("password must be at least 8 characters long")
			}
			return nil
		},
		setFunc: func(form url.Values) error {
			form.Set(formPassword.Value(), password)
			return nil
		},
	}
}

// --- Enum or special options ----------------------------------------------------

// WithActivityTypeIDs returns an option to set multiple `activityTypeId[]` query parameters.
func (s *OptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return &apiOption{
		t: queryActivityTypeIDs,
		checkFunc: func() error {
			for _, id := range typeIDs {
				if err := validateActivityID(id, "activityTypeIds"); err != nil {
					return err
				}
			}
			return nil
		},
		setFunc: func(query url.Values) error {
			for _, id := range typeIDs {
				query.Add(queryActivityTypeIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithOrder returns an option to set the `order` query parameter.
func (s *OptionService) WithOrder(order Order) RequestOption {
	return &apiOption{
		t: queryOrder,
		checkFunc: func() error {
			if order != OrderAsc && order != OrderDesc {
				msg := fmt.Sprintf("order must be only '%s' or '%s'", string(OrderAsc), string(OrderDesc))
				return newValidationError(msg)
			}
			return nil
		},
		setFunc: func(query url.Values) error {
			query.Set(queryOrder.Value(), string(order))
			return nil
		},
	}
}

// WithRoleType returns a form option that sets the `roleType` field.
func (s *OptionService) WithRoleType(roleType Role) RequestOption {
	return &apiOption{
		t: formRoleType,
		checkFunc: func() error {
			if roleType < 1 || 6 < roleType {
				return newValidationError("roleType must be between 1 and 6")
			}
			return nil
		},
		setFunc: func(form url.Values) error {
			form.Set(formRoleType.Value(), strconv.Itoa(int(roleType)))
			return nil
		},
	}
}

// WithTextFormattingRule returns a form option that sets the `textFormattingRule` field.
func (s *OptionService) WithTextFormattingRule(format Format) RequestOption {
	return &apiOption{
		t: formTextFormattingRule,
		checkFunc: func() error {
			if format != FormatBacklog && format != FormatMarkdown {
				msg := fmt.Sprintf("format must be only '%s' or '%s'", string(FormatBacklog), string(FormatMarkdown))
				return newValidationError(msg)
			}
			return nil
		},
		setFunc: func(form url.Values) error {
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
type ActivityOptionService struct {
	registry *optionRegistry
}

func (s *ActivityOptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return s.registry.option.WithActivityTypeIDs(typeIDs)
}

func (s *ActivityOptionService) WithMinID(id int) RequestOption {
	return s.registry.option.WithMinID(id)
}

func (s *ActivityOptionService) WithMaxID(id int) RequestOption {
	return s.registry.option.WithMaxID(id)
}

func (s *ActivityOptionService) WithCount(count int) RequestOption {
	return s.registry.option.WithCount(count)
}

func (s *ActivityOptionService) WithOrder(order Order) RequestOption {
	return s.registry.option.WithOrder(order)
}

//
// ──────────────────────────────────────────────────────────────
//  ProjectOptionService
// ──────────────────────────────────────────────────────────────
//

// ProjectOptionService provides a domain-specific set of option builders
// for operations within the ProjectService.
type ProjectOptionService struct {
	registry *optionRegistry
}

func (s *ProjectOptionService) WithAll(enabled bool) RequestOption {
	return s.registry.option.WithAll(enabled)
}

func (s *ProjectOptionService) WithArchived(enabled bool) RequestOption {
	return s.registry.option.WithArchived(enabled)
}

func (s *ProjectOptionService) WithChartEnabled(enabled bool) RequestOption {
	return s.registry.option.WithChartEnabled(enabled)
}

func (s *ProjectOptionService) WithKey(key string) RequestOption {
	return s.registry.option.WithKey(key)
}

func (s *ProjectOptionService) WithName(name string) RequestOption {
	return s.registry.option.WithName(name)
}

func (s *ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return s.registry.option.WithProjectLeaderCanEditProjectLeader(enabled)
}

func (s *ProjectOptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return s.registry.option.WithSubtaskingEnabled(enabled)
}

func (s *ProjectOptionService) WithTextFormattingRule(format Format) RequestOption {
	return s.registry.option.WithTextFormattingRule(format)
}

//
// ──────────────────────────────────────────────────────────────
//  UserOptionService
// ──────────────────────────────────────────────────────────────
//

// UserOptionService provides a domain-specific set of option builders
// for operations within the UserService.
type UserOptionService struct {
	registry *optionRegistry
}

func (s *UserOptionService) WithMailAddress(mail string) RequestOption {
	return s.registry.option.WithMailAddress(mail)
}

func (s *UserOptionService) WithName(name string) RequestOption {
	return s.registry.option.WithName(name)
}

func (s *UserOptionService) WithPassword(password string) RequestOption {
	return s.registry.option.WithPassword(password)
}

func (s *UserOptionService) WithRoleType(role Role) RequestOption {
	return s.registry.option.WithRoleType(role)
}

func (s *UserOptionService) WithSendMail(enabled bool) RequestOption {
	return s.registry.option.WithSendMail(enabled)
}

func (s *UserOptionService) WithUserID(id int) RequestOption {
	return s.registry.option.WithUserID(id)
}

//
// ──────────────────────────────────────────────────────────────
//  WikiOptionService
// ──────────────────────────────────────────────────────────────
//

// WikiOptionService provides a domain-specific set of option builders
// for operations within the WikiService.
type WikiOptionService struct {
	registry *optionRegistry
}

func (s *WikiOptionService) WithKeyword(keyword string) RequestOption {
	return s.registry.option.WithKeyword(keyword)
}

func (s *WikiOptionService) WithContent(content string) RequestOption {
	return s.registry.option.WithContent(content)
}

func (s *WikiOptionService) WithMailNotify(enabled bool) RequestOption {
	return s.registry.option.WithMailNotify(enabled)
}

func (s *WikiOptionService) WithName(name string) RequestOption {
	return s.registry.option.WithName(name)
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
func validateActivityID(id int, key string) error {
	if id < 1 || id > maxActivityTypeID {
		return newValidationError(fmt.Sprintf("invalid %s: must be between 1 and %d", key, maxActivityTypeID))
	}
	return nil
}

// hasRequiredOption checks whether the provided options include at least one of the required form types.
func hasRequiredOption(options []RequestOption, requiredTypes []formType) bool {
	for _, opt := range options {
		ao, ok := opt.(*apiOption)
		if !ok {
			continue
		}
		t, ok := ao.t.(formType)
		if !ok {
			continue
		}
		for _, requiredType := range requiredTypes {
			if t == requiredType {
				return true
			}
		}
	}
	return false
}
