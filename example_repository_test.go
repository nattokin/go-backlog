package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleRepositoryService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRepositoryList),
	)

	repos, _ := c.Repository.List(context.Background(), "TEST")
	fmt.Printf("Count: %d, ID: %d, Name: %s\n", len(repos), repos[0].ID, repos[0].Name)
	// Output:
	// Count: 2, ID: 5, Name: foo
}

func ExampleRepositoryService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRepositorySingle),
	)

	repo, _ := c.Repository.One(context.Background(), "TEST", "foo")
	fmt.Printf("ID: %d, Name: %s, Description: %s\n", repo.ID, repo.Name, repo.Description)
	// Output:
	// ID: 5, Name: foo, Description: test repo
}
