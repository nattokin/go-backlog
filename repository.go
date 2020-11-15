package backlog

import (
	"errors"
	"strconv"
)

// RepositoryIDOrKeyGetter has method to get RepositoryIDOrKey and validation errror.
type RepositoryIDOrKeyGetter interface {
	getRepositoryIDOrKey() (string, error)
}

// RepositoryID implements RepositoryIDOrKeyGetter interface.
type RepositoryID int

// RepositoryName implements RepositoryIDOrKeyGetter interface.
type RepositoryName string

func (i RepositoryID) getRepositoryIDOrKey() (string, error) {
	if i <= 0 {
		return "", errors.New("id must be greater than 0")
	}
	return strconv.Itoa(int(i)), nil
}

func (k RepositoryName) getRepositoryIDOrKey() (string, error) {
	if k == "" {
		return "", errors.New("key must not be empty")
	}
	return string(k), nil
}

// RepositoryService has methods for Repository.
type RepositoryService struct {
	method *method
}
