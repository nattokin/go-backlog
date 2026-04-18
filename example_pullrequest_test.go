package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// PullRequestAttachmentService
	doerPullRequestAttachmentList   = newMockDoer(fixture.Attachment.ListJSON)
	doerPullRequestAttachmentRemove = newMockDoer(fixture.Attachment.SingleJSON)
)

func ExamplePullRequestAttachmentService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestAttachmentList),
	)

	attachments, _ := c.PullRequest.Attachment.List(context.Background(), "TEST", "myrepo", 1)
	fmt.Printf("ID: %d, Name: %s\n", attachments[0].ID, attachments[0].Name)
	// Output:
	// ID: 2, Name: A.png
}

func ExamplePullRequestAttachmentService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestAttachmentRemove),
	)

	attachment, _ := c.PullRequest.Attachment.Remove(context.Background(), "TEST", "myrepo", 1, 8)
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 8, Name: IMG0088.png
}
