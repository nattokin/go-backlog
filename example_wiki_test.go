package backlog_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

// fixedDoer is a Doer that always returns a fixed response body with HTTP 200.
type fixedDoer struct {
	body string
}

func (d *fixedDoer) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(d.body)),
	}, nil
}

var (
	doerWikiAll    = &fixedDoer{body: fixture.Wiki.ListJSON}
	doerWikiCount  = &fixedDoer{body: `{"count": 5}`}
	doerWikiOne    = &fixedDoer{body: fixture.Wiki.MinimumJSON}
	doerWikiCreate = &fixedDoer{body: fixture.Wiki.MinimumJSON}
	doerWikiUpdate = &fixedDoer{body: fixture.Wiki.MinimumJSON}
	doerWikiDelete = &fixedDoer{body: fixture.Wiki.MinimumJSON}

	doerWikiAttachmentAttach = &fixedDoer{body: fixture.Attachment.ListJSON}
	doerWikiAttachmentList   = &fixedDoer{body: fixture.Attachment.ListJSON}
	doerWikiAttachmentRemove = &fixedDoer{body: fixture.Attachment.SingleJSON}
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
