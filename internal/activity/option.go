package activity

import (
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// ActivityOptionService provides a domain-specific set of option builders
// for operations within the ActivityService.
type ActivityOptionService struct {
	base *core.OptionService
}

func (s *ActivityOptionService) WithActivityTypeIDs(typeIDs []int) core.RequestOption {
	return s.base.WithActivityTypeIDs(typeIDs)
}

func (s *ActivityOptionService) WithMinID(id int) core.RequestOption {
	return s.base.WithMinID(id)
}

func (s *ActivityOptionService) WithMaxID(id int) core.RequestOption {
	return s.base.WithMaxID(id)
}

func (s *ActivityOptionService) WithCount(count int) core.RequestOption {
	return s.base.WithCount(count)
}

func (s *ActivityOptionService) WithOrder(order model.Order) core.RequestOption {
	return s.base.WithOrder(order)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewActivityOptionService(option *core.OptionService) *ActivityOptionService {
	return &ActivityOptionService{base: option}
}
