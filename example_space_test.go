package backlog_test

import (
	"context"
	"fmt"
	"strings"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleSpaceService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSpaceSpace),
	)

	space, _ := c.Space.One(context.Background())
	fmt.Printf("SpaceKey: %s, Name: %s\n", space.SpaceKey, space.Name)
	// Output:
	// SpaceKey: nulab, Name: Nulab Inc.
}

func ExampleSpaceService_DiskUsage() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSpaceDiskUsage),
	)

	diskUsage, _ := c.Space.DiskUsage(context.Background())
	fmt.Printf("Capacity: %d, Issue: %d\n", diskUsage.Capacity, diskUsage.Issue)
	// Output:
	// Capacity: 1073741824, Issue: 119511
}

func ExampleSpaceService_Notification() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSpaceNotification),
	)

	notification, _ := c.Space.Notification(context.Background())
	fmt.Printf("Content: %s\n", notification.Content)
	// Output:
	// Content: Backlog is a project management tool.
}

func ExampleSpaceService_UpdateNotification() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerSpaceNotification),
	)

	notification, _ := c.Space.UpdateNotification(context.Background(), "Backlog is a project management tool.")
	fmt.Printf("Content: %s\n", notification.Content)
	// Output:
	// Content: Backlog is a project management tool.
}

func ExampleSpaceActivityService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerActivityList),
	)

	activities, _ := c.Space.Activity.List(context.Background())
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleSpaceActivityService_Get() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerActivitySingle),
	)

	activity, _ := c.Space.Activity.Get(context.Background(), 3153)
	fmt.Printf("ID: %d, Type: %d\n", activity.ID, activity.Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleSpaceAttachmentService_Upload() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerAttachmentUpload),
	)

	attachment, _ := c.Space.Attachment.Upload(context.Background(), "test.txt", strings.NewReader("hello"))
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 1, Name: test.txt
}
