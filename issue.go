package backlog

import (
	"errors"
	"strconv"
)

// IssueIDOrKeyGetter has method to get IssueIDOrKey and validation errror.
type IssueIDOrKeyGetter interface {
	getIssueIDOrKey() (string, error)
}

// IssueID implements IssueIDOrKeyGetter interface.
type IssueID int

func (i IssueID) validate() error {
	if i < 1 {
		return errors.New("issueID must not be less than 1")
	}
	return nil
}

func (i IssueID) String() string {
	return strconv.Itoa(int(i))
}

func (i IssueID) getIssueIDOrKey() (string, error) {
	if err := i.validate(); err != nil {
		return "", err
	}
	return i.String(), nil
}

// IssueKey implements IssueIDOrKeyGetter interface.
type IssueKey string

func (k IssueKey) getIssueIDOrKey() (string, error) {
	if k == "" {
		return "", errors.New("key must not be empty")
	}
	return string(k), nil
}

// IssueService has methods for Issue.
type IssueService struct {
	method *method

	Attachment *IssueAttachmentService
}
