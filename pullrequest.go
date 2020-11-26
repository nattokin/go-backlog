package backlog

import (
	"errors"
	"strconv"
)

// PRNumber implements IssueIDOrKeyGetter interface.
type PRNumber int

func (n PRNumber) validate() error {
	if n < 1 {
		return errors.New("prNumber must not be less than 1")
	}
	return nil
}

func (n PRNumber) String() string {
	return strconv.Itoa(int(n))
}

// PullRequestService has methods for Issue.
type PullRequestService struct {
	method *method

	Attachment *PullRequestAttachmentService
}
