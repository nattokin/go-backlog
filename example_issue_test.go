package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleIssueService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueList),
	)

	issues, _ := c.Issue.List(context.Background())
	fmt.Printf("Count: %d, ID: %d, Summary: %s\n", len(issues), issues[0].ID, issues[0].Summary)
	// Output:
	// Count: 2, ID: 1, Summary: First issue
}

func ExampleIssueService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueList),
	)

	seq, err := c.Issue.All(context.Background(), 100)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	var ids []int
	for issue, err := range seq {
		if err != nil {
			break
		}
		ids = append(ids, issue.ID)
	}
	fmt.Printf("Count: %d, IDs: %v\n", len(ids), ids)
	// Output:
	// Count: 2, IDs: [1 2]
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
		backlog.WithDoer(doerIssueSingle),
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
		backlog.WithDoer(doerIssueSingle),
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
		backlog.WithDoer(doerIssueSingle),
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
		backlog.WithDoer(doerIssueSingle),
	)

	issue, _ := c.Issue.Delete(context.Background(), "PRJ-1")
	fmt.Printf("ID: %d, IssueKey: %s\n", issue.ID, issue.IssueKey)
	// Output:
	// ID: 1, IssueKey: PRJ-1
}

func ExampleIssueService_Participants() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserList),
	)

	users, _ := c.Issue.Participants(context.Background(), "PRJ-1")
	fmt.Printf("Count: %d, ID: %d, Name: %s\n", len(users), users[0].ID, users[0].Name)
	// Output:
	// Count: 4, ID: 1, Name: admin
}

func ExampleIssueAttachmentService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerAttachmentList),
	)

	attachments, _ := c.Issue.Attachment.List(context.Background(), "TEST-1")
	fmt.Printf("ID: %d, Name: %s\n", attachments[0].ID, attachments[0].Name)
	// Output:
	// ID: 2, Name: A.png
}

func ExampleIssueAttachmentService_Download() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerAttachmentDownload),
	)

	file, _ := c.Issue.Attachment.Download(context.Background(), "TEST-1", 2)
	fmt.Printf("ContentType: %s, FileName: %s\n", file.ContentType, file.Filename)
	// Output:
	// ContentType: image/png, FileName: A.png
}

func ExampleIssueAttachmentService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerAttachmentSingle),
	)

	attachment, _ := c.Issue.Attachment.Remove(context.Background(), "TEST-1", 8)
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 8, Name: IMG0088.png
}

func ExampleIssueCommentService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerCommentList),
	)

	comments, _ := c.Issue.Comment.List(context.Background(), "PRJ-1")
	fmt.Printf("Count: %d, ID: %d, Content: %s\n", len(comments), comments[0].ID, comments[0].Content)
	// Output:
	// Count: 2, ID: 1, Content: This is a comment.
}

func ExampleIssueCommentService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerCommentSingle),
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
		backlog.WithDoer(doerCommentSingle),
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
		backlog.WithDoer(doerCommentSingle),
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
		backlog.WithDoer(doerCommentSingle),
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
		backlog.WithDoer(doerCommentSingle),
	)

	comment, _ := c.Issue.Comment.Update(context.Background(), "PRJ-1", 1, "This is a comment.")
	fmt.Printf("ID: %d, Content: %s\n", comment.ID, comment.Content)
	// Output:
	// ID: 1, Content: This is a comment.
}

func ExampleIssueSharedFileService_Link() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSharedFileList),
	)

	files, _ := c.Issue.SharedFile.Link(context.Background(), "TEST-1", []int{454403, 454404})
	fmt.Printf("ID: %d, Name: %s\n", files[0].ID, files[0].Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}

func ExampleIssueSharedFileService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSharedFileList),
	)

	files, _ := c.Issue.SharedFile.List(context.Background(), "TEST-1")
	fmt.Printf("ID: %d, Name: %s\n", files[0].ID, files[0].Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}

func ExampleIssueSharedFileService_Unlink() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSharedFileSingle),
	)

	file, _ := c.Issue.SharedFile.Unlink(context.Background(), "TEST-1", 454403)
	fmt.Printf("ID: %d, Name: %s\n", file.ID, file.Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}

func ExampleIssueStarService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerNoContent),
	)

	err := c.Issue.Star.Add(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}

func ExampleIssueStarService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerNoContent),
	)

	err := c.Issue.Star.Remove(context.Background(), 42)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}
