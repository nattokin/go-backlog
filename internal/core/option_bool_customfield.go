package core

// WithRequired returns an option to set the `required` parameter.
func (s *OptionService) WithRequired(required bool) RequestOption {
	return boolOption(ParamRequired, required)
}
