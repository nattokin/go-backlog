package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// PullRequestService
	doerPullRequestAll    = newMockDoer(fixture.PullRequest.ListJSON)
	doerPullRequestCount  = newMockDoer(`{"count":2}`)
	doerPullRequestOne    = newMockDoer(fixture.PullRequest.SingleJSON)
	doerPullRequestCreate = newMockDoer(fixture.PullRequest.SingleJSON)
	doerPullRequestUpdate = newMockDoer(fixture.PullRequest.SingleJSON)

	// PullRequestAttachmentService
	doerPullRequestAttachmentList   = newMockDoer(fixture.Attachment.ListJSON)
	doerPullRequestAttachmentRemove = newMockDoer(fixture.Attachment.SingleJSON)
)

func ExamplePullRequestService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestAll),
	)

	prs, _ := c.PullRequest.All(context.Background(), "TEST", "myrepo")
	fmt.Printf("Count: %d, ID: %d, Summary: %s\n", len(prs), prs[0].ID, prs[0].Summary)
	// Output:
	// Count: 2, ID: 2, Summary: test PR
}

func ExamplePullRequestService_Count() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestCount),
	)

	count, _ := c.PullRequest.Count(context.Background(), "TEST", "myrepo")
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 2
}

func ExamplePullRequestService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestOne),
	)

	pr, _ := c.PullRequest.One(context.Background(), "TEST", "myrepo", 1)
	fmt.Printf("ID: %d, Summary: %s\n", pr.ID, pr.Summary)
	// Output:
	// ID: 2, Summary: test PR
}

func ExamplePullRequestService_Create() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestCreate),
	)

	pr, _ := c.PullRequest.Create(context.Background(), "TEST", "myrepo", "test PR", "test description", "main", "feature/foo")
	fmt.Printf("ID: %d, Summary: %s\n", pr.ID, pr.Summary)
	// Output:
	// ID: 2, Summary: test PR
}

func ExamplePullRequestService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestUpdate),
	)

	pr, _ := c.PullRequest.Update(context.Background(), "TEST", "myrepo", 1, c.PullRequest.Option.WithSummary("test PR"))
	fmt.Printf("ID: %d, Summary: %s\n", pr.ID, pr.Summary)
	// Output:
	// ID: 2, Summary: test PR
}

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
