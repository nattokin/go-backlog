package core

import "fmt"

// WithActivityTypeIDs returns an option to set multiple `activityTypeId[]` parameters.
func (s *OptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return &APIParamOption{
		Type: ParamActivityTypeIDs,
		CheckFunc: func() error {
			for _, id := range typeIDs {
				if err := validateActivityTypeID(id, "activityTypeIds"); err != nil {
					return err
				}
			}
			return nil
		},
		SetFunc: addIntFunc(ParamActivityTypeIDs, typeIDs),
	}
}

// WithApplicableIssueTypeIDs returns an option to set the `applicableIssueTypes[]` parameter.
func (s *OptionService) WithApplicableIssueTypeIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamApplicableIssueTypeIDs, "applicableIssueTypes", ids)
}

// WithAttachmentIDs returns an option to set multiple `attachmentId[]` parameters.
func (s *OptionService) WithAttachmentIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamAttachmentIDs, "attachmentId", ids)
}

// WithItems returns an option to set the `items[]` parameter for List type custom fields.
// Each string becomes a selectable list item and must not be empty.
func (s *OptionService) WithItems(items []string) RequestOption {
	return &APIParamOption{
		Type: ParamItems,
		CheckFunc: func() error {
			for i, item := range items {
				if item == "" {
					return NewValidationError(fmt.Sprintf("items[%d] must not be empty", i))
				}
			}
			return nil
		},
		SetFunc: addStringFunc(ParamItems, items),
	}
}

// WithProjectIDs returns an option to filter by project IDs.
func (s *OptionService) WithProjectIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamProjectIDs, "projectId", ids)
}

// WithIssueTypeIDs returns an option to filter by issue type IDs.
func (s *OptionService) WithIssueTypeIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamIssueTypeIDs, "issueTypeId", ids)
}

// WithCategoryIDs returns an option to filter by category IDs.
func (s *OptionService) WithCategoryIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamCategoryIDs, "categoryId", ids)
}

// WithVersionIDs returns an option to filter by version IDs.
func (s *OptionService) WithVersionIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamVersionIDs, "versionId", ids)
}

// WithMilestoneIDs returns an option to filter by milestone IDs.
func (s *OptionService) WithMilestoneIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamMilestoneIDs, "milestoneId", ids)
}

// WithIssueIDs returns an option to filter by issue IDs.
func (s *OptionService) WithIssueIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamIssueIDs, "issueId", ids)
}

// WithNotifiedUserIDs returns an option to set multiple `notifiedUserId[]` parameters.
func (s *OptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamNotifiedUserIDs, "notifiedUserId", ids)
}

// WithStatusIDs returns an option to filter by status IDs.
func (s *OptionService) WithStatusIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamStatusIDs, "statusId", ids)
}

// WithPriorityIDs returns an option to filter by priority IDs.
func (s *OptionService) WithPriorityIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamPriorityIDs, "priorityId", ids)
}

// WithAssigneeIDs returns an option to filter by assignee user IDs.
func (s *OptionService) WithAssigneeIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamAssigneeIDs, "assigneeId", ids)
}

// WithCreatedUserIDs returns an option to filter by created user IDs.
func (s *OptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamCreatedUserIDs, "createdUserId", ids)
}

// WithResolutionIDs returns an option to filter by resolution IDs.
func (s *OptionService) WithResolutionIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamResolutionIDs, "resolutionId", ids)
}

// WithIDs returns an option to filter by issue IDs.
func (s *OptionService) WithIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamIDs, "id", ids)
}

// WithParentIssueIDs returns an option to filter by parent issue IDs.
func (s *OptionService) WithParentIssueIDs(ids []int) RequestOption {
	return positiveIntSliceOption(ParamParentIssueIDs, "parentIssueId", ids)
}

//
// ──────────────────────────────────────────────────────────────
//  Option builder helpers
// ──────────────────────────────────────────────────────────────
//

// positiveIntSliceOption builds a RequestOption that validates and adds multiple ints as repeated query params.
func positiveIntSliceOption(paramType APIParamOptionType, paramName string, values []int) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			return validatePositiveInts(values, paramName)
		},
		SetFunc: addIntFunc(paramType, values),
	}
}

//
// ──────────────────────────────────────────────────────────────
//  Validation helpers
// ──────────────────────────────────────────────────────────────
//

// validateActivityTypeID ensures that the given activity type ID is within the valid range [1, 26].
func validateActivityTypeID(id int, key string) error {
	if id < 1 || id > MaxActivityTypeID {
		return NewValidationError(fmt.Sprintf("invalid %s: must be between 1 and %d", key, MaxActivityTypeID))
	}
	return nil
}

// validatePositiveInts checks that all values in the slice are >= 1.
// paramName is used in the error message (e.g. "projectId").
func validatePositiveInts(values []int, paramName string) error {
	for _, v := range values {
		if v < 1 {
			return NewValidationError(fmt.Sprintf("invalid %s: %d must not be less than 1", paramName, v))
		}
	}
	return nil
}
