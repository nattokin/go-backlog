package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	doerWikiAll    = newMockDoer(fixture.Wiki.ListJSON)
	doerWikiCount  = newMockDoer(`{"count": 5}`)
	doerWikiOne    = newMockDoer(fixture.Wiki.MinimumJSON)
	doerWikiCreate = newMockDoer(fixture.Wiki.MinimumJSON)
	doerWikiUpdate = newMockDoer(fixture.Wiki.MinimumJSON)
	doerWikiDelete = newMockDoer(fixture.Wiki.MinimumJSON)

	doerWikiAttachmentAttach = newMockDoer(fixture.Attachment.ListJSON)
	doerWikiAttachmentList   = newMockDoer(fixture.Attachment.ListJSON)
	doerWikiAttachmentRemove = newMockDoer(fixture.Attachment.SingleJSON)
)

func ExampleWikiService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiAll),
	)

	wikis, _ := c.Wiki.All(context.Background(), "MYPROJECT")
	fmt.Printf("ID: %d, Name: %s\n", wikis[0].ID, wikis[0].Name)
	// Output:
	// ID: 112, Name: test1
}

func ExampleWikiService_Count() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiCount),
	)

	count, _ := c.Wiki.Count(context.Background(), "MYPROJECT")
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 5
}

func ExampleWikiService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiOne),
	)

	wiki, _ := c.Wiki.One(context.Background(), 34)
	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
	// Output:
	// ID: 34, Name: Minimum Wiki Page
}

func ExampleWikiService_Create() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiCreate),
	)

	wiki, _ := c.Wiki.Create(context.Background(), 56, "Minimum Wiki Page", "This is a minimal wiki page.")
	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
	// Output:
	// ID: 34, Name: Minimum Wiki Page
}

func ExampleWikiService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiUpdate),
	)

	wiki, _ := c.Wiki.Update(
		context.Background(),
		34,
		c.Wiki.Option.WithName("Minimum Wiki Page"),
	)
	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
	// Output:
	// ID: 34, Name: Minimum Wiki Page
}

func ExampleWikiService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiDelete),
	)

	wiki, _ := c.Wiki.Delete(context.Background(), 34)
	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
	// Output:
	// ID: 34, Name: Minimum Wiki Page
}

func ExampleWikiAttachmentService_Attach() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiAttachmentAttach),
	)

	attachments, _ := c.Wiki.Attachment.Attach(context.Background(), 34, []int{2, 5})
	fmt.Printf("ID: %d, Name: %s\n", attachments[0].ID, attachments[0].Name)
	// Output:
	// ID: 2, Name: A.png
}

func ExampleWikiAttachmentService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiAttachmentList),
	)

	attachments, _ := c.Wiki.Attachment.List(context.Background(), 34)
	fmt.Printf("ID: %d, Name: %s\n", attachments[0].ID, attachments[0].Name)
	// Output:
	// ID: 2, Name: A.png
}

func ExampleWikiAttachmentService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiAttachmentRemove),
	)

	attachment, _ := c.Wiki.Attachment.Remove(context.Background(), 34, 8)
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 8, Name: IMG0088.png
}
