package backlog

import (
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/pullrequest"
	"github.com/nattokin/go-backlog/internal/space"
	"github.com/nattokin/go-backlog/internal/wiki"
)

type IssueAttachmentService = attachment.IssueAttachmentService

type IssueService = issue.IssueService

type PullRequestAttachmentService = attachment.PullRequestAttachmentService

type PullRequestService = pullrequest.PullRequestService

type SpaceAttachmentService = attachment.SpaceAttachmentService

type SpaceService = space.SpaceService

type WikiAttachmentService = attachment.WikiAttachmentService

type WikiService = wiki.WikiService
