package backlog_test

import (
	"strconv"
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestOptionType_String(t *testing.T) {
	cases := map[string]struct {
		optionType backlog.ExportOptionType
		want       string
	}{
		"ActivityTypeIDs": {
			optionType: backlog.ExportOptionActivityTypeIDs,
			want:       "ActivityTypeIDs",
		},
		"All": {
			optionType: backlog.ExportOptionAll,
			want:       "All",
		},
		"Archived": {
			optionType: backlog.ExportOptionArchived,
			want:       "Archived",
		},
		"ChartEnabled": {
			optionType: backlog.ExportOptionChartEnabled,
			want:       "ChartEnabled",
		},
		"Content": {
			optionType: backlog.ExportOptionContent,
			want:       "Content",
		},
		"Count": {
			optionType: backlog.ExportOptionCount,
			want:       "Count",
		},
		"Key": {
			optionType: backlog.ExportOptionKey,
			want:       "Key",
		},
		"Keyword": {
			optionType: backlog.ExportOptionKeyword,
			want:       "Keyword",
		},
		"Name": {
			optionType: backlog.ExportOptionName,
			want:       "Name",
		},
		"MailAddress": {
			optionType: backlog.ExportOptionMailAddress,
			want:       "MailAddress",
		},
		"MailNotify": {
			optionType: backlog.ExportOptionMailNotify,
			want:       "MailNotify",
		},
		"MaxID": {
			optionType: backlog.ExportOptionMaxID,
			want:       "MaxID",
		},
		"MinID": {
			optionType: backlog.ExportOptionMinID,
			want:       "MinID",
		},
		"Order": {
			optionType: backlog.ExportOptionOrder,
			want:       "Order",
		},
		"Password": {
			optionType: backlog.ExportOptionPassword,
			want:       "Password",
		},
		"ProjectLeaderCanEditProjectLeader": {
			optionType: backlog.ExportOptionProjectLeaderCanEditProjectLeader,
			want:       "ProjectLeaderCanEditProjectLeader",
		},
		"RoleType": {
			optionType: backlog.ExportOptionRoleType,
			want:       "RoleType",
		},
		"SubtaskingEnabled": {
			optionType: backlog.ExportOptionSubtaskingEnabled,
			want:       "SubtaskingEnabled",
		},
		"TextFormattingRule": {
			optionType: backlog.ExportOptionTextFormattingRule,
			want:       "TextFormattingRule",
		},
		"unknown": {
			optionType: backlog.ExportOptionType(0),
			want:       "unknown",
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.optionType.String(), tc.want)
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
			option := o.WithActivityTypeIDs(tc.typeIDs)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportActivityOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				v := *params.ExportURLValues()
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
			option := o.WithMinID(tc.minID)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportActivityOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.minID), params.Get("minId"))
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
			option := o.WithMaxID(tc.maxID)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportActivityOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.maxID), params.Get("maxId"))
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
			option := o.WithCount(tc.count)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportActivityOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.count), params.Get("count"))
			}
		})
	}
}

func TestActivityOptionService_WithOrder(t *testing.T) {
	o := backlog.ActivityOptionService{}

	cases := map[string]struct {
		order     backlog.ExportOrder
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
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportActivityOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(tc.order), params.Get("order"))
			}
		})
	}
}

func TestProjectOptionService_WithAll(t *testing.T) {
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
			option := o.WithAll(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := backlog.ExportProjectOptionSet(option, params)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("all"))
		})
	}
}

func TestProjectOptionService_WithKey(t *testing.T) {
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
			option := o.WithKey(tc.key)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportProjectOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.key, params.Get("key"))
			}
		})
	}
}

func TestProjectOptionService_WithName(t *testing.T) {
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
			option := o.WithName(tc.name)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportProjectOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, params.Get("name"))
			}
		})
	}
}

func TestProjectOptionService_WithChartEnabled(t *testing.T) {
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
			option := o.WithChartEnabled(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := backlog.ExportProjectOptionSet(option, params)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("chartEnabled"))
		})
	}
}

func TestProjectOptionService_WithSubtaskingEnabled(t *testing.T) {
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
			option := o.WithSubtaskingEnabled(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := backlog.ExportProjectOptionSet(option, params)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("subtaskingEnabled"))
		})
	}
}

func TestProjectOptionService_WithProjectLeaderCanEditProjectLeader(t *testing.T) {
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
			option := o.WithProjectLeaderCanEditProjectLeader(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := backlog.ExportProjectOptionSet(option, params)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("projectLeaderCanEditProjectLeader"))
		})
	}
}

func TestProjectOptionService_WithTextFormattingRule(t *testing.T) {
	o := backlog.ProjectOptionService{}

	cases := map[string]struct {
		format    backlog.ExportFormat
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
			option := o.WithTextFormattingRule(tc.format)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportProjectOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(tc.format), params.Get("textFormattingRule"))
			}
		})
	}
}

func TestProjectOptionService_WithArchived(t *testing.T) {
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
			option := o.WithArchived(tc.archived)
			params := backlog.ExportNewRequestParams()
			err := backlog.ExportProjectOptionSet(option, params)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.archived), params.Get("archived"))
		})
	}
}

func TestUserOptionService_WithPassword(t *testing.T) {
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
			option := o.WithPassword(tc.password)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportUserOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.password, params.Get("password"))
			}
		})
	}
}

func TestUserOptionService_WithName(t *testing.T) {
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
			option := o.WithName(tc.name)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportUserOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, params.Get("name"))
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
			option := o.WithMailAddress(tc.mailAddress)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportUserOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mailAddress, params.Get("mailAddress"))
			}
		})
	}
}

func TestUserOptionService_WithRoleType(t *testing.T) {
	o := backlog.UserOptionService{}

	cases := map[string]struct {
		roleType  backlog.ExportRole
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
			option := o.WithRoleType(tc.roleType)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportUserOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, params.Get("roleType"))
			}
		})
	}
}

func TestWikiOptionService_WithKeyword(t *testing.T) {
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
			option := o.WithKeyword(tc.keyword)
			params := backlog.ExportNewRequestParams()
			err := backlog.ExportWikiOptionSet(option, params)
			assert.NoError(t, err)
			assert.Equal(t, tc.keyword, params.Get("keyword"))
		})
	}
}

func TestWikiOptionService_WithName(t *testing.T) {
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
			option := o.WithName(tc.name)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportWikiOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, params.Get("name"))
			}
		})
	}
}

func TestWikiOptionService_WithContent(t *testing.T) {
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
			option := o.WithContent(tc.content)
			params := backlog.ExportNewRequestParams()

			if err := backlog.ExportWikiOptionSet(option, params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.content, params.Get("content"))
			}
		})
	}
}

func TestWikiOptionService_WithMailNotify(t *testing.T) {
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
			option := o.WithMailNotify(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := backlog.ExportWikiOptionSet(option, params)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("mailNotify"))
		})
	}
}
