package pullrequest

import (
	"github.com/nattokin/go-backlog/internal/core"
)

type PullRequestService struct {
	method *core.Method
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewPullRequestService(method *core.Method) *PullRequestService {
	return &PullRequestService{
		method: method,
	}
}
