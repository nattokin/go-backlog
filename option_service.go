package backlog

//
// ──────────────────────────────────────────────────────────────
//  ActivityOptionService
// ──────────────────────────────────────────────────────────────
//

// ActivityOptionService provides a domain-specific set of option builders
// for operations within the ActivityService.
type ActivityOptionService struct {
	registry *OptionService
}

func (s *ActivityOptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return s.registry.WithActivityTypeIDs(typeIDs)
}

func (s *ActivityOptionService) WithMinID(id int) RequestOption {
	return s.registry.WithMinID(id)
}

func (s *ActivityOptionService) WithMaxID(id int) RequestOption {
	return s.registry.WithMaxID(id)
}

func (s *ActivityOptionService) WithCount(count int) RequestOption {
	return s.registry.WithCount(count)
}

func (s *ActivityOptionService) WithOrder(order Order) RequestOption {
	return s.registry.WithOrder(order)
}

//
// ──────────────────────────────────────────────────────────────
//  ProjectOptionService
// ──────────────────────────────────────────────────────────────
//

// ProjectOptionService provides a domain-specific set of option builders
// for operations within the ProjectService.
type ProjectOptionService struct {
	registry *OptionService
}

func (s *ProjectOptionService) WithAll(enabled bool) RequestOption {
	return s.registry.WithAll(enabled)
}

func (s *ProjectOptionService) WithArchived(enabled bool) RequestOption {
	return s.registry.WithArchived(enabled)
}

func (s *ProjectOptionService) WithChartEnabled(enabled bool) RequestOption {
	return s.registry.WithChartEnabled(enabled)
}

func (s *ProjectOptionService) WithKey(key string) RequestOption {
	return s.registry.WithKey(key)
}

func (s *ProjectOptionService) WithName(name string) RequestOption {
	return s.registry.WithName(name)
}

func (s *ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return s.registry.WithProjectLeaderCanEditProjectLeader(enabled)
}

func (s *ProjectOptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return s.registry.WithSubtaskingEnabled(enabled)
}

func (s *ProjectOptionService) WithTextFormattingRule(format Format) RequestOption {
	return s.registry.WithTextFormattingRule(format)
}

//
// ──────────────────────────────────────────────────────────────
//  UserOptionService
// ──────────────────────────────────────────────────────────────
//

// UserOptionService provides a domain-specific set of option builders
// for operations within the UserService.
type UserOptionService struct {
	registry *OptionService
}

func (s *UserOptionService) WithMailAddress(mail string) RequestOption {
	return s.registry.WithMailAddress(mail)
}

func (s *UserOptionService) WithName(name string) RequestOption {
	return s.registry.WithName(name)
}

func (s *UserOptionService) WithPassword(password string) RequestOption {
	return s.registry.WithPassword(password)
}

func (s *UserOptionService) WithRoleType(role Role) RequestOption {
	return s.registry.WithRoleType(role)
}

func (s *UserOptionService) WithSendMail(enabled bool) RequestOption {
	return s.registry.WithSendMail(enabled)
}

func (s *UserOptionService) WithUserID(id int) RequestOption {
	return s.registry.WithUserID(id)
}

//
// ──────────────────────────────────────────────────────────────
//  WikiOptionService
// ──────────────────────────────────────────────────────────────
//

// WikiOptionService provides a domain-specific set of option builders
// for operations within the WikiService.
type WikiOptionService struct {
	registry *OptionService
}

func (s *WikiOptionService) WithKeyword(keyword string) RequestOption {
	return s.registry.WithKeyword(keyword)
}

func (s *WikiOptionService) WithContent(content string) RequestOption {
	return s.registry.WithContent(content)
}

func (s *WikiOptionService) WithMailNotify(enabled bool) RequestOption {
	return s.registry.WithMailNotify(enabled)
}

func (s *WikiOptionService) WithName(name string) RequestOption {
	return s.registry.WithName(name)
}
