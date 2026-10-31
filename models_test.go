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
		"basic_fields_only": {
			input: &model.ChangeLog{Field: "status", NewValue: "4", OriginalValue: "1"},
			want:  &ChangeLog{Field: "status", NewValue: "4", OriginalValue: "1"},
		},
		"with_attachment_info": {
			input: &model.ChangeLog{
				Field:          "attachment",
				AttachmentInfo: &model.AttachmentInfo{ID: 10, Name: "file.txt"},
			},
			want: &ChangeLog{
				Field: "attachment",
				AttachmentInfo: &struct {
					ID   int
					Name string
				}{ID: 10, Name: "file.txt"},
			},
		},
		"with_attribute_info": {
			input: &model.ChangeLog{
				Field:         "attribute",
				AttributeInfo: &model.AttributeInfo{ID: 5, TypeID: "text"},
			},
			want: &ChangeLog{
				Field: "attribute",
				AttributeInfo: &struct {
					ID     int
					TypeID string
				}{ID: 5, TypeID: "text"},
			},
		},
		"with_notification_info": {
			input: &model.ChangeLog{
				Field:            "notification",
				NotificationInfo: &model.NotificationInfo{Type: "issueAssigned"},
			},
			want: &ChangeLog{
				Field: "notification",
				NotificationInfo: &struct {
					Type string
				}{Type: "issueAssigned"},
			},
		},
		"with_all_info_fields": {
			input: &model.ChangeLog{
				Field:            "status",
				NewValue:         "4",
				OriginalValue:    "1",
				AttachmentInfo:   &model.AttachmentInfo{ID: 10, Name: "file.txt"},
				AttributeInfo:    &model.AttributeInfo{ID: 5, TypeID: "text"},
				NotificationInfo: &model.NotificationInfo{Type: "issueAssigned"},
			},
			want: &ChangeLog{
				Field:         "status",
				NewValue:      "4",
				OriginalValue: "1",
				AttachmentInfo: &struct {
					ID   int
					Name string
				}{ID: 10, Name: "file.txt"},
				AttributeInfo: &struct {
					ID     int
					TypeID string
				}{ID: 5, TypeID: "text"},
				NotificationInfo: &struct {
					Type string
				}{Type: "issueAssigned"},
			},
		},
		"nil_info_fields_remain_nil": {
			input: &model.ChangeLog{
				Field:            "status",
				AttachmentInfo:   nil,
				AttributeInfo:    nil,
				NotificationInfo: nil,
			},
			want: &ChangeLog{
				Field:            "status",
				AttachmentInfo:   nil,
				AttributeInfo:    nil,
				NotificationInfo: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, changeLogFromModel(tc.input))
		})
	}
}

func Test_customFieldFromModel(t *testing.T) {
	minNumber, maxNumber, initialValue := 0.0, 100.0, 8.0
	minDate := time.Date(2024, 1, 14, 15, 0, 0, 0, time.UTC)
	maxDate := time.Date(2024, 12, 30, 15, 0, 0, 0, time.UTC)
	initialDate := time.Date(2024, 3, 31, 15, 0, 0, 0, time.UTC)

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
		"number_type": {
			input: &model.CustomField{
				ID:           3,
				TypeID:       3,
				Min:          &model.CustomFieldBound{Number: &minNumber},
				Max:          &model.CustomFieldBound{Number: &maxNumber},
				InitialValue: &initialValue,
				Unit:         "h",
			},
			want: &CustomField{
				ID:           3,
				TypeID:       3,
				Items:        []*CustomFieldItem{},
				Min:          &minNumber,
				Max:          &maxNumber,
				InitialValue: &initialValue,
				Unit:         "h",
			},
		},
		"date_type": {
			input: &model.CustomField{
				ID:          4,
				TypeID:      4,
				Min:         &model.CustomFieldBound{Date: &minDate},
				Max:         &model.CustomFieldBound{Date: &maxDate},
				InitialDate: &model.CustomFieldInitialDate{ID: 3, Date: initialDate},
			},
			want: &CustomField{
				ID:          4,
				TypeID:      4,
				Items:       []*CustomFieldItem{},
				MinDate:     Timestamp{minDate},
				MaxDate:     Timestamp{maxDate},
				InitialDate: &CustomFieldInitialDate{InitialValueType: 3, Date: Timestamp{initialDate}},
			},
		},
		"list_type_allow_input": {
			input: &model.CustomField{
				ID:         5,
				TypeID:     7,
				AllowInput: true,
			},
			want: &CustomField{
				ID:         5,
				TypeID:     7,
				Items:      []*CustomFieldItem{},
				AllowInput: true,
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
	startDate := time.Date(2023, 12, 31, 15, 0, 0, 0, time.UTC)
	dueDate := time.Date(2024, 5, 31, 15, 0, 0, 0, time.UTC)

	user := &model.User{ID: 1, UserID: "admin", Name: "admin", RoleType: 1, Lang: "ja", MailAddress: "admin@example.com"}
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
				Resolution:  &model.Resolution{ID: 0, Name: "Fixed"},
				Priority:    &model.Priority{ID: 3, Name: "Normal"},
				Status:      &model.Status{ID: 1, ProjectID: 10, Name: "Open", Color: "#ed8077", DisplayOrder: 1000},
				Assignee:    user,
				Category: []*model.Category{
					{ID: 5, Name: "Frontend", DisplayOrder: 0},
				},
				Versions: []*model.Version{
					{ID: 30, Name: "v1.0"},
				},
				Milestone: []*model.Version{
					{ID: 31, Name: "v1.1"},
				},
				StartDate:      startDate,
				DueDate:        dueDate,
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
				ChildIssueSummary: &model.ChildIssueSummary{Total: 3, Closed: 1},
			},
			want: &Issue{
				ID:          1,
				ProjectID:   10,
				IssueKey:    "BLG-1",
				KeyID:       1,
				IssueType:   &IssueType{ID: 2, ProjectID: 10, Name: "Task", Color: "#7ea800", DisplayOrder: 0},
				Summary:     "test issue",
				Description: "desc",
				Resolution:  &Resolution{ID: 0, Name: "Fixed"},
				Priority:    &Priority{ID: 3, Name: "Normal"},
				Status:      &Status{ID: 1, ProjectID: 10, Name: "Open", Color: "#ed8077", DisplayOrder: 1000},
				Assignee:    wantUser,
				Category: []*Category{
					{ID: 5, Name: "Frontend", DisplayOrder: 0},
				},
				Versions: []*Version{
					{ID: 30, Name: "v1.0"},
				},
				Milestone: []*Version{
					{ID: 31, Name: "v1.1"},
				},
				StartDate:      Timestamp{startDate},
				DueDate:        Timestamp{dueDate},
				EstimatedHours: 8,
				ActualHours:    4,
				ParentIssueID:  0,
				CreatedUser:    wantUser,
				Created:        Timestamp{created},
				UpdatedUser:    wantUser,
				Updated:        Timestamp{updated},
				CustomFields: []*CustomField{
					{ID: 1, TypeID: 6, Name: "OS", Items: []*CustomFieldItem{
						{ID: 1, Name: "Windows", DisplayOrder: 0},
					}},
				},
				Attachments: []*Attachment{
					{ID: 10, Name: "file.txt", Size: 100, CreatedUser: wantUser, Created: Timestamp{created}},
				},
				SharedFiles: []*SharedFile{
					{ID: 20, Type: "file", Dir: "/", Name: "shared.txt", Size: 200, CreatedUser: wantUser, Created: Timestamp{created}},
				},
				Stars: []*Star{
					{ID: 75, Comment: "good", URL: "https://example.com", Title: "title", Presenter: wantUser, Created: Timestamp{created}},
				},
				ChildIssueSummary: &ChildIssueSummary{Total: 3, Closed: 1},
			},
		},
		"nil_elements": {
			input: &model.Issue{
				Category:     []*model.Category{nil},
				Versions:     []*model.Version{nil},
				Milestone:    []*model.Version{nil},
				CustomFields: []*model.CustomField{nil},
				Attachments:  []*model.Attachment{nil},
				SharedFiles:  []*model.SharedFile{nil},
				Stars:        []*model.Star{nil},
			},
			want: &Issue{
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

func Test_notificationFromModel(t *testing.T) {
	t.Parallel()

	created := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	user := &model.User{ID: 1, UserID: "admin", Name: "admin", RoleType: 1, Lang: "ja", MailAddress: "admin@example.com"}
	wantUser := &User{ID: 1, UserID: "admin", Name: "admin", RoleType: RoleAdministrator, Lang: "ja", MailAddress: "admin@example.com"}

	input := &model.Notification{
		ID:                  10,
		AlreadyRead:         false,
		Reason:              2,
		ResourceAlreadyRead: true,
		Sender:              user,
		User:                user,
		Created:             created,
	}
	want := &Notification{
		ID:                  10,
		AlreadyRead:         false,
		Reason:              2,
		ResourceAlreadyRead: true,
		Sender:              wantUser,
		User:                wantUser,
		Created:             Timestamp{created},
	}
	assert.Equal(t, want, notificationFromModel(input))
}

func Test_pullRequestFromModel(t *testing.T) {
	t.Parallel()

	created := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	closeAt := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	mergeAt := time.Date(2024, 5, 2, 0, 0, 0, 0, time.UTC)

	user := &model.User{ID: 1, UserID: "admin", Name: "admin", RoleType: 1, Lang: "ja", MailAddress: "admin@example.com"}
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
		Status:       &model.Status{ID: 1, ProjectID: 3, Name: "Open", Color: "#ed8077", DisplayOrder: 1000},
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
		Status:       &Status{ID: 1, ProjectID: 3, Name: "Open", Color: "#ed8077", DisplayOrder: 1000},
		Assignee:     wantUser,
		Issue: &Issue{
			ID:      10,
			Summary: "related issue",
		},
		BaseCommit:   "abc123",
		BranchCommit: "def456",
		CloseAt:      Timestamp{closeAt},
		MergeAt:      Timestamp{mergeAt},
		CreatedUser:  wantUser,
		Created:      Timestamp{created},
		UpdatedUser:  wantUser,
		Updated:      Timestamp{updated},
		Attachments: []*Attachment{
			{ID: 10, Name: "file.txt", Size: 100, CreatedUser: wantUser, Created: Timestamp{created}},
		},
		Stars: []*Star{
			{ID: 75, Comment: "good", URL: "https://example.com", Title: "title", Presenter: wantUser, Created: Timestamp{created}},
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
		TextFormattingRule: "markdown",
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
		Created:            Timestamp{created},
		Updated:            Timestamp{updated},
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
		Updated: Timestamp{updated},
	}
	assert.Equal(t, want, spaceNotificationFromModel(input))
}

func Test_statusFromModel(t *testing.T) {
	t.Parallel()

	input := &model.Status{ID: 1, ProjectID: 6, Name: "Open", Color: "#ed8077", DisplayOrder: 1000}
	want := &Status{ID: 1, ProjectID: 6, Name: "Open", Color: "#ed8077", DisplayOrder: 1000}
	assert.Equal(t, want, statusFromModel(input))
}

func Test_versionFromModel(t *testing.T) {
	t.Parallel()

	startDate := time.Date(2023, 12, 31, 15, 0, 0, 0, time.UTC)
	releaseDueDate := time.Date(2024, 3, 30, 15, 0, 0, 0, time.UTC)

	input := &model.Version{
		ID:             30,
		ProjectID:      1,
		Name:           "v1.0",
		Description:    "first release",
		StartDate:      startDate,
		ReleaseDueDate: releaseDueDate,
		Archived:       false,
		DisplayOrder:   0,
	}
	want := &Version{
		ID:             30,
		ProjectID:      1,
		Name:           "v1.0",
		Description:    "first release",
		StartDate:      Timestamp{startDate},
		ReleaseDueDate: Timestamp{releaseDueDate},
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
		"activity":                 {call: func() any { return activityFromModel(nil) }},
		"activity_content":         {call: func() any { return activityContentFromModel(nil) }},
		"attachment":               {call: func() any { return attachmentFromModel(nil) }},
		"change_log":               {call: func() any { return changeLogFromModel(nil) }},
		"child_issue_summary":      {call: func() any { return childIssueSummaryFromModel(nil) }},
		"comment":                  {call: func() any { return commentFromModel(nil) }},
		"custom_field":             {call: func() any { return customFieldFromModel(nil) }},
		"custom_field_item":        {call: func() any { return customFieldItemFromModel(nil) }},
		"custom_field_init_date":   {call: func() any { return customFieldInitialDateFromModel(nil) }},
		"disk_usage_project":       {call: func() any { return diskUsageProjectFromModel(nil) }},
		"disk_usage_space":         {call: func() any { return diskUsageSpaceFromModel(nil) }},
		"issue":                    {call: func() any { return issueFromModel(nil) }},
		"notification":             {call: func() any { return notificationFromModel(nil) }},
		"project":                  {call: func() any { return projectFromModel(nil) }},
		"pull_request":             {call: func() any { return pullRequestFromModel(nil) }},
		"recently_viewed_issue":    {call: func() any { return recentlyViewedIssueFromModel(nil) }},
		"recently_viewed_issues":   {call: func() any { return recentlyViewedIssuesFromModel(nil) }},
		"recently_viewed_project":  {call: func() any { return recentlyViewedProjectFromModel(nil) }},
		"recently_viewed_projects": {call: func() any { return recentlyViewedProjectsFromModel(nil) }},
		"recently_viewed_wiki":     {call: func() any { return recentlyViewedWikiFromModel(nil) }},
		"recently_viewed_wikis":    {call: func() any { return recentlyViewedWikisFromModel(nil) }},
		"resolution":               {call: func() any { return resolutionFromModel(nil) }},
		"shared_file":              {call: func() any { return sharedFileFromModel(nil) }},
		"space":                    {call: func() any { return spaceFromModel(nil) }},
		"space_notification":       {call: func() any { return spaceNotificationFromModel(nil) }},
		"star":                     {call: func() any { return starFromModel(nil) }},
		"stars":                    {call: func() any { return starsFromModel(nil) }},
		"status":                   {call: func() any { return statusFromModel(nil) }},
		"tag":                      {call: func() any { return tagFromModel(nil) }},
		"user":                     {call: func() any { return userFromModel(nil) }},
		"version":                  {call: func() any { return versionFromModel(nil) }},
		"versions":                 {call: func() any { return versionsFromModel(nil) }},
		"wiki":                     {call: func() any { return wikiFromModel(nil) }},
		"wiki_history":             {call: func() any { return wikiHistoryFromModel(nil) }},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Nil(t, tc.call())
		})
	}
}
