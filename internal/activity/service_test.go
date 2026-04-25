package activity_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestProjectActivityService_List(t *testing.T) {
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
	s := activity.NewProjectService(method)

	_, err := s.List(context.Background(), projectKey)
	assert.Error(t, err)
}

func TestProjectActivityService_List_projectIDOrKeyIsEmpty(t *testing.T) {
	t.Parallel()

	projectKey := ""
	method := mock.NewMethod(t)
	s := activity.NewProjectService(method)

	_, err := s.List(context.Background(), projectKey)
	assert.Error(t, err)
}

func TestProjectActivityService_List_invalidJson(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
		}
		return resp, nil
	}
	s := activity.NewProjectService(method)

	projects, err := s.List(context.Background(), "TEST")
	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestSpaceActivityService_List(t *testing.T) {
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
	s := activity.NewSpaceService(method)

	_, err := s.List(context.Background())
	assert.Error(t, err)
}

func TestSpaceActivityService_Get(t *testing.T) {
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
	s := activity.NewSpaceService(method)

	_, err := s.Get(context.Background(), activityID)
	assert.Error(t, err)
}

func TestSpaceActivityService_Get_invalidID(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	s := activity.NewSpaceService(method)

	_, err := s.Get(context.Background(), 0)
	assert.Error(t, err)
}

func TestSpaceActivityService_Get_invalidJson(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
		}
		return resp, nil
	}
	s := activity.NewSpaceService(method)

	got, err := s.Get(context.Background(), 1)
	assert.Nil(t, got)
	assert.Error(t, err)
}

func TestUserActivityService_List(t *testing.T) {
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
	s := activity.NewUserService(method)

	_, err := s.List(context.Background(), id)
	assert.Error(t, err)
}

func TestUserActivityService_List_invalidID(t *testing.T) {
	t.Parallel()

	id := 0
	method := mock.NewMethod(t)
	s := activity.NewUserService(method)

	_, err := s.List(context.Background(), id)
	assert.Error(t, err)
}

func TestBaseActivityService_GetList(t *testing.T) {
	o := &core.OptionService{}
	type want struct {
		activityTypeID []string
		minID          string
		maxID          string
		count          string
		order          string
	}
	cases := map[string]struct {
		opts      []core.RequestOption
		wantError bool
		want      want
	}{
		"success-no-option": {
			opts:      []core.RequestOption{},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"success-withActivityTypeIDs": {
			opts: []core.RequestOption{
				o.WithActivityTypeIDs([]int{1}),
			},
			wantError: false,
			want: want{
				activityTypeID: []string{"1"},
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"success-withMinID": {
			opts: []core.RequestOption{
				o.WithMinID(1),
			},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "1",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"success-withMaxID": {
			opts: []core.RequestOption{
				o.WithMaxID(1),
			},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "1",
				count:          "",
				order:          "",
			},
		},
		"success-withCount": {
			opts: []core.RequestOption{
				o.WithCount(1),
			},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "1",
				order:          "",
			},
		},
		"success-withOrder": {
			opts: []core.RequestOption{
				o.WithOrder(model.OrderAsc),
			},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "asc",
			},
		},
		"success-multiple-options": {
			opts: []core.RequestOption{
				o.WithActivityTypeIDs([]int{1, 2}),
				o.WithMinID(1),
				o.WithMaxID(26),
				o.WithCount(20),
				o.WithOrder(model.OrderAsc),
			},
			wantError: false,
			want: want{
				activityTypeID: []string{"1", "2"},
				minID:          "1",
				maxID:          "26",
				count:          "20",
				order:          "asc",
			},
		},
		"error-option-invalid-value": {
			opts: []core.RequestOption{
				o.WithCount(0),
			},
			wantError: true,
			want:      want{},
		},
		"error-option-invalid-type": {
			opts:      []core.RequestOption{mock.NewInvalidTypeOption()},
			wantError: true,
			want:      want{},
		},
		"error-option-set-failed": {
			opts: []core.RequestOption{
				mock.NewFailingSetOption(core.ParamCount),
			},
			wantError: true,
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, tc.want.activityTypeID, (query)["activityTypeId[]"])
				assert.Equal(t, tc.want.minID, query.Get("minId"))
				assert.Equal(t, tc.want.maxID, query.Get("maxId"))
				assert.Equal(t, tc.want.count, query.Get("count"))
				assert.Equal(t, tc.want.order, query.Get("order"))

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Activity.ListJSON))),
				}
				return resp, nil
			}
			s := activity.NewSpaceService(method)

			if resp, err := s.List(context.Background(), tc.opts...); tc.wantError {
				require.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func Test_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeMockFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"ProjectActivityService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := activity.NewProjectService(m)
			s.List(ctx, "TEST") //nolint:errcheck
		}},
		{"SpaceActivityService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := activity.NewSpaceService(m)
			s.List(ctx) //nolint:errcheck
		}},
		{"SpaceActivityService.Get", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := activity.NewSpaceService(m)
			s.Get(ctx, 1) //nolint:errcheck
		}},
		{"UserActivityService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := activity.NewUserService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
