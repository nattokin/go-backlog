package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// IssueAttachmentService
	doerIssueAttachmentList   = newMockDoer(fixture.Attachment.ListJSON)
	doerIssueAttachmentRemove = newMockDoer(fixture.Attachment.SingleJSON)
)

func ExampleIssueAttachmentService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueAttachmentList),
	)

	attachments, _ := c.Issue.Attachment.List(context.Background(), "TEST-1")
	fmt.Printf("ID: %d, Name: %s\n", attachments[0].ID, attachments[0].Name)
	// Output:
	// ID: 2, Name: A.png
}

func ExampleIssueAttachmentService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueAttachmentRemove),
	)

	attachment, _ := c.Issue.Attachment.Remove(context.Background(), "TEST-1", 8)
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 8, Name: IMG0088.png
}
