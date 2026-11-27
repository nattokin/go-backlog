package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleNotificationService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	notifications, err := c.Notification.List(
		context.Background(),
		c.Notification.Option.WithCount(20),
		c.Notification.Option.WithOrder(backlog.OrderDesc),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, n := range notifications {
		fmt.Printf("ID: %d, Issue: %s, Sender: %s\n", n.ID, n.Issue.IssueKey, n.Sender.Name)
	}
	// Output:
	// ID: 299, Issue: TEST2-17, Sender: eguchi
}

func ExampleNotificationService_Count() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	count, err := c.Notification.Count(
		context.Background(),
		c.Notification.Option.WithResourceAlreadyRead(false),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Count:", count)
	// Output:
	// Count: 138
}

func ExampleNotificationService_ResetUnreadCount() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	count, err := c.Notification.ResetUnreadCount(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Count:", count)
	// Output:
	// Count: 4
}

func ExampleNotificationService_MarkAsRead() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = c.Notification.MarkAsRead(context.Background(), 299)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}
