package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleRecentlyViewedService_ListIssues() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	issues, err := c.RecentlyViewed.ListIssues(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("IssueKey: %s\n", issues[0].Issue.IssueKey)
	// Output:
	// IssueKey: TEST-1
}

func ExampleRecentlyViewedService_AddIssue() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	issue, err := c.RecentlyViewed.AddIssue(context.Background(), "TEST-1")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("IssueKey: %s\n", issue.IssueKey)
	// Output:
	// IssueKey: TEST-1
}

func ExampleRecentlyViewedService_ListProjects() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	projects, err := c.RecentlyViewed.ListProjects(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ProjectKey: %s\n", projects[0].Project.ProjectKey)
	// Output:
	// ProjectKey: TEST
}

func ExampleRecentlyViewedService_ListWikis() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	wikis, err := c.RecentlyViewed.ListWikis(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Name: %s\n", wikis[0].Page.Name)
	// Output:
	// Name: Home
}

func ExampleRecentlyViewedService_AddWiki() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	wiki, err := c.RecentlyViewed.AddWiki(context.Background(), 10)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Name: %s\n", wiki.Name)
	// Output:
	// Name: Home
}
