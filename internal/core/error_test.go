package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestValidationError_Error(t *testing.T) {
	msg := "validation error"
	e := core.NewValidationError("someTarget", msg)
	assert.EqualError(t, e, msg)
}

func TestValidationError_Fields(t *testing.T) {
	e := core.NewValidationError("offset", "offset must not be negative")
	assert.Equal(t, "offset", e.Target())
	assert.Equal(t, "offset must not be negative", e.Message())
}

func TestValidationErrors_Error_single(t *testing.T) {
	t.Parallel()
	ves := core.ValidationErrors{
		core.NewValidationError("count", "count must be greater than 0"),
	}
	assert.Equal(t, "count must be greater than 0", ves.Error())
}

func TestValidationErrors_Error_multiple(t *testing.T) {
	t.Parallel()
	ves := core.ValidationErrors{
		core.NewValidationError("count", "count must be greater than 0"),
		core.NewValidationError("order", "order must be asc or desc"),
	}
	assert.Equal(t, "count must be greater than 0\norder must be asc or desc", ves.Error())
}

func TestValidationError_errorsAs(t *testing.T) {
	err := core.NewValidationError("key", "invalid argument")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.ValidationError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "invalid argument", target.Error())
	assert.Equal(t, "key", target.Target())
}
