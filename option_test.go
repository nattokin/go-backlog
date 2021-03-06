package backlog_test

import (
	"strconv"
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestQueryOptionService_WithActivityTypeIDs(t *testing.T) {
	o := backlog.QueryOptionService{}

	cases := map[string]struct {
		typeIDs   []int
		want      []string
		wantError bool
	}{
		"Valid-1": {
			typeIDs:   []int{1},
			want:      []string{"1"},
			wantError: false,
		},
		"Valid-2": {
			typeIDs:   []int{26},
			want:      []string{"26"},
			wantError: false,
		},
		"Valid-3": {
			typeIDs: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
			want: []string{
				"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13",
				"14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26",
			},
			wantError: false,
		},
		"Invalid-1": {
			typeIDs:   []int{0},
			want:      nil,
			wantError: true,
		},
		"Invalid-2": {
			typeIDs:   []int{-1},
			want:      nil,
			wantError: true,
		},
		"Invalid-3": {
			typeIDs:   []int{27},
			want:      nil,
			wantError: true,
		},
		"Invalid-4": {
			typeIDs:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27},
			want:      nil,
			wantError: true,
		},
		"Invalid-5": {
			typeIDs:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
			want:      nil,
			wantError: true,
		},
		"Empty": {
			typeIDs:   []int{},
			want:      nil,
			wantError: false,
		},
		"duplicate": {
			typeIDs:   []int{1, 1},
			want:      []string{"1", "1"},
			wantError: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithActivityTypeIDs(tc.typeIDs)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryParamsWithOptions(query, []*backlog.QueryOption{option}, backlog.ExportQueryActivityTypeIDs); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, (*query.Values)["activityTypeId[]"])
			}
		})
	}
}

func TestQueryOptionService_WithActivityTypeIDs_invalidOption(t *testing.T) {
	o := backlog.QueryOptionService{}

	activityIDs := []int{1, 2}
	validTypes := []backlog.ExportQueryType{
		backlog.ExportQueryAll,
		backlog.ExportQueryArchived,
		backlog.ExportQueryCount,
		backlog.ExportQueryKey,
		backlog.ExportQueryOrder,
	}

	option := o.WithActivityTypeIDs(activityIDs)
	query := backlog.NewQueryParams()
	err := backlog.ExportQueryParamsWithOptions(query, []*backlog.QueryOption{option}, validTypes...)
	assert.IsType(t, &backlog.InvalidQueryOptionError{}, err)
}

func TestQueryOptionService_WithAll(t *testing.T) {
	o := backlog.QueryOptionService{}

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithAll(tc.enabeld)
			query := backlog.NewQueryParams()
			err := backlog.ExportQueryParamsWithOptions(query, []*backlog.QueryOption{option}, backlog.ExportQueryAll)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), query.Get("all"))
		})
	}
}

func TestQueryOptionService_WithArchived(t *testing.T) {
	o := backlog.QueryOptionService{}

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithArchived(tc.enabeld)
			query := backlog.NewQueryParams()
			err := backlog.ExportQueryParamsWithOptions(query, []*backlog.QueryOption{option}, backlog.ExportQueryArchived)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), query.Get("archived"))
		})
	}
}

func TestQueryOptionService_WithCount(t *testing.T) {
	o := backlog.QueryOptionService{}

	cases := map[string]struct {
		count     int
		wantError bool
	}{
		"Valid-1": {
			count:     1,
			wantError: false,
		},
		"Valid-2": {
			count:     100,
			wantError: false,
		},
		"Invalid-1": {
			count:     0,
			wantError: true,
		},
		"Invalid-2": {
			count:     -1,
			wantError: true,
		},
		"Invalid-3": {
			count:     101,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithCount(tc.count)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryParamsWithOptions(query, []*backlog.QueryOption{option}, backlog.ExportQueryCount); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.count), query.Get("count"))
			}
		})
	}
}

func TestQueryOptionService_WithKeyword(t *testing.T) {
	o := backlog.QueryOptionService{}

	cases := map[string]struct {
		keyword string
	}{
		"Valid": {
			keyword: "test",
		},
		"Empty": {
			keyword: "",
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithKeyword(tc.keyword)
			query := backlog.NewQueryParams()

			err := backlog.ExportQueryParamsWithOptions(query, []*backlog.QueryOption{option}, backlog.ExportQueryKeyword)
			assert.NoError(t, err)
			assert.Equal(t, tc.keyword, query.Get("keyword"))
		})
	}
}

func TestQueryOptionService_WithMinID(t *testing.T) {
	o := backlog.QueryOptionService{}

	cases := map[string]struct {
		minID     int
		wantError bool
	}{
		"Valid": {
			minID:     1,
			wantError: false,
		},
		"Invalid-1": {
			minID:     0,
			wantError: true,
		},
		"Invalid-2": {
			minID:     -1,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithMinID(tc.minID)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryOptionSet(option, query); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.minID), query.Get("minId"))
			}
		})
	}
}

func TestQueryOptionService_WithMaxID(t *testing.T) {
	o := backlog.QueryOptionService{}

	cases := map[string]struct {
		maxID     int
		wantError bool
	}{
		"Valid": {
			maxID:     1,
			wantError: false,
		},
		"Invalid-1": {
			maxID:     0,
			wantError: true,
		},
		"Invalid-2": {
			maxID:     -1,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithMaxID(tc.maxID)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryOptionSet(option, query); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.maxID), query.Get("maxId"))
			}
		})
	}
}

func TestQueryOptionService_WithOrder(t *testing.T) {
	o := backlog.QueryOptionService{}

	cases := map[string]struct {
		order     backlog.Order
		wantError bool
	}{
		"asc": {
			order:     backlog.OrderAsc,
			wantError: false,
		},
		"desc": {
			order:     backlog.OrderDesc,
			wantError: false,
		},
		"Invalid": {
			order:     "test",
			wantError: true,
		},
		"Empty": {
			order:     "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithOrder(tc.order)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryParamsWithOptions(query, []*backlog.QueryOption{option}, backlog.ExportQueryOrder); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(tc.order), query.Get("order"))
			}
		})
	}
}

func TestActivityOptionService_WithActivityTypeIDs(t *testing.T) {
	o := backlog.ActivityOptionService{}

	cases := map[string]struct {
		typeIDs   []int
		want      []string
		wantError bool
	}{
		"Valid-1": {
			typeIDs:   []int{1},
			want:      []string{"1"},
			wantError: false,
		},
		"Valid-2": {
			typeIDs:   []int{26},
			want:      []string{"26"},
			wantError: false,
		},
		"Valid-3": {
			typeIDs: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
			want: []string{
				"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13",
				"14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26",
			},
			wantError: false,
		},
		"Invalid-1": {
			typeIDs:   []int{0},
			want:      nil,
			wantError: true,
		},
		"Invalid-2": {
			typeIDs:   []int{-1},
			want:      nil,
			wantError: true,
		},
		"Invalid-3": {
			typeIDs:   []int{27},
			want:      nil,
			wantError: true,
		},
		"Invalid-4": {
			typeIDs:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27},
			want:      nil,
			wantError: true,
		},
		"Invalid-5": {
			typeIDs:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
			want:      nil,
			wantError: true,
		},
		"Empty": {
			typeIDs:   []int{},
			want:      nil,
			wantError: false,
		},
		"duplicate": {
			typeIDs:   []int{1, 1},
			want:      []string{"1", "1"},
			wantError: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithQueryActivityTypeIDs(tc.typeIDs)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryOptionSet(option, query); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				v := *query.Values
				assert.Equal(t, tc.want, v["activityTypeId[]"])
			}
		})
	}
}

func TestActivityOptionService_WithMinID(t *testing.T) {
	o := backlog.ActivityOptionService{}

	cases := map[string]struct {
		minID     int
		wantError bool
	}{
		"Valid": {
			minID:     1,
			wantError: false,
		},
		"Invalid-1": {
			minID:     0,
			wantError: true,
		},
		"Invalid-2": {
			minID:     -1,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithQueryMinID(tc.minID)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryOptionSet(option, query); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.minID), query.Get("minId"))
			}
		})
	}
}

func TestActivityOptionService_WithMaxID(t *testing.T) {
	o := backlog.ActivityOptionService{}

	cases := map[string]struct {
		maxID     int
		wantError bool
	}{
		"Valid": {
			maxID:     1,
			wantError: false,
		},
		"Invalid-1": {
			maxID:     0,
			wantError: true,
		},
		"Invalid-2": {
			maxID:     -1,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithQueryMaxID(tc.maxID)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryOptionSet(option, query); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.maxID), query.Get("maxId"))
			}
		})
	}
}

func TestActivityOptionService_WithCount(t *testing.T) {
	o := backlog.ActivityOptionService{}

	cases := map[string]struct {
		count     int
		wantError bool
	}{
		"Valid-1": {
			count:     1,
			wantError: false,
		},
		"Valid-2": {
			count:     100,
			wantError: false,
		},
		"Invalid-1": {
			count:     0,
			wantError: true,
		},
		"Invalid-2": {
			count:     -1,
			wantError: true,
		},
		"Invalid-3": {
			count:     101,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithQueryCount(tc.count)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryOptionSet(option, query); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.count), query.Get("count"))
			}
		})
	}
}

func TestActivityOptionService_WithOrder(t *testing.T) {
	o := backlog.ActivityOptionService{}

	cases := map[string]struct {
		order     backlog.Order
		wantError bool
	}{
		"asc": {
			order:     backlog.OrderAsc,
			wantError: false,
		},
		"desc": {
			order:     backlog.OrderDesc,
			wantError: false,
		},
		"Invalid": {
			order:     "test",
			wantError: true,
		},
		"Empty": {
			order:     "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithQueryOrder(tc.order)
			query := backlog.NewQueryParams()

			if err := backlog.ExportQueryOptionSet(option, query); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(tc.order), query.Get("order"))
			}
		})
	}
}

func TestProjectOptionService_WithQueryAll(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithQueryAll(tc.enabeld)
			query := backlog.NewQueryParams()
			err := backlog.ExportQueryOptionSet(option, query)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), query.Get("all"))
		})
	}
}

func TestProjectOptionService_WithFormKey(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		key       string
		wantError bool
	}{
		"Valid": {
			key:       "TEST",
			wantError: false,
		},
		"Empty": {
			key:       "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormKey(tc.key)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.key, form.Get("key"))
			}
		})
	}
}

func TestProjectOptionService_WithFormName(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		name      string
		wantError bool
	}{
		"Valid": {
			name:      "test",
			wantError: false,
		},
		"Empty": {
			name:      "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormName(tc.name)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, form.Get("name"))
			}
		})
	}
}

func TestProjectOptionService_WithFormChartEnabled(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormChartEnabled(tc.enabeld)
			form := backlog.NewFormParams()
			err := backlog.ExportFormOptionSet(option, form)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), form.Get("chartEnabled"))
		})
	}
}

func TestProjectOptionService_WithFormSubtaskingEnabled(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormSubtaskingEnabled(tc.enabeld)
			form := backlog.NewFormParams()
			err := backlog.ExportFormOptionSet(option, form)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), form.Get("subtaskingEnabled"))
		})
	}
}

func TestProjectOptionService_WithFormProjectLeaderCanEditProjectLeader(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormProjectLeaderCanEditProjectLeader(tc.enabeld)
			form := backlog.NewFormParams()
			err := backlog.ExportFormOptionSet(option, form)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), form.Get("projectLeaderCanEditProjectLeader"))
		})
	}
}

func TestProjectOptionService_WithFormTextFormattingRule(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		format    backlog.Format
		wantError bool
	}{
		"backlog": {
			format:    backlog.FormatBacklog,
			wantError: false,
		},
		"markdown": {
			format:    backlog.FormatMarkdown,
			wantError: false,
		},
		"Invalid": {
			format:    "test",
			wantError: true,
		},
		"Empty": {
			format:    "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormTextFormattingRule(tc.format)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(tc.format), form.Get("textFormattingRule"))
			}
		})
	}
}

func TestProjectOptionService_WithFormArchived(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		archived bool
	}{
		"true": {
			archived: true,
		},
		"false": {
			archived: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithQueryArchived(tc.archived)
			query := backlog.NewQueryParams()
			err := backlog.ExportQueryOptionSet(option, query)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.archived), query.Get("archived"))
		})
	}
}

func TestUserOptionService_WithFormPassword(t *testing.T) {
	o := backlog.UserOptionService{}

	cases := map[string]struct {
		password  string
		wantError bool
	}{
		"Valid-1": {
			password:  "password",
			wantError: false,
		},
		"Valid-2": {
			password:  "@password#1234",
			wantError: false,
		},
		"Empty": {
			password:  "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormPassword(tc.password)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.password, form.Get("password"))
			}
		})
	}
}

func TestUserOptionService_WithFormName(t *testing.T) {
	o := backlog.UserOptionService{}

	cases := map[string]struct {
		name      string
		wantError bool
	}{
		"Valid": {
			name:      "test",
			wantError: false,
		},
		"Empty": {
			name:      "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormName(tc.name)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, form.Get("name"))
			}
		})
	}
}

func TestUserOptionService_withMailAddress(t *testing.T) {
	o := backlog.UserOptionService{}

	cases := map[string]struct {
		mailAddress string
		wantError   bool
	}{
		"Valid-1": {
			mailAddress: "mail@test.com",
			wantError:   false,
		},
		"Valid-2": {
			mailAddress: "mail_test@test.com",
			wantError:   false,
		},
		"Valid-3": {
			mailAddress: "mail-test@test.com",
			wantError:   false,
		},
		// TODO
		// "inalid": {
		// 	mailAddress:  "test",
		// 	wantError: true,
		// },
		"Empty": {
			mailAddress: "",
			wantError:   true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormMailAddress(tc.mailAddress)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mailAddress, form.Get("mailAddress"))
			}
		})
	}
}

func TestUserOptionService_WithFormRoleType(t *testing.T) {
	o := backlog.UserOptionService{}

	cases := map[string]struct {
		roleType  backlog.Role
		want      string
		wantError bool
	}{
		"RoleAdministrator": {
			roleType:  backlog.RoleAdministrator,
			want:      "1",
			wantError: false,
		},
		"RoleNormalUser": {
			roleType:  backlog.RoleNormalUser,
			want:      "2",
			wantError: false,
		},
		"RoleReporter": {
			roleType:  backlog.RoleReporter,
			want:      "3",
			wantError: false,
		},
		"Viewer": {
			roleType:  backlog.RoleViewer,
			want:      "4",
			wantError: false,
		},
		"RoleGuestReporter": {
			roleType:  backlog.RoleGuestReporter,
			want:      "5",
			wantError: false,
		},
		"RoleGuestViewer": {
			roleType:  backlog.RoleGuestViewer,
			want:      "6",
			wantError: false,
		},
		"Invalid-1": {
			roleType:  0,
			want:      "6",
			wantError: true,
		},
		"Invalid-2": {
			roleType:  -1,
			want:      "6",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormRoleType(tc.roleType)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, form.Get("roleType"))
			}
		})
	}
}

func TestWikiOptionService_WithFormKeyword(t *testing.T) {
	o := backlog.WikiOptionService{}

	cases := map[string]struct {
		keyword string
	}{
		"Valid": {
			keyword: "test",
		},
		"Empty": {
			keyword: "",
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithQueryKeyword(tc.keyword)
			query := backlog.NewQueryParams()
			err := backlog.ExportQueryOptionSet(option, query)
			assert.NoError(t, err)
			assert.Equal(t, tc.keyword, query.Get("keyword"))
		})
	}
}

func TestWikiOptionService_WithFormName(t *testing.T) {
	o := backlog.WikiOptionService{}

	cases := map[string]struct {
		name      string
		wantError bool
	}{
		"Valid": {
			name:      "test",
			wantError: false,
		},
		"Empty": {
			name:      "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormName(tc.name)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, form.Get("name"))
			}
		})
	}
}

func TestWikiOptionService_WithFormContent(t *testing.T) {
	o := backlog.WikiOptionService{}

	cases := map[string]struct {
		content   string
		wantError bool
	}{
		"Valid": {
			content:   "content",
			wantError: false,
		},
		"Empty": {
			content:   "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormContent(tc.content)
			form := backlog.NewFormParams()

			if err := backlog.ExportFormOptionSet(option, form); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.content, form.Get("content"))
			}
		})
	}
}

func TestWikiOptionService_WithFormMailNotify(t *testing.T) {
	o := backlog.WikiOptionService{}

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := o.WithFormMailNotify(tc.enabeld)
			form := backlog.NewFormParams()
			err := backlog.ExportFormOptionSet(option, form)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), form.Get("mailNotify"))
		})
	}
}
