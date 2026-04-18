package backlog_test

import (
	"context"
	"fmt"
	"strings"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// SpaceActivityService
	doerSpaceActivityList = newMockDoer(fixture.Activity.ListJSON)

	// SpaceAttachmentService
	doerSpaceAttachmentUpload = newMockDoer(fixture.Attachment.UploadJSON)
)

func ExampleSpaceActivityService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSpaceActivityList),
	)

	activities, _ := c.Space.Activity.List(context.Background())
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleSpaceAttachmentService_Upload() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSpaceAttachmentUpload),
	)

	attachment, _ := c.Space.Attachment.Upload(context.Background(), "test.txt", strings.NewReader("hello"))
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 1, Name: test.txt
}
