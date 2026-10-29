package backlog

import (
	"net/url"

	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validation"
)

// RequestOption defines a common interface for all option types.
// It allows unified validation and application handling across different request-level options.
// Callers can implement this interface to provide custom options (e.g. for mocking in tests).
//
// A method applies its options in order. For each option, it checks Key, then
// calls Check, and calls Set only when Check reports no error.
type RequestOption interface {
	// Key returns the name of the API parameter the option sets, such as
	// "keyword" or "categoryId[]". If the key is not one the method accepts,
	// the method returns an *InvalidOptionKeyError without sending a request.
	Key() string
	// Check validates the option's value and returns nil if it is valid.
	// Validation errors from all options are collected and returned together;
	// use [NewValidationError] to create one.
	Check() *ValidationError
	// Set writes the option's value to the request parameters. If Set returns
	// an error, the method returns it without sending a request.
	Set(url.Values) error
}

// requestOption is the internal implementation of RequestOption that wraps
// a *option.APIParamOption, converting *validation.Error to *ValidationError.
type requestOption struct {
	opt *option.APIParamOption
}

func (o *requestOption) Key() string { return o.opt.Key() }
func (o *requestOption) Check() *ValidationError {
	if ve := o.opt.Check(); ve != nil {
		return &ValidationError{target: ve.Target(), message: ve.Message()}
	}
	return nil
}
func (o *requestOption) Set(v url.Values) error { return o.opt.Set(v) }

// ──────────────────────────────────────────────────────────────
//  ActivityOptionService
// ──────────────────────────────────────────────────────────────

// ActivityOptionService provides option builders for activity list operations.
type ActivityOptionService struct {
	base *option.OptionService
}

// WithActivityTypeIDs filters activities by type IDs.
func (s *ActivityOptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return &requestOption{opt: s.base.WithActivityTypeIDs(typeIDs)}
}

// WithMinID filters activities whose ID is greater than or equal to id.
func (s *ActivityOptionService) WithMinID(id int) RequestOption {
	return &requestOption{opt: s.base.WithMinID(id)}
}

// WithMaxID filters activities whose ID is less than or equal to id.
func (s *ActivityOptionService) WithMaxID(id int) RequestOption {
	return &requestOption{opt: s.base.WithMaxID(id)}
}

// WithCount sets the number of activities to retrieve (1-100).
func (s *ActivityOptionService) WithCount(count int) RequestOption {
	return &requestOption{opt: s.base.WithCount(count)}
}

// WithOrder sets the sort order of results.
func (s *ActivityOptionService) WithOrder(order Order) RequestOption {
	return &requestOption{opt: s.base.WithOrder(string(order))}
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newActivityOptionService(option *option.OptionService) *ActivityOptionService {
	return &ActivityOptionService{base: option}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

// toInnerOptions converts a slice of RequestOption to []*option.APIParamOption.
// A nil RequestOption is passed through as a nil *option.APIParamOption so that
// the internal layer can detect it and return InvalidOptionError.
func toInnerOptions(opts []RequestOption) []*option.APIParamOption {
	innerOpts := make([]*option.APIParamOption, len(opts))
	for i, o := range opts {
		if o == nil {
			innerOpts[i] = nil
			continue
		}
		innerOpts[i] = toInnerOption(o)
	}
	return innerOpts
}

// toInnerOption converts a backlog.RequestOption to *option.APIParamOption so it
// can be passed to internal domain service endpoints.
func toInnerOption(opt RequestOption) *option.APIParamOption {
	return &option.APIParamOption{
		KeyFunc: opt.Key,
		CheckFunc: func() *validation.Error {
			if ve := opt.Check(); ve != nil {
				return validation.NewError(ve.Target(), ve.Message())
			}
			return nil
		},
		SetFunc: opt.Set,
	}
}
