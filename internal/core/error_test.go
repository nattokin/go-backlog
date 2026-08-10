package core_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestError_Error(t *testing.T) {
	e := &core.Error{
		Message:  "No project.",
		Code:     6,
		MoreInfo: "more info",
	}
	want := "Message:No project., Code:6, MoreInfo:more info"

	assert.Equal(t, want, e.Error())
}

func TestAPIResponseError_Error(t *testing.T) {
	e := &core.APIResponseError{
		StatusCode: 404,
		Errors: []*core.Error{
			{
				Message:  "1st error",
				Code:     5,
				MoreInfo: "more info 1",
			},
			{
				Message:  "2nd error",
				Code:     9,
				MoreInfo: "more info 2",
			},
		},
	}
	want := "Status Code:404\nMessage:1st error, Code:5, MoreInfo:more info 1\nMessage:2nd error, Code:9, MoreInfo:more info 2"

	assert.Equal(t, want, e.Error())
}

func TestInvalidOptionKeyError_Error_form(t *testing.T) {
	e := &core.InvalidOptionKeyError{
		Invalid: core.ParamKey.Value(),
		ValidList: []string{
			core.ParamName.Value(),
			core.ParamKey.Value(),
			core.ParamChartEnabled.Value(),
		},
	}
	assert.EqualError(t, e, "invalid option key:key, allowed option keys:name,key,chartEnabled")
}

func TestInvalidOptionKeyError_Error_query(t *testing.T) {
	e := &core.InvalidOptionKeyError{
		Invalid: core.ParamActivityTypeIDs.Value(),
		ValidList: []string{
			core.ParamAll.Value(),
			core.ParamArchived.Value(),
			core.ParamOrder.Value(),
		},
	}
	assert.EqualError(t, e, "invalid option key:activityTypeId[], allowed option keys:all,archived,order")
}

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

// ──────────────────────────────────────────────────────────────
//  errors.As assertion tests
// ──────────────────────────────────────────────────────────────

func TestAPIResponseError_errorsAs(t *testing.T) {
	resp := &http.Response{
		StatusCode: 404,
		Body:       nil,
	}
	_, err := core.CheckResponse(resp)
	require.Error(t, err)

	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.APIResponseError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, 404, target.StatusCode)
}

func TestValidationError_errorsAs(t *testing.T) {
	err := core.NewValidationError("key", "invalid argument")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.ValidationError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "invalid argument", target.Error())
	assert.Equal(t, "key", target.Target())
}

func TestInvalidOptionKeyError_errorsAs_query(t *testing.T) {
	err := core.NewInvalidOptionKeyError(core.ParamActivityTypeIDs.Value(), []core.APIParamOptionType{core.ParamAll, core.ParamArchived})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.InvalidOptionKeyError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, core.ParamActivityTypeIDs.Value(), target.Invalid)
}

func TestInvalidOptionKeyError_errorsAs_form(t *testing.T) {
	err := core.NewInvalidOptionKeyError(core.ParamKey.Value(), []core.APIParamOptionType{core.ParamName, core.ParamChartEnabled})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.InvalidOptionKeyError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, core.ParamKey.Value(), target.Invalid)
}

func TestInternalClientError_errorsAs(t *testing.T) {
	err := core.NewInternalClientError("missing token")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.InternalClientError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "missing token", target.Error())
}

func TestInvalidOptionError_errorsAs(t *testing.T) {
	err := core.NewInvalidOptionError("nil option is not allowed")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.InvalidOptionError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "nil option is not allowed", target.Error())
}
