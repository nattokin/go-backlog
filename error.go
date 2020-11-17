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

// InvalidOptionError is an invalid option error.
type InvalidOptionError struct {
	Invalid   optionType
	ValidList []optionType
}

func newInvalidOptionError(invalid optionType, validList []optionType) *InvalidOptionError {
	return &InvalidOptionError{
		Invalid:   invalid,
		ValidList: validList,
	}
}

func (e InvalidOptionError) validListString() string {
	var types []string
	for _, v := range e.ValidList {
		types = append(types, v.String())
	}
	return strings.Join(types, ",")
}

func (e *InvalidOptionError) Error() string {
	return fmt.Sprintf("invalid option error. option:%s, allowd options:%s", e.Invalid, e.validListString())
}
