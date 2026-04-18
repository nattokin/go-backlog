package core

// WithAll returns an option to set the `all` parameter.
func (s *OptionService) WithAll(enabled bool) RequestOption {
	return boolOption(ParamAll, enabled)
}

// WithArchived returns an option to set the `archived` parameter.
func (s *OptionService) WithArchived(enabled bool) RequestOption {
	// apiArchived and queryArchived share the same string value "archived",
	// so we use apiArchived as the canonical type here and accept both in services.
	return boolOption(ParamArchived, enabled)
}

// WithAttachment returns an option to include only issues with attachments.
func (s *OptionService) WithAttachment(enabled bool) RequestOption {
	return boolOption(ParamAttachment, enabled)
}

// WithChartEnabled returns a option that sets the `chartEnabled` field.
func (s *OptionService) WithChartEnabled(enabled bool) RequestOption {
	return boolOption(ParamChartEnabled, enabled)
}

// WithHasDueDate returns an option to exclude issues without a due date.
// Note: Setting this to true is not supported by the Backlog API and will result in an error.
func (s *OptionService) WithHasDueDate(enabled bool) RequestOption {
	return boolOption(ParamHasDueDate, enabled)
}

// WithMailNotify returns a option that sets the `mailNotify` field.
func (s *OptionService) WithMailNotify(enabled bool) RequestOption {
	return boolOption(ParamMailNotify, enabled)
}

// WithProjectLeaderCanEditProjectLeader returns a option.
func (s *OptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return boolOption(ParamProjectLeaderCanEditProjectLeader, enabled)
}

// WithSendMail returns a option to specify whether to send an invitation email.
func (s *OptionService) WithSendMail(enabled bool) RequestOption {
	return boolOption(ParamSendMail, enabled)
}

// WithSharedFile returns an option to include only issues with shared files.
func (s *OptionService) WithSharedFile(enabled bool) RequestOption {
	return boolOption(ParamSharedFile, enabled)
}

// WithSubtaskingEnabled returns a option that sets the `subtaskingEnabled` field.
func (s *OptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return boolOption(ParamSubtaskingEnabled, enabled)
}
