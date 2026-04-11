package backlog

import (
	"github.com/nattokin/go-backlog/internal/core"
)

type Error = core.Error

type InternalClientError = core.InternalClientError

type APIResponseError = core.APIResponseError

// InvalidOptionKeyError represents an error for an invalid option value.
type InvalidOptionKeyError = core.InvalidOptionKeyError

type ValidationError = core.ValidationError
