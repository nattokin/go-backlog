package core

import (
	"fmt"
	"net/url"
	"strconv"
)

// WithActualHours returns an option to set the `actualHours` parameter.
// Any positive float64 value is accepted (e.g. 2.5).
func (s *OptionService) WithActualHours(hours float64) RequestOption {
	return positiveFloat64Option(ParamActualHours, hours)
}

// WithEstimatedHours returns an option to set the `estimatedHours` parameter.
// Any positive float64 value is accepted (e.g. 2.5).
func (s *OptionService) WithEstimatedHours(hours float64) RequestOption {
	return positiveFloat64Option(ParamEstimatedHours, hours)
}

// WithInitialValue sets `initialValue` for Number type custom fields.
// Any float64 value is accepted, including zero and negative values.
func (s *OptionService) WithInitialValue(initialValue float64) RequestOption {
	return &APIParamOption{
		Type:    ParamInitialValue,
		SetFunc: setFloat64Func(ParamInitialValue, initialValue),
	}
}

// WithMax sets `max` for Number type custom fields.
// Any float64 value is accepted, including zero and negative values.
func (s *OptionService) WithMax(max float64) RequestOption {
	return &APIParamOption{
		Type:    ParamMax,
		SetFunc: setFloat64Func(ParamMax, max),
	}
}

// WithMin sets `min` for Number type custom fields.
// Any float64 value is accepted, including zero and negative values.
func (s *OptionService) WithMin(min float64) RequestOption {
	return &APIParamOption{
		Type:    ParamMin,
		SetFunc: setFloat64Func(ParamMin, min),
	}
}

// positiveFloat64Option builds a RequestOption that validates a float64 is > 0 and sets it.
func positiveFloat64Option(paramType APIParamOptionType, value float64) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if value <= 0 {
				return NewValidationError(fmt.Sprintf("invalid %s: must be greater than 0", paramType.Value()))
			}
			return nil
		},
		SetFunc: setFloat64Func(paramType, value),
	}
}

// setFloat64Func returns a SetFunc that serializes a float64 without trailing zeros.
func setFloat64Func(key APIParamOptionType, value float64) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), strconv.FormatFloat(value, 'f', -1, 64))
		return nil
	}
}
