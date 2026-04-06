package backlog

import (
	"fmt"
	"net/url"
	"strconv"
)

//
// ──────────────────────────────────────────────────────────────
//  API Option Type
// ──────────────────────────────────────────────────────────────
//

// apiParamOptionType represents the distinct parameter keys for Backlog API requests.
type apiParamOptionType string

// Value returns the string representation of the parameter key for the API request.
func (t apiParamOptionType) Value() string {
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
	Key() string
	Check() error
	Set(url.Values) error
}

//
// ──────────────────────────────────────────────────────────────
//  apiOption — unified internal option type
// ──────────────────────────────────────────────────────────────
//

// apiParamOption is the internal implementation of RequestOption.
//
// It encapsulates the parameter type together with optional validation
// and the logic required to apply the value to API request parameters.
// OptionService builder methods return instances of this struct.
type apiParamOption struct {
	t         apiParamOptionType     // canonical API parameter type
	checkFunc func() error           // optional validation executed before applying the option
	setFunc   func(url.Values) error // applies the option value to the provided values
}

// Key returns the API parameter key associated with this option.
func (o *apiParamOption) Key() string {
	return o.t.Value()
}

// Check validates the option by executing its checkFunc, if defined.
func (o *apiParamOption) Check() error {
	if o.checkFunc != nil {
		return o.checkFunc()
	}
	return nil
}

// Set applies the option value to the given url.Values.
func (o *apiParamOption) Set(v url.Values) error {
	if o.setFunc == nil {
		panic("option has no setter")
	}
	return o.setFunc(v)
}

//
// ──────────────────────────────────────────────────────────────
//  OptionService — unified builder
// ──────────────────────────────────────────────────────────────
//

// OptionService provides builders for request options.
// Each XxxOptionService selectively exposes only the valid methods.
type OptionService struct{}

// --- Boolean options ------------------------------------------------------------

// WithAll returns an option to set the `all` parameter.
func (s *OptionService) WithAll(enabled bool) RequestOption {
	return &apiParamOption{
		t: paramAll,
		setFunc: func(v url.Values) error {
			v.Set(paramAll.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithArchived returns an option to set the `archived` parameter.
func (s *OptionService) WithArchived(enabled bool) RequestOption {
	// apiArchived and queryArchived share the same string value "archived",
	// so we use apiArchived as the canonical type here and accept both in services.
	return &apiParamOption{
		t: paramArchived,
		setFunc: func(v url.Values) error {
			v.Set(paramArchived.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithChartEnabled returns a option that sets the `chartEnabled` field.
func (s *OptionService) WithChartEnabled(enabled bool) RequestOption {
	return &apiParamOption{
		t: paramChartEnabled,
		setFunc: func(v url.Values) error {
			v.Set(paramChartEnabled.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithMailNotify returns a option that sets the `mailNotify` field.
func (s *OptionService) WithMailNotify(enabled bool) RequestOption {
	return &apiParamOption{
		t: paramMailNotify,
		setFunc: func(v url.Values) error {
			v.Set(paramMailNotify.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithProjectLeaderCanEditProjectLeader returns a option.
func (s *OptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return &apiParamOption{
		t: paramProjectLeaderCanEditProjectLeader,
		setFunc: func(v url.Values) error {
			v.Set(paramProjectLeaderCanEditProjectLeader.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithSendMail returns a option to specify whether to send an invitation email.
func (s *OptionService) WithSendMail(enabled bool) RequestOption {
	return &apiParamOption{
		t: paramSendMail,
		setFunc: func(v url.Values) error {
			v.Set(paramSendMail.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithSubtaskingEnabled returns a option that sets the `subtaskingEnabled` field.
func (s *OptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return &apiParamOption{
		t: paramSubtaskingEnabled,
		setFunc: func(v url.Values) error {
			v.Set(paramSubtaskingEnabled.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// --- Integer options ------------------------------------------------------------

// WithCount returns an option to set the `count` parameter.
func (s *OptionService) WithCount(count int) RequestOption {
	return &apiParamOption{
		t: paramCount,
		checkFunc: func() error {
			if count < 1 || 100 < count {
				return newValidationError("count must be between 1 and 100")
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramCount.Value(), strconv.Itoa(count))
			return nil
		},
	}
}

// WithMaxID returns an option to set the `maxId` parameter.
func (s *OptionService) WithMaxID(id int) RequestOption {
	return &apiParamOption{
		t: paramMaxID,
		checkFunc: func() error {
			return validateActivityID(id, "maxID")
		},
		setFunc: func(v url.Values) error {
			v.Set(paramMaxID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// WithMinID returns an option to set the `minId` parameter.
func (s *OptionService) WithMinID(id int) RequestOption {
	return &apiParamOption{
		t: paramMinID,
		checkFunc: func() error {
			return validateActivityID(id, "minID")
		},
		setFunc: func(v url.Values) error {
			v.Set(paramMinID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// WithUserID returns a option to set the user's ID.
func (s *OptionService) WithUserID(id int) RequestOption {
	return &apiParamOption{
		t: paramUserID,
		checkFunc: func() error {
			return validateID(id, paramUserID.Value())
		},
		setFunc: func(v url.Values) error {
			v.Set(paramUserID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// --- String options ------------------------------------------------------------

// WithContent returns a option that sets the `content` field.
func (s *OptionService) WithContent(content string) RequestOption {
	return &apiParamOption{
		t: paramContent,
		checkFunc: func() error {
			if content == "" {
				return newValidationError("content must not be empty")
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramContent.Value(), content)
			return nil
		},
	}
}

// WithKey returns a option that sets the `key` field.
func (s *OptionService) WithKey(key string) RequestOption {
	return &apiParamOption{
		t: paramKey,
		checkFunc: func() error {
			if key == "" {
				return newValidationError("key must not be empty")
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramKey.Value(), key)
			return nil
		},
	}
}

// WithKeyword returns an option to set the `keyword` parameter.
func (s *OptionService) WithKeyword(keyword string) RequestOption {
	return &apiParamOption{
		t: paramKeyword,
		setFunc: func(v url.Values) error {
			v.Set(paramKeyword.Value(), keyword)
			return nil
		},
	}
}

// WithMailAddress returns a option that sets the `mailAddress` field.
func (s *OptionService) WithMailAddress(mailAddress string) RequestOption {
	// ToDo: validate mailAddress (Note: The validation remains as simple not-empty check)
	return &apiParamOption{
		t: paramMailAddress,
		checkFunc: func() error {
			if mailAddress == "" {
				return newValidationError("mailAddress must not be empty")
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramMailAddress.Value(), mailAddress)
			return nil
		},
	}
}

// WithName returns a option that sets the `name` field.
func (s *OptionService) WithName(name string) RequestOption {
	return &apiParamOption{
		t: paramName,
		checkFunc: func() error {
			if name == "" {
				return newValidationError("name must not be empty")
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramName.Value(), name)
			return nil
		},
	}
}

// WithPassword returns a option that sets the `password` field.
func (s *OptionService) WithPassword(password string) RequestOption {
	return &apiParamOption{
		t: paramPassword,
		checkFunc: func() error {
			if len(password) < 8 {
				return newValidationError("password must be at least 8 characters long")
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramPassword.Value(), password)
			return nil
		},
	}
}

// --- Enum or special options ----------------------------------------------------

// WithActivityTypeIDs returns an option to set multiple `activityTypeId[]` parameters.
func (s *OptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return &apiParamOption{
		t: paramActivityTypeIDs,
		checkFunc: func() error {
			for _, id := range typeIDs {
				if err := validateActivityID(id, "activityTypeIds"); err != nil {
					return err
				}
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			for _, id := range typeIDs {
				v.Add(paramActivityTypeIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithOrder returns an option to set the `order` parameter.
func (s *OptionService) WithOrder(order Order) RequestOption {
	return &apiParamOption{
		t: paramOrder,
		checkFunc: func() error {
			if order != OrderAsc && order != OrderDesc {
				msg := fmt.Sprintf("order must be only '%s' or '%s'", string(OrderAsc), string(OrderDesc))
				return newValidationError(msg)
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramOrder.Value(), string(order))
			return nil
		},
	}
}

// WithRoleType returns a option that sets the `roleType` field.
func (s *OptionService) WithRoleType(roleType Role) RequestOption {
	return &apiParamOption{
		t: paramRoleType,
		checkFunc: func() error {
			if roleType < 1 || 6 < roleType {
				return newValidationError("roleType must be between 1 and 6")
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramRoleType.Value(), strconv.Itoa(int(roleType)))
			return nil
		},
	}
}

// WithTextFormattingRule returns a option that sets the `textFormattingRule` field.
func (s *OptionService) WithTextFormattingRule(format Format) RequestOption {
	return &apiParamOption{
		t: paramTextFormattingRule,
		checkFunc: func() error {
			if format != FormatBacklog && format != FormatMarkdown {
				msg := fmt.Sprintf("format must be only '%s' or '%s'", string(FormatBacklog), string(FormatMarkdown))
				return newValidationError(msg)
			}
			return nil
		},
		setFunc: func(v url.Values) error {
			v.Set(paramTextFormattingRule.Value(), string(format))
			return nil
		},
	}
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

// validateOption checks whether the given option key is allowed
// for the current API operation.
func validateOption(optionKey string, validOptions []apiParamOptionType) error {
	for _, valid := range validOptions {
		if optionKey == valid.Value() {
			return nil
		}
	}
	return newInvalidOptionKeyError(optionKey, validOptions)
}

// applyOptions validates and applies request options to the given url.Values.
func applyOptions(v url.Values, validTypes []apiParamOptionType, opts ...RequestOption) error {
	for _, opt := range opts {
		if err := validateOption(opt.Key(), validTypes); err != nil {
			return err
		}
		if err := opt.Check(); err != nil {
			return err
		}
		if err := opt.Set(v); err != nil {
			return err
		}
	}
	return nil
}

// hasRequiredOption checks whether the provided options include at least one of the required form types.
func hasRequiredOption(options []RequestOption, requiredTypes []apiParamOptionType) bool {
	for _, opt := range options {
		for _, requiredType := range requiredTypes {
			if opt.Key() == requiredType.Value() {
				return true
			}
		}
	}
	return false
}
