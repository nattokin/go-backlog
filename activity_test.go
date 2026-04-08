package backlog

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectActivityService_List(t *testing.T) {
	t.Parallel()

	want := struct {
		spath string
	}{
		spath: "projects/TEST/activities",
	}

	projectKey := "TEST"

	s := newProjectActivityService()
	s.method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}

	_, err := s.List(context.Background(), projectKey)
	assert.Error(t, err)
}

func TestProjectActivityService_List_projectIDOrKeyIsEmpty(t *testing.T) {
	t.Parallel()

	projectKey := ""
	s := newProjectActivityService()
	s.method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		t.Error("s.method.Get must never be called")
		return nil, errors.New("error")
	}

	_, err := s.List(context.Background(), projectKey)
	assert.Error(t, err)
}

func TestProjectActivityService_List_invalidJson(t *testing.T) {
	t.Parallel()

	s := newProjectActivityService()
	s.method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

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

	s := newSpaceActivityService()
	s.method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}

	_, err := s.List(context.Background())
	assert.Error(t, err)
}

func TestUserActivityService_List(t *testing.T) {
	t.Parallel()

	id := 1
	want := struct {
		spath string
	}{
		spath: "users/1/activities",
	}

	s := newUserActivityService()
	s.method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}

	_, err := s.List(context.Background(), id)
	assert.Error(t, err)
}

func TestUserActivityService_List_invalidID(t *testing.T) {
	t.Parallel()

	id := 0
	s := newUserActivityService()
	s.method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		t.Error("s.method.Get must never be called")
		return nil, errors.New("error")
	}

	_, err := s.List(context.Background(), id)
	assert.Error(t, err)
}

func TestBaseActivityService_GetList(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		opts []RequestOption

		want struct {
			activityTypeID []string
			minID          string
			maxID          string
			count          string
			order          string
		}

		wantError bool
	}{
		"no-option": {
			want: struct {
				activityTypeID []string
				minID          string
				maxID          string
				count          string
				order          string
			}{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
			wantError: true,
		},
		"with-activityTypeID": {
			opts: []RequestOption{
				newSpaceActivityOptionService().WithActivityTypeIDs([]int{1, 2}),
			},
			want: struct {
				activityTypeID []string
				minID          string
				maxID          string
				count          string
				order          string
			}{
				activityTypeID: []string{"1", "2"},
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
			wantError: true,
		},
		"with-minID": {
			opts: []RequestOption{
				newSpaceActivityOptionService().WithMinID(10),
			},
			want: struct {
				activityTypeID []string
				minID          string
				maxID          string
				count          string
				order          string
			}{
				activityTypeID: nil,
				minID:          "10",
				maxID:          "",
				count:          "",
				order:          "",
			},
			wantError: true,
		},
		"with-maxID": {
			opts: []RequestOption{
				newSpaceActivityOptionService().WithMaxID(100),
			},
			want: struct {
				activityTypeID []string
				minID          string
				maxID          string
				count          string
				order          string
			}{
				activityTypeID: nil,
				minID:          "",
				maxID:          "100",
				count:          "",
				order:          "",
			},
			wantError: true,
		},
		"with-count": {
			opts: []RequestOption{
				newSpaceActivityOptionService().WithCount(20),
			},
			want: struct {
				activityTypeID []string
				minID          string
				maxID          string
				count          string
				order          string
			}{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "20",
				order:          "",
			},
			wantError: true,
		},
		"with-order": {
			opts: []RequestOption{
				newSpaceActivityOptionService().WithOrder(OrderAsc),
			},
			want: struct {
				activityTypeID []string
				minID          string
				maxID          string
				count          string
				order          string
			}{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "asc",
			},
			wantError: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newSpaceActivityService()
			s.method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, tc.want.activityTypeID, (query)["activityTypeId[]"])
				assert.Equal(t, tc.want.minID, query.Get("minId"))
				assert.Equal(t, tc.want.maxID, query.Get("maxId"))
				assert.Equal(t, tc.want.count, query.Get("count"))
				assert.Equal(t, tc.want.order, query.Get("order"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			}

			if resp, err := s.List(context.Background(), tc.opts...); tc.wantError {
				require.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

// TestActivityService_contextPropagation verifies that the context passed to each
// activity service method is correctly relayed to the underlying method call.
func TestActivityService_contextPropagation(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	cases := []struct {
		name string
		call func(t *testing.T)
	}{
		{"ProjectActivityService.List", func(t *testing.T) {
			s := newProjectActivityService()
			s.method.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s.List(ctx, "TEST") //nolint:errcheck
		}},
		{"SpaceActivityService.List", func(t *testing.T) {
			s := newSpaceActivityService()
			s.method.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s.List(ctx) //nolint:errcheck
		}},
		{"UserActivityService.List", func(t *testing.T) {
			s := newUserActivityService()
			s.method.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s.List(ctx, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t)
		})
	}
}
