package core

import (
	"fmt"
	"net/url"
	"strconv"
)

func (s *OptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return &APIParamOption{
		Type: ParamActivityTypeIDs,
		CheckFunc: func() ValidationResult {
			for _, id := range typeIDs {
				if result := validateActivityTypeID(id, "activityTypeIds"); !result.Valid() {
					return result
				}
			}
			return OK
		},
		SetFunc: addIntFunc(ParamActivityTypeIDs, typeIDs),
	}
}

func (s *OptionService) WithApplicableIssueTypeIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamApplicableIssueTypeIDs, "applicableIssueTypes", ids)
}

func (s *OptionService) WithAttachmentIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamAttachmentIDs, "attachmentId", ids)
}

// WithItems sets `items[]` for List type custom fields.
// Each string becomes a selectable list item and must not be empty.
func (s *OptionService) WithItems(items []string) RequestOption {
	return &APIParamOption{
		Type: ParamItems,
		CheckFunc: func() ValidationResult {
			for i, item := range items {
				if item == "" {
					return NewValidationError(ParamItems.Value(), fmt.Sprintf("items[%d] must not be empty", i))
				}
			}
			return OK
		},
		SetFunc: addStringFunc(ParamItems, items),
	}
}

func (s *OptionService) WithProjectIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamProjectIDs, "projectId", ids)
}

func (s *OptionService) WithIssueTypeIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamIssueTypeIDs, "issueTypeId", ids)
}

func (s *OptionService) WithCategoryIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamCategoryIDs, "categoryId", ids)
}

func (s *OptionService) WithVersionIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamVersionIDs, "versionId", ids)
}

func (s *OptionService) WithMilestoneIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamMilestoneIDs, "milestoneId", ids)
}

func (s *OptionService) WithIssueIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamIssueIDs, "issueId", ids)
}

func (s *OptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamNotifiedUserIDs, "notifiedUserId", ids)
}

func (s *OptionService) WithStatusIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamStatusIDs, "statusId", ids)
}

func (s *OptionService) WithPriorityIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamPriorityIDs, "priorityId", ids)
}

func (s *OptionService) WithAssigneeIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamAssigneeIDs, "assigneeId", ids)
}

func (s *OptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamCreatedUserIDs, "createdUserId", ids)
}

func (s *OptionService) WithResolutionIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamResolutionIDs, "resolutionId", ids)
}

func (s *OptionService) WithIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamIDs, "id", ids)
}

func (s *OptionService) WithParentIssueIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamParentIssueIDs, "parentIssueId", ids)
}

// positiveIntSliceOption builds a RequestOption that validates and adds multiple ints as repeated query params.
func positiveIntSliceOption(paramType APIParamOptionType, paramName string, values []int) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() ValidationResult {
			return validatePositiveInts(values, paramName)
		},
		SetFunc: addIntFunc(paramType, values),
	}
}

// addIntFunc returns a SetFunc that calls v.Add for each int in the slice.
func addIntFunc(key APIParamOptionType, values []int) func(url.Values) error {
	return func(v url.Values) error {
		for _, val := range values {
			v.Add(key.Value(), strconv.Itoa(val))
		}
		return nil
	}
}

// addStringFunc returns a SetFunc that calls v.Add for each string in the slice.
func addStringFunc(key APIParamOptionType, values []string) func(url.Values) error {
	return func(v url.Values) error {
		for _, val := range values {
			v.Add(key.Value(), val)
		}
		return nil
	}
}

// validateActivityTypeID ensures the ID is within the valid range [1, 26].
func validateActivityTypeID(id int, key string) ValidationResult {
	if id < 1 || id > MaxActivityTypeID {
		return NewValidationError(key, fmt.Sprintf("invalid %s: must be between 1 and %d", key, MaxActivityTypeID))
	}
	return OK
}

// validatePositiveInts checks that all values in the slice are >= 1.
// paramName is used in the error message (e.g. "projectId").
func validatePositiveInts(values []int, paramName string) ValidationResult {
	for _, v := range values {
		if v < 1 {
			return NewValidationError(paramName, fmt.Sprintf("invalid %s: %d must not be less than 1", paramName, v))
		}
	}
	return OK
}
