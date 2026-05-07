package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// IssueService
	doerIssueAll          = newMockDoer(fixture.Issue.ListJSON)
	doerIssueCount        = newMockDoer(`{"count":2}`)
	doerIssueOne          = newMockDoer(fixture.Issue.SingleJSON)
	doerIssueCreate       = newMockDoer(fixture.Issue.SingleJSON)
	doerIssueUpdate       = newMockDoer(fixture.Issue.SingleJSON)
	doerIssueDelete       = newMockDoer(fixture.Issue.SingleJSON)
	doerIssueParticipants = newMockDoer(fixture.User.ListJSON)

	// IssueAttachmentService
	doerIssueAttachmentList     = newMockDoer(fixture.Attachment.ListJSON)
	doerIssueAttachmentDownload = newMockBinaryDoer("image/png", "A.png", []byte("PNG"))
	doerIssueAttachmentRemove   = newMockDoer(fixture.Attachment.SingleJSON)

	// IssueCommentService
	doerIssueCommentAll           = newMockDoer(fixture.Comment.ListJSON)
	doerIssueCommentAdd           = newMockDoer(fixture.Comment.SingleJSON)
	doerIssueCommentCount         = newMockDoer(`{"count":2}`)
	doerIssueCommentDelete        = newMockDoer(fixture.Comment.SingleJSON)
	doerIssueCommentNotifications = newMockDoer(`[{"id":25,"alreadyRead":false,"reason":2,"resourceAlreadyRead":false}]`)
	doerIssueCommentNotify        = newMockDoer(fixture.Comment.SingleJSON)
	doerIssueCommentOne           = newMockDoer(fixture.Comment.SingleJSON)
	doerIssueCommentUpdate        = newMockDoer(fixture.Comment.SingleJSON)

	// IssueSharedFileService
	doerIssueSharedFileLink   = newMockDoer(fixture.SharedFile.ListJSON)
	doerIssueSharedFileList   = newMockDoer(fixture.SharedFile.ListJSON)
	doerIssueSharedFileUnlink = newMockDoer(fixture.SharedFile.SingleJSON)
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

func ExampleIssueService_Participants() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueParticipants),
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
		backlog.WithDoer(doerIssueAttachmentList),
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
		backlog.WithDoer(doerIssueAttachmentDownload),
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
		backlog.WithDoer(doerIssueAttachmentRemove),
	)

	attachment, _ := c.Issue.Attachment.Remove(context.Background(), "TEST-1", 8)
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 8, Name: IMG0088.png
}

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

func ExampleIssueSharedFileService_Link() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueSharedFileLink),
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
		backlog.WithDoer(doerIssueSharedFileList),
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
		backlog.WithDoer(doerIssueSharedFileUnlink),
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
