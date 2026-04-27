package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// RepositoryService
	doerRepositoryAll = newMockDoer(fixture.Repository.ListJSON)
	doerRepositoryOne = newMockDoer(fixture.Repository.SingleJSON)
)

func ExampleRepositoryService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRepositoryAll),
	)

	repos, _ := c.Repository.All(context.Background(), "TEST")
	fmt.Printf("Count: %d, ID: %d, Name: %s\n", len(repos), repos[0].ID, repos[0].Name)
	// Output:
	// Count: 2, ID: 5, Name: foo
}

func ExampleRepositoryService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRepositoryOne),
	)

	repo, _ := c.Repository.One(context.Background(), "TEST", "foo")
	fmt.Printf("ID: %d, Name: %s, Description: %s\n", repo.ID, repo.Name, repo.Description)
	// Output:
	// ID: 5, Name: foo, Description: test repo
}
