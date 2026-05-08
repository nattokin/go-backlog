package backlog

// CustomFieldType represents the type identifier of a custom field.
type CustomFieldType int

// CustomFieldType constants for use with ProjectCustomFieldService.Create.
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
