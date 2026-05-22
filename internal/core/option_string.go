package core

import (
	"fmt"
	"net/mail"
	"net/url"
)

func (s *OptionService) WithBase(base string) *APIParamOption {
	return nonEmptyStringOption(ParamBase, base)
}

func (s *OptionService) WithBranch(branch string) *APIParamOption {
	return nonEmptyStringOption(ParamBranch, branch)
}

func (s *OptionService) WithColor(color string) *APIParamOption {
	return nonEmptyStringOption(ParamColor, color)
}

func (s *OptionService) WithComment(comment string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamComment,
		SetFunc: setStringFunc(ParamComment, comment),
	}
}

func (s *OptionService) WithContent(content string) *APIParamOption {
	return nonEmptyStringOption(ParamContent, content)
}

func (s *OptionService) WithDescription(description string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamDescription,
		SetFunc: setStringFunc(ParamDescription, description),
	}
}

func (s *OptionService) WithHookURL(hookURL string) *APIParamOption {
	return nonEmptyStringOption(ParamHookURL, hookURL)
}

func (s *OptionService) WithKey(key string) *APIParamOption {
	return nonEmptyStringOption(ParamKey, key)
}

func (s *OptionService) WithKeyword(keyword string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamKeyword,
		SetFunc: setStringFunc(ParamKeyword, keyword),
	}
}

var validIssueSorts = []string{
	"issueType", "category", "version", "milestone", "summary", "status",
	"priority", "attachment", "sharedFile", "created", "createdUser",
	"updated", "updatedUser", "assignee", "startDate", "dueDate",
	"estimatedHours", "actualHours", "childIssue",
}

func (s *OptionService) WithIssueSort(sort string) *APIParamOption {
	return &APIParamOption{
		Type: ParamSort,
		CheckFunc: func() *ValidationError {
			for _, v := range validIssueSorts {
				if sort == v {
					return nil
				}
			}
			return NewValidationError(ParamSort.Value(), fmt.Sprintf("invalid sort value: %q", sort))
		},
		SetFunc: setStringFunc(ParamSort, sort),
	}
}

func (s *OptionService) WithMailAddress(mailAddress string) *APIParamOption {
	return &APIParamOption{
		Type: ParamMailAddress,
		CheckFunc: func() *ValidationError {
			addr, err := mail.ParseAddress(mailAddress)
			if err != nil || addr.Address != mailAddress {
				return NewValidationError(ParamMailAddress.Value(), fmt.Sprintf("mailAddress %q is not a valid email address", mailAddress))
			}
			return nil
		},
		SetFunc: setStringFunc(ParamMailAddress, mailAddress),
	}
}

func (s *OptionService) WithName(name string) *APIParamOption {
	return nonEmptyStringOption(ParamName, name)
}

func (s *OptionService) WithOrder(order string) *APIParamOption {
	return &APIParamOption{
		Type: ParamOrder,
		CheckFunc: func() *ValidationError {
			if order != "asc" && order != "desc" {
				return NewValidationError(ParamOrder.Value(), "order must be only 'asc' or 'desc'")
			}
			return nil
		},
		SetFunc: setStringFunc(ParamOrder, order),
	}
}

func (s *OptionService) WithPassword(password string) *APIParamOption {
	return &APIParamOption{
		Type: ParamPassword,
		CheckFunc: func() *ValidationError {
			if len(password) < 8 {
				return NewValidationError(ParamPassword.Value(), "password must be at least 8 characters long")
			}
			return nil
		},
		SetFunc: setStringFunc(ParamPassword, password),
	}
}

func (s *OptionService) WithSummary(summary string) *APIParamOption {
	return nonEmptyStringOption(ParamSummary, summary)
}

func (s *OptionService) WithTemplateDescription(description string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamTemplateDescription,
		SetFunc: setStringFunc(ParamTemplateDescription, description),
	}
}

func (s *OptionService) WithTemplateSummary(summary string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamTemplateSummary,
		SetFunc: setStringFunc(ParamTemplateSummary, summary),
	}
}

var validFormats = []string{"backlog", "markdown"}

func (s *OptionService) WithTextFormattingRule(format string) *APIParamOption {
	return &APIParamOption{
		Type: ParamTextFormattingRule,
		CheckFunc: func() *ValidationError {
			for _, v := range validFormats {
				if format == v {
					return nil
				}
			}
			return NewValidationError(ParamTextFormattingRule.Value(), "format must be only 'backlog' or 'markdown'")
		},
		SetFunc: setStringFunc(ParamTextFormattingRule, format),
	}
}

func (s *OptionService) WithUnit(unit string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamUnit,
		SetFunc: setStringFunc(ParamUnit, unit),
	}
}

// nonEmptyStringOption builds a *APIParamOption that rejects empty strings.
func nonEmptyStringOption(paramType APIParamOptionType, value string) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *ValidationError {
			if value == "" {
				return NewValidationError(paramType.Value(), fmt.Sprintf("%s must not be empty", paramType.Value()))
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
