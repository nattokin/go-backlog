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

func Test_versionsFromModel(t *testing.T) {
	t.Parallel()

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
		"nil": {
			input: nil,
			want:  nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, versionsFromModel(tc.input))
		})
	}
}

func Test_issueFromModel(t *testing.T) {
	t.Parallel()

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
			// Resolutions と Category のforループ内 nil チェックをカバー
			input: &model.Issue{
				Resolutions: []*model.Resolution{nil},
				Category:    []*model.Category{nil},
			},
			want: &Issue{
				Resolutions:  []*Resolution{nil},
				Category:     []*Category{nil},
				CustomFields: []*CustomField{},
				SharedFiles:  []*SharedFile{},
				Stars:        []*Star{},
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

	cases := map[string]struct {
		input *model.PullRequest
		want  *PullRequest
	}{
		"full": {
			input: &model.PullRequest{
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
			},
			want: &PullRequest{
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
					ID:           10,
					Summary:      "related issue",
					Resolutions:  []*Resolution{},
					Category:     []*Category{},
					CustomFields: []*CustomField{},
					SharedFiles:  []*SharedFile{},
					Stars:        []*Star{},
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
			assert.Equal(t, tc.want, pullRequestFromModel(tc.input))
		})
	}
}

func Test_starFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, starFromModel(nil))
}

func Test_starsFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, starsFromModel(nil))
}

func Test_attachmentFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, attachmentFromModel(nil))
}

func Test_commentFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, commentFromModel(nil))
}

func Test_notificationFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, notificationFromModel(nil))
}

func Test_sharedFileFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, sharedFileFromModel(nil))
}

func Test_tagFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, tagFromModel(nil))
}

func Test_activityContentFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, activityContentFromModel(nil))
}

func Test_activityFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, activityFromModel(nil))
}

func Test_userFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, userFromModel(nil))
}

func Test_projectFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, projectFromModel(nil))
}

func Test_wikiFromModel_nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, wikiFromModel(nil))
}
