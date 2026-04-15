package issue

import (
	"github.com/nattokin/go-backlog/internal/core"
)

type Service struct {
	method *core.Method
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
