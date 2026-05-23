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
	Check() ValidationResult
	Set(url.Values) error
}

type ValidationResult interface {
	Valid() bool
	Target() string
	Message() string
}

type apiParamOption struct {
	key   func() string
	check func() core.ValidationResult
	set   func(url.Values) error
}

func (o *apiParamOption) Key() string                  { return o.key() }
func (o *apiParamOption) Check() core.ValidationResult { return o.check() }
func (o *apiParamOption) Set(v url.Values) error       { return o.set(v) }

// ──────────────────────────────────────────────────────────────
//  ActivityOptionService
// ──────────────────────────────────────────────────────────────

// ActivityOptionService provides option builders for activity list operations.
type ActivityOptionService struct {
	base *core.OptionService
}

// WithActivityTypeIDs filters activities by type IDs.
func (s *ActivityOptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	option := s.base.WithActivityTypeIDs(typeIDs)
	return &apiParamOption{
		key:   option.Key,
		check: func() ValidationResult { return option.Check() },
		set:   option.Set,
	}
}

// WithMinID filters activities whose ID is greater than or equal to id.
func (s *ActivityOptionService) WithMinID(id int) RequestOption {
	return s.base.WithMinActivityTypeID(id)
}

// WithMaxID filters activities whose ID is less than or equal to id.
func (s *ActivityOptionService) WithMaxID(id int) RequestOption {
	return s.base.WithMaxActivityTypeID(id)
}

// WithCount sets the number of activities to retrieve.
func (s *ActivityOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithOrder sets the sort order of results.
func (s *ActivityOptionService) WithOrder(order Order) RequestOption {
	return s.base.WithOrder(string(order))
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

func toCoreOptions(opts []RequestOption) []core.RequestOption {
	coreOpts := make([]core.RequestOption, len(opts))
	for i, o := range opts {
		coreOpts[i] = toCoreOption(o)
	}
	return coreOpts
}

func toCoreOption(option RequestOption) core.RequestOption {
	return &apiParamOption{
		key:   option.Key,
		check: func() core.ValidationResult { return option.Check() },
		set:   option.Set,
	}
}
