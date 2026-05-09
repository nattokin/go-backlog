package core

import (
	"fmt"
	"net/url"

	"github.com/nattokin/go-backlog/internal/model"
)

func (s *OptionService) WithBase(base string) RequestOption {
	return nonEmptyStringOption(ParamBase, base)
}

func (s *OptionService) WithBranch(branch string) RequestOption {
	return nonEmptyStringOption(ParamBranch, branch)
}

func (s *OptionService) WithColor(color string) RequestOption {
	return nonEmptyStringOption(ParamColor, color)
}

func (s *OptionService) WithComment(comment string) RequestOption {
	return &APIParamOption{
		Type:    ParamComment,
		SetFunc: setStringFunc(ParamComment, comment),
	}
}

func (s *OptionService) WithContent(content string) RequestOption {
	return nonEmptyStringOption(ParamContent, content)
}

func (s *OptionService) WithDescription(description string) RequestOption {
	return &APIParamOption{
		Type:    ParamDescription,
		SetFunc: setStringFunc(ParamDescription, description),
	}
}

func (s *OptionService) WithHookURL(hookURL string) RequestOption {
	return nonEmptyStringOption(ParamHookURL, hookURL)
}

// WithInitialDate sets `initialDate` for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDate(date string) RequestOption {
	return dateFormatStringOption(ParamInitialDate, date)
}

// WithInitialDateMax sets `max` for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDateMax(date string) RequestOption {
	return dateFormatStringOption(ParamMax, date)
}

// WithInitialDateMin sets `min` for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDateMin(date string) RequestOption {
	return dateFormatStringOption(ParamMin, date)
}

func (s *OptionService) WithKey(key string) RequestOption {
	return nonEmptyStringOption(ParamKey, key)
}

func (s *OptionService) WithKeyword(keyword string) RequestOption {
	return &APIParamOption{
		Type:    ParamKeyword,
		SetFunc: setStringFunc(ParamKeyword, keyword),
	}
}

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

func (s *OptionService) WithMailAddress(mailAddress string) RequestOption {
	// ToDo: validate mailAddress (Note: The validation remains as simple not-empty check)
	return nonEmptyStringOption(ParamMailAddress, mailAddress)
}

func (s *OptionService) WithName(name string) RequestOption {
	return nonEmptyStringOption(ParamName, name)
}

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

func (s *OptionService) WithSummary(summary string) RequestOption {
	return nonEmptyStringOption(ParamSummary, summary)
}

func (s *OptionService) WithTemplateDescription(description string) RequestOption {
	return &APIParamOption{
		Type:    ParamTemplateDescription,
		SetFunc: setStringFunc(ParamTemplateDescription, description),
	}
}

func (s *OptionService) WithTemplateSummary(summary string) RequestOption {
	return &APIParamOption{
		Type:    ParamTemplateSummary,
		SetFunc: setStringFunc(ParamTemplateSummary, summary),
	}
}

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

func (s *OptionService) WithUnit(unit string) RequestOption {
	return &APIParamOption{
		Type:    ParamUnit,
		SetFunc: setStringFunc(ParamUnit, unit),
	}
}

// nonEmptyStringOption builds a RequestOption that rejects empty strings.
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

func setStringFunc(key APIParamOptionType, value string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), value)
		return nil
	}
}
