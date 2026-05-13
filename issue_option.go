package backlog

import (
	"time"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
)

// IssueOptionService provides a domain-specific set of option builders
// for operations within the IssueService.
type IssueOptionService struct {
	base *core.OptionService
}

// WithActualHours returns an option to set the `actualHours` parameter.
func (s *IssueOptionService) WithActualHours(hours float64) RequestOption {
	return s.base.WithActualHours(hours)
}

// WithAssigneeID returns an option to set the `assigneeId` parameter.
func (s *IssueOptionService) WithAssigneeID(id int) RequestOption {
	return s.base.WithAssigneeID(id)
}

// WithAssigneeIDs filters issues by assignee user IDs.
func (s *IssueOptionService) WithAssigneeIDs(ids []int) RequestOption {
	return s.base.WithAssigneeIDs(ids)
}

// WithAttachment filters to include only issues with attachments.
func (s *IssueOptionService) WithAttachment(enabled bool) RequestOption {
	return s.base.WithAttachment(enabled)
}

// WithAttachmentIDs returns an option to set multiple `attachmentId[]` parameters.
func (s *IssueOptionService) WithAttachmentIDs(ids []int) RequestOption {
	return s.base.WithAttachmentIDs(ids)
}

// WithCategoryIDs filters issues by category IDs.
func (s *IssueOptionService) WithCategoryIDs(ids []int) RequestOption {
	return s.base.WithCategoryIDs(ids)
}

// WithComment returns an option to set the `comment` parameter.
func (s *IssueOptionService) WithComment(comment string) RequestOption {
	return s.base.WithComment(comment)
}

// WithCount sets the number of issues to retrieve (1-100).
func (s *IssueOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithCreatedSince filters issues created on or after the given date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithCreatedSince(date string) RequestOption {
	return s.base.WithCreatedSince(date)
}

// WithCreatedUntil filters issues created on or before the given date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithCreatedUntil(date string) RequestOption {
	return s.base.WithCreatedUntil(date)
}

// WithCreatedUserIDs filters issues by created user IDs.
func (s *IssueOptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return s.base.WithCreatedUserIDs(ids)
}

// WithCustomFieldItems returns an option to set predefined item selections for a
// list-type custom field (Single list, Multiple list, Checkbox, Radio).
//
// The parameter name is dynamically generated as "customField_{id}". Multiple
// item IDs are sent as repeated values under the same key, which the Backlog API
// interprets as a list.
//
// Returns an error if id is less than 1.
func (s *IssueOptionService) WithCustomFieldItems(id int, itemIDs []int) RequestOption {
	return issue.WithCustomFieldItems(id, itemIDs)
}

// WithCustomFieldNum returns an option to set a Number type custom field value.
//
// The parameter name is dynamically generated as "customField_{id}".
// Both integer and fractional values are supported (e.g. 1.0, -3.5).
//
// Returns an error if id is less than 1.
func (s *IssueOptionService) WithCustomFieldNum(id int, value float64) RequestOption {
	return issue.WithCustomField(id, value)
}

// WithCustomFieldOther returns an option to set the free-text "Other" value for a
// list-type custom field where allowInput is enabled.
//
// The parameter name is dynamically generated as "customField_{id}_otherValue".
//
// Returns an error if id is less than 1.
func (s *IssueOptionService) WithCustomFieldOther(id int, value string) RequestOption {
	return issue.WithCustomFieldOther(id, value)
}

// WithCustomFieldString returns an option to set a Text or Sentence type custom
// field value.
//
// The parameter name is dynamically generated as "customField_{id}".
//
// Returns an error if id is less than 1 or value is empty.
func (s *IssueOptionService) WithCustomFieldString(id int, value string) RequestOption {
	return issue.WithCustomField(id, value)
}

// WithCustomFieldTime returns an option to set a Date type custom field value.
//
// The parameter name is dynamically generated as "customField_{id}".
// The value is formatted as "yyyy-MM-dd".
//
// Returns an error if id is less than 1 or value is a zero time.Time.
func (s *IssueOptionService) WithCustomFieldTime(id int, value time.Time) RequestOption {
	return issue.WithCustomField(id, value)
}

// WithDescription returns an option to set the `description` parameter.
func (s *IssueOptionService) WithDescription(description string) RequestOption {
	return s.base.WithDescription(description)
}

// WithDueDate returns an option to set the `dueDate` parameter.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithDueDate(date string) RequestOption {
	return s.base.WithDueDate(date)
}

// WithDueDateSince filters issues with a due date on or after the given date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithDueDateSince(date string) RequestOption {
	return s.base.WithDueDateSince(date)
}

// WithDueDateUntil filters issues with a due date on or before the given date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithDueDateUntil(date string) RequestOption {
	return s.base.WithDueDateUntil(date)
}

// WithEstimatedHours returns an option to set the `estimatedHours` parameter.
func (s *IssueOptionService) WithEstimatedHours(hours float64) RequestOption {
	return s.base.WithEstimatedHours(hours)
}

// WithHasDueDate filters to exclude issues without a due date.
// Note: Setting this to true is not supported by the Backlog API and will result in an error.
func (s *IssueOptionService) WithHasDueDate(enabled bool) RequestOption {
	return s.base.WithHasDueDate(enabled)
}

// WithIDs filters issues by issue IDs.
func (s *IssueOptionService) WithIDs(ids []int) RequestOption {
	return s.base.WithIDs(ids)
}

// WithIssueSort sets the field to sort issue list results by.
func (s *IssueOptionService) WithIssueSort(sort IssueSort) RequestOption {
	return s.base.WithIssueSort(string(sort))
}

// WithIssueTypeID returns an option to set the `issueTypeId` parameter.
func (s *IssueOptionService) WithIssueTypeID(id int) RequestOption {
	return s.base.WithIssueTypeID(id)
}

// WithIssueTypeIDs filters issues by issue type IDs.
func (s *IssueOptionService) WithIssueTypeIDs(ids []int) RequestOption {
	return s.base.WithIssueTypeIDs(ids)
}

// WithKeyword filters issues by keyword.
func (s *IssueOptionService) WithKeyword(keyword string) RequestOption {
	return s.base.WithKeyword(keyword)
}

// WithMilestoneIDs filters issues by milestone IDs.
func (s *IssueOptionService) WithMilestoneIDs(ids []int) RequestOption {
	return s.base.WithMilestoneIDs(ids)
}

// WithNotifiedUserIDs returns an option to set multiple `notifiedUserId[]` parameters.
func (s *IssueOptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return s.base.WithNotifiedUserIDs(ids)
}

// WithOffset sets the number of issues to skip.
func (s *IssueOptionService) WithOffset(offset int) RequestOption {
	return s.base.WithOffset(offset)
}

// WithOrder sets the sort order of results.
func (s *IssueOptionService) WithOrder(order Order) RequestOption {
	return s.base.WithOrder(string(order))
}

// WithParentChild filters issues by subtask relationship.
// 0: All, 1: Exclude Child Issue, 2: Child Issue, 3: Neither Parent nor Child, 4: Parent Issue.
func (s *IssueOptionService) WithParentChild(parentChild int) RequestOption {
	return s.base.WithParentChild(parentChild)
}

// WithParentIssueID returns an option to set the `parentIssueId` parameter.
func (s *IssueOptionService) WithParentIssueID(id int) RequestOption {
	return s.base.WithParentIssueID(id)
}

// WithParentIssueIDs filters issues by parent issue IDs.
func (s *IssueOptionService) WithParentIssueIDs(ids []int) RequestOption {
	return s.base.WithParentIssueIDs(ids)
}

// WithPriorityID returns an option to set the `priorityId` parameter.
func (s *IssueOptionService) WithPriorityID(id int) RequestOption {
	return s.base.WithPriorityID(id)
}

// WithPriorityIDs filters issues by priority IDs.
func (s *IssueOptionService) WithPriorityIDs(ids []int) RequestOption {
	return s.base.WithPriorityIDs(ids)
}

// WithProjectIDs filters issues by project IDs.
func (s *IssueOptionService) WithProjectIDs(ids []int) RequestOption {
	return s.base.WithProjectIDs(ids)
}

// WithResolutionID returns an option to set the `resolutionId` parameter.
func (s *IssueOptionService) WithResolutionID(id int) RequestOption {
	return s.base.WithResolutionID(id)
}

// WithResolutionIDs filters issues by resolution IDs.
func (s *IssueOptionService) WithResolutionIDs(ids []int) RequestOption {
	return s.base.WithResolutionIDs(ids)
}

// WithSharedFile filters to include only issues with shared files.
func (s *IssueOptionService) WithSharedFile(enabled bool) RequestOption {
	return s.base.WithSharedFile(enabled)
}

// WithStartDate returns an option to set the `startDate` parameter.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithStartDate(date string) RequestOption {
	return s.base.WithStartDate(date)
}

// WithStartDateSince filters issues with a start date on or after the given date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithStartDateSince(date string) RequestOption {
	return s.base.WithStartDateSince(date)
}

// WithStartDateUntil filters issues with a start date on or before the given date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithStartDateUntil(date string) RequestOption {
	return s.base.WithStartDateUntil(date)
}

// WithStatusID returns an option to set the `statusId` parameter.
func (s *IssueOptionService) WithStatusID(id int) RequestOption {
	return s.base.WithStatusID(id)
}

// WithStatusIDs filters issues by status IDs.
func (s *IssueOptionService) WithStatusIDs(ids []int) RequestOption {
	return s.base.WithStatusIDs(ids)
}

// WithSummary returns an option to set the `summary` parameter.
func (s *IssueOptionService) WithSummary(summary string) RequestOption {
	return s.base.WithSummary(summary)
}

// WithUpdatedSince filters issues updated on or after the given date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithUpdatedSince(date string) RequestOption {
	return s.base.WithUpdatedSince(date)
}

// WithUpdatedUntil filters issues updated on or before the given date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *IssueOptionService) WithUpdatedUntil(date string) RequestOption {
	return s.base.WithUpdatedUntil(date)
}

// WithVersionIDs filters issues by version IDs.
func (s *IssueOptionService) WithVersionIDs(ids []int) RequestOption {
	return s.base.WithVersionIDs(ids)
}

func newIssueOptionService(option *core.OptionService) *IssueOptionService {
	return &IssueOptionService{base: option}
}
