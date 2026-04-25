package backlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/model"
)

func Test_changeLogFromModel(t *testing.T) {
	t.Parallel()

	input := &model.ChangeLog{Field: "status", NewValue: "4", OriginalValue: "1"}
	want := &ChangeLog{Field: "status", NewValue: "4", OriginalValue: "1"}
	assert.Equal(t, want, changeLogFromModel(input))
}

func Test_customFieldFromModel(t *testing.T) {
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
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, customFieldFromModel(tc.input))
		})
	}
}

func Test_customFieldItemFromModel(t *testing.T) {
	t.Parallel()

	input := &model.CustomFieldItem{ID: 1, Name: "Windows 8", DisplayOrder: 0}
	want := &CustomFieldItem{ID: 1, Name: "Windows 8", DisplayOrder: 0}
	assert.Equal(t, want, customFieldItemFromModel(input))
}

func Test_diskUsageProjectFromModel(t *testing.T) {
	t.Parallel()

	input := &model.DiskUsageProject{
		DiskUsageBase: model.DiskUsageBase{
			Issue:      11931,
			Wiki:       0,
			File:       512,
			Subversion: 0,
			Git:        1024,
			GitLFS:     0,
		},
		ProjectID: 1,
	}
	want := &DiskUsageProject{
		ProjectID:  1,
		Issue:      11931,
		Wiki:       0,
		File:       512,
		Subversion: 0,
		Git:        1024,
		GitLFS:     0,
	}
	assert.Equal(t, want, diskUsageProjectFromModel(input))
}

func Test_diskUsageSpaceFromModel(t *testing.T) {
	cases := map[string]struct {
		input *model.DiskUsageSpace
		want  *DiskUsageSpace
	}{
		"with_details": {
			input: &model.DiskUsageSpace{
				DiskUsageBase: model.DiskUsageBase{
					Issue:      119511,
					Wiki:       0,
					File:       0,
					Subversion: 0,
					Git:        0,
					GitLFS:     0,
				},
				Capacity: 1073741824,
				Details: []*model.DiskUsageProject{
					{
						DiskUsageBase: model.DiskUsageBase{Issue: 11931},
						ProjectID:     1,
					},
				},
			},
			want: &DiskUsageSpace{
				Capacity:   1073741824,
				Issue:      119511,
				Wiki:       0,
				File:       0,
				Subversion: 0,
				Git:        0,
				GitLFS:     0,
				Details: []*DiskUsageProject{
					{ProjectID: 1, Issue: 11931},
				},
			},
		},
		"empty_details": {
			input: &model.DiskUsageSpace{
				Capacity: 512,
				Details:  []*model.DiskUsageProject{},
			},
			want: &DiskUsageSpace{
				Capacity: 512,
				Details:  []*DiskUsageProject{},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, diskUsageSpaceFromModel(tc.input))
		})
	}
}

func Test_issueFromModel(t *testing.T) {
	created := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	user := &model.User{ID: 1, UserID: "admin", Name: "admin", RoleType: model.RoleAdministrator, Lang: "ja", MailAddress: "admin@example.com"}
	wantUser := &User{ID: 1, UserID: "admin", Name: "admin", RoleType: RoleAdministrator, Lang: "ja", MailAddress: "admin@example.com"}

	cases := map[string]struct {
		input *model.Issue
		want  *Issue
	}{
		"full": {
			input: &model.Issue{
				ID:          1,
				ProjectID:   10,
				IssueKey:    "BLG-1",
				KeyID:       1,
				IssueType:   &model.IssueType{ID: 2, ProjectID: 10, Name: "Task", Color: "#7ea800", DisplayOrder: 0},
				Summary:     "test issue",
				Description: "desc",
				Resolutions: []*model.Resolution{
					{ID: 0, Name: "Fixed"},
					{ID: 1, Name: "Won't Fix"},
				},
				Priority: &model.Priority{ID: 3, Name: "Normal"},
				Status:   &model.Status{ID: 1, Name: "Open"},
				Assignee: user,
				Category: []*model.Category{
					{ID: 5, Name: "Frontend", DisplayOrder: 0},
				},
				Versions: []*model.Version{
					{ID: 30, Name: "v1.0"},
				},
				Milestone: []*model.Version{
					{ID: 31, Name: "v1.1"},
				},
				StartDate:      created,
				DueDate:        updated,
				EstimatedHours: 8,
				ActualHours:    4,
				ParentIssueID:  0,
				CreatedUser:    user,
				Created:        created,
				UpdatedUser:    user,
				Updated:        updated,
				CustomFields: []*model.CustomField{
					{ID: 1, TypeID: 6, Name: "OS", Items: []*model.CustomFieldItem{
						{ID: 1, Name: "Windows", DisplayOrder: 0},
					}},
				},
				Attachments: []*model.Attachment{
					{ID: 10, Name: "file.txt", Size: 100, CreatedUser: user, Created: created},
				},
				SharedFiles: []*model.SharedFile{
					{ID: 20, Type: "file", Dir: "/", Name: "shared.txt", Size: 200, CreatedUser: user, Created: created},
				},
				Stars: []*model.Star{
					{ID: 75, Comment: "good", URL: "https://example.com", Title: "title", Presenter: user, Created: created},
				},
			},
			want: &Issue{
				ID:          1,
				ProjectID:   10,
				IssueKey:    "BLG-1",
				KeyID:       1,
				IssueType:   &IssueType{ID: 2, ProjectID: 10, Name: "Task", Color: "#7ea800", DisplayOrder: 0},
				Summary:     "test issue",
				Description: "desc",
				Resolutions: []*Resolution{
					{ID: 0, Name: "Fixed"},
					{ID: 1, Name: "Won't Fix"},
				},
				Priority: &Priority{ID: 3, Name: "Normal"},
				Status:   &Status{ID: 1, Name: "Open"},
				Assignee: wantUser,
				Category: []*Category{
					{ID: 5, Name: "Frontend", DisplayOrder: 0},
				},
				Versions: []*Version{
					{ID: 30, Name: "v1.0"},
				},
				Milestone: []*Version{
					{ID: 31, Name: "v1.1"},
				},
				StartDate:      created,
				DueDate:        updated,
				EstimatedHours: 8,
				ActualHours:    4,
				ParentIssueID:  0,
				CreatedUser:    wantUser,
				Created:        created,
				UpdatedUser:    wantUser,
				Updated:        updated,
				CustomFields: []*CustomField{
					{ID: 1, TypeID: 6, Name: "OS", Items: []*CustomFieldItem{
						{ID: 1, Name: "Windows", DisplayOrder: 0},
					}},
				},
				Attachments: []*Attachment{
					{ID: 10, Name: "file.txt", Size: 100, CreatedUser: wantUser, Created: created},
				},
				SharedFiles: []*SharedFile{
					{ID: 20, Type: "file", Dir: "/", Name: "shared.txt", Size: 200, CreatedUser: wantUser, Created: created},
				},
				Stars: []*Star{
					{ID: 75, Comment: "good", URL: "https://example.com", Title: "title", Presenter: wantUser, Created: created},
				},
			},
		},
		"nil_elements": {
			input: &model.Issue{
				Resolutions:  []*model.Resolution{nil},
				Category:     []*model.Category{nil},
				Versions:     []*model.Version{nil},
				Milestone:    []*model.Version{nil},
				CustomFields: []*model.CustomField{nil},
				Attachments:  []*model.Attachment{nil},
				SharedFiles:  []*model.SharedFile{nil},
				Stars:        []*model.Star{nil},
			},
			want: &Issue{
				Resolutions:  []*Resolution{nil},
				Category:     []*Category{nil},
				Versions:     []*Version{nil},
				Milestone:    []*Version{nil},
				CustomFields: []*CustomField{nil},
				Attachments:  []*Attachment{nil},
				SharedFiles:  []*SharedFile{nil},
				Stars:        []*Star{nil},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, issueFromModel(tc.input))
		})
	}
}

func Test_pullRequestFromModel(t *testing.T) {
	t.Parallel()

	created := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	closeAt := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	mergeAt := time.Date(2024, 5, 2, 0, 0, 0, 0, time.UTC)

	user := &model.User{ID: 1, UserID: "admin", Name: "admin", RoleType: model.RoleAdministrator, Lang: "ja", MailAddress: "admin@example.com"}
	wantUser := &User{ID: 1, UserID: "admin", Name: "admin", RoleType: RoleAdministrator, Lang: "ja", MailAddress: "admin@example.com"}

	input := &model.PullRequest{
		ID:           2,
		ProjectID:    3,
		RepositoryID: 5,
		Number:       1,
		Summary:      "test PR",
		Description:  "PR desc",
		Base:         "main",
		Branch:       "feature",
		Status:       &model.Status{ID: 1, Name: "Open"},
		Assignee:     user,
		Issue:        &model.Issue{ID: 10, Summary: "related issue"},
		BaseCommit:   "abc123",
		BranchCommit: "def456",
		CloseAt:      closeAt,
		MergeAt:      mergeAt,
		CreatedUser:  user,
		Created:      created,
		UpdatedUser:  user,
		Updated:      updated,
		Attachments: []*model.Attachment{
			{ID: 10, Name: "file.txt", Size: 100, CreatedUser: user, Created: created},
		},
		Stars: []*model.Star{
			{ID: 75, Comment: "good", URL: "https://example.com", Title: "title", Presenter: user, Created: created},
		},
	}
	want := &PullRequest{
		ID:           2,
		ProjectID:    3,
		RepositoryID: 5,
		Number:       1,
		Summary:      "test PR",
		Description:  "PR desc",
		Base:         "main",
		Branch:       "feature",
		Status:       &Status{ID: 1, Name: "Open"},
		Assignee:     wantUser,
		Issue: &Issue{
			ID:      10,
			Summary: "related issue",
		},
		BaseCommit:   "abc123",
		BranchCommit: "def456",
		CloseAt:      closeAt,
		MergeAt:      mergeAt,
		CreatedUser:  wantUser,
		Created:      created,
		UpdatedUser:  wantUser,
		Updated:      updated,
		Attachments: []*Attachment{
			{ID: 10, Name: "file.txt", Size: 100, CreatedUser: wantUser, Created: created},
		},
		Stars: []*Star{
			{ID: 75, Comment: "good", URL: "https://example.com", Title: "title", Presenter: wantUser, Created: created},
		},
	}
	assert.Equal(t, want, pullRequestFromModel(input))
}

func Test_spaceFromModel(t *testing.T) {
	t.Parallel()

	created := time.Date(2008, 7, 6, 15, 0, 0, 0, time.UTC)
	updated := time.Date(2013, 6, 18, 7, 55, 37, 0, time.UTC)
	input := &model.Space{
		SpaceKey:           "nulab",
		Name:               "Nulab Inc.",
		OwnerID:            1,
		Lang:               "ja",
		Timezone:           "Asia/Tokyo",
		ReportSendTime:     "08:00:00",
		TextFormattingRule: model.FormatMarkdown,
		Created:            created,
		Updated:            updated,
	}
	want := &Space{
		SpaceKey:           "nulab",
		Name:               "Nulab Inc.",
		OwnerID:            1,
		Lang:               "ja",
		Timezone:           "Asia/Tokyo",
		ReportSendTime:     "08:00:00",
		TextFormattingRule: FormatMarkdown,
		Created:            created,
		Updated:            updated,
	}
	assert.Equal(t, want, spaceFromModel(input))
}

func Test_spaceNotificationFromModel(t *testing.T) {
	t.Parallel()

	updated := time.Date(2013, 6, 18, 7, 55, 37, 0, time.UTC)
	input := &model.SpaceNotification{
		Content: "Backlog is a project management tool.",
		Updated: updated,
	}
	want := &SpaceNotification{
		Content: "Backlog is a project management tool.",
		Updated: updated,
	}
	assert.Equal(t, want, spaceNotificationFromModel(input))
}

func Test_statusFromModel(t *testing.T) {
	t.Parallel()

	input := &model.Status{ID: 1, Name: "Open"}
	want := &Status{ID: 1, Name: "Open"}
	assert.Equal(t, want, statusFromModel(input))
}

func Test_versionFromModel(t *testing.T) {
	t.Parallel()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	due := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	input := &model.Version{
		ID:             30,
		ProjectID:      1,
		Name:           "v1.0",
		Description:    "first release",
		StartDate:      start,
		ReleaseDueDate: due,
		Archived:       false,
		DisplayOrder:   0,
	}
	want := &Version{
		ID:             30,
		ProjectID:      1,
		Name:           "v1.0",
		Description:    "first release",
		StartDate:      start,
		ReleaseDueDate: due,
		Archived:       false,
		DisplayOrder:   0,
	}
	assert.Equal(t, want, versionFromModel(input))
}

func Test_versionsFromModel(t *testing.T) {
	cases := map[string]struct {
		input []*model.Version
		want  []*Version
	}{
		"with_elements": {
			input: []*model.Version{
				{ID: 1, Name: "v1.0"},
				{ID: 2, Name: "v2.0"},
			},
			want: []*Version{
				{ID: 1, Name: "v1.0"},
				{ID: 2, Name: "v2.0"},
			},
		},
		"empty": {
			input: []*model.Version{},
			want:  []*Version{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, versionsFromModel(tc.input))
		})
	}
}

func Test_fromModel_nil(t *testing.T) {
	cases := map[string]struct {
		call func() any
	}{
		"activity":           {call: func() any { return activityFromModel(nil) }},
		"activity_content":   {call: func() any { return activityContentFromModel(nil) }},
		"attachment":         {call: func() any { return attachmentFromModel(nil) }},
		"change_log":         {call: func() any { return changeLogFromModel(nil) }},
		"comment":            {call: func() any { return commentFromModel(nil) }},
		"custom_field":       {call: func() any { return customFieldFromModel(nil) }},
		"custom_field_item":  {call: func() any { return customFieldItemFromModel(nil) }},
		"disk_usage_project": {call: func() any { return diskUsageProjectFromModel(nil) }},
		"disk_usage_space":   {call: func() any { return diskUsageSpaceFromModel(nil) }},
		"issue":              {call: func() any { return issueFromModel(nil) }},
		"notification":       {call: func() any { return notificationFromModel(nil) }},
		"project":            {call: func() any { return projectFromModel(nil) }},
		"pull_request":       {call: func() any { return pullRequestFromModel(nil) }},
		"shared_file":        {call: func() any { return sharedFileFromModel(nil) }},
		"space":              {call: func() any { return spaceFromModel(nil) }},
		"space_notification": {call: func() any { return spaceNotificationFromModel(nil) }},
		"star":               {call: func() any { return starFromModel(nil) }},
		"stars":              {call: func() any { return starsFromModel(nil) }},
		"status":             {call: func() any { return statusFromModel(nil) }},
		"tag":                {call: func() any { return tagFromModel(nil) }},
		"user":               {call: func() any { return userFromModel(nil) }},
		"version":            {call: func() any { return versionFromModel(nil) }},
		"versions":           {call: func() any { return versionsFromModel(nil) }},
		"wiki":               {call: func() any { return wikiFromModel(nil) }},
		"wiki_history":       {call: func() any { return wikiHistoryFromModel(nil) }},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Nil(t, tc.call())
		})
	}
}
