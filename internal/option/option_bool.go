package option

import (
	"net/url"
	"strconv"
)

// WithAll sets `all`.
func (s *OptionService) WithAll(enabled bool) *APIParamOption {
	return boolOption(ParamAll, enabled)
}

// WithAllEvent sets `allEvent`.
func (s *OptionService) WithAllEvent(enabled bool) *APIParamOption {
	return boolOption(ParamAllEvent, enabled)
}

// WithAllowAddItem sets `allowAddItem`.
func (s *OptionService) WithAllowAddItem(allowAddItem bool) *APIParamOption {
	return boolOption(ParamAllowAddItem, allowAddItem)
}

// WithAllowInput sets `allowInput`.
func (s *OptionService) WithAllowInput(allowInput bool) *APIParamOption {
	return boolOption(ParamAllowInput, allowInput)
}

// WithArchived sets `archived`.
func (s *OptionService) WithArchived(enabled bool) *APIParamOption {
	return boolOption(ParamArchived, enabled)
}

// WithAttachment sets `attachment`.
func (s *OptionService) WithAttachment(enabled bool) *APIParamOption {
	return boolOption(ParamAttachment, enabled)
}

// WithChartEnabled sets `chartEnabled`.
func (s *OptionService) WithChartEnabled(enabled bool) *APIParamOption {
	return boolOption(ParamChartEnabled, enabled)
}

// WithExcludeGroupMembers sets `excludeGroupMembers`.
func (s *OptionService) WithExcludeGroupMembers(enabled bool) *APIParamOption {
	return boolOption(ParamExcludeGroupMembers, enabled)
}

// WithGrandchildIssueEnabled sets `grandchildIssueEnabled`.
func (s *OptionService) WithGrandchildIssueEnabled(enabled bool) *APIParamOption {
	return boolOption(ParamGrandchildIssueEnabled, enabled)
}

// WithHasDueDate sets `hasDueDate`.
func (s *OptionService) WithHasDueDate(enabled bool) *APIParamOption {
	return boolOption(ParamHasDueDate, enabled)
}

// WithMailNotify sets `mailNotify`.
func (s *OptionService) WithMailNotify(enabled bool) *APIParamOption {
	return boolOption(ParamMailNotify, enabled)
}

// WithProjectLeaderCanEditProjectLeader sets `projectLeaderCanEditProjectLeader`.
func (s *OptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) *APIParamOption {
	return boolOption(ParamProjectLeaderCanEditProjectLeader, enabled)
}

// WithRequired sets `required`.
func (s *OptionService) WithRequired(required bool) *APIParamOption {
	return boolOption(ParamRequired, required)
}

// WithSharedFile sets `sharedFile`.
func (s *OptionService) WithSharedFile(enabled bool) *APIParamOption {
	return boolOption(ParamSharedFile, enabled)
}

// WithSubtaskingEnabled sets `subtaskingEnabled`.
func (s *OptionService) WithSubtaskingEnabled(enabled bool) *APIParamOption {
	return boolOption(ParamSubtaskingEnabled, enabled)
}

func boolOption(paramType APIParamOptionType, enabled bool) *APIParamOption {
	return &APIParamOption{
		Type:    paramType,
		SetFunc: setBoolFunc(paramType, enabled),
	}
}

func setBoolFunc(key APIParamOptionType, value bool) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), strconv.FormatBool(value))
		return nil
	}
}
