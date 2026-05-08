package project_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestActivityService_List(t *testing.T) {
	t.Parallel()

	projectKey := "TEST"

	want := struct {
		spath string
	}{
		spath: "projects/" + projectKey + "/activities",
	}

	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}
	s := project.NewActivityService(method)

	_, err := s.List(context.Background(), projectKey)
	assert.Error(t, err)
}

func TestActivityService_List_projectIDOrKeyIsEmpty(t *testing.T) {
	t.Parallel()

	projectKey := ""
	method := mock.NewMethod(t)
	s := project.NewActivityService(method)

	_, err := s.List(context.Background(), projectKey)
	assert.Error(t, err)
}

func TestActivityService_List_invalidJson(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		return mock.NewJSONResponse(fixture.InvalidJSON), nil
	}
	s := project.NewActivityService(method)

	projects, err := s.List(context.Background(), "TEST")
	assert.Nil(t, projects)
	assert.Error(t, err)
}
