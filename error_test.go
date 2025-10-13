package backlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestInvalidFormOptionError_Error(t *testing.T) {
	e := &InvalidFormOptionError{
		Invalid: formKey,
		ValidList: []formType{
			formName,
			formKey,
			formChartEnabled,
		},
	}
	assert.EqualError(t, e, "invalid option:key, allowed options:name,key,chartEnabled")
}

func TestInvalidQueryOptionError_Error(t *testing.T) {
	e := &InvalidQueryOptionError{
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
		Message: msg,
	}
	assert.EqualError(t, e, msg)
}
