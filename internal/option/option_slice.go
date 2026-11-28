package option

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// WithActivityTypeIDs sets `activityTypeId[]` for activity list endpoints.
func (s *OptionService) WithActivityTypeIDs(typeIDs []int) *APIParamOption {
	return positiveIntSliceOption(ParamActivityTypeIDs, "activityTypeIds", typeIDs)
}

// WithWebhookActivityTypeIDs sets `activityTypeIds[]` for the Add/Update Webhook endpoints.
func (s *OptionService) WithWebhookActivityTypeIDs(typeIDs []int) *APIParamOption {
	return positiveIntSliceOption(ParamWebhookActivityTypeIDs, "activityTypeIds", typeIDs)
}

// WithApplicableIssueTypeIDs sets `applicableIssueTypes[]`.
// Each ID must be at least 1.
func (s *OptionService) WithApplicableIssueTypeIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamApplicableIssueTypeIDs, "applicableIssueTypes", ids)
}

// WithAttachmentIDs sets `attachmentId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithAttachmentIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamAttachmentIDs, "attachmentId", ids)
}

// WithIssueExpand sets `expand[]` for the issue list endpoints.
func (s *OptionService) WithIssueExpand(expands []string) *APIParamOption {
	return &APIParamOption{
		Type: ParamExpand,
		CheckFunc: func() *validation.Error {
			for _, expand := range expands {
				if ve := validate.ValidateIssueExpand(ParamExpand.Value(), expand); ve != nil {
					return ve
				}
			}
			return nil
		},
		SetFunc: addStringFunc(ParamExpand, expands),
	}
}

// WithItems sets `items[]` for List type custom fields.
func (s *OptionService) WithItems(items []string) *APIParamOption {
	return &APIParamOption{
		Type: ParamItems,
		CheckFunc: func() *validation.Error {
			for i, item := range items {
				if item == "" {
					return validation.NewError(ParamItems.Value(), fmt.Sprintf("items[%d] must not be empty", i))
				}
			}
			return nil
		},
		SetFunc: addStringFunc(ParamItems, items),
	}
}

// WithProjectIDs sets `projectId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithProjectIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamProjectIDs, "projectId", ids)
}

// WithIssueTypeIDs sets `issueTypeId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithIssueTypeIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamIssueTypeIDs, "issueTypeId", ids)
}

// WithCategoryIDs sets `categoryId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithCategoryIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamCategoryIDs, "categoryId", ids)
}

// WithVersionIDs sets `versionId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithVersionIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamVersionIDs, "versionId", ids)
}

// WithMilestoneIDs sets `milestoneId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithMilestoneIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamMilestoneIDs, "milestoneId", ids)
}

// WithIssueIDs sets `issueId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithIssueIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamIssueIDs, "issueId", ids)
}

// WithNotifiedUserIDs sets `notifiedUserId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithNotifiedUserIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamNotifiedUserIDs, "notifiedUserId", ids)
}

// WithStatusIDs sets `statusId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithStatusIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamStatusIDs, "statusId", ids)
}

// WithPriorityIDs sets `priorityId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithPriorityIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamPriorityIDs, "priorityId", ids)
}

// WithAssigneeIDs sets `assigneeId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithAssigneeIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamAssigneeIDs, "assigneeId", ids)
}

// WithCreatedUserIDs sets `createdUserId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithCreatedUserIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamCreatedUserIDs, "createdUserId", ids)
}

// WithResolutionIDs sets `resolutionId[]`.
// Each ID must be at least 0.
func (s *OptionService) WithResolutionIDs(ids []int) *APIParamOption {
	return nonNegativeIntSliceOption(ParamResolutionIDs, "resolutionId", ids)
}

// WithIDs sets `id[]`.
// Each ID must be at least 1.
func (s *OptionService) WithIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamIDs, "id", ids)
}

// WithParentIssueIDs sets `parentIssueId[]`.
// Each ID must be at least 1.
func (s *OptionService) WithParentIssueIDs(ids []int) *APIParamOption {
	return positiveIntSliceOption(ParamParentIssueIDs, "parentIssueId", ids)
}

func nonNegativeIntSliceOption(paramType APIParamOptionType, paramName string, values []int) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *validation.Error {
			return validate.ValidateNonNegativeInts(paramName, values)
		},
		SetFunc: addIntFunc(paramType, values),
	}
}

func positiveIntSliceOption(paramType APIParamOptionType, paramName string, values []int) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *validation.Error {
			return validate.ValidatePositiveInts(paramName, values)
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
