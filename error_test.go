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
	want := "Massage:No project., Code:6, MoreInfo:more info"

	assert.Equal(t, want, e.Error())
}

func TestAPIResponseError_Error(t *testing.T) {
	e := &backlog.APIResponseError{
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
	want := "Massage:1st error, Code:5, MoreInfo:more info 1\nMassage:2nd error, Code:9, MoreInfo:more info 2"

	assert.Equal(t, want, e.Error())
}

func TestInvalidOptionError_Error(t *testing.T) {
	e := &backlog.InvalidOptionError{
		Invalid: backlog.ExportOptionActivityTypeIDs,
		ValidList: []backlog.ExportOptionType{
			backlog.ExportOptionActivityTypeIDs,
			backlog.ExportOptionArchived,
			backlog.ExportOptionChartEnabled,
		},
	}
	assert.EqualError(t, e, "invalid option error. option:ActivityTypeIDs, allowd options:ActivityTypeIDs,Archived,ChartEnabled")
}
