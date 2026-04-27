package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	doerIssueSharedFileLink   = newMockDoer(fixture.SharedFile.ListJSON)
	doerIssueSharedFileList   = newMockDoer(fixture.SharedFile.ListJSON)
	doerIssueSharedFileUnlink = newMockDoer(fixture.SharedFile.SingleJSON)
)

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
