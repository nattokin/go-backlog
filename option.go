package backlog

import (
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
)

// RequestOption defines a common interface for all option types.
// It allows unified validation and application handling across different request-level options.
// Callers can implement this interface to provide custom options (e.g. for mocking in tests).
type RequestOption interface {
	Key() string
	Check() *ValidationError
	Set(url.Values) error
}

// ──────────────────────────────────────────────────────────────
//  ActivityOptionService
// ──────────────────────────────────────────────────────────────

// ActivityOptionService provides option builders for activity list operations.
type ActivityOptionService struct {
	base *core.OptionService
}

// WithActivityTypeIDs filters activities by type IDs.
func (s *ActivityOptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return fromCoreOption(s.base.WithActivityTypeIDs(typeIDs))
}

// WithMinID filters activities whose ID is greater than or equal to id.
func (s *ActivityOptionService) WithMinID(id int) RequestOption {
	return fromCoreOption(s.base.WithMinActivityTypeID(id))
}

// WithMaxID filters activities whose ID is less than or equal to id.
func (s *ActivityOptionService) WithMaxID(id int) RequestOption {
	return fromCoreOption(s.base.WithMaxActivityTypeID(id))
}

// WithCount sets the number of activities to retrieve.
func (s *ActivityOptionService) WithCount(count int) RequestOption {
	return fromCoreOption(s.base.WithCount(count))
}

// WithOrder sets the sort order of results.
func (s *ActivityOptionService) WithOrder(order Order) RequestOption {
	return fromCoreOption(s.base.WithOrder(string(order)))
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newActivityOptionService(option *core.OptionService) *ActivityOptionService {
	return &ActivityOptionService{base: option}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

// fromCoreOption wraps a *core.APIParamOption as a RequestOption, converting
// *core.ValidationError to *ValidationError in Check().
func fromCoreOption(opt *core.APIParamOption) RequestOption {
	return &coreOptionWrapper{opt: opt}
}

type coreOptionWrapper struct {
	opt *core.APIParamOption
}

func (w *coreOptionWrapper) Key() string { return w.opt.Key() }
func (w *coreOptionWrapper) Check() *ValidationError {
	if ve := w.opt.Check(); ve != nil {
		return &ValidationError{target: ve.Target(), message: ve.Message()}
	}
	return nil
}
func (w *coreOptionWrapper) Set(v url.Values) error { return w.opt.Set(v) }

// toCoreOptions converts a slice of RequestOption to []*core.APIParamOption.
func toCoreOptions(opts []RequestOption) []*core.APIParamOption {
	coreOpts := make([]*core.APIParamOption, len(opts))
	for i, o := range opts {
		coreOpts[i] = toCoreOption(o)
	}
	return coreOpts
}

// toCoreOption converts a backlog.RequestOption to *core.APIParamOption so it
// can be passed to internal domain service endpoints.
func toCoreOption(option RequestOption) *core.APIParamOption {
	return &core.APIParamOption{
		KeyFunc: option.Key,
		CheckFunc: func() *core.ValidationError {
			if ve := option.Check(); ve != nil {
				return core.NewValidationError(ve.Target(), ve.Message())
			}
			return nil
		},
		SetFunc: option.Set,
	}
}
