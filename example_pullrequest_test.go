package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExamplePullRequestService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	prs, err := c.PullRequest.List(context.Background(), "TEST", "myrepo")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Count: %d, ID: %d, Summary: %s\n", len(prs), prs[0].ID, prs[0].Summary)
	// Output:
	// Count: 2, ID: 2, Summary: test PR
}

func ExamplePullRequestService_All() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	seq, err := c.PullRequest.All(context.Background(), 100, "TEST", "myrepo")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	var ids []int
	for pr, err := range seq {
		if err != nil {
			break
		}
		ids = append(ids, pr.ID)
	}
	fmt.Printf("Count: %d, IDs: %v\n", len(ids), ids)
	// Output:
	// Count: 2, IDs: [2 3]
}

func ExamplePullRequestService_Count() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	count, err := c.PullRequest.Count(context.Background(), "TEST", "myrepo")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 2
}

func ExamplePullRequestService_One() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	pr, err := c.PullRequest.One(context.Background(), "TEST", "myrepo", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Summary: %s\n", pr.ID, pr.Summary)
	// Output:
	// ID: 2, Summary: test PR
}

func ExamplePullRequestService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	pr, err := c.PullRequest.Create(context.Background(), "TEST", "myrepo", "test PR", "test description", "main", "feature/foo")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Summary: %s\n", pr.ID, pr.Summary)
	// Output:
	// ID: 2, Summary: test PR
}

func ExamplePullRequestService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	pr, err := c.PullRequest.Update(context.Background(), "TEST", "myrepo", 1, c.PullRequest.Option.WithSummary("test PR"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Summary: %s\n", pr.ID, pr.Summary)
	// Output:
	// ID: 2, Summary: test PR
}

func ExamplePullRequestAttachmentService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	attachments, err := c.PullRequest.Attachment.List(context.Background(), "TEST", "myrepo", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", attachments[0].ID, attachments[0].Name)
	// Output:
	// ID: 2, Name: A.png
}

func ExamplePullRequestAttachmentService_Download() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	file, err := c.PullRequest.Attachment.Download(context.Background(), "TEST", "myrepo", 1, 2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ContentType: %s, FileName: %s\n", file.ContentType, file.Filename)
	// Output:
	// ContentType: image/png, FileName: A.png
}

func ExamplePullRequestAttachmentService_Remove() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	attachment, err := c.PullRequest.Attachment.Remove(context.Background(), "TEST", "myrepo", 1, 8)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 8, Name: IMG0088.png
}

func ExamplePullRequestCommentService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	comments, err := c.PullRequest.Comment.List(context.Background(), "TEST", "myrepo", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Count: %d, ID: %d, Content: %s\n", len(comments), comments[0].ID, comments[0].Content)
	// Output:
	// Count: 2, ID: 1, Content: This is a comment.
}

func ExamplePullRequestCommentService_Add() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	comment, err := c.PullRequest.Comment.Add(context.Background(), "TEST", "myrepo", 1, "This is a comment.")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExamplePullRequestCommentService_Count() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	count, err := c.PullRequest.Comment.Count(context.Background(), "TEST", "myrepo", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 2
}

func ExamplePullRequestCommentService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	comment, err := c.PullRequest.Comment.Update(context.Background(), "TEST", "myrepo", 1, 1, "This is a comment.")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExamplePullRequestStarService_Add() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = c.PullRequest.Star.Add(context.Background(), 2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}

func ExamplePullRequestStarService_Remove() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = c.PullRequest.Star.Remove(context.Background(), 42)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}
