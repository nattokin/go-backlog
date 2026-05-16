package core

import (
	"net/url"
	"strconv"
)

func (s *OptionService) WithAll(enabled bool) RequestOption {
	return boolOption(ParamAll, enabled)
}

func (s *OptionService) WithAllEvent(enabled bool) RequestOption {
	return boolOption(ParamAllEvent, enabled)
}

// WithAllowAddItem sets `allowAddItem` for List type custom fields.
// When true, users can add new items to the list from the issue form.
func (s *OptionService) WithAllowAddItem(allowAddItem bool) RequestOption {
	return boolOption(ParamAllowAddItem, allowAddItem)
}

// WithAllowInput sets `allowInput` for List type custom fields.
// When true, users can enter a free-text value in addition to selecting from the list.
func (s *OptionService) WithAllowInput(allowInput bool) RequestOption {
	return boolOption(ParamAllowInput, allowInput)
}

func (s *OptionService) WithArchived(enabled bool) RequestOption {
	return boolOption(ParamArchived, enabled)
}

func (s *OptionService) WithAttachment(enabled bool) RequestOption {
	return boolOption(ParamAttachment, enabled)
}

func (s *OptionService) WithChartEnabled(enabled bool) RequestOption {
	return boolOption(ParamChartEnabled, enabled)
}

// WithExcludeGroupMembers sets `excludeGroupMembers`.
// When true, users who joined the project only via group membership are excluded from the result.
func (s *OptionService) WithExcludeGroupMembers(enabled bool) RequestOption {
	return boolOption(ParamExcludeGroupMembers, enabled)
}

// WithHasDueDate sets `hasDueDate`.
// Note: Setting this to true is not supported by the Backlog API and will result in an error.
func (s *OptionService) WithHasDueDate(enabled bool) RequestOption {
	return boolOption(ParamHasDueDate, enabled)
}

func (s *OptionService) WithMailNotify(enabled bool) RequestOption {
	return boolOption(ParamMailNotify, enabled)
}

func (s *OptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return boolOption(ParamProjectLeaderCanEditProjectLeader, enabled)
}

func (s *OptionService) WithRequired(required bool) RequestOption {
	return boolOption(ParamRequired, required)
}

func (s *OptionService) WithSendMail(enabled bool) RequestOption {
	return boolOption(ParamSendMail, enabled)
}

func (s *OptionService) WithSharedFile(enabled bool) RequestOption {
	return boolOption(ParamSharedFile, enabled)
}

func (s *OptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return boolOption(ParamSubtaskingEnabled, enabled)
}

func boolOption(paramType APIParamOptionType, enabled bool) RequestOption {
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
