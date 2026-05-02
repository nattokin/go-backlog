package core

// WithTypeID returns an option to set the `typeId` parameter.
func (s *OptionService) WithTypeID(typeID int) RequestOption {
	return positiveIntOption(ParamTypeID, typeID)
}
