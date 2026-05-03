package core

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/nattokin/go-backlog/internal/model"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// WithBase returns an option that sets the `base` field (merge base branch name).
func (s *OptionService) WithBase(base string) RequestOption {
	return nonEmptyStringOption(ParamBase, base)
}

// WithBranch returns an option that sets the `branch` field (merging branch name).
func (s *OptionService) WithBranch(branch string) RequestOption {
	return nonEmptyStringOption(ParamBranch, branch)
}

// WithColor returns an option that sets the `color` field.
func (s *OptionService) WithColor(color string) RequestOption {
	return nonEmptyStringOption(ParamColor, color)
}

// WithComment returns an option to set the `comment` parameter.
func (s *OptionService) WithComment(comment string) RequestOption {
	return &APIParamOption{
		Type:    ParamComment,
		SetFunc: setStringFunc(ParamComment, comment),
	}
}

// WithContent returns a option that sets the `content` field.
func (s *OptionService) WithContent(content string) RequestOption {
	return nonEmptyStringOption(ParamContent, content)
}

// WithDescription returns an option to set the `description` parameter.
func (s *OptionService) WithDescription(description string) RequestOption {
	return &APIParamOption{
		Type:    ParamDescription,
		SetFunc: setStringFunc(ParamDescription, description),
	}
}

// WithHookURL returns an option that sets the `hookUrl` parameter.
func (s *OptionService) WithHookURL(hookURL string) RequestOption {
	return nonEmptyStringOption(ParamHookURL, hookURL)
}

// WithInitialDate returns an option to set the `initialDate` parameter for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDate(date string) RequestOption {
	return dateFormatStringOption(ParamInitialDate, date)
}

// WithInitialDateMax returns an option to set the `max` parameter for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDateMax(date string) RequestOption {
	return dateFormatStringOption(ParamMax, date)
}

// WithInitialDateMin returns an option to set the `min` parameter for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDateMin(date string) RequestOption {
	return dateFormatStringOption(ParamMin, date)
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

// WithSummary returns an option to set the `summary` parameter.
func (s *OptionService) WithSummary(summary string) RequestOption {
	return nonEmptyStringOption(ParamSummary, summary)
}

// WithTemplateDescription returns an option to set the `templateDescription` parameter.
func (s *OptionService) WithTemplateDescription(description string) RequestOption {
	return &APIParamOption{
		Type:    ParamTemplateDescription,
		SetFunc: setStringFunc(ParamTemplateDescription, description),
	}
}

// WithTemplateSummary returns an option to set the `templateSummary` parameter.
func (s *OptionService) WithTemplateSummary(summary string) RequestOption {
	return &APIParamOption{
		Type:    ParamTemplateSummary,
		SetFunc: setStringFunc(ParamTemplateSummary, summary),
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

// WithUnit returns an option to set the `unit` parameter for Number type custom fields.
func (s *OptionService) WithUnit(unit string) RequestOption {
	return &APIParamOption{
		Type:    ParamUnit,
		SetFunc: setStringFunc(ParamUnit, unit),
	}
}

//
// ──────────────────────────────────────────────────────────────
//  Option builder helpers
// ──────────────────────────────────────────────────────────────
//

// dateFormatStringOption builds a RequestOption that validates the string matches
// "yyyy-MM-dd" format and sets it.
func dateFormatStringOption(paramType APIParamOptionType, date string) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if !datePattern.MatchString(date) {
				return NewValidationError(fmt.Sprintf("%s must be formatted as yyyy-MM-dd, got %q", paramType.Value(), date))
			}
			return nil
		},
		SetFunc: setStringFunc(paramType, date),
	}
}

// nonEmptyStringOption builds a RequestOption that validates the string is not empty and sets it.
func nonEmptyStringOption(paramType APIParamOptionType, value string) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if value == "" {
				return NewValidationError(fmt.Sprintf("%s must not be empty", paramType.Value()))
			}
			return nil
		},
		SetFunc: setStringFunc(paramType, value),
	}
}

//
// ──────────────────────────────────────────────────────────────
//  SetFunc factories
// ──────────────────────────────────────────────────────────────
//

// setStringFunc returns a SetFunc that calls v.Set with the given string value.
func setStringFunc(key APIParamOptionType, value string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), value)
		return nil
	}
}
