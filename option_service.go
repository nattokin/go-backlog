package backlog

import (
	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/user"
	"github.com/nattokin/go-backlog/internal/wiki"
)

//
// ──────────────────────────────────────────────────────────────
//  ActivityOptionService
// ──────────────────────────────────────────────────────────────
//

type ActivityOptionService = activity.ActivityOptionService

//
// ──────────────────────────────────────────────────────────────
//  ProjectOptionService
// ──────────────────────────────────────────────────────────────
//

// ProjectOptionService provides a domain-specific set of option builders
// for operations within the ProjectService.
type ProjectOptionService struct {
	base *core.OptionService
}

func (s *ProjectOptionService) WithAll(enabled bool) RequestOption {
	return s.base.WithAll(enabled)
}

func (s *ProjectOptionService) WithArchived(enabled bool) RequestOption {
	return s.base.WithArchived(enabled)
}

func (s *ProjectOptionService) WithChartEnabled(enabled bool) RequestOption {
	return s.base.WithChartEnabled(enabled)
}

func (s *ProjectOptionService) WithKey(key string) RequestOption {
	return s.base.WithKey(key)
}

func (s *ProjectOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

func (s *ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return s.base.WithProjectLeaderCanEditProjectLeader(enabled)
}

func (s *ProjectOptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return s.base.WithSubtaskingEnabled(enabled)
}

func (s *ProjectOptionService) WithTextFormattingRule(format Format) RequestOption {
	return s.base.WithTextFormattingRule(format)
}

//
// ──────────────────────────────────────────────────────────────
//  UserOptionService
// ──────────────────────────────────────────────────────────────
//

type UserOptionService = user.UserOptionService

//
// ──────────────────────────────────────────────────────────────
//  WikiOptionService
// ──────────────────────────────────────────────────────────────
//

type WikiOptionService = wiki.WikiOptionService
