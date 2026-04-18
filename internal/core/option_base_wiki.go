package core

// WithContent returns a option that sets the `content` field.
func (s *OptionService) WithContent(content string) RequestOption {
	return nonEmptyStringOption(ParamContent, content)
}

// WithMailNotify returns a option that sets the `mailNotify` field.
func (s *OptionService) WithMailNotify(enabled bool) RequestOption {
	return boolOption(ParamMailNotify, enabled)
}
