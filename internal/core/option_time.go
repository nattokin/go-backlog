package core

import (
	"net/url"
	"time"
)

// issueDateFormat is the date layout used by the Backlog API for date parameters.
const issueDateFormat = "2006-01-02"

func (s *OptionService) WithCreatedSince(t time.Time) RequestOption {
	return timeOption(ParamCreatedSince, t, issueDateFormat)
}

func (s *OptionService) WithCreatedUntil(t time.Time) RequestOption {
	return timeOption(ParamCreatedUntil, t, issueDateFormat)
}

func (s *OptionService) WithDueDate(t time.Time) RequestOption {
	return timeOption(ParamDueDate, t, issueDateFormat)
}

func (s *OptionService) WithUpdatedSince(t time.Time) RequestOption {
	return timeOption(ParamUpdatedSince, t, issueDateFormat)
}

func (s *OptionService) WithUpdatedUntil(t time.Time) RequestOption {
	return timeOption(ParamUpdatedUntil, t, issueDateFormat)
}

func (s *OptionService) WithStartDate(t time.Time) RequestOption {
	return timeOption(ParamStartDate, t, issueDateFormat)
}

func (s *OptionService) WithStartDateSince(t time.Time) RequestOption {
	return timeOption(ParamStartDateSince, t, issueDateFormat)
}

func (s *OptionService) WithStartDateUntil(t time.Time) RequestOption {
	return timeOption(ParamStartDateUntil, t, issueDateFormat)
}

func (s *OptionService) WithDueDateSince(t time.Time) RequestOption {
	return timeOption(ParamDueDateSince, t, issueDateFormat)
}

func (s *OptionService) WithDueDateUntil(t time.Time) RequestOption {
	return timeOption(ParamDueDateUntil, t, issueDateFormat)
}

func (s *OptionService) WithReleaseDueDate(t time.Time) RequestOption {
	return timeOption(ParamReleaseDueDate, t, issueDateFormat)
}

// timeOption builds a RequestOption that formats a time.Time value and sets it.
func timeOption(paramType APIParamOptionType, t time.Time, format string) RequestOption {
	return &APIParamOption{
		Type:    paramType,
		SetFunc: setTimeFunc(paramType, t, format),
	}
}

func setTimeFunc(key APIParamOptionType, t time.Time, format string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), t.Format(format))
		return nil
	}
}
