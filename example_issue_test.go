package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// IssueService
	doerIssueAll    = newMockDoer(fixture.Issue.ListJSON)
	doerIssueCount  = newMockDoer(`{"count":2}`)
	doerIssueOne    = newMockDoer(fixture.Issue.SingleJSON)
	doerIssueCreate = newMockDoer(fixture.Issue.SingleJSON)
	doerIssueUpdate = newMockDoer(fixture.Issue.SingleJSON)
	doerIssueDelete = newMockDoer(fixture.Issue.SingleJSON)

	// IssueAttachmentService
	doerIssueAttachmentList   = newMockDoer(fixture.Attachment.ListJSON)
	doerIssueAttachmentRemove = newMockDoer(fixture.Attachment.SingleJSON)
)

func ExampleIssueService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueAll),
	)

	issues, _ := c.Issue.All(context.Background())
	fmt.Printf("Count: %d, ID: %d, Summary: %s\n", len(issues), issues[0].ID, issues[0].Summary)
	// Output:
	// Count: 2, ID: 1, Summary: First issue
}

func ExampleIssueService_Count() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCount),
	)

	count, _ := c.Issue.Count(context.Background())
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 2
}

func ExampleIssueService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueOne),
	)

	issue, _ := c.Issue.One(context.Background(), "PRJ-1")
	fmt.Printf("ID: %d, IssueKey: %s, Summary: %s\n", issue.ID, issue.IssueKey, issue.Summary)
	// Output:
	// ID: 1, IssueKey: PRJ-1, Summary: First issue
}

func ExampleIssueService_Create() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCreate),
	)

	issue, _ := c.Issue.Create(context.Background(), 10, "First issue", 2, 3)
	fmt.Printf("ID: %d, Summary: %s\n", issue.ID, issue.Summary)
	// Output:
	// ID: 1, Summary: First issue
}

func ExampleIssueService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueUpdate),
	)

	issue, _ := c.Issue.Update(context.Background(), "PRJ-1", c.Issue.Option.WithSummary("First issue"))
	fmt.Printf("ID: %d, Summary: %s\n", issue.ID, issue.Summary)
	// Output:
	// ID: 1, Summary: First issue
}

func ExampleIssueService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueDelete),
	)

	issue, _ := c.Issue.Delete(context.Background(), "PRJ-1")
	fmt.Printf("ID: %d, IssueKey: %s\n", issue.ID, issue.IssueKey)
	// Output:
	// ID: 1, IssueKey: PRJ-1
}

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
