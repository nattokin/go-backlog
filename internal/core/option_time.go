package core

import "time"

const issueDateFormat = "2006-01-02"

// WithCreatedSince returns an option to filter issues created on or after the given date.
func (s *OptionService) WithCreatedSince(t time.Time) RequestOption {
	return timeOption(ParamCreatedSince, t, issueDateFormat)
}

// WithCreatedUntil returns an option to filter issues created on or before the given date.
func (s *OptionService) WithCreatedUntil(t time.Time) RequestOption {
	return timeOption(ParamCreatedUntil, t, issueDateFormat)
}

// WithDueDate returns an option to set the `dueDate` parameter.
func (s *OptionService) WithDueDate(t time.Time) RequestOption {
	return timeOption(ParamDueDate, t, issueDateFormat)
}

// WithUpdatedSince returns an option to filter issues updated on or after the given date.
func (s *OptionService) WithUpdatedSince(t time.Time) RequestOption {
	return timeOption(ParamUpdatedSince, t, issueDateFormat)
}

// WithUpdatedUntil returns an option to filter issues updated on or before the given date.
func (s *OptionService) WithUpdatedUntil(t time.Time) RequestOption {
	return timeOption(ParamUpdatedUntil, t, issueDateFormat)
}

// WithStartDate returns an option to set the `startDate` parameter.
func (s *OptionService) WithStartDate(t time.Time) RequestOption {
	return timeOption(ParamStartDate, t, issueDateFormat)
}

// WithStartDateSince returns an option to filter issues with a start date on or after the given date.
func (s *OptionService) WithStartDateSince(t time.Time) RequestOption {
	return timeOption(ParamStartDateSince, t, issueDateFormat)
}

// WithStartDateUntil returns an option to filter issues with a start date on or before the given date.
func (s *OptionService) WithStartDateUntil(t time.Time) RequestOption {
	return timeOption(ParamStartDateUntil, t, issueDateFormat)
}

// WithDueDateSince returns an option to filter issues with a due date on or after the given date.
func (s *OptionService) WithDueDateSince(t time.Time) RequestOption {
	return timeOption(ParamDueDateSince, t, issueDateFormat)
}

// WithDueDateUntil returns an option to filter issues with a due date on or before the given date.
func (s *OptionService) WithDueDateUntil(t time.Time) RequestOption {
	return timeOption(ParamDueDateUntil, t, issueDateFormat)
}
