package backlog

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
	e := &Error{
		Message:  "No project.",
		Code:     6,
		MoreInfo: "more info",
	}
	want := "Message:No project., Code:6, MoreInfo:more info"

	assert.Equal(t, want, e.Error())
}

func TestAPIResponseError_Error(t *testing.T) {
	e := &APIResponseError{
		StatusCode: 404,
		Errors: []*Error{
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
	e := &InvalidOptionKeyError{
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
	e := &InvalidOptionKeyError{
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
	e := core.NewValidationError(msg)
	assert.EqualError(t, e, msg)
}

// ──────────────────────────────────────────────────────────────
//  errors.As assertion tests
// ──────────────────────────────────────────────────────────────

// TestAPIResponseError_errorsAs verifies that APIResponseError returned from
// checkResponse can be unwrapped with errors.As by callers.
func TestAPIResponseError_errorsAs(t *testing.T) {
	resp := &http.Response{
		StatusCode: 404,
		Body:       nil,
	}
	_, err := core.CheckResponse(resp)
	require.Error(t, err)

	wrapped := fmt.Errorf("wrap: %w", err)

	var target *APIResponseError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, 404, target.StatusCode)
}

// TestValidationError_errorsAs verifies that ValidationError can be unwrapped
// with errors.As by callers.
func TestValidationError_errorsAs(t *testing.T) {
	err := core.NewValidationError("invalid argument")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *ValidationError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "invalid argument", target.Error())
}

// TestInvalidOptionKeyError_errorsAs_query verifies that InvalidOptionKeyError
// can be unwrapped with errors.As by callers.
func TestInvalidOptionKeyError_errorsAs_query(t *testing.T) {
	err := core.NewInvalidOptionKeyError(core.ParamActivityTypeIDs.Value(), []core.APIParamOptionType{core.ParamAll, core.ParamArchived})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *InvalidOptionKeyError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, core.ParamActivityTypeIDs.Value(), target.Invalid)
}

// TestInvalidOptionKeyError_errorsAs_form verifies that InvalidOptionKeyError
// can be unwrapped with errors.As by callers.
func TestInvalidOptionKeyError_errorsAs_form(t *testing.T) {
	err := core.NewInvalidOptionKeyError(core.ParamKey.Value(), []core.APIParamOptionType{core.ParamName, core.ParamChartEnabled})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *InvalidOptionKeyError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, core.ParamKey.Value(), target.Invalid)
}

// TestInternalClientError_errorsAs verifies that InternalClientError can be
// unwrapped with errors.As by callers.
func TestInternalClientError_errorsAs(t *testing.T) {
	err := core.NewInternalClientError("missing token")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *InternalClientError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "missing token", target.Error())
}
