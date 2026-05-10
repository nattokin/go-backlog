package backlog_test

import (
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// Activity
	doerActivityList   = newMockDoer(fixture.Activity.ListJSON)
	doerActivitySingle = newMockDoer(fixture.Activity.SingleJSON)

	// Attachment
	doerAttachmentList   = newMockDoer(fixture.Attachment.ListJSON)
	doerAttachmentSingle = newMockDoer(fixture.Attachment.SingleJSON)
	doerAttachmentUpload = newMockDoer(fixture.Attachment.UploadJSON)

	// Category
	doerCategoryList   = newMockDoer(fixture.Category.ListJSON)
	doerCategorySingle = newMockDoer(fixture.Category.SingleJSON)

	// Comment
	doerCommentList   = newMockDoer(fixture.Comment.ListJSON)
	doerCommentSingle = newMockDoer(fixture.Comment.SingleJSON)

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
	doerPullRequestList   = newMockDoer(fixture.PullRequest.ListJSON)
	doerPullRequestSingle = newMockDoer(fixture.PullRequest.SingleJSON)

	// RecentlyViewed
	doerRecentlyViewedIssueList    = newMockDoer(fixture.RecentlyViewed.IssueListJSON)
	doerRecentlyViewedIssueSingle  = newMockDoer(fixture.RecentlyViewed.IssueSingleJSON)
	doerRecentlyViewedProjectList  = newMockDoer(fixture.RecentlyViewed.ProjectListJSON)
	doerRecentlyViewedWikiList     = newMockDoer(fixture.RecentlyViewed.WikiListJSON)
	doerRecentlyViewedWikiSingle   = newMockDoer(fixture.RecentlyViewed.WikiSingleJSON)

	// Repository
	doerRepositoryList   = newMockDoer(fixture.Repository.ListJSON)
	doerRepositorySingle = newMockDoer(fixture.Repository.SingleJSON)

	// SharedFile
	doerSharedFileList   = newMockDoer(fixture.SharedFile.ListJSON)
	doerSharedFileSingle = newMockDoer(fixture.SharedFile.SingleJSON)

	// Space
	doerSpaceSpace        = newMockDoer(fixture.Space.SpaceJSON)
	doerSpaceDiskUsage    = newMockDoer(fixture.Space.DiskUsageJSON)
	doerSpaceNotification = newMockDoer(fixture.Space.NotificationJSON)

	// Star
	doerStarList  = newMockDoer(fixture.Star.ListJSON)
	doerStarCount = newMockDoer(fixture.Star.CountJSON)

	// Status
	doerStatusList   = newMockDoer(fixture.Status.ListJSON)
	doerStatusSingle = newMockDoer(fixture.Status.SingleJSON)

	// User
	doerUserList   = newMockDoer(fixture.User.ListJSON)
	doerUserSingle = newMockDoer(fixture.User.SingleJSON)
	doerUserIcon   = newMockBinaryDoer("image/png", "icon.png", []byte("PNG"))

	// Version
	doerVersionList   = newMockDoer(fixture.Version.ListJSON)
	doerVersionSingle = newMockDoer(fixture.Version.SingleJSON)

	// Webhook
	doerWebhookList     = newMockDoer(fixture.Webhook.ListJSON)
	doerWebhookAllEvent = newMockDoer(fixture.Webhook.AllEventJSON)

	// Wiki
	doerWikiList        = newMockDoer(fixture.Wiki.ListJSON)
	doerWikiSingle      = newMockDoer(fixture.Wiki.MinimumJSON)
	doerWikiHistoryList = newMockDoer(fixture.WikiHistory.ListJSON)
	doerWikiAttachmentDownload = newMockBinaryDoer("image/png", "A.png", []byte("PNG"))
	doerSharedFileGetFile      = newMockBinaryDoer("image/png", "shared.png", []byte("PNG"))
	doerAttachmentDownload     = newMockBinaryDoer("image/png", "A.png", []byte("PNG"))
)
