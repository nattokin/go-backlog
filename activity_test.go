package backlog

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectActivityService_List(t *testing.T) {
	t.Parallel()

	projectKey := "TEST"

	want := struct {
		spath string
	}{
		spath: "projects/" + projectKey + "/activities",
	}

	s := newProjectActivityService()
	s.method.Get = func(spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}

	_, err := s.List(projectKey)
	assert.Error(t, err)
}

func TestProjectActivityService_List_projectIDOrKeyIsEmpty(t *testing.T) {
	t.Parallel()

	projectKey := ""
	s := newProjectActivityService()
	s.method.Get = func(spath string, query url.Values) (*http.Response, error) {
		t.Error("s.method.Get must never be called")
		return nil, errors.New("error")
	}

	_, err := s.List(projectKey)
	assert.Error(t, err)
}

func TestProjectActivityService_List_invalidJson(t *testing.T) {
	t.Parallel()

	s := newProjectActivityService()
	s.method.Get = func(spath string, query url.Values) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	projects, err := s.List("TEST")
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
	s.method.Get = func(spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}

	_, err := s.List()
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

	s := newUserActivityService()
	s.method.Get = func(spath string, query url.Values) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		return nil, errors.New("error")
	}

	_, err := s.List(id)
	assert.Error(t, err)
}

func TestUserActivityService_List_invalidID(t *testing.T) {
	t.Parallel()

	id := 0
	s := newUserActivityService()
	s.method.Get = func(spath string, query url.Values) (*http.Response, error) {
		t.Error("s.method.Get must never be called")
		return nil, errors.New("error")
	}

	_, err := s.List(id)
	assert.Error(t, err)
}

func TestBaseActivityService_GetList(t *testing.T) {
	o := newActivityOptionService()
	type want struct {
		activityTypeID []string
		minID          string
		maxID          string
		count          string
		order          string
	}
	cases := map[string]struct {
		opts      []RequestOption
		wantError bool
		want      want
	}{
		"success-no-option": {
			opts:      []RequestOption{},
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
			opts: []RequestOption{
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
			opts: []RequestOption{
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
			opts: []RequestOption{
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
			opts: []RequestOption{
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
			opts: []RequestOption{
				o.WithOrder(OrderAsc),
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
			opts: []RequestOption{
				o.WithActivityTypeIDs([]int{1, 2}),
				o.WithMinID(1),
				o.WithMaxID(26),
				o.WithCount(20),
				o.WithOrder(OrderAsc),
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
			opts: []RequestOption{
				o.WithCount(0),
			},
			wantError: true,
			want:      want{},
		},
		"error-option-invalid-type": {
			opts:      []RequestOption{newInvalidTypeOption()},
			wantError: true,
			want:      want{},
		},
		"error-option-set-failed": {
			opts: []RequestOption{
				newFailingSetOption(paramCount),
			},
			wantError: true,
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := newSpaceActivityService()
			s.method.Get = func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, tc.want.activityTypeID, (query)["activityTypeId[]"])
				assert.Equal(t, tc.want.minID, query.Get("minId"))
				assert.Equal(t, tc.want.maxID, query.Get("maxId"))
				assert.Equal(t, tc.want.count, query.Get("count"))
				assert.Equal(t, tc.want.order, query.Get("order"))

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataActivityListJSON))),
				}
				return resp, nil
			}

			if resp, err := s.List(tc.opts...); tc.wantError {
				require.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
