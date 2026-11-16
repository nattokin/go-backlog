package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleRepositoryService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	repos, err := c.Repository.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Count: %d, ID: %d, Name: %s\n", len(repos), repos[0].ID, repos[0].Name)
	// Output:
	// Count: 2, ID: 5, Name: foo
}

func ExampleRepositoryService_One() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	repo, err := c.Repository.One(context.Background(), "TEST", "foo")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s, Description: %s\n", repo.ID, repo.Name, repo.Description)
	// Output:
	// ID: 5, Name: foo, Description: test repo
}
