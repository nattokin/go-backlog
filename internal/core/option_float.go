package core

import (
	"net/url"
	"strconv"
)

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

// setFloat64Func returns a SetFunc that serializes a float64 without trailing zeros.
func setFloat64Func(key APIParamOptionType, value float64) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), strconv.FormatFloat(value, 'f', -1, 64))
		return nil
	}
}
