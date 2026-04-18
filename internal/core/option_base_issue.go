package core

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/nattokin/go-backlog/internal/model"
)

const issueDateFormat = "2006-01-02"

// WithProjectIDs returns an option to filter by project IDs.
func (s *OptionService) WithProjectIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamProjectIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid projectId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamProjectIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithIssueTypeIDs returns an option to filter by issue type IDs.
func (s *OptionService) WithIssueTypeIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamIssueTypeIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid issueTypeId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamIssueTypeIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithCategoryIDs returns an option to filter by category IDs.
func (s *OptionService) WithCategoryIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamCategoryIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid categoryId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamCategoryIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithVersionIDs returns an option to filter by version IDs.
func (s *OptionService) WithVersionIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamVersionIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid versionId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamVersionIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithMilestoneIDs returns an option to filter by milestone IDs.
func (s *OptionService) WithMilestoneIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamMilestoneIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid milestoneId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamMilestoneIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithStatusIDs returns an option to filter by status IDs.
func (s *OptionService) WithStatusIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamStatusIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid statusId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamStatusIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithPriorityIDs returns an option to filter by priority IDs.
func (s *OptionService) WithPriorityIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamPriorityIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid priorityId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamPriorityIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithAssigneeIDs returns an option to filter by assignee user IDs.
func (s *OptionService) WithAssigneeIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamAssigneeIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid assigneeId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamAssigneeIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithCreatedUserIDs returns an option to filter by created user IDs.
func (s *OptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamCreatedUserIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid createdUserId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamCreatedUserIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithResolutionIDs returns an option to filter by resolution IDs.
func (s *OptionService) WithResolutionIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamResolutionIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid resolutionId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamResolutionIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithParentChild returns an option to set the `parentChild` parameter.
// 0: All, 1: Exclude Child Issue, 2: Child Issue, 3: Neither Parent nor Child, 4: Parent Issue.
func (s *OptionService) WithParentChild(parentChild int) RequestOption {
	return &APIParamOption{
		Type: ParamParentChild,
		CheckFunc: func() error {
			if parentChild < 0 || parentChild > 4 {
				return NewValidationError("parentChild must be between 0 and 4")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamParentChild.Value(), strconv.Itoa(parentChild))
			return nil
		},
	}
}

// WithAttachment returns an option to include only issues with attachments.
func (s *OptionService) WithAttachment(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamAttachment,
		SetFunc: func(v url.Values) error {
			v.Set(ParamAttachment.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithSharedFile returns an option to include only issues with shared files.
func (s *OptionService) WithSharedFile(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamSharedFile,
		SetFunc: func(v url.Values) error {
			v.Set(ParamSharedFile.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithIssueSort returns an option to set the `sort` parameter for issue list.
func (s *OptionService) WithIssueSort(sort model.IssueSort) RequestOption {
	validSorts := []model.IssueSort{
		model.IssueSortIssueType, model.IssueSortCategory, model.IssueSortVersion,
		model.IssueSortMilestone, model.IssueSortSummary, model.IssueSortStatus,
		model.IssueSortPriority, model.IssueSortAttachment, model.IssueSortSharedFile,
		model.IssueSortCreated, model.IssueSortCreatedUser, model.IssueSortUpdated,
		model.IssueSortUpdatedUser, model.IssueSortAssignee, model.IssueSortStartDate,
		model.IssueSortDueDate, model.IssueSortEstimatedHours, model.IssueSortActualHours,
		model.IssueSortChildIssue,
	}
	return &APIParamOption{
		Type: ParamSort,
		CheckFunc: func() error {
			for _, v := range validSorts {
				if sort == v {
					return nil
				}
			}
			return NewValidationError(fmt.Sprintf("invalid sort value: %q", string(sort)))
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamSort.Value(), string(sort))
			return nil
		},
	}
}

// WithOffset returns an option to set the `offset` parameter.
func (s *OptionService) WithOffset(offset int) RequestOption {
	return &APIParamOption{
		Type: ParamOffset,
		CheckFunc: func() error {
			if offset < 0 {
				return NewValidationError("offset must not be negative")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamOffset.Value(), strconv.Itoa(offset))
			return nil
		},
	}
}

// WithCreatedSince returns an option to filter issues created on or after the given date.
func (s *OptionService) WithCreatedSince(t time.Time) RequestOption {
	return &APIParamOption{
		Type: ParamCreatedSince,
		SetFunc: func(v url.Values) error {
			v.Set(ParamCreatedSince.Value(), t.Format(issueDateFormat))
			return nil
		},
	}
}

// WithCreatedUntil returns an option to filter issues created on or before the given date.
func (s *OptionService) WithCreatedUntil(t time.Time) RequestOption {
	return &APIParamOption{
		Type: ParamCreatedUntil,
		SetFunc: func(v url.Values) error {
			v.Set(ParamCreatedUntil.Value(), t.Format(issueDateFormat))
			return nil
		},
	}
}

// WithUpdatedSince returns an option to filter issues updated on or after the given date.
func (s *OptionService) WithUpdatedSince(t time.Time) RequestOption {
	return &APIParamOption{
		Type: ParamUpdatedSince,
		SetFunc: func(v url.Values) error {
			v.Set(ParamUpdatedSince.Value(), t.Format(issueDateFormat))
			return nil
		},
	}
}

// WithUpdatedUntil returns an option to filter issues updated on or before the given date.
func (s *OptionService) WithUpdatedUntil(t time.Time) RequestOption {
	return &APIParamOption{
		Type: ParamUpdatedUntil,
		SetFunc: func(v url.Values) error {
			v.Set(ParamUpdatedUntil.Value(), t.Format(issueDateFormat))
			return nil
		},
	}
}

// WithStartDateSince returns an option to filter issues with a start date on or after the given date.
func (s *OptionService) WithStartDateSince(t time.Time) RequestOption {
	return &APIParamOption{
		Type: ParamStartDateSince,
		SetFunc: func(v url.Values) error {
			v.Set(ParamStartDateSince.Value(), t.Format(issueDateFormat))
			return nil
		},
	}
}

// WithStartDateUntil returns an option to filter issues with a start date on or before the given date.
func (s *OptionService) WithStartDateUntil(t time.Time) RequestOption {
	return &APIParamOption{
		Type: ParamStartDateUntil,
		SetFunc: func(v url.Values) error {
			v.Set(ParamStartDateUntil.Value(), t.Format(issueDateFormat))
			return nil
		},
	}
}

// WithDueDateSince returns an option to filter issues with a due date on or after the given date.
func (s *OptionService) WithDueDateSince(t time.Time) RequestOption {
	return &APIParamOption{
		Type: ParamDueDateSince,
		SetFunc: func(v url.Values) error {
			v.Set(ParamDueDateSince.Value(), t.Format(issueDateFormat))
			return nil
		},
	}
}

// WithDueDateUntil returns an option to filter issues with a due date on or before the given date.
func (s *OptionService) WithDueDateUntil(t time.Time) RequestOption {
	return &APIParamOption{
		Type: ParamDueDateUntil,
		SetFunc: func(v url.Values) error {
			v.Set(ParamDueDateUntil.Value(), t.Format(issueDateFormat))
			return nil
		},
	}
}

// WithHasDueDate returns an option to exclude issues without a due date.
// Note: Setting this to true is not supported by the Backlog API and will result in an error.
func (s *OptionService) WithHasDueDate(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamHasDueDate,
		SetFunc: func(v url.Values) error {
			v.Set(ParamHasDueDate.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithIDs returns an option to filter by issue IDs.
func (s *OptionService) WithIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid id: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithParentIssueIDs returns an option to filter by parent issue IDs.
func (s *OptionService) WithParentIssueIDs(ids []int) RequestOption {
	return &APIParamOption{
		Type: ParamParentIssueIDs,
		CheckFunc: func() error {
			for _, id := range ids {
				if id < 1 {
					return NewValidationError(fmt.Sprintf("invalid parentIssueId: %d must not be less than 1", id))
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range ids {
				v.Add(ParamParentIssueIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}
