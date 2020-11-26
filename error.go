package backlog

import (
	"fmt"
	"strings"
)

// Error represents one of Backlog API response errors.
type Error struct {
	Message  string `json:"message,omitempty"`
	Code     int    `json:"code,omitempty"`
	MoreInfo string `json:"moreInfo,omitempty"`
}

// Error message converted from API error is returned.
func (e *Error) Error() string {
	msg := fmt.Sprint("Massage:", e.Message)
	msg += fmt.Sprint(", Code:", e.Code)

	if len(e.MoreInfo) != 0 {
		msg += fmt.Sprint(", MoreInfo:", e.MoreInfo)
	}

	return msg
}

// APIResponseError represents Error Response of Backlog API.
type APIResponseError struct {
	Errors []*Error `json:"errors,omitempty"`
}

// All error massages converted to APIResponseError is returned.
func (e *APIResponseError) Error() string {
	len := len(e.Errors)
	msgs := make([]string, len)

	for i := 0; i < len; i++ {
		msgs[i] = e.Errors[i].Error()
	}

	return strings.Join(msgs[:], "\n")
}

// InvalidQueryOptionError is an invalid option error.
type InvalidQueryOptionError struct {
	Invalid   queryType
	ValidList []queryType
}

func newInvalidQueryOptionError(invalid queryType, validList []queryType) *InvalidQueryOptionError {
	return &InvalidQueryOptionError{
		Invalid:   invalid,
		ValidList: validList,
	}
}

func (e *InvalidQueryOptionError) Error() string {
	types := make([]string, len(e.ValidList))
	for k, v := range e.ValidList {
		types[k] = v.Value()
	}

	return fmt.Sprintf("invalid option:%s, allowd options:%s", e.Invalid.Value(), strings.Join(types, ","))
}

// InvalidFormOptionError is an invalid option error.
type InvalidFormOptionError struct {
	Invalid   formType
	ValidList []formType
}

func newInvalidFormOptionError(invalid formType, validList []formType) *InvalidFormOptionError {
	return &InvalidFormOptionError{
		Invalid:   invalid,
		ValidList: validList,
	}
}

func (e *InvalidFormOptionError) Error() string {
	types := make([]string, len(e.ValidList))
	for k, v := range e.ValidList {
		types[k] = v.Value()
	}

	return fmt.Sprintf("invalid option:%s, allowd options:%s", e.Invalid.Value(), strings.Join(types, ","))
}
