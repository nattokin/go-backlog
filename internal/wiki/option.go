package wiki

import "github.com/nattokin/go-backlog/internal/core"

//
// ──────────────────────────────────────────────────────────────
//  WikiOptionService
// ──────────────────────────────────────────────────────────────
//

// WikiOptionService provides a domain-specific set of option builders
// for operations within the WikiService.
type WikiOptionService struct {
	base *core.OptionService
}

func (s *WikiOptionService) WithKeyword(keyword string) core.RequestOption {
	return s.base.WithKeyword(keyword)
}

func (s *WikiOptionService) WithContent(content string) core.RequestOption {
	return s.base.WithContent(content)
}

func (s *WikiOptionService) WithMailNotify(enabled bool) core.RequestOption {
	return s.base.WithMailNotify(enabled)
}

func (s *WikiOptionService) WithName(name string) core.RequestOption {
	return s.base.WithName(name)
}
