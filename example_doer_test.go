package backlog_test

import (
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

var (
	// Activity
	doerActivityList = newMockDoer(fixture.Activity.ListJSON)

	// Category
	doerCategoryList   = newMockDoer(fixture.Category.ListJSON)
	doerCategorySingle = newMockDoer(fixture.Category.SingleJSON)

	// Comment
	doerCommentList = newMockDoer(fixture.Comment.ListJSON)

	// CustomField
	doerCustomFieldList   = newMockDoer(fixture.CustomField.ListJSON)
	doerCustomFieldSingle = newMockDoer(fixture.CustomField.SingleJSON)

	// Issue
	doerIssueList   = newMockDoer(fixture.Issue.ListJSON)
	doerIssueSingle = newMockDoer(fixture.Issue.SingleJSON)

	// IssueType
	doerIssueTypeList   = newMockDoer(fixture.IssueType.ListJSON)
	doerIssueTypeSingle = newMockDoer(fixture.IssueType.SingleJSON)

	// Project
	doerProjectList      = newMockDoer(fixture.Project.ListJSON)
	doerProjectSingle    = newMockDoer(fixture.Project.SingleJSON)
	doerProjectDiskUsage = newMockDoer(fixture.Project.DiskUsageJSON)
	doerProjectIcon      = newMockBinaryDoer("image/png", "test.png", []byte("PNG"))

	// PullRequest
	doerPullRequestList = newMockDoer(fixture.PullRequest.ListJSON)

	// RecentlyViewed
	doerRecentlyViewedIssueList = newMockDoer(fixture.RecentlyViewed.IssueListJSON)

	// SharedFile
	doerSharedFileList = newMockDoer(fixture.SharedFile.ListJSON)

	// Star
	doerStarList = newMockDoer(fixture.Star.ListJSON)

	// Status
	doerStatusList   = newMockDoer(fixture.Status.ListJSON)
	doerStatusSingle = newMockDoer(fixture.Status.SingleJSON)

	// User
	doerUserList   = newMockDoer(fixture.User.ListJSON)
	doerUserSingle = newMockDoer(fixture.User.SingleJSON)

	// Version
	doerVersionList   = newMockDoer(fixture.Version.ListJSON)
	doerVersionSingle = newMockDoer(fixture.Version.SingleJSON)

	// Webhook
	doerWebhookList     = newMockDoer(fixture.Webhook.ListJSON)
	doerWebhookAllEvent = newMockDoer(fixture.Webhook.AllEventJSON)

	// Wiki
	doerWikiSingle            = newMockDoer(fixture.Wiki.MinimumJSON)
	doerSharedFileGetDownload = newMockBinaryDoer("image/png", "shared.png", []byte("PNG"))
)

// newMockDoer returns a *mock.Doer that always responds with HTTP 200 and the given JSON body.
// Used as a lightweight Doer for Example tests, which run without *testing.T.
func newMockDoer(json string) *mock.Doer {
	return &mock.Doer{
		DoFunc: mock.NewDoFunc(json),
	}
}

// newMockBinaryDoer returns a *mock.Doer that always responds with HTTP 200
// and a binary file download response with the given Content-Type, filename, and body.
// Used as a lightweight Doer for Example tests, which run without *testing.T.
func newMockBinaryDoer(contentType, filename string, body []byte) *mock.Doer {
	return &mock.Doer{
		DoFunc: mock.NewBinaryDoFunc(filename, contentType, body),
	}
}

// doerNoContent is a *mock.Doer that always responds with HTTP 204 No Content.
// Used as a lightweight Doer for Example tests, which run without *testing.T.
var doerNoContent = &mock.Doer{
	DoFunc: mock.NewNoContentDoFunc(),
}
