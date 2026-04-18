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
	return idSliceOption(ParamProjectIDs, "projectId", ids)
}

// WithIssueTypeIDs returns an option to filter by issue type IDs.
func (s *OptionService) WithIssueTypeIDs(ids []int) RequestOption {
	return idSliceOption(ParamIssueTypeIDs, "issueTypeId", ids)
}

// WithCategoryIDs returns an option to filter by category IDs.
func (s *OptionService) WithCategoryIDs(ids []int) RequestOption {
	return idSliceOption(ParamCategoryIDs, "categoryId", ids)
}

// WithVersionIDs returns an option to filter by version IDs.
func (s *OptionService) WithVersionIDs(ids []int) RequestOption {
	return idSliceOption(ParamVersionIDs, "versionId", ids)
}

// WithMilestoneIDs returns an option to filter by milestone IDs.
func (s *OptionService) WithMilestoneIDs(ids []int) RequestOption {
	return idSliceOption(ParamMilestoneIDs, "milestoneId", ids)
}

// WithStatusIDs returns an option to filter by status IDs.
func (s *OptionService) WithStatusIDs(ids []int) RequestOption {
	return idSliceOption(ParamStatusIDs, "statusId", ids)
}

// WithPriorityIDs returns an option to filter by priority IDs.
func (s *OptionService) WithPriorityIDs(ids []int) RequestOption {
	return idSliceOption(ParamPriorityIDs, "priorityId", ids)
}

// WithAssigneeIDs returns an option to filter by assignee user IDs.
func (s *OptionService) WithAssigneeIDs(ids []int) RequestOption {
	return idSliceOption(ParamAssigneeIDs, "assigneeId", ids)
}

// WithCreatedUserIDs returns an option to filter by created user IDs.
func (s *OptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return idSliceOption(ParamCreatedUserIDs, "createdUserId", ids)
}

// WithResolutionIDs returns an option to filter by resolution IDs.
func (s *OptionService) WithResolutionIDs(ids []int) RequestOption {
	return idSliceOption(ParamResolutionIDs, "resolutionId", ids)
}

// WithIDs returns an option to filter by issue IDs.
func (s *OptionService) WithIDs(ids []int) RequestOption {
	return idSliceOption(ParamIDs, "id", ids)
}

// WithParentIssueIDs returns an option to filter by parent issue IDs.
func (s *OptionService) WithParentIssueIDs(ids []int) RequestOption {
	return idSliceOption(ParamParentIssueIDs, "parentIssueId", ids)
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
	return boolOption(ParamAttachment, enabled)
}

// WithSharedFile returns an option to include only issues with shared files.
func (s *OptionService) WithSharedFile(enabled bool) RequestOption {
	return boolOption(ParamSharedFile, enabled)
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
	return boolOption(ParamHasDueDate, enabled)
}
