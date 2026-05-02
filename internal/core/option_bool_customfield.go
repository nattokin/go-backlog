package core

// WithRequired returns an option to set the `required` parameter.
func (s *OptionService) WithRequired(required bool) RequestOption {
	return boolOption(ParamRequired, required)
}

// WithAllowInput returns an option to set the `allowInput` parameter for List type custom fields.
// When true, users can enter a free-text value in addition to selecting from the list.
func (s *OptionService) WithAllowInput(allowInput bool) RequestOption {
	return boolOption(ParamAllowInput, allowInput)
}

// WithAllowAddItem returns an option to set the `allowAddItem` parameter for List type custom fields.
// When true, users can add new items to the list from the issue form.
func (s *OptionService) WithAllowAddItem(allowAddItem bool) RequestOption {
	return boolOption(ParamAllowAddItem, allowAddItem)
}
