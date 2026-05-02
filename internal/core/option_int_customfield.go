package core

// WithTypeID returns an option to set the `typeId` parameter.
func (s *OptionService) WithTypeID(typeID int) RequestOption {
	return positiveIntOption(ParamTypeID, typeID)
}

// WithInitialValueType returns an option to set the `initialValueType` parameter for Date type custom fields.
// 0: Today, 1: Specified date, 2: Today + initialShift days.
func (s *OptionService) WithInitialValueType(initialValueType int) RequestOption {
	return intRangeOption(ParamInitialValueType, initialValueType, 0, 2)
}

// WithInitialShift returns an option to set the `initialShift` parameter for Date type custom fields.
// Used when initialValueType is 2 (today + N days).
func (s *OptionService) WithInitialShift(days int) RequestOption {
	return &APIParamOption{
		Type:    ParamInitialShift,
		SetFunc: setIntFunc(ParamInitialShift, days),
	}
}
