package core

import "net/url"

// WithContent returns a option that sets the `content` field.
func (s *OptionService) WithContent(content string) RequestOption {
	return &APIParamOption{
		Type: ParamContent,
		CheckFunc: func() error {
			if content == "" {
				return NewValidationError("content must not be empty")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamContent.Value(), content)
			return nil
		},
	}
}

// WithMailNotify returns a option that sets the `mailNotify` field.
func (s *OptionService) WithMailNotify(enabled bool) RequestOption {
	return boolOption(ParamMailNotify, enabled)
}
