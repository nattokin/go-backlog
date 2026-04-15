package issue

import (
	"github.com/nattokin/go-backlog/internal/core"
)

type IssueService struct {
	method *core.Method
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewIssueService(method *core.Method, option *core.OptionService) *IssueService {
	return &IssueService{
		method: method,
	}
}
