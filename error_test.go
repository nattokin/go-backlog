package backlog

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestInvalidOptionError_Error_form(t *testing.T) {
	e := &InvalidOptionError[formType]{
		Invalid: formKey,
		ValidList: []formType{
			formName,
			formKey,
			formChartEnabled,
		},
	}
	assert.EqualError(t, e, "invalid option:key, allowed options:name,key,chartEnabled")
}

func TestInvalidOptionError_Error_query(t *testing.T) {
	e := &InvalidOptionError[queryType]{
		Invalid: queryActivityTypeIDs,
		ValidList: []queryType{
			queryAll,
			queryArchived,
			queryOrder,
		},
	}
	assert.EqualError(t, e, "invalid option:activityTypeId[], allowed options:all,archived,order")
}

func TestValidationError_Error(t *testing.T) {
	msg := "validation error"
	e := &ValidationError{
		message: msg,
	}
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
	_, err := checkResponse(resp)
	require.Error(t, err)

	wrapped := fmt.Errorf("wrap: %w", err)

	var target *APIResponseError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, 404, target.StatusCode)
}

// TestValidationError_errorsAs verifies that ValidationError can be unwrapped
// with errors.As by callers.
func TestValidationError_errorsAs(t *testing.T) {
	err := newValidationError("invalid argument")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *ValidationError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "invalid argument", target.Error())
}

// TestInvalidOptionError_errorsAs_query verifies that InvalidOptionError[queryType]
// can be unwrapped with errors.As by callers.
func TestInvalidOptionError_errorsAs_query(t *testing.T) {
	err := newInvalidOptionError(queryActivityTypeIDs, []queryType{queryAll, queryArchived})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *InvalidOptionError[queryType]
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, queryActivityTypeIDs, target.Invalid)
}

// TestInvalidOptionError_errorsAs_form verifies that InvalidOptionError[formType]
// can be unwrapped with errors.As by callers.
func TestInvalidOptionError_errorsAs_form(t *testing.T) {
	err := newInvalidOptionError(formKey, []formType{formName, formChartEnabled})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *InvalidOptionError[formType]
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, formKey, target.Invalid)
}

// TestInternalClientError_errorsAs verifies that InternalClientError can be
// unwrapped with errors.As by callers.
func TestInternalClientError_errorsAs(t *testing.T) {
	err := newInternalClientError("missing token")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *InternalClientError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "missing token", target.Error())
}
