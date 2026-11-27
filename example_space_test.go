package backlog_test

import (
	"context"
	"fmt"
	"strings"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleSpaceService_Info() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	space, err := c.Space.Info(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("SpaceKey: %s, Name: %s\n", space.SpaceKey, space.Name)
	// Output:
	// SpaceKey: nulab, Name: Nulab Inc.
}

func ExampleSpaceService_DiskUsage() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	diskUsage, err := c.Space.DiskUsage(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Capacity: %d, Issue: %d\n", diskUsage.Capacity, diskUsage.Issue)
	// Output:
	// Capacity: 1073741824, Issue: 119511
}

func ExampleSpaceService_Notification() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	notification, err := c.Space.Notification(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Content: %s\n", notification.Content)
	// Output:
	// Content: Backlog is a project management tool.
}

func ExampleSpaceService_UpdateNotification() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	notification, err := c.Space.UpdateNotification(context.Background(), "Backlog is a project management tool.")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Content: %s\n", notification.Content)
	// Output:
	// Content: Backlog is a project management tool.
}

func ExampleSpaceService_Icon() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	icon, err := c.Space.Icon(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer icon.Body.Close()

	fmt.Printf("ContentType: %s, FileName: %s\n", icon.ContentType, icon.Filename)
	// Output:
	// ContentType: image/png, FileName: space.png
}

func ExampleSpaceActivityService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	activities, err := c.Space.Activity.List(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleSpaceActivityService_One() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	activity, err := c.Space.Activity.One(context.Background(), 3153)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Type: %d\n", activity.ID, activity.Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleSpaceAttachmentService_Upload() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	attachment, err := c.Space.Attachment.Upload(context.Background(), "test.txt", strings.NewReader("hello"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", attachment.ID, attachment.Name)
	// Output:
	// ID: 1, Name: test.txt
}
