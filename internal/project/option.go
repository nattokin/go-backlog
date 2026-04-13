package project

import (
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// ProjectOptionService provides a domain-specific set of option builders
// for operations within the ProjectService.
type ProjectOptionService struct {
	base *core.OptionService
}

func (s *ProjectOptionService) WithAll(enabled bool) core.RequestOption {
	return s.base.WithAll(enabled)
}

func (s *ProjectOptionService) WithArchived(enabled bool) core.RequestOption {
	return s.base.WithArchived(enabled)
}

func (s *ProjectOptionService) WithChartEnabled(enabled bool) core.RequestOption {
	return s.base.WithChartEnabled(enabled)
}

func (s *ProjectOptionService) WithKey(key string) core.RequestOption {
	return s.base.WithKey(key)
}

func (s *ProjectOptionService) WithName(name string) core.RequestOption {
	return s.base.WithName(name)
}

func (s *ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) core.RequestOption {
	return s.base.WithProjectLeaderCanEditProjectLeader(enabled)
}

func (s *ProjectOptionService) WithSubtaskingEnabled(enabled bool) core.RequestOption {
	return s.base.WithSubtaskingEnabled(enabled)
}

func (s *ProjectOptionService) WithTextFormattingRule(format model.Format) core.RequestOption {
	return s.base.WithTextFormattingRule(format)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

// NewProjectOptionService returns a new ProjectOptionService.
func NewProjectOptionService(option *core.OptionService) *ProjectOptionService {
	return &ProjectOptionService{
		base: option,
	}
}
