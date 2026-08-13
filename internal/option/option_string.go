package option

import (
	"net/url"

	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
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

func (s *OptionService) WithIssueSort(sort string) *APIParamOption {
	return &APIParamOption{
		Type: ParamSort,
		CheckFunc: func() *validation.Error {
			return validate.ValidateIssueSort(ParamSort.Value(), sort)
		},
		SetFunc: setStringFunc(ParamSort, sort),
	}
}

func (s *OptionService) WithMailAddress(mailAddress string) *APIParamOption {
	return &APIParamOption{
		Type: ParamMailAddress,
		CheckFunc: func() *validation.Error {
			return validate.ValidateEmail(ParamMailAddress.Value(), mailAddress)
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
		CheckFunc: func() *validation.Error {
			return validate.ValidateOrder(ParamOrder.Value(), order)
		},
		SetFunc: setStringFunc(ParamOrder, order),
	}
}

func (s *OptionService) WithPassword(password string) *APIParamOption {
	return &APIParamOption{
		Type: ParamPassword,
		CheckFunc: func() *validation.Error {
			return validate.ValidatePassword(ParamPassword.Value(), password)
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

func (s *OptionService) WithTextFormattingRule(format string) *APIParamOption {
	return &APIParamOption{
		Type: ParamTextFormattingRule,
		CheckFunc: func() *validation.Error {
			return validate.ValidateTextFormattingRule(ParamTextFormattingRule.Value(), format)
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

// nonEmptyStringOption builds a *APIParamOption that rejects empty or whitespace-only strings.
func nonEmptyStringOption(paramType APIParamOptionType, value string) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *validation.Error {
			return validate.ValidateNonEmptyString(paramType.Value(), value)
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
