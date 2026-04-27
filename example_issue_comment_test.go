package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// IssueCommentService
	doerIssueCommentAll           = newMockDoer(fixture.Comment.ListJSON)
	doerIssueCommentAdd           = newMockDoer(fixture.Comment.SingleJSON)
	doerIssueCommentCount         = newMockDoer(`{"count":2}`)
	doerIssueCommentDelete        = newMockDoer(fixture.Comment.SingleJSON)
	doerIssueCommentNotifications = newMockDoer(`[{"id":25,"alreadyRead":false,"reason":2,"resourceAlreadyRead":false}]`)
	doerIssueCommentNotify        = newMockDoer(fixture.Comment.SingleJSON)
	doerIssueCommentOne           = newMockDoer(fixture.Comment.SingleJSON)
	doerIssueCommentUpdate        = newMockDoer(fixture.Comment.SingleJSON)
)

func ExampleIssueCommentService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentAll),
	)

	comments, _ := c.Issue.Comment.All(context.Background(), "PRJ-1")
	fmt.Printf("Count: %d, ID: %d, Content: %s\n", len(comments), comments[0].ID, comments[0].Content)
	// Output:
	// Count: 2, ID: 1, Content: This is a comment.
}

func ExampleIssueCommentService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentAdd),
	)

	comment, _ := c.Issue.Comment.Add(context.Background(), "PRJ-1", "This is a comment.")
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExampleIssueCommentService_Count() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentCount),
	)

	count, _ := c.Issue.Comment.Count(context.Background(), "PRJ-1")
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 2
}

func ExampleIssueCommentService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentDelete),
	)

	comment, _ := c.Issue.Comment.Delete(context.Background(), "PRJ-1", 1)
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExampleIssueCommentService_Notifications() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentNotifications),
	)

	notifications, _ := c.Issue.Comment.Notifications(context.Background(), "PRJ-1", 1)
	fmt.Printf("Count: %d, ID: %d\n", len(notifications), notifications[0].ID)
	// Output:
	// Count: 1, ID: 25
}

func ExampleIssueCommentService_Notify() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentNotify),
	)

	comment, _ := c.Issue.Comment.Notify(context.Background(), "PRJ-1", 1, []int{5686})
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExampleIssueCommentService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentOne),
	)

	comment, _ := c.Issue.Comment.One(context.Background(), "PRJ-1", 1)
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExampleIssueCommentService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentUpdate),
	)

	comment, _ := c.Issue.Comment.Update(context.Background(), "PRJ-1", 1, "This is a comment.")
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}
