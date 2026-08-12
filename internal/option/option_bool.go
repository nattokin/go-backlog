package option

import (
	"net/url"
	"strconv"
)

func (s *OptionService) WithAll(enabled bool) *APIParamOption {
	return boolOption(ParamAll, enabled)
}

func (s *OptionService) WithAllEvent(enabled bool) *APIParamOption {
	return boolOption(ParamAllEvent, enabled)
}

func (s *OptionService) WithAllowAddItem(allowAddItem bool) *APIParamOption {
	return boolOption(ParamAllowAddItem, allowAddItem)
}

func (s *OptionService) WithAllowInput(allowInput bool) *APIParamOption {
	return boolOption(ParamAllowInput, allowInput)
}

func (s *OptionService) WithArchived(enabled bool) *APIParamOption {
	return boolOption(ParamArchived, enabled)
}

func (s *OptionService) WithAttachment(enabled bool) *APIParamOption {
	return boolOption(ParamAttachment, enabled)
}

func (s *OptionService) WithChartEnabled(enabled bool) *APIParamOption {
	return boolOption(ParamChartEnabled, enabled)
}

func (s *OptionService) WithExcludeGroupMembers(enabled bool) *APIParamOption {
	return boolOption(ParamExcludeGroupMembers, enabled)
}

func (s *OptionService) WithHasDueDate(enabled bool) *APIParamOption {
	return boolOption(ParamHasDueDate, enabled)
}

func (s *OptionService) WithMailNotify(enabled bool) *APIParamOption {
	return boolOption(ParamMailNotify, enabled)
}

func (s *OptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) *APIParamOption {
	return boolOption(ParamProjectLeaderCanEditProjectLeader, enabled)
}

func (s *OptionService) WithRequired(required bool) *APIParamOption {
	return boolOption(ParamRequired, required)
}

func (s *OptionService) WithSendMail(enabled bool) *APIParamOption {
	return boolOption(ParamSendMail, enabled)
}

func (s *OptionService) WithSharedFile(enabled bool) *APIParamOption {
	return boolOption(ParamSharedFile, enabled)
}

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
