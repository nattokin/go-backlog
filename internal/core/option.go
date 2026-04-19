package core

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	ParamActivityTypeIDs                   APIParamOptionType = "activityTypeId[]"
	ParamActualHours                       APIParamOptionType = "actualHours"
	ParamAll                               APIParamOptionType = "all"
	ParamArchived                          APIParamOptionType = "archived"
	ParamAssigneeID                        APIParamOptionType = "assigneeId"
	ParamAssigneeIDs                       APIParamOptionType = "assigneeId[]"
	ParamAttachment                        APIParamOptionType = "attachment"
	ParamAttachmentIDs                     APIParamOptionType = "attachmentId[]"
	ParamCategoryIDs                       APIParamOptionType = "categoryId[]"
	ParamChartEnabled                      APIParamOptionType = "chartEnabled"
	ParamContent                           APIParamOptionType = "content"
	ParamCount                             APIParamOptionType = "count"
	ParamCreatedSince                      APIParamOptionType = "createdSince"
	ParamCreatedUntil                      APIParamOptionType = "createdUntil"
	ParamCreatedUserIDs                    APIParamOptionType = "createdUserId[]"
	ParamDescription                       APIParamOptionType = "description"
	ParamDueDate                           APIParamOptionType = "dueDate"
	ParamDueDateSince                      APIParamOptionType = "dueDateSince"
	ParamDueDateUntil                      APIParamOptionType = "dueDateUntil"
	ParamEstimatedHours                    APIParamOptionType = "estimatedHours"
	ParamHasDueDate                        APIParamOptionType = "hasDueDate"
	ParamIDs                               APIParamOptionType = "id[]"
	ParamIssueTypeID                       APIParamOptionType = "issueTypeId"
	ParamIssueTypeIDs                      APIParamOptionType = "issueTypeId[]"
	ParamKey                               APIParamOptionType = "key"
	ParamKeyword                           APIParamOptionType = "keyword"
	ParamMailAddress                       APIParamOptionType = "mailAddress"
	ParamMailNotify                        APIParamOptionType = "mailNotify"
	ParamMaxID                             APIParamOptionType = "maxId"
	ParamMilestoneIDs                      APIParamOptionType = "milestoneId[]"
	ParamMinID                             APIParamOptionType = "minId"
	ParamName                              APIParamOptionType = "name"
	ParamNotifiedUserIDs                   APIParamOptionType = "notifiedUserId[]"
	ParamOffset                            APIParamOptionType = "offset"
	ParamOrder                             APIParamOptionType = "order"
	ParamParentChild                       APIParamOptionType = "parentChild"
	ParamParentIssueID                     APIParamOptionType = "parentIssueId"
	ParamParentIssueIDs                    APIParamOptionType = "parentIssueId[]"
	ParamPassword                          APIParamOptionType = "password"
	ParamPriorityID                        APIParamOptionType = "priorityId"
	ParamPriorityIDs                       APIParamOptionType = "priorityId[]"
	ParamProjectIDs                        APIParamOptionType = "projectId[]"
	ParamProjectLeaderCanEditProjectLeader APIParamOptionType = "projectLeaderCanEditProjectLeader"
	ParamResolutionID                      APIParamOptionType = "resolutionId"
	ParamResolutionIDs                     APIParamOptionType = "resolutionId[]"
	ParamRoleType                          APIParamOptionType = "roleType"
	ParamSendMail                          APIParamOptionType = "sendMail"
	ParamSharedFile                        APIParamOptionType = "sharedFile"
	ParamSort                              APIParamOptionType = "sort"
	ParamStartDate                         APIParamOptionType = "startDate"
	ParamStartDateSince                    APIParamOptionType = "startDateSince"
	ParamStartDateUntil                    APIParamOptionType = "startDateUntil"
	ParamStatusID                          APIParamOptionType = "statusId"
	ParamStatusIDs                         APIParamOptionType = "statusId[]"
	ParamSubtaskingEnabled                 APIParamOptionType = "subtaskingEnabled"
	ParamSummary                           APIParamOptionType = "summary"
	ParamTextFormattingRule                APIParamOptionType = "textFormattingRule"
	ParamUpdatedSince                      APIParamOptionType = "updatedSince"
	ParamUpdatedUntil                      APIParamOptionType = "updatedUntil"
	ParamUserID                            APIParamOptionType = "userId"
	ParamVersionIDs                        APIParamOptionType = "versionId[]"
)

const MaxActivityTypeID = 26

//
// ──────────────────────────────────────────────────────────────
//  API Option Type
// ──────────────────────────────────────────────────────────────
//

// APIParamOptionType represents the distinct parameter keys for Backlog API requests.
type APIParamOptionType string

// Value returns the string representation of the parameter key for the API request.
func (t APIParamOptionType) Value() string {
	return string(t)
}

//
// ──────────────────────────────────────────────────────────────
//  RequestOption interface
// ──────────────────────────────────────────────────────────────
//

type RequestOption interface {
	Key() string
	Check() error
	Set(url.Values) error
}

//
// ──────────────────────────────────────────────────────────────
//  OptionService
// ──────────────────────────────────────────────────────────────
//

// OptionService provides builders for request options.
// Each XxxOptionService selectively exposes only the valid methods.
type OptionService struct{}

//
// ──────────────────────────────────────────────────────────────
//  APIParamOption
// ──────────────────────────────────────────────────────────────
//

// APIParamOption is the internal implementation of RequestOption.
//
// It encapsulates the parameter type together with optional validation
// and the logic required to apply the value to API request parameters.
// OptionService builder methods return instances of this struct.
type APIParamOption struct {
	Type      APIParamOptionType     // canonical API parameter type
	CheckFunc func() error           // optional validation executed before applying the option
	SetFunc   func(url.Values) error // applies the option value to the provided values
}

// Key returns the API parameter key associated with this option.
func (o *APIParamOption) Key() string {
	return o.Type.Value()
}

// Check validates the option by executing its checkFunc, if defined.
func (o *APIParamOption) Check() error {
	if o.CheckFunc != nil {
		return o.CheckFunc()
	}
	return nil
}

// Set applies the option value to the given url.Values.
func (o *APIParamOption) Set(v url.Values) error {
	if o.SetFunc == nil {
		panic("option has no setter")
	}
	return o.SetFunc(v)
}

//
// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────
//

// ValidateOption checks whether the given option key is allowed
// for the current API operation.
func ValidateOption(optionKey string, validOptions []APIParamOptionType) error {
	for _, valid := range validOptions {
		if optionKey == valid.Value() {
			return nil
		}
	}
	return NewInvalidOptionKeyError(optionKey, validOptions)
}

// ApplyOptions validates and applies request options to the given url.Values.
func ApplyOptions(v url.Values, validTypes []APIParamOptionType, opts ...RequestOption) error {
	for _, opt := range opts {
		if err := ValidateOption(opt.Key(), validTypes); err != nil {
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

// HasRequiredOption checks whether the provided options include at least one of the required form types.
func HasRequiredOption(options []RequestOption, requiredTypes []APIParamOptionType) bool {
	for _, opt := range options {
		for _, requiredType := range requiredTypes {
			if opt.Key() == requiredType.Value() {
				return true
			}
		}
	}
	return false
}

//
// ──────────────────────────────────────────────────────────────
//  SetFunc factories
// ──────────────────────────────────────────────────────────────
//

// setStringFunc returns a SetFunc that calls v.Set with the given string value.
func setStringFunc(key APIParamOptionType, value string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), value)
		return nil
	}
}

// setIntFunc returns a SetFunc that calls v.Set with the int converted to a string.
func setIntFunc(key APIParamOptionType, value int) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), strconv.Itoa(value))
		return nil
	}
}

// setBoolFunc returns a SetFunc that calls v.Set with the bool converted to a string.
func setBoolFunc(key APIParamOptionType, value bool) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), strconv.FormatBool(value))
		return nil
	}
}

// setTimeFunc returns a SetFunc that calls v.Set with the time formatted by the given layout.
func setTimeFunc(key APIParamOptionType, t time.Time, format string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), t.Format(format))
		return nil
	}
}

// addIntFunc returns a SetFunc that calls v.Add for each int in the slice.
func addIntFunc(key APIParamOptionType, values []int) func(url.Values) error {
	return func(v url.Values) error {
		for _, val := range values {
			v.Add(key.Value(), strconv.Itoa(val))
		}
		return nil
	}
}

//
// ──────────────────────────────────────────────────────────────
//  Option builder helpers
// ──────────────────────────────────────────────────────────────
//

// boolOption builds a RequestOption that sets a boolean parameter.
func boolOption(paramType APIParamOptionType, enabled bool) RequestOption {
	return &APIParamOption{
		Type:    paramType,
		SetFunc: setBoolFunc(paramType, enabled),
	}
}

// timeOption builds a RequestOption that formats a time.Time value and sets it.
func timeOption(paramType APIParamOptionType, t time.Time, format string) RequestOption {
	return &APIParamOption{
		Type:    paramType,
		SetFunc: setTimeFunc(paramType, t, format),
	}
}

// nonEmptyStringOption builds a RequestOption that validates the string is not empty and sets it.
func nonEmptyStringOption(paramType APIParamOptionType, value string) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if value == "" {
				return NewValidationError(fmt.Sprintf("%s must not be empty", paramType.Value()))
			}
			return nil
		},
		SetFunc: setStringFunc(paramType, value),
	}
}

// positiveIntOption builds a RequestOption that validates an int is >= 1 and sets it.
func positiveIntOption(paramType APIParamOptionType, value int) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if value < 1 {
				return NewValidationError(fmt.Sprintf("invalid %s: must not be less than 1", paramType.Value()))
			}
			return nil
		},
		SetFunc: setIntFunc(paramType, value),
	}
}

// intRangeOption builds a RequestOption that validates an int is within [min, max] and sets it.
func intRangeOption(paramType APIParamOptionType, value, min, max int) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if value < min || value > max {
				return NewValidationError(fmt.Sprintf("%s must be between %d and %d", paramType.Value(), min, max))
			}
			return nil
		},
		SetFunc: setIntFunc(paramType, value),
	}
}

// validatePositiveInts checks that all values in the slice are >= 1.
// paramName is used in the error message (e.g. "projectId").
func validatePositiveInts(values []int, paramName string) error {
	for _, v := range values {
		if v < 1 {
			return NewValidationError(fmt.Sprintf("invalid %s: %d must not be less than 1", paramName, v))
		}
	}
	return nil
}

// intSliceOption builds a RequestOption that validates and adds multiple ints as repeated query params.
func intSliceOption(paramType APIParamOptionType, paramName string, values []int) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			return validatePositiveInts(values, paramName)
		},
		SetFunc: addIntFunc(paramType, values),
	}
}
