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
// This type is used for JSON deserialization of Space.TextFormattingRule and
// Project.TextFormattingRule, so it is kept as a named type.
type Format string

const (
	FormatMarkdown Format = "markdown"
	FormatBacklog  Format = "backlog"
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
