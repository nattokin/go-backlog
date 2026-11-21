package option

import (
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// WithCreatedSince sets `createdSince`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithCreatedSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamCreatedSince, date)
}

// WithCreatedUntil sets `createdUntil`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithCreatedUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamCreatedUntil, date)
}

// WithDueDate sets `dueDate`.
// The date must be formatted as yyyy-MM-dd; an empty string clears the date.
func (s *OptionService) WithDueDate(date string) *APIParamOption {
	return clearableDateFormatStringOption(ParamDueDate, date)
}

// WithDueDateSince sets `dueDateSince`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithDueDateSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamDueDateSince, date)
}

// WithDueDateUntil sets `dueDateUntil`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithDueDateUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamDueDateUntil, date)
}

// WithInitialDate sets `initialDate` for Date type custom fields.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithInitialDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamInitialDate, date)
}

// WithMaxDate sets `max` for Date type custom fields.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithMaxDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamMax, date)
}

// WithMinDate sets `min` for Date type custom fields.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithMinDate(date string) *APIParamOption {
	return dateFormatStringOption(ParamMin, date)
}

// WithReleaseDueDate sets `releaseDueDate`.
// The date must be formatted as yyyy-MM-dd; an empty string clears the date.
func (s *OptionService) WithReleaseDueDate(date string) *APIParamOption {
	return clearableDateFormatStringOption(ParamReleaseDueDate, date)
}

// WithSince sets `since`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamSince, date)
}

// WithStartDate sets `startDate`.
// The date must be formatted as yyyy-MM-dd; an empty string clears the date.
func (s *OptionService) WithStartDate(date string) *APIParamOption {
	return clearableDateFormatStringOption(ParamStartDate, date)
}

// WithStartDateSince sets `startDateSince`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithStartDateSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamStartDateSince, date)
}

// WithStartDateUntil sets `startDateUntil`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithStartDateUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamStartDateUntil, date)
}

// WithUntil sets `until`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithUntil(date string) *APIParamOption {
	return dateFormatStringOption(ParamUntil, date)
}

// WithUpdatedSince sets `updatedSince`.
// The date must be formatted as yyyy-MM-dd.
func (s *OptionService) WithUpdatedSince(date string) *APIParamOption {
	return dateFormatStringOption(ParamUpdatedSince, date)
}

// WithUpdatedUntil sets `updatedUntil`.
// The date must be formatted as yyyy-MM-dd.
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
