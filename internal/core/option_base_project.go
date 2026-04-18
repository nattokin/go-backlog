package core

import (
	"fmt"
	"net/url"

	"github.com/nattokin/go-backlog/internal/model"
)

// WithAll returns an option to set the `all` parameter.
func (s *OptionService) WithAll(enabled bool) RequestOption {
	return boolOption(ParamAll, enabled)
}

// WithArchived returns an option to set the `archived` parameter.
func (s *OptionService) WithArchived(enabled bool) RequestOption {
	// apiArchived and queryArchived share the same string value "archived",
	// so we use apiArchived as the canonical type here and accept both in services.
	return boolOption(ParamArchived, enabled)
}

// WithChartEnabled returns a option that sets the `chartEnabled` field.
func (s *OptionService) WithChartEnabled(enabled bool) RequestOption {
	return boolOption(ParamChartEnabled, enabled)
}

// WithKey returns a option that sets the `key` field.
func (s *OptionService) WithKey(key string) RequestOption {
	return nonEmptyStringOption(ParamKey, key)
}

// WithProjectLeaderCanEditProjectLeader returns a option.
func (s *OptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return boolOption(ParamProjectLeaderCanEditProjectLeader, enabled)
}

// WithSubtaskingEnabled returns a option that sets the `subtaskingEnabled` field.
func (s *OptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return boolOption(ParamSubtaskingEnabled, enabled)
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
