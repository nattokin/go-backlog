package option

import (
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/validate"
)

// WithActualHours returns an option to set the `actualHours` parameter.
func (s *OptionService) WithActualHours(hours float64) *APIParamOption {
	return positiveFloat64Option(ParamActualHours, hours)
}

// WithEstimatedHours returns an option to set the `estimatedHours` parameter.
func (s *OptionService) WithEstimatedHours(hours float64) *APIParamOption {
	return positiveFloat64Option(ParamEstimatedHours, hours)
}

// WithInitialValue sets `initialValue` for Number type custom fields.
func (s *OptionService) WithInitialValue(initialValue float64) *APIParamOption {
	return &APIParamOption{
		Type:    ParamInitialValue,
		SetFunc: setFloat64Func(ParamInitialValue, initialValue),
	}
}

// WithMax sets `max` for Number type custom fields.
func (s *OptionService) WithMax(max float64) *APIParamOption {
	return &APIParamOption{
		Type:    ParamMax,
		SetFunc: setFloat64Func(ParamMax, max),
	}
}

// WithMin sets `min` for Number type custom fields.
func (s *OptionService) WithMin(min float64) *APIParamOption {
	return &APIParamOption{
		Type:    ParamMin,
		SetFunc: setFloat64Func(ParamMin, min),
	}
}

func positiveFloat64Option(paramType APIParamOptionType, value float64) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *core.ValidationError {
			return validate.ValidatePositiveFloat64(paramType.Value(), value)
		},
		SetFunc: setFloat64Func(paramType, value),
	}
}

func setFloat64Func(key APIParamOptionType, value float64) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), strconv.FormatFloat(value, 'f', -1, 64))
		return nil
	}
}
