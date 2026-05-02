package core

// WithApplicableIssueTypeIDs returns an option to set the `applicableIssueTypes[]` parameter.
func (s *OptionService) WithApplicableIssueTypeIDs(ids []int) RequestOption {
	return intSliceOption(ParamApplicableIssueTypeIDs, "applicableIssueTypes", ids)
}

// WithItems returns an option to set the `items[]` parameter for List type custom fields.
// Each string becomes a selectable list item.
func (s *OptionService) WithItems(items []string) RequestOption {
	return &APIParamOption{
		Type:    ParamItems,
		SetFunc: addStringFunc(ParamItems, items),
	}
}
