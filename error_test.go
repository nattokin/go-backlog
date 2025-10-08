package backlog_test

import (
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	e := &backlog.Error{
		Message:  "No project.",
		Code:     6,
		MoreInfo: "more info",
	}
	want := "Message:No project., Code:6, MoreInfo:more info"

	assert.Equal(t, want, e.Error())
}

func TestAPIResponseError_Error(t *testing.T) {
	e := &backlog.APIResponseError{
		StatusCode: 404,
		Errors: []*backlog.Error{
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
	e := &backlog.InvalidFormOptionError{
		Invalid: backlog.ExportFormKey,
		ValidList: []backlog.ExportFormType{
			backlog.ExportFormName,
			backlog.ExportFormKey,
			backlog.ExportFormChartEnabled,
		},
	}
	assert.EqualError(t, e, "invalid option:key, allowed options:name,key,chartEnabled")
}

func TestInvalidQueryOptionError_Error(t *testing.T) {
	e := &backlog.InvalidQueryOptionError{
		Invalid: backlog.ExportQueryActivityTypeIDs,
		ValidList: []backlog.ExportQueryType{
			backlog.ExportQueryAll,
			backlog.ExportQueryArchived,
			backlog.ExportQueryOrder,
		},
	}
	assert.EqualError(t, e, "invalid option:activityTypeId[], allowed options:all,archived,order")
}

func TestValidationError_Error(t *testing.T) {
	msg := "validation error"
	e := &backlog.ValidationError{
		Message: msg,
	}
	assert.EqualError(t, e, msg)
}
