package option

import (
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

func (s *OptionService) WithCreatedSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamCreatedSince, date)
}

func (s *OptionService) WithCreatedUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamCreatedUntil, date)
}

func (s *OptionService) WithDueDate(date string) *APIParamOption {
	return clearableDateFormatStringOption(ParamDueDate, date)
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

func (s *OptionService) WithMaxDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamMax, date)
}

func (s *OptionService) WithMinDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamMin, date)
}

func (s *OptionService) WithReleaseDueDate(date string) *APIParamOption {
	return clearableDateFormatStringOption(ParamReleaseDueDate, date)
}

func (s *OptionService) WithStartDate(date string) *APIParamOption {
	return clearableDateFormatStringOption(ParamStartDate, date)
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
		CheckFunc: func() *validation.Error {
			return validate.ValidateDateFormat(paramType.Value(), date)
		},
		SetFunc: setStringFunc(paramType, date),
	}
}

// clearableDateFormatStringOption is like dateFormatStringOption but also
// accepts an empty string, which the Backlog API interprets as clearing the date.
func clearableDateFormatStringOption(paramType APIParamOptionType, date string) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *validation.Error {
			if date == "" {
				return nil
			}
			return validate.ValidateDateFormat(paramType.Value(), date)
		},
		SetFunc: setStringFunc(paramType, date),
	}
}
