package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/user"
)

func TestActivityService_List(t *testing.T) {
	t.Parallel()

	id := 1234

	want := struct {
		spath string
	}{
		spath: "users/" + strconv.Itoa(id) + "/activities",
	}

	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}
	s := user.NewActivityService(method)

	_, err := s.List(context.Background(), id)
	assert.Error(t, err)
}

func TestActivityService_List_invalidID(t *testing.T) {
	t.Parallel()

	id := 0
	method := mock.NewMethod(t)
	s := user.NewActivityService(method)

	_, err := s.List(context.Background(), id)
	assert.Error(t, err)
}
