package core

// WithApplicableIssueTypeIDs returns an option to set the `applicableIssueTypes[]` parameter.
func (s *OptionService) WithApplicableIssueTypeIDs(ids []int) RequestOption {
	return intSliceOption(ParamApplicableIssueTypeIDs, "applicableIssueTypes", ids)
}
