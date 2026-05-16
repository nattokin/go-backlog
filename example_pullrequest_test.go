package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExamplePullRequestService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestList),
	)

	prs, _ := c.PullRequest.List(context.Background(), "TEST", "myrepo")
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
		backlog.WithDoer(doerPullRequestSingle),
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
		backlog.WithDoer(doerPullRequestSingle),
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
		backlog.WithDoer(doerPullRequestSingle),
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
		backlog.WithDoer(doerAttachmentList),
	)

	attachments, _ := c.PullRequest.Attachment.List(context.Background(), "TEST", "myrepo", 1)
	fmt.Printf("ID: %d, Name: %s\n", attachments[0].ID, attachments[0].Name)
	// Output:
	// ID: 2, Name: A.png
}

func ExamplePullRequestAttachmentService_Download() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerAttachmentDownload),
	)

	file, _ := c.PullRequest.Attachment.Download(context.Background(), "TEST", "myrepo", 1, 2)
	fmt.Printf("ContentType: %s, FileName: %s\n", file.ContentType, file.Filename)
	// Output:
	// ContentType: image/png, FileName: A.png
}

func ExamplePullRequestAttachmentService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerAttachmentSingle),
	)

	attachment, _ := c.PullRequest.Attachment.Remove(context.Background(), "TEST", "myrepo", 1, 8)
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 8, Name: IMG0088.png
}

func ExamplePullRequestCommentService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerCommentList),
	)

	comments, _ := c.PullRequest.Comment.List(context.Background(), "TEST", "myrepo", 1)
	fmt.Printf("Count: %d, ID: %d, Content: %s\n", len(comments), comments[0].ID, comments[0].Content)
	// Output:
	// Count: 2, ID: 1, Content: This is a comment.
}

func ExamplePullRequestCommentService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerCommentSingle),
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
		backlog.WithDoer(doerCommentSingle),
	)

	comment, _ := c.PullRequest.Comment.Update(context.Background(), "TEST", "myrepo", 1, 1, "This is a comment.")
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExamplePullRequestStarService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerNoContent),
	)

	err := c.PullRequest.Star.Add(context.Background(), 2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}

func ExamplePullRequestStarService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerNoContent),
	)

	err := c.PullRequest.Star.Remove(context.Background(), 42)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}
