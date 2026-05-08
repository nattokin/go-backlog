package space_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/space"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestActivityService_List(t *testing.T) {
	t.Parallel()

	want := struct {
		spath string
	}{
		spath: "space/activities",
	}

	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}
	s := space.NewActivityService(method)

	_, err := s.List(context.Background())
	assert.Error(t, err)
}

func TestActivityService_Get(t *testing.T) {
	t.Parallel()

	activityID := 3153

	want := struct {
		spath string
	}{
		spath: "activities/" + strconv.Itoa(activityID),
	}

	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.Nil(t, query)
		return nil, errors.New("error")
	}
	s := space.NewActivityService(method)

	_, err := s.Get(context.Background(), activityID)
	assert.Error(t, err)
}

func TestActivityService_Get_invalidID(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	s := space.NewActivityService(method)

	_, err := s.Get(context.Background(), 0)
	assert.Error(t, err)
}

func TestActivityService_Get_invalidJson(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		return mock.NewJSONResponse(fixture.InvalidJSON), nil
	}
	s := space.NewActivityService(method)

	got, err := s.Get(context.Background(), 1)
	assert.Nil(t, got)
	assert.Error(t, err)
}
