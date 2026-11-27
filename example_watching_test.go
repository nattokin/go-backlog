package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleWatchingService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	watchings, err := c.Watching.List(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, w := range watchings {
		fmt.Printf("ID: %d, IssueKey: %s, AlreadyRead: %t\n", w.ID, w.Issue.IssueKey, w.AlreadyRead)
	}
	// Output:
	// ID: 1, IssueKey: PRJ-1, AlreadyRead: true
	// ID: 2, IssueKey: PRJ-1, AlreadyRead: false
}

func ExampleWatchingService_Count() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	count, err := c.Watching.Count(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Count:", count)
	// Output:
	// Count: 138
}

func ExampleWatchingService_One() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	w, err := c.Watching.One(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Note: %s\n", w.ID, w.Note)
	// Output:
	// ID: 1, Note: Check before release
}

func ExampleWatchingService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	w, err := c.Watching.Create(
		context.Background(),
		"PRJ-1",
		c.Watching.Option.WithNote("Check before release"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, IssueKey: %s\n", w.ID, w.Issue.IssueKey)
	// Output:
	// ID: 1, IssueKey: PRJ-1
}

func ExampleWatchingService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	w, err := c.Watching.Update(context.Background(), 1, "Check before release")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Note: %s\n", w.ID, w.Note)
	// Output:
	// ID: 1, Note: Check before release
}

func ExampleWatchingService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	w, err := c.Watching.Delete(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d\n", w.ID)
	// Output:
	// ID: 1
}

func ExampleWatchingService_MarkAsRead() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = c.Watching.MarkAsRead(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}
