package backlog

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivityOptionService_WithActivityTypeIDs(t *testing.T) {
	o := newActivityOptionService()

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
			t.Parallel()

			option := o.WithQueryActivityTypeIDs(tc.typeIDs)
			query := NewQueryParams()

			if err := option.set(query); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, query.Values)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, (*query.Values)["activityTypeId[]"])
			}
		})

	}
}

func TestActivityOptionService_WithMinID(t *testing.T) {
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
			t.Parallel()

			o := newActivityOptionService()
			option := o.WithQueryMinID(tc.minID)
			query := NewQueryParams()

			if err := option.set(query); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, query.Get("minId"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.minID), query.Get("minId"))
			}
		})

	}
}

func TestActivityOptionService_WithMaxID(t *testing.T) {
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
			t.Parallel()

			o := newActivityOptionService()
			option := o.WithQueryMaxID(tc.maxID)
			query := NewQueryParams()

			if err := option.set(query); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, query.Get("maxId"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.maxID), query.Get("maxId"))
			}
		})

	}
}

func TestActivityOptionService_WithCount(t *testing.T) {
	o := newActivityOptionService()

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
			t.Parallel()

			option := o.WithQueryCount(tc.count)
			query := NewQueryParams()

			if err := option.set(query); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, query.Get("count"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.count), query.Get("count"))
			}
		})

	}
}

func TestActivityOptionService_WithOrder(t *testing.T) {
	o := newActivityOptionService()

	cases := map[string]struct {
		order     Order
		wantError bool
	}{
		"asc": {
			order:     OrderAsc,
			wantError: false,
		},
		"desc": {
			order:     OrderDesc,
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
			t.Parallel()

			option := o.WithQueryOrder(tc.order)
			query := NewQueryParams()

			if err := option.set(query); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, query.Get("order"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(tc.order), query.Get("order"))
			}
		})

	}
}

func TestProjectOptionService_WithQueryAll(t *testing.T) {
	o := newProjectOptionService()

	cases := map[string]struct {
		enabled bool
	}{
		"true": {
			enabled: true,
		},
		"false": {
			enabled: false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			option := o.WithQueryAll(tc.enabled)
			query := NewQueryParams()
			err := option.set(query)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabled), query.Get("all"))
		})

	}
}

func TestProjectOptionService_WithFormKey(t *testing.T) {
	o := newProjectOptionService()

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
			t.Parallel()

			option := o.WithFormKey(tc.key)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("key"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.key, form.Get("key"))
			}
		})

	}
}

func TestProjectOptionService_WithFormName(t *testing.T) {
	o := newProjectOptionService()

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
			t.Parallel()

			option := o.WithFormName(tc.name)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("name"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, form.Get("name"))
			}
		})

	}
}

func TestProjectOptionService_WithFormChartEnabled(t *testing.T) {
	o := newProjectOptionService()

	cases := map[string]struct {
		enabled bool
	}{
		"true": {
			enabled: true,
		},
		"false": {
			enabled: false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			option := o.WithFormChartEnabled(tc.enabled)
			form := NewFormParams()
			err := option.set(form)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabled), form.Get("chartEnabled"))
		})

	}
}

func TestProjectOptionService_WithFormSubtaskingEnabled(t *testing.T) {
	o := newProjectOptionService()

	cases := map[string]struct {
		enabled bool
	}{
		"true": {
			enabled: true,
		},
		"false": {
			enabled: false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			option := o.WithFormSubtaskingEnabled(tc.enabled)
			form := NewFormParams()
			err := option.set(form)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabled), form.Get("subtaskingEnabled"))
		})

	}
}

func TestProjectOptionService_WithFormProjectLeaderCanEditProjectLeader(t *testing.T) {
	o := newProjectOptionService()

	cases := map[string]struct {
		enabled bool
	}{
		"true": {
			enabled: true,
		},
		"false": {
			enabled: false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			option := o.WithFormProjectLeaderCanEditProjectLeader(tc.enabled)
			form := NewFormParams()
			err := option.set(form)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabled), form.Get("projectLeaderCanEditProjectLeader"))
		})

	}
}

func TestProjectOptionService_WithFormTextFormattingRule(t *testing.T) {
	o := newProjectOptionService()

	cases := map[string]struct {
		format    Format
		wantError bool
	}{
		"backlog": {
			format:    FormatBacklog,
			wantError: false,
		},
		"markdown": {
			format:    FormatMarkdown,
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
			t.Parallel()

			option := o.WithFormTextFormattingRule(tc.format)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("textFormattingRule"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(tc.format), form.Get("textFormattingRule"))
			}
		})

	}
}

func TestProjectOptionService_WithFormArchived(t *testing.T) {
	o := newProjectOptionService()

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
			t.Parallel()

			option := o.WithQueryArchived(tc.archived)
			query := NewQueryParams()
			err := option.set(query)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.archived), query.Get("archived"))
		})

	}
}

func TestUserOptionService(t *testing.T) {
	o := newUserOptionService()

	t.Run("Integer", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue bool
		}{
			"WithFormSendMail": {
				option:    o.WithFormSendMail(true),
				key:       formSendMail.Value(),
				wantValue: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	t.Run("Boolean", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue int
		}{
			"WithFormUserID": {
				option:    o.WithFormUserID(1),
				key:       formUserID.Value(),
				wantValue: 1,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.wantValue), form.Get(tc.key))
			})
		}
	})
}

func TestUserOptionService_WithFormPassword(t *testing.T) {
	o := newUserOptionService()

	cases := map[string]struct {
		password  string
		wantError bool
	}{
		"Valid-8char": {
			password:  "password",
			wantError: false,
		},
		"Valid-more": {
			password:  "@password#1234",
			wantError: false,
		},
		"Invalid-1": {
			password:  "1234567",
			wantError: true,
		},
		"Invalid-2": {
			password:  "asdEdfd",
			wantError: true,
		},
		"Invalid-empty": {
			password:  "",
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			option := o.WithFormPassword(tc.password)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("password"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.password, form.Get("password"))
			}
		})

	}
}

func TestUserOptionService_WithFormName(t *testing.T) {
	o := newUserOptionService()

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
			t.Parallel()

			option := o.WithFormName(tc.name)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("name"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, form.Get("name"))
			}
		})

	}
}

func TestUserOptionService_withMailAddress(t *testing.T) {
	o := newUserOptionService()

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
			t.Parallel()

			option := o.WithFormMailAddress(tc.mailAddress)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("mailAddress"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mailAddress, form.Get("mailAddress"))
			}
		})

	}
}

func TestUserOptionService_WithFormRoleType(t *testing.T) {
	o := newUserOptionService()

	cases := map[string]struct {
		roleType  Role
		want      string
		wantError bool
	}{
		"RoleAdministrator": {
			roleType:  RoleAdministrator,
			want:      "1",
			wantError: false,
		},
		"RoleNormalUser": {
			roleType:  RoleNormalUser,
			want:      "2",
			wantError: false,
		},
		"RoleReporter": {
			roleType:  RoleReporter,
			want:      "3",
			wantError: false,
		},
		"Viewer": {
			roleType:  RoleViewer,
			want:      "4",
			wantError: false,
		},
		"RoleGuestReporter": {
			roleType:  RoleGuestReporter,
			want:      "5",
			wantError: false,
		},
		"RoleGuestViewer": {
			roleType:  RoleGuestViewer,
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
			t.Parallel()

			option := o.WithFormRoleType(tc.roleType)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("roleType"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, form.Get("roleType"))
			}
		})

	}
}

func TestWikiOptionService_WithFormKeyword(t *testing.T) {
	o := newWikiOptionService()

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
			t.Parallel()

			option := o.WithQueryKeyword(tc.keyword)
			query := NewQueryParams()
			err := option.set(query)
			assert.NoError(t, err)
			assert.Equal(t, tc.keyword, query.Get("keyword"))
		})

	}
}

func TestWikiOptionService_WithFormName(t *testing.T) {
	o := newWikiOptionService()

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
			t.Parallel()

			option := o.WithFormName(tc.name)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("name"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, form.Get("name"))
			}
		})

	}
}

func TestWikiOptionService_WithFormContent(t *testing.T) {
	o := newWikiOptionService()

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
			t.Parallel()

			option := o.WithFormContent(tc.content)
			form := NewFormParams()

			if err := option.set(form); tc.wantError {
				assert.Error(t, err)
				assert.Empty(t, form.Get("content"))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.content, form.Get("content"))
			}
		})

	}
}

func TestWikiOptionService_WithFormMailNotify(t *testing.T) {
	o := newWikiOptionService()

	cases := map[string]struct {
		enabled bool
	}{
		"true": {
			enabled: true,
		},
		"false": {
			enabled: false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			option := o.WithFormMailNotify(tc.enabled)
			form := NewFormParams()
			err := option.set(form)
			assert.NoError(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabled), form.Get("mailNotify"))
		})

	}
}

// --- Option Service Helpers ---

// newQueryOptionService returns a test instance of QueryOptionService.
func newQueryOptionService() *QueryOptionService {
	return &QueryOptionService{}
}

// newFormOptionService returns a test instance of FormOptionService.
func newFormOptionService() *FormOptionService {
	return &FormOptionService{}
}

// newActivityOptionService returns a test instance of ActivityOptionService.
func newActivityOptionService() *ActivityOptionService {
	return &ActivityOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newProjectOptionService returns a test instance of ProjectOptionService.
func newProjectOptionService() *ProjectOptionService {
	return &ProjectOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newUserOptionService returns a test instance of UserOptionService.
func newUserOptionService() *UserOptionService {
	return &UserOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newWikiOptionService returns a test instance of WikiOptionService.
func newWikiOptionService() *WikiOptionService {
	return &WikiOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}
