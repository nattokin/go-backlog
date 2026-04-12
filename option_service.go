package backlog

import (
	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
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

// UserOptionService provides a domain-specific set of option builders
// for operations within the UserService.
type UserOptionService struct {
	base *core.OptionService
}

func (s *UserOptionService) WithMailAddress(mail string) RequestOption {
	return s.base.WithMailAddress(mail)
}

func (s *UserOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

func (s *UserOptionService) WithPassword(password string) RequestOption {
	return s.base.WithPassword(password)
}

func (s *UserOptionService) WithRoleType(role Role) RequestOption {
	return s.base.WithRoleType(role)
}

func (s *UserOptionService) WithSendMail(enabled bool) RequestOption {
	return s.base.WithSendMail(enabled)
}

func (s *UserOptionService) WithUserID(id int) RequestOption {
	return s.base.WithUserID(id)
}

//
// ──────────────────────────────────────────────────────────────
//  WikiOptionService
// ──────────────────────────────────────────────────────────────
//

// WikiOptionService provides a domain-specific set of option builders
// for operations within the WikiService.
type WikiOptionService = wiki.WikiOptionService
