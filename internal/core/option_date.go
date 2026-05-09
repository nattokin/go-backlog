package core

import (
	"fmt"
	"regexp"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

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

// dateFormatStringOption builds a RequestOption that validates the string matches
// "yyyy-MM-dd" format before setting it.
func dateFormatStringOption(paramType APIParamOptionType, date string) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if !datePattern.MatchString(date) {
				return NewValidationError(fmt.Sprintf("%s must be formatted as yyyy-MM-dd, got %q", paramType.Value(), date))
			}
			return nil
		},
		SetFunc: setStringFunc(paramType, date),
	}
}
