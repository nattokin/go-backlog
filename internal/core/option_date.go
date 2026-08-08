package core

import (
	"fmt"
	"regexp"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func (s *OptionService) WithCreatedSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamCreatedSince, date)
}

func (s *OptionService) WithCreatedUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamCreatedUntil, date)
}

func (s *OptionService) WithDueDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamDueDate, date)
}

func (s *OptionService) WithDueDateSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamDueDateSince, date)
}

func (s *OptionService) WithDueDateUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamDueDateUntil, date)
}

func (s *OptionService) WithInitialDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamInitialDate, date)
}

func (s *OptionService) WithInitialDateMax(date string) *APIParamOption {
	return dateFormatStringOption(ParamMax, date)
}

func (s *OptionService) WithInitialDateMin(date string) *APIParamOption {
	return dateFormatStringOption(ParamMin, date)
}

func (s *OptionService) WithReleaseDueDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamReleaseDueDate, date)
}

func (s *OptionService) WithStartDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamStartDate, date)
}

func (s *OptionService) WithStartDateSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamStartDateSince, date)
}

func (s *OptionService) WithStartDateUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamStartDateUntil, date)
}

func (s *OptionService) WithUpdatedSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamUpdatedSince, date)
}

func (s *OptionService) WithUpdatedUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamUpdatedUntil, date)
}

func dateFormatStringOption(paramType APIParamOptionType, date string) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *ValidationError {
			if !datePattern.MatchString(date) {
				return NewValidationError(paramType.Value(), fmt.Sprintf("%s must be formatted as yyyy-MM-dd, got %q", paramType.Value(), date))
			}
			return nil
		},
		SetFunc: setStringFunc(paramType, date),
	}
}
