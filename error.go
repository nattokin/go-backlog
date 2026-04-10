package backlog

import (
	"fmt"
	"strings"

	"github.com/nattokin/go-backlog/internal/core"
)

type Error = core.Error

type InternalClientError = core.InternalClientError

type APIResponseError = core.APIResponseError

// InvalidOptionKeyError represents an error for an invalid option value.
type InvalidOptionKeyError struct {
	Invalid   string
	ValidList []string
}

func newInvalidOptionKeyError(invalid string, validList []apiParamOptionType) *InvalidOptionKeyError {
	validKeys := []string{}
	for _, v := range validList {
		validKeys = append(validKeys, v.Value())
	}

	return &InvalidOptionKeyError{
		Invalid:   invalid,
		ValidList: validKeys,
	}
}

func (e *InvalidOptionKeyError) Error() string {
	return fmt.Sprintf("invalid option key:%s, allowed option keys:%s", e.Invalid, strings.Join(e.ValidList, ","))
}

// ValidationError represents an argument validation error.
type ValidationError struct {
	message string
}

func newValidationError(msg string) *ValidationError {
	return &ValidationError{
		message: msg,
	}
}

func (e *ValidationError) Error() string {
	return e.message
}
