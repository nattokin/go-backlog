package core

import (
	"net/url"
)

const (
	ParamActivityTypeIDs                   APIParamOptionType = "activityTypeId[]"
	ParamActualHours                       APIParamOptionType = "actualHours"
	ParamAll                               APIParamOptionType = "all"
	ParamAllEvent                          APIParamOptionType = "allEvent"
	ParamAllowAddItem                      APIParamOptionType = "allowAddItem"
	ParamAllowInput                        APIParamOptionType = "allowInput"
	ParamApplicableIssueTypeIDs            APIParamOptionType = "applicableIssueTypes[]"
	ParamArchived                          APIParamOptionType = "archived"
	ParamAssigneeID                        APIParamOptionType = "assigneeId"
	ParamAssigneeIDs                       APIParamOptionType = "assigneeId[]"
	ParamAttachment                        APIParamOptionType = "attachment"
	ParamAttachmentIDs                     APIParamOptionType = "attachmentId[]"
	ParamBase                              APIParamOptionType = "base"
	ParamBranch                            APIParamOptionType = "branch"
	ParamCategoryIDs                       APIParamOptionType = "categoryId[]"
	ParamChartEnabled                      APIParamOptionType = "chartEnabled"
	ParamColor                             APIParamOptionType = "color"
	ParamComment                           APIParamOptionType = "comment"
	ParamCommentID                         APIParamOptionType = "commentId"
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
	ParamHookURL                           APIParamOptionType = "hookUrl"
	ParamIDs                               APIParamOptionType = "id[]"
	ParamInitialDate                       APIParamOptionType = "initialDate"
	ParamInitialShift                      APIParamOptionType = "initialShift"
	ParamInitialValue                      APIParamOptionType = "initialValue"
	ParamInitialValueType                  APIParamOptionType = "initialValueType"
	ParamIssueID                           APIParamOptionType = "issueId"
	ParamIssueIDs                          APIParamOptionType = "issueId[]"
	ParamIssueTypeID                       APIParamOptionType = "issueTypeId"
	ParamIssueTypeIDs                      APIParamOptionType = "issueTypeId[]"
	ParamItems                             APIParamOptionType = "items[]"
	ParamKey                               APIParamOptionType = "key"
	ParamKeyword                           APIParamOptionType = "keyword"
	ParamMailAddress                       APIParamOptionType = "mailAddress"
	ParamMailNotify                        APIParamOptionType = "mailNotify"
	ParamMax                               APIParamOptionType = "max"
	ParamMaxID                             APIParamOptionType = "maxId"
	ParamMilestoneIDs                      APIParamOptionType = "milestoneId[]"
	ParamMin                               APIParamOptionType = "min"
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
	ParamPullRequestCommentID              APIParamOptionType = "pullRequestCommentId"
	ParamPullRequestID                     APIParamOptionType = "pullRequestId"
	ParamReleaseDueDate                    APIParamOptionType = "releaseDueDate"
	ParamRequired                          APIParamOptionType = "required"
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
	ParamTemplateDescription               APIParamOptionType = "templateDescription"
	ParamTemplateSummary                   APIParamOptionType = "templateSummary"
	ParamTextFormattingRule                APIParamOptionType = "textFormattingRule"
	ParamTypeID                            APIParamOptionType = "typeId"
	ParamUnit                              APIParamOptionType = "unit"
	ParamUpdatedSince                      APIParamOptionType = "updatedSince"
	ParamUpdatedUntil                      APIParamOptionType = "updatedUntil"
	ParamUserID                            APIParamOptionType = "userId"
	ParamVersionIDs                        APIParamOptionType = "versionId[]"
	ParamWikiID                            APIParamOptionType = "wikiId"
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
		if opt == nil {
			return NewValidationError("nil option is not allowed")
		}
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
