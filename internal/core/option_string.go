package core

import (
	"fmt"
	"net/mail"
	"net/url"
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

func (s *OptionService) WithKey(key string) RequestOption {
	return nonEmptyStringOption(ParamKey, key)
}

func (s *OptionService) WithKeyword(keyword string) RequestOption {
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

func (s *OptionService) WithIssueSort(sort string) RequestOption {
	return &APIParamOption{
		Type: ParamSort,
		CheckFunc: func() error {
			for _, v := range validIssueSorts {
				if sort == v {
					return nil
				}
			}
			return NewValidationError(fmt.Sprintf("invalid sort value: %q", sort))
		},
		SetFunc: setStringFunc(ParamSort, sort),
	}
}

func (s *OptionService) WithMailAddress(mailAddress string) RequestOption {
	return &APIParamOption{
		Type: ParamMailAddress,
		CheckFunc: func() error {
			addr, err := mail.ParseAddress(mailAddress)
			if err != nil || addr.Address != mailAddress {
				return NewValidationError(fmt.Sprintf("mailAddress %q is not a valid email address", mailAddress))
			}
			return nil
		},
		SetFunc: setStringFunc(ParamMailAddress, mailAddress),
	}
}

func (s *OptionService) WithName(name string) RequestOption {
	return nonEmptyStringOption(ParamName, name)
}

func (s *OptionService) WithOrder(order string) RequestOption {
	return &APIParamOption{
		Type: ParamOrder,
		CheckFunc: func() error {
			if order != "asc" && order != "desc" {
				return NewValidationError("order must be only 'asc' or 'desc'")
			}
			return nil
		},
		SetFunc: setStringFunc(ParamOrder, order),
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

var validFormats = []string{"backlog", "markdown"}

func (s *OptionService) WithTextFormattingRule(format string) RequestOption {
	return &APIParamOption{
		Type: ParamTextFormattingRule,
		CheckFunc: func() error {
			for _, v := range validFormats {
				if format == v {
					return nil
				}
			}
			return NewValidationError("format must be only 'backlog' or 'markdown'")
		},
		SetFunc: setStringFunc(ParamTextFormattingRule, format),
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
