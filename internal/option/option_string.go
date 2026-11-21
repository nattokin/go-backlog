package option

import (
	"net/url"

	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// WithBase sets `base`.
// The value must not be empty or whitespace-only.
func (s *OptionService) WithBase(base string) *APIParamOption {
	return nonEmptyStringOption(ParamBase, base)
}

// WithBranch sets `branch`.
// The value must not be empty or whitespace-only.
func (s *OptionService) WithBranch(branch string) *APIParamOption {
	return nonEmptyStringOption(ParamBranch, branch)
}

// WithColor sets `color`.
// The value must not be empty or whitespace-only.
func (s *OptionService) WithColor(color string) *APIParamOption {
	return nonEmptyStringOption(ParamColor, color)
}

// WithComment sets `comment`.
func (s *OptionService) WithComment(comment string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamComment,
		SetFunc: setStringFunc(ParamComment, comment),
	}
}

// WithContent sets `content`.
// The value must not be empty or whitespace-only.
func (s *OptionService) WithContent(content string) *APIParamOption {
	return nonEmptyStringOption(ParamContent, content)
}

// WithDescription sets `description`.
func (s *OptionService) WithDescription(description string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamDescription,
		SetFunc: setStringFunc(ParamDescription, description),
	}
}

// WithHookURL sets `hookUrl`.
// The value must not be empty or whitespace-only.
func (s *OptionService) WithHookURL(hookURL string) *APIParamOption {
	return nonEmptyStringOption(ParamHookURL, hookURL)
}

// WithKey sets `key`.
// The value must not be empty or whitespace-only.
func (s *OptionService) WithKey(key string) *APIParamOption {
	return nonEmptyStringOption(ParamKey, key)
}

// WithKeyword sets `keyword`.
func (s *OptionService) WithKeyword(keyword string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamKeyword,
		SetFunc: setStringFunc(ParamKeyword, keyword),
	}
}

// WithIssueSort sets `sort`.
// The value must be a sort key accepted by the issue list endpoints.
func (s *OptionService) WithIssueSort(sort string) *APIParamOption {
	return &APIParamOption{
		Type: ParamSort,
		CheckFunc: func() *validation.Error {
			return validate.ValidateIssueSort(ParamSort.Value(), sort)
		},
		SetFunc: setStringFunc(ParamSort, sort),
	}
}

// WithMailAddress sets `mailAddress`.
// The value must be a single, well-formed email address.
func (s *OptionService) WithMailAddress(mailAddress string) *APIParamOption {
	return &APIParamOption{
		Type: ParamMailAddress,
		CheckFunc: func() *validation.Error {
			return validate.ValidateEmail(ParamMailAddress.Value(), mailAddress)
		},
		SetFunc: setStringFunc(ParamMailAddress, mailAddress),
	}
}

// WithName sets `name`.
// The value must not be empty or whitespace-only.
func (s *OptionService) WithName(name string) *APIParamOption {
	return nonEmptyStringOption(ParamName, name)
}

// WithOrder sets `order`.
// The value must be "asc" or "desc".
func (s *OptionService) WithOrder(order string) *APIParamOption {
	return &APIParamOption{
		Type: ParamOrder,
		CheckFunc: func() *validation.Error {
			return validate.ValidateOrder(ParamOrder.Value(), order)
		},
		SetFunc: setStringFunc(ParamOrder, order),
	}
}

// WithPassword sets `password`.
// The value must be at least 8 characters long.
func (s *OptionService) WithPassword(password string) *APIParamOption {
	return &APIParamOption{
		Type: ParamPassword,
		CheckFunc: func() *validation.Error {
			return validate.ValidatePassword(ParamPassword.Value(), password)
		},
		SetFunc: setStringFunc(ParamPassword, password),
	}
}

// WithSummary sets `summary`.
// The value must not be empty or whitespace-only.
func (s *OptionService) WithSummary(summary string) *APIParamOption {
	return nonEmptyStringOption(ParamSummary, summary)
}

// WithTemplateDescription sets `templateDescription`.
func (s *OptionService) WithTemplateDescription(description string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamTemplateDescription,
		SetFunc: setStringFunc(ParamTemplateDescription, description),
	}
}

// WithTemplateSummary sets `templateSummary`.
func (s *OptionService) WithTemplateSummary(summary string) *APIParamOption {
	return &APIParamOption{
		Type:    ParamTemplateSummary,
		SetFunc: setStringFunc(ParamTemplateSummary, summary),
	}
}

// WithTextFormattingRule sets `textFormattingRule`.
// The value must be "backlog" or "markdown".
func (s *OptionService) WithTextFormattingRule(format string) *APIParamOption {
	return &APIParamOption{
		Type: ParamTextFormattingRule,
		CheckFunc: func() *validation.Error {
			return validate.ValidateTextFormattingRule(ParamTextFormattingRule.Value(), format)
		},
		SetFunc: setStringFunc(ParamTextFormattingRule, format),
	}
}

// WithUnit sets `unit` for Number type custom fields.
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
