package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExamplePriorityService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	priorities, err := c.Priority.List(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, p := range priorities {
		fmt.Printf("ID: %d, Name: %s\n", p.ID, p.Name)
	}
	// Output:
	// ID: 2, Name: High
	// ID: 3, Name: Normal
	// ID: 4, Name: Low
}
