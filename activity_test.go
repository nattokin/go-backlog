package backlog_test

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestProjectActivityService_List(t *testing.T) {
	projectKey := "TEST"
	want := struct {
		spath string
	}{
		spath: "projects/" + projectKey + "/activities",
	}
	s := &backlog.ProjectActivityService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	})
	s.List(backlog.ProjectKey(projectKey))
}

func TestProjectActivityService_List_projectIDOrKeyIsEmpty(t *testing.T) {
	projectKey := ""
	s := &backlog.ProjectActivityService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			t.Error("s.method.Get must never be called")
			return nil, errors.New("error")
		},
	})
	s.List(backlog.ProjectKey(projectKey))
}

func TestProjectActivityService_List_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.ProjectActivityService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	projects, err := s.List(backlog.ProjectKey("TEST"))
	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestSpaceActivityService_List(t *testing.T) {
	want := struct {
		spath string
	}{
		spath: "space/activities",
	}
	s := &backlog.SpaceActivityService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	})
	s.List()
}

func TestUserActivityService_List(t *testing.T) {
	id := 1234
	want := struct {
		spath string
	}{
		spath: "users/" + strconv.Itoa(id) + "/activities",
	}
	s := &backlog.UserActivityService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	})
	s.List(id)
}

func TestUserActivityService_List_invaliedID(t *testing.T) {
	id := 0
	s := &backlog.UserActivityService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			t.Error("s.method.Get must never be called")
			return nil, errors.New("error")
		},
	})
	s.List(id)
}

func TestBaseActivityService_GetList(t *testing.T) {
	o := &backlog.ActivityOptionService{}
	type want struct {
		activityTypeID []string
		minID          string
		maxID          string
		count          string
		order          string
	}
	cases := map[string]struct {
		options   []backlog.ActivityOption
		wantError bool
		want      want
	}{
		"NoOption": {
			options:   []backlog.ActivityOption{},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"WithActivityTypeIDs_valid": {
			options: []backlog.ActivityOption{
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
		"WithActivityTypeIDs_invalid": {
			options: []backlog.ActivityOption{
				o.WithActivityTypeIDs([]int{0}),
			},
			wantError: true,
			want:      want{},
		},
		"WithMinID_valid": {
			options: []backlog.ActivityOption{
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
		"WithMinID_invalid": {
			options: []backlog.ActivityOption{
				o.WithMinID(0),
			},
			wantError: true,
			want:      want{},
		},
		"WithMaxID_valid": {
			options: []backlog.ActivityOption{
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
		"WithMaxID_invalid": {
			options: []backlog.ActivityOption{
				o.WithMaxID(0),
			},
			wantError: true,
			want:      want{},
		},
		"WithCount_valid": {
			options: []backlog.ActivityOption{
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
		"WithCount_invalid": {
			options: []backlog.ActivityOption{
				o.WithCount(0),
			},
			wantError: true,
			want:      want{},
		},
		"WithOrder_valid": {
			options: []backlog.ActivityOption{
				o.WithOrder(backlog.OrderAsc),
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
		"WithOrder_invalid": {
			options: []backlog.ActivityOption{
				o.WithOrder("test"),
			},
			wantError: true,
			want:      want{},
		},
		"MultipleOptions": {
			options: []backlog.ActivityOption{
				o.WithActivityTypeIDs([]int{1, 2}),
				o.WithMinID(1),
				o.WithMaxID(100),
				o.WithCount(20),
				o.WithOrder(backlog.OrderAsc),
			},
			wantError: false,
			want: want{
				activityTypeID: []string{"1", "2"},
				minID:          "1",
				maxID:          "100",
				count:          "20",
				order:          "asc",
			},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/activity.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.SpaceActivityService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					v := *params.ExportURLValues()
					assert.Equal(t, tc.want.activityTypeID, v["activityTypeId[]"])
					assert.Equal(t, tc.want.minID, params.Get("minId"))
					assert.Equal(t, tc.want.maxID, params.Get("maxId"))
					assert.Equal(t, tc.want.count, params.Get("count"))
					assert.Equal(t, tc.want.order, params.Get("order"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			})

			if _, err := s.List(tc.options...); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
