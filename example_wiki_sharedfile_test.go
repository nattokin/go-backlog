package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	doerWikiSharedFileLink   = newMockDoer(fixture.SharedFile.ListJSON)
	doerWikiSharedFileList   = newMockDoer(fixture.SharedFile.ListJSON)
	doerWikiSharedFileUnlink = newMockDoer(fixture.SharedFile.SingleJSON)
)

func ExampleWikiSharedFileService_Link() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiSharedFileLink),
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
		backlog.WithDoer(doerWikiSharedFileList),
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
		backlog.WithDoer(doerWikiSharedFileUnlink),
	)

	file, _ := c.Wiki.SharedFile.Unlink(context.Background(), 34, 454403)
	fmt.Printf("ID: %d, Name: %s\n", file.ID, file.Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}
