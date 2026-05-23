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
	ParamCustomField                       APIParamOptionType = "customField"
	ParamDescription                       APIParamOptionType = "description"
	ParamDueDate                           APIParamOptionType = "dueDate"
	ParamDueDateSince                      APIParamOptionType = "dueDateSince"
	ParamDueDateUntil                      APIParamOptionType = "dueDateUntil"
	ParamEstimatedHours                    APIParamOptionType = "estimatedHours"
	ParamExcludeGroupMembers               APIParamOptionType = "excludeGroupMembers"
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

// MaxActivityTypeID is the upper bound of valid activity type IDs in the Backlog API.
const MaxActivityTypeID = 26

// APIParamOptionType represents a distinct parameter key for Backlog API requests.
type APIParamOptionType string

func (t APIParamOptionType) Value() string {
	return string(t)
}

// OptionService provides builder methods for constructing *APIParamOption values.
// Each XxxOptionService selectively exposes only the valid methods for its API endpoint.
type OptionService struct{}

// APIParamOption is the option type used for all API request parameters.
//
// It pairs an API parameter key with optional validation (CheckFunc) and
// the logic to write the value into url.Values (SetFunc).
// OptionService builder methods return instances of this struct.
type APIParamOption struct {
	Type      APIParamOptionType      // canonical API parameter key
	CheckFunc func() *ValidationError // optional validation; nil means no validation
	SetFunc   func(url.Values) error  // applies the value to the request parameters
}

func (o *APIParamOption) Key() string {
	return o.Type.Value()
}

// Validate runs the option's validation and returns a *ValidationError if it
// fails, or nil if it passes. Callers within internal packages should use this
// method directly instead of Check().
func (o *APIParamOption) Validate() *ValidationError {
	if o.CheckFunc != nil {
		return o.CheckFunc()
	}
	return nil
}

// Check runs Validate and returns the result as a ValidationResult.
// Returns OK when validation passes (including when CheckFunc is nil).
func (o *APIParamOption) Check() ValidationResult {
	if ve := o.Validate(); ve != nil {
		return ve
	}
	return OK
}

func (o *APIParamOption) Set(v url.Values) error {
	if o.SetFunc == nil {
		panic("option has no setter")
	}
	return o.SetFunc(v)
}

// ValidateOption checks whether the given option key is permitted for the current API operation.
func ValidateOption(optionKey string, validOptions []APIParamOptionType) error {
	for _, valid := range validOptions {
		if optionKey == valid.Value() {
			return nil
		}
	}
	return NewInvalidOptionKeyError(optionKey, validOptions)
}

// ApplyOptions validates and applies request options to the given url.Values.
// Validation errors from Check() are collected into ValidationErrors and returned
// together so callers can inspect all invalid inputs at once.
<<<<<<< Updated upstream
// InvalidOptionKeyError, nil options, and nil ValidationResult are returned immediately.
func ApplyOptions(v url.Values, validTypes []APIParamOptionType, opts ...RequestOption) error {
=======
// InvalidOptionKeyError and nil options are returned immediately.
func ApplyOptions(v url.Values, validTypes []APIParamOptionType, opts ...*APIParamOption) error {
>>>>>>> Stashed changes
	var errs ValidationErrors

	for _, opt := range opts {
		if opt == nil {
			return NewInvalidOptionError("nil option is not allowed")
		}
		if err := ValidateOption(opt.Key(), validTypes); err != nil {
			return err
		}

		result := opt.Check()
		if result == nil {
			return NewInvalidOptionError("Check() must not return nil")
		}
		if !result.Valid() {
			errs = append(errs, NewValidationError(result.Target(), result.Message()))
			continue
		}

		if err := opt.Set(v); err != nil {
			return err
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
