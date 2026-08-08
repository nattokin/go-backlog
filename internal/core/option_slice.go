package core

import (
	"fmt"
	"net/url"
	"strconv"
)

func (s *OptionService) WithActivityTypeIDs(typeIDs []int) *APIParamOption {
	return &APIParamOption{
		Type: ParamActivityTypeIDs,
		CheckFunc: func() *ValidationError {
			for _, id := range typeIDs {
				if ve := validateActivityTypeID(id, "activityTypeIds"); ve != nil {
					return ve
				}
			}
			return nil
		},
		SetFunc: addIntFunc(ParamActivityTypeIDs, typeIDs),
	}
}

func (s *OptionService) WithApplicableIssueTypeIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamApplicableIssueTypeIDs, "applicableIssueTypes", ids)
}

func (s *OptionService) WithAttachmentIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamAttachmentIDs, "attachmentId", ids)
}

// WithItems sets `items[]` for List type custom fields.
func (s *OptionService) WithItems(items []string) *APIParamOption {
	return &APIParamOption{
		Type: ParamItems,
		CheckFunc: func() *ValidationError {
			for i, item := range items {
				if item == "" {
					return NewValidationError(ParamItems.Value(), fmt.Sprintf("items[%d] must not be empty", i))
				}
			}
			return nil
		},
		SetFunc: addStringFunc(ParamItems, items),
	}
}

func (s *OptionService) WithProjectIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamProjectIDs, "projectId", ids)
}

func (s *OptionService) WithIssueTypeIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamIssueTypeIDs, "issueTypeId", ids)
}

func (s *OptionService) WithCategoryIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamCategoryIDs, "categoryId", ids)
}

func (s *OptionService) WithVersionIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamVersionIDs, "versionId", ids)
}

func (s *OptionService) WithMilestoneIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamMilestoneIDs, "milestoneId", ids)
}

func (s *OptionService) WithIssueIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamIssueIDs, "issueId", ids)
}

func (s *OptionService) WithNotifiedUserIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamNotifiedUserIDs, "notifiedUserId", ids)
}

func (s *OptionService) WithStatusIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamStatusIDs, "statusId", ids)
}

func (s *OptionService) WithPriorityIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamPriorityIDs, "priorityId", ids)
}

func (s *OptionService) WithAssigneeIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamAssigneeIDs, "assigneeId", ids)
}

func (s *OptionService) WithCreatedUserIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamCreatedUserIDs, "createdUserId", ids)
}

func (s *OptionService) WithResolutionIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamResolutionIDs, "resolutionId", ids)
}

func (s *OptionService) WithIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamIDs, "id", ids)
}

func (s *OptionService) WithParentIssueIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamParentIssueIDs, "parentIssueId", ids)
}

func positiveIntSliceOption(paramType APIParamOptionType, paramName string, values []int) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *ValidationError {
			return validatePositiveInts(values, paramName)
		},
		SetFunc: addIntFunc(paramType, values),
	}
}

func addIntFunc(key APIParamOptionType, values []int) func(url.Values) error {
	return func(v url.Values) error {
		for _, val := range values {
			v.Add(key.Value(), strconv.Itoa(val))
		}
		return nil
	}
}

func addStringFunc(key APIParamOptionType, values []string) func(url.Values) error {
	return func(v url.Values) error {
		for _, val := range values {
			v.Add(key.Value(), val)
		}
		return nil
	}
}

func validateActivityTypeID(id int, key string) *ValidationError {
	if id < 1 || id > MaxActivityTypeID {
		return NewValidationError(key, fmt.Sprintf("invalid %s: must be between 1 and %d", key, MaxActivityTypeID))
	}
	return nil
}

func validatePositiveInts(values []int, paramName string) *ValidationError {
	for _, v := range values {
		if v < 1 {
			return NewValidationError(paramName, fmt.Sprintf("invalid %s: %d must not be less than 1", paramName, v))
		}
	}
	return nil
}
