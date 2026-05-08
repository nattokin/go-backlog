package model

// CustomFieldType represents the type identifier of a custom field.
type CustomFieldType int

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

// Format defines the text formatting rule for a Backlog space or project.
type Format string

const (
	FormatMarkdown Format = "markdown"
	FormatBacklog  Format = "backlog"
)

// IssueSort represents the field to sort issue list results by.
type IssueSort string

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

// Order defines the sort order for list results.
type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// Role defines the user role type within a project.
type Role int

const (
	_ Role = iota
	RoleAdministrator
	RoleNormalUser
	RoleReporter
	RoleViewer
	RoleGuestReporter
	RoleGuestViewer
)
