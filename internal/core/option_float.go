package core

import (
	"net/url"
	"strconv"
)

// WithInitialValue returns an option to set the `initialValue` parameter for Number type custom fields.
func (s *OptionService) WithInitialValue(initialValue float64) RequestOption {
	return &APIParamOption{
		Type:    ParamInitialValue,
		SetFunc: setFloat64Func(ParamInitialValue, initialValue),
	}
}

// WithMax returns an option to set the `max` parameter for Number type custom fields.
func (s *OptionService) WithMax(max float64) RequestOption {
	return &APIParamOption{
		Type:    ParamMax,
		SetFunc: setFloat64Func(ParamMax, max),
	}
}

// WithMin returns an option to set the `min` parameter for Number type custom fields.
func (s *OptionService) WithMin(min float64) RequestOption {
	return &APIParamOption{
		Type:    ParamMin,
		SetFunc: setFloat64Func(ParamMin, min),
	}
}

//
// ──────────────────────────────────────────────────────────────
//  SetFunc factories
// ──────────────────────────────────────────────────────────────
//

// setFloat64Func returns a SetFunc that calls v.Set with the float64 formatted without trailing zeros.
func setFloat64Func(key APIParamOptionType, value float64) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), strconv.FormatFloat(value, 'f', -1, 64))
		return nil
	}
}
