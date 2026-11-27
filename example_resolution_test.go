package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleResolutionService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	resolutions, err := c.Resolution.List(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, r := range resolutions {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
	// Output:
	// ID: 0, Name: Fixed
	// ID: 1, Name: Won't Fix
	// ID: 2, Name: Invalid
	// ID: 3, Name: Duplication
	// ID: 4, Name: Cannot Reproduce
}
