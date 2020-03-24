package backlog

import (
	"fmt"
	"strings"
)

// APIResponseError represents Error Response of Backlog API.
type APIResponseError struct {
	Errors []*Error `json:"errors,omitempty"`
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

// All error massages converted to APIResponseError is returned.
func (e *APIResponseError) Error() string {
	len := len(e.Errors)
	msgs := make([]string, len)

	for i := 0; i < len; i++ {
		msgs[i] = e.Errors[i].Error()
	}

	return strings.Join(msgs[:], "\n")
}
