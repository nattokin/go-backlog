package core

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/nattokin/go-backlog/internal/model"
)

//
// ──────────────────────────────────────────────────────────────
//  OptionService — unified builder
// ──────────────────────────────────────────────────────────────
//

// OptionService provides builders for request options.
// Each XxxOptionService selectively exposes only the valid methods.
type OptionService struct{}

// --- Boolean options ------------------------------------------------------------

// WithAll returns an option to set the `all` parameter.
func (s *OptionService) WithAll(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamAll,
		SetFunc: func(v url.Values) error {
			v.Set(ParamAll.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithArchived returns an option to set the `archived` parameter.
func (s *OptionService) WithArchived(enabled bool) RequestOption {
	// apiArchived and queryArchived share the same string value "archived",
	// so we use apiArchived as the canonical type here and accept both in services.
	return &APIParamOption{
		Type: ParamArchived,
		SetFunc: func(v url.Values) error {
			v.Set(ParamArchived.Value(), strconv.FormatBool(enabled))
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

// WithChartEnabled returns a option that sets the `chartEnabled` field.
func (s *OptionService) WithChartEnabled(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamChartEnabled,
		SetFunc: func(v url.Values) error {
			v.Set(ParamChartEnabled.Value(), strconv.FormatBool(enabled))
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

// WithMailNotify returns a option that sets the `mailNotify` field.
func (s *OptionService) WithMailNotify(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamMailNotify,
		SetFunc: func(v url.Values) error {
			v.Set(ParamMailNotify.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithProjectLeaderCanEditProjectLeader returns a option.
func (s *OptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamProjectLeaderCanEditProjectLeader,
		SetFunc: func(v url.Values) error {
			v.Set(ParamProjectLeaderCanEditProjectLeader.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// WithSendMail returns a option to specify whether to send an invitation email.
func (s *OptionService) WithSendMail(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamSendMail,
		SetFunc: func(v url.Values) error {
			v.Set(ParamSendMail.Value(), strconv.FormatBool(enabled))
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

// WithSubtaskingEnabled returns a option that sets the `subtaskingEnabled` field.
func (s *OptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return &APIParamOption{
		Type: ParamSubtaskingEnabled,
		SetFunc: func(v url.Values) error {
			v.Set(ParamSubtaskingEnabled.Value(), strconv.FormatBool(enabled))
			return nil
		},
	}
}

// --- Integer options ------------------------------------------------------------

// WithCount returns an option to set the `count` parameter.
func (s *OptionService) WithCount(count int) RequestOption {
	return &APIParamOption{
		Type: ParamCount,
		CheckFunc: func() error {
			if count < 1 || 100 < count {
				return NewValidationError("count must be between 1 and 100")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamCount.Value(), strconv.Itoa(count))
			return nil
		},
	}
}

// WithMaxID returns an option to set the `maxId` parameter.
func (s *OptionService) WithMaxID(id int) RequestOption {
	return &APIParamOption{
		Type: ParamMaxID,
		CheckFunc: func() error {
			return validateActivityID(id, "maxID")
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamMaxID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// WithMinID returns an option to set the `minId` parameter.
func (s *OptionService) WithMinID(id int) RequestOption {
	return &APIParamOption{
		Type: ParamMinID,
		CheckFunc: func() error {
			return validateActivityID(id, "minID")
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamMinID.Value(), strconv.Itoa(id))
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

// WithUserID returns a option to set the user's ID.
func (s *OptionService) WithUserID(id int) RequestOption {
	return &APIParamOption{
		Type: ParamUserID,
		CheckFunc: func() error {
			if id < 1 {
				return NewValidationError(fmt.Sprintf("invalid %s: must not be less than 1", ParamUserID))
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamUserID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// --- Integer slice options -------------------------------------------------------

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

// --- String options ------------------------------------------------------------

// WithContent returns a option that sets the `content` field.
func (s *OptionService) WithContent(content string) RequestOption {
	return &APIParamOption{
		Type: ParamContent,
		CheckFunc: func() error {
			if content == "" {
				return NewValidationError("content must not be empty")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamContent.Value(), content)
			return nil
		},
	}
}

// WithKey returns a option that sets the `key` field.
func (s *OptionService) WithKey(key string) RequestOption {
	return &APIParamOption{
		Type: ParamKey,
		CheckFunc: func() error {
			if key == "" {
				return NewValidationError("key must not be empty")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamKey.Value(), key)
			return nil
		},
	}
}

// WithKeyword returns an option to set the `keyword` parameter.
func (s *OptionService) WithKeyword(keyword string) RequestOption {
	return &APIParamOption{
		Type: ParamKeyword,
		SetFunc: func(v url.Values) error {
			v.Set(ParamKeyword.Value(), keyword)
			return nil
		},
	}
}

// WithMailAddress returns a option that sets the `mailAddress` field.
func (s *OptionService) WithMailAddress(mailAddress string) RequestOption {
	// ToDo: validate mailAddress (Note: The validation remains as simple not-empty check)
	return &APIParamOption{
		Type: ParamMailAddress,
		CheckFunc: func() error {
			if mailAddress == "" {
				return NewValidationError("mailAddress must not be empty")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamMailAddress.Value(), mailAddress)
			return nil
		},
	}
}

// WithName returns a option that sets the `name` field.
func (s *OptionService) WithName(name string) RequestOption {
	return &APIParamOption{
		Type: ParamName,
		CheckFunc: func() error {
			if name == "" {
				return NewValidationError("name must not be empty")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamName.Value(), name)
			return nil
		},
	}
}

// WithPassword returns a option that sets the `password` field.
func (s *OptionService) WithPassword(password string) RequestOption {
	return &APIParamOption{
		Type: ParamPassword,
		CheckFunc: func() error {
			if len(password) < 8 {
				return NewValidationError("password must be at least 8 characters long")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamPassword.Value(), password)
			return nil
		},
	}
}

// --- Date options (time.Time) ---------------------------------------------------

const issueDateFormat = "2006-01-02"

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

// --- Enum or special options ----------------------------------------------------

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
		SetFunc: func(v url.Values) error {
			for _, id := range typeIDs {
				v.Add(ParamActivityTypeIDs.Value(), strconv.Itoa(id))
			}
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

// WithOrder returns an option to set the `order` parameter.
func (s *OptionService) WithOrder(order model.Order) RequestOption {
	return &APIParamOption{
		Type: ParamOrder,
		CheckFunc: func() error {
			if order != model.OrderAsc && order != model.OrderDesc {
				msg := fmt.Sprintf("order must be only '%s' or '%s'", string(model.OrderAsc), string(model.OrderDesc))
				return NewValidationError(msg)
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamOrder.Value(), string(order))
			return nil
		},
	}
}

// WithRoleType returns a option that sets the `roleType` field.
func (s *OptionService) WithRoleType(roleType model.Role) RequestOption {
	return &APIParamOption{
		Type: ParamRoleType,
		CheckFunc: func() error {
			if roleType < 1 || 6 < roleType {
				return NewValidationError("roleType must be between 1 and 6")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamRoleType.Value(), strconv.Itoa(int(roleType)))
			return nil
		},
	}
}

// WithTextFormattingRule returns a option that sets the `textFormattingRule` field.
func (s *OptionService) WithTextFormattingRule(format model.Format) RequestOption {
	return &APIParamOption{
		Type: ParamTextFormattingRule,
		CheckFunc: func() error {
			if format != model.FormatBacklog && format != model.FormatMarkdown {
				msg := fmt.Sprintf("format must be only '%s' or '%s'", string(model.FormatBacklog), string(model.FormatMarkdown))
				return NewValidationError(msg)
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamTextFormattingRule.Value(), string(format))
			return nil
		},
	}
}

//
// ──────────────────────────────────────────────────────────────
//  Internal Helpers
// ──────────────────────────────────────────────────────────────
//

// validateActivityID ensures that the given activity ID is within the valid range [1, 26].
func validateActivityID(id int, key string) error {
	if id < 1 || id > MaxActivityTypeID {
		return NewValidationError(fmt.Sprintf("invalid %s: must be between 1 and %d", key, MaxActivityTypeID))
	}
	return nil
}
