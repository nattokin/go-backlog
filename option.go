package backlog

import (
	"net/url"
)

//
// ──────────────────────────────────────────────────────────────
//  API Option Type
// ──────────────────────────────────────────────────────────────
//

// apiParamOptionType represents the distinct parameter keys for Backlog API requests.
type apiParamOptionType string

// Value returns the string representation of the parameter key for the API request.
func (t apiParamOptionType) Value() string {
	return string(t)
}

//
// ──────────────────────────────────────────────────────────────
//  RequestOption interface
// ──────────────────────────────────────────────────────────────
//

// RequestOption defines a common interface for all option types.
// It allows unified validation and application handling across different request-level options.
// Callers can implement this interface to provide custom options (e.g. for mocking in tests).
type RequestOption interface {
	Key() string
	Check() error
	Set(url.Values) error
}

//
// ──────────────────────────────────────────────────────────────
//  apiOption — unified internal option type
// ──────────────────────────────────────────────────────────────
//

// apiParamOption is the internal implementation of RequestOption.
//
// It encapsulates the parameter type together with optional validation
// and the logic required to apply the value to API request parameters.
// OptionService builder methods return instances of this struct.
type apiParamOption struct {
	t         apiParamOptionType     // canonical API parameter type
	checkFunc func() error           // optional validation executed before applying the option
	setFunc   func(url.Values) error // applies the option value to the provided values
}

// Key returns the API parameter key associated with this option.
func (o *apiParamOption) Key() string {
	return o.t.Value()
}

// Check validates the option by executing its checkFunc, if defined.
func (o *apiParamOption) Check() error {
	if o.checkFunc != nil {
		return o.checkFunc()
	}
	return nil
}

// Set applies the option value to the given url.Values.
func (o *apiParamOption) Set(v url.Values) error {
	if o.setFunc == nil {
		panic("option has no setter")
	}
	return o.setFunc(v)
}

//
// ──────────────────────────────────────────────────────────────
//  Internal Helpers
// ──────────────────────────────────────────────────────────────
//

// validateOption checks whether the given option key is allowed
// for the current API operation.
func validateOption(optionKey string, validOptions []apiParamOptionType) error {
	for _, valid := range validOptions {
		if optionKey == valid.Value() {
			return nil
		}
	}
	return newInvalidOptionKeyError(optionKey, validOptions)
}

// applyOptions validates and applies request options to the given url.Values.
func applyOptions(v url.Values, validTypes []apiParamOptionType, opts ...RequestOption) error {
	for _, opt := range opts {
		if err := validateOption(opt.Key(), validTypes); err != nil {
			return err
		}
		if err := opt.Check(); err != nil {
			return err
		}
		if err := opt.Set(v); err != nil {
			return err
		}
	}
	return nil
}

// hasRequiredOption checks whether the provided options include at least one of the required form types.
func hasRequiredOption(options []RequestOption, requiredTypes []apiParamOptionType) bool {
	for _, opt := range options {
		for _, requiredType := range requiredTypes {
			if opt.Key() == requiredType.Value() {
				return true
			}
		}
	}
	return false
}
