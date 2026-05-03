package backlog

import "github.com/nattokin/go-backlog/internal/model"

// CustomFieldType represents the type identifier of a custom field.
type CustomFieldType = model.CustomFieldType

// CustomFieldType constants for use with ProjectCustomFieldService.Create.
const (
	CustomFieldTypeText         CustomFieldType = model.CustomFieldTypeText
	CustomFieldTypeSentence     CustomFieldType = model.CustomFieldTypeSentence
	CustomFieldTypeNumber       CustomFieldType = model.CustomFieldTypeNumber
	CustomFieldTypeDate         CustomFieldType = model.CustomFieldTypeDate
	CustomFieldTypeSingleList   CustomFieldType = model.CustomFieldTypeSingleList
	CustomFieldTypeMultipleList CustomFieldType = model.CustomFieldTypeMultipleList
	CustomFieldTypeCheckbox     CustomFieldType = model.CustomFieldTypeCheckbox
	CustomFieldTypeRadio        CustomFieldType = model.CustomFieldTypeRadio
)
