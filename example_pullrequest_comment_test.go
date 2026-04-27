package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// PullRequestCommentService
	doerPullRequestCommentAll    = newMockDoer(fixture.Comment.ListJSON)
	doerPullRequestCommentAdd    = newMockDoer(fixture.Comment.SingleJSON)
	doerPullRequestCommentCount  = newMockDoer(`{"count":2}`)
	doerPullRequestCommentUpdate = newMockDoer(fixture.Comment.SingleJSON)
)

func ExamplePullRequestCommentService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestCommentAll),
	)

	comments, _ := c.PullRequest.Comment.All(context.Background(), "TEST", "myrepo", 1)
	fmt.Printf("Count: %d, ID: %d, Content: %s\n", len(comments), comments[0].ID, comments[0].Content)
	// Output:
	// Count: 2, ID: 1, Content: This is a comment.
}

func ExamplePullRequestCommentService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestCommentAdd),
	)

	comment, _ := c.PullRequest.Comment.Add(context.Background(), "TEST", "myrepo", 1, "This is a comment.")
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExamplePullRequestCommentService_Count() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestCommentCount),
	)

	count, _ := c.PullRequest.Comment.Count(context.Background(), "TEST", "myrepo", 1)
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 2
}

func ExamplePullRequestCommentService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestCommentUpdate),
	)

	comment, _ := c.PullRequest.Comment.Update(context.Background(), "TEST", "myrepo", 1, 1, "This is a comment.")
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}
