package model

// Order represents the sort order for API list requests.
type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// Role represents a user role.
type Role int

// Format represents a text formatting rule.
type Format string

const (
	FormatBacklog  Format = "backlog"
	FormatMarkdown Format = "markdown"
)

// IssueSort represents the field to sort issue list results by.
type IssueSort string

const (
	IssueSortIssueType    IssueSort = "issueType"
	IssueSortCategory     IssueSort = "category"
	IssueSortVersion      IssueSort = "version"
	IssueSortMilestone    IssueSort = "milestone"
	IssueSortSummary      IssueSort = "summary"
	IssueSortStatus       IssueSort = "status"
	IssueSortPriority     IssueSort = "priority"
	IssueSortAttachment   IssueSort = "attachment"
	IssueSortSharedFile   IssueSort = "sharedFile"
	IssueSortCreated      IssueSort = "created"
	IssueSortCreatedUser  IssueSort = "createdUser"
	IssueSortUpdated      IssueSort = "updated"
	IssueSortUpdatedUser  IssueSort = "updatedUser"
	IssueSortAssignee     IssueSort = "assignee"
	IssueSortStartDate    IssueSort = "startDate"
	IssueSortDueDate      IssueSort = "dueDate"
	IssueSortEstimatedHours IssueSort = "estimatedHours"
	IssueSortActualHours  IssueSort = "actualHours"
	IssueSortChildIssue   IssueSort = "childIssue"
)
