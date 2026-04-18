package backlog_test

import (
	"context"
	"fmt"
	"strings"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// SpaceAttachmentService
	doerSpaceAttachmentUpload = newMockDoer(fixture.Attachment.UploadJSON)
)

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
