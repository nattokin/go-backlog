package backlog_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/nattokin/go-backlog"
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

	s := backlog.ExportNewProjectActivityService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	})

	s.List(projectKey)
}

func TestProjectActivityService_List_projectIDOrKeyIsEmpty(t *testing.T) {
	t.Parallel()

	projectKey := ""
	s := backlog.ExportNewProjectActivityService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			t.Error("s.method.Get must never be called")
			return nil, errors.New("error")
		},
	})

	s.List(projectKey)
}

func TestProjectActivityService_List_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewProjectActivityService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

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

	s := backlog.ExportNewSpaceActivityService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	})

	s.List()
}

func TestUserActivityService_List(t *testing.T) {
	t.Parallel()

	id := 1234

	want := struct {
		spath string
	}{
		spath: "users/" + strconv.Itoa(id) + "/activities",
	}

	s := backlog.ExportNewUserActivityService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	})

	s.List(id)
}

func TestUserActivityService_List_invalidID(t *testing.T) {
	t.Parallel()

	id := 0
	s := backlog.ExportNewUserActivityService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			t.Error("s.method.Get must never be called")
			return nil, errors.New("error")
		},
	})

	s.List(id)
}

func TestBaseActivityService_GetList(t *testing.T) {
	o := backlog.ExportNewActivityOptionService()
	type want struct {
		activityTypeID []string
		minID          string
		maxID          string
		count          string
		order          string
	}
	cases := map[string]struct {
		options   []*backlog.QueryOption
		wantError bool
		want      want
	}{
		"NoOption": {
			options:   []*backlog.QueryOption{},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"WithActivityTypeIDs": {
			options: []*backlog.QueryOption{
				o.WithQueryActivityTypeIDs([]int{1}),
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
		"WithMinID": {
			options: []*backlog.QueryOption{
				o.WithQueryMinID(1),
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
		"WithMaxID": {
			options: []*backlog.QueryOption{
				o.WithQueryMaxID(1),
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
		"WithCount": {
			options: []*backlog.QueryOption{
				o.WithQueryCount(1),
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
		"WithOrder": {
			options: []*backlog.QueryOption{
				o.WithQueryOrder(backlog.OrderAsc),
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
		"MultipleOptions": {
			options: []*backlog.QueryOption{
				o.WithQueryActivityTypeIDs([]int{1, 2}),
				o.WithQueryMinID(1),
				o.WithQueryMaxID(26),
				o.WithQueryCount(20),
				o.WithQueryOrder(backlog.OrderAsc),
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
		"OptionError": {
			options: []*backlog.QueryOption{
				o.WithQueryCount(0),
			},
			wantError: true,
			want:      want{},
		},
		"InvalidOption": {
			options: []*backlog.QueryOption{
				backlog.ExportNewQueryOption(0, nil, func(p *backlog.QueryParams) error {
					return nil
				}),
			},
			wantError: true,
			want:      want{},
		},
		"SetErr": {
			options: []*backlog.QueryOption{
				backlog.ExportNewQueryOption(
					backlog.ExportQueryCount,
					nil,
					func(p *backlog.QueryParams) error {
						return errors.New("set error")
					}),
			},
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := backlog.ExportNewSpaceActivityService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					assert.Equal(t, tc.want.activityTypeID, (*query.Values)["activityTypeId[]"])
					assert.Equal(t, tc.want.minID, query.Get("minId"))
					assert.Equal(t, tc.want.maxID, query.Get("maxId"))
					assert.Equal(t, tc.want.count, query.Get("count"))
					assert.Equal(t, tc.want.order, query.Get("order"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataActivityListJSON))),
					}
					return resp, nil
				},
			})

			if resp, err := s.List(tc.options...); tc.wantError {
				require.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}
