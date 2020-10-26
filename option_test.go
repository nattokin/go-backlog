package backlog_test

import (
	"strconv"
	"testing"

	backlog "github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestActivityOptionService_WithActivityTypeIDs(t *testing.T) {
	s := backlog.ExportNewActivityService(nil).Option

	cases := map[string]struct {
		typeIDs   []int
		want      []string
		wantError bool
	}{
		"valid-1": {
			typeIDs:   []int{1},
			want:      []string{"1"},
			wantError: false,
		},
		"valid-2": {
			typeIDs:   []int{26},
			want:      []string{"26"},
			wantError: false,
		},
		"valid-3": {
			typeIDs: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
			want: []string{
				"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13",
				"14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26",
			},
			wantError: false,
		},
		"invalid-1": {
			typeIDs:   []int{0},
			want:      nil,
			wantError: true,
		},
		"invalid-2": {
			typeIDs:   []int{-1},
			want:      nil,
			wantError: true,
		},
		"invalid-3": {
			typeIDs:   []int{27},
			want:      nil,
			wantError: true,
		},
		"invalid-4": {
			typeIDs:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27},
			want:      nil,
			wantError: true,
		},
		"invalid-5": {
			typeIDs:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
			want:      nil,
			wantError: true,
		},
		"empty": {
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
			option := s.WithActivityTypeIDs(tc.typeIDs)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				v := *params.ExportURLValues()
				assert.Equal(t, tc.want, v["activityTypeId[]"])
			}
		})
	}
}

func TestActivityOptionService_WithMinID(t *testing.T) {
	s := backlog.ExportNewActivityService(nil).Option

	cases := map[string]struct {
		minID     int
		wantError bool
	}{
		"valid": {
			minID:     1,
			wantError: false,
		},
		"invalid-1": {
			minID:     0,
			wantError: true,
		},
		"invalid-2": {
			minID:     -1,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithMinID(tc.minID)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, strconv.Itoa(tc.minID), params.Get("minId"))
			}
		})
	}
}

func TestActivityOptionService_WithMaxID(t *testing.T) {
	s := backlog.ExportNewActivityService(nil).Option

	cases := map[string]struct {
		maxID     int
		wantError bool
	}{
		"valid": {
			maxID:     1,
			wantError: false,
		},
		"invalid-1": {
			maxID:     0,
			wantError: true,
		},
		"invalid-2": {
			maxID:     -1,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithMaxID(tc.maxID)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, strconv.Itoa(tc.maxID), params.Get("maxId"))
			}
		})
	}
}

func TestActivityOptionService_WithCount(t *testing.T) {
	s := backlog.ExportNewActivityService(nil).Option

	cases := map[string]struct {
		count     int
		wantError bool
	}{
		"valid-1": {
			count:     1,
			wantError: false,
		},
		"valid-2": {
			count:     100,
			wantError: false,
		},
		"invalid-1": {
			count:     0,
			wantError: true,
		},
		"invalid-2": {
			count:     -1,
			wantError: true,
		},
		"invalid-3": {
			count:     101,
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithCount(tc.count)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, strconv.Itoa(tc.count), params.Get("count"))
			}
		})
	}
}

func TestActivityOptionService_WithOrder(t *testing.T) {
	s := backlog.ExportNewActivityService(nil).Option

	cases := map[string]struct {
		order     string
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
		"invalid": {
			order:     "test",
			wantError: true,
		},
		"empty": {
			order:     "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithOrder(tc.order)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.order, params.Get("order"))
			}
		})
	}
}

func TestProjectOptionService_WithKey(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

	cases := map[string]struct {
		key       string
		wantError bool
	}{
		"valid": {
			key:       "TEST",
			wantError: false,
		},
		"empty": {
			key:       "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithKey(tc.key)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.key, params.Get("key"))
			}
		})
	}
}

func TestProjectOptionService_WithName(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

	cases := map[string]struct {
		name      string
		wantError bool
	}{
		"valid": {
			name:      "test",
			wantError: false,
		},
		"empty": {
			name:      "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithName(tc.name)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.name, params.Get("name"))
			}
		})
	}
}

func TestProjectOptionService_WithChartEnabled(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

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
			option := s.WithChartEnabled(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("chartEnabled"))
		})
	}
}

func TestProjectOptionService_WithSubtaskingEnabled(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

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
			option := s.WithSubtaskingEnabled(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("subtaskingEnabled"))
		})
	}
}

func TestProjectOptionService_WithProjectLeaderCanEditProjectLeader(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

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
			option := s.WithProjectLeaderCanEditProjectLeader(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("projectLeaderCanEditProjectLeader"))
		})
	}
}

func TestProjectOptionService_WithTextFormattingRule(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option
	// TODO
	cases := map[string]struct {
		format    string
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
		"invalid": {
			format:    "test",
			wantError: true,
		},
		"empty": {
			format:    "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithTextFormattingRule(tc.format)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.format, params.Get("textFormattingRule"))
			}
		})
	}
}

func TestProjectOptionService_WithArchived(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

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
			option := s.WithArchived(tc.archived)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.archived), params.Get("archived"))
		})
	}
}

func TestWikiOptionService_WithName(t *testing.T) {
	s := backlog.ExportNewWikiService(nil).Option

	cases := map[string]struct {
		name      string
		wantError bool
	}{
		"valid": {
			name:      "test",
			wantError: false,
		},
		"empty": {
			name:      "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithName(tc.name)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.name, params.Get("name"))
			}
		})
	}
}

func TestWikiOptionService_WithContent(t *testing.T) {
	s := backlog.ExportNewWikiService(nil).Option

	cases := map[string]struct {
		content   string
		wantError bool
	}{
		"valid": {
			content:   "content",
			wantError: false,
		},
		"empty": {
			content:   "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithContent(tc.content)
			params := backlog.ExportNewRequestParams()

			if err := option(params); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.content, params.Get("content"))
			}
		})
	}
}

func TestWikiOptionService_WithMailNotify(t *testing.T) {
	s := backlog.ExportNewWikiService(nil).Option

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
			option := s.WithMailNotify(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("mailNotify"))
		})
	}
}
