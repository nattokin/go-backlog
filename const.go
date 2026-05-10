package backlog

import "fmt"

// CustomFieldType represents the type identifier of a custom field.
type CustomFieldType int

// CustomFieldType constants for use with [ProjectCustomFieldService.Create].
const (
	CustomFieldTypeText         CustomFieldType = 1
	CustomFieldTypeSentence     CustomFieldType = 2
	CustomFieldTypeNumber       CustomFieldType = 3
	CustomFieldTypeDate         CustomFieldType = 4
	CustomFieldTypeSingleList   CustomFieldType = 5
	CustomFieldTypeMultipleList CustomFieldType = 6
	CustomFieldTypeCheckbox     CustomFieldType = 7
	CustomFieldTypeRadio        CustomFieldType = 8
)

// Format defines the text formatting rule for the Backlog wiki.
type Format string

// Available text formatting rules.
const (
	FormatMarkdown Format = "markdown"
	FormatBacklog  Format = "backlog"
)

func (f Format) String() string {
	switch f {
	case FormatMarkdown:
		return "Markdown"
	case FormatBacklog:
		return "Backlog"
	default:
		return fmt.Sprintf("unknown Format type %s", string(f))
	}
}

// IssueSort defines the field to sort issue list results by.
type IssueSort string

// Available sort fields for issue list operations.
const (
	IssueSortIssueType      IssueSort = "issueType"
	IssueSortCategory       IssueSort = "category"
	IssueSortVersion        IssueSort = "version"
	IssueSortMilestone      IssueSort = "milestone"
	IssueSortSummary        IssueSort = "summary"
	IssueSortStatus         IssueSort = "status"
	IssueSortPriority       IssueSort = "priority"
	IssueSortAttachment     IssueSort = "attachment"
	IssueSortSharedFile     IssueSort = "sharedFile"
	IssueSortCreated        IssueSort = "created"
	IssueSortCreatedUser    IssueSort = "createdUser"
	IssueSortUpdated        IssueSort = "updated"
	IssueSortUpdatedUser    IssueSort = "updatedUser"
	IssueSortAssignee       IssueSort = "assignee"
	IssueSortStartDate      IssueSort = "startDate"
	IssueSortDueDate        IssueSort = "dueDate"
	IssueSortEstimatedHours IssueSort = "estimatedHours"
	IssueSortActualHours    IssueSort = "actualHours"
	IssueSortChildIssue     IssueSort = "childIssue"
)

// Order defines the sort order (ascending or descending).
type Order string

// Available sort orders.
const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

func (o Order) String() string {
	switch o {
	case OrderAsc:
		return "Asc"
	case OrderDesc:
		return "Desc"
	default:
		return fmt.Sprintf("unknown Order type %s", string(o))
	}
}

// Role defines the type of user role within a project.
type Role int

// Available user roles within a project.
const (
	RoleAdministrator Role = 1
	RoleNormalUser    Role = 2
	RoleReporter      Role = 3
	RoleViewer        Role = 4
	RoleGuestReporter Role = 5
	RoleGuestViewer   Role = 6
)

func (r Role) String() string {
	switch r {
	case RoleAdministrator:
		return "Administrator"
	case RoleNormalUser:
		return "NormalUser"
	case RoleReporter:
		return "Reporter"
	case RoleViewer:
		return "Viewer"
	case RoleGuestReporter:
		return "GuestReporter"
	case RoleGuestViewer:
		return "GuestViewer"
	default:
		return fmt.Sprintf("unknown Role type %d", r)
	}
}
