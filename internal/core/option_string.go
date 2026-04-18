package core

import (
	"fmt"

	"github.com/nattokin/go-backlog/internal/model"
)

// WithContent returns a option that sets the `content` field.
func (s *OptionService) WithContent(content string) RequestOption {
	return nonEmptyStringOption(ParamContent, content)
}

// WithKey returns a option that sets the `key` field.
func (s *OptionService) WithKey(key string) RequestOption {
	return nonEmptyStringOption(ParamKey, key)
}

// WithKeyword returns an option to set the `keyword` parameter.
func (s *OptionService) WithKeyword(keyword string) RequestOption {
	return &APIParamOption{
		Type:    ParamKeyword,
		SetFunc: setStringFunc(ParamKeyword, keyword),
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
		SetFunc: setStringFunc(ParamSort, string(sort)),
	}
}

// WithMailAddress returns a option that sets the `mailAddress` field.
func (s *OptionService) WithMailAddress(mailAddress string) RequestOption {
	// ToDo: validate mailAddress (Note: The validation remains as simple not-empty check)
	return nonEmptyStringOption(ParamMailAddress, mailAddress)
}

// WithName returns a option that sets the `name` field.
func (s *OptionService) WithName(name string) RequestOption {
	return nonEmptyStringOption(ParamName, name)
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
		SetFunc: setStringFunc(ParamOrder, string(order)),
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
		SetFunc: setStringFunc(ParamPassword, password),
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
		SetFunc: setStringFunc(ParamTextFormattingRule, string(format)),
	}
}
