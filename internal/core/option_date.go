package core

func (s *OptionService) WithCreatedSince(date string) RequestOption {
	return dateFormatStringOption(ParamCreatedSince, date)
}

func (s *OptionService) WithCreatedUntil(date string) RequestOption {
	return dateFormatStringOption(ParamCreatedUntil, date)
}

func (s *OptionService) WithDueDate(date string) RequestOption {
	return dateFormatStringOption(ParamDueDate, date)
}

func (s *OptionService) WithUpdatedSince(date string) RequestOption {
	return dateFormatStringOption(ParamUpdatedSince, date)
}

func (s *OptionService) WithUpdatedUntil(date string) RequestOption {
	return dateFormatStringOption(ParamUpdatedUntil, date)
}

func (s *OptionService) WithStartDate(date string) RequestOption {
	return dateFormatStringOption(ParamStartDate, date)
}

func (s *OptionService) WithStartDateSince(date string) RequestOption {
	return dateFormatStringOption(ParamStartDateSince, date)
}

func (s *OptionService) WithStartDateUntil(date string) RequestOption {
	return dateFormatStringOption(ParamStartDateUntil, date)
}

func (s *OptionService) WithDueDateSince(date string) RequestOption {
	return dateFormatStringOption(ParamDueDateSince, date)
}

func (s *OptionService) WithDueDateUntil(date string) RequestOption {
	return dateFormatStringOption(ParamDueDateUntil, date)
}

func (s *OptionService) WithReleaseDueDate(date string) RequestOption {
	return dateFormatStringOption(ParamReleaseDueDate, date)
}
