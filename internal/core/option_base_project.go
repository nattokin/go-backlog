package core

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/model"
)

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
