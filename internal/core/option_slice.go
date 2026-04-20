package core

import "fmt"

// WithActivityTypeIDs returns an option to set multiple `activityTypeId[]` parameters.
func (s *OptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return &APIParamOption{
		Type: ParamActivityTypeIDs,
		CheckFunc: func() error {
			for _, id := range typeIDs {
				if err := validateActivityID(id, "activityTypeIds"); err != nil {
					return err
				}
			}
			return nil
		},
		SetFunc: addIntFunc(ParamActivityTypeIDs, typeIDs),
	}
}

// WithAttachmentIDs returns an option to set multiple `attachmentId[]` parameters.
func (s *OptionService) WithAttachmentIDs(ids []int) RequestOption {
	return intSliceOption(ParamAttachmentIDs, "attachmentId", ids)
}

// WithProjectIDs returns an option to filter by project IDs.
func (s *OptionService) WithProjectIDs(ids []int) RequestOption {
	return intSliceOption(ParamProjectIDs, "projectId", ids)
}

// WithIssueTypeIDs returns an option to filter by issue type IDs.
func (s *OptionService) WithIssueTypeIDs(ids []int) RequestOption {
	return intSliceOption(ParamIssueTypeIDs, "issueTypeId", ids)
}

// WithCategoryIDs returns an option to filter by category IDs.
func (s *OptionService) WithCategoryIDs(ids []int) RequestOption {
	return intSliceOption(ParamCategoryIDs, "categoryId", ids)
}

// WithVersionIDs returns an option to filter by version IDs.
func (s *OptionService) WithVersionIDs(ids []int) RequestOption {
	return intSliceOption(ParamVersionIDs, "versionId", ids)
}

// WithMilestoneIDs returns an option to filter by milestone IDs.
func (s *OptionService) WithMilestoneIDs(ids []int) RequestOption {
	return intSliceOption(ParamMilestoneIDs, "milestoneId", ids)
}

// WithIssueIDs returns an option to filter by issue IDs.
func (s *OptionService) WithIssueIDs(ids []int) RequestOption {
	return intSliceOption(ParamIssueIDs, "issueId", ids)
}

// WithNotifiedUserIDs returns an option to set multiple `notifiedUserId[]` parameters.
func (s *OptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return intSliceOption(ParamNotifiedUserIDs, "notifiedUserId", ids)
}

// WithStatusIDs returns an option to filter by status IDs.
func (s *OptionService) WithStatusIDs(ids []int) RequestOption {
	return intSliceOption(ParamStatusIDs, "statusId", ids)
}

// WithPriorityIDs returns an option to filter by priority IDs.
func (s *OptionService) WithPriorityIDs(ids []int) RequestOption {
	return intSliceOption(ParamPriorityIDs, "priorityId", ids)
}

// WithAssigneeIDs returns an option to filter by assignee user IDs.
func (s *OptionService) WithAssigneeIDs(ids []int) RequestOption {
	return intSliceOption(ParamAssigneeIDs, "assigneeId", ids)
}

// WithCreatedUserIDs returns an option to filter by created user IDs.
func (s *OptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return intSliceOption(ParamCreatedUserIDs, "createdUserId", ids)
}

// WithResolutionIDs returns an option to filter by resolution IDs.
func (s *OptionService) WithResolutionIDs(ids []int) RequestOption {
	return intSliceOption(ParamResolutionIDs, "resolutionId", ids)
}

// WithIDs returns an option to filter by issue IDs.
func (s *OptionService) WithIDs(ids []int) RequestOption {
	return intSliceOption(ParamIDs, "id", ids)
}

// WithParentIssueIDs returns an option to filter by parent issue IDs.
func (s *OptionService) WithParentIssueIDs(ids []int) RequestOption {
	return intSliceOption(ParamParentIssueIDs, "parentIssueId", ids)
}

// validateActivityID ensures that the given activity ID is within the valid range [1, 26].
func validateActivityID(id int, key string) error {
	if id < 1 || id > MaxActivityTypeID {
		return NewValidationError(fmt.Sprintf("invalid %s: must be between 1 and %d", key, MaxActivityTypeID))
	}
	return nil
}
