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

// Role defines the user role type within a project.
type Role int

const (
	RoleAdministrator Role = 1
	RoleNormalUser    Role = 2
	RoleReporter      Role = 3
	RoleViewer        Role = 4
	RoleGuestReporter Role = 5
	RoleGuestViewer   Role = 6
)
