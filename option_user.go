package backlog

import (
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// UserOptionService provides a domain-specific set of option builders
// for operations within the UserService.
type UserOptionService struct {
	base *core.OptionService
}

func (s *UserOptionService) WithMailAddress(mail string) core.RequestOption {
	return s.base.WithMailAddress(mail)
}

func (s *UserOptionService) WithName(name string) core.RequestOption {
	return s.base.WithName(name)
}

func (s *UserOptionService) WithPassword(password string) core.RequestOption {
	return s.base.WithPassword(password)
}

func (s *UserOptionService) WithRoleType(role model.Role) core.RequestOption {
	return s.base.WithRoleType(role)
}

func (s *UserOptionService) WithSendMail(enabled bool) core.RequestOption {
	return s.base.WithSendMail(enabled)
}

func (s *UserOptionService) WithUserID(id int) core.RequestOption {
	return s.base.WithUserID(id)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newUserOptionService(option *core.OptionService) *UserOptionService {
	return &UserOptionService{
		base: option,
	}
}
