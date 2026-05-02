package core

// WithInitialValue returns an option to set the `initialValue` parameter for Number type custom fields.
func (s *OptionService) WithInitialValue(initialValue float64) RequestOption {
	return &APIParamOption{
		Type:    ParamInitialValue,
		SetFunc: setFloat64Func(ParamInitialValue, initialValue),
	}
}

// WithMax returns an option to set the `max` parameter for Number type custom fields.
func (s *OptionService) WithMax(max float64) RequestOption {
	return &APIParamOption{
		Type:    ParamMax,
		SetFunc: setFloat64Func(ParamMax, max),
	}
}

// WithMin returns an option to set the `min` parameter for Number type custom fields.
func (s *OptionService) WithMin(min float64) RequestOption {
	return &APIParamOption{
		Type:    ParamMin,
		SetFunc: setFloat64Func(ParamMin, min),
	}
}
