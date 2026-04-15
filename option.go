package backlog

import (
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

type RequestOption interface {
	Key() string
	Check() error
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
	return s.base.WithActivityTypeIDs(typeIDs)
}

// WithMinID filters activities whose ID is greater than or equal to id.
func (s *ActivityOptionService) WithMinID(id int) RequestOption {
	return s.base.WithMinID(id)
}

// WithMaxID filters activities whose ID is less than or equal to id.
func (s *ActivityOptionService) WithMaxID(id int) RequestOption {
	return s.base.WithMaxID(id)
}

// WithCount sets the number of activities to retrieve.
func (s *ActivityOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithOrder sets the sort order of results.
func (s *ActivityOptionService) WithOrder(order model.Order) RequestOption {
	return s.base.WithOrder(order)
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
		coreOpts[i] = o
	}
	return coreOpts
}
