package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	doerUserStarCount = newMockDoer(fixture.Star.CountJSON)
	doerUserStarList  = newMockDoer(fixture.Star.ListJSON)
)

func ExampleUserStarService_Count() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserStarCount),
	)

	count, _ := c.User.Star.Count(context.Background(), 1)
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 42
}

func ExampleUserStarService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserStarList),
	)

	stars, _ := c.User.Star.List(context.Background(), 1)
	fmt.Printf("ID: %d, Title: %s\n", stars[0].ID, stars[0].Title)
	// Output:
	// ID: 10, Title: [TEST-1] first issue
}
