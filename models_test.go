package backlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/model"
)

func Test_changeLogFromModel(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input *model.ChangeLog
		want  *ChangeLog
	}{
		"normal": {
			input: &model.ChangeLog{
				Field:         "status",
				NewValue:      "4",
				OriginalValue: "1",
			},
			want: &ChangeLog{
				Field:         "status",
				NewValue:      "4",
				OriginalValue: "1",
			},
		},
		"nil": {
			input: nil,
			want:  nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, changeLogFromModel(tc.input))
		})
	}
}

func Test_statusFromModel(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input *model.Status
		want  *Status
	}{
		"normal": {
			input: &model.Status{ID: 1, Name: "Open"},
			want:  &Status{ID: 1, Name: "Open"},
		},
		"nil": {
			input: nil,
			want:  nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, statusFromModel(tc.input))
		})
	}
}

func Test_customFieldItemFromModel(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input *model.CustomFieldItem
		want  *CustomFieldItem
	}{
		"normal": {
			input: &model.CustomFieldItem{ID: 1, Name: "Windows 8", DisplayOrder: 0},
			want:  &CustomFieldItem{ID: 1, Name: "Windows 8", DisplayOrder: 0},
		},
		"nil": {
			input: nil,
			want:  nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, customFieldItemFromModel(tc.input))
		})
	}
}

func Test_customFieldFromModel(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input *model.CustomField
		want  *CustomField
	}{
		"with_items": {
			input: &model.CustomField{
				ID:                     1,
				TypeID:                 6,
				Name:                   "custom",
				Description:            "desc",
				Required:               true,
				ApplicableIssueTypeIDs: []int{1, 2},
				AllowAddItem:           true,
				Items: []*model.CustomFieldItem{
					{ID: 1, Name: "Windows 8", DisplayOrder: 0},
					{ID: 2, Name: "Windows 10", DisplayOrder: 1},
				},
			},
			want: &CustomField{
				ID:                     1,
				TypeID:                 6,
				Name:                   "custom",
				Description:            "desc",
				Required:               true,
				ApplicableIssueTypeIDs: []int{1, 2},
				AllowAddItem:           true,
				Items: []*CustomFieldItem{
					{ID: 1, Name: "Windows 8", DisplayOrder: 0},
					{ID: 2, Name: "Windows 10", DisplayOrder: 1},
				},
			},
		},
		"empty_items": {
			input: &model.CustomField{
				ID:    2,
				Name:  "no items",
				Items: []*model.CustomFieldItem{},
			},
			want: &CustomField{
				ID:    2,
				Name:  "no items",
				Items: []*CustomFieldItem{},
			},
		},
		"nil": {
			input: nil,
			want:  nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, customFieldFromModel(tc.input))
		})
	}
}

func Test_versionFromModel(t *testing.T) {
	t.Parallel()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	due := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		input *model.Version
		want  *Version
	}{
		"normal": {
			input: &model.Version{
				ID:             30,
				ProjectID:      1,
				Name:           "v1.0",
				Description:    "first release",
				StartDate:      start,
				ReleaseDueDate: due,
				Archived:       false,
				DisplayOrder:   0,
			},
			want: &Version{
				ID:             30,
				ProjectID:      1,
				Name:           "v1.0",
				Description:    "first release",
				StartDate:      start,
				ReleaseDueDate: due,
				Archived:       false,
				DisplayOrder:   0,
			},
		},
		"nil": {
			input: nil,
			want:  nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, versionFromModel(tc.input))
		})
	}
}
