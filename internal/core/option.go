package core

import "net/url"

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

// RequestOption is implemented by all option types that can be applied to an API request.
type RequestOption interface {
	Key() string
	Check() error
	Set(url.Values) error
}

// OptionService provides builder methods for constructing RequestOption values.
// Each XxxOptionService selectively exposes only the valid methods for its API endpoint.
type OptionService struct{}

// APIParamOption is the internal implementation of RequestOption.
//
// It pairs an API parameter key with optional validation (CheckFunc) and
// the logic to write the value into url.Values (SetFunc).
// OptionService builder methods return instances of this struct.
type APIParamOption struct {
	Type      APIParamOptionType     // canonical API parameter key
	CheckFunc func() error           // optional validation executed before applying the option
	SetFunc   func(url.Values) error // applies the value to the request parameters
}

func (o *APIParamOption) Key() string {
	return o.Type.Value()
}

func (o *APIParamOption) Check() error {
	if o.CheckFunc != nil {
		return o.CheckFunc()
	}
	return nil
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
