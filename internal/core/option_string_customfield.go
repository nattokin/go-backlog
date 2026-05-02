package core

// WithUnit returns an option to set the `unit` parameter for Number type custom fields.
func (s *OptionService) WithUnit(unit string) RequestOption {
	return &APIParamOption{
		Type:    ParamUnit,
		SetFunc: setStringFunc(ParamUnit, unit),
	}
}

// WithInitialDateMin returns an option to set the `min` parameter for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDateMin(date string) RequestOption {
	return &APIParamOption{
		Type:    ParamMin,
		SetFunc: setStringFunc(ParamMin, date),
	}
}

// WithInitialDateMax returns an option to set the `max` parameter for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDateMax(date string) RequestOption {
	return &APIParamOption{
		Type:    ParamMax,
		SetFunc: setStringFunc(ParamMax, date),
	}
}

// WithInitialDate returns an option to set the `initialDate` parameter for Date type custom fields.
// The value must be formatted as "yyyy-MM-dd".
func (s *OptionService) WithInitialDate(date string) RequestOption {
	return &APIParamOption{
		Type:    ParamInitialDate,
		SetFunc: setStringFunc(ParamInitialDate, date),
	}
}
