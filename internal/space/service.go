package space

import (
	"github.com/nattokin/go-backlog/internal/core"
)

type SpaceService struct {
	method *core.Method
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewSpaceService(method *core.Method, option *core.OptionService) *SpaceService {
	return &SpaceService{
		method: method,
	}
}
