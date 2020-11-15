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

// IssueKey implements IssueIDOrKeyGetter interface.
type IssueKey string

func (i IssueID) getIssueIDOrKey() (string, error) {
	if i <= 0 {
		return "", errors.New("id must be greater than 0")
	}
	return strconv.Itoa(int(i)), nil
}

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
