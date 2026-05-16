package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleWikiService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiList),
	)

	wikis, _ := c.Wiki.List(context.Background(), "MYPROJECT")
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
		backlog.WithDoer(doerWikiSingle),
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
		backlog.WithDoer(doerWikiSingle),
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
		backlog.WithDoer(doerWikiSingle),
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
		backlog.WithDoer(doerWikiSingle),
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
		backlog.WithDoer(doerAttachmentList),
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
		backlog.WithDoer(doerAttachmentList),
	)

	attachments, _ := c.Wiki.Attachment.List(context.Background(), 34)
	fmt.Printf("ID: %d, Name: %s\n", attachments[0].ID, attachments[0].Name)
	// Output:
	// ID: 2, Name: A.png
}

func ExampleWikiAttachmentService_Download() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiAttachmentDownload),
	)

	file, _ := c.Wiki.Attachment.Download(context.Background(), 34, 2)
	fmt.Printf("ContentType: %s, FileName: %s\n", file.ContentType, file.Filename)
	// Output:
	// ContentType: image/png, FileName: A.png
}

func ExampleWikiAttachmentService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerAttachmentSingle),
	)

	attachment, _ := c.Wiki.Attachment.Remove(context.Background(), 34, 8)
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 8, Name: IMG0088.png
}

func ExampleWikiHistoryService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiHistoryList),
	)

	entries, _ := c.Wiki.History.List(context.Background(), 34)
	fmt.Printf("PageID: %d, Version: %d, Name: %s\n", entries[0].PageID, entries[0].Version, entries[0].Name)
	// Output:
	// PageID: 34, Version: 2, Name: Home
}

func ExampleWikiSharedFileService_Link() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSharedFileList),
	)

	files, _ := c.Wiki.SharedFile.Link(context.Background(), 34, []int{454403, 454404})
	fmt.Printf("ID: %d, Name: %s\n", files[0].ID, files[0].Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}

func ExampleWikiSharedFileService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSharedFileList),
	)

	files, _ := c.Wiki.SharedFile.List(context.Background(), 34)
	fmt.Printf("ID: %d, Name: %s\n", files[0].ID, files[0].Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}

func ExampleWikiSharedFileService_Unlink() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSharedFileSingle),
	)

	file, _ := c.Wiki.SharedFile.Unlink(context.Background(), 34, 454403)
	fmt.Printf("ID: %d, Name: %s\n", file.ID, file.Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}

func ExampleWikiStarService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerStarList),
	)

	stars, _ := c.Wiki.Star.List(context.Background(), 34)
	fmt.Printf("ID: %d, Title: %s\n", stars[0].ID, stars[0].Title)
	// Output:
	// ID: 10, Title: [TEST-1] first issue
}

func ExampleWikiStarService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerNoContent),
	)

	err := c.Wiki.Star.Add(context.Background(), 34)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}

func ExampleWikiStarService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerNoContent),
	)

	err := c.Wiki.Star.Remove(context.Background(), 42)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}
